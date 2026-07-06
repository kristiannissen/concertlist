const CACHE_NAME = 'concertlist-v4';
const ASSETS_TO_CACHE = ['/', '/index.html', '/styles.css', '/app.js', '/manifest.json', '/sw.js', '/changelog.json'];

self.addEventListener('install', (event) => {
    event.waitUntil(caches.open(CACHE_NAME).then((cache) => cache.addAll(ASSETS_TO_CACHE)));
});

self.addEventListener('fetch', (event) => {
    event.respondWith(
        caches.match(event.request).then((response) => {
            const fetchPromise = fetch(event.request).then((networkResponse) => {
                if (event.request.method === 'GET' && networkResponse.ok) {
                    const responseClone = networkResponse.clone();
                    caches.open(CACHE_NAME).then((cache) => cache.put(event.request, responseClone));
                }
                return networkResponse;
            });
            return response || fetchPromise;
        })
    );
});

self.addEventListener('activate', (event) => {
    const cacheWhitelist = [CACHE_NAME];
    event.waitUntil(
        caches.keys().then((cacheNames) => {
            return Promise.all(cacheNames.map((cacheName) => {
                if (cacheWhitelist.indexOf(cacheName) === -1) return caches.delete(cacheName);
            }));
        })
    );
});

self.addEventListener('message', (event) => {
    if (event.data && event.data.type === 'SKIP_WAITING') self.skipWaiting();
});

self.addEventListener('install', (event) => event.waitUntil(self.skipWaiting()));

self.addEventListener('updatefound', (event) => {
    if (event.oldVersion) {
        self.clients.matchAll().then((clients) => {
            clients.forEach((client) => client.postMessage({ type: 'NEW_VERSION_AVAILABLE' }));
        });
    }
});