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
	{#if moveHistory}
		<div class="absolute bottom-2 flex w-4/5 justify-around">
			<button
				aria-label="go to first move"
				class="flex w-1/5 cursor-pointer items-center justify-center py-1 hover:bg-gray-700"
				onclick={() => {
					setActiveIndex(0);
				}}
			>
				<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 16 16"
					><path fill="currentColor" d="M2 14V2h2v5.5l5-5v5l5-5v11l-5-5v5l-5-5V14z" /></svg
				>
			</button>
			<button
				aria-label="go to previous move"
				class="flex w-1/5 cursor-pointer items-center justify-center py-1 hover:bg-gray-700"
				onclick={() => {
					setActiveIndex(activeIndex - 1);
				}}
				><svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 16 16"
					><path fill="currentColor" d="M4 14V2h2v5.5l5-5v11l-5-5V14z" /></svg
				></button
			>
			<button
				aria-label="go to next move"
				class="flex w-1/5 cursor-pointer items-center justify-center py-1 hover:bg-gray-700"
				onclick={() => {
					setActiveIndex(activeIndex + 1);
				}}
				><svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 16 16"
					><path fill="currentColor" d="M12 2v12h-2V8.5l-5 5v-11l5 5V2z" /></svg
				></button
			>
			<button
				aria-label="go to latest move"
				class="flex w-1/5 cursor-pointer items-center justify-center py-1 hover:bg-gray-700"
				onclick={() => {
					setActiveIndex(moveHistory.length - 1);
				}}
			>
				<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 16 16"
					><path fill="currentColor" d="M14 2v12h-2V8.5l-5 5v-5l-5 5v-11l5 5v-5l5 5V2z" /></svg
				>
			</button>
		</div>
	{/if}
</div>
