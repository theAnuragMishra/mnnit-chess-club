interface Env {
	BACKEND_URL?: string;
}

export const onRequest: PagesFunction<Env> = async (context) => {
	// Target backend URL: configured in Cloudflare Pages environment variables, or fallback to your Render API URL
	const backendUrl = context.env.BACKEND_URL || 'https://mnnit-chess-club-api.onrender.com';

	const clientUrl = new URL(context.request.url);

	// Strip '/api' prefix: e.g. /api/me -> /me, /api/ws -> /ws
	const targetPath = clientUrl.pathname.replace(/^\/api/, '') || '/';
	const targetUrl = new URL(targetPath + clientUrl.search, backendUrl);

	// Clone request headers and update host headers
	const headers = new Headers(context.request.headers);
	headers.set('Host', targetUrl.host);
	headers.set('X-Forwarded-Host', clientUrl.host);
	headers.set('X-Forwarded-Proto', clientUrl.protocol.replace(':', ''));

	// Handle WebSocket upgrade
	const isWebSocket = context.request.headers.get('Upgrade')?.toLowerCase() === 'websocket';
	if (isWebSocket) {
		return fetch(targetUrl.toString(), {
			method: context.request.method,
			headers: headers
		});
    }

	// Forward HTTP request with manual redirect handling so browser receives Set-Cookie & Location
	const init: RequestInit = {
		method: context.request.method,
		headers: headers,
		redirect: 'manual'
	};

	if (context.request.method !== 'GET' && context.request.method !== 'HEAD') {
		init.body = context.request.body;
	}

	return fetch(targetUrl.toString(), init);
};
