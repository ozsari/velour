import { NavLink } from 'react-router-dom';
import {
  LayoutDashboard,
  Package,
  Server,
  Settings,
  Zap,
  LogOut,
  Menu,
  X,
  Cpu,
  MemoryStick,
  HardDrive,
  ArrowUp,
  ArrowDown,
  Heart,
  Star,
  ExternalLink,
  Sun,
  Moon,
  ChevronsLeft,
  ChevronsRight,
} from 'lucide-react';
import { useState, useEffect } from 'react';
import { useSystemInfo, formatBytes, formatSpeed } from '../hooks/useSystemInfo';
import { useTheme } from '../hooks/useTheme';
import type { SystemInfoWithSpeed } from '../hooks/useSystemInfo';

interface SidebarProps {
  onLogout: () => void;
  collapsed: boolean;
  onToggleCollapse: () => void;
}

const EXPANDED_W = 256;
const COLLAPSED_W = 68;

const navItems = [
  { to: '/', icon: LayoutDashboard, label: 'Dashboard' },
  { to: '/services', icon: Package, label: 'Services' },
  { to: '/catalog', icon: Server, label: 'App Catalog' },
  { to: '/automation', icon: Zap, label: 'Automation' },
  { to: '/settings', icon: Settings, label: 'Settings' },
];

function MiniStat({ icon, label, value, percent, color }: {
  icon: React.ReactNode; label: string; value: string; percent?: number; color: string;
}) {
  const barColor = percent && percent > 90 ? 'var(--accent-red)' : percent && percent > 70 ? 'var(--accent-yellow)' : color;
  return (
    <div style={{ display: 'flex', alignItems: 'center', gap: 10, padding: '6px 0' }}>
      <div style={{ color, flexShrink: 0, display: 'flex' }}>{icon}</div>
      <div style={{ flex: 1, minWidth: 0 }}>
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'baseline', marginBottom: 3 }}>
          <span style={{ fontSize: 11, color: 'var(--text-tertiary)' }}>{label}</span>
          <span style={{ fontSize: 12, fontWeight: 600, color: 'var(--text-primary)', fontFamily: 'monospace' }}>{value}</span>
        </div>
        {percent !== undefined && (
          <div style={{ height: 3, background: 'var(--bg-tertiary)', borderRadius: 2, overflow: 'hidden' }}>
            <div style={{
              height: '100%', borderRadius: 2, background: barColor,
              width: `${Math.min(percent, 100)}%`, transition: 'width 0.5s ease',
            }} />
          </div>
        )}
      </div>
    </div>
  );
}

function SidebarStats({ system }: { system: SystemInfoWithSpeed }) {
  return (
    <div style={{ padding: '8px 16px', borderTop: '1px solid var(--border-primary)' }}>
      <span style={{ fontSize: 10, fontWeight: 600, color: 'var(--text-muted)', textTransform: 'uppercase', letterSpacing: 1 }}>System</span>
      <MiniStat icon={<Cpu size={14} />} label="CPU" value={`${system.cpu.usage.toFixed(0)}%`} percent={system.cpu.usage} color="var(--accent-blue-light)" />
      <MiniStat icon={<MemoryStick size={14} />} label="RAM" value={formatBytes(system.memory.used)} percent={system.memory.usage_percent} color="var(--accent-purple)" />
      <MiniStat icon={<HardDrive size={14} />} label="Disk" value={formatBytes(system.disk.used)} percent={system.disk.usage_percent} color="var(--accent-yellow)" />
      <div style={{ display: 'flex', alignItems: 'center', gap: 10, padding: '6px 0' }}>
        <div style={{ display: 'flex', alignItems: 'center', gap: 4, flex: 1 }}>
          <ArrowUp size={12} style={{ color: 'var(--accent-red)', flexShrink: 0 }} />
          <span style={{ fontSize: 11, color: 'var(--text-tertiary)' }}>Up</span>
          <span style={{ fontSize: 12, fontWeight: 600, color: 'var(--text-primary)', fontFamily: 'monospace', marginLeft: 'auto' }}>
            {formatSpeed(system.network_speed.upload)}
          </span>
        </div>
      </div>
      <div style={{ display: 'flex', alignItems: 'center', gap: 10, padding: '6px 0' }}>
        <div style={{ display: 'flex', alignItems: 'center', gap: 4, flex: 1 }}>
          <ArrowDown size={12} style={{ color: 'var(--accent-green)', flexShrink: 0 }} />
          <span style={{ fontSize: 11, color: 'var(--text-tertiary)' }}>Down</span>
          <span style={{ fontSize: 12, fontWeight: 600, color: 'var(--text-primary)', fontFamily: 'monospace', marginLeft: 'auto' }}>
            {formatSpeed(system.network_speed.download)}
          </span>
        </div>
      </div>
    </div>
  );
}

export default function Sidebar({ onLogout, collapsed, onToggleCollapse }: SidebarProps) {
  const [mobileOpen, setMobileOpen] = useState(false);
  const [isDesktop, setIsDesktop] = useState(window.innerWidth >= 1024);
  const system = useSystemInfo();
  const { theme, toggleTheme } = useTheme();

  useEffect(() => {
    const handler = () => {
      setIsDesktop(window.innerWidth >= 1024);
      if (window.innerWidth >= 1024) setMobileOpen(false);
    };
    window.addEventListener('resize', handler);
    return () => window.removeEventListener('resize', handler);
  }, []);

  const sidebarVisible = isDesktop || mobileOpen;
  const isCollapsed = isDesktop && collapsed;
  const sidebarWidth = isCollapsed ? COLLAPSED_W : EXPANDED_W;

  const linkStyle = (isActive: boolean): React.CSSProperties => ({
    display: 'flex',
    alignItems: 'center',
    gap: isCollapsed ? 0 : 12,
    justifyContent: isCollapsed ? 'center' : 'flex-start',
    padding: isCollapsed ? '12px' : '12px 16px',
    borderRadius: 10,
    fontSize: 14,
    fontWeight: 500,
    textDecoration: 'none',
    transition: 'all 0.15s',
    border: isActive ? '1px solid rgba(59,130,246,0.2)' : '1px solid transparent',
    background: isActive ? 'rgba(59,130,246,0.1)' : 'transparent',
    color: isActive ? 'var(--accent-blue-light)' : 'var(--text-tertiary)',
  });

  return (
    <>
      {/* Mobile hamburger */}
      {!isDesktop && (
        <button
          onClick={() => setMobileOpen(!mobileOpen)}
          style={{
            position: 'fixed', top: 16, left: 16, zIndex: 50,
            padding: 10, borderRadius: 10, border: 'none', cursor: 'pointer',
            background: 'var(--bg-secondary)', color: 'var(--text-primary)', display: 'flex',
            boxShadow: '0 2px 8px rgba(0,0,0,0.3)',
          }}
        >
          {mobileOpen ? <X size={22} /> : <Menu size={22} />}
        </button>
      )}

      {/* Overlay */}
      {mobileOpen && !isDesktop && (
        <div
          onClick={() => setMobileOpen(false)}
          style={{
            position: 'fixed', inset: 0, background: 'var(--bg-overlay)', zIndex: 30,
          }}
        />
      )}

      {/* Sidebar */}
      {sidebarVisible && (
        <aside style={{
          position: 'fixed', top: 0, left: 0, height: '100vh', width: sidebarWidth,
          background: 'var(--sidebar-bg)', borderRight: '1px solid var(--border-primary)',
          zIndex: 40, display: 'flex', flexDirection: 'column',
          overflowY: 'auto', overflowX: 'hidden',
          transition: 'width 0.2s ease',
        }}>
          {/* Logo */}
          <div style={{
            padding: isCollapsed ? '24px 12px 20px' : '24px 24px 20px',
            borderBottom: '1px solid var(--border-primary)',
            display: 'flex', alignItems: 'center', justifyContent: isCollapsed ? 'center' : 'flex-start',
          }}>
            <div style={{ display: 'flex', alignItems: 'center', gap: 10 }}>
              <div style={{
                width: 36, height: 36, borderRadius: 10, flexShrink: 0,
                background: 'linear-gradient(135deg, #3b82f6, #06b6d4)',
                display: 'flex', alignItems: 'center', justifyContent: 'center',
                boxShadow: '0 4px 12px rgba(59,130,246,0.3)',
              }}>
                <Server size={18} color="white" />
              </div>
              {!isCollapsed && (
                <div>
                  <span style={{
                    fontSize: 22, fontWeight: 700,
                    background: 'linear-gradient(90deg, #60a5fa, #22d3ee)',
                    WebkitBackgroundClip: 'text', WebkitTextFillColor: 'transparent',
                  }}>
                    Velour
                  </span>
                  <p style={{ fontSize: 11, color: 'var(--text-muted)', marginTop: 2 }}>
                    Server Management Panel
                  </p>
                </div>
              )}
            </div>
          </div>

          {/* Navigation */}
          <nav style={{ flex: 1, padding: isCollapsed ? '16px 10px' : '16px', display: 'flex', flexDirection: 'column', gap: 4 }}>
            {navItems.map(({ to, icon: Icon, label }) => (
              <NavLink
                key={to}
                to={to}
                onClick={() => setMobileOpen(false)}
                style={({ isActive }) => linkStyle(isActive)}
                title={isCollapsed ? label : undefined}
              >
                <Icon size={20} />
                {!isCollapsed && <span>{label}</span>}
              </NavLink>
            ))}
          </nav>

          {/* System stats in sidebar - mobile only */}
          {!isDesktop && system && <SidebarStats system={system} />}

          {/* Support & Links */}
          {!isCollapsed && (
            <div style={{ padding: '8px 16px', display: 'flex', flexDirection: 'column', gap: 6 }}>
              <a
                href="https://github.com/ozsari/velour"
                target="_blank"
                rel="noopener noreferrer"
                style={{
                  display: 'flex', alignItems: 'center', gap: 10, width: '100%',
                  padding: '9px 14px', borderRadius: 10, textDecoration: 'none',
                  background: 'rgba(255,255,255,0.03)',
                  border: '1px solid var(--border-primary)',
                  color: 'var(--text-tertiary)', fontSize: 13, fontWeight: 500,
                  transition: 'all 0.15s',
                }}
                onMouseEnter={(e) => {
                  e.currentTarget.style.background = 'rgba(255,255,255,0.06)';
                  e.currentTarget.style.color = 'var(--text-secondary)';
                }}
                onMouseLeave={(e) => {
                  e.currentTarget.style.background = 'rgba(255,255,255,0.03)';
                  e.currentTarget.style.color = 'var(--text-tertiary)';
                }}
              >
                <Star size={15} />
                <span style={{ flex: 1 }}>Star on GitHub</span>
                <ExternalLink size={12} style={{ opacity: 0.5 }} />
              </a>
              <a
                href="https://buymeacoffee.com/velour"
                target="_blank"
                rel="noopener noreferrer"
                style={{
                  display: 'flex', alignItems: 'center', gap: 10, width: '100%',
                  padding: '9px 14px', borderRadius: 10, textDecoration: 'none',
                  background: 'linear-gradient(135deg, rgba(251,191,36,0.1), rgba(245,158,11,0.05))',
                  border: '1px solid rgba(251,191,36,0.2)',
                  color: '#fbbf24', fontSize: 13, fontWeight: 600,
                  transition: 'all 0.15s',
                }}
                onMouseEnter={(e) => {
                  e.currentTarget.style.background = 'linear-gradient(135deg, rgba(251,191,36,0.2), rgba(245,158,11,0.1))';
                  e.currentTarget.style.borderColor = 'rgba(251,191,36,0.4)';
                }}
                onMouseLeave={(e) => {
                  e.currentTarget.style.background = 'linear-gradient(135deg, rgba(251,191,36,0.1), rgba(245,158,11,0.05))';
                  e.currentTarget.style.borderColor = 'rgba(251,191,36,0.2)';
                }}
              >
                <Heart size={15} style={{ fill: '#fbbf24', color: '#fbbf24' }} />
                <span>Buy Me a Coffee</span>
              </a>
            </div>
          )}

          {/* Collapsed: icon-only links */}
          {isCollapsed && (
            <div style={{ padding: '8px 10px', display: 'flex', flexDirection: 'column', gap: 6, alignItems: 'center' }}>
              <a
                href="https://github.com/ozsari/velour"
                target="_blank"
                rel="noopener noreferrer"
                title="Star on GitHub"
                style={{
                  width: 36, height: 36, borderRadius: 8, display: 'flex',
                  alignItems: 'center', justifyContent: 'center', textDecoration: 'none',
                  background: 'rgba(255,255,255,0.03)', border: '1px solid var(--border-primary)',
                  color: 'var(--text-tertiary)', transition: 'all 0.15s',
                }}
              >
                <Star size={15} />
              </a>
              <a
                href="https://buymeacoffee.com/velour"
                target="_blank"
                rel="noopener noreferrer"
                title="Buy Me a Coffee"
                style={{
                  width: 36, height: 36, borderRadius: 8, display: 'flex',
                  alignItems: 'center', justifyContent: 'center', textDecoration: 'none',
                  background: 'rgba(251,191,36,0.08)', border: '1px solid rgba(251,191,36,0.2)',
                  color: '#fbbf24', transition: 'all 0.15s',
                }}
              >
                <Heart size={15} style={{ fill: '#fbbf24', color: '#fbbf24' }} />
              </a>
            </div>
          )}

          {/* Bottom: Logout + Collapse toggle + Version */}
          <div style={{ padding: isCollapsed ? '12px 10px' : '16px', borderTop: '1px solid var(--border-primary)' }}>
            <div style={{ display: 'flex', flexDirection: isCollapsed ? 'column' : 'row', gap: 6, alignItems: 'center' }}>
              <button
                onClick={onLogout}
                title="Logout"
                style={{
                  display: 'flex', alignItems: 'center', justifyContent: isCollapsed ? 'center' : 'flex-start',
                  gap: 12, flex: isCollapsed ? undefined : 1,
                  padding: isCollapsed ? '10px' : '12px 16px', borderRadius: 10, border: 'none',
                  background: 'transparent', color: 'var(--text-tertiary)', fontSize: 14,
                  fontWeight: 500, cursor: 'pointer', transition: 'all 0.15s',
                  width: isCollapsed ? 40 : undefined, height: isCollapsed ? 40 : undefined,
                }}
                onMouseEnter={(e) => {
                  e.currentTarget.style.background = 'rgba(239,68,68,0.1)';
                  e.currentTarget.style.color = 'var(--accent-red)';
                }}
                onMouseLeave={(e) => {
                  e.currentTarget.style.background = 'transparent';
                  e.currentTarget.style.color = 'var(--text-tertiary)';
                }}
              >
                <LogOut size={20} />
                {!isCollapsed && <span>Logout</span>}
              </button>

              {/* Mobile theme toggle */}
              {!isDesktop && (
                <button
                  onClick={toggleTheme}
                  style={{
                    display: 'flex', alignItems: 'center', justifyContent: 'center',
                    width: 44, borderRadius: 10, border: 'none',
                    background: 'transparent', color: 'var(--accent-yellow)',
                    cursor: 'pointer', transition: 'all 0.15s',
                  }}
                  title={theme === 'dark' ? 'Light Mode' : 'Dark Mode'}
                >
                  {theme === 'dark' ? <Sun size={18} /> : <Moon size={18} />}
                </button>
              )}

              {/* Desktop collapse toggle */}
              {isDesktop && (
                <button
                  onClick={onToggleCollapse}
                  title={collapsed ? 'Expand sidebar' : 'Collapse sidebar'}
                  style={{
                    display: 'flex', alignItems: 'center', justifyContent: 'center',
                    width: 36, height: 36, borderRadius: 8, border: '1px solid var(--border-primary)',
                    background: 'var(--bg-tertiary)', color: 'var(--text-tertiary)',
                    cursor: 'pointer', transition: 'all 0.15s', flexShrink: 0,
                  }}
                  onMouseEnter={(e) => {
                    e.currentTarget.style.color = 'var(--text-primary)';
                    e.currentTarget.style.background = 'var(--bg-secondary)';
                  }}
                  onMouseLeave={(e) => {
                    e.currentTarget.style.color = 'var(--text-tertiary)';
                    e.currentTarget.style.background = 'var(--bg-tertiary)';
                  }}
                >
                  {collapsed ? <ChevronsRight size={16} /> : <ChevronsLeft size={16} />}
                </button>
              )}
            </div>
            {!isCollapsed && (
              <div style={{
                textAlign: 'center', marginTop: 8,
                fontSize: 11, color: 'var(--text-faint)', fontFamily: 'monospace',
              }}>
                Velour v0.1.0
              </div>
            )}
          </div>
        </aside>
      )}
    </>
  );
}
