<script lang="ts">
	import { websocketStore } from '$lib/websocket';
	import { onDestroy, onMount } from 'svelte';
	interface MessageInterface {
		sender: string;
		receiver: string;
		text: string;
		gameID: number;
	}
	let messages: MessageInterface[] = $state([]);
	const { username, userID, gameID, whiteID, blackID, whiteUsername, blackUsername } = $props();
	let text = $state('');

	const scrollToBottom = (node: HTMLDivElement, watch: () => MessageInterface[]) => {
		$effect(() => {
			void watch?.();
			node.scrollIntoView({
				behavior: 'smooth'
			});
		});
	};

	const handleSend = () => {
		if (text.trim() != '') {
			websocketStore.sendMessage({
				type: 'chat',
				payload: {
					sender: userID,
					receiver: userID === whiteID ? blackID : whiteID,
					senderUsername: username,
					receiverUsername: username === whiteUsername ? blackUsername : whiteUsername,
					text,
					gameID
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

<div class="rounded border p-4 shadow-lg">
	<div class="mb-2 h-64 overflow-y-auto border-b p-2 text-lg">
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
	<div class="flex w-full gap-2 text-sm">
		<input
			type="text"
			bind:value={text}
			class="flex-1 rounded border p-2"
			placeholder="Type a message..."
		/>
		<button onclick={handleSend} class="rounded bg-blue-500 px-4 py-2 text-white"> Send </button>
	</div>
</div>
