import { getBaseURL } from '$lib/utils';
import { redirect } from '@sveltejs/kit';

export async function load({ params, parent }) {
	const { user } = await parent();
	if (!user) redirect(303, '/');
	const response = await fetch(`${getBaseURL()}/game/${params.gameID}`, {
		credentials: 'include'
	});
	if (!response.ok) {
		throw new Error('Failed to fetch game data');
	}
	return { gameData: await response.json() };
}
