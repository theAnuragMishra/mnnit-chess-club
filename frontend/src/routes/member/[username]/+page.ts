import { getBaseURL } from '$lib/utils';

export async function load({ params }) {
	const response = await fetch(`${getBaseURL()}/profile/${params.username}`, {
		credentials: 'include'
	});
	if (!response.ok) {
		return { member: null };
	}
	return { member: await response.json() };
}
