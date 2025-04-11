import { goto } from '$app/navigation';
class WebSocketStore {
	private url: string;
	private ws: WebSocket | null = null;
	// private reconnectDelay: number = 2000;
	private listeners: Map<string, ((data: any) => void)[]> = new Map();

	constructor(url: string) {
		this.url = url;
	}

	tryConnect(): Promise<void> {
		if (this.ws?.readyState === WebSocket.OPEN) return Promise.resolve();
		console.log('inside tryconnect');
		let settled = false;
		return new Promise((resolve, reject) => {
			this.ws = new WebSocket(this.url);

			this.ws.onopen = () => {
				if (settled) return;
				settled = true;
				console.log('✅ WebSocket Connected');
				resolve();
			};
			this.ws.onmessage = (event: MessageEvent) => this.handleMessage(event);
			this.ws.onclose = () => {
				console.warn('⚠️ WebSocket Disconnected');
				// setTimeout(() => this.connect(), this.reconnectDelay);
			};
			this.ws.onerror = (error: Event) => {
				if (settled) return;
				settled = true;
				console.error('WebSocket Error:', error);

				reject(error);
			};
		});
	}

	async retryConnect(retries = 5): Promise<void> {
		console.log('insde retryconnect');
		for (let i = 0; i < retries; i++) {
			try {
				await this.tryConnect();
				return; // success
			} catch (error) {
				console.error(error);
				console.warn(`Retrying... (${i + 1}/${retries})`);
			}
		}
	}

	async connect(): Promise<void> {
		try {
			await this.tryConnect();
		} catch (err) {
			console.log(err, 'Initial connect failed. Retrying...');
			await this.retryConnect();
		}
	}

	private handleMessage(event: MessageEvent): void {
		const data = JSON.parse(event.data);
		const { type, payload } = data;

		if (type === 'Init_Game') {
			console.log('init game received');
			goto(`/game/${payload.GameID}`);
		}

		if (this.listeners.has(type)) {
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
			console.log('message sent');
		} else {
			console.warn('WebSocket not open. Message not sent.');
		}
	}

	get socket(): WebSocket | null {
		return this.ws;
	}
}

export const websocketStore = new WebSocketStore('ws://localhost:8080/ws');
