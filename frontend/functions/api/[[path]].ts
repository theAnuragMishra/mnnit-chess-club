interface Env {
	BACKEND_URL?: string;
	RATE_LIMIT_PROXY_SECRET?: string;
}

export const onRequest: PagesFunction<Env> = async (context) => {
	// Target backend URL: configured in Cloudflare Pages environment variables, or fallback to your Render API URL
	const backendUrl = context.env.BACKEND_URL || 'https://mnnit-chess-club-api.onrender.com';

	const clientUrl = new URL(context.request.url);

	// Strip '/api' prefix: e.g. /api/me -> /me, /api/ws -> /ws
	const targetPath = clientUrl.pathname.replace(/^\/api/, '') || '/';
	const targetUrl = new URL(backendUrl);
	targetUrl.pathname = targetPath;
	targetUrl.search = clientUrl.search;

	// Clone request headers and update host headers
	const headers = new Headers(context.request.headers);
	headers.set('Host', targetUrl.host);
	headers.set('X-Forwarded-Host', clientUrl.host);
	headers.set('X-Forwarded-Proto', clientUrl.protocol.replace(':', ''));

	// Only the Pages function can assert a visitor IP to the Render API.
	const visitorIP = context.request.headers.get('CF-Connecting-IP');
	const proxySecret = context.env.RATE_LIMIT_PROXY_SECRET;
	headers.delete('X-Rate-Limit-Client-IP');
	headers.delete('X-Rate-Limit-Proxy-Secret');
	if (visitorIP && proxySecret) {
		headers.set('X-Rate-Limit-Client-IP', visitorIP);
		headers.set('X-Rate-Limit-Proxy-Secret', proxySecret);
	}

	// Handle WebSocket upgrade
	const isWebSocket = context.request.headers.get('Upgrade')?.toLowerCase() === 'websocket';
	if (isWebSocket) {
		return fetch(targetUrl.toString(), {
			method: context.request.method,
			headers,
			redirect: 'manual'
		});
	}

	// Forward HTTP request with manual redirect handling so browser receives Set-Cookie & Location
	const init: RequestInit = {
		method: context.request.method,
		headers,
		redirect: 'manual'
	};

	if (context.request.method !== 'GET' && context.request.method !== 'HEAD') {
		init.body = context.request.body;
	}

	return fetch(targetUrl.toString(), init);
};
