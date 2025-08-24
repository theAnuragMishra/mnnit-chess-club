import { getBaseURL } from '$lib/utils.js';

export async function load() {
	const scheduledRes = await fetch(`${getBaseURL()}/tournament/scheduled`, {
		credentials: 'include'
	});
	const liveRes = await fetch(`${getBaseURL()}/tournament/live`, {
		credentials: 'include'
	});
	const scheduled = scheduledRes.ok ? await scheduledRes.json() : null;
	const live = liveRes.ok ? await liveRes.json() : null;

	return { scheduled, live };
}
