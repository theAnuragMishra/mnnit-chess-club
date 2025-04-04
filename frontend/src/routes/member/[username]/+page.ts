import { getBaseURL } from '$lib/utils';

export async function load({ params, url }) {
	let page = Number(url.searchParams.get('page'));
	if (!page) page = 1;
	// console.log(page);

	const response = await fetch(`${getBaseURL()}/profile/${params.username}?page=${page}`, {
		credentials: 'include'
	});
	if (!response.ok) {
		// console.log('hi');
		return { member: null, page, hasMore: false };
	}
	const member = await response.json();
	// console.log(member);
	return { member, page, hasMore: !!member };
}
