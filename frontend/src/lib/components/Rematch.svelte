<script lang="ts">
	import { user } from '$lib/user.svelte';
	import { websocketStore } from '$lib/websocket.svelte';
	import { onDestroy, onMount } from 'svelte';

	let { canRematch, offeredBy } = $props();

	let offer = $state(offeredBy !== 0 && offeredBy !== user.id);
	let offered = $state(offeredBy === user.id);
	let disabled = $state(!canRematch);

	const handleRematch = () => {
		if (disabled) return;
		offered = !offered;
		websocketStore.sendMessage({
			type: 'rematch'
		});
	};
	const handleRematchOffer = (payload: any) => {
		offer = !offer;
	};

	const handleGameDeleted = (payload: any) => {
		disabled = true;
		offer = false;
	};

	onMount(() => {
		websocketStore.onMessage('rematchOffer', handleRematchOffer);
		websocketStore.onMessage('GameDeleted', handleGameDeleted);
	});
	onDestroy(() => {
		websocketStore.offMessage('rematchOffer', handleRematchOffer);
		websocketStore.offMessage('GameDeleted', handleGameDeleted);
	});
</script>

<button
	{disabled}
	class="button w-full"
	class:button-secondary={!offer}
	class:button-accent={offer}
	onclick={handleRematch}>{offered ? 'Cancel Rematch' : 'Rematch'}</button
>
