class WebSocketStore {
	private url: string;
	private ws: WebSocket | null = null;
	//private reconnectDelay: number = 2000;
	private listeners: Map<string, ((data: any) => void)[]> = new Map();

	constructor(url: string) {
		this.url = url;
	}

	connect(): Promise<void> {
		if (this.ws?.readyState === WebSocket.OPEN) return Promise.resolve();

		return new Promise((resolve, reject) => {
			this.ws = new WebSocket(this.url);

			this.ws.onopen = () => {
				console.log('✅ WebSocket Connected');
				resolve();
			};
			this.ws.onmessage = (event: MessageEvent) => this.handleMessage(event);
			this.ws.onclose = () => {
				console.warn('⚠️ WebSocket Disconnected');
				// setTimeout(() => this.connect(), this.reconnectDelay);
			};
			this.ws.onerror = (error: Event) => {
				console.error('WebSocket Error:', error);
				reject(error);
			};
		});
	}

	private handleMessage(event: MessageEvent): void {
		const data = JSON.parse(event.data);
		const { type, payload } = data;

		if (type === 'GoTo') {
			window.location.href = `/${payload.Type}/${payload.ID}`;
		} else if (type === 'Refresh') {
			if (window.location.pathname === `/${payload.Type}/${payload.ID}`) window.location.reload();
		} else if (this.listeners.has(type)) {
			this.listeners.get(type)?.forEach((callback) => callback(payload));
		}
	}

	onMessage(type: string, callback: (data: any) => void): void {
		if (!this.listeners.has(type)) {
			this.listeners.set(type, []);
		}
		this.listeners.get(type)?.push(callback);
	}

	offMessage(type: string, callback: (data: any) => void): void {
		if (this.listeners.has(type)) {
			const newListeners = this.listeners.get(type)?.filter((cb) => cb !== callback) || [];
			if (newListeners.length > 0) {
				this.listeners.set(type, newListeners);
			} else {
				this.listeners.delete(type);
			}
		}
	}

	sendMessage(message: unknown): void {
		if (this.ws && this.ws.readyState === WebSocket.OPEN) {
			this.ws.send(JSON.stringify(message));
		} else {
			console.warn('WebSocket not open. Message not sent.');
		}
	}

	get socket(): WebSocket | null {
		return this.ws;
	}
}

export const websocketStore = new WebSocketStore('ws://localhost:8080/ws');
