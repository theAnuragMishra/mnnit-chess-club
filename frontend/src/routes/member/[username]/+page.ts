import { getBaseURL } from '$lib/utils';
import { error } from '@sveltejs/kit';

export async function load({ params }) {
	const r1 = await fetch(`${getBaseURL()}/profile/${params.username}`, { credentials: 'include' });

	if (!r1.ok) error(404, { message: 'Member not found' });

	const profile = await r1.json();

	const r2 = await fetch(`${getBaseURL()}/games/${params.username}?page=1`, {
		credentials: 'include'
	});
	if (!r2.ok) return { profile, games: null };
	const games = await r2.json();
	return { profile, games };
}
