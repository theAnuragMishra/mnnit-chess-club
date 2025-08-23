<script lang="ts">
	import { page } from '$app/state';
	import AcceptChallenge from '$lib/components/AcceptChallenge.svelte';
	import Game from '$lib/components/game.svelte';
	import { websocketStore } from '$lib/websocket';
	import { onDestroy, onMount } from 'svelte';
	let { data } = $props();
	const gameID = page.params.gameID;

	onMount(() => {
		websocketStore.sendMessage({ type: 'room_change', payload: { room: gameID } });
	});
	onDestroy(() => {
		websocketStore.sendMessage({ type: 'leave_room' });
	});
</script>

{#if data.gameData.Creator}
	<AcceptChallenge {data} />
{:else}
	<Game {data} />
{/if}
