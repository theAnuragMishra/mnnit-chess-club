<script lang="ts">
	import { websocketStore } from '$lib/websocket';
	import { onDestroy, onMount } from 'svelte';

	const { userID, gameID, setResultReason } = $props();
	let offer = $state(false);

	const handleDrawResign = (dr: string) => {
		websocketStore.sendMessage({
			type: dr,
			payload: {
				playerID: userID,
				gameID: gameID
			}
		});
	};

	const handleGameDrawn = (payload: any) => {
		setResultReason(payload.Result, payload.Reason);
	};

	const handleDrawOffer = (payload: any) => {
		if (payload.gameID != gameID) return;
		offer = true;
	};

	const handleCancelDraw = (payload: any) => {
		offer = false;
	};

	onMount(() => {
		websocketStore.onMessage('gameDrawn', handleGameDrawn);
		websocketStore.onMessage('drawOffer', handleDrawOffer);
		websocketStore.onMessage('Move_Response', handleCancelDraw);
	});
	onDestroy(() => {
		websocketStore.offMessage('gameDrawn', handleGameDrawn);
		websocketStore.offMessage('drawOffer', handleDrawOffer);
		websocketStore.offMessage('Move_Response', handleCancelDraw);
	});
</script>

<div class="flex w-full items-center justify-center gap-2 text-white">
	<button
		class={`rounded-lg px-4 py-2 hover:bg-gray-600 ${offer && 'animate-pulse bg-blue-600'}`}
		onclick={() => handleDrawResign('draw')}
	>
		1/2
	</button>
	<button
		aria-label="resign"
		onclick={() => handleDrawResign('resign')}
		class="rounded-lg px-4 py-2 hover:bg-gray-600"
	>
		<svg
			xmlns="http://www.w3.org/2000/svg"
			width="32"
			height="32"
			fill="#ffffff"
			viewBox="0 0 256 256"
			><path
				d="M42.76,50A8,8,0,0,0,40,56V224a8,8,0,0,0,16,0V179.77c26.79-21.16,49.87-9.75,76.45,3.41,16.4,8.11,34.06,16.85,53,16.85,13.93,0,28.54-4.75,43.82-18a8,8,0,0,0,2.76-6V56A8,8,0,0,0,218.76,50c-28,24.23-51.72,12.49-79.21-1.12C111.07,34.76,78.78,18.79,42.76,50ZM216,172.25c-26.79,21.16-49.87,9.74-76.45-3.41-25-12.35-52.81-26.13-83.55-8.4V59.79c26.79-21.16,49.87-9.75,76.45,3.4,25,12.35,52.82,26.13,83.55,8.4Z"
			></path></svg
		>
	</button>
</div>
