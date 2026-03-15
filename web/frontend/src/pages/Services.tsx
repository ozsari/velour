import { useEffect, useState } from 'react';
import { Package, Play, Square, RotateCcw, Trash2, ExternalLink } from 'lucide-react';
import { api, type Service } from '../lib/api';
import { APP_LOGOS } from '../lib/logos';

const catStyle = (cat: string): React.CSSProperties => {
  const m: Record<string, { bg: string; color: string; border: string }> = {
    media: { bg: 'rgba(168,85,247,0.1)', color: '#c084fc', border: 'rgba(168,85,247,0.2)' },
    download: { bg: 'rgba(34,197,94,0.1)', color: '#4ade80', border: 'rgba(34,197,94,0.2)' },
    system: { bg: 'rgba(59,130,246,0.1)', color: '#60a5fa', border: 'rgba(59,130,246,0.2)' },
    network: { bg: 'rgba(245,158,11,0.1)', color: '#fbbf24', border: 'rgba(245,158,11,0.2)' },
  };
  const c = m[cat] || { bg: 'rgba(100,116,139,0.1)', color: '#94a3b8', border: 'rgba(100,116,139,0.2)' };
  return { fontSize: 11, padding: '2px 8px', borderRadius: 6, background: c.bg, color: c.color, border: `1px solid ${c.border}`, fontWeight: 500 };
};

export default function Services() {
  const [services, setServices] = useState<Service[]>([]);
  const [loading, setLoading] = useState(true);

  const fetchServices = async () => {
    try { setServices((await api.listServices()) || []); } catch {} finally { setLoading(false); }
  };

  useEffect(() => {
    fetchServices();
    const interval = setInterval(fetchServices, 5000);
    return () => clearInterval(interval);
  }, []);

  const act = async (fn: (id: string) => Promise<any>, id: string) => {
    try { await fn(id); setTimeout(fetchServices, 1000); } catch {}
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
        <h2 style={{ fontSize: 24, fontWeight: 700, color: 'var(--text-primary)', marginBottom: 4 }}>Installed Services</h2>
        <p style={{ fontSize: 14, color: 'var(--text-tertiary)' }}>Manage your running applications</p>
      </div>

      {services.length === 0 ? (
        <div style={{ textAlign: 'center', padding: '80px 0' }}>
          <Package size={48} style={{ margin: '0 auto 16px', color: 'var(--text-faint)' }} />
          <h3 style={{ fontSize: 18, fontWeight: 500, color: 'var(--text-secondary)', marginBottom: 4 }}>No services installed</h3>
          <p style={{ fontSize: 14, color: 'var(--text-muted)' }}>Head over to the App Catalog to install your first service</p>
        </div>
      ) : (
        <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fill, minmax(380px, 1fr))', gap: 16 }}>
          {services.map((s) => {
            const running = s.status === 'running';
            const logo = APP_LOGOS[s.id];
            return (
              <div key={s.id} style={{
                background: 'var(--bg-card)', border: '1px solid var(--border-primary)',
                borderRadius: 14, padding: 20,
              }}>
                {/* Header */}
                <div style={{ display: 'flex', alignItems: 'flex-start', justifyContent: 'space-between', marginBottom: 16 }}>
                  <div style={{ display: 'flex', alignItems: 'center', gap: 12 }}>
                    <div style={{
                      width: 32, height: 32, borderRadius: 8, background: 'var(--bg-primary)', flexShrink: 0,
                      display: 'flex', alignItems: 'center', justifyContent: 'center',
                      fontSize: 14, fontWeight: 700, color: 'var(--text-tertiary)', overflow: 'hidden',
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
                    <div>
                      <h3 style={{ fontSize: 15, fontWeight: 600, color: 'var(--text-primary)' }}>{s.name}</h3>
                      <p style={{ fontSize: 12, color: 'var(--text-tertiary)', marginTop: 2 }}>{s.description}</p>
                    </div>
                  </div>
                  <div style={{ display: 'flex', alignItems: 'center', gap: 6 }}>
                    <span style={{
                      width: 8, height: 8, borderRadius: '50%',
                      background: running ? '#22c55e' : 'var(--text-muted)',
                      boxShadow: running ? '0 0 8px rgba(34,197,94,0.5)' : 'none',
                    }} />
                    <span style={{ fontSize: 12, fontWeight: 500, color: running ? 'var(--accent-green)' : 'var(--text-tertiary)' }}>
                      {s.status}
                    </span>
                  </div>
                </div>

                {/* Tags */}
                <div style={{ display: 'flex', alignItems: 'center', gap: 8, marginBottom: 16 }}>
                  <span style={catStyle(s.category)}>{s.category}</span>
                  <span style={{ fontSize: 12, color: 'var(--text-muted)' }}>Port: {s.port}</span>
                </div>

                {/* Actions */}
                <div style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
                  {running ? (
                    <button onClick={() => act(api.stopService, s.id)} style={{
                      display: 'flex', alignItems: 'center', gap: 6, padding: '6px 12px', borderRadius: 8,
                      background: 'rgba(239,68,68,0.1)', color: 'var(--accent-red)', border: 'none', fontSize: 13, cursor: 'pointer',
                    }}><Square size={14} /> Stop</button>
                  ) : (
                    <button onClick={() => act(api.startService, s.id)} style={{
                      display: 'flex', alignItems: 'center', gap: 6, padding: '6px 12px', borderRadius: 8,
                      background: 'rgba(34,197,94,0.1)', color: 'var(--accent-green)', border: 'none', fontSize: 13, cursor: 'pointer',
                    }}><Play size={14} /> Start</button>
                  )}
                  <button onClick={() => act(api.restartService, s.id)} style={{
                    display: 'flex', alignItems: 'center', padding: '6px 10px', borderRadius: 8,
                    background: 'rgba(51,65,85,0.3)', color: 'var(--text-secondary)', border: 'none', cursor: 'pointer',
                  }}><RotateCcw size={14} /></button>
                  {running && s.web_url && (
                    <a href={s.web_url.replace('localhost', window.location.hostname)} target="_blank" rel="noopener noreferrer" style={{
                      display: 'flex', alignItems: 'center', gap: 6, padding: '6px 12px', borderRadius: 8,
                      background: 'rgba(59,130,246,0.1)', color: 'var(--accent-blue-light)', fontSize: 13, textDecoration: 'none',
                    }}><ExternalLink size={14} /> Open</a>
                  )}
                  <button onClick={() => act(api.removeService, s.id)} style={{
                    display: 'flex', alignItems: 'center', padding: '6px 10px', borderRadius: 8,
                    background: 'transparent', color: 'var(--text-muted)', border: 'none', cursor: 'pointer', marginLeft: 'auto',
                  }}><Trash2 size={14} /></button>
                </div>
              </div>
            );
          })}
        </div>
      )}
    </div>
  );
}
