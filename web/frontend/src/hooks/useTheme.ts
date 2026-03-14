import { createContext, useContext, useEffect, useState } from 'react';

export type Theme = 'dark' | 'light';

interface ThemeContextValue {
  theme: Theme;
  toggleTheme: () => void;
}

export const ThemeContext = createContext<ThemeContextValue>({
  theme: 'dark',
  toggleTheme: () => {},
});

export function useTheme() {
  return useContext(ThemeContext);
}

export function useThemeProvider(): ThemeContextValue {
  const [theme, setTheme] = useState<Theme>(() => {
    const saved = localStorage.getItem('velour-theme');
    return (saved === 'light' ? 'light' : 'dark') as Theme;
  });

  useEffect(() => {
    localStorage.setItem('velour-theme', theme);
    document.documentElement.setAttribute('data-theme', theme);
  }, [theme]);

  const toggleTheme = () => setTheme(t => (t === 'dark' ? 'light' : 'dark'));

  return { theme, toggleTheme };
}
