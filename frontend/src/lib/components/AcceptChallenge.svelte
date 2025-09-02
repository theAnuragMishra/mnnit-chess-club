<script>
	import { page } from '$app/state';
	import { PUBLIC_FRONTEND_URL } from '$env/static/public';
	import { user } from '$lib/user.svelte';
	import { websocketStore } from '$lib/websocket.svelte';
	const { data } = $props();
	const gameLink = `${PUBLIC_FRONTEND_URL}/game/${page.params.gameID}/`;
	let copied = $state(false);
</script>

<svelte:head>
	<title
		>Challenge - {data.gameData.TimeControl.baseTime / 60}+{data.gameData.TimeControl
			.increment}</title
	>
</svelte:head>
<div class="flex flex-col items-center justify-center gap-3 text-lg md:text-2xl">
	<p class="text-3xl">
		{data.gameData.TimeControl.baseTime / 60}+{data.gameData.TimeControl.increment}
	</p>
	{#if user.username == data.gameData.CreatorUsername}
		<p class="text-center">The first person to accept the challenge via the link will play you.</p>
		<a class="text-sm text-blue-600 underline sm:text-2xl" href={gameLink}>{gameLink}</a><button
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
