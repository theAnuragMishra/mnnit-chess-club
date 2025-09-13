<script lang="ts">
	const { data } = $props();
</script>

<svelte:head>
	<title>Tournaments - MCC</title>
</svelte:head>
<div class="flex w-full flex-col gap-5 px-5 py-5">
	{#if data.live.length > 0}
		<div class="w-full p-2">
			<h2 class="mb-2 text-center text-lg md:text-2xl">Live Arenas</h2>
			<div class="titles flex flex-col gap-3">
				<!-- <p class="mb-[5px] flex justify-around text-sm md:text-xl">
					<span>Name</span><span>Time Control</span><span>Start Time</span><span
						>Duration (hours)</span
					>
				</p> -->

				{#each data.live as tournament, i}
					<p class={`flex justify-around bg-[#212632] py-3 md:text-lg`}>
						<span><a href={`/tournament/${tournament.ID}`}>{tournament.Name}</a></span><span
							>{tournament.BaseTime / 60}+{tournament.Increment}</span
						><span
							>{new Date(tournament.StartTime).toLocaleString('en-IN', {
								year: '2-digit',
								month: 'short',
								day: '2-digit',
								hour: '2-digit',
								minute: '2-digit',
								hour12: true
							})}</span
						><span>{tournament.Duration / 3600} hrs</span>
					</p>
				{/each}
			</div>
		</div>
	{/if}
	{#if data.scheduled.length > 0}
		<div class="w-full p-2">
			<h2 class="mb-2 text-center text-lg md:text-2xl">Coming up</h2>
			<div class="titles flex flex-col gap-3">
				<!-- <p class="mb-[5px] flex justify-around text-sm md:text-xl">
					<span>Name</span><span>Time Control</span><span>Start Time</span><span
						>Duration (hrs)</span
					>
				</p> -->

				{#each data.scheduled as tournament, i}
					<p class={`flex justify-around bg-[#212632] py-3 md:text-lg`}>
						<span><a href={`/tournament/${tournament.ID}`}>{tournament.Name}</a></span><span
							>{tournament.BaseTime / 60}+{tournament.Increment}</span
						><span
							>{new Date(tournament.StartTime).toLocaleString('en-IN', {
								year: '2-digit',
								month: 'short',
								day: '2-digit',
								hour: '2-digit',
								minute: '2-digit',
								hour12: true
							})}</span
						><span>{tournament.Duration / 3600} hrs</span>
					</p>
				{/each}
			</div>
		</div>
	{/if}
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
