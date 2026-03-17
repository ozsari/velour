import { useEffect, useState } from 'react';
import {
  Package, ExternalLink, RotateCcw, Square, Play,
  ArrowUp, ArrowDown, Download,
  Clock, ChevronLeft, ChevronRight, Calendar,
  Tv, Film, CheckCircle, Database,
} from 'lucide-react';
import { useSystemInfo, formatBytes, formatSpeed } from '../hooks/useSystemInfo';
import {
  api, type Service, type MonthlyNetStats,
  type DownloadItem, type SonarrCalendarEntry, type RadarrCalendarEntry,
} from '../lib/api';
import { APP_LOGOS } from '../lib/logos';

// ── Demo Data ──

const DEMO_DOWNLOADS: DownloadItem[] = [
  { name: 'Ubuntu.24.04.LTS.Desktop.amd64.iso', size: 5100000000, progress: 0.73, dlspeed: 12500000, upspeed: 2100000, state: 'downloading', eta: 280, client: 'qbittorrent', added_on: Date.now() / 1000 - 3600, seeds: 142, peers: 38 },
  { name: 'Big.Buck.Bunny.2008.1080p.BluRay.x264', size: 8200000000, progress: 0.45, dlspeed: 8700000, upspeed: 950000, state: 'downloading', eta: 620, client: 'qbittorrent', added_on: Date.now() / 1000 - 7200, seeds: 87, peers: 21 },
  { name: 'Arch.Linux.2024.03.01.x86_64.iso', size: 920000000, progress: 1.0, dlspeed: 0, upspeed: 350000, state: 'seeding', eta: 0, client: 'transmission', added_on: Date.now() / 1000 - 86400, seeds: 0, peers: 5 },
  { name: 'Sintel.2010.4K.Remaster.mkv', size: 14500000000, progress: 0.12, dlspeed: 5200000, upspeed: 0, state: 'downloading', eta: 2400, client: 'deluge', added_on: Date.now() / 1000 - 1800, seeds: 34, peers: 12 },
  { name: 'Tears.of.Steel.2012.1080p.mkv', size: 3400000000, progress: 0.0, dlspeed: 0, upspeed: 0, state: 'paused', eta: 0, client: 'qbittorrent', added_on: Date.now() / 1000 - 172800, seeds: 0, peers: 0 },
  { name: 'Linux.Mint.21.3.Cinnamon.iso', size: 2700000000, progress: 0.91, dlspeed: 15000000, upspeed: 1800000, state: 'downloading', eta: 45, client: 'transmission', added_on: Date.now() / 1000 - 600, seeds: 256, peers: 14 },
  { name: 'Cosmos.Laundromat.First.Cycle.4K', size: 6100000000, progress: 1.0, dlspeed: 0, upspeed: 120000, state: 'seeding', eta: 0, client: 'qbittorrent', added_on: Date.now() / 1000 - 259200, seeds: 0, peers: 3 },
  { name: 'Latest.NZB.Package.Complete', size: 4800000000, progress: 0.67, dlspeed: 22000000, upspeed: 0, state: 'downloading', eta: 180, client: 'sabnzbd', added_on: Date.now() / 1000 - 900, seeds: 0, peers: 0 },
];

const DEMO_SONARR: SonarrCalendarEntry[] = [
  { seriesId: 1, seasonNumber: 5, episodeNumber: 3, title: 'The One Where Everything Changes', airDateUtc: new Date().toISOString(), hasFile: false, monitored: true, series: { title: 'The Last of Us', images: [{ coverType: 'poster', remoteUrl: 'https://image.tmdb.org/t/p/w200/uKvVjHNqB5VmOrdxqAt2F7J78ED.jpg' }] } },
  { seriesId: 2, seasonNumber: 2, episodeNumber: 7, title: 'Aftermath', airDateUtc: new Date(Date.now() + 86400000).toISOString(), hasFile: false, monitored: true, series: { title: 'Severance', images: [{ coverType: 'poster', remoteUrl: 'https://image.tmdb.org/t/p/w200/pPHpeI2X1qEd1CS1SeyrdhZ4qnT.jpg' }] } },
  { seriesId: 3, seasonNumber: 3, episodeNumber: 1, title: 'New Beginnings', airDateUtc: new Date(Date.now() + 86400000 * 2).toISOString(), hasFile: false, monitored: true, series: { title: 'The Bear', images: [{ coverType: 'poster', remoteUrl: 'https://image.tmdb.org/t/p/w200/eKfVzzEazSIjJMrw9ADa2x8ksLz.jpg' }] } },
  { seriesId: 4, seasonNumber: 1, episodeNumber: 12, title: 'Endgame', airDateUtc: new Date(Date.now() + 86400000 * 3).toISOString(), hasFile: true, monitored: true, series: { title: 'Shogun', images: [{ coverType: 'poster', remoteUrl: 'https://image.tmdb.org/t/p/w200/7O4iVfOMQmdCSxhOg1WnzG1AgYT.jpg' }] } },
  { seriesId: 5, seasonNumber: 7, episodeNumber: 4, title: 'The Reckoning', airDateUtc: new Date(Date.now() + 86400000 * 5).toISOString(), hasFile: false, monitored: true, series: { title: 'Better Call Saul', images: [{ coverType: 'poster', remoteUrl: 'https://image.tmdb.org/t/p/w200/fC2HDm5t0kHl7mTm7jxMR31b7by.jpg' }] } },
];

const DEMO_RADARR: RadarrCalendarEntry[] = [
  { title: 'Dune: Part Two', year: 2024, digitalRelease: new Date(Date.now() + 86400000 * 7).toISOString(), hasFile: false, monitored: true, images: [{ coverType: 'poster', remoteUrl: 'https://image.tmdb.org/t/p/w200/tihf8Trht9zP3scmUQfvGlAY9FU.jpg' }] },
  { title: 'The Batman', year: 2022, inCinemas: new Date(Date.now() + 86400000 * 14).toISOString(), hasFile: false, monitored: true, images: [{ coverType: 'poster', remoteUrl: 'https://image.tmdb.org/t/p/w200/3B19RaHTSpyNEcjWtsAPqgJsQ0h.jpg' }] },
  { title: 'Oppenheimer', year: 2023, digitalRelease: new Date(Date.now() - 86400000).toISOString(), hasFile: true, monitored: true, images: [{ coverType: 'poster', remoteUrl: 'https://image.tmdb.org/t/p/w200/8Gxv8gSFCU0XGDykEGv7zR1n2ua.jpg' }] },
];

// ── Helpers ──

function MonthLabel(month: string): string {
  const [y, m] = month.split('-');
  const date = new Date(parseInt(y), parseInt(m) - 1);
  return date.toLocaleDateString('en-US', { month: 'long', year: 'numeric' });
}

function getCurrentMonth(): string {
  const now = new Date();
  return `${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, '0')}`;
}

function shiftMonth(month: string, delta: number): string {
  const [y, m] = month.split('-').map(Number);
  const d = new Date(y, m - 1 + delta);
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}`;
}

function formatEta(seconds: number): string {
  if (seconds <= 0 || seconds >= 8640000) return '∞';
  const h = Math.floor(seconds / 3600);
  const m = Math.floor((seconds % 3600) / 60);
  if (h > 0) return `${h}h ${m}m`;
  const s = seconds % 60;
  if (m > 0) return `${m}m ${s}s`;
  return `${s}s`;
}

function stateInfo(state: string): { label: string; color: string } {
  const map: Record<string, { label: string; color: string }> = {
    downloading: { label: 'Downloading', color: 'var(--accent-green)' },
    seeding: { label: 'Seeding', color: 'var(--accent-blue-light)' },
    paused: { label: 'Paused', color: 'var(--text-muted)' },
    queued: { label: 'Queued', color: 'var(--text-muted)' },
    stalled: { label: 'Stalled', color: 'var(--accent-yellow)' },
    checking: { label: 'Checking', color: 'var(--accent-yellow)' },
    completed: { label: 'Completed', color: 'var(--accent-blue-light)' },
    processing: { label: 'Processing', color: 'var(--accent-cyan)' },
    moving: { label: 'Moving', color: 'var(--accent-cyan)' },
    error: { label: 'Error', color: 'var(--accent-red)' },
  };
  return map[state] || { label: state, color: 'var(--text-muted)' };
}

const CLIENT_LABELS: Record<string, { name: string; color: string }> = {
  qbittorrent: { name: 'qBit', color: '#4499ee' },
  transmission: { name: 'Trans', color: '#cc3333' },
  deluge: { name: 'Deluge', color: '#4488dd' },
  sabnzbd: { name: 'SAB', color: '#e8a83c' },
  nzbget: { name: 'NZBGet', color: '#44bb44' },
};

function getPosterUrl(images?: { coverType: string; remoteUrl: string }[]): string | null {
  if (!images) return null;
  const poster = images.find(img => img.coverType === 'poster');
  if (poster?.remoteUrl) return poster.remoteUrl;
  // fallback to any image
  return images[0]?.remoteUrl || null;
}

function relativeDay(dateStr: string): string {
  const d = new Date(dateStr);
  const now = new Date();
  const today = new Date(now.getFullYear(), now.getMonth(), now.getDate());
  const target = new Date(d.getFullYear(), d.getMonth(), d.getDate());
  const diff = Math.round((target.getTime() - today.getTime()) / 86400000);
  if (diff === 0) return 'Today';
  if (diff === 1) return 'Tomorrow';
  if (diff === -1) return 'Yesterday';
  if (diff > 0 && diff <= 7) return d.toLocaleDateString('en-US', { weekday: 'long' });
  return d.toLocaleDateString('en-US', { month: 'short', day: 'numeric' });
}

// ── Widget wrapper ──

function Widget({ title, icon, children, headerRight }: {
  title: string; icon: React.ReactNode; children: React.ReactNode; headerRight?: React.ReactNode;
}) {
  return (
    <div style={{
      background: 'var(--bg-card)', border: '1px solid var(--border-primary)',
      borderRadius: 14, padding: 20, display: 'flex', flexDirection: 'column',
    }}>
      <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', marginBottom: 16 }}>
        <div style={{ display: 'flex', alignItems: 'center', gap: 10 }}>
          {icon}
          <span style={{ fontSize: 14, fontWeight: 600, color: 'var(--text-primary)' }}>{title}</span>
        </div>
        {headerRight}
      </div>
      {children}
    </div>
  );
}

// ── Main Dashboard ──

export default function Dashboard() {
  const system = useSystemInfo();
  const [services, setServices] = useState<Service[]>([]);
  const [selectedMonth, setSelectedMonth] = useState(getCurrentMonth);
  const [monthStats, setMonthStats] = useState<MonthlyNetStats | null>(null);
  const [availableMonths, setAvailableMonths] = useState<string[]>([]);
  const [downloads, setDownloads] = useState<DownloadItem[]>([]);
  const [dlClients, setDlClients] = useState<string[]>([]);
  const [sonarrCal, setSonarrCal] = useState<SonarrCalendarEntry[]>([]);
  const [radarrCal, setRadarrCal] = useState<RadarrCalendarEntry[]>([]);
  const [demoMode, setDemoMode] = useState(false);

  const fetchServices = () => {
    api.listServices().then((s) => setServices(s || [])).catch(() => {});
  };
  const fetchMonthStats = () => {
    api.networkMonth(selectedMonth).then(setMonthStats).catch(() => {});
  };
  const fetchDownloads = () => {
    api.downloads().then(r => {
      setDownloads(r.items || []);
      setDlClients(r.clients || []);
    }).catch(() => { setDownloads([]); setDlClients([]); });
  };
  const fetchCalendar = () => {
    api.sonarrCalendar().then(setSonarrCal).catch(() => setSonarrCal([]));
    api.radarrCalendar().then(setRadarrCal).catch(() => setRadarrCal([]));
  };

  useEffect(() => {
    fetchServices();
    fetchDownloads();
    fetchCalendar();
    api.networkMonths().then(m => setAvailableMonths((m || []).map(s => s.month))).catch(() => {});

    const svcIv = setInterval(fetchServices, 5000);
    const dlIv = setInterval(fetchDownloads, 3000);
    const calIv = setInterval(fetchCalendar, 60000);
    return () => { clearInterval(svcIv); clearInterval(dlIv); clearInterval(calIv); };
  }, []);

  useEffect(() => {
    fetchMonthStats();
    const iv = setInterval(fetchMonthStats, 60000);
    return () => clearInterval(iv);
  }, [selectedMonth]);

  if (!system) {
    return (
      <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'center', height: 256 }}>
        <div style={{
          width: 32, height: 32, border: '3px solid rgba(59,130,246,0.2)',
          borderTopColor: '#3b82f6', borderRadius: '50%', animation: 'spin 0.8s linear infinite',
        }} />
      </div>
    );
  }

  const handleStart = async (id: string) => {
    try { await api.startService(id); setTimeout(fetchServices, 1500); } catch {}
  };
  const handleRestart = async (id: string) => {
    try { await api.restartService(id); setTimeout(fetchServices, 1500); } catch {}
  };
  const handleStop = async (id: string) => {
    try { await api.stopService(id); setTimeout(fetchServices, 1500); } catch {}
  };

  // Use demo data if enabled or use real data
  const displayDownloads = demoMode ? DEMO_DOWNLOADS : downloads;
  const displaySonarr = demoMode ? DEMO_SONARR : sonarrCal;
  const displayRadarr = demoMode ? DEMO_RADARR : radarrCal;
  const displayClients = demoMode ? ['qbittorrent', 'transmission', 'deluge', 'sabnzbd'] : dlClients;

  const runningCount = services.filter(s => s.status === 'running').length;
  const activeDL = displayDownloads.filter(t => t.state === 'downloading');

  // Detect which download clients are installed
  const dlClientIDs = ['qbittorrent', 'transmission', 'deluge', 'sabnzbd', 'nzbget', 'flood', 'rutorrent', 'pyload', 'jdownloader2'];
  const hasAnyDLClient = demoMode || services.some(s => dlClientIDs.includes(s.id)) || displayClients.length > 0;

  const sonarrInstalled = demoMode || services.some(s => s.id === 'sonarr');
  const radarrInstalled = demoMode || services.some(s => s.id === 'radarr');
  const hasCalendar = sonarrInstalled || radarrInstalled;

  const hasWidgets = hasAnyDLClient || hasCalendar;

  return (
    <div>
      {/* Header */}
      <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', marginBottom: 20 }}>
        <div>
          <h2 style={{ fontSize: 22, fontWeight: 700, color: 'var(--text-primary)', marginBottom: 2 }}>Dashboard</h2>
          <p style={{ fontSize: 12, color: 'var(--text-muted)' }}>
            {system.hostname} &middot; {system.os} &middot; {system.kernel}
          </p>
        </div>
        <div style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
          {/* Demo mode toggle */}
          <button
            onClick={() => setDemoMode(d => !d)}
            title={demoMode ? 'Switch to live data' : 'Show demo data'}
            style={{
              display: 'flex', alignItems: 'center', gap: 5,
              padding: '5px 10px', borderRadius: 8, border: '1px solid var(--border-primary)',
              background: demoMode ? 'rgba(168,85,247,0.1)' : 'var(--bg-card)',
              color: demoMode ? 'var(--accent-purple)' : 'var(--text-muted)',
              fontSize: 11, fontWeight: 500, cursor: 'pointer', transition: 'all 0.15s',
            }}
          >
            <Database size={12} />
            {demoMode ? 'Demo' : 'Live'}
          </button>

          <div style={{
            display: 'flex', alignItems: 'center', gap: 8,
            padding: '5px 12px', borderRadius: 8,
            background: 'var(--bg-card)', border: '1px solid var(--border-primary)',
          }}>
            <Clock size={13} style={{ color: 'var(--text-muted)' }} />
            <span style={{ fontSize: 12, color: 'var(--text-tertiary)', fontFamily: 'monospace' }}>
              {system.uptime_human}
            </span>
          </div>
        </div>
      </div>

      {/* ── Widgets grid ── */}
      {hasWidgets && (
        <div style={{ display: 'grid', gridTemplateColumns: 'repeat(2, 1fr)', gap: 12, marginBottom: 12 }}>

          {/* Downloads Widget (unified: all clients) */}
          {hasAnyDLClient && (
            <Widget
              title="Downloads"
              icon={<div style={{ width: 28, height: 28, borderRadius: 7, display: 'flex', alignItems: 'center', justifyContent: 'center', background: 'rgba(34,197,94,0.1)', color: 'var(--accent-green)' }}>
                <Download size={16} />
              </div>}
              headerRight={
                <div style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
                  {/* Client badges */}
                  {displayClients.map(c => {
                    const cl = CLIENT_LABELS[c];
                    return cl ? (
                      <span key={c} style={{
                        fontSize: 9, fontWeight: 600, padding: '2px 5px', borderRadius: 4,
                        background: `${cl.color}18`, color: cl.color, letterSpacing: 0.3,
                      }}>{cl.name}</span>
                    ) : null;
                  })}
                  <span style={{ fontSize: 11, color: activeDL.length > 0 ? 'var(--accent-green)' : 'var(--text-muted)', fontWeight: 500 }}>
                    {activeDL.length > 0 ? `${activeDL.length} active` : 'Idle'}
                  </span>
                </div>
              }
            >
              {displayDownloads.length === 0 ? (
                <div style={{ textAlign: 'center', padding: '24px 0', color: 'var(--text-muted)', fontSize: 13 }}>
                  No downloads
                </div>
              ) : (
                <div style={{ display: 'flex', flexDirection: 'column', gap: 8, maxHeight: 360, overflowY: 'auto' }}>
                  {displayDownloads.slice(0, 15).map((t, i) => {
                    const st = stateInfo(t.state);
                    const pct = Math.round(t.progress * 100);
                    const cl = CLIENT_LABELS[t.client];
                    return (
                      <div key={i} style={{
                        padding: '10px 12px', borderRadius: 8,
                        background: 'var(--bg-primary)', border: '1px solid var(--border-secondary)',
                      }}>
                        <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', marginBottom: 6 }}>
                          <div style={{ display: 'flex', alignItems: 'center', gap: 6, minWidth: 0, flex: 1 }}>
                            {cl && (
                              <span style={{
                                fontSize: 8, fontWeight: 700, padding: '1px 4px', borderRadius: 3,
                                background: `${cl.color}18`, color: cl.color, flexShrink: 0, letterSpacing: 0.3,
                              }}>{cl.name}</span>
                            )}
                            <span style={{
                              fontSize: 12, fontWeight: 500, color: 'var(--text-primary)',
                              overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap',
                            }}>
                              {t.name}
                            </span>
                          </div>
                          <span style={{ fontSize: 10, color: st.color, fontWeight: 600, flexShrink: 0, marginLeft: 8 }}>{st.label}</span>
                        </div>
                        {/* Progress bar */}
                        <div style={{ height: 3, background: 'var(--bg-tertiary)', borderRadius: 2, overflow: 'hidden', marginBottom: 6 }}>
                          <div style={{
                            height: '100%', borderRadius: 2, transition: 'width 0.5s ease',
                            background: pct >= 100 ? 'var(--accent-blue-light)' : 'var(--accent-green)',
                            width: `${pct}%`,
                          }} />
                        </div>
                        <div style={{ display: 'flex', alignItems: 'center', gap: 12, fontSize: 10, color: 'var(--text-muted)' }}>
                          <span>{pct}%</span>
                          <span>{formatBytes(t.size)}</span>
                          {t.dlspeed > 0 && <span style={{ color: 'var(--accent-green)' }}>↓ {formatSpeed(t.dlspeed)}</span>}
                          {t.upspeed > 0 && <span style={{ color: 'var(--accent-blue-light)' }}>↑ {formatSpeed(t.upspeed)}</span>}
                          {t.eta > 0 && t.eta < 8640000 && <span>ETA {formatEta(t.eta)}</span>}
                          {(t.seeds > 0 || t.peers > 0) && (
                            <span style={{ marginLeft: 'auto' }}>{t.seeds}S / {t.peers}P</span>
                          )}
                        </div>
                      </div>
                    );
                  })}
                </div>
              )}
            </Widget>
          )}

          {/* Calendar Widget (Sonarr + Radarr) */}
          {hasCalendar && (
            <Widget
              title="Calendar"
              icon={<div style={{ width: 28, height: 28, borderRadius: 7, display: 'flex', alignItems: 'center', justifyContent: 'center', background: 'rgba(168,85,247,0.1)', color: 'var(--accent-purple)' }}>
                <Calendar size={16} />
              </div>}
              headerRight={
                <span style={{ fontSize: 11, color: 'var(--text-muted)', fontWeight: 500 }}>
                  {displaySonarr.length + displayRadarr.length} upcoming
                </span>
              }
            >
              {displaySonarr.length === 0 && displayRadarr.length === 0 ? (
                <div style={{ textAlign: 'center', padding: '24px 0', color: 'var(--text-muted)', fontSize: 13 }}>
                  No upcoming episodes or movies
                </div>
              ) : (
                <div style={{ display: 'flex', flexDirection: 'column', gap: 6, maxHeight: 360, overflowY: 'auto' }}>
                  {/* Sonarr entries */}
                  {displaySonarr.map((ep, i) => {
                    const day = relativeDay(ep.airDateUtc);
                    const isToday = day === 'Today';
                    const poster = getPosterUrl(ep.series?.images);
                    return (
                      <div key={`s${i}`} style={{
                        display: 'flex', alignItems: 'center', gap: 10,
                        padding: '8px 12px', borderRadius: 8,
                        background: isToday ? 'rgba(168,85,247,0.05)' : 'var(--bg-primary)',
                        border: `1px solid ${isToday ? 'rgba(168,85,247,0.2)' : 'var(--border-secondary)'}`,
                      }}>
                        <div style={{
                          width: 36, height: 50, borderRadius: 6, flexShrink: 0, overflow: 'hidden',
                          background: 'rgba(96,165,250,0.1)', display: 'flex',
                          alignItems: 'center', justifyContent: 'center',
                        }}>
                          {poster ? (
                            <img src={poster} alt="" style={{ width: '100%', height: '100%', objectFit: 'cover' }}
                              onError={(e) => {
                                const el = e.target as HTMLImageElement;
                                el.style.display = 'none';
                                el.parentElement!.innerHTML = '<svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" style="color:var(--accent-blue-light)"><rect width="20" height="15" x="2" y="7" rx="2" ry="2"/><polyline points="17 2 12 7 7 2"/></svg>';
                              }}
                            />
                          ) : (
                            <Tv size={14} style={{ color: 'var(--accent-blue-light)' }} />
                          )}
                        </div>
                        <div style={{ flex: 1, minWidth: 0 }}>
                          <div style={{ display: 'flex', alignItems: 'center', gap: 6 }}>
                            <span style={{
                              fontSize: 12, fontWeight: 600, color: 'var(--text-primary)',
                              overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap',
                            }}>
                              {ep.series?.title || 'Unknown'}
                            </span>
                            {ep.hasFile && <CheckCircle size={11} style={{ color: 'var(--accent-green)', flexShrink: 0 }} />}
                          </div>
                          <span style={{ fontSize: 10, color: 'var(--text-muted)' }}>
                            S{String(ep.seasonNumber).padStart(2, '0')}E{String(ep.episodeNumber).padStart(2, '0')} - {ep.title}
                          </span>
                        </div>
                        <span style={{
                          fontSize: 10, fontWeight: 600, color: isToday ? 'var(--accent-purple)' : 'var(--text-muted)',
                          whiteSpace: 'nowrap',
                        }}>
                          {day}
                        </span>
                      </div>
                    );
                  })}
                  {/* Radarr entries */}
                  {displayRadarr.map((mov, i) => {
                    const dateStr = mov.digitalRelease || mov.physicalRelease || mov.inCinemas || '';
                    const day = dateStr ? relativeDay(dateStr) : '—';
                    const isToday = day === 'Today';
                    const poster = getPosterUrl(mov.images);
                    return (
                      <div key={`r${i}`} style={{
                        display: 'flex', alignItems: 'center', gap: 10,
                        padding: '8px 12px', borderRadius: 8,
                        background: isToday ? 'rgba(245,158,11,0.05)' : 'var(--bg-primary)',
                        border: `1px solid ${isToday ? 'rgba(245,158,11,0.2)' : 'var(--border-secondary)'}`,
                      }}>
                        <div style={{
                          width: 36, height: 50, borderRadius: 6, flexShrink: 0, overflow: 'hidden',
                          background: 'rgba(245,158,11,0.1)', display: 'flex',
                          alignItems: 'center', justifyContent: 'center',
                        }}>
                          {poster ? (
                            <img src={poster} alt="" style={{ width: '100%', height: '100%', objectFit: 'cover' }}
                              onError={(e) => {
                                const el = e.target as HTMLImageElement;
                                el.style.display = 'none';
                                el.parentElement!.innerHTML = '<svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" style="color:var(--accent-yellow)"><rect width="18" height="18" x="3" y="3" rx="2"/><circle cx="12" cy="12" r="3"/><path d="m16 2-4 4-4-4"/></svg>';
                              }}
                            />
                          ) : (
                            <Film size={14} style={{ color: 'var(--accent-yellow)' }} />
                          )}
                        </div>
                        <div style={{ flex: 1, minWidth: 0 }}>
                          <div style={{ display: 'flex', alignItems: 'center', gap: 6 }}>
                            <span style={{
                              fontSize: 12, fontWeight: 600, color: 'var(--text-primary)',
                              overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap',
                            }}>
                              {mov.title}
                            </span>
                            {mov.hasFile && <CheckCircle size={11} style={{ color: 'var(--accent-green)', flexShrink: 0 }} />}
                          </div>
                          <span style={{ fontSize: 10, color: 'var(--text-muted)' }}>
                            {mov.year}
                          </span>
                        </div>
                        <span style={{
                          fontSize: 10, fontWeight: 600, color: isToday ? 'var(--accent-yellow)' : 'var(--text-muted)',
                          whiteSpace: 'nowrap',
                        }}>
                          {day}
                        </span>
                      </div>
                    );
                  })}
                </div>
              )}
            </Widget>
          )}

          {/* Monthly Data */}
          <Widget
            title="Monthly Data"
            icon={<div style={{ width: 28, height: 28, borderRadius: 7, display: 'flex', alignItems: 'center', justifyContent: 'center', background: 'rgba(6,182,212,0.1)', color: 'var(--accent-cyan)' }}>
              <Calendar size={16} />
            </div>}
            headerRight={
              <div style={{ display: 'flex', alignItems: 'center', gap: 4 }}>
                <button onClick={() => setSelectedMonth(m => shiftMonth(m, -1))} style={{
                  display: 'flex', alignItems: 'center', justifyContent: 'center',
                  width: 22, height: 22, borderRadius: 5, border: 'none',
                  background: 'var(--bg-tertiary)', color: 'var(--text-tertiary)', cursor: 'pointer', padding: 0,
                }}><ChevronLeft size={12} /></button>
                <span style={{ fontSize: 11, fontWeight: 600, color: 'var(--text-secondary)', minWidth: 100, textAlign: 'center', fontFamily: 'monospace' }}>
                  {MonthLabel(selectedMonth)}
                </span>
                <button
                  onClick={() => setSelectedMonth(m => { const n = shiftMonth(m, 1); return n <= getCurrentMonth() ? n : m; })}
                  disabled={selectedMonth >= getCurrentMonth()}
                  style={{
                    display: 'flex', alignItems: 'center', justifyContent: 'center',
                    width: 22, height: 22, borderRadius: 5, border: 'none',
                    background: selectedMonth >= getCurrentMonth() ? 'transparent' : 'var(--bg-tertiary)',
                    color: selectedMonth >= getCurrentMonth() ? 'var(--text-faint)' : 'var(--text-tertiary)',
                    cursor: selectedMonth >= getCurrentMonth() ? 'default' : 'pointer', padding: 0,
                  }}
                ><ChevronRight size={12} /></button>
              </div>
            }
          >
            <div style={{ display: 'flex', gap: 20 }}>
              <div style={{ flex: 1 }}>
                <div style={{ display: 'flex', alignItems: 'center', gap: 6, marginBottom: 4 }}>
                  <ArrowDown size={13} style={{ color: 'var(--accent-cyan)' }} />
                  <span style={{ fontSize: 11, color: 'var(--text-muted)' }}>Downloaded</span>
                </div>
                <div style={{ fontSize: 20, fontWeight: 700, color: 'var(--text-primary)', fontFamily: 'monospace' }}>
                  {monthStats ? formatBytes(monthStats.bytes_recv) : '—'}
                </div>
              </div>
              <div style={{ width: 1, background: 'var(--border-secondary)' }} />
              <div style={{ flex: 1 }}>
                <div style={{ display: 'flex', alignItems: 'center', gap: 6, marginBottom: 4 }}>
                  <ArrowUp size={13} style={{ color: 'var(--accent-cyan)' }} />
                  <span style={{ fontSize: 11, color: 'var(--text-muted)' }}>Uploaded</span>
                </div>
                <div style={{ fontSize: 20, fontWeight: 700, color: 'var(--text-primary)', fontFamily: 'monospace' }}>
                  {monthStats ? formatBytes(monthStats.bytes_sent) : '—'}
                </div>
              </div>
            </div>
            {availableMonths.length > 1 && (
              <div style={{ display: 'flex', gap: 4, marginTop: 10, flexWrap: 'wrap' }}>
                {availableMonths.slice(0, 6).map(m => (
                  <button key={m} onClick={() => setSelectedMonth(m)} style={{
                    padding: '2px 6px', borderRadius: 4, border: 'none', fontSize: 9,
                    fontFamily: 'monospace', cursor: 'pointer',
                    background: m === selectedMonth ? 'rgba(6,182,212,0.15)' : 'var(--bg-tertiary)',
                    color: m === selectedMonth ? 'var(--accent-cyan)' : 'var(--text-muted)',
                    fontWeight: m === selectedMonth ? 600 : 400,
                  }}>{m.split('-')[1]}/{m.split('-')[0].slice(2)}</button>
                ))}
              </div>
            )}
          </Widget>
        </div>
      )}

      {/* ── Services Grid (Homarr-style app tiles) ── */}
      <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', marginBottom: 10 }}>
        <span style={{ fontSize: 14, fontWeight: 600, color: 'var(--text-primary)' }}>Services</span>
        {services.length > 0 && (
          <span style={{ fontSize: 11, color: 'var(--text-muted)' }}>
            <span style={{ display: 'inline-block', width: 5, height: 5, borderRadius: '50%', background: 'var(--accent-green)', marginRight: 4 }} />
            {runningCount}/{services.length} running
          </span>
        )}
      </div>

      {services.length > 0 ? (
        <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fill, minmax(140px, 1fr))', gap: 10 }}>
          {services.map((s) => {
            const logo = APP_LOGOS[s.id];
            const isRunning = s.status === 'running';
            return (
              <div key={s.id} style={{
                background: 'var(--bg-card)', border: '1px solid var(--border-primary)',
                borderRadius: 12, padding: '16px 12px', textAlign: 'center',
                opacity: isRunning ? 1 : 0.5, transition: 'opacity 0.2s, transform 0.15s',
                position: 'relative',
              }}>
                {/* Status dot */}
                <div style={{
                  position: 'absolute', top: 8, right: 8,
                  width: 6, height: 6, borderRadius: '50%',
                  background: isRunning ? 'var(--accent-green)' : 'var(--text-muted)',
                  boxShadow: isRunning ? '0 0 6px rgba(34,197,94,0.4)' : 'none',
                }} />

                {/* Logo */}
                <div style={{
                  width: 40, height: 40, borderRadius: 10, margin: '0 auto 8px',
                  background: 'var(--bg-primary)', display: 'flex', alignItems: 'center', justifyContent: 'center',
                  fontSize: 16, fontWeight: 700, color: 'var(--text-tertiary)', overflow: 'hidden',
                }}>
                  {logo ? (
                    <img src={logo} alt="" style={{ width: 32, height: 32, objectFit: 'contain' }}
                      onError={(e) => {
                        const el = e.target as HTMLImageElement;
                        el.style.display = 'none';
                        el.parentElement!.textContent = s.name.charAt(0);
                      }} />
                  ) : s.name.charAt(0)}
                </div>

                {/* Name */}
                <p style={{
                  fontSize: 12, fontWeight: 600, color: 'var(--text-primary)',
                  marginBottom: 2, overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap',
                }}>{s.name}</p>
                <p style={{ fontSize: 10, color: 'var(--text-muted)', marginBottom: 10 }}>:{s.port}</p>

                {/* Actions */}
                <div style={{ display: 'flex', gap: 4, justifyContent: 'center' }}>
                  {isRunning ? (
                    <>
                      {s.web_url && (
                        <a href={s.web_url.replace('localhost', window.location.hostname)} target="_blank" rel="noreferrer" title="Open" style={{
                          padding: '4px 8px', borderRadius: 5, display: 'flex', alignItems: 'center',
                          background: 'rgba(59,130,246,0.08)', color: 'var(--accent-blue-light)', textDecoration: 'none',
                        }}><ExternalLink size={11} /></a>
                      )}
                      <button onClick={() => handleRestart(s.id)} title="Restart" style={{
                        padding: '4px 6px', borderRadius: 5, border: 'none', cursor: 'pointer',
                        background: 'rgba(100,116,139,0.1)', color: 'var(--text-tertiary)', display: 'flex', alignItems: 'center',
                      }}><RotateCcw size={11} /></button>
                      <button onClick={() => handleStop(s.id)} title="Stop" style={{
                        padding: '4px 6px', borderRadius: 5, border: 'none', cursor: 'pointer',
                        background: 'rgba(239,68,68,0.08)', color: 'var(--accent-red)', display: 'flex', alignItems: 'center',
                      }}><Square size={11} /></button>
                    </>
                  ) : (
                    <button onClick={() => handleStart(s.id)} title="Start" style={{
                      padding: '4px 8px', borderRadius: 5, border: 'none', cursor: 'pointer',
                      background: 'rgba(34,197,94,0.08)', color: 'var(--accent-green)', display: 'flex', alignItems: 'center', gap: 4,
                      fontSize: 11,
                    }}><Play size={11} /> Start</button>
                  )}
                </div>
              </div>
            );
          })}
        </div>
      ) : (
        <div style={{
          textAlign: 'center', padding: '40px 20px',
          background: 'var(--bg-card)', border: '1px solid var(--border-secondary)', borderRadius: 14,
        }}>
          <Package size={22} style={{ margin: '0 auto 10px', color: 'var(--accent-blue-light)' }} />
          <p style={{ fontSize: 14, fontWeight: 600, color: 'var(--text-secondary)', marginBottom: 4 }}>No services installed</p>
          <p style={{ fontSize: 12, color: 'var(--text-muted)' }}>Go to App Catalog to install your first app</p>
        </div>
      )}
    </div>
  );
}
