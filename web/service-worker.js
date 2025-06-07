// 缓存名称和版本
const CACHE_NAME = 'factiories-cache-v1';

// 需要缓存的资源列表
const CACHED_RESOURCES = [
  '/',
  '/index.html',
  '/flutter.js',
  '/main.dart.js',
  '/manifest.json',
  '/favicon.png',
  '/icons/Icon-192.png',
  '/icons/Icon-512.png',
];

// 安装Service Worker
self.addEventListener('install', (event) => {
  event.waitUntil(
    caches.open(CACHE_NAME).then((cache) => {
      return cache.addAll(CACHED_RESOURCES);
    })
  );
});

// 激活Service Worker
self.addEventListener('activate', (event) => {
  event.waitUntil(
    caches.keys().then((cacheNames) => {
      return Promise.all(
        cacheNames.map((cacheName) => {
          if (cacheName !== CACHE_NAME) {
            return caches.delete(cacheName);
          }
        })
      );
    })
  );
});

// 处理请求
self.addEventListener('fetch', (event) => {
  event.respondWith(
    caches.match(event.request).then((response) => {
      // 如果在缓存中找到响应，则返回缓存的响应
      if (response) {
        return response;
      }

      // 否则，发起网络请求
      return fetch(event.request).then((response) => {
        // 检查是否是有效的响应
        if (!response || response.status !== 200 || response.type !== 'basic') {
          return response;
        }

        // 克隆响应，因为响应流只能使用一次
        const responseToCache = response.clone();

        // 将响应添加到缓存
        caches.open(CACHE_NAME).then((cache) => {
          cache.put(event.request, responseToCache);
        });

        return response;
      });
    })
  );
});

// 处理推送通知
self.addEventListener('push', (event) => {
  const options = {
    body: event.data.text(),
    icon: '/icons/Icon-192.png',
    badge: '/icons/Icon-192.png',
  };

  event.waitUntil(
    self.registration.showNotification('Factiories', options)
  );
});

// 处理通知点击
self.addEventListener('notificationclick', (event) => {
  event.notification.close();
  event.waitUntil(
    clients.openWindow('/')
  );
}); 