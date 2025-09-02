<script lang="ts">
	import { page } from '$app/state';
	import AcceptChallenge from '$lib/components/AcceptChallenge.svelte';
	import Game from '$lib/components/game.svelte';
	import { websocketStore } from '$lib/websocket.svelte.js';
	let { data } = $props();

	$effect(() => {
		websocketStore.sendMessage({ type: 'room_change', payload: { room: page.params.gameID } });
		return () => websocketStore.sendMessage({ type: 'leave_room' });
	});
</script>

{#key page.params.gameID}
	{#if data.gameData.Creator}
		<AcceptChallenge {data} />
	{:else}
		<Game {data} />
	{/if}
{/key}
