<script lang="ts">
	const { setActiveIndex, activeIndex, moveHistory, highlightLastArrow } = $props();
	const moves: (HTMLButtonElement | null)[] = $state([]);
	$effect(() => {
		const activeButton = moves[activeIndex];
		activeButton?.scrollIntoView({ behavior: 'smooth', block: 'nearest' });
	});
</script>

<div
	class="movescontainer w-[80vw h-[40px] bg-gray-800 px-3 text-sm md:h-[310px] md:w-fit md:text-lg"
>
	<div
		class="moves flex w-full items-center overflow-x-auto md:grid md:h-[250px] md:grid-cols-[1fr_16fr_16fr] md:place-items-center md:content-start md:overflow-x-hidden md:overflow-y-auto"
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
				class={`h-fit w-fit cursor-pointer px-3 py-1.5 ${activeIndex == index ? 'bg-gray-700' : ''}`}
			>
				{move.MoveNotation}
			</button>
		{/each}
	</div>
	{#if moveHistory}
		<button
			aria-label="go to first move"
			class="first flex w-1/5 cursor-pointer items-center justify-center py-1 hover:bg-gray-700"
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
			class="prev flex w-1/5 cursor-pointer items-center justify-center py-1 hover:bg-gray-700"
			onclick={() => {
				setActiveIndex(activeIndex - 1);
			}}
			><svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 16 16"
				><path fill="currentColor" d="M4 14V2h2v5.5l5-5v11l-5-5V14z" /></svg
			></button
		>
		<button
			aria-label="go to next move"
			class="next flex w-1/5 cursor-pointer items-center justify-center py-1 hover:bg-gray-700"
			onclick={() => {
				setActiveIndex(activeIndex + 1);
			}}
			><svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 16 16"
				><path fill="currentColor" d="M12 2v12h-2V8.5l-5 5v-11l5 5V2z" /></svg
			></button
		>
		<button
			aria-label="go to latest move"
			class={`last flex w-1/5 cursor-pointer items-center justify-center py-1 hover:bg-gray-700 ${highlightLastArrow && 'animate-pulse bg-blue-700'}`}
			onclick={() => {
				setActiveIndex(moveHistory.length - 1);
			}}
		>
			<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 16 16"
				><path fill="currentColor" d="M14 2v12h-2V8.5l-5 5v-5l-5 5v-11l5 5v-5l5 5V2z" /></svg
			>
		</button>
	{/if}
</div>

<style>
	.movescontainer {
		display: grid;
		/* grid-template-columns: 1fr 4fr 1fr; */
		grid-template-areas: 'prev moves next';
	}
	.moves {
		grid-area: moves;
	}
	.first {
		grid-area: first;
		display: none;
	}
	.next {
		grid-area: next;
		justify-self: end;
	}
	.last {
		grid-area: last;
		display: none;
	}
	.prev {
		grid-area: prev;
		justify-self: start;
	}

	@media (width >=768px) {
		.movescontainer {
			grid-template-areas:
				'first prev next last'
				'moves moves moves moves';
			justify-content: space-between;
		}
		.last,
		.first {
			display: inline-block;
		}
	}
</style>
