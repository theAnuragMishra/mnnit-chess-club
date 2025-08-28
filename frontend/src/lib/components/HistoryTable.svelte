<script lang="ts">
	import firstImg from '$lib/assets/icons/first.svg';
	import lastImg from '$lib/assets/icons/last.svg';
	import nextImg from '$lib/assets/icons/next.svg';
	import previousImg from '$lib/assets/icons/previous.svg';
	const { setActiveIndex, activeIndex, moveHistory, highlightLastArrow } = $props();
	const moveButtons: (HTMLButtonElement | null)[] = $state([]);
	$effect(() => {
		const activeButton = moveButtons[activeIndex];
		activeButton?.scrollIntoView({ behavior: 'smooth', block: 'end' });
	});
</script>

<div class="w-full bg-gray-800 text-sm md:w-[300px] md:text-lg">
	<div class="flex w-full items-center justify-between bg-[#232327] px-5 md:px-2">
		{#if moveHistory.length > 0}
			<div class="flex gap-[10px]">
				<button
					aria-label="go to first move"
					class="flex cursor-pointer items-center justify-center py-1 hover:bg-gray-700"
					onclick={() => {
						setActiveIndex(-1);
					}}
				>
					<img src={firstImg} alt="go to first move" class="h-[16px] md:h-[24px]" />
				</button>
				<button
					aria-label="go to previous move"
					class="flex cursor-pointer items-center justify-center py-1 hover:bg-gray-700"
					onclick={() => {
						setActiveIndex(activeIndex - 1);
					}}
					><img src={previousImg} class="h-[16px] md:h-[24px]" alt="go to previous move" /></button
				>
			</div>
			<div class="flex gap-[10px]">
				<button
					aria-label="go to next move"
					class="flex cursor-pointer items-center justify-center py-1 hover:bg-gray-700"
					onclick={() => {
						setActiveIndex(activeIndex + 1);
					}}><img src={nextImg} alt="go to next move" class="h-[16px] md:h-[24px]" /></button
				>
				<button
					aria-label="go to latest move"
					class={`flex cursor-pointer items-center justify-center py-1 hover:bg-gray-700 ${highlightLastArrow && 'animate-pulse bg-blue-700'}`}
					onclick={() => {
						setActiveIndex(moveHistory.length - 1);
					}}
				>
					<img src={lastImg} alt="go to last move" class="h-[16px] md:h-[24px]" />
				</button>
			</div>
		{/if}
	</div>
	<div
		class="grid h-[50px] grid-cols-[1fr_5fr_5fr] place-items-center overflow-y-auto px-2 py-2 md:h-[70px] md:overflow-x-hidden"
	>
		{#each moveHistory as move, index}
			{#if index % 2 == 0}
				<span>{index / 2 + 1}.</span>
			{/if}
			<button
				bind:this={moveButtons[index]}
				onclick={() => {
					setActiveIndex(index);
				}}
				class={`h-fit w-fit cursor-pointer px-1 md:px-3 ${activeIndex == index ? 'bg-gray-700' : ''}`}
			>
				{move.MoveNotation}
			</button>
		{/each}
	</div>
</div>

<style>
</style>
