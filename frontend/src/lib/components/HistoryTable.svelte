<script lang="ts">
	const { setActiveIndex, activeIndex, moveHistory, highlightLastArrow } = $props();
	const moves: (HTMLButtonElement | null)[] = $state([]);
	$effect(() => {
		const activeButton = moves[activeIndex];
		activeButton?.scrollIntoView({ behavior: 'smooth', block: 'nearest' });
	});
</script>

<div class="relative flex h-[60px] w-full bg-gray-800 px-4 py-2 text-lg md:h-[310px]">
	{#if moveHistory}
		<button
			aria-label="go to previous move"
			class="py-1 md:hidden"
			onclick={() => {
				setActiveIndex(activeIndex - 1);
			}}
			><svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 16 16"
				><path fill="currentColor" d="M4 14V2h2v5.5l5-5v11l-5-5V14z" /></svg
			></button
		>
	{/if}
	<div
		class="flex items-center justify-between overflow-x-auto md:grid md:h-[250px] md:grid-cols-[1fr_16fr_16fr] md:place-items-center md:content-start md:overflow-x-hidden md:overflow-y-auto"
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
		<button
			aria-label="go to next move"
			class="py-1 hover:bg-gray-700 md:hidden"
			onclick={() => {
				setActiveIndex(activeIndex + 1);
			}}
			><svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 16 16"
				><path fill="currentColor" d="M12 2v12h-2V8.5l-5 5v-11l5 5V2z" /></svg
			></button
		>
	{/if}
	{#if moveHistory}
		<div class="absolute bottom-2 hidden w-4/5 justify-around md:flex">
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
				class={`flex w-1/5 cursor-pointer items-center justify-center py-1 hover:bg-gray-700 ${highlightLastArrow && 'animate-pulse bg-blue-700'}`}
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
