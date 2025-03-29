import { getBaseURL } from '$lib/utils';

export async function load({ params }) {
	const response = await fetch(`${getBaseURL()}/game/${params.gameID}`, {
		credentials: 'include'
	});
	if (!response.ok) {
		throw new Error('Failed to fetch game data');
	}
	return { gameData: await response.json() };
}
