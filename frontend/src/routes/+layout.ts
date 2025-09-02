import { user } from '$lib/user.svelte.js';
import { redirect } from '@sveltejs/kit';

export const ssr = false;
export const trailingSlash = 'always';
export const prerender = false;

export async function load({ route }) {
	if (route.id !== '/set-username' && user.id && !user.username) redirect(303, '/set-username');
}
