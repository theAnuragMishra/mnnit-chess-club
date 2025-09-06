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
	class={`mb-2 flex h-[50px] w-full cursor-pointer items-center justify-center text-xl disabled:cursor-not-allowed md:text-2xl ${disabled ? 'bg-[#1f2125] text-[#63666b]' : offer ? 'animate-pulse bg-blue-600' : 'bg-[#2a2e34] '}`}
	onclick={handleRematch}>{offered ? 'Cancel Rematch' : 'Rematch'}</button
>
