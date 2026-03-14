import { useEffect, useState } from 'react';
import { Plus, Trash2, Zap, Play, Pause, Clock, CheckCircle, XCircle, ChevronDown, Settings2 } from 'lucide-react';
import { api, type AutomationRule, type ActionTemplate, type CreateRuleRequest, type Service } from '../lib/api';

const triggerTypes = [
  { value: 'torrent_done', label: 'Torrent Finished', desc: 'Fires when a download completes' },
  { value: 'service_start', label: 'Service Started', desc: 'Fires when a service starts' },
  { value: 'service_stop', label: 'Service Stopped', desc: 'Fires when a service stops' },
];

const torrentClients = ['qbittorrent', 'deluge', 'rtorrent', 'transmission'];

const actionTypes = [
  { value: 'exec_in_service', label: 'Run Command in Service' },
  { value: 'restart_service', label: 'Restart a Service' },
];

export default function Automation() {
  const [rules, setRules] = useState<AutomationRule[]>([]);
  const [_templates, setTemplates] = useState<ActionTemplate[]>([]);
  const [services, setServices] = useState<Service[]>([]);
  const [loading, setLoading] = useState(true);
  const [showCreate, setShowCreate] = useState(false);
  const [expandedLog, setExpandedLog] = useState<string | null>(null);

  // Form state
  const [name, setName] = useState('');
  const [triggerType, setTriggerType] = useState('torrent_done');
  const [triggerService, setTriggerService] = useState('qbittorrent');
  const [actionType, setActionType] = useState('exec_in_service');
  const [actionService, setActionService] = useState('filebot');
  const [_selectedTemplate, setSelectedTemplate] = useState('filebot_amc');
  const [customCommand, setCustomCommand] = useState('');
  const [customArgs, setCustomArgs] = useState('');
  const [advancedMode, setAdvancedMode] = useState(false);
  const [selectedPreset, setSelectedPreset] = useState<string | null>(null);

  const fetchAll = async () => {
    try {
      const [r, t, s] = await Promise.all([
        api.listRules(),
        api.listTemplates(),
        api.listServices(),
      ]);
      setRules(r || []);
      setTemplates(t || []);
      setServices(s || []);
    } catch {} finally {
      setLoading(false);
    }
  };

  useEffect(() => { fetchAll(); }, []);

  const handleCreate = async () => {
    let rule: CreateRuleRequest;

    if (!advancedMode && selectedPreset) {
      // Simple mode: use pre-defined preset
      const presetMap: Record<string, { template: string; serviceId: string; command: string; args: string[] }> = {
        filebot_amc: {
          template: 'filebot_amc',
          serviceId: 'filebot',
          command: 'filebot',
          args: ['-script', 'fn:amc', '--output', '/data/media', '--action', 'hardlink', '-non-strict', '--def', 'ut_dir={torrent_path}', 'ut_title={torrent_name}', 'ut_label={torrent_category}'],
        },
        unpackerr: {
          template: 'unpackerr',
          serviceId: 'unpackerr',
          command: 'unpackerr',
          args: ['--path', '{torrent_path}'],
        },
        notify_sonarr: {
          template: 'notify_sonarr',
          serviceId: 'sonarr',
          command: 'curl',
          args: ['-X', 'POST', 'http://localhost:8989/api/v3/command', '-H', 'Content-Type: application/json', '-d', '{"name":"DownloadedEpisodesScan"}'],
        },
        rclone_sync: {
          template: 'rclone_sync',
          serviceId: 'rclone',
          command: 'rclone',
          args: ['sync', '/data/media', 'remote:', '--progress', '--transfers', '4'],
        },
        rclone_move: {
          template: 'rclone_move',
          serviceId: 'rclone',
          command: 'rclone',
          args: ['move', '{torrent_path}', 'remote:uploads/', '--progress'],
        },
      };
      const preset = presetMap[selectedPreset];
      rule = {
        name: name || selectedPreset + ' on ' + triggerService,
        trigger: { type: 'torrent_done', service_id: triggerService },
        action: {
          type: 'exec_in_service',
          service_id: preset.serviceId,
          command: preset.command,
          args: preset.args,
          template: preset.template,
        },
      };
    } else {
      // Advanced mode: custom config
      rule = {
        name: name || 'New Rule',
        trigger: {
          type: triggerType as any,
          service_id: triggerService,
        },
        action: {
          type: actionType as any,
          service_id: actionService,
          command: customCommand,
          args: customArgs.split(' ').filter(Boolean),
        },
      };
      if (actionType === 'restart_service') {
        rule.action = { type: 'restart_service', service_id: actionService };
      }
    }

    try {
      await api.createRule(rule);
      setShowCreate(false);
      resetForm();
      setAdvancedMode(false);
      setSelectedPreset(null);
      fetchAll();
    } catch {}
  };

  const resetForm = () => {
    setName('');
    setTriggerType('torrent_done');
    setTriggerService('qbittorrent');
    setActionType('exec_in_service');
    setActionService('filebot');
    setSelectedTemplate('filebot_amc');
    setCustomCommand('');
    setCustomArgs('');
  };

  const handleToggle = async (id: string, enabled: boolean) => {
    try {
      await api.toggleRule(id, !enabled);
      fetchAll();
    } catch {}
  };

  const handleDelete = async (id: string) => {
    try {
      await api.deleteRule(id);
      fetchAll();
    } catch {}
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
      {/* Header */}
      <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', marginBottom: 24 }}>
        <div>
          <h2 style={{ fontSize: 24, fontWeight: 700, color: 'var(--text-primary)', marginBottom: 4, display: 'flex', alignItems: 'center', gap: 10 }}>
            <Zap size={24} style={{ color: 'var(--accent-yellow)' }} />
            Automation
          </h2>
          <p style={{ fontSize: 14, color: 'var(--text-tertiary)' }}>Post-processing rules for your media pipeline</p>
        </div>
        <button
          onClick={() => setShowCreate(true)}
          style={{
            display: 'flex', alignItems: 'center', gap: 8, padding: '10px 20px', borderRadius: 10,
            background: 'linear-gradient(135deg, #3b82f6, #2563eb)', color: '#fff', border: 'none',
            fontSize: 14, fontWeight: 600, cursor: 'pointer',
            boxShadow: '0 4px 12px rgba(59,130,246,0.3)',
          }}
        >
          <Plus size={18} /> New Rule
        </button>
      </div>

      {/* How it works */}
      <div style={{
        background: 'rgba(59,130,246,0.05)', border: '1px solid rgba(59,130,246,0.15)',
        borderRadius: 12, padding: 20, marginBottom: 24,
      }}>
        <h3 style={{ fontSize: 14, fontWeight: 600, color: 'var(--accent-blue-light)', marginBottom: 8 }}>How it works</h3>
        <div style={{ display: 'flex', gap: 32, flexWrap: 'wrap' }}>
          {[
            { step: '1', title: 'Trigger', desc: 'Velour polls your torrent client API for completed downloads' },
            { step: '2', title: 'Match', desc: 'When a new completion is detected, the rule fires' },
            { step: '3', title: 'Execute', desc: 'Velour runs the action command inside the target service' },
          ].map(s => (
            <div key={s.step} style={{ display: 'flex', gap: 10, flex: 1, minWidth: 200 }}>
              <div style={{
                width: 28, height: 28, borderRadius: '50%', flexShrink: 0,
                background: 'rgba(59,130,246,0.15)', color: 'var(--accent-blue-light)',
                display: 'flex', alignItems: 'center', justifyContent: 'center',
                fontSize: 13, fontWeight: 700,
              }}>{s.step}</div>
              <div>
                <div style={{ fontSize: 13, fontWeight: 600, color: 'var(--text-secondary)' }}>{s.title}</div>
                <div style={{ fontSize: 12, color: 'var(--text-tertiary)', marginTop: 2 }}>{s.desc}</div>
              </div>
            </div>
          ))}
        </div>
      </div>

      {/* Create Rule Modal */}
      {showCreate && (
        <div style={{
          position: 'fixed', inset: 0, background: 'var(--bg-overlay)', zIndex: 100,
          display: 'flex', alignItems: 'center', justifyContent: 'center',
        }} onClick={() => setShowCreate(false)}>
          <div onClick={e => e.stopPropagation()} style={{
            background: 'var(--bg-secondary)', borderRadius: 16, padding: 32, width: '100%', maxWidth: 600,
            border: '1px solid var(--border-primary)', maxHeight: '90vh', overflowY: 'auto',
          }}>
            <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', marginBottom: 24 }}>
              <h3 style={{ fontSize: 20, fontWeight: 700, color: 'var(--text-primary)', display: 'flex', alignItems: 'center', gap: 8 }}>
                <Zap size={20} style={{ color: 'var(--accent-yellow)' }} /> Create Automation Rule
              </h3>
              <button
                onClick={() => { setAdvancedMode(!advancedMode); setSelectedPreset(null); }}
                style={{
                  display: 'flex', alignItems: 'center', gap: 6, padding: '6px 12px', borderRadius: 8,
                  border: '1px solid var(--border-primary)', fontSize: 12, cursor: 'pointer',
                  background: advancedMode ? 'rgba(168,85,247,0.1)' : 'transparent',
                  color: advancedMode ? 'var(--accent-purple)' : 'var(--text-muted)',
                }}
              >
                <Settings2 size={14} /> {advancedMode ? 'Simple Mode' : 'Advanced'}
              </button>
            </div>

            {!advancedMode ? (
              /* ── SIMPLE MODE: Pre-defined presets ── */
              <>
                <p style={{ fontSize: 13, color: 'var(--text-tertiary)', marginBottom: 20 }}>
                  Pick a ready-made automation. Just select your torrent client and go.
                </p>

                {/* Torrent client selector */}
                <div style={{ marginBottom: 20 }}>
                  <label style={{ display: 'block', fontSize: 13, fontWeight: 500, color: 'var(--text-tertiary)', marginBottom: 8 }}>Your Torrent Client</label>
                  <div style={{ display: 'flex', gap: 8, flexWrap: 'wrap' }}>
                    {torrentClients.map(c => (
                      <button key={c} onClick={() => setTriggerService(c)} style={{
                        padding: '8px 16px', borderRadius: 8, cursor: 'pointer', fontSize: 13, fontWeight: 500,
                        border: triggerService === c ? '2px solid var(--accent-blue-light)' : '1px solid var(--border-primary)',
                        background: triggerService === c ? 'rgba(59,130,246,0.1)' : 'var(--bg-primary)',
                        color: triggerService === c ? 'var(--accent-blue-light)' : 'var(--text-tertiary)',
                      }}>
                        {c.charAt(0).toUpperCase() + c.slice(1)}
                      </button>
                    ))}
                  </div>
                </div>

                {/* Pre-defined presets */}
                <div style={{ display: 'flex', flexDirection: 'column', gap: 10, marginBottom: 24 }}>
                  {[
                    {
                      id: 'filebot_amc',
                      name: 'FileBot AMC - Auto Organize Media',
                      desc: 'When a torrent finishes, FileBot creates hardlinks in your media library (Movies/TV Shows/Music). Hardlinks use no extra disk space and keep seeding intact. The most popular post-processing setup.',
                      trigger: 'Torrent Finished',
                      action: 'Run FileBot AMC script',
                      command: 'filebot -script fn:amc --output /data/media --action hardlink -non-strict',
                      color: '#4ade80',
                      borderColor: 'rgba(34,197,94,0.3)',
                      bgColor: 'rgba(34,197,94,0.05)',
                    },
                    {
                      id: 'unpackerr',
                      name: 'Unpackerr - Auto Extract Archives',
                      desc: 'Automatically extracts RAR/ZIP archives from completed downloads. Great for scene releases that come in compressed format.',
                      trigger: 'Torrent Finished',
                      action: 'Run Unpackerr extract',
                      command: 'unpackerr --path {torrent_path}',
                      color: '#60a5fa',
                      borderColor: 'rgba(59,130,246,0.3)',
                      bgColor: 'rgba(59,130,246,0.05)',
                    },
                    {
                      id: 'notify_sonarr',
                      name: 'Notify Sonarr/Radarr - Trigger Import',
                      desc: 'Tells Sonarr or Radarr to scan and import the completed download immediately, instead of waiting for the next periodic scan.',
                      trigger: 'Torrent Finished',
                      action: 'Restart Sonarr/Radarr import scan',
                      command: 'curl -X POST http://localhost:8989/api/v3/command -d \'{"name":"DownloadedEpisodesScan"}\'',
                      color: '#fbbf24',
                      borderColor: 'rgba(234,179,8,0.3)',
                      bgColor: 'rgba(234,179,8,0.05)',
                    },
                    {
                      id: 'rclone_sync',
                      name: 'Rclone Sync - Upload to Cloud',
                      desc: 'Sync your media library to a cloud remote (Google Drive, OneDrive, Dropbox). Runs after torrent completion or can be scheduled via Advanced mode.',
                      trigger: 'Torrent Finished',
                      action: 'Run Rclone sync',
                      command: 'rclone sync /data/media remote: --progress --transfers 4',
                      color: '#22d3ee',
                      borderColor: 'rgba(6,182,212,0.3)',
                      bgColor: 'rgba(6,182,212,0.05)',
                    },
                    {
                      id: 'rclone_move',
                      name: 'Rclone Move - Upload & Delete Local',
                      desc: 'Move completed downloads to cloud storage and free up local disk space. Files are removed locally after successful upload.',
                      trigger: 'Torrent Finished',
                      action: 'Run Rclone move',
                      command: 'rclone move {torrent_path} remote:uploads/ --progress',
                      color: '#c084fc',
                      borderColor: 'rgba(168,85,247,0.3)',
                      bgColor: 'rgba(168,85,247,0.05)',
                    },
                  ].map(preset => (
                    <button
                      key={preset.id}
                      onClick={() => {
                        setSelectedPreset(selectedPreset === preset.id ? null : preset.id);
                        setSelectedTemplate(preset.id);
                        setName(preset.name.split(' - ')[0] + ' on ' + triggerService);
                      }}
                      style={{
                        padding: 20, borderRadius: 12, cursor: 'pointer', textAlign: 'left',
                        border: selectedPreset === preset.id
                          ? `2px solid ${preset.color}`
                          : `1px solid ${preset.borderColor}`,
                        background: selectedPreset === preset.id
                          ? preset.bgColor
                          : 'var(--bg-primary)',
                        color: 'var(--text-primary)', transition: 'all 0.15s',
                      }}
                    >
                      <div style={{ display: 'flex', alignItems: 'center', gap: 8, marginBottom: 8 }}>
                        <div style={{
                          width: 8, height: 8, borderRadius: '50%',
                          background: selectedPreset === preset.id ? preset.color : 'var(--text-faint)',
                        }} />
                        <span style={{ fontSize: 15, fontWeight: 600, color: selectedPreset === preset.id ? preset.color : 'var(--text-secondary)' }}>
                          {preset.name}
                        </span>
                      </div>
                      <p style={{ fontSize: 12, color: 'var(--text-tertiary)', lineHeight: 1.5, marginBottom: 12 }}>
                        {preset.desc}
                      </p>
                      {selectedPreset === preset.id && (
                        <div style={{
                          padding: '10px 14px', borderRadius: 8, background: 'var(--bg-primary)',
                          border: '1px solid var(--border-primary)',
                        }}>
                          <div style={{ display: 'flex', alignItems: 'center', gap: 8, marginBottom: 6 }}>
                            <span style={{ fontSize: 11, fontWeight: 700, color: 'var(--accent-yellow)', padding: '1px 6px', background: 'rgba(234,179,8,0.1)', borderRadius: 3 }}>WHEN</span>
                            <span style={{ fontSize: 12, color: 'var(--text-secondary)' }}>{preset.trigger} in <span style={{ color: 'var(--accent-blue-light)', fontWeight: 600 }}>{triggerService}</span></span>
                          </div>
                          <div style={{ display: 'flex', alignItems: 'center', gap: 8, marginBottom: 8 }}>
                            <span style={{ fontSize: 11, fontWeight: 700, color: 'var(--accent-green)', padding: '1px 6px', background: 'rgba(34,197,94,0.1)', borderRadius: 3 }}>THEN</span>
                            <span style={{ fontSize: 12, color: 'var(--text-secondary)' }}>{preset.action}</span>
                          </div>
                          <div style={{ fontFamily: 'monospace', fontSize: 11, color: 'var(--accent-green)', wordBreak: 'break-all' }}>
                            $ {preset.command}
                          </div>
                        </div>
                      )}
                    </button>
                  ))}
                </div>
              </>
            ) : (
              /* ── ADVANCED MODE: Full control ── */
              <>
                {/* Rule Name */}
                <div style={{ marginBottom: 20 }}>
                  <label style={{ display: 'block', fontSize: 13, fontWeight: 500, color: 'var(--text-tertiary)', marginBottom: 6 }}>Rule Name</label>
                  <input
                    value={name}
                    onChange={e => setName(e.target.value)}
                    placeholder="e.g. FileBot AMC on torrent complete"
                    style={{
                      width: '100%', padding: '10px 14px', borderRadius: 8, border: '1px solid var(--border-primary)',
                      background: 'var(--bg-input)', color: 'var(--text-primary)', fontSize: 14, outline: 'none',
                      boxSizing: 'border-box',
                    }}
                  />
                </div>

                {/* Trigger Section */}
                <div style={{
                  background: 'rgba(234,179,8,0.05)', border: '1px solid rgba(234,179,8,0.15)',
                  borderRadius: 12, padding: 20, marginBottom: 20,
                }}>
                  <h4 style={{ fontSize: 14, fontWeight: 600, color: 'var(--accent-yellow)', marginBottom: 16 }}>WHEN (Trigger)</h4>

                  <div style={{ marginBottom: 14 }}>
                    <label style={{ display: 'block', fontSize: 12, color: 'var(--text-tertiary)', marginBottom: 4 }}>Event Type</label>
                    <select
                      value={triggerType}
                      onChange={e => setTriggerType(e.target.value)}
                      style={{
                        width: '100%', padding: '10px 14px', borderRadius: 8, border: '1px solid var(--border-primary)',
                        background: 'var(--bg-input)', color: 'var(--text-primary)', fontSize: 14, outline: 'none',
                      }}
                    >
                      {triggerTypes.map(t => (
                        <option key={t.value} value={t.value}>{t.label}</option>
                      ))}
                    </select>
                    <p style={{ fontSize: 11, color: 'var(--text-muted)', marginTop: 4 }}>
                      {triggerTypes.find(t => t.value === triggerType)?.desc}
                    </p>
                  </div>

                  <div>
                    <label style={{ display: 'block', fontSize: 12, color: 'var(--text-tertiary)', marginBottom: 4 }}>Source Service</label>
                    <select
                      value={triggerService}
                      onChange={e => setTriggerService(e.target.value)}
                      style={{
                        width: '100%', padding: '10px 14px', borderRadius: 8, border: '1px solid var(--border-primary)',
                        background: 'var(--bg-input)', color: 'var(--text-primary)', fontSize: 14, outline: 'none',
                      }}
                    >
                      {triggerType === 'torrent_done'
                        ? torrentClients.map(c => <option key={c} value={c}>{c.charAt(0).toUpperCase() + c.slice(1)}</option>)
                        : services.map(s => <option key={s.id} value={s.id}>{s.name}</option>)
                      }
                    </select>
                  </div>
                </div>

                {/* Action Section */}
                <div style={{
                  background: 'rgba(34,197,94,0.05)', border: '1px solid rgba(34,197,94,0.15)',
                  borderRadius: 12, padding: 20, marginBottom: 24,
                }}>
                  <h4 style={{ fontSize: 14, fontWeight: 600, color: 'var(--accent-green)', marginBottom: 16 }}>THEN (Action)</h4>

                  <div style={{ marginBottom: 14 }}>
                    <label style={{ display: 'block', fontSize: 12, color: 'var(--text-tertiary)', marginBottom: 4 }}>Action Type</label>
                    <select
                      value={actionType}
                      onChange={e => setActionType(e.target.value)}
                      style={{
                        width: '100%', padding: '10px 14px', borderRadius: 8, border: '1px solid var(--border-primary)',
                        background: 'var(--bg-input)', color: 'var(--text-primary)', fontSize: 14, outline: 'none',
                      }}
                    >
                      {actionTypes.map(a => (
                        <option key={a.value} value={a.value}>{a.label}</option>
                      ))}
                    </select>
                  </div>

                  {actionType === 'exec_in_service' && (
                    <>
                      <div style={{ marginBottom: 14 }}>
                        <label style={{ display: 'block', fontSize: 12, color: 'var(--text-tertiary)', marginBottom: 4 }}>Target Service</label>
                        <select
                          value={actionService}
                          onChange={e => setActionService(e.target.value)}
                          style={{
                            width: '100%', padding: '10px 14px', borderRadius: 8,
                            border: '1px solid var(--border-primary)',
                            background: 'var(--bg-input)', color: 'var(--text-primary)', fontSize: 14, outline: 'none',
                          }}
                        >
                          {services.map(s => <option key={s.id} value={s.id}>{s.name}</option>)}
                        </select>
                      </div>
                      <div style={{ marginBottom: 14 }}>
                        <label style={{ display: 'block', fontSize: 12, color: 'var(--text-tertiary)', marginBottom: 4 }}>Command</label>
                        <input
                          value={customCommand}
                          onChange={e => setCustomCommand(e.target.value)}
                          placeholder="e.g. filebot, /usr/bin/my-script"
                          style={{
                            width: '100%', padding: '10px 14px', borderRadius: 8,
                            border: '1px solid var(--border-primary)',
                            background: 'var(--bg-input)', color: 'var(--text-primary)', fontSize: 14, outline: 'none',
                            fontFamily: 'monospace', boxSizing: 'border-box',
                          }}
                        />
                      </div>
                      <div>
                        <label style={{ display: 'block', fontSize: 12, color: 'var(--text-tertiary)', marginBottom: 4 }}>
                          Arguments
                        </label>
                        <input
                          value={customArgs}
                          onChange={e => setCustomArgs(e.target.value)}
                          placeholder="e.g. -script fn:amc --output /data/media"
                          style={{
                            width: '100%', padding: '10px 14px', borderRadius: 8,
                            border: '1px solid var(--border-primary)',
                            background: 'var(--bg-input)', color: 'var(--text-primary)', fontSize: 14, outline: 'none',
                            fontFamily: 'monospace', boxSizing: 'border-box',
                          }}
                        />
                        <p style={{ fontSize: 11, color: 'var(--text-muted)', marginTop: 4 }}>
                          Variables: {'{torrent_path}'} {'{torrent_name}'} {'{torrent_hash}'} {'{torrent_category}'}
                        </p>
                      </div>
                    </>
                  )}

                  {actionType === 'restart_service' && (
                    <div>
                      <label style={{ display: 'block', fontSize: 12, color: 'var(--text-tertiary)', marginBottom: 4 }}>Service to Restart</label>
                      <select
                        value={actionService}
                        onChange={e => setActionService(e.target.value)}
                        style={{
                          width: '100%', padding: '10px 14px', borderRadius: 8,
                          border: '1px solid var(--border-primary)',
                          background: 'var(--bg-input)', color: 'var(--text-primary)', fontSize: 14, outline: 'none',
                        }}
                      >
                        {services.map(s => <option key={s.id} value={s.id}>{s.name}</option>)}
                      </select>
                    </div>
                  )}
                </div>

                {/* Summary */}
                <div style={{
                  background: 'var(--bg-primary)', borderRadius: 10, padding: 16, marginBottom: 24,
                  border: '1px solid var(--border-primary)',
                }}>
                  <div style={{ fontSize: 12, color: 'var(--text-muted)', marginBottom: 8 }}>Rule Summary</div>
                  <div style={{ fontSize: 14, color: 'var(--text-secondary)' }}>
                    <span style={{ color: 'var(--accent-yellow)', fontWeight: 600 }}>WHEN</span>{' '}
                    {triggerTypes.find(t => t.value === triggerType)?.label} in{' '}
                    <span style={{ color: 'var(--accent-blue-light)', fontWeight: 600 }}>{triggerService}</span>
                    {' → '}
                    <span style={{ color: 'var(--accent-green)', fontWeight: 600 }}>THEN</span>{' '}
                    {actionType === 'exec_in_service'
                      ? <>run command in <span style={{ color: 'var(--accent-blue-light)', fontWeight: 600 }}>{actionService}</span></>
                      : <>restart <span style={{ color: 'var(--accent-blue-light)', fontWeight: 600 }}>{actionService}</span></>
                    }
                  </div>
                </div>
              </>
            )}

            {/* Buttons */}
            <div style={{ display: 'flex', gap: 12, justifyContent: 'flex-end' }}>
              <button
                onClick={() => { setShowCreate(false); resetForm(); setAdvancedMode(false); setSelectedPreset(null); }}
                style={{
                  padding: '10px 20px', borderRadius: 8, border: '1px solid var(--border-primary)',
                  background: 'transparent', color: 'var(--text-tertiary)', fontSize: 14, cursor: 'pointer',
                }}
              >Cancel</button>
              <button
                onClick={handleCreate}
                disabled={!advancedMode && !selectedPreset}
                style={{
                  padding: '10px 24px', borderRadius: 8, border: 'none',
                  background: (!advancedMode && !selectedPreset)
                    ? 'var(--bg-tertiary)'
                    : 'linear-gradient(135deg, #3b82f6, #2563eb)',
                  color: (!advancedMode && !selectedPreset) ? 'var(--text-muted)' : '#fff',
                  fontSize: 14, fontWeight: 600,
                  cursor: (!advancedMode && !selectedPreset) ? 'not-allowed' : 'pointer',
                  boxShadow: (!advancedMode && !selectedPreset) ? 'none' : '0 4px 12px rgba(59,130,246,0.3)',
                }}
              >Create Rule</button>
            </div>
          </div>
        </div>
      )}

      {/* Rules List */}
      {rules.length === 0 ? (
        <div style={{ textAlign: 'center', padding: '80px 0' }}>
          <Zap size={48} style={{ margin: '0 auto 16px', color: 'var(--text-faint)' }} />
          <h3 style={{ fontSize: 18, fontWeight: 500, color: 'var(--text-secondary)', marginBottom: 4 }}>No automation rules</h3>
          <p style={{ fontSize: 14, color: 'var(--text-muted)' }}>Create your first rule to automate post-processing</p>
        </div>
      ) : (
        <div style={{ display: 'flex', flexDirection: 'column', gap: 12 }}>
          {rules.map(rule => (
            <div key={rule.id} style={{
              background: 'var(--bg-card)', border: '1px solid var(--border-primary)',
              borderRadius: 14, padding: 20, opacity: rule.enabled ? 1 : 0.5,
              transition: 'opacity 0.2s',
            }}>
              {/* Rule header */}
              <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', marginBottom: 12 }}>
                <div style={{ display: 'flex', alignItems: 'center', gap: 12 }}>
                  <Zap size={18} style={{ color: rule.enabled ? 'var(--accent-yellow)' : 'var(--text-faint)' }} />
                  <h3 style={{ fontSize: 16, fontWeight: 600, color: 'var(--text-primary)' }}>{rule.name}</h3>
                </div>
                <div style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
                  {/* Toggle */}
                  <button
                    onClick={() => handleToggle(rule.id, rule.enabled)}
                    style={{
                      display: 'flex', alignItems: 'center', gap: 6, padding: '6px 12px',
                      borderRadius: 8, border: 'none', fontSize: 13, cursor: 'pointer',
                      background: rule.enabled ? 'rgba(34,197,94,0.1)' : 'rgba(100,116,139,0.1)',
                      color: rule.enabled ? 'var(--accent-green)' : 'var(--text-tertiary)',
                    }}
                  >
                    {rule.enabled ? <><Pause size={14} /> Pause</> : <><Play size={14} /> Enable</>}
                  </button>
                  <button
                    onClick={() => handleDelete(rule.id)}
                    style={{
                      display: 'flex', alignItems: 'center', padding: '6px 10px', borderRadius: 8,
                      background: 'transparent', color: 'var(--text-muted)', border: 'none', cursor: 'pointer',
                    }}
                  ><Trash2 size={14} /></button>
                </div>
              </div>

              {/* Rule flow */}
              <div style={{
                display: 'flex', alignItems: 'center', gap: 12, flexWrap: 'wrap',
                padding: '12px 16px', background: 'var(--bg-primary)', borderRadius: 10,
                border: '1px solid var(--border-secondary)', marginBottom: 12,
              }}>
                <span style={{
                  fontSize: 11, fontWeight: 700, color: 'var(--accent-yellow)', padding: '2px 8px',
                  background: 'rgba(234,179,8,0.1)', borderRadius: 4,
                }}>WHEN</span>
                <span style={{ fontSize: 13, color: 'var(--text-secondary)' }}>
                  {rule.trigger.type === 'torrent_done' ? 'Torrent Finished' : rule.trigger.type} in
                </span>
                <span style={{
                  fontSize: 12, fontWeight: 600, color: 'var(--accent-blue-light)', padding: '2px 10px',
                  background: 'rgba(59,130,246,0.1)', borderRadius: 6,
                  border: '1px solid rgba(59,130,246,0.2)',
                }}>{rule.trigger.service_id}</span>
                <span style={{ fontSize: 16, color: 'var(--text-faint)' }}>→</span>
                <span style={{
                  fontSize: 11, fontWeight: 700, color: 'var(--accent-green)', padding: '2px 8px',
                  background: 'rgba(34,197,94,0.1)', borderRadius: 4,
                }}>THEN</span>
                <span style={{ fontSize: 13, color: 'var(--text-secondary)' }}>
                  {rule.action.type === 'exec_in_service' ? 'Run' : 'Restart'}
                </span>
                <span style={{
                  fontSize: 12, fontWeight: 600, color: 'var(--accent-purple)', padding: '2px 10px',
                  background: 'rgba(168,85,247,0.1)', borderRadius: 6,
                  border: '1px solid rgba(168,85,247,0.2)',
                }}>{rule.action.template || rule.action.command || rule.action.service_id}</span>
                {rule.action.service_id && rule.action.type === 'exec_in_service' && (
                  <>
                    <span style={{ fontSize: 13, color: 'var(--text-tertiary)' }}>in</span>
                    <span style={{
                      fontSize: 12, fontWeight: 600, color: 'var(--accent-blue-light)', padding: '2px 10px',
                      background: 'rgba(59,130,246,0.1)', borderRadius: 6,
                      border: '1px solid rgba(59,130,246,0.2)',
                    }}>{rule.action.service_id}</span>
                  </>
                )}
              </div>

              {/* Stats */}
              <div style={{ display: 'flex', alignItems: 'center', gap: 20, fontSize: 12 }}>
                <div style={{ display: 'flex', alignItems: 'center', gap: 4, color: 'var(--text-tertiary)' }}>
                  <Clock size={13} />
                  {rule.last_run_at
                    ? `Last: ${new Date(rule.last_run_at).toLocaleString()}`
                    : 'Never run'
                  }
                </div>
                {rule.last_run_ok !== undefined && (
                  <div style={{ display: 'flex', alignItems: 'center', gap: 4, color: rule.last_run_ok ? 'var(--accent-green)' : 'var(--accent-red)' }}>
                    {rule.last_run_ok ? <CheckCircle size={13} /> : <XCircle size={13} />}
                    {rule.last_run_ok ? 'Success' : 'Failed'}
                  </div>
                )}
                <div style={{ color: 'var(--text-muted)' }}>
                  {rule.run_count} run{rule.run_count !== 1 ? 's' : ''} total
                </div>
                {rule.last_run_log && (
                  <button
                    onClick={() => setExpandedLog(expandedLog === rule.id ? null : rule.id)}
                    style={{
                      display: 'flex', alignItems: 'center', gap: 4, padding: '2px 8px',
                      borderRadius: 4, border: 'none', background: 'rgba(100,116,139,0.1)',
                      color: 'var(--text-tertiary)', fontSize: 12, cursor: 'pointer',
                    }}
                  >
                    Log <ChevronDown size={12} style={{ transform: expandedLog === rule.id ? 'rotate(180deg)' : 'none', transition: 'transform 0.2s' }} />
                  </button>
                )}
              </div>

              {/* Log output */}
              {expandedLog === rule.id && rule.last_run_log && (
                <div style={{
                  marginTop: 12, padding: 12, borderRadius: 8, background: 'var(--bg-primary)',
                  border: '1px solid var(--border-secondary)', fontFamily: 'monospace', fontSize: 12,
                  color: 'var(--text-tertiary)', whiteSpace: 'pre-wrap', maxHeight: 200, overflowY: 'auto',
                }}>
                  {rule.last_run_log}
                </div>
              )}
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
