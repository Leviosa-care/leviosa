/// <reference types="@sveltejs/kit" />
/// <reference lib="webworker" />

declare let self: ServiceWorkerGlobalScope;

import { build, files, version } from '$service-worker';

const CACHE = `cache-${version}`;

const OAUTH_PATHS = ['/oauth/google/callback', '/oauth/apple/callback'];
const ASSETS = [...build, ...files].filter(
    (path) => !OAUTH_PATHS.some((oauthPath) => path.startsWith(oauthPath))
);

// install service worker
self.addEventListener('install', (event) => {
    async function addFilesToCache() {
        const cache = await caches.open(CACHE);
        await cache.addAll(ASSETS);
    }
    event.waitUntil(addFilesToCache());
});

// activate the service worker
self.addEventListener('activate', (event) => {
    console.log('Service Worker activated with version:', version);
    async function deleteOldCaches() {
        for (const key of await caches.keys()) {
            if (key !== CACHE) {
                await caches.delete(key);
            }
        }
    }
    event.waitUntil(deleteOldCaches());
});

// Listen to fetch events and override the url with your cache so that we can get offline mode
self.addEventListener('fetch', (event) => {
    if (event.request.method !== 'GET') return;

    async function respond() {
        const url = new URL(event.request.url);
        const cache = await caches.open(CACHE);

        // Add special handling for OAuth callback
        const pathnameCondition =
            url.pathname.includes('/oauth/google/callback') ||
            url.pathname.includes('/oauth/apple/callback');
        if (pathnameCondition) {
            try {
                // Ensure we're using the right fetch options
                const response = await fetch(event.request, {
                    method: 'GET',
                    credentials: 'include', // Important for OAuth flows
                    mode: 'cors' // or 'same-origin' depending on your setup
                });
                return response;
            } catch (error) {
                console.error('OAuth callback detailed error:', {
                    message: error.message,
                    stack: error.stack,
                    type: error.constructor.name,
                    url: event.request.url
                });

                // Forward the request to the network without service worker intervention
                return fetch(event.request);
            }
        }

        // Regular asset handling
        if (ASSETS.includes(url.pathname)) {
            const cachedResponse = await cache.match(url.pathname);
            if (cachedResponse) {
                return cachedResponse;
            }
        }

        // Network first strategy
        try {
            const response = await fetch(event.request);
            const isNotExtension = url.protocol === 'http:';
            const isSuccess = response.status === 200;

            if (isNotExtension && isSuccess) {
                cache.put(event.request, response.clone());
            }
            return response;
        } catch {
            const cachedResponse = await cache.match(event.request);
            if (cachedResponse) {
                return cachedResponse;
            }
        }

        return new Response('Not found', { status: 404 });
    }

    event.respondWith(respond());
});

self.addEventListener('message', (event) => {
    if (event.data && event.data.type === 'SKIP WAITING') {
        self.skipWaiting();
    }
});
