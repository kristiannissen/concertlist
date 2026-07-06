// Service worker caching is disabled while the site is under active
// development and deploys are still stabilizing. A caching service worker
// that outlives a broken deploy will keep serving that broken response
// (including crash/error pages) to returning visitors indefinitely, which
// is worse than having no cache at all during this phase.
//
// This version acts as a kill switch: it clears any cache left behind by
// a previous version of this file, unregisters itself, and reloads any
// open tabs so they fall through to a normal, uncached network request.
// It intentionally has no 'fetch' handler, so once unregistered, requests
// are no longer intercepted at all.
//
// Re-introduce real caching (see git history for the previous
// cache-first implementation) once deploys are stable.

self.addEventListener('install', () => {
    self.skipWaiting();
});

self.addEventListener('activate', (event) => {
    event.waitUntil(
        caches.keys()
            .then((cacheNames) => Promise.all(cacheNames.map((name) => caches.delete(name))))
            .then(() => self.registration.unregister())
            .then(() => self.clients.matchAll({ type: 'window' }))
            .then((clients) => {
                clients.forEach((client) => client.navigate(client.url));
            })
    );
});
