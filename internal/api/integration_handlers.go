package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/ozsari/velour/internal/services"
)

// ── Unified Download Item ──

type DownloadItem struct {
	Name     string  `json:"name"`
	Size     int64   `json:"size"`
	Progress float64 `json:"progress"` // 0.0 - 1.0
	DlSpeed  int64   `json:"dlspeed"`
	UpSpeed  int64   `json:"upspeed"`
	State    string  `json:"state"` // downloading, seeding, paused, queued, completed, error, checking
	Eta      int64   `json:"eta"`
	Client   string  `json:"client"` // qbittorrent, transmission, deluge, sabnzbd, nzbget
	AddedOn  int64   `json:"added_on"`
	Seeds    int     `json:"seeds"`
	Peers    int     `json:"peers"`
}

type DownloadsResponse struct {
	Items   []DownloadItem `json:"items"`
	Clients []string       `json:"clients"` // which clients contributed data
}

func (s *Server) handleDownloads(w http.ResponseWriter, r *http.Request) {
	var (
		mu      sync.Mutex
		wg      sync.WaitGroup
		items   []DownloadItem
		clients []string
	)

	// List of fetchers: each checks if service is installed and fetches
	type fetcher struct {
		id    string
		fetch func(port int) []DownloadItem
	}

	fetchers := []fetcher{
		{"qbittorrent", fetchQbitDownloads},
		{"transmission", fetchTransmissionDownloads},
		{"deluge", fetchDelugeDownloads},
		{"sabnzbd", func(port int) []DownloadItem { return fetchSabnzbdDownloads(port, s.cfg.DataDir) }},
		{"nzbget", fetchNzbgetDownloads},
	}

	for _, f := range fetchers {
		port := getServicePort(f.id)
		if port == 0 {
			continue
		}
		wg.Add(1)
		go func(id string, port int, fn func(int) []DownloadItem) {
			defer wg.Done()
			result := fn(port)
			if result == nil {
				return
			}
			mu.Lock()
			items = append(items, result...)
			clients = append(clients, id)
			mu.Unlock()
		}(f.id, port, f.fetch)
	}

	wg.Wait()

	if items == nil {
		items = []DownloadItem{}
	}
	if clients == nil {
		clients = []string{}
	}

	jsonResponse(w, http.StatusOK, DownloadsResponse{Items: items, Clients: clients})
}

// ── qBittorrent ──

type qbitTorrentRaw struct {
	Name        string  `json:"name"`
	Size        int64   `json:"size"`
	Progress    float64 `json:"progress"`
	DlSpeed     int64   `json:"dlspeed"`
	UpSpeed     int64   `json:"upspeed"`
	State       string  `json:"state"`
	Eta         int64   `json:"eta"`
	AddedOn     int64   `json:"added_on"`
	NumSeeds    int     `json:"num_seeds"`
	NumLeechers int     `json:"num_leechs"`
}

func fetchQbitDownloads(port int) []DownloadItem {
	url := fmt.Sprintf("http://localhost:%d/api/v2/torrents/info?filter=all&limit=50&sort=added_on&reverse=true", port)
	data, err := proxyGet(url)
	if err != nil {
		return nil
	}
	var raw []qbitTorrentRaw
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil
	}
	items := make([]DownloadItem, 0, len(raw))
	for _, t := range raw {
		items = append(items, DownloadItem{
			Name:     t.Name,
			Size:     t.Size,
			Progress: t.Progress,
			DlSpeed:  t.DlSpeed,
			UpSpeed:  t.UpSpeed,
			State:    normalizeQbitState(t.State),
			Eta:      t.Eta,
			Client:   "qbittorrent",
			AddedOn:  t.AddedOn,
			Seeds:    t.NumSeeds,
			Peers:    t.NumLeechers,
		})
	}
	return items
}

func normalizeQbitState(state string) string {
	switch state {
	case "downloading", "forcedDL", "metaDL":
		return "downloading"
	case "uploading", "forcedUP", "stalledUP":
		return "seeding"
	case "pausedDL", "pausedUP":
		return "paused"
	case "queuedDL", "queuedUP":
		return "queued"
	case "stalledDL":
		return "stalled"
	case "checkingDL", "checkingUP", "checkingResumeData":
		return "checking"
	case "missingFiles", "error":
		return "error"
	case "moving":
		return "moving"
	default:
		return state
	}
}

// ── Transmission ──

type transmissionRequest struct {
	Method    string                 `json:"method"`
	Arguments map[string]interface{} `json:"arguments"`
}

type transmissionResponse struct {
	Result    string `json:"result"`
	Arguments struct {
		Torrents []transmissionTorrent `json:"torrents"`
	} `json:"arguments"`
}

type transmissionTorrent struct {
	Name            string  `json:"name"`
	TotalSize       int64   `json:"totalSize"`
	PercentDone     float64 `json:"percentDone"`
	RateDownload    int64   `json:"rateDownload"`
	RateUpload      int64   `json:"rateUpload"`
	Status          int     `json:"status"`
	Eta             int64   `json:"eta"`
	AddedDate       int64   `json:"addedDate"`
	PeersSendingToUs int    `json:"peersSendingToUs"`
	PeersGettingFromUs int  `json:"peersGettingFromUs"`
}

func fetchTransmissionDownloads(port int) []DownloadItem {
	url := fmt.Sprintf("http://localhost:%d/transmission/rpc", port)
	reqBody := transmissionRequest{
		Method: "torrent-get",
		Arguments: map[string]interface{}{
			"fields": []string{
				"name", "totalSize", "percentDone", "rateDownload", "rateUpload",
				"status", "eta", "addedDate", "peersSendingToUs", "peersGettingFromUs",
			},
		},
	}
	body, _ := json.Marshal(reqBody)

	client := &http.Client{Timeout: 5 * time.Second}

	// Transmission requires X-Transmission-Session-Id; get it from 409 response
	resp, err := client.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return nil
	}
	resp.Body.Close()

	sessionID := ""
	if resp.StatusCode == 409 {
		sessionID = resp.Header.Get("X-Transmission-Session-Id")
		if sessionID == "" {
			return nil
		}
	}

	req, _ := http.NewRequest("POST", url, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	if sessionID != "" {
		req.Header.Set("X-Transmission-Session-Id", sessionID)
	}

	resp, err = client.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil || resp.StatusCode != 200 {
		return nil
	}

	var result transmissionResponse
	if err := json.Unmarshal(respBody, &result); err != nil || result.Result != "success" {
		return nil
	}

	items := make([]DownloadItem, 0, len(result.Arguments.Torrents))
	for _, t := range result.Arguments.Torrents {
		items = append(items, DownloadItem{
			Name:     t.Name,
			Size:     t.TotalSize,
			Progress: t.PercentDone,
			DlSpeed:  t.RateDownload,
			UpSpeed:  t.RateUpload,
			State:    normalizeTransmissionState(t.Status),
			Eta:      t.Eta,
			Client:   "transmission",
			AddedOn:  t.AddedDate,
			Seeds:    t.PeersSendingToUs,
			Peers:    t.PeersGettingFromUs,
		})
	}
	return items
}

// Transmission status codes: 0=stopped, 1=check-wait, 2=check, 3=dl-wait, 4=downloading, 5=seed-wait, 6=seeding
func normalizeTransmissionState(status int) string {
	switch status {
	case 0:
		return "paused"
	case 1, 2:
		return "checking"
	case 3:
		return "queued"
	case 4:
		return "downloading"
	case 5:
		return "queued"
	case 6:
		return "seeding"
	default:
		return "unknown"
	}
}

// ── Deluge ──

type delugeRPCRequest struct {
	Method string        `json:"method"`
	Params []interface{} `json:"params"`
	ID     int           `json:"id"`
}

type delugeRPCResponse struct {
	Result interface{} `json:"result"`
	Error  interface{} `json:"error"`
	ID     int         `json:"id"`
}

func fetchDelugeDownloads(port int) []DownloadItem {
	url := fmt.Sprintf("http://localhost:%d/json", port)
	client := &http.Client{Timeout: 5 * time.Second}

	// Authenticate (default password)
	authReq := delugeRPCRequest{Method: "auth.login", Params: []interface{}{"deluge"}, ID: 1}
	authBody, _ := json.Marshal(authReq)
	resp, err := client.Post(url, "application/json", bytes.NewReader(authBody))
	if err != nil {
		return nil
	}
	// Save cookies for session
	cookies := resp.Cookies()
	resp.Body.Close()

	// Get torrents
	torrentReq := delugeRPCRequest{
		Method: "web.update_ui",
		Params: []interface{}{
			[]string{"name", "total_size", "progress", "download_payload_rate", "upload_payload_rate",
				"state", "eta", "time_added", "num_seeds", "num_peers"},
			map[string]interface{}{},
		},
		ID: 2,
	}
	torrentBody, _ := json.Marshal(torrentReq)
	req, _ := http.NewRequest("POST", url, bytes.NewReader(torrentBody))
	req.Header.Set("Content-Type", "application/json")
	for _, c := range cookies {
		req.AddCookie(c)
	}

	resp, err = client.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil || resp.StatusCode != 200 {
		return nil
	}

	var result delugeRPCResponse
	if err := json.Unmarshal(respBody, &result); err != nil || result.Error != nil {
		return nil
	}

	// Parse the nested result structure
	resultMap, ok := result.Result.(map[string]interface{})
	if !ok {
		return nil
	}
	torrentsMap, ok := resultMap["torrents"].(map[string]interface{})
	if !ok {
		return nil
	}

	items := make([]DownloadItem, 0)
	for _, v := range torrentsMap {
		t, ok := v.(map[string]interface{})
		if !ok {
			continue
		}
		items = append(items, DownloadItem{
			Name:     jsonStr(t, "name"),
			Size:     jsonInt64(t, "total_size"),
			Progress: jsonFloat(t, "progress") / 100.0,
			DlSpeed:  jsonInt64(t, "download_payload_rate"),
			UpSpeed:  jsonInt64(t, "upload_payload_rate"),
			State:    normalizeDelugeState(jsonStr(t, "state")),
			Eta:      jsonInt64(t, "eta"),
			Client:   "deluge",
			AddedOn:  jsonInt64(t, "time_added"),
			Seeds:    int(jsonInt64(t, "num_seeds")),
			Peers:    int(jsonInt64(t, "num_peers")),
		})
	}
	return items
}

func normalizeDelugeState(state string) string {
	switch strings.ToLower(state) {
	case "downloading":
		return "downloading"
	case "seeding":
		return "seeding"
	case "paused":
		return "paused"
	case "queued":
		return "queued"
	case "checking":
		return "checking"
	case "error":
		return "error"
	default:
		return state
	}
}

// ── SABnzbd ──

type sabnzbdQueue struct {
	Queue struct {
		Slots []sabnzbdSlot `json:"slots"`
	} `json:"queue"`
}

type sabnzbdSlot struct {
	Filename  string `json:"filename"`
	MB        string `json:"mb"`
	MBLeft    string `json:"mbleft"`
	Percentage string `json:"percentage"`
	Status    string `json:"status"`
	TimeLeft  string `json:"timeleft"`
}

type sabnzbdHistory struct {
	History struct {
		Slots []sabnzbdHistSlot `json:"slots"`
	} `json:"history"`
}

type sabnzbdHistSlot struct {
	Name       string  `json:"name"`
	Bytes      int64   `json:"bytes"`
	Status     string  `json:"status"`
	CompletedTime int64 `json:"completed"`
}

func fetchSabnzbdDownloads(port int, dataDir string) []DownloadItem {
	apiKey := getSabnzbdAPIKey(dataDir)

	// Fetch queue
	queueURL := fmt.Sprintf("http://localhost:%d/api?mode=queue&output=json", port)
	if apiKey != "" {
		queueURL += "&apikey=" + apiKey
	}

	data, err := proxyGet(queueURL)
	if err != nil {
		return nil
	}

	var queue sabnzbdQueue
	if err := json.Unmarshal(data, &queue); err != nil {
		return nil
	}

	items := make([]DownloadItem, 0)
	for _, slot := range queue.Queue.Slots {
		pct := parseFloat(slot.Percentage) / 100.0
		totalMB := parseFloat(slot.MB)
		items = append(items, DownloadItem{
			Name:     slot.Filename,
			Size:     int64(totalMB * 1024 * 1024),
			Progress: pct,
			State:    normalizeSabnzbdState(slot.Status),
			Client:   "sabnzbd",
		})
	}

	// Fetch recent history (last 10)
	histURL := fmt.Sprintf("http://localhost:%d/api?mode=history&limit=10&output=json", port)
	if apiKey != "" {
		histURL += "&apikey=" + apiKey
	}
	data, err = proxyGet(histURL)
	if err == nil {
		var hist sabnzbdHistory
		if json.Unmarshal(data, &hist) == nil {
			for _, slot := range hist.History.Slots {
				items = append(items, DownloadItem{
					Name:     slot.Name,
					Size:     slot.Bytes,
					Progress: 1.0,
					State:    normalizeSabnzbdHistState(slot.Status),
					Client:   "sabnzbd",
					AddedOn:  slot.CompletedTime,
				})
			}
		}
	}

	return items
}

func normalizeSabnzbdState(status string) string {
	switch strings.ToLower(status) {
	case "downloading":
		return "downloading"
	case "paused":
		return "paused"
	case "queued", "idle":
		return "queued"
	default:
		return "downloading"
	}
}

func normalizeSabnzbdHistState(status string) string {
	switch strings.ToLower(status) {
	case "completed":
		return "completed"
	case "failed":
		return "error"
	default:
		return "completed"
	}
}

func getSabnzbdAPIKey(dataDir string) string {
	paths := []string{
		fmt.Sprintf("%s/sabnzbd/config/sabnzbd.ini", dataDir),
		fmt.Sprintf("%s/data/sabnzbd/sabnzbd.ini", dataDir),
	}
	for _, path := range paths {
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		for _, line := range strings.Split(string(data), "\n") {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "api_key") {
				parts := strings.SplitN(line, "=", 2)
				if len(parts) == 2 {
					return strings.TrimSpace(parts[1])
				}
			}
		}
	}
	return ""
}

// ── NZBGet ──

type nzbgetResult struct {
	Result []nzbgetGroup `json:"result"`
}

type nzbgetGroup struct {
	NZBID          int    `json:"NZBID"`
	NZBName        string `json:"NZBName"`
	FileSizeLo     int64  `json:"FileSizeLo"`
	FileSizeHi     int64  `json:"FileSizeHi"`
	RemainingSizeLo int64 `json:"RemainingSizeLo"`
	RemainingSizeHi int64 `json:"RemainingSizeHi"`
	DownloadRate   int64  `json:"DownloadRate"`
	Status         string `json:"Status"`
}

func fetchNzbgetDownloads(port int) []DownloadItem {
	// NZBGet uses JSON-RPC, default auth: nzbget:tegbzn6789
	url := fmt.Sprintf("http://nzbget:tegbzn6789@localhost:%d/jsonrpc/listgroups", port)
	data, err := proxyGet(url)
	if err != nil {
		return nil
	}

	var result nzbgetResult
	if err := json.Unmarshal(data, &result); err != nil {
		return nil
	}

	items := make([]DownloadItem, 0, len(result.Result))
	for _, g := range result.Result {
		totalSize := g.FileSizeHi*4294967296 + g.FileSizeLo
		remainSize := g.RemainingSizeHi*4294967296 + g.RemainingSizeLo
		progress := 0.0
		if totalSize > 0 {
			progress = float64(totalSize-remainSize) / float64(totalSize)
		}
		items = append(items, DownloadItem{
			Name:     g.NZBName,
			Size:     totalSize,
			Progress: progress,
			DlSpeed:  g.DownloadRate,
			State:    normalizeNzbgetState(g.Status),
			Client:   "nzbget",
		})
	}
	return items
}

func normalizeNzbgetState(status string) string {
	s := strings.ToLower(status)
	if strings.Contains(s, "download") {
		return "downloading"
	}
	if strings.Contains(s, "pause") {
		return "paused"
	}
	if strings.Contains(s, "queue") {
		return "queued"
	}
	if strings.Contains(s, "post") || strings.Contains(s, "unpack") {
		return "processing"
	}
	return s
}

// ── Legacy qBittorrent endpoints (kept for backward compat) ──

type QbitTorrent struct {
	Name        string  `json:"name"`
	Size        int64   `json:"size"`
	Progress    float64 `json:"progress"`
	DlSpeed     int64   `json:"dlspeed"`
	UpSpeed     int64   `json:"upspeed"`
	State       string  `json:"state"`
	Eta         int64   `json:"eta"`
	Category    string  `json:"category"`
	AddedOn     int64   `json:"added_on"`
	NumSeeds    int     `json:"num_seeds"`
	NumLeechers int     `json:"num_leechs"`
}

type QbitTransferInfo struct {
	DlSpeed       int64 `json:"dl_info_speed"`
	UpSpeed       int64 `json:"up_info_speed"`
	DlTotal       int64 `json:"dl_info_data"`
	UpTotal       int64 `json:"up_info_data"`
	DlSessionData int64 `json:"dl_session_data"`
	UpSessionData int64 `json:"up_session_data"`
}

func (s *Server) handleQbitTorrents(w http.ResponseWriter, r *http.Request) {
	port := getServicePort("qbittorrent")
	if port == 0 {
		jsonError(w, http.StatusNotFound, "qBittorrent not in catalog")
		return
	}

	url := fmt.Sprintf("http://localhost:%d/api/v2/torrents/info?filter=all&limit=50&sort=added_on&reverse=true", port)
	data, err := proxyGet(url)
	if err != nil {
		jsonResponse(w, http.StatusOK, []QbitTorrent{})
		return
	}

	var torrents []QbitTorrent
	if err := json.Unmarshal(data, &torrents); err != nil {
		jsonResponse(w, http.StatusOK, []QbitTorrent{})
		return
	}
	jsonResponse(w, http.StatusOK, torrents)
}

func (s *Server) handleQbitTransfer(w http.ResponseWriter, r *http.Request) {
	port := getServicePort("qbittorrent")
	if port == 0 {
		jsonError(w, http.StatusNotFound, "qBittorrent not in catalog")
		return
	}

	url := fmt.Sprintf("http://localhost:%d/api/v2/transfer/info", port)
	data, err := proxyGet(url)
	if err != nil {
		jsonResponse(w, http.StatusOK, QbitTransferInfo{})
		return
	}

	var info QbitTransferInfo
	if err := json.Unmarshal(data, &info); err != nil {
		jsonResponse(w, http.StatusOK, QbitTransferInfo{})
		return
	}
	jsonResponse(w, http.StatusOK, info)
}

// ── Sonarr Integration ──

type SonarrCalendarEntry struct {
	SeriesID      int           `json:"seriesId"`
	EpisodeFileID int           `json:"episodeFileId"`
	SeasonNumber  int           `json:"seasonNumber"`
	EpisodeNumber int           `json:"episodeNumber"`
	Title         string        `json:"title"`
	AirDateUTC    string        `json:"airDateUtc"`
	HasFile       bool          `json:"hasFile"`
	Monitored     bool          `json:"monitored"`
	Series        *SonarrSeries `json:"series,omitempty"`
}

type SonarrSeries struct {
	Title  string        `json:"title"`
	Images []SonarrImage `json:"images,omitempty"`
}

type SonarrImage struct {
	CoverType string `json:"coverType"`
	RemoteURL string `json:"remoteUrl"`
}

func (s *Server) handleSonarrCalendar(w http.ResponseWriter, r *http.Request) {
	port := getServicePort("sonarr")
	if port == 0 {
		jsonResponse(w, http.StatusOK, []SonarrCalendarEntry{})
		return
	}

	apiKey := r.URL.Query().Get("apikey")
	if apiKey == "" {
		apiKey = getSonarrAPIKey(s.cfg.DataDir)
	}

	now := time.Now()
	start := now.AddDate(0, 0, -1).Format("2006-01-02")
	end := now.AddDate(0, 0, 14).Format("2006-01-02")

	url := fmt.Sprintf("http://localhost:%d/api/v3/calendar?start=%s&end=%s&includeSeries=true", port, start, end)

	client := &http.Client{Timeout: 5 * time.Second}
	req, _ := http.NewRequest("GET", url, nil)
	if apiKey != "" {
		req.Header.Set("X-Api-Key", apiKey)
	}

	resp, err := client.Do(req)
	if err != nil {
		jsonResponse(w, http.StatusOK, []SonarrCalendarEntry{})
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil || resp.StatusCode != 200 {
		jsonResponse(w, http.StatusOK, []SonarrCalendarEntry{})
		return
	}

	var entries []SonarrCalendarEntry
	if err := json.Unmarshal(body, &entries); err != nil {
		jsonResponse(w, http.StatusOK, []SonarrCalendarEntry{})
		return
	}
	jsonResponse(w, http.StatusOK, entries)
}

// ── Radarr Integration ──

type RadarrCalendarEntry struct {
	Title           string        `json:"title"`
	Year            int           `json:"year"`
	PhysicalRelease string        `json:"physicalRelease,omitempty"`
	DigitalRelease  string        `json:"digitalRelease,omitempty"`
	InCinemas       string        `json:"inCinemas,omitempty"`
	HasFile         bool          `json:"hasFile"`
	Monitored       bool          `json:"monitored"`
	Images          []RadarrImage `json:"images,omitempty"`
}

type RadarrImage struct {
	CoverType string `json:"coverType"`
	RemoteURL string `json:"remoteUrl"`
}

func (s *Server) handleRadarrCalendar(w http.ResponseWriter, r *http.Request) {
	port := getServicePort("radarr")
	if port == 0 {
		jsonResponse(w, http.StatusOK, []RadarrCalendarEntry{})
		return
	}

	apiKey := r.URL.Query().Get("apikey")
	if apiKey == "" {
		apiKey = getRadarrAPIKey(s.cfg.DataDir)
	}

	now := time.Now()
	start := now.AddDate(0, 0, -1).Format("2006-01-02")
	end := now.AddDate(0, 0, 30).Format("2006-01-02")

	url := fmt.Sprintf("http://localhost:%d/api/v3/calendar?start=%s&end=%s", port, start, end)

	client := &http.Client{Timeout: 5 * time.Second}
	req, _ := http.NewRequest("GET", url, nil)
	if apiKey != "" {
		req.Header.Set("X-Api-Key", apiKey)
	}

	resp, err := client.Do(req)
	if err != nil {
		jsonResponse(w, http.StatusOK, []RadarrCalendarEntry{})
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil || resp.StatusCode != 200 {
		jsonResponse(w, http.StatusOK, []RadarrCalendarEntry{})
		return
	}

	var entries []RadarrCalendarEntry
	if err := json.Unmarshal(body, &entries); err != nil {
		jsonResponse(w, http.StatusOK, []RadarrCalendarEntry{})
		return
	}
	jsonResponse(w, http.StatusOK, entries)
}

// ── Helpers ──

func getServicePort(id string) int {
	def := services.FindByID(id)
	if def == nil || len(def.Ports) == 0 {
		return 0
	}
	return def.Ports[0].Host
}

func proxyGet(url string) ([]byte, error) {
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func getSonarrAPIKey(dataDir string) string {
	return readAPIKeyFromXML(dataDir, "sonarr")
}

func getRadarrAPIKey(dataDir string) string {
	return readAPIKeyFromXML(dataDir, "radarr")
}

func readAPIKeyFromXML(dataDir, service string) string {
	paths := []string{
		fmt.Sprintf("%s/%s/config/config.xml", dataDir, service),
		fmt.Sprintf("%s/data/%s/config.xml", dataDir, service),
	}

	for _, path := range paths {
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		content := string(data)
		start := strings.Index(content, "<ApiKey>")
		if start == -1 {
			continue
		}
		start += len("<ApiKey>")
		end := strings.Index(content[start:], "</ApiKey>")
		if end == -1 {
			continue
		}
		return content[start : start+end]
	}
	return ""
}

// JSON helper functions for parsing dynamic maps
func jsonStr(m map[string]interface{}, key string) string {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func jsonInt64(m map[string]interface{}, key string) int64 {
	if v, ok := m[key]; ok {
		switch n := v.(type) {
		case float64:
			return int64(n)
		case int64:
			return n
		}
	}
	return 0
}

func jsonFloat(m map[string]interface{}, key string) float64 {
	if v, ok := m[key]; ok {
		if f, ok := v.(float64); ok {
			return f
		}
	}
	return 0
}

func parseFloat(s string) float64 {
	var f float64
	fmt.Sscanf(s, "%f", &f)
	return f
}

// Suppress unused import warning
var _ = log.Printf
