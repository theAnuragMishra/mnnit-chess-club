import { getBaseURL } from '$lib/utils.js';

export async function load() {
	const response = await fetch(`${getBaseURL()}/tournament/upcoming`, {
		credentials: 'include'
	});

	if (!response.ok) {
		return { tournaments: undefined };
	}
	return { tournaments: await response.json() };
}
