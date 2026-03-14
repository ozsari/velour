import { Cpu, MemoryStick, HardDrive, ArrowUp, ArrowDown, Sun, Moon } from 'lucide-react';
import { useSystemInfo, formatBytes, formatSpeed } from '../hooks/useSystemInfo';
import { useTheme } from '../hooks/useTheme';

function Indicator({ icon, value, percent, color }: {
  icon: React.ReactNode; value: string; percent?: number; color: string;
}) {
  const barColor = percent && percent > 90 ? 'var(--accent-red)' : percent && percent > 70 ? 'var(--accent-yellow)' : color;
  return (
    <div style={{
      display: 'flex', alignItems: 'center', gap: 6,
      padding: '5px 10px', borderRadius: 8,
      background: 'var(--topbar-bg)', border: '1px solid var(--border-primary)',
      backdropFilter: 'blur(8px)',
    }}>
      <div style={{ color, display: 'flex', flexShrink: 0 }}>{icon}</div>
      <span style={{ fontSize: 11, fontWeight: 600, color: 'var(--text-secondary)', fontFamily: 'monospace', whiteSpace: 'nowrap' }}>
        {value}
      </span>
      {percent !== undefined && (
        <div style={{ width: 28, height: 3, background: 'var(--bg-tertiary)', borderRadius: 2, overflow: 'hidden' }}>
          <div style={{
            height: '100%', borderRadius: 2, background: barColor,
            width: `${Math.min(percent, 100)}%`, transition: 'width 0.5s ease',
          }} />
        </div>
      )}
    </div>
  );
}

export default function TopBar() {
  const system = useSystemInfo();
  const { theme, toggleTheme } = useTheme();

  if (!system) return null;

  return (
    <div style={{
      position: 'fixed', top: 10, right: 16, zIndex: 20,
      display: 'flex', gap: 5, alignItems: 'center',
    }}>
      <Indicator icon={<Cpu size={12} />} value={`${system.cpu.usage.toFixed(0)}%`} percent={system.cpu.usage} color="var(--accent-blue-light)" />
      <Indicator icon={<MemoryStick size={12} />} value={formatBytes(system.memory.used)} percent={system.memory.usage_percent} color="var(--accent-purple)" />
      <Indicator icon={<HardDrive size={12} />} value={`${system.disk.usage_percent.toFixed(0)}%`} percent={system.disk.usage_percent} color="var(--accent-yellow)" />

      {/* Network Speed */}
      <div style={{
        display: 'flex', alignItems: 'center', gap: 8,
        padding: '5px 10px', borderRadius: 8,
        background: 'var(--topbar-bg)', border: '1px solid var(--border-primary)',
        backdropFilter: 'blur(8px)',
      }}>
        <div style={{ display: 'flex', alignItems: 'center', gap: 3 }}>
          <ArrowUp size={10} style={{ color: 'var(--accent-red)' }} />
          <span style={{ fontSize: 11, fontWeight: 600, color: 'var(--text-secondary)', fontFamily: 'monospace', whiteSpace: 'nowrap' }}>
            {formatSpeed(system.network_speed.upload)}
          </span>
        </div>
        <div style={{ width: 1, height: 12, background: 'var(--border-primary)' }} />
        <div style={{ display: 'flex', alignItems: 'center', gap: 3 }}>
          <ArrowDown size={10} style={{ color: 'var(--accent-green)' }} />
          <span style={{ fontSize: 11, fontWeight: 600, color: 'var(--text-secondary)', fontFamily: 'monospace', whiteSpace: 'nowrap' }}>
            {formatSpeed(system.network_speed.download)}
          </span>
        </div>
      </div>

      {/* Theme Toggle */}
      <button
        onClick={toggleTheme}
        style={{
          display: 'flex', alignItems: 'center', justifyContent: 'center',
          width: 32, height: 32, borderRadius: 8, border: '1px solid var(--border-primary)',
          background: 'var(--topbar-bg)', backdropFilter: 'blur(8px)',
          color: 'var(--accent-yellow)', cursor: 'pointer',
          transition: 'all 0.15s',
        }}
        title={theme === 'dark' ? 'Switch to Light Mode' : 'Switch to Dark Mode'}
      >
        {theme === 'dark' ? <Sun size={14} /> : <Moon size={14} />}
      </button>
    </div>
  );
}
