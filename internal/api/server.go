package api

import (
	"database/sql"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/ozsari/velour/internal/auth"
	"github.com/ozsari/velour/internal/automation"
	"github.com/ozsari/velour/internal/config"
	"github.com/ozsari/velour/internal/monitor"
	"github.com/ozsari/velour/internal/services"
	"github.com/ozsari/velour/web"

	_ "github.com/mattn/go-sqlite3"
)

type Server struct {
	cfg        *config.Config
	auth       *auth.AuthService
	monitor    *monitor.Monitor
	netTracker *monitor.NetworkTracker
	docker     *services.DockerManager
	native     *services.NativeManager
	automation *automation.Orchestrator
	db         *sql.DB
}

func NewServer(cfg *config.Config) *Server {
	return &Server{
		cfg:     cfg,
		monitor: monitor.New(),
	}
}

func (s *Server) Start() error {
	// Init database
	db, err := sql.Open("sqlite3", s.cfg.DBPath)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	s.db = db

	// Init auth
	s.auth = auth.New(db, s.cfg.JWTSecret)
	if err := s.auth.InitDB(); err != nil {
		return fmt.Errorf("failed to init auth: %w", err)
	}

	// Create unified data directory structure for hardlink support
	for _, dir := range []string{
		"data/downloads", "data/media/tv", "data/media/movies",
		"data/media/music", "data/media/books", "data/media/comics",
		"data/media/audiobooks", "data/usenet",
	} {
		os.MkdirAll(filepath.Join(s.cfg.DataDir, dir), 0755)
	}

	// Init Docker manager
	docker, err := services.NewDockerManager(s.cfg.DataDir)
	if err != nil {
		log.Printf("Warning: Docker not available: %v", err)
	} else {
		s.docker = docker
	}

	// Init Native manager (Linux only)
	s.native = services.NewNativeManager(s.cfg.DataDir)

	// Init Network tracker (monthly stats persistence)
	s.netTracker = monitor.NewNetworkTracker(db)
	if err := s.netTracker.InitDB(); err != nil {
		return fmt.Errorf("failed to init network tracker: %w", err)
	}
	s.netTracker.Start()

	// Init Automation orchestrator
	s.automation = automation.New(db, s.docker, s.native, s.cfg.InstallMode)
	if err := s.automation.InitDB(); err != nil {
		return fmt.Errorf("failed to init automation: %w", err)
	}
	s.automation.Start()

	// Setup routes
	mux := http.NewServeMux()
	s.registerRoutes(mux)

	addr := fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.Port)
	log.Printf("Starting server on %s", addr)
	return http.ListenAndServe(addr, corsMiddleware(mux))
}

func (s *Server) registerRoutes(mux *http.ServeMux) {
	// Public routes
	mux.HandleFunc("GET /api/health", s.handleHealth)
	mux.HandleFunc("GET /api/config", s.handleConfig)
	mux.HandleFunc("GET /api/setup/status", s.handleSetupStatus)
	mux.HandleFunc("POST /api/setup", s.handleSetup)
	mux.HandleFunc("POST /api/auth/login", s.handleLogin)

	// Protected routes
	mux.HandleFunc("GET /api/system", s.authMiddleware(s.handleSystemInfo))
	mux.HandleFunc("GET /api/services", s.authMiddleware(s.handleListServices))
	mux.HandleFunc("GET /api/services/catalog", s.authMiddleware(s.handleServiceCatalog))
	mux.HandleFunc("POST /api/services/{id}/install", s.authMiddleware(s.handleInstallService))
	mux.HandleFunc("POST /api/services/{id}/start", s.authMiddleware(s.handleStartService))
	mux.HandleFunc("POST /api/services/{id}/stop", s.authMiddleware(s.handleStopService))
	mux.HandleFunc("POST /api/services/{id}/restart", s.authMiddleware(s.handleRestartService))
	mux.HandleFunc("DELETE /api/services/{id}", s.authMiddleware(s.handleRemoveService))

	// Network stats
	mux.HandleFunc("GET /api/network/months", s.authMiddleware(s.handleNetworkMonths))
	mux.HandleFunc("GET /api/network/month/{month}", s.authMiddleware(s.handleNetworkMonth))

	// Integrations (service API proxies)
	mux.HandleFunc("GET /api/integrations/downloads", s.authMiddleware(s.handleDownloads))
	mux.HandleFunc("GET /api/integrations/qbit/torrents", s.authMiddleware(s.handleQbitTorrents))
	mux.HandleFunc("GET /api/integrations/qbit/transfer", s.authMiddleware(s.handleQbitTransfer))
	mux.HandleFunc("GET /api/integrations/sonarr/calendar", s.authMiddleware(s.handleSonarrCalendar))
	mux.HandleFunc("GET /api/integrations/radarr/calendar", s.authMiddleware(s.handleRadarrCalendar))

	// Automation rules
	mux.HandleFunc("GET /api/automation/rules", s.authMiddleware(s.handleListRules))
	mux.HandleFunc("POST /api/automation/rules", s.authMiddleware(s.handleCreateRule))
	mux.HandleFunc("GET /api/automation/rules/{id}", s.authMiddleware(s.handleGetRule))
	mux.HandleFunc("PUT /api/automation/rules/{id}", s.authMiddleware(s.handleUpdateRule))
	mux.HandleFunc("DELETE /api/automation/rules/{id}", s.authMiddleware(s.handleDeleteRule))
	mux.HandleFunc("POST /api/automation/rules/{id}/toggle", s.authMiddleware(s.handleToggleRule))
	mux.HandleFunc("GET /api/automation/templates", s.authMiddleware(s.handleListTemplates))

	// Serve embedded frontend
	distFS, err := fs.Sub(web.DistFS, "dist")
	if err != nil {
		log.Fatalf("Failed to load embedded frontend: %v", err)
	}
	mux.Handle("/", http.FileServer(http.FS(distFS)))
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
