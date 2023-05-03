/**
 * service worker
 */
var cacheName = '';
var apiCacheName = '';
var cacheFiles = [
    '/',
    './index.html',
];

// 监听install事件
self.addEventListener('install', function (e) {
    console.log('Service Worker 状态： install');
});

// 监听activate事件
self.addEventListener('activate', function (e) {
    console.log('Service Worker 状态： activate');
});

self.addEventListener('fetch', function (e) {
    // 需要缓存的xhr请求
    console.log("fetch 事件")
});
