import { user } from '$lib/user.svelte';
import { getBaseURL } from '$lib/utils';
import { error, redirect } from '@sveltejs/kit';

export async function load({ params }) {
	if (!user.id) redirect(303, '/');
	const response = await fetch(`${getBaseURL()}/tournament/${params.tournamentID}`, {
		credentials: 'include'
	});
	if (!response.ok) {
		error(404, { message: 'Tournament not found' });
	}
	const data = await response.json();
	// console.log(data);
	return { tournamentData: data };
}
