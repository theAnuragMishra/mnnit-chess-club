import { getBaseURL } from '$lib/utils';

export async function load({ params }) {
	const response = await fetch(`${getBaseURL()}/profile/${params.username}?page=1`, {
		credentials: 'include'
	});
	if (!response.ok) return { member: null };
	return { member: await response.json() };
}
