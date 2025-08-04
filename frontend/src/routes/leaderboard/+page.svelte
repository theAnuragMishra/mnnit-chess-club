<script lang="ts">
	const { data } = $props();
	const topSize = 40;
</script>

<div class="p-10 text-xl">
	<h1 class="mb-3 text-4xl">MCC's Top Rated Players</h1>
	<hr class="mb-2" />
	{#each data.leaderboard as player, i}
		<p style={`font-size: ${topSize - 2 * i}px`} class={`flex w-full justify-center gap-2`}>
			<span>{player.Username}</span><span>{Math.floor(player.Rating)}</span>
		</p>
	{/each}
</div>
