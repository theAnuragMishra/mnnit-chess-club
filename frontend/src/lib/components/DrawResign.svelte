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

<div
	class="panel flex flex-wrap items-center justify-center gap-1 p-1"
	role="group"
	aria-label="Game actions"
>
	<button
		class="button button-secondary button-compact button-icon"
		onclick={handleFlip}
		aria-label="Flip the board"
	>
		<img src={flipImg} alt="" class="h-5" />
	</button>
	<button
		class="button button-compact button-icon"
		class:button-secondary={!offer && !offered}
		class:button-accent={offer || offered}
		aria-label={offered ? 'Cancel draw offer' : offer ? 'Accept draw' : 'Offer draw'}
		onclick={handleDraw}
		disabled={isDisabled}
	>
		{#if !offer && offered}<img src={crossImg} alt="" class="h-5" />{:else}½{/if}
	</button>
	<button
		class="button button-compact button-icon"
		class:button-secondary={!resignSelected}
		class:button-danger={resignSelected}
		aria-label={resignSelected ? 'Confirm resignation' : 'Resign'}
		onclick={handleResign}
		disabled={isDisabled}
	>
		<img src={resignImg} alt="" class="h-5" />
	</button>
	{#if resignSelected}
		<button
			class="button button-secondary button-compact button-icon"
			aria-label="Cancel resignation"
			onclick={() => (resignSelected = false)}><img src={crossImg} alt="" class="h-5" /></button
		>
	{/if}
	{#if showBerserkButton}
		<button
			class="button button-secondary button-compact button-icon"
			onclick={sendBerserk}
			aria-label="Berserk"
		>
			<img src={berserkImg} alt="" class="h-5" />
		</button>
	{/if}
</div>
