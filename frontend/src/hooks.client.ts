import { websocketStore } from '$lib/websocket';
import type { ClientInit } from '@sveltejs/kit';

export const init: ClientInit = async () => {
	await websocketStore.connect();
};
