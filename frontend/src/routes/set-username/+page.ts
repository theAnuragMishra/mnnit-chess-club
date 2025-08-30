import { redirect } from '@sveltejs/kit';

export async function load({ parent }) {
	const { user } = await parent();
	//console.log(user);
	if (!user || user.username) redirect(303, '/');
}
