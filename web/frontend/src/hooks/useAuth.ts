import { useState, useEffect, useCallback } from 'react';
import { api } from '../lib/api';

export function useAuth() {
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const [isLoading, setIsLoading] = useState(true);
  const [needsSetup, setNeedsSetup] = useState(false);
  const [user, setUser] = useState<any>(null);

  useEffect(() => {
    checkAuth();
  }, []);

  const checkAuth = async () => {
    try {
      const { needs_setup } = await api.setupStatus();
      if (needs_setup) {
        localStorage.removeItem('velour_token');
        localStorage.removeItem('velour_user');
        setNeedsSetup(true);
        setIsLoading(false);
        return;
      }

      const token = localStorage.getItem('velour_token');
      if (token) {
        try {
          await api.systemInfo();
          setIsAuthenticated(true);
          setUser(JSON.parse(localStorage.getItem('velour_user') || '{}'));
        } catch {
          localStorage.removeItem('velour_token');
          localStorage.removeItem('velour_user');
        }
      }
    } catch {
      // API unreachable
    } finally {
      setIsLoading(false);
    }
  };

  const login = useCallback(async (username: string, password: string) => {
    const { token, user } = await api.login(username, password);
    localStorage.setItem('velour_token', token);
    localStorage.setItem('velour_user', JSON.stringify(user));
    setIsAuthenticated(true);
    setNeedsSetup(false);
    setUser(user);
  }, []);

  const setup = useCallback(async (username: string, password: string) => {
    await api.setup(username, password);
    await login(username, password);
  }, [login]);

  const logout = useCallback(() => {
    localStorage.removeItem('velour_token');
    localStorage.removeItem('velour_user');
    setIsAuthenticated(false);
    setUser(null);
  }, []);

  return { isAuthenticated, isLoading, needsSetup, user, login, setup, logout };
}
