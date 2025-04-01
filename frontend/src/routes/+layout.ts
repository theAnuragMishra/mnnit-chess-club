import { getBaseURL } from '$lib/utils';
import { redirect } from '@sveltejs/kit';

export const ssr = false;

export async function load() {
	const res = await fetch(`${getBaseURL()}/me`, { credentials: 'include' });

	if (!res.ok) {
		return { user: null };
	}

	const user = await res.json();
	if (!user.username) redirect(303, '/set-username');

	return { user };
}
