<script lang="ts">
	import { scrollIntoContainerView } from '$lib/utils';
	import { websocketStore } from '$lib/websocket.svelte';
	import { onDestroy, onMount } from 'svelte';
	interface MessageInterface {
		sender: string;
		text: string;
	}
	let messages: MessageInterface[] = $state([]);
	const { hei } = $props();
	let text = $state('');

	let chatContainer: HTMLDivElement;
	let chatInput: HTMLInputElement;

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
		if (messages) messages = [...messages, payload];
		else messages = [payload];
	};

	$effect(() => {
		window.addEventListener('scroll', (e) => {
			if (document.activeElement === chatInput) {
				chatInput.blur();
			}
		});
	});

	onMount(() => {
		websocketStore.onMessage('chat', handleMessage);
	});
	onDestroy(() => {
		websocketStore.offMessage('chat', handleMessage);
	});
</script>

<div class="flex-1 rounded bg-[#1c1d1e] p-4 shadow-lg">
	<div
		style={`height: ${hei}`}
		bind:this={chatContainer}
		class={`mb-2 overflow-y-auto border-b p-2 text-lg`}
	>
		{#each messages as msg}
			<div class="mb-1">
				<strong>{msg.sender}:</strong>
				{msg.text}
			</div>
		{/each}
		<div use:scrollToBottom={() => messages}></div>
	</div>
	<div class="flex text-sm">
		<input
			type="text"
			bind:value={text}
			bind:this={chatInput}
			class="flex-1 rounded border p-2"
			placeholder="Chat ki maryadaon ka palan krein..."
			onkeydown={(e) => {
				if (e.key === 'Enter') handleSend();
			}}
		/>
		<!-- <button onclick={handleSend} class="rounded bg-blue-500 px-4 py-2 text-white"> Send </button> -->
	</div>
</div>
