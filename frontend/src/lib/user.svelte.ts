import { invalidateAll } from '$app/navigation';
import { getBaseURL } from './utils';

interface User {
	id: number;
	username: string;
	role: number;
}

export const user: User = $state({} as User);

export async function initUser() {
	const res = await fetch(`${getBaseURL()}/me`, { credentials: 'include' });
	if (!res.ok) return;
	const data = await res.json();
	if (data) Object.assign(user, data);
}

export async function logout() {
	await fetch(`${getBaseURL()}/logout`, {
		method: 'POST',
		credentials: 'include'
	});
	Object.assign(user, Object.fromEntries(Object.keys(user).map((key) => [key, undefined])));
	invalidateAll();
}
