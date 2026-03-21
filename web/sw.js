const CACHE_NAME = 'habitflow-v11';
const API_CACHE = 'habitflow-api-v7';
const ASSETS = [
  '/',
  '/index.html',
  '/css/tailwind.css',
  '/css/app.css',
  '/js/api.js',
  '/js/icons.js',
  '/js/app.js',
  '/js/sw-register.js',
  '/manifest.json',
  '/manifest.json?v=3',
  '/icons/favicon.svg',
  '/icons/favicon.svg?v=3',
  '/favicon.ico',
];

// Install — cache core assets
self.addEventListener('install', (event) => {
  event.waitUntil(
    caches.open(CACHE_NAME).then((cache) => cache.addAll(ASSETS))
  );
  self.skipWaiting();
});

// Activate — remove old caches
self.addEventListener('activate', (event) => {
  event.waitUntil(
    caches.keys().then((keys) =>
      Promise.all(
        keys
          .filter((key) => key !== CACHE_NAME && key !== API_CACHE)
          .map((key) => caches.delete(key))
      )
    )
  );
  self.clients.claim();
});

// Fetch — network first for static, never cache authenticated API responses.
self.addEventListener('fetch', (event) => {
  if (event.request.method !== 'GET') return;

  // Never intercept API traffic; let browser/network stack handle it directly.
  // This prevents SW from masking CORS/mixed-content diagnostics.
  if (event.request.url.includes('/api/')) {
    return;
  }

  // Network first for static assets
  event.respondWith(
    fetch(event.request)
      .then((response) => {
        const clone = response.clone();
        caches.open(CACHE_NAME).then((cache) => cache.put(event.request, clone));
        return response;
      })
      .catch(() => caches.match(event.request).then((r) => r || caches.match('/index.html')))
  );
});

// ─── Push notification handler ──────────────────────────────────
self.addEventListener('push', function(event) {
  let data = {
    title: 'HabitFlow',
    body: 'Waktunya cek kebiasaanmu!',
    icon: '/icons/icon-192.png',
    badge: '/icons/badge-72.png',
    data: { url: '/habits' }
  };

  if (event.data) {
    try {
      data = event.data.json();
    } catch (e) {
      data.body = event.data.text();
    }
  }

  event.waitUntil(
    self.registration.showNotification(data.title, {
      body: data.body,
      icon: data.icon || '/icons/icon-192.png',
      badge: data.badge || '/icons/badge-72.png',
      vibrate: [200, 100, 200],
      data: data.data || {},
      actions: [
        { action: 'checkin', title: '✅ Selesai' },
        { action: 'later', title: '⏰ Nanti' }
      ]
    })
  );
});

// ─── Notification click handler ─────────────────────────────────
self.addEventListener('notificationclick', function(event) {
  event.notification.close();
  const notifData = event.notification.data || {};
  const targetUrl = normalizeAppUrl(notifData.url || '/habits');

  if (event.action === 'checkin' && notifData.habit_id) {
    // Direct check-in from notification
    event.waitUntil(
      getStoredToken().then((token) => {
        if (!token) {
          return clients.openWindow(targetUrl);
        }
        return fetch('/api/v1/habits/' + notifData.habit_id + '/check', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
            'Authorization': 'Bearer ' + token
          }
        }).then((resp) => {
          if (resp.ok) {
            return self.registration.showNotification('HabitFlow ✅', {
              body: 'Habit berhasil dicheckin! 🔥',
              icon: '/icons/icon-192.png',
              badge: '/icons/badge-72.png',
            });
          } else {
            return clients.openWindow(targetUrl);
          }
        }).catch(() => {
          return clients.openWindow(targetUrl);
        });
      })
    );
  } else {
    // Open the app to the relevant page
    event.waitUntil(
      clients.matchAll({ type: 'window', includeUncontrolled: true }).then((windowClients) => {
        // If an app window is already open, focus it
        for (const client of windowClients) {
          if (client.url.includes(self.location.origin) && 'focus' in client) {
            client.navigate(targetUrl);
            return client.focus();
          }
        }
        // Otherwise open a new window
        return clients.openWindow(targetUrl);
      })
    );
  }
});

function normalizeAppUrl(rawUrl) {
  if (!rawUrl) return '/#/habits';
  if (rawUrl.startsWith('/#/') || rawUrl.startsWith('#/')) {
    return rawUrl.startsWith('#/') ? '/' + rawUrl : rawUrl;
  }
  if (rawUrl.startsWith('/')) {
    return '/#' + rawUrl;
  }
  return '/#/' + rawUrl.replace(/^#?\/?/, '');
}

// ─── Helper: retrieve stored JWT from IndexedDB ────────────────
function getStoredToken() {
  return new Promise((resolve) => {
    try {
      const request = indexedDB.open('habitflow_auth', 1);

      request.onupgradeneeded = (event) => {
        const db = event.target.result;
        if (!db.objectStoreNames.contains('tokens')) {
          db.createObjectStore('tokens');
        }
      };

      request.onsuccess = (event) => {
        const db = event.target.result;
        if (!db.objectStoreNames.contains('tokens')) {
          resolve(null);
          return;
        }
        const tx = db.transaction('tokens', 'readonly');
        const store = tx.objectStore('tokens');
        const getReq = store.get('access_token');
        getReq.onsuccess = () => resolve(getReq.result || null);
        getReq.onerror = () => resolve(null);
      };

      request.onerror = () => resolve(null);
    } catch (e) {
      resolve(null);
    }
  });
}
