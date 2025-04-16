import { getBaseURL } from '$lib/utils';
import { error, redirect } from '@sveltejs/kit';

export async function load({ params, parent }) {
	const { user } = await parent();
	if (!user) redirect(303, '/');
	const response = await fetch(`${getBaseURL()}/game/${params.gameID}`, {
		credentials: 'include'
	});
	if (!response.ok) {
		error(404, { message: 'Game not found' });
	}
	return { gameData: await response.json() };
}
