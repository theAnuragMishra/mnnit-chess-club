<script lang="ts">
	import { websocketStore } from '$lib/websocket';
	import { onDestroy, onMount } from 'svelte';

	const { baseTime, increment, opponentID } = $props();
	let offer = $state(false);
	let rematchID: string;

	const handleRematch = () => {
		if (offer) {
			websocketStore.sendMessage({
				type: 'accept_rematch',
				payload: {
					GameId: rematchID
				}
			});
		} else {
			websocketStore.sendMessage({
				type: 'create_rematch',
				payload: {
					//to fix
					timeControl: {
						baseTime,
						increment
					},
					opponentID
				}
			});
		}
	};
	const handleRematchOffer = (payload: any) => {
		offer = true;
		rematchID = payload.rematchID;
	};
	onMount(() => {
		websocketStore.onMessage('rematchOffer', handleRematchOffer);
	});
	onDestroy(() => {
		websocketStore.offMessage('rematchOffer', handleRematchOffer);
	});
</script>

<button
	class={`mb-2 flex h-[50px] w-full cursor-pointer items-center justify-center text-2xl ${offer ? 'animate-pulse bg-blue-600' : 'bg-gray-400 '}`}
	onclick={handleRematch}>Rematch</button
>
