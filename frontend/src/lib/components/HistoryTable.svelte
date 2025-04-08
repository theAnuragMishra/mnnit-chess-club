<script lang="ts">
	const { setActiveIndex, activeIndex, moveHistory, highlightLastArrow } = $props();
	const moves: (HTMLButtonElement | null)[] = $state([]);
	$effect(() => {
		const activeButton = moves[activeIndex];
		activeButton?.scrollIntoView({ behavior: 'smooth', block: 'nearest' });
	});
</script>

<div
	class="movescontainer relative z-0 h-[40px] w-[80vw] bg-gray-800 text-sm md:h-[200px] md:w-[250px] md:px-3 md:text-lg"
>
	<div
		class="moves flex w-full items-center overflow-x-auto md:grid md:h-[150px] md:grid-cols-[1fr_5fr_5fr] md:place-items-center md:content-start md:items-center md:overflow-x-hidden md:overflow-y-auto"
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
				class={`h-fit w-fit cursor-pointer px-1 py-1.5 md:px-3 ${activeIndex == index ? 'bg-gray-700' : ''}`}
			>
				{move.MoveNotation}
			</button>
		{/each}
	</div>
	{#if moveHistory}
		<button
			aria-label="go to first move"
			class="first flex cursor-pointer items-center justify-center py-1 hover:bg-gray-700"
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
			class="prev flex cursor-pointer items-center justify-center py-1 hover:bg-gray-700"
			onclick={() => {
				setActiveIndex(activeIndex - 1);
			}}
			><svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 16 16"
				><path fill="currentColor" d="M4 14V2h2v5.5l5-5v11l-5-5V14z" /></svg
			></button
		>
		<button
			aria-label="go to next move"
			class="next flex cursor-pointer items-center justify-center py-1 hover:bg-gray-700"
			onclick={() => {
				setActiveIndex(activeIndex + 1);
			}}
			><svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 16 16"
				><path fill="currentColor" d="M12 2v12h-2V8.5l-5 5v-11l5 5V2z" /></svg
			></button
		>
		<button
			aria-label="go to latest move"
			class={`last flex cursor-pointer items-center justify-center py-1 hover:bg-gray-700 ${highlightLastArrow && 'animate-pulse bg-blue-700'}`}
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
			gap: 10px;
			/* grid-template-columns: 1fr 1fr 1fr 1fr; */
			align-items: center;
			align-content: start;
			grid-template-areas:
				'first prev next last'
				'moves moves moves moves';
			padding-top: 2px;
		}
		.movescontainer::before {
			content: '';
			position: absolute;
			background-color: #111e31;
			width: 100%;
			height: 18%;
			top: 0;
			z-index: -1;
		}
		.last,
		.first {
			display: flex;
		}
	}
</style>
