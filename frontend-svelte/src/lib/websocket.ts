import { goto } from '$app/navigation';

class WebSocketStore {
	private url: string;
	private ws: WebSocket | null = null;
	private reconnectDelay: number = 2000;

	constructor(url: string) {
		this.url = url;
		this.connect();
	}

	private connect(): void {
		this.ws = new WebSocket(this.url);

		this.ws.onopen = () => console.log('✅ WebSocket Connected');
		this.ws.onmessage = (event: MessageEvent) => {
			const data = JSON.parse(event.data);

			if (data.type === 'Init_Game') {
				goto(`/game/${data.payload.GameID}`);
			}
		};
		this.ws.onclose = () => {
			console.warn('⚠️ WebSocket Disconnected. Reconnecting...');
			setTimeout(() => this.connect(), this.reconnectDelay);
		};
		this.ws.onerror = (error: Event) => console.error('WebSocket Error:', error);
	}

	sendMessage = (message: unknown): void => {
		if (this.ws && this.ws.readyState === WebSocket.OPEN) {
			this.ws.send(JSON.stringify(message));
		} else {
			console.warn('WebSocket not open. Message not sent.');
		}
	};

	get socket() {
		return this.ws;
	}
}

export const websocketStore = new WebSocketStore('ws://localhost:8080/ws');
