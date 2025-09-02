import { user } from '$lib/user.svelte';
import { redirect } from '@sveltejs/kit';

export async function load() {
	//console.log(user);
	if (!user.id) redirect(303, '/login');
}
