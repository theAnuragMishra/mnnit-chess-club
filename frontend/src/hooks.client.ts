import { websocketStore } from '$lib/websocket';
import type { ClientInit } from '@sveltejs/kit';

export const init: ClientInit = async () => {
	try {
		await websocketStore.connect();
	} catch (e) {
		console.log(e);
	}
};
