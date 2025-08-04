<script lang="ts">
	const { data } = $props();
</script>

<div class="w-full p-10">
	<h2 class="text-2xl">Upcoming:</h2>
	<div>
		<p class="mb-[5px] flex justify-around text-xl">
			<span>Name</span><span>Time Control</span><span>Start Time</span><span>Duration (hours)</span>
		</p>
		<hr class="mb-4" />
		{#each data.tournaments as tournament, i}
			<p class={`flex justify-around text-xl ${i % 2 == 0 ? '' : ''}`}>
				<span><a href={`/tournament/${tournament.ID}`}>{tournament.Name}</a></span><span
					>{tournament.BaseTime / 60}+{tournament.Increment}</span
				><span
					>{new Date(tournament.StartTime).toLocaleString('en-IN', {
						year: 'numeric',
						month: 'short',
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
<div class="w-full p-10">
	<h2 class="text-2xl">Past:</h2>
</div>
