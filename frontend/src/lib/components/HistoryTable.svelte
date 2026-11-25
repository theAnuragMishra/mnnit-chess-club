<script lang="ts">
	import { scrollIntoContainerView } from '$lib/utils';
	import firstImg from '$lib/assets/icons/first.svg';
	import lastImg from '$lib/assets/icons/last.svg';
	import nextImg from '$lib/assets/icons/next.svg';
	import previousImg from '$lib/assets/icons/previous.svg';
	const { setActiveIndex, activeIndex, moveHistory, highlightLastArrow } = $props();
	let container: HTMLDivElement;
	const moveButtons: (HTMLButtonElement | null)[] = $state([]);
	$effect(() => {
		const activeButton = moveButtons[activeIndex];
		if (container && activeButton) scrollIntoContainerView(container, activeButton);
	});
</script>

<div
	class="w-full overflow-hidden rounded-lg border border-gray-700 bg-gray-800 text-sm"
	role="region"
	aria-label="Move history"
>
	<div class="flex w-full items-center justify-between border-b border-gray-700 px-2">
		{#if moveHistory.length > 0}
			<div class="flex gap-1">
				<button
					aria-label="go to first move"
					class="flex min-h-8 min-w-8 cursor-pointer items-center justify-center rounded-md hover:bg-gray-700"
					onclick={() => {
						setActiveIndex(-1);
						container.scrollTo({
							top: 0,
							behavior: 'smooth'
						});
					}}
				>
					<img src={firstImg} alt="go to first move" class="h-4" />
				</button>
				<button
					aria-label="go to previous move"
					class="flex min-h-8 min-w-8 cursor-pointer items-center justify-center rounded-md hover:bg-gray-700"
					onclick={() => {
						setActiveIndex(activeIndex - 1);
					}}><img src={previousImg} class="h-4" alt="go to previous move" /></button
				>
			</div>
			<div class="flex gap-1">
				<button
					aria-label="go to next move"
					class="flex min-h-8 min-w-8 cursor-pointer items-center justify-center rounded-md hover:bg-gray-700"
					onclick={() => {
						setActiveIndex(activeIndex + 1);
					}}><img src={nextImg} alt="go to next move" class="h-4" /></button
				>
				<button
					aria-label="go to latest move"
					class="flex min-h-8 min-w-8 items-center justify-center rounded-md hover:bg-gray-700"
					class:button-accent={highlightLastArrow}
					onclick={() => {
						setActiveIndex(moveHistory.length - 1);
					}}
				>
					<img src={lastImg} alt="go to last move" class="h-4" />
				</button>
			</div>
		{/if}
	</div>
	<div
		bind:this={container}
		class="grid max-h-[clamp(4rem,22dvh,9rem)] grid-cols-[2rem_1fr_1fr] place-items-center overflow-y-auto px-2 py-1"
	>
		{#if moveHistory.length === 0}
			<p class="col-span-full p-2 text-gray-300">No moves yet</p>
		{/if}
		{#each moveHistory as move, index}
			{#if index % 2 == 0}
				<span>{index / 2 + 1}.</span>
			{/if}
			<button
				bind:this={moveButtons[index]}
				onclick={() => {
					setActiveIndex(index);
				}}
				class={`h-fit w-fit cursor-pointer rounded px-2 py-0.5 ${activeIndex == index ? 'bg-gray-700' : ''}`}
			>
				{move.MoveNotation}
			</button>
		{/each}
	</div>
</div>
