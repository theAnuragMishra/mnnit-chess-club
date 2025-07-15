import { getBaseURL } from '$lib/utils.js';
import { redirect } from '@sveltejs/kit';

export async function load({ parent }) {
	const { user } = await parent();
	//console.log(user);
	if (user.role != 2) redirect(303, '/');

	const response = await fetch(`${getBaseURL()}/tournament/upcoming`, {
		credentials: 'include'
	});

	if (!response.ok) {
		return { tournaments: undefined };
	}
	return { tournaments: await response.json() };
}
