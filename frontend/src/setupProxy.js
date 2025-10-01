const { createProxyMiddleware } = require('http-proxy-middleware');

module.exports = function(app) {
  app.use(
    '/api',
    createProxyMiddleware({
      target: 'http://localhost:8080',
      changeOrigin: true,
      secure: false,
      onProxyReq: (proxyReq, req, res) => {
        // Add CORS headers
        proxyReq.setHeader('Access-Control-Allow-Origin', 'http://localhost:3000');
        proxyReq.setHeader('Access-Control-Allow-Methods', 'GET, POST, PUT, DELETE, OPTIONS');
        proxyReq.setHeader('Access-Control-Allow-Headers', 'Content-Type, Authorization');
      },
      onProxyRes: (proxyRes, req, res) => {
        // Add CORS headers to the response
        proxyRes.headers['Access-Control-Allow-Origin'] = 'http://localhost:3000';
        proxyRes.headers['Access-Control-Allow-Credentials'] = 'true';
      },
    })
  );
};
