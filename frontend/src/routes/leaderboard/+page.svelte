<script lang="ts">
	const { data } = $props();
</script>

<svelte:head>
	<title>MCC Top 10</title>
</svelte:head>
<div class="p-10 text-xl">
	<h1 class="mb-3 text-center text-2xl">MCC's Top Rated Players</h1>
	<hr class="mb-2" />
	{#each data.leaderboard as player, i}
		<a
			href={`/member/${player.Username}`}
			class={`flex w-full justify-center gap-2 text-xl md:text-2xl`}
		>
			<span>{player.Username}</span><span>{Math.floor(player.Rating)}</span>
		</a>
	{/each}
</div>
