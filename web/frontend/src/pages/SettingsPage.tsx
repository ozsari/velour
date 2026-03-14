import { useEffect, useState } from 'react';
import { Shield, Info } from 'lucide-react';
import { api } from '../lib/api';

export default function SettingsPage() {
  const [version, setVersion] = useState('');

  useEffect(() => {
    api.health().then((data) => setVersion(data.version));
  }, []);

  return (
    <div>
      <div style={{ marginBottom: 24 }}>
        <h2 style={{ fontSize: 24, fontWeight: 700, color: 'var(--text-primary)', marginBottom: 4 }}>Settings</h2>
        <p style={{ fontSize: 14, color: 'var(--text-tertiary)' }}>Configure your Velour instance</p>
      </div>

      <div style={{ display: 'flex', flexDirection: 'column', gap: 16 }}>
        {/* About */}
        <div style={{
          background: 'var(--bg-card)', border: '1px solid var(--border-primary)',
          borderRadius: 14, padding: 24,
        }}>
          <div style={{ display: 'flex', alignItems: 'center', gap: 12, marginBottom: 20 }}>
            <div style={{ padding: 10, borderRadius: 10, background: 'rgba(59,130,246,0.1)', color: 'var(--accent-blue-light)', display: 'flex' }}>
              <Info size={20} />
            </div>
            <h3 style={{ fontSize: 18, fontWeight: 600, color: 'var(--text-primary)' }}>About Velour</h3>
          </div>
          <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: 16 }}>
            <div>
              <span style={{ fontSize: 12, color: 'var(--text-tertiary)' }}>Version</span>
              <p style={{ fontSize: 14, fontWeight: 500, color: 'var(--text-primary)', marginTop: 2 }}>v{version}</p>
            </div>
            <div>
              <span style={{ fontSize: 12, color: 'var(--text-tertiary)' }}>License</span>
              <p style={{ fontSize: 14, fontWeight: 500, color: 'var(--text-primary)', marginTop: 2 }}>MIT</p>
            </div>
          </div>
        </div>

        {/* Security */}
        <div style={{
          background: 'var(--bg-card)', border: '1px solid var(--border-primary)',
          borderRadius: 14, padding: 24,
        }}>
          <div style={{ display: 'flex', alignItems: 'center', gap: 12, marginBottom: 20 }}>
            <div style={{ padding: 10, borderRadius: 10, background: 'rgba(245,158,11,0.1)', color: 'var(--accent-yellow)', display: 'flex' }}>
              <Shield size={20} />
            </div>
            <h3 style={{ fontSize: 18, fontWeight: 600, color: 'var(--text-primary)' }}>Security</h3>
          </div>
          <p style={{ fontSize: 14, color: 'var(--text-tertiary)', marginBottom: 16 }}>
            Manage authentication and access settings
          </p>
          <div style={{ display: 'flex', flexDirection: 'column', gap: 10 }}>
            {[
              { name: 'JWT Authentication', desc: 'Token-based session management', status: 'Active', color: 'var(--accent-green)', bg: 'rgba(34,197,94,0.1)', border: 'rgba(34,197,94,0.2)' },
              { name: 'OAuth / LDAP', desc: 'External authentication providers', status: 'Coming Soon', color: 'var(--text-tertiary)', bg: 'rgba(100,116,139,0.1)', border: 'rgba(100,116,139,0.2)' },
            ].map((item) => (
              <div key={item.name} style={{
                display: 'flex', alignItems: 'center', justifyContent: 'space-between',
                padding: 14, background: 'var(--bg-primary)', borderRadius: 10,
              }}>
                <div>
                  <p style={{ fontSize: 14, fontWeight: 500, color: 'var(--text-primary)' }}>{item.name}</p>
                  <p style={{ fontSize: 12, color: 'var(--text-muted)', marginTop: 2 }}>{item.desc}</p>
                </div>
                <span style={{
                  fontSize: 11, padding: '3px 10px', borderRadius: 6,
                  background: item.bg, color: item.color, border: `1px solid ${item.border}`,
                  fontWeight: 500,
                }}>
                  {item.status}
                </span>
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  );
}
