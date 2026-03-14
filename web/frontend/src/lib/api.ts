const API_BASE = '/api';

async function request<T>(path: string, options?: RequestInit): Promise<T> {
  const token = localStorage.getItem('velour_token');
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    ...(token ? { Authorization: `Bearer ${token}` } : {}),
  };

  const res = await fetch(`${API_BASE}${path}`, {
    ...options,
    headers: { ...headers, ...options?.headers },
  });

  if (res.status === 401) {
    localStorage.removeItem('velour_token');
    localStorage.removeItem('velour_user');
    throw new Error('Unauthorized');
  }

  if (!res.ok) {
    const error = await res.json().catch(() => ({ error: 'Unknown error' }));
    throw new Error(error.error || 'Request failed');
  }

  return res.json();
}

export const api = {
  // Health
  health: () => request<{ status: string; version: string }>('/health'),

  // Setup
  setupStatus: () => request<{ needs_setup: boolean }>('/setup/status'),
  setup: (username: string, password: string) =>
    request('/setup', { method: 'POST', body: JSON.stringify({ username, password }) }),

  // Auth
  login: (username: string, password: string) =>
    request<{ token: string; user: any }>('/auth/login', {
      method: 'POST',
      body: JSON.stringify({ username, password }),
    }),

  // System
  systemInfo: () => request<SystemInfo>('/system'),

  // Services
  listServices: () => request<Service[]>('/services'),
  serviceCatalog: () => request<ServiceDefinition[]>('/services/catalog'),
  installService: (id: string) => request(`/services/${id}/install`, { method: 'POST' }),
  startService: (id: string) => request(`/services/${id}/start`, { method: 'POST' }),
  stopService: (id: string) => request(`/services/${id}/stop`, { method: 'POST' }),
  restartService: (id: string) => request(`/services/${id}/restart`, { method: 'POST' }),
  removeService: (id: string) => request(`/services/${id}`, { method: 'DELETE' }),

  // Config
  getConfig: () => request<{ version: string; install_mode: string }>('/config'),

  // Automation
  listRules: () => request<AutomationRule[]>('/automation/rules'),
  createRule: (rule: CreateRuleRequest) =>
    request<AutomationRule>('/automation/rules', { method: 'POST', body: JSON.stringify(rule) }),
  getRule: (id: string) => request<AutomationRule>(`/automation/rules/${id}`),
  updateRule: (id: string, rule: CreateRuleRequest) =>
    request<AutomationRule>(`/automation/rules/${id}`, { method: 'PUT', body: JSON.stringify(rule) }),
  deleteRule: (id: string) => request(`/automation/rules/${id}`, { method: 'DELETE' }),
  toggleRule: (id: string, enabled: boolean) =>
    request(`/automation/rules/${id}/toggle`, { method: 'POST', body: JSON.stringify({ enabled }) }),
  listTemplates: () => request<ActionTemplate[]>('/automation/templates'),

  // Network stats
  networkMonths: () => request<MonthlyNetStats[]>('/network/months'),
  networkMonth: (month: string) => request<MonthlyNetStats>(`/network/month/${month}`),

  // Integrations
  downloads: () => request<DownloadsResponse>('/integrations/downloads'),
  qbitTorrents: () => request<QbitTorrent[]>('/integrations/qbit/torrents'),
  qbitTransfer: () => request<QbitTransferInfo>('/integrations/qbit/transfer'),
  sonarrCalendar: () => request<SonarrCalendarEntry[]>('/integrations/sonarr/calendar'),
  radarrCalendar: () => request<RadarrCalendarEntry[]>('/integrations/radarr/calendar'),
};

export interface SystemInfo {
  hostname: string;
  os: string;
  platform: string;
  kernel: string;
  uptime: number;
  uptime_human: string;
  cpu: { model: string; cores: number; threads: number; usage: number };
  memory: { total: number; used: number; free: number; usage_percent: number };
  disk: { total: number; used: number; free: number; usage_percent: number };
  network: { bytes_sent: number; bytes_recv: number };
}

export interface Service {
  id: string;
  name: string;
  description: string;
  icon: string;
  category: string;
  port: number;
  web_url: string;
  status: 'running' | 'stopped' | 'installing' | 'error';
  type: string;
  image: string;
  installed: boolean;
}

export interface AutomationRule {
  id: string;
  name: string;
  enabled: boolean;
  trigger: {
    type: 'torrent_done' | 'service_start' | 'service_stop' | 'schedule';
    service_id: string;
    cron?: string;
  };
  action: {
    type: 'exec_in_service' | 'webhook' | 'restart_service';
    service_id?: string;
    command?: string;
    args?: string[];
    url?: string;
    template?: string;
  };
  created_at: string;
  updated_at: string;
  last_run_at?: string;
  last_run_ok?: boolean;
  last_run_log?: string;
  run_count: number;
}

export interface CreateRuleRequest {
  name: string;
  trigger: AutomationRule['trigger'];
  action: AutomationRule['action'];
}

export interface ActionTemplate {
  id: string;
  name: string;
  description: string;
  service_id: string;
  command: string;
  args: string[];
  icon: string;
}

export interface MonthlyNetStats {
  month: string;
  bytes_sent: number;
  bytes_recv: number;
}

export interface DownloadItem {
  name: string;
  size: number;
  progress: number;
  dlspeed: number;
  upspeed: number;
  state: string;
  eta: number;
  client: string;
  added_on: number;
  seeds: number;
  peers: number;
}

export interface DownloadsResponse {
  items: DownloadItem[];
  clients: string[];
}

export interface QbitTorrent {
  name: string;
  size: number;
  progress: number;
  dlspeed: number;
  upspeed: number;
  state: string;
  eta: number;
  category: string;
  added_on: number;
  num_seeds: number;
  num_leechs: number;
}

export interface QbitTransferInfo {
  dl_info_speed: number;
  up_info_speed: number;
  dl_info_data: number;
  up_info_data: number;
}

export interface SonarrCalendarEntry {
  seriesId: number;
  seasonNumber: number;
  episodeNumber: number;
  title: string;
  airDateUtc: string;
  hasFile: boolean;
  monitored: boolean;
  series?: {
    title: string;
    images?: { coverType: string; remoteUrl: string }[];
  };
}

export interface RadarrCalendarEntry {
  title: string;
  year: number;
  physicalRelease?: string;
  digitalRelease?: string;
  inCinemas?: string;
  hasFile: boolean;
  monitored: boolean;
  images?: { coverType: string; remoteUrl: string }[];
}

export interface ServiceDefinition {
  id: string;
  name: string;
  description: string;
  icon: string;
  category: string;
  image: string;
  ports: { host: number; container: number; protocol: string }[];
  volumes: { host: string; container: string }[];
  env: Record<string, string>;
  install_types: ('docker' | 'native')[];
  native?: {
    method: string;
    service_name: string;
    port: number;
  };
}
