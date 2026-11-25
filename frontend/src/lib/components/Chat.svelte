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

	const mouchEvents: (keyof HTMLElementEventMap)[] = ['mousedown', 'touchstart'];
	const mouchListener: any = (e: MouseEvent) => {
		if (!e.shiftKey && e.buttons !== 2 && e.button !== 2 && e.target !== chatInput)
			chatInput.blur();
	};

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

	onMount(() => {
		websocketStore.onMessage('chat', handleMessage);
	});
	onDestroy(() => {
		websocketStore.offMessage('chat', handleMessage);
		mouchEvents.forEach((event) =>
			document.body.removeEventListener(event, mouchListener, { capture: true })
		);
	});
</script>

<div class="flex-1 rounded-xl border border-gray-700 bg-gray-800 p-2">
	<div
		style={`height: ${hei}`}
		bind:this={chatContainer}
		class={`mb-2 overflow-y-auto border-b border-gray-700 p-2 text-sm break-words`}
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
			class="form-field min-w-0 flex-1"
			aria-label="Chat message"
			placeholder="Type a message"
			onkeydown={(e) => {
				if (e.key === 'Enter') handleSend();
			}}
			onfocus={() => {
				mouchEvents.forEach((event) =>
					document.body.addEventListener(event, mouchListener, { passive: true, capture: true })
				);
			}}
			onblur={() => {
				mouchEvents.forEach((event) =>
					document.body.removeEventListener(event, mouchListener, { capture: true })
				);
			}}
		/>
	</div>
</div>
