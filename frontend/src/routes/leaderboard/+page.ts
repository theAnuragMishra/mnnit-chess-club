import { getBaseURL } from '$lib/utils.js';
import { error } from '@sveltejs/kit';

export async function load() {
	const response = await fetch(`${getBaseURL()}/get-leaderboard`, {
		credentials: 'include'
	});

	if (!response.ok) {
		error(500, { message: 'Error fetching leaderboard' });
	}
	return { leaderboard: await response.json() };
}
