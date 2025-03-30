<script lang="ts">
	const { setActiveIndex, activeIndex, moveHistory } = $props();
	const moves: (HTMLButtonElement | null)[] = $state([]);
	$effect(() => {
		const activeButton = moves[activeIndex];
		activeButton?.scrollIntoView({ behavior: 'smooth', block: 'nearest' });
	});
</script>

<div class="relative h-[310px] bg-gray-800 px-4 py-2 text-lg">
	<div
		class="grid h-[250px] w-full grid-cols-[1fr_16fr_16fr] place-items-center content-start overflow-x-hidden overflow-y-auto"
	>
		{#each moveHistory as move, index}
			{#if index % 2 == 0}
				<span>{index / 2 + 1}.</span>
			{/if}
			<button
				bind:this={moves[index]}
				onclick={() => {
					setActiveIndex(index);
				}}
				class={`h-fit w-fit cursor-pointer px-4 py-2 ${activeIndex == index ? 'bg-gray-700' : ''}`}
			>
				{move.MoveNotation}
			</button>
		{/each}
	</div>
</div>
