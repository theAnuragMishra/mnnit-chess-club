import type { Handle } from '@sveltejs/kit';
import { getBaseURL } from './utils/urlUtils';

export const handle: Handle = async ({ event, resolve }) => {
	try {
		const sessionTokenCookie = event.cookies.get('session_token');
		const res = await fetch(`${getBaseURL()}/me`, {
			credentials: 'include',
			headers: {
				Cookie: `session_token=${sessionTokenCookie}`
			}
		});
		console.log(res);
	} catch (e) {
		console.log(e);
	}

	const response = resolve(event);
	return response;
};
