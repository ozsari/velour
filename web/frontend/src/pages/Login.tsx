import { useState, useEffect } from 'react';
import { LogIn, UserPlus, Server, Shield, Zap, Eye, EyeOff, HelpCircle, Terminal, X } from 'lucide-react';

interface LoginProps {
  needsSetup: boolean;
  onLogin: (username: string, password: string) => Promise<void>;
  onSetup: (username: string, password: string) => Promise<void>;
}

export default function Login({ needsSetup, onLogin, onSetup }: LoginProps) {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const [showPassword, setShowPassword] = useState(false);
  const [mounted, setMounted] = useState(false);
  const [isDesktop, setIsDesktop] = useState(window.innerWidth >= 1024);
  const [showRecovery, setShowRecovery] = useState(false);

  useEffect(() => {
    setMounted(true);
    const handler = () => setIsDesktop(window.innerWidth >= 1024);
    window.addEventListener('resize', handler);
    return () => window.removeEventListener('resize', handler);
  }, []);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setLoading(true);

    try {
      if (needsSetup) {
        await onSetup(username, password);
      } else {
        await onLogin(username, password);
      }
    } catch (err: any) {
      setError(err.message || 'Something went wrong');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div
      style={{
        minHeight: '100vh',
        display: 'flex',
        position: 'relative',
        overflow: 'hidden',
        background: '#0f172a',
      }}
    >
      {/* Background orbs */}
      <div style={{
        position: 'absolute', width: 400, height: 400, borderRadius: '50%',
        background: '#2563eb', filter: 'blur(120px)', opacity: 0.15,
        top: -80, left: -80, animation: 'pulse 4s ease-in-out infinite',
      }} />
      <div style={{
        position: 'absolute', width: 350, height: 350, borderRadius: '50%',
        background: '#06b6d4', filter: 'blur(120px)', opacity: 0.12,
        bottom: -60, right: -60, animation: 'pulse 4s ease-in-out infinite 1s',
      }} />
      <div style={{
        position: 'absolute', width: 250, height: 250, borderRadius: '50%',
        background: '#7c3aed', filter: 'blur(120px)', opacity: 0.1,
        top: '50%', left: '50%', transform: 'translate(-50%, -50%)',
        animation: 'pulse 4s ease-in-out infinite 2s',
      }} />

      {/* Grid overlay */}
      <div style={{
        position: 'absolute', inset: 0, opacity: 0.03,
        backgroundImage: 'linear-gradient(rgba(255,255,255,.1) 1px, transparent 1px), linear-gradient(90deg, rgba(255,255,255,.1) 1px, transparent 1px)',
        backgroundSize: '60px 60px',
      }} />

      {/* Left panel - branding (desktop only) */}
      {isDesktop && <div style={{
        width: '50%', position: 'relative', display: 'flex', flexDirection: 'column',
        justifyContent: 'center', padding: '60px 80px',
      }}>
        <div style={{
          opacity: mounted ? 1 : 0, transform: mounted ? 'translateY(0)' : 'translateY(30px)',
          transition: 'all 0.7s ease-out',
        }}>
          <div style={{ display: 'flex', alignItems: 'center', gap: 14, marginBottom: 40 }}>
            <div style={{
              width: 48, height: 48, borderRadius: 14,
              background: 'linear-gradient(135deg, #3b82f6, #06b6d4)',
              display: 'flex', alignItems: 'center', justifyContent: 'center',
              boxShadow: '0 8px 24px rgba(59,130,246,0.3)',
            }}>
              <Server size={24} color="white" />
            </div>
            <span style={{
              fontSize: 30, fontWeight: 700,
              background: 'linear-gradient(90deg, #60a5fa, #22d3ee)',
              WebkitBackgroundClip: 'text', WebkitTextFillColor: 'transparent',
            }}>
              Velour
            </span>
          </div>

          <h2 style={{ fontSize: 48, fontWeight: 700, color: '#fff', lineHeight: 1.15, marginBottom: 24 }}>
            Your server,<br />
            <span style={{
              background: 'linear-gradient(90deg, #60a5fa, #22d3ee, #60a5fa)',
              WebkitBackgroundClip: 'text', WebkitTextFillColor: 'transparent',
            }}>
              your rules.
            </span>
          </h2>

          <p style={{ fontSize: 18, color: '#94a3b8', lineHeight: 1.7, marginBottom: 48, maxWidth: 420 }}>
            Manage all your media services from a single, beautiful dashboard.
            Install apps in one click, monitor your system in real-time.
          </p>

          <div style={{ display: 'flex', flexDirection: 'column', gap: 12 }}>
            {[
              { icon: Zap, text: 'One-click app installation', bg: 'rgba(245,158,11,0.1)', border: 'rgba(245,158,11,0.2)', color: '#fbbf24' },
              { icon: Shield, text: 'Secure authentication', bg: 'rgba(34,197,94,0.1)', border: 'rgba(34,197,94,0.2)', color: '#4ade80' },
              { icon: Server, text: 'Real-time system monitoring', bg: 'rgba(59,130,246,0.1)', border: 'rgba(59,130,246,0.2)', color: '#60a5fa' },
            ].map(({ icon: Icon, text, bg, border, color }) => (
              <div key={text} style={{
                display: 'inline-flex', alignItems: 'center', gap: 12,
                padding: '10px 18px', borderRadius: 50,
                background: bg, border: `1px solid ${border}`, color, width: 'fit-content',
              }}>
                <Icon size={16} />
                <span style={{ fontSize: 14, fontWeight: 500 }}>{text}</span>
              </div>
            ))}
          </div>
        </div>
      </div>}

      {/* Right panel - form */}
      <div style={{
        flex: 1, display: 'flex', alignItems: 'center', justifyContent: 'center',
        padding: '24px',
      }}>
        <div style={{
          width: '100%', maxWidth: 420,
          opacity: mounted ? 1 : 0, transform: mounted ? 'translateY(0)' : 'translateY(30px)',
          transition: 'all 0.7s ease-out 0.15s',
        }}>
          {/* Mobile logo */}
          {!isDesktop && <div style={{ display: 'flex', alignItems: 'center', gap: 12, marginBottom: 32 }}>
            <div style={{
              width: 42, height: 42, borderRadius: 12,
              background: 'linear-gradient(135deg, #3b82f6, #06b6d4)',
              display: 'flex', alignItems: 'center', justifyContent: 'center',
              boxShadow: '0 8px 24px rgba(59,130,246,0.3)',
            }}>
              <Server size={20} color="white" />
            </div>
            <span style={{
              fontSize: 24, fontWeight: 700,
              background: 'linear-gradient(90deg, #60a5fa, #22d3ee)',
              WebkitBackgroundClip: 'text', WebkitTextFillColor: 'transparent',
            }}>
              Velour
            </span>
          </div>}

          {/* Card with glow */}
          <div style={{ position: 'relative' }}>
            <div style={{
              position: 'absolute', inset: -2, borderRadius: 24,
              background: 'linear-gradient(135deg, rgba(59,130,246,0.15), rgba(6,182,212,0.15), rgba(124,58,237,0.15))',
              filter: 'blur(20px)', opacity: 0.6,
            }} />

            <div style={{
              position: 'relative',
              background: 'rgba(15, 23, 42, 0.9)',
              backdropFilter: 'blur(20px)',
              border: '1px solid rgba(51, 65, 85, 0.5)',
              borderRadius: 20,
              padding: 36,
              boxShadow: '0 25px 50px rgba(0,0,0,0.4)',
            }}>
              {/* Header */}
              <div style={{ marginBottom: 28 }}>
                <h3 style={{ fontSize: 26, fontWeight: 700, color: '#fff', marginBottom: 6 }}>
                  {needsSetup ? 'Welcome' : 'Welcome back'}
                </h3>
                <p style={{ fontSize: 14, color: '#94a3b8' }}>
                  {needsSetup
                    ? 'Set up your admin account to get started'
                    : 'Sign in to access your dashboard'}
                </p>
              </div>

              <form onSubmit={handleSubmit}>
                {/* Error */}
                {error && (
                  <div style={{
                    display: 'flex', alignItems: 'center', gap: 10,
                    background: 'rgba(239,68,68,0.08)', border: '1px solid rgba(239,68,68,0.2)',
                    borderRadius: 12, padding: '12px 16px', marginBottom: 20,
                    fontSize: 14, color: '#f87171',
                    animation: 'shake 0.5s ease-in-out',
                  }}>
                    <div style={{ width: 6, height: 6, borderRadius: '50%', background: '#f87171', flexShrink: 0 }} />
                    {error}
                  </div>
                )}

                {/* Username */}
                <div style={{ marginBottom: 20 }}>
                  <label style={{ display: 'block', fontSize: 13, fontWeight: 500, color: '#cbd5e1', marginBottom: 8 }}>
                    Username
                  </label>
                  <input
                    type="text"
                    value={username}
                    onChange={(e) => setUsername(e.target.value)}
                    placeholder="Enter username"
                    required
                    autoFocus
                    style={{
                      width: '100%', padding: '14px 16px', borderRadius: 12,
                      background: 'rgba(30, 41, 59, 0.5)', border: '1px solid rgba(51, 65, 85, 0.5)',
                      color: '#fff', fontSize: 15, outline: 'none',
                      transition: 'border-color 0.2s, box-shadow 0.2s',
                    }}
                    onFocus={(e) => {
                      e.target.style.borderColor = 'rgba(59,130,246,0.5)';
                      e.target.style.boxShadow = '0 0 0 3px rgba(59,130,246,0.1)';
                    }}
                    onBlur={(e) => {
                      e.target.style.borderColor = 'rgba(51, 65, 85, 0.5)';
                      e.target.style.boxShadow = 'none';
                    }}
                  />
                </div>

                {/* Password */}
                <div style={{ marginBottom: 24 }}>
                  <label style={{ display: 'block', fontSize: 13, fontWeight: 500, color: '#cbd5e1', marginBottom: 8 }}>
                    Password
                  </label>
                  <div style={{ position: 'relative' }}>
                    <input
                      type={showPassword ? 'text' : 'password'}
                      value={password}
                      onChange={(e) => setPassword(e.target.value)}
                      placeholder={needsSetup ? 'Min. 6 characters' : 'Enter password'}
                      required
                      minLength={6}
                      style={{
                        width: '100%', padding: '14px 48px 14px 16px', borderRadius: 12,
                        background: 'rgba(30, 41, 59, 0.5)', border: '1px solid rgba(51, 65, 85, 0.5)',
                        color: '#fff', fontSize: 15, outline: 'none',
                        transition: 'border-color 0.2s, box-shadow 0.2s',
                      }}
                      onFocus={(e) => {
                        e.target.style.borderColor = 'rgba(59,130,246,0.5)';
                        e.target.style.boxShadow = '0 0 0 3px rgba(59,130,246,0.1)';
                      }}
                      onBlur={(e) => {
                        e.target.style.borderColor = 'rgba(51, 65, 85, 0.5)';
                        e.target.style.boxShadow = 'none';
                      }}
                    />
                    <button
                      type="button"
                      onClick={() => setShowPassword(!showPassword)}
                      style={{
                        position: 'absolute', right: 12, top: '50%', transform: 'translateY(-50%)',
                        background: 'none', border: 'none', cursor: 'pointer',
                        color: '#64748b', padding: 4, display: 'flex',
                        transition: 'color 0.2s',
                      }}
                      onMouseEnter={(e) => (e.currentTarget.style.color = '#cbd5e1')}
                      onMouseLeave={(e) => (e.currentTarget.style.color = '#64748b')}
                    >
                      {showPassword ? <EyeOff size={18} /> : <Eye size={18} />}
                    </button>
                  </div>
                </div>

                {/* Submit */}
                <button
                  type="submit"
                  disabled={loading}
                  style={{
                    width: '100%', display: 'flex', alignItems: 'center', justifyContent: 'center',
                    gap: 10, padding: '14px 20px', borderRadius: 12, border: 'none',
                    background: 'linear-gradient(135deg, #2563eb, #3b82f6)',
                    color: '#fff', fontSize: 15, fontWeight: 600, cursor: 'pointer',
                    boxShadow: '0 8px 24px rgba(37,99,235,0.3)',
                    transition: 'all 0.2s', opacity: loading ? 0.6 : 1,
                  }}
                  onMouseEnter={(e) => {
                    if (!loading) {
                      e.currentTarget.style.background = 'linear-gradient(135deg, #3b82f6, #60a5fa)';
                      e.currentTarget.style.boxShadow = '0 8px 32px rgba(59,130,246,0.4)';
                      e.currentTarget.style.transform = 'translateY(-1px)';
                    }
                  }}
                  onMouseLeave={(e) => {
                    e.currentTarget.style.background = 'linear-gradient(135deg, #2563eb, #3b82f6)';
                    e.currentTarget.style.boxShadow = '0 8px 24px rgba(37,99,235,0.3)';
                    e.currentTarget.style.transform = 'translateY(0)';
                  }}
                >
                  {loading ? (
                    <div style={{
                      width: 20, height: 20, border: '2px solid rgba(255,255,255,0.3)',
                      borderTopColor: '#fff', borderRadius: '50%',
                      animation: 'spin 0.8s linear infinite',
                    }} />
                  ) : needsSetup ? (
                    <>
                      <UserPlus size={18} /> Create Account
                    </>
                  ) : (
                    <>
                      <LogIn size={18} /> Sign In
                    </>
                  )}
                </button>
              </form>

              {!needsSetup && (
                <button
                  type="button"
                  onClick={() => setShowRecovery(true)}
                  style={{
                    display: 'block', margin: '16px auto 0', padding: 0, border: 'none',
                    background: 'none', color: '#64748b', fontSize: 13, cursor: 'pointer',
                    transition: 'color 0.2s',
                  }}
                  onMouseEnter={e => (e.currentTarget.style.color = '#94a3b8')}
                  onMouseLeave={e => (e.currentTarget.style.color = '#64748b')}
                >
                  Forgot password?
                </button>
              )}

              {needsSetup && (
                <p style={{
                  fontSize: 12, color: '#64748b', marginTop: 20,
                  textAlign: 'center', lineHeight: 1.6,
                }}>
                  This will create the initial admin account.<br />
                  You can add more users later from Settings.
                </p>
              )}
            </div>
          </div>

          {/* Footer */}
          <p style={{ textAlign: 'center', fontSize: 12, color: '#475569', marginTop: 24 }}>
            Velour v0.1.0 &middot; Server Management Panel
          </p>
        </div>
      </div>

      {/* Password Recovery Modal */}
      {showRecovery && (
        <div style={{
          position: 'fixed', inset: 0, background: 'rgba(0,0,0,0.7)', zIndex: 100,
          display: 'flex', alignItems: 'center', justifyContent: 'center', padding: 24,
        }} onClick={() => setShowRecovery(false)}>
          <div onClick={e => e.stopPropagation()} style={{
            background: '#1e293b', borderRadius: 16, padding: 32, width: '100%', maxWidth: 460,
            border: '1px solid rgba(51,65,85,0.5)', position: 'relative',
          }}>
            <button
              onClick={() => setShowRecovery(false)}
              style={{
                position: 'absolute', top: 16, right: 16, padding: 4,
                background: 'none', border: 'none', color: '#64748b', cursor: 'pointer',
              }}
            ><X size={18} /></button>

            <div style={{ display: 'flex', alignItems: 'center', gap: 10, marginBottom: 20 }}>
              <div style={{
                width: 40, height: 40, borderRadius: 10, background: 'rgba(245,158,11,0.1)',
                display: 'flex', alignItems: 'center', justifyContent: 'center',
              }}>
                <HelpCircle size={20} style={{ color: '#fbbf24' }} />
              </div>
              <div>
                <h3 style={{ fontSize: 18, fontWeight: 700, color: '#fff' }}>Password Recovery</h3>
                <p style={{ fontSize: 12, color: '#94a3b8' }}>Reset your password via SSH</p>
              </div>
            </div>

            <p style={{ fontSize: 14, color: '#cbd5e1', lineHeight: 1.7, marginBottom: 20 }}>
              Since Velour runs on your own server, password recovery works through the command line.
              Connect to your server via SSH and run:
            </p>

            <div style={{
              background: '#0f172a', borderRadius: 10, padding: 16,
              border: '1px solid rgba(51,65,85,0.5)', marginBottom: 16,
            }}>
              <div style={{ display: 'flex', alignItems: 'center', gap: 8, marginBottom: 12 }}>
                <Terminal size={14} style={{ color: '#4ade80' }} />
                <span style={{ fontSize: 12, fontWeight: 600, color: '#4ade80' }}>Terminal</span>
              </div>
              <code style={{
                display: 'block', fontFamily: 'monospace', fontSize: 14, color: '#e2e8f0',
                wordBreak: 'break-all',
              }}>
                <span style={{ color: '#64748b' }}>$</span> sudo velour reset-password <span style={{ color: '#fbbf24' }}>username</span> <span style={{ color: '#60a5fa' }}>newpassword</span>
              </code>
            </div>

            <div style={{
              background: 'rgba(59,130,246,0.05)', border: '1px solid rgba(59,130,246,0.15)',
              borderRadius: 10, padding: 14,
            }}>
              <p style={{ fontSize: 12, color: '#94a3b8', lineHeight: 1.6 }}>
                <span style={{ fontWeight: 600, color: '#60a5fa' }}>Example:</span><br />
                <code style={{ fontFamily: 'monospace', color: '#cbd5e1' }}>
                  sudo velour reset-password admin MyNewPassword123
                </code>
              </p>
            </div>

            <p style={{ fontSize: 12, color: '#64748b', marginTop: 16, lineHeight: 1.5 }}>
              After resetting, come back here and sign in with your new password.
              The password must be at least 6 characters.
            </p>
          </div>
        </div>
      )}
    </div>
  );
}
