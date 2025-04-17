<script lang="ts">
	import { scrollIntoContainerView } from '$lib/utils';
	import { websocketStore } from '$lib/websocket';
	import { onDestroy, onMount } from 'svelte';
	interface MessageInterface {
		sender: string;
		text: string;
		gameID: string;
	}
	let messages: MessageInterface[] = $state([]);
	const { gameID } = $props();
	let text = $state('');

	let chatContainer: HTMLDivElement;

	const scrollToBottom = (node: HTMLDivElement, watch: () => MessageInterface[]) => {
		$effect(() => {
			void watch?.();
			scrollIntoContainerView(chatContainer, node);
		});
	};

	const handleSend = () => {
		if (text.trim() != '') {
			websocketStore.sendMessage({
				type: 'chat',
				payload: {
					text
				}
			});
			text = '';
		}
	};

	const handleMessage = (payload: any) => {
		if (payload.gameID != gameID) return;
		if (messages) messages = [...messages, payload];
		else messages = [payload];
	};

	onMount(() => {
		websocketStore.onMessage('chat', handleMessage);
	});
	onDestroy(() => {
		websocketStore.offMessage('chat', handleMessage);
	});
</script>

<div class="w-full rounded bg-[#1c1d1e] p-4 shadow-lg">
	<div bind:this={chatContainer} class="mb-2 h-64 overflow-y-auto border-b p-2 text-lg">
		{#each messages as msg}
			{#if msg.gameID == gameID}
				<div class="mb-1">
					<strong>{msg.sender}:</strong>
					{msg.text}
				</div>
			{/if}
		{/each}
		<div use:scrollToBottom={() => messages}></div>
	</div>
	<div class="flex text-sm">
		<input
			type="text"
			bind:value={text}
			class="flex-1 rounded border p-2"
			placeholder="Chat ki maryadaon ka palan krein..."
			onkeydown={(e) => {
				if (e.key === 'Enter') handleSend();
			}}
		/>
		<!-- <button onclick={handleSend} class="rounded bg-blue-500 px-4 py-2 text-white"> Send </button> -->
	</div>
</div>
