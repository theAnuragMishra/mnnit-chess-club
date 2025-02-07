import { getBaseURL } from '../../utils/urlUtils.js';

export const ssr = false;

export const actions = {
	default: async ({ cookies, request }) => {
		const data = await request.formData();
		const username = data.get('username');
		const password = data.get('password');

		const res = await fetch(`${getBaseURL()}/login`, {
			method: 'POST',
			body: JSON.stringify({ username, password }),
			headers: {
				'Content-Type': 'application/json'
			},
			credentials: 'include'
		});
		console.log(res);
	}
};
