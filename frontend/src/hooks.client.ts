import { initUser } from '$lib/user.svelte';
import { websocketStore } from '$lib/websocket.svelte';
import type { ClientInit } from '@sveltejs/kit';

export const init: ClientInit = async () => {
	try {
		await initUser();
		await websocketStore.connect();
	} catch (e) {
		console.error(e);
	}
};
