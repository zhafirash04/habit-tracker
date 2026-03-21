if ('serviceWorker' in navigator) {
  navigator.serviceWorker.register('/sw.js')
    .then(r => console.log('SW registered:', r.scope))
    .catch(e => console.error('SW error:', e));
}
