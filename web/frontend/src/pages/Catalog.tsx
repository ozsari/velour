import { useEffect, useState } from 'react';
import { Download, Check, ExternalLink, Github } from 'lucide-react';
import { api, type Service } from '../lib/api';

interface AppMeta {
  name: string;
  logo: string;
  github: string;
  version: string;
  stars: string;
  subtitle: string;
  port: number;
  image: string;
  description: string;
  category: string;
}

const APPS: Record<string, AppMeta> = {
  // ── Download Automation ──
  autobrr: { name: 'Autobrr', logo: '/logos/autobrr.png', github: 'https://github.com/autobrr/autobrr', version: 'v1.74.0', stars: '2.5k', subtitle: 'Download Automation', port: 7474, image: 'ghcr.io/autobrr/autobrr', description: 'Modern automation tool for torrents and usenet. Monitors IRC, RSS feeds and more.', category: 'download' },
  recyclarr: { name: 'Recyclarr', logo: '/logos/recyclarr.png', github: 'https://github.com/recyclarr/recyclarr', version: 'v7.4.0', stars: '1.5k', subtitle: 'TRaSH Guides Sync', port: 0, image: 'ghcr.io/recyclarr/recyclarr', description: 'Automatically sync TRaSH Guides recommended settings to Sonarr and Radarr.', category: 'download' },

  // ── Download Clients ──
  deluge: { name: 'Deluge', logo: '/logos/deluge.png', github: 'https://github.com/deluge-torrent/deluge', version: 'v2.1.1', stars: '1.7k', subtitle: 'BitTorrent Client', port: 8112, image: 'lscr.io/linuxserver/deluge', description: 'Lightweight, free, cross-platform BitTorrent client with a web interface.', category: 'client' },
  flood: { name: 'Flood', logo: '/logos/flood.png', github: 'https://github.com/jesec/flood', version: 'v4.8.5', stars: '2.2k', subtitle: 'Torrent Web UI', port: 3001, image: 'jesec/flood', description: 'Modern web UI for rTorrent, qBittorrent, and Transmission.', category: 'client' },
  qbittorrent: { name: 'qBittorrent', logo: '/logos/qbittorrent.png', github: 'https://github.com/qbittorrent/qBittorrent', version: 'v5.1.0', stars: '30k', subtitle: 'BitTorrent Client', port: 8085, image: 'lscr.io/linuxserver/qbittorrent', description: 'Free, open-source BitTorrent client with a feature-rich web interface.', category: 'client' },
  qui: { name: 'qui', logo: '/logos/qui.png', github: 'https://github.com/autobrr/qui', version: 'latest', stars: '-', subtitle: 'qBittorrent Web UI', port: 7476, image: 'ghcr.io/autobrr/qui', description: 'Fast, modern web interface for qBittorrent with multi-instance support and automations.', category: 'client' },
  rutorrent: { name: 'ruTorrent', logo: '/logos/rutorrent.png', github: 'https://github.com/Novik/ruTorrent', version: 'v4.3.8', stars: '2.2k', subtitle: 'rTorrent + autodl-irssi', port: 8087, image: 'crazymax/rtorrent-rutorrent', description: 'Complete torrent solution: rTorrent + ruTorrent web UI + autodl-irssi IRC automation.', category: 'client' },
  transmission: { name: 'Transmission', logo: '/logos/transmission.png', github: 'https://github.com/transmission/transmission', version: 'v4.0.6', stars: '13k', subtitle: 'BitTorrent Client', port: 9091, image: 'lscr.io/linuxserver/transmission', description: 'Fast, easy and free BitTorrent client with minimal resource usage.', category: 'client' },
  nzbget: { name: 'NZBGet', logo: '/logos/nzbget.png', github: 'https://github.com/nzbgetcom/nzbget', version: 'v24.6', stars: '1.5k', subtitle: 'Usenet Downloader', port: 6789, image: 'lscr.io/linuxserver/nzbget', description: 'Efficient usenet downloader written in C++ for maximum performance.', category: 'client' },
  sabnzbd: { name: 'SABnzbd', logo: '/logos/sabnzbd.png', github: 'https://github.com/sabnzbd/sabnzbd', version: 'v4.4.1', stars: '2k', subtitle: 'Usenet Downloader', port: 8088, image: 'lscr.io/linuxserver/sabnzbd', description: 'Free, open-source usenet downloader with web-based interface.', category: 'client' },
  pyload: { name: 'pyLoad', logo: '/logos/pyload.png', github: 'https://github.com/pyload/pyload', version: 'v0.5.0', stars: '3.4k', subtitle: 'Download Manager', port: 8089, image: 'lscr.io/linuxserver/pyload-ng', description: 'Free, open-source download manager for HTTP, FTP and other protocols.', category: 'client' },
  jdownloader2: { name: 'JDownloader 2', logo: '/logos/jdownloader2.png', github: 'https://github.com/mirror/jdownloader', version: 'latest', stars: '-', subtitle: 'Download Manager', port: 5801, image: 'jlesage/jdownloader-2', description: 'Open-source download manager with web interface for direct downloads.', category: 'client' },

  // ── Media Management ──
  sonarr: { name: 'Sonarr', logo: '/logos/sonarr.png', github: 'https://github.com/Sonarr/Sonarr', version: 'v4.0.14', stars: '11.5k', subtitle: 'TV Series', port: 8989, image: 'lscr.io/linuxserver/sonarr', description: 'Smart PVR for newsgroup and bittorrent users. Monitors RSS feeds for new episodes.', category: 'media' },
  sonarr2: { name: 'Sonarr (2nd)', logo: '/logos/sonarr.png', github: 'https://github.com/Sonarr/Sonarr', version: 'v4.0.14', stars: '11.5k', subtitle: 'TV Series (2nd Instance)', port: 8990, image: 'lscr.io/linuxserver/sonarr', description: 'Second instance of Sonarr for managing a separate TV library.', category: 'media' },
  radarr: { name: 'Radarr', logo: '/logos/radarr.png', github: 'https://github.com/Radarr/Radarr', version: 'v5.20.2', stars: '11k', subtitle: 'Movies', port: 7878, image: 'lscr.io/linuxserver/radarr', description: 'Movie collection manager. Automatically searches, downloads and manages movies.', category: 'media' },
  radarr2: { name: 'Radarr (2nd)', logo: '/logos/radarr.png', github: 'https://github.com/Radarr/Radarr', version: 'v5.20.2', stars: '11k', subtitle: 'Movies (2nd Instance)', port: 7879, image: 'lscr.io/linuxserver/radarr', description: 'Second instance of Radarr for managing a separate movie library.', category: 'media' },
  bazarr: { name: 'Bazarr', logo: '/logos/bazarr.png', github: 'https://github.com/morpheus65535/bazarr', version: 'v1.5.1', stars: '3k', subtitle: 'Subtitles', port: 6767, image: 'lscr.io/linuxserver/bazarr', description: 'Companion to Sonarr and Radarr for managing and downloading subtitles.', category: 'media' },
  lidarr: { name: 'Lidarr', logo: '/logos/lidarr.png', github: 'https://github.com/Lidarr/Lidarr', version: 'v2.9.6', stars: '3.9k', subtitle: 'Music', port: 8686, image: 'lscr.io/linuxserver/lidarr', description: 'Music collection manager. Monitors RSS feeds for new albums.', category: 'media' },
  medusa: { name: 'Medusa', logo: '/logos/medusa.png', github: 'https://github.com/pymedusa/Medusa', version: 'v1.0.21', stars: '1.8k', subtitle: 'TV Shows', port: 8081, image: 'lscr.io/linuxserver/medusa', description: 'Automatic video library manager for TV shows with multi-source support.', category: 'media' },
  sickchill: { name: 'SickChill', logo: '/logos/sickchill.png', github: 'https://github.com/SickChill/SickChill', version: 'v2024.3.1', stars: '2.4k', subtitle: 'TV Shows', port: 8083, image: 'lscr.io/linuxserver/sickchill', description: 'Automatic video library manager for TV shows from various sources.', category: 'media' },
  sickgear: { name: 'SickGear', logo: '/logos/sickgear.png', github: 'https://github.com/SickGear/SickGear', version: 'v3.32.0', stars: '650', subtitle: 'TV Shows', port: 8082, image: 'lscr.io/linuxserver/sickgear', description: 'SickBeard fork with improved stability and features.', category: 'media' },
  mylar3: { name: 'Mylar3', logo: '/logos/mylar3.png', github: 'https://github.com/mylar3/mylar3', version: 'v0.7.8', stars: '1.1k', subtitle: 'Comics', port: 8090, image: 'lscr.io/linuxserver/mylar3', description: 'Automated comic book downloader for usenet and torrent users.', category: 'media' },
  filebot: { name: 'FileBot', logo: '/logos/filebot.png', github: 'https://github.com/filebot/filebot', version: 'v5.1.7', stars: '900', subtitle: 'File Organizer', port: 5800, image: 'jlesage/filebot', description: 'Ultimate tool for organizing and renaming movies, TV shows and music.', category: 'media' },
  lazylibrarian: { name: 'LazyLibrarian', logo: '/logos/lazylibrarian.png', github: 'https://gitlab.com/LazyLibrarian/LazyLibrarian', version: 'latest', stars: '-', subtitle: 'Book Manager', port: 5299, image: 'lscr.io/linuxserver/lazylibrarian', description: 'Automated book downloader for eBooks and audiobooks from usenet and torrents.', category: 'media' },
  maintainerr: { name: 'Maintainerr', logo: '/logos/maintainerr.png', github: 'https://github.com/jorenn92/Maintainerr', version: 'v2.2.1', stars: '1.2k', subtitle: 'Plex Maintenance', port: 6246, image: 'ghcr.io/jorenn92/maintainerr', description: 'Automated media maintenance for Plex. Clean up old or unwatched content.', category: 'media' },
  unpackerr: { name: 'Unpackerr', logo: '/logos/unpackerr.png', github: 'https://github.com/Unpackerr/unpackerr', version: 'v0.14.5', stars: '1k', subtitle: 'Archive Extractor', port: 0, image: 'golift/unpackerr', description: 'Extracts downloaded archives for Sonarr, Radarr, Lidarr and Readarr.', category: 'media' },
  readarr: { name: 'Readarr', logo: '/logos/readarr.png', github: 'https://github.com/Readarr/Readarr', version: 'v0.4.5', stars: '3k', subtitle: 'Books & Audiobooks', port: 8787, image: 'lscr.io/linuxserver/readarr', description: 'Book, audiobook and comic collection manager for usenet and torrent users.', category: 'media' },
  whisparr: { name: 'Whisparr', logo: '/logos/whisparr.png', github: 'https://github.com/whisparr/whisparr', version: 'nightly', stars: '-', subtitle: 'Adult Content', port: 6969, image: 'cr.hotio.dev/hotio/whisparr', description: 'Adult content collection manager. Automatically searches, downloads and manages content.', category: 'media' },
  notifiarr: { name: 'Notifiarr', logo: '/logos/notifiarr.png', github: 'https://github.com/Notifiarr/notifiarr', version: 'latest', stars: '800', subtitle: 'Notifications', port: 5454, image: 'golift/notifiarr', description: 'Unified notification client for Sonarr, Radarr, Lidarr, Readarr, Prowlarr, Plex and Tautulli.', category: 'media' },
  tdarr: { name: 'Tdarr', logo: '/logos/tdarr.png', github: 'https://github.com/HaveAGitGat/Tdarr', version: 'v2.27.2', stars: '3k', subtitle: 'Transcode Manager', port: 8265, image: 'ghcr.io/haveagitgat/tdarr', description: 'Distributed transcoding system for automating media library transcode and remux management.', category: 'media' },
  handbrake: { name: 'HandBrake', logo: '/logos/handbrake.png', github: 'https://github.com/HandBrake/HandBrake', version: 'v1.8.2', stars: '18k', subtitle: 'Video Transcoder', port: 5800, image: 'jlesage/handbrake', description: 'Open-source video transcoder with web GUI. Convert videos to modern formats.', category: 'media' },

  komga: { name: 'Komga', logo: '/logos/komga.png', github: 'https://github.com/gotson/komga', version: 'v1.14.1', stars: '4.5k', subtitle: 'Comic/Book Server', port: 25600, image: 'gotson/komga', description: 'Free and open-source manga, comic and book media server with OPDS support and web reader.', category: 'server' },
  stash: { name: 'Stash', logo: '/logos/stash.png', github: 'https://github.com/stashapp/stash', version: 'v0.27.2', stars: '9k', subtitle: 'Media Organizer', port: 9999, image: 'stashapp/stash', description: 'Organizer and metadata scraper for your personal media collection with tagging and filtering.', category: 'server' },

  // ── Media Servers ──
  plex: { name: 'Plex', logo: '/logos/plex.png', github: 'https://github.com/plexinc/pms-docker', version: 'latest', stars: '3k', subtitle: 'Media Server', port: 32400, image: 'lscr.io/linuxserver/plex', description: 'Stream movies, TV, music and more to any device.', category: 'server' },
  jellyfin: { name: 'Jellyfin', logo: '/logos/jellyfin.png', github: 'https://github.com/jellyfin/jellyfin', version: 'v10.10.6', stars: '38k', subtitle: 'Media Server', port: 8097, image: 'lscr.io/linuxserver/jellyfin', description: 'Free, open-source media server. Stream to any device.', category: 'server' },
  emby: { name: 'Emby', logo: '/logos/emby.png', github: 'https://github.com/MediaBrowser/Emby', version: 'v4.8.11', stars: '3.7k', subtitle: 'Media Server', port: 8096, image: 'lscr.io/linuxserver/emby', description: 'Media server that organizes and streams your personal collection.', category: 'server' },
  airsonic: { name: 'Airsonic', logo: '/logos/airsonic.png', github: 'https://github.com/airsonic-advanced/airsonic-advanced', version: 'v11.1.4', stars: '1.3k', subtitle: 'Music Server', port: 4040, image: 'lscr.io/linuxserver/airsonic-advanced', description: 'Free, web-based media streamer for your music collection.', category: 'server' },
  navidrome: { name: 'Navidrome', logo: '/logos/navidrome.png', github: 'https://github.com/navidrome/navidrome', version: 'v0.54.5', stars: '13k', subtitle: 'Music Server', port: 4533, image: 'deluan/navidrome', description: 'Modern music server compatible with Subsonic/Airsonic clients.', category: 'server' },
  calibreweb: { name: 'Calibre-Web', logo: '/logos/calibreweb.png', github: 'https://github.com/janeczku/calibre-web', version: 'v0.6.24', stars: '13.5k', subtitle: 'eBook Server', port: 8084, image: 'lscr.io/linuxserver/calibre-web', description: 'Web app for browsing, reading and downloading eBooks.', category: 'server' },
  tautulli: { name: 'Tautulli', logo: '/logos/tautulli.png', github: 'https://github.com/Tautulli/Tautulli', version: 'v2.15.1', stars: '5.9k', subtitle: 'Plex Monitoring', port: 8181, image: 'lscr.io/linuxserver/tautulli', description: 'Monitoring and tracking tool for Plex Media Server statistics.', category: 'server' },
  audiobookshelf: { name: 'Audiobookshelf', logo: '/logos/audiobookshelf.png', github: 'https://github.com/advplyr/audiobookshelf', version: 'v2.17.5', stars: '7.5k', subtitle: 'Audiobook Server', port: 13378, image: 'ghcr.io/advplyr/audiobookshelf', description: 'Self-hosted audiobook and podcast server with web-based player.', category: 'server' },
  kavita: { name: 'Kavita', logo: '/logos/kavita.png', github: 'https://github.com/Kareadita/Kavita', version: 'v0.8.4', stars: '7k', subtitle: 'Book/Comic Reader', port: 5000, image: 'jvmilazz0/kavita', description: 'Fast, feature-rich manga, comic, and book reader server.', category: 'server' },
  ubooquity: { name: 'Ubooquity', logo: '/logos/ubooquity.png', github: 'https://github.com/vaemendis/ubooquity', version: 'latest', stars: '-', subtitle: 'Comic/eBook Server', port: 2202, image: 'lscr.io/linuxserver/ubooquity', description: 'Home server for comics and eBooks with OPDS feed and web reader.', category: 'server' },

  // ── Indexers ──
  prowlarr: { name: 'Prowlarr', logo: '/logos/prowlarr.png', github: 'https://github.com/Prowlarr/Prowlarr', version: 'v1.31.2', stars: '4k', subtitle: 'Indexer Manager', port: 9696, image: 'lscr.io/linuxserver/prowlarr', description: 'Indexer manager that integrates with Sonarr, Radarr, Lidarr.', category: 'indexer' },
  jackett: { name: 'Jackett', logo: '/logos/jackett.png', github: 'https://github.com/Jackett/Jackett', version: 'v0.22.1', stars: '13k', subtitle: 'Indexer Proxy', port: 9117, image: 'lscr.io/linuxserver/jackett', description: 'Proxy server translating queries for tracker-site-specific searches.', category: 'indexer' },
  nzbhydra2: { name: 'NZBHydra2', logo: '/logos/nzbhydra2.png', github: 'https://github.com/theotherp/nzbhydra2', version: 'v7.5.0', stars: '1.1k', subtitle: 'Usenet Meta Search', port: 5076, image: 'lscr.io/linuxserver/nzbhydra2', description: 'Meta search for usenet indexers. Unified access to multiple newznab indexers.', category: 'indexer' },
  flaresolverr: { name: 'FlareSolverr', logo: '/logos/flaresolverr.png', github: 'https://github.com/FlareSolverr/FlareSolverr', version: 'v3.3.21', stars: '8.5k', subtitle: 'Cloudflare Bypass', port: 8191, image: 'ghcr.io/flaresolverr/flaresolverr', description: 'Proxy server to bypass Cloudflare protection for scrapers.', category: 'indexer' },

  // ── Request Management ──
  ombi: { name: 'Ombi', logo: '/logos/ombi.png', github: 'https://github.com/Ombi-app/Ombi', version: 'v4.44.1', stars: '3.8k', subtitle: 'Media Requests', port: 3579, image: 'lscr.io/linuxserver/ombi', description: 'Self-hosted app for Plex/Emby users to request content.', category: 'request' },
  overseerr: { name: 'Overseerr', logo: '/logos/overseerr.png', github: 'https://github.com/sct/overseerr', version: 'v1.33.2', stars: '4k', subtitle: 'Media Requests', port: 5055, image: 'lscr.io/linuxserver/overseerr', description: 'Modern media request and discovery tool for Plex, Sonarr and Radarr.', category: 'request' },
  doplarr: { name: 'Doplarr', logo: '/logos/doplarr.png', github: 'https://github.com/kiranshila/Doplarr', version: 'latest', stars: '500', subtitle: 'Discord Bot', port: 0, image: 'lscr.io/linuxserver/doplarr', description: 'Discord bot for requesting media through Overseerr or Ombi.', category: 'request' },
  requestrr: { name: 'Requestrr', logo: '/logos/requestrr.png', github: 'https://github.com/thomst08/requestrr', version: 'v2.1.6', stars: '800', subtitle: 'Discord Bot', port: 4545, image: 'thomst08/requestrr', description: 'Discord chatbot for requesting movies and TV shows via Sonarr, Radarr and Overseerr.', category: 'request' },

  // ── Sync & Cloud ──
  resilio: { name: 'Resilio Sync', logo: '/logos/resilio.svg', github: 'https://github.com/bt-sync', version: 'latest', stars: '-', subtitle: 'P2P File Sync', port: 8888, image: 'lscr.io/linuxserver/resilio-sync', description: 'Fast, reliable file sync using peer-to-peer technology.', category: 'sync' },
  nextcloud: { name: 'Nextcloud', logo: '/logos/nextcloud.png', github: 'https://github.com/nextcloud/server', version: 'v30.0.6', stars: '29k', subtitle: 'Cloud Storage', port: 8443, image: 'lscr.io/linuxserver/nextcloud', description: 'Self-hosted cloud storage, contacts, calendar and collaboration.', category: 'sync' },
  syncthing: { name: 'Syncthing', logo: '/logos/syncthing.png', github: 'https://github.com/syncthing/syncthing', version: 'v1.29.2', stars: '68k', subtitle: 'File Sync', port: 8384, image: 'lscr.io/linuxserver/syncthing', description: 'Continuous, decentralized file synchronization between devices.', category: 'sync' },
  rclone: { name: 'Rclone', logo: '/logos/rclone.png', github: 'https://github.com/rclone/rclone', version: 'v1.69.1', stars: '49k', subtitle: 'Cloud Storage CLI', port: 5572, image: 'rclone/rclone', description: 'Manage files on cloud storage with web GUI.', category: 'sync' },
  filebrowser: { name: 'File Browser', logo: '/logos/filebrowser.png', github: 'https://github.com/filebrowser/filebrowser', version: 'v2.31.2', stars: '28k', subtitle: 'File Manager', port: 8092, image: 'filebrowser/filebrowser:latest', description: 'Web-based file manager with a clean interface for uploading, deleting, previewing and sharing files.', category: 'sync' },

  // ── Network ──
  wireguard: { name: 'WireGuard', logo: '/logos/wireguard.png', github: 'https://github.com/WireGuard', version: 'latest', stars: '-', subtitle: 'VPN', port: 51820, image: 'lscr.io/linuxserver/wireguard', description: 'Fast, modern, secure VPN tunnel.', category: 'network' },
  thelounge: { name: 'The Lounge', logo: '/logos/thelounge.png', github: 'https://github.com/thelounge/thelounge', version: 'v4.4.3', stars: '5.8k', subtitle: 'IRC Client', port: 9000, image: 'lscr.io/linuxserver/thelounge', description: 'Modern, self-hosted IRC client with always-on connectivity and web interface.', category: 'network' },
  znc: { name: 'ZNC', logo: '/logos/znc.png', github: 'https://github.com/znc/znc', version: 'v1.9.1', stars: '2k', subtitle: 'IRC Bouncer', port: 6501, image: 'lscr.io/linuxserver/znc', description: 'Advanced IRC bouncer that stays connected and buffers messages.', category: 'network' },

  // ── System ──
  uptimekuma: { name: 'Uptime Kuma', logo: '/logos/uptimekuma.png', github: 'https://github.com/louislam/uptime-kuma', version: 'v1.23.16', stars: '62k', subtitle: 'Monitoring', port: 3010, image: 'louislam/uptime-kuma', description: 'Self-hosted monitoring tool with beautiful status pages and notifications.', category: 'system' },
  mariadb: { name: 'MariaDB', logo: '/logos/mariadb.png', github: 'https://github.com/MariaDB/server', version: 'v11.7', stars: '5.8k', subtitle: 'Database', port: 3306, image: 'lscr.io/linuxserver/mariadb', description: 'Community-developed MySQL fork. Reliable, high performance database server.', category: 'system' },
};

const CATEGORIES: { key: string; label: string; color: string }[] = [
  { key: 'download', label: 'Download Automation', color: '#4ade80' },
  { key: 'client', label: 'Download Clients', color: '#34d399' },
  { key: 'media', label: 'Media Management', color: '#c084fc' },
  { key: 'server', label: 'Media Servers', color: '#818cf8' },
  { key: 'indexer', label: 'Indexers', color: '#fbbf24' },
  { key: 'request', label: 'Request Management', color: '#60a5fa' },
  { key: 'sync', label: 'Sync & Cloud', color: '#2dd4bf' },
  { key: 'network', label: 'Network', color: '#f472b6' },
  { key: 'system', label: 'System', color: '#f97316' },
];

export default function Catalog() {
  const [services, setServices] = useState<Service[]>([]);
  const [loading, setLoading] = useState(true);
  const [installingId, setInstallingId] = useState<string | null>(null);
  const [installMode, setInstallMode] = useState<string>('docker');

  useEffect(() => {
    Promise.all([
      api.listServices().then((s) => s || []).catch(() => [] as Service[]),
      api.getConfig().catch(() => ({ version: '', install_mode: 'docker' })),
    ]).then(([svc, cfg]) => {
      setServices(svc);
      setInstallMode(cfg.install_mode || 'docker');
    }).finally(() => setLoading(false));
  }, []);

  useEffect(() => {
    if (!installingId) return;
    const iv = setInterval(async () => {
      try {
        const svc = await api.listServices() || [];
        setServices(svc);
        if (svc.some((s) => s.id === installingId)) setInstallingId(null);
      } catch { /* ignore */ }
    }, 3000);
    return () => clearInterval(iv);
  }, [installingId]);

  const handleInstall = async (id: string) => {
    setInstallingId(id);
    try { await api.installService(id); } catch { setInstallingId(null); }
  };

  if (loading) {
    return (
      <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'center', height: 256 }}>
        <div style={{ width: 32, height: 32, border: '3px solid rgba(59,130,246,0.2)', borderTopColor: '#3b82f6', borderRadius: '50%', animation: 'spin 0.8s linear infinite' }} />
      </div>
    );
  }

  return (
    <div>
      <div style={{ marginBottom: 24 }}>
        <h2 style={{ fontSize: 22, fontWeight: 700, color: 'var(--text-primary)', marginBottom: 4 }}>App Catalog</h2>
        <p style={{ fontSize: 13, color: 'var(--text-tertiary)' }}>
          {Object.keys(APPS).length} applications available
          <span style={{
            marginLeft: 8, padding: '2px 8px', borderRadius: 4, fontSize: 10, fontWeight: 600,
            background: installMode === 'native' ? 'rgba(74,222,128,0.1)' : 'rgba(96,165,250,0.1)',
            border: installMode === 'native' ? '1px solid rgba(74,222,128,0.2)' : '1px solid rgba(96,165,250,0.2)',
            color: installMode === 'native' ? 'var(--accent-green)' : 'var(--accent-blue-light)',
          }}>
            {installMode === 'native' ? 'Native Mode' : 'Docker Mode'}
          </span>
        </p>
      </div>

      {CATEGORIES.map((cat) => {
        const apps = Object.entries(APPS).filter(([, m]) => m.category === cat.key);
        if (apps.length === 0) return null;
        return (
          <div key={cat.key} style={{ marginBottom: 24 }}>
            <div style={{ display: 'flex', alignItems: 'center', gap: 8, marginBottom: 10 }}>
              <div style={{ width: 3, height: 14, borderRadius: 2, background: cat.color }} />
              <h3 style={{ fontSize: 13, fontWeight: 600, color: 'var(--text-secondary)' }}>{cat.label}</h3>
              <span style={{ fontSize: 11, color: 'var(--text-muted)' }}>({apps.length})</span>
            </div>

            <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fill, minmax(320px, 1fr))', gap: 10 }}>
              {apps.map(([id, meta]) => {
                const svc = services.find((s) => s.id === id);
                const installed = !!svc;
                const isInstalling = installingId === id;

                return (
                  <div key={id} style={{
                    background: 'var(--bg-card)', border: '1px solid var(--border-primary)',
                    borderRadius: 10, padding: 14,
                  }}>
                    <div style={{ display: 'flex', alignItems: 'center', gap: 10, marginBottom: 8 }}>
                      <div style={{
                        width: 32, height: 32, borderRadius: 8, background: 'var(--bg-primary)', flexShrink: 0,
                        display: 'flex', alignItems: 'center', justifyContent: 'center',
                        fontSize: 14, fontWeight: 700, color: cat.color, overflow: 'hidden',
                      }}>
                        <img src={meta.logo} alt="" style={{
                          width: 32, height: 32, objectFit: 'contain',
                        }} onError={(e) => {
                          const el = e.target as HTMLImageElement;
                          el.style.display = 'none';
                          el.parentElement!.textContent = meta.name.charAt(0);
                        }} />
                      </div>
                      <div style={{ flex: 1, minWidth: 0 }}>
                        <div style={{ display: 'flex', alignItems: 'center', gap: 5, flexWrap: 'wrap' }}>
                          <span style={{ fontSize: 13, fontWeight: 700, color: 'var(--text-primary)' }}>{meta.name}</span>
                          {installed && (
                            <div style={{
                              display: 'flex', alignItems: 'center', gap: 3, padding: '1px 6px', borderRadius: 8,
                              background: 'rgba(34,197,94,0.1)', border: '1px solid rgba(34,197,94,0.2)',
                            }}>
                              <div style={{ width: 4, height: 4, borderRadius: '50%', background: 'var(--accent-green)' }} />
                              <span style={{ fontSize: 9, color: 'var(--accent-green)', fontWeight: 500 }}>{svc?.status || 'installed'}</span>
                            </div>
                          )}
                        </div>
                        <p style={{ fontSize: 10, color: 'var(--text-muted)' }}>{meta.subtitle}</p>
                      </div>
                    </div>

                    <p style={{ fontSize: 11, color: 'var(--text-tertiary)', lineHeight: 1.4, marginBottom: 10 }}>{meta.description}</p>

                    <div style={{ display: 'flex', flexWrap: 'wrap', gap: 5, marginBottom: 10 }}>
                      <span style={{ padding: '2px 6px', borderRadius: 4, background: 'rgba(59,130,246,0.08)', border: '1px solid rgba(59,130,246,0.12)', fontSize: 9, color: 'var(--accent-blue-light)', fontFamily: 'monospace', fontWeight: 600 }}>{meta.version}</span>
                      {meta.stars !== '-' && <span style={{ padding: '2px 6px', borderRadius: 4, background: 'rgba(245,158,11,0.08)', border: '1px solid rgba(245,158,11,0.12)', fontSize: 9, color: 'var(--accent-yellow)', fontWeight: 600 }}>{meta.stars}</span>}
                      <span style={{ padding: '2px 6px', borderRadius: 4, background: 'rgba(100,116,139,0.08)', border: '1px solid rgba(100,116,139,0.12)', fontSize: 9, color: 'var(--text-tertiary)', fontFamily: 'monospace' }}>:{meta.port}</span>
                    </div>

                    <div style={{ display: 'flex', gap: 5 }}>
                      {installed ? (
                        <>
                          <a href={`http://${window.location.hostname}:${meta.port}`} target="_blank" rel="noreferrer" style={{
                            display: 'flex', alignItems: 'center', gap: 4, padding: '5px 10px', borderRadius: 6,
                            border: '1px solid rgba(59,130,246,0.3)', background: 'rgba(59,130,246,0.1)',
                            color: 'var(--accent-blue-light)', fontSize: 11, fontWeight: 600, textDecoration: 'none',
                          }}><ExternalLink size={11} /> Open</a>
                          <div style={{
                            display: 'flex', alignItems: 'center', gap: 4, padding: '5px 10px', borderRadius: 6,
                            background: 'rgba(34,197,94,0.06)', border: '1px solid rgba(34,197,94,0.15)',
                            color: 'var(--accent-green)', fontSize: 11, fontWeight: 500,
                          }}><Check size={11} /> Installed</div>
                        </>
                      ) : (
                        <button onClick={() => handleInstall(id)} disabled={isInstalling} style={{
                          display: 'flex', alignItems: 'center', gap: 4, padding: '5px 14px', borderRadius: 6,
                          border: 'none', background: isInstalling ? '#1e40af' : '#2563eb', color: '#fff',
                          fontSize: 11, fontWeight: 600, cursor: 'pointer', opacity: isInstalling ? 0.7 : 1,
                        }}>
                          {isInstalling ? (
                            <><div style={{ width: 11, height: 11, border: '2px solid rgba(255,255,255,0.3)', borderTopColor: '#fff', borderRadius: '50%', animation: 'spin 0.8s linear infinite' }} /> Installing...</>
                          ) : (
                            <><Download size={11} /> Install</>
                          )}
                        </button>
                      )}
                      <a href={meta.github} target="_blank" rel="noreferrer" style={{
                        display: 'flex', alignItems: 'center', padding: '5px 8px', borderRadius: 6,
                        background: 'rgba(100,116,139,0.08)', border: '1px solid rgba(100,116,139,0.15)',
                        color: 'var(--text-tertiary)', textDecoration: 'none',
                      }}><Github size={11} /></a>
                    </div>
                  </div>
                );
              })}
            </div>
          </div>
        );
      })}
    </div>
  );
}
