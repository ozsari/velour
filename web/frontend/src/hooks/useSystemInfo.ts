import { useEffect, useState } from 'react';
import { api, type SystemInfo } from '../lib/api';

export interface SystemInfoWithSpeed extends SystemInfo {
  network_speed: {
    upload: number;   // bytes per second
    download: number; // bytes per second
  };
}

let cachedSystem: SystemInfoWithSpeed | null = null;
let prevBytes: { sent: number; recv: number; time: number } | null = null;
let listeners: Array<(s: SystemInfoWithSpeed | null) => void> = [];
let intervalId: ReturnType<typeof setInterval> | null = null;

async function fetchSystem() {
  try {
    const raw = await api.systemInfo();
    const now = Date.now();

    let uploadSpeed = 0;
    let downloadSpeed = 0;

    if (prevBytes) {
      const elapsed = (now - prevBytes.time) / 1000; // seconds
      if (elapsed > 0) {
        const sentDelta = raw.network.bytes_sent - prevBytes.sent;
        const recvDelta = raw.network.bytes_recv - prevBytes.recv;
        // Only calculate if delta is positive (handles counter resets)
        uploadSpeed = sentDelta > 0 ? sentDelta / elapsed : 0;
        downloadSpeed = recvDelta > 0 ? recvDelta / elapsed : 0;
      }
    }

    prevBytes = {
      sent: raw.network.bytes_sent,
      recv: raw.network.bytes_recv,
      time: now,
    };

    cachedSystem = {
      ...raw,
      network_speed: {
        upload: uploadSpeed,
        download: downloadSpeed,
      },
    };

    listeners.forEach((fn) => fn(cachedSystem));
  } catch {
    // keep last known
  }
}

function subscribe(fn: (s: SystemInfoWithSpeed | null) => void) {
  listeners.push(fn);
  if (listeners.length === 1) {
    fetchSystem();
    intervalId = setInterval(fetchSystem, 5000);
  }
  return () => {
    listeners = listeners.filter((l) => l !== fn);
    if (listeners.length === 0 && intervalId) {
      clearInterval(intervalId);
      intervalId = null;
    }
  };
}

export function useSystemInfo() {
  const [system, setSystem] = useState<SystemInfoWithSpeed | null>(cachedSystem);

  useEffect(() => {
    const unsub = subscribe(setSystem);
    return unsub;
  }, []);

  return system;
}

export function formatBytes(bytes: number): string {
  if (bytes === 0) return '0 B';
  const k = 1024;
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return `${(bytes / Math.pow(k, i)).toFixed(1)} ${sizes[i]}`;
}

export function formatSpeed(bytesPerSec: number): string {
  if (bytesPerSec < 1) return '0 B/s';
  const k = 1024;
  const sizes = ['B/s', 'KB/s', 'MB/s', 'GB/s'];
  const i = Math.floor(Math.log(bytesPerSec) / Math.log(k));
  return `${(bytesPerSec / Math.pow(k, i)).toFixed(1)} ${sizes[i]}`;
}
