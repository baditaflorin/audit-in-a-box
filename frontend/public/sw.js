const CACHE = "audit-in-a-box-v1";
const SHELL = [
  "/audit-in-a-box/",
  "/audit-in-a-box/index.html",
  "/audit-in-a-box/manifest.webmanifest",
];

self.addEventListener("install", (event) => {
  event.waitUntil(caches.open(CACHE).then((cache) => cache.addAll(SHELL)));
});

self.addEventListener("activate", (event) => {
  event.waitUntil(
    caches
      .keys()
      .then((keys) =>
        Promise.all(
          keys.filter((key) => key !== CACHE).map((key) => caches.delete(key)),
        ),
      ),
  );
});

self.addEventListener("fetch", (event) => {
  const request = event.request;
  if (
    request.method !== "GET" ||
    !new URL(request.url).pathname.startsWith("/audit-in-a-box/")
  ) {
    return;
  }
  event.respondWith(
    fetch(request)
      .then((response) => {
        const copy = response.clone();
        caches.open(CACHE).then((cache) => cache.put(request, copy));
        return response;
      })
      .catch(() =>
        caches
          .match(request)
          .then((response) => response || caches.match("/audit-in-a-box/")),
      ),
  );
});
