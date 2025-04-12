<script lang="ts">
	import { websocketStore } from '$lib/websocket';
	import { onDestroy, onMount } from 'svelte';

	const { userID, gameID, setResultReason, isDisabled } = $props();
	let offer = $state(false);
	let resignSelected = $state(false);

	const handleDraw = () => {
		websocketStore.sendMessage({
			type: 'draw',
			payload: {
				playerID: userID,
				gameID: gameID
			}
		});
	};

	const handleResign = () => {
		if (!resignSelected) {
			resignSelected = true;
			return;
		}
		websocketStore.sendMessage({
			type: 'resign',
			payload: {
				playerID: userID,
				gameID: gameID
			}
		});
	};

	const handleGameDrawn = (payload: any) => {
		if (payload.gameID != gameID) return;
		// console.log('game drawn', payload.Result, payload.Reason);
		setResultReason(payload.Result, payload.Reason);
	};

	const handleDrawOffer = (payload: any) => {
		if (payload.gameID != gameID) return;
		offer = true;
	};

	const handleCancelDraw = (payload: any) => {
		if (payload.gameID != gameID) return;
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

<div class="flex items-center justify-center gap-2 text-xl text-white md:text-2xl">
	<button
		class={`cursor-pointer rounded-lg px-2 py-1 hover:bg-gray-600 md:px-4 md:py-2 ${offer && 'animate-pulse bg-blue-600'} disabled:cursor-not-allowed disabled:hover:bg-transparent`}
		onclick={handleDraw}
		disabled={isDisabled}
	>
		1/2
	</button>
	<button
		aria-label="resign"
		onclick={handleResign}
		disabled={isDisabled}
		class={`cursor-pointer rounded-lg px-2 py-1 hover:bg-red-600 md:px-4 md:py-2 ${resignSelected && 'bg-red-600'} disabled:cursor-not-allowed disabled:hover:bg-transparent`}
	>
		<svg
			class="h-[24px] w-[24px] md:h-[32px] md:w-[32px]"
			xmlns="http://www.w3.org/2000/svg"
			fill="#ffffff"
			viewBox="0 0 256 256"
			><path
				d="M42.76,50A8,8,0,0,0,40,56V224a8,8,0,0,0,16,0V179.77c26.79-21.16,49.87-9.75,76.45,3.41,16.4,8.11,34.06,16.85,53,16.85,13.93,0,28.54-4.75,43.82-18a8,8,0,0,0,2.76-6V56A8,8,0,0,0,218.76,50c-28,24.23-51.72,12.49-79.21-1.12C111.07,34.76,78.78,18.79,42.76,50ZM216,172.25c-26.79,21.16-49.87,9.74-76.45-3.41-25-12.35-52.81-26.13-83.55-8.4V59.79c26.79-21.16,49.87-9.75,76.45,3.4,25,12.35,52.82,26.13,83.55,8.4Z"
			></path></svg
		>
	</button>
	{#if resignSelected}
		<button
			aria-label="cancel resignation attempt"
			class="rounded-lg px-1.5 py-1 hover:bg-gray-600 md:py-2"
			onclick={() => (resignSelected = false)}
			><svg
				class="h-[16px] w-[16px] md:h-[24px] md:w-[24px]"
				xmlns="http://www.w3.org/2000/svg"
				viewBox="0 0 16 16"
				><path
					fill="currentColor"
					d="M15.854 12.854L11 8l4.854-4.854a.503.503 0 0 0 0-.707L13.561.146a.5.5 0 0 0-.707 0L8 5L3.146.146a.5.5 0 0 0-.707 0L.146 2.439a.5.5 0 0 0 0 .707L5 8L.146 12.854a.5.5 0 0 0 0 .707l2.293 2.293a.5.5 0 0 0 .707 0L8 11l4.854 4.854a.5.5 0 0 0 .707 0l2.293-2.293a.5.5 0 0 0 0-.707"
				/></svg
			></button
		>
	{/if}
</div>
