import { useState, useEffect } from 'react';
import { Routes, Route, Navigate } from 'react-router-dom';
import { useAuth } from './hooks/useAuth';
import { ThemeContext, useThemeProvider } from './hooks/useTheme';
import Sidebar from './components/Sidebar';
import TopBar from './components/TopBar';
import Login from './pages/Login';
import Dashboard from './pages/Dashboard';
import Services from './pages/Services';
import Catalog from './pages/Catalog';
import SettingsPage from './pages/SettingsPage';
import Automation from './pages/Automation';

const SIDEBAR_EXPANDED = 256;
const SIDEBAR_COLLAPSED = 68;

export default function App() {
  const { isAuthenticated, isLoading, needsSetup, login, setup, logout } = useAuth();
  const themeValue = useThemeProvider();
  const [isDesktop, setIsDesktop] = useState(window.innerWidth >= 1024);
  const [sidebarCollapsed, setSidebarCollapsed] = useState(() => {
    return localStorage.getItem('velour_sidebar_collapsed') === 'true';
  });

  useEffect(() => {
    const handler = () => setIsDesktop(window.innerWidth >= 1024);
    window.addEventListener('resize', handler);
    return () => window.removeEventListener('resize', handler);
  }, []);

  const toggleSidebar = () => {
    setSidebarCollapsed(prev => {
      const next = !prev;
      localStorage.setItem('velour_sidebar_collapsed', String(next));
      return next;
    });
  };

  if (isLoading) {
    return (
      <div style={{
        minHeight: '100vh', display: 'flex', alignItems: 'center', justifyContent: 'center',
      }}>
        <div style={{ textAlign: 'center' }}>
          <div style={{
            width: 40, height: 40, border: '3px solid rgba(59,130,246,0.2)',
            borderTopColor: '#3b82f6', borderRadius: '50%',
            animation: 'spin 0.8s linear infinite', margin: '0 auto 16px',
          }} />
          <p style={{ color: '#94a3b8', fontSize: 14 }}>Loading Velour...</p>
        </div>
      </div>
    );
  }

  if (needsSetup || !isAuthenticated) {
    return <Login needsSetup={needsSetup} onLogin={login} onSetup={setup} />;
  }

  const sidebarWidth = isDesktop ? (sidebarCollapsed ? SIDEBAR_COLLAPSED : SIDEBAR_EXPANDED) : 0;

  return (
    <ThemeContext.Provider value={themeValue}>
      <div style={{ minHeight: '100vh' }}>
        <Sidebar
          onLogout={logout}
          collapsed={sidebarCollapsed}
          onToggleCollapse={toggleSidebar}
        />
        {isDesktop && <TopBar />}
        <main style={{
          marginLeft: sidebarWidth,
          padding: isDesktop ? '32px 32px 32px' : '68px 20px 20px',
          minHeight: '100vh',
          transition: 'margin-left 0.2s ease',
        }}>
          <Routes>
            <Route path="/" element={<Dashboard />} />
            <Route path="/services" element={<Services />} />
            <Route path="/catalog" element={<Catalog />} />
            <Route path="/automation" element={<Automation />} />
            <Route path="/settings" element={<SettingsPage />} />
            <Route path="*" element={<Navigate to="/" replace />} />
          </Routes>
        </main>
      </div>
    </ThemeContext.Provider>
  );
}
