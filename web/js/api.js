// ─── HabitFlow API Client ───────────────────────────────────────

const API = {
  BASE: null,

  // Resolve API base once at runtime to support local dev and reverse proxy deployments.
  _resolveBase() {
    if (API.BASE) return API.BASE;

    const rawBase =
      (typeof window !== 'undefined' && typeof window.HABITFLOW_API_BASE === 'string' && window.HABITFLOW_API_BASE) ||
      (typeof localStorage !== 'undefined' && localStorage.getItem('hf_api_base')) ||
      '';

    let normalized = String(rawBase || '').trim();
    if (!normalized) {
      // If opened directly from file://, prefer local API server for easier debugging.
      normalized = location.protocol === 'file:' ? 'http://localhost:8080/api/v1' : '/api/v1';
    }

    API.BASE = normalized.replace(/\/+$/, '');
    return API.BASE;
  },

  _buildUrl(path) {
    const base = API._resolveBase();
    if (location.protocol === 'https:' && /^http:\/\//i.test(base)) {
      throw new Error('Mixed Content terdeteksi: frontend HTTPS tidak bisa memanggil API HTTP. Gunakan API HTTPS.');
    }
    const suffix = path.startsWith('/') ? path : '/' + path;
    return base + suffix;
  },

  // ── Token management ──
  getToken()    { return localStorage.getItem('hf_access_token'); },
  setTokens(access) {
    localStorage.setItem('hf_access_token', access);
    // Also store in IndexedDB for service worker access
    API._storeTokenIDB(access);
  },
  clearTokens() {
    localStorage.removeItem('hf_access_token');
    localStorage.removeItem('hf_user');
  },
  getUser()     { try { return JSON.parse(localStorage.getItem('hf_user')); } catch { return null; } },
  setUser(u)    { localStorage.setItem('hf_user', JSON.stringify(u)); },
  isLoggedIn()  { return !!API.getToken(); },

  // Store token in IndexedDB for service worker
  _storeTokenIDB(token) {
    try {
      const req = indexedDB.open('habitflow_auth', 1);
      req.onupgradeneeded = (e) => {
        const db = e.target.result;
        if (!db.objectStoreNames.contains('tokens')) db.createObjectStore('tokens');
      };
      req.onsuccess = (e) => {
        const db = e.target.result;
        const tx = db.transaction('tokens', 'readwrite');
        tx.objectStore('tokens').put(token, 'access_token');
      };
    } catch(e) { /* ignore */ }
  },

  // ── HTTP helpers ──
  async _fetch(path, opts = {}) {
    const headers = { 'Content-Type': 'application/json', ...opts.headers };
    const token = API.getToken();
    if (token) headers['Authorization'] = 'Bearer ' + token;

    const url = API._buildUrl(path);
    const body = opts.body != null && typeof opts.body === 'object' ? JSON.stringify(opts.body) : opts.body;

    let res;
    try {
      res = await fetch(url, { ...opts, body, headers, credentials: 'include' });
    } catch (err) {
      API._logNetworkError(path, url, err);
      throw API._networkError(path, url);
    }

    // Auto-refresh on 401
    if (res.status === 401) {
      const refreshed = await API._refresh();
      if (refreshed) {
        headers['Authorization'] = 'Bearer ' + API.getToken();
        try {
          res = await fetch(url, { ...opts, body, headers, credentials: 'include' });
        } catch (err) {
          API._logNetworkError(path, url, err);
          throw API._networkError(path, url);
        }
      } else {
        API.clearTokens();
        location.hash = '#/login';
        return null;
      }
    }

    const data = await API._readJsonSafe(res);
    if (!res.ok) throw new Error(data.message || 'Request failed');
    return data;
  },

  async _readJsonSafe(res) {
    const text = await res.text();
    if (!text) return {};
    try {
      return JSON.parse(text);
    } catch {
      return { message: text.slice(0, 250) || 'Respons server tidak valid' };
    }
  },

  _logNetworkError(path, url, err) {
    console.error('[HabitFlow API] Network request failed', {
      path,
      url,
      base: API._resolveBase(),
      origin: location.origin,
      protocol: location.protocol,
      online: navigator.onLine,
      error: err && err.message ? err.message : String(err),
    });
  },

  _networkError(path, url) {
    return new Error(
      'Koneksi ke server gagal. Periksa internet/CORS, pastikan API aktif, lalu coba lagi. ' +
      '(endpoint: ' + url + ')'
    );
  },

  async _refresh() {
    try {
      const res = await fetch(API._buildUrl('/auth/refresh'), {
        method: 'POST',
        credentials: 'include',
      });
      if (!res.ok) return false;
      const data = await res.json();
      if (data.success && data.data) {
        API.setTokens(data.data.access_token);
        return true;
      }
      return false;
    } catch { return false; }
  },

  // ── Auth ──
  async register(name, email, password) {
    const data = await API._fetch('/auth/register', {
      method: 'POST',
      body: JSON.stringify({ name, email, password }),
    });
    if (data?.data) {
      API.setTokens(data.data.access_token);
      API.setUser(data.data.user);
    }
    return data;
  },

  async login(email, password) {
    const data = await API._fetch('/auth/login', {
      method: 'POST',
      body: JSON.stringify({ email, password }),
    });
    if (data?.data) {
      API.setTokens(data.data.access_token);
      API.setUser(data.data.user);
    }
    return data;
  },

  async logout() {
    try {
      await API._fetch('/auth/logout', { method: 'POST' });
    } catch (e) {
      // Best effort revoke; still clear local session if request fails.
    }
    API.clearTokens();
    location.hash = '#/login';
  },

  // ── Habits ──
  getHabits()               { return API._fetch('/habits'); },
  getHabit(id)              { return API._fetch('/habits/' + id); },
  createHabit(body)         { return API._fetch('/habits', { method: 'POST', body: JSON.stringify(body) }); },
  updateHabit(id, body)     { return API._fetch('/habits/' + id, { method: 'PUT', body: JSON.stringify(body) }); },
  deleteHabit(id)           { return API._fetch('/habits/' + id, { method: 'DELETE' }); },

  // ── Check-in ──
  getToday()                { return API._fetch('/habits/today'); },
  checkin(habitId, note)    { return API._fetch('/habits/' + habitId + '/check', { method: 'POST', body: JSON.stringify({ note }) }); },
  undoCheckin(habitId)      { return API._fetch('/habits/' + habitId + '/check', { method: 'DELETE' }); },

  // ── Reports ──
  getWeeklyReport()         { return API._fetch('/reports/weekly'); },
  getDailyReport(date)      { return API._fetch('/reports/daily' + (date ? '?date=' + date : '')); },
  getScore(days)            { return API._fetch('/reports/score?days=' + (days || 7)); },
  getInsights()             { return API._fetch('/reports/insights'); },

  // ── Push ──
  getVAPIDKey()             { return API._fetch('/push/vapid-key'); },
  subscribePush(sub)        { return API._fetch('/push/subscribe', { method: 'POST', body: JSON.stringify(sub) }); },
  unsubscribePush(endpoint) { return API._fetch('/push/unsubscribe', { method: 'DELETE', body: JSON.stringify({ endpoint }) }); },
};
