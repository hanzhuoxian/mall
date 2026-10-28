const backend = { target: 'http://localhost:9091', changeOrigin: true };
const apiRoutes = ['/v1/', '/login', '/register', '/logout', '/refresh', '/captcha'];
const local = Object.fromEntries(apiRoutes.map((path) => [path, backend]));
export default { dev: local, test: local };
