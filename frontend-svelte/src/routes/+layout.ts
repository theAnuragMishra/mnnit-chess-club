import { getBaseURL } from '$lib/utils';

export const ssr = false;

export async function load() {
  const res = await fetch(`${getBaseURL()}/me`, { credentials: 'include' });

  if (!res.ok) {
    return { user: null };
  }

  return { user: await res.json() };
}
