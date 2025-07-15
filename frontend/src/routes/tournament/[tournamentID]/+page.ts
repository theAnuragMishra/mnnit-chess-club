import { getBaseURL } from '$lib/utils';
import { error, redirect } from '@sveltejs/kit';

export async function load({ params, parent }) {
	const { user } = await parent();
	if (!user) redirect(303, '/');
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
