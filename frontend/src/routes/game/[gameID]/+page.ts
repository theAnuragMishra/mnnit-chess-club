import { user } from '$lib/user.svelte';
import { getBaseURL } from '$lib/utils';
import { error, redirect } from '@sveltejs/kit';

export async function load({ params }) {
	if (!user.id) redirect(303, '/');
	const response = await fetch(`${getBaseURL()}/game/${params.gameID}`, {
		credentials: 'include'
	});
	if (!response.ok) {
		error(404, { message: 'Game not found' });
	}
	return { gameData: await response.json() };
}
