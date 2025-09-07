<script lang="ts">
	import { onDestroy, onMount } from 'svelte';
	import resignImg from '$lib/assets/icons/flag.svg';
	import crossImg from '$lib/assets/icons/cross.svg';
	import berserkImg from '$lib/assets/icons/kill.svg';
	import flipImg from '$lib/assets/icons/flip.svg';
	import { websocketStore } from '$lib/websocket.svelte';
	import { user } from '$lib/user.svelte';

	const { gameID, isDisabled, showBerserkButton, handleFlip, drawOfferedBy } = $props();
	let offer = $state(drawOfferedBy !== user.id && drawOfferedBy !== 0);
	let resignSelected = $state(false);
	let offered = $derived(drawOfferedBy === user.id);

	const handleDraw = () => {
		offered = !offered;
		websocketStore.sendMessage({
			type: 'draw',
			payload: {
				GameID: gameID
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
				GameID: gameID
			}
		});
	};

	const handleDrawOffer = (payload: any) => {
		offer = !offer;
	};

	const handleCancelDraw = (payload: any) => {
		offer = false;
		offered = false;
	};

	const sendBerserk = () => {
		websocketStore.sendMessage({
			type: 'berserk'
		});
	};

	onMount(() => {
		websocketStore.onMessage('drawOffer', handleDrawOffer);
		websocketStore.onMessage('Move_Response', handleCancelDraw);
	});
	onDestroy(() => {
		websocketStore.offMessage('drawOffer', handleDrawOffer);
		websocketStore.offMessage('Move_Response', handleCancelDraw);
	});
</script>

<div class="flex items-center justify-center gap-2 text-xl text-white md:text-2xl">
	<button
		onclick={handleFlip}
		class="cursor-pointer rounded-lg p-1 hover:bg-gray-600 md:px-4 md:py-2"
		><img src={flipImg} alt="flip the board" class="h-[32px]" /></button
	>
	<button
		class={`cursor-pointer rounded-lg p-1  md:px-4 md:py-2 ${offer ? 'animate-pulse bg-blue-600' : offered ? 'bg-red-600' : 'hover:bg-gray-600'} disabled:cursor-not-allowed disabled:hover:bg-transparent`}
		onclick={handleDraw}
		disabled={isDisabled}
	>
		{#if !offer && offered}
			<img src={crossImg} alt="cancel resignation" class="h-[24px]" />
		{:else}
			1/2
		{/if}
	</button>
	<button
		aria-label="resign"
		onclick={handleResign}
		disabled={isDisabled}
		class={`cursor-pointer rounded-lg p-1 hover:bg-red-600 md:px-4 md:py-2 ${resignSelected && 'bg-red-600'} disabled:cursor-not-allowed disabled:hover:bg-transparent`}
	>
		<img class="h-[32px]" src={resignImg} alt="resign button" />
	</button>
	{#if resignSelected}
		<button
			aria-label="cancel resignation attempt"
			class="cursor-pointer rounded-lg p-1 hover:bg-gray-600 md:py-2"
			onclick={() => (resignSelected = false)}
			><img src={crossImg} alt="cancel resignation" class="h-[24px]" /></button
		>
	{/if}
	{#if showBerserkButton}
		<button onclick={sendBerserk} class="cursor-pointer rounded-lg p-1 hover:bg-gray-600 md:py-2">
			<img src={berserkImg} alt="berserk button" class="h-[32px]" />
		</button>
	{/if}
</div>
