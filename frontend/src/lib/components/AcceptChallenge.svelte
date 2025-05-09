<script>
	import { page } from '$app/state';
	import { websocketStore } from '$lib/websocket';
	const { data } = $props();
	const gameLink = `localhost:5173/game/${page.params.gameID}`;
	let copied = $state(false);
</script>

<div class="flex flex-col items-center justify-center gap-3 text-2xl">
	<p class="text-3xl">
		{data.gameData.TimeControl.baseTime / 60}+{data.gameData.TimeControl.increment}
	</p>
	{#if data.user.username == data.gameData.CreatorUsername}
		<p>The first person to accept the challenge via the link will play you.</p>
		<a class="text-blue-600 underline" href={gameLink}>{gameLink}</a><button
			class="cursor-pointer rounded-md bg-gray-800 px-3 py-2"
			onclick={async () => {
				await navigator.clipboard.writeText(gameLink);
				copied = true;
				setTimeout(() => (copied = false), 1500);
			}}
			>{#if copied}
				✅ Copied!
			{:else}
				📋 Copy
			{/if}</button
		>
	{:else}
		<p>Challenge by {data.gameData.CreatorUsername}</p>
		<button
			class="cursor-pointer rounded-md bg-gray-800 px-3 py-2"
			onclick={() =>
				websocketStore.sendMessage({
					type: 'accept_challenge',
					payload: { GameID: page.params.gameID }
				})}>Accept</button
		>
	{/if}
</div>
