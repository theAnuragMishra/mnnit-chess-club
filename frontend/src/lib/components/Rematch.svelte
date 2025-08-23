<script lang="ts">
	import { websocketStore } from '$lib/websocket';
	import { onDestroy, onMount } from 'svelte';

	let { canRematch } = $props();

	let offer = $state(false);
	let disabled = $state(!canRematch);

	const handleRematch = () => {
		if (disabled) return;
		websocketStore.sendMessage({
			type: 'rematch'
		});
	};
	const handleRematchOffer = (payload: any) => {
		offer = true;
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
	onclick={handleRematch}>Rematch</button
>
