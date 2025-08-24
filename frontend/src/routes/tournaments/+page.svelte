<script lang="ts">
	const { data } = $props();
</script>

<svelte:head>
	<title>Tournaments - MCC</title>
</svelte:head>
<div class="w-full p-2 md:p-10">
	<h2 class="text-lg md:text-2xl">Live:</h2>
	<div class="titles">
		<p class="mb-[5px] flex justify-around text-lg md:text-xl">
			<span>Name</span><span>Time Control</span><span>Start Time</span><span>Duration (hours)</span>
		</p>
		<hr class="mb-4" />
		{#each data.live as tournament, i}
			<p class={`flex justify-around text-lg ${i % 2 == 0 ? '' : ''}`}>
				<span><a href={`/tournament/${tournament.ID}`}>{tournament.Name}</a></span><span
					>{tournament.BaseTime / 60}+{tournament.Increment}</span
				><span
					>{new Date(tournament.StartTime).toLocaleString('en-IN', {
						year: '2-digit',
						month: '2-digit',
						day: '2-digit',
						hour: '2-digit',
						minute: '2-digit',
						hour12: false
					})}</span
				><span>{tournament.Duration / 3600}</span>
			</p>
		{/each}
	</div>
</div>
<div class="w-full p-2 md:p-10">
	<h2 class="text-lg md:text-2xl">Scheduled:</h2>
	<div class="titles">
		<p class="mb-[5px] flex justify-around text-lg md:text-xl">
			<span>Name</span><span>Time Control</span><span>Start Time</span><span>Duration (hours)</span>
		</p>
		<hr class="mb-4" />
		{#each data.scheduled as tournament, i}
			<p class={`flex justify-around text-lg ${i % 2 == 0 ? '' : ''}`}>
				<span><a href={`/tournament/${tournament.ID}`}>{tournament.Name}</a></span><span
					>{tournament.BaseTime / 60}+{tournament.Increment}</span
				><span
					>{new Date(tournament.StartTime).toLocaleString('en-IN', {
						year: '2-digit',
						month: '2-digit',
						day: '2-digit',
						hour: '2-digit',
						minute: '2-digit',
						hour12: false
					})}</span
				><span>{tournament.Duration / 3600}</span>
			</p>
		{/each}
	</div>
</div>

<!-- <div class="w-full p-2 md:p-10">
	<h2 class="text-2xl">Past:</h2>
</div> -->

<style>
	.titles > p > span {
		width: 25%;
		text-align: center;
	}
</style>
