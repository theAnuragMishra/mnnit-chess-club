<script lang="ts">
	import { page } from '$app/state';
	import { websocketStore } from '$lib/websocket.svelte';
	import { getTimeControl } from '$lib/utils';
	import Icon from '$lib/components/Icon.svelte';
	const timeControls = [
		{ baseTime: 60, increment: 0 },
		{ baseTime: 60, increment: 1 },
		{ baseTime: 120, increment: 1 },
		{ baseTime: 180, increment: 0 },
		{ baseTime: 180, increment: 2 },
		{ baseTime: 300, increment: 0 },
		{ baseTime: 300, increment: 3 },
		{ baseTime: 600, increment: 0 },
		{ baseTime: 600, increment: 5 },
		{ baseTime: 900, increment: 10 },
		{ baseTime: 1800, increment: 0 },
		{ baseTime: 1800, increment: 20 }
	];
	const baseTimes = [
		0.25, 0.5, 0.75, 1, 1.5, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 25,
		30, 35, 40, 45, 60, 75, 90, 105, 120, 135, 150, 165, 180
	];
	const increments = [
		0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 25, 30, 35, 40, 45,
		60, 90, 120, 150, 180
	];
	const filters = [
		'All games',
		...new Set(timeControls.map((control) => getTimeControl(control.baseTime, control.increment)))
	];
	let selectedFilter = $state('All games');
	$effect(() => {
		selectedFilter =
			filters.find((filter) => filter.toLowerCase() === page.url.searchParams.get('mode')) ??
			'All games';
	});
	let activeIndex = $state(-1);
	let baseIndex = $state(6);
	let incrementIndex = $state(0);
	function handleInitGame(index: number) {
		activeIndex = activeIndex === index ? -1 : index;
		websocketStore.sendMessage({ type: 'init_game', payload: index });
	}
</script>

<svelte:head><title>Play chess — MNNIT Chess Club</title></svelte:head>

<div class="play-page">
	<div class="page-title">
		<h1>Play chess</h1>
		<a href="#challenge" class="text-link jump-link">Challenge a friend</a>
	</div>

	<div class="play-layout">
		<section class="pair-panel panel" aria-labelledby="pair-heading">
			<div class="panel-heading">
				<h2 id="pair-heading">Quick pairing</h2>
				<div class="filters" role="group" aria-label="Time control">
					{#each filters as filter}
						<button
							class:chosen={selectedFilter === filter}
							aria-pressed={selectedFilter === filter}
							onclick={() => (selectedFilter = filter)}
							>{filter === 'All games' ? 'All' : filter}</button
						>
					{/each}
				</div>
			</div>
			<div class="time-grid">
				{#each timeControls as control, index}
					{@const format = getTimeControl(control.baseTime, control.increment)}
					{#if selectedFilter === 'All games' || selectedFilter === format}
						<button
							class="time-button"
							class:searching={activeIndex === index}
							aria-pressed={activeIndex === index}
							aria-label={activeIndex === index
								? 'Cancel matchmaking'
								: `Play ${control.baseTime / 60} minutes plus ${control.increment} seconds, ${format}`}
							onclick={() => handleInitGame(index)}
						>
							<span class="time-value"
								>{control.baseTime / 60}<span> + </span>{control.increment}</span
							>
							<span class="time-format">{activeIndex === index ? 'Cancel search' : format}</span>
						</button>
					{/if}
				{/each}
			</div>

			<div class="queue-status" aria-live="polite">
				{#if activeIndex !== -1}
					<p>Looking for an opponent…</p>
					<button class="text-link" onclick={() => handleInitGame(activeIndex)}>Cancel</button>
				{/if}
			</div>
		</section>

		<section id="challenge" class="challenge-panel panel" aria-labelledby="challenge-heading">
			<h2 id="challenge-heading">Challenge a friend</h2>
			<div class="range-field">
				<label for="minutes">Minutes per side <strong>{baseTimes[baseIndex]}</strong></label>
				<input
					id="minutes"
					type="range"
					bind:value={baseIndex}
					min="0"
					max={baseTimes.length - 1}
					aria-valuetext={`${baseTimes[baseIndex]} minutes`}
				/>
			</div>
			<div class="range-field">
				<label for="increment">Increment <strong>{increments[incrementIndex]} sec</strong></label>
				<input
					id="increment"
					type="range"
					bind:value={incrementIndex}
					min="0"
					max={increments.length - 1}
					aria-valuetext={`${increments[incrementIndex]} seconds`}
				/>
			</div>
			<button
				class="button button-primary"
				onclick={() =>
					websocketStore.sendMessage({
						type: 'create_challenge',
						payload: { baseTime: baseTimes[baseIndex] * 60, increment: increments[incrementIndex] }
					})}><Icon name="link" size={18} />Create challenge link</button
			>
		</section>
	</div>
</div>

<style>
	.play-page {
		width: min(64rem, calc(100% - 2rem));
		margin-inline: auto;
		padding-block: clamp(0.75rem, 3dvh, 1.75rem);
	}
	.page-title {
		display: flex;
		align-items: baseline;
		justify-content: space-between;
		gap: 0.75rem;
		margin-bottom: 0.875rem;
	}
	h1 {
		font-size: clamp(1.375rem, 4vw, 1.875rem);
		font-weight: 650;
	}
	h2 {
		font-size: 1rem;
		font-weight: 650;
	}
	.jump-link {
		display: none;
	}
	.play-layout {
		display: grid;
		grid-template-columns: minmax(0, 1fr) minmax(14rem, 0.42fr);
		align-items: start;
		gap: 0.875rem;
	}
	.pair-panel,
	.challenge-panel {
		padding: clamp(0.75rem, 2vw, 1.25rem);
		border-radius: 12px;
	}
	.panel-heading {
		display: flex;
		flex-wrap: wrap;
		align-items: center;
		justify-content: space-between;
		gap: 0.25rem 1rem;
		margin-bottom: 0.625rem;
	}
	.filters {
		display: flex;
		flex-wrap: wrap;
		gap: 0.125rem;
	}
	.filters button {
		min-height: 2.125rem;
		padding: 0.25rem 0.625rem;
		border-radius: 6px;
		color: var(--muted);
		font-size: 0.875rem;
	}
	.filters button:hover,
	.filters .chosen {
		color: var(--accent);
		background: #2c3d2f;
	}
	.time-grid {
		display: grid;
		grid-template-columns: repeat(3, minmax(0, 1fr));
		gap: 0.5rem;
	}
	.time-button {
		display: flex;
		min-height: clamp(3.375rem, 11dvh, 4.5rem);
		flex-direction: column;
		align-items: center;
		justify-content: center;
		gap: 0.25rem;
		padding-inline: 0.125rem;
		border: 1px solid #354238;
		border-radius: 8px;
		background: #232d25;
	}
	.time-button:hover,
	.time-button.searching {
		border-color: #88ab8b;
		background: #304735;
	}
	.time-value {
		font-size: clamp(1.125rem, 3.5vw, 1.375rem);
		font-weight: 650;
		line-height: 1.15;
		font-variant-numeric: tabular-nums;
	}
	.time-value > span,
	.time-format {
		color: var(--muted);
		font-weight: 400;
	}
	.time-format {
		font-size: 0.8125rem;
		line-height: 1.15;
	}
	.queue-status:not(:empty) {
		display: flex;
		align-items: center;
		justify-content: space-between;
		gap: 0.5rem;
		margin-top: 0.5rem;
		color: var(--accent);
		font-size: 0.875rem;
	}
	.queue-status button {
		min-height: 2.125rem;
	}
	.challenge-panel {
		display: grid;
		gap: 0.875rem;
		scroll-margin-top: 0.75rem;
	}
	.range-field label {
		display: flex;
		justify-content: space-between;
		margin-bottom: 0.25rem;
		color: var(--muted);
		font-size: 0.875rem;
	}
	.range-field strong {
		color: var(--accent);
	}
	.challenge-panel .button {
		min-height: 2.5rem;
		padding: 0.5rem;
	}
	@media (max-width: 639px) {
		.play-layout {
			grid-template-columns: minmax(0, 1fr);
		}
		.time-grid {
			gap: 0.375rem;
		}
		.time-button {
			min-height: 3.375rem;
		}
		.time-value {
			font-size: 1rem;
			white-space: nowrap;
		}
		.time-format {
			font-size: 0.75rem;
		}
		.jump-link {
			display: inline;
		}
	}
	@media (max-height: 440px) {
		.page-title {
			margin-bottom: 0.5rem;
		}
		h1 {
			line-height: 1.15;
		}
		.panel-heading {
			margin-bottom: 0.375rem;
		}
		.filters button {
			min-height: 1.875rem;
			padding-inline: 0.5rem;
		}
		.time-grid {
			gap: 0.375rem;
		}
		.time-button {
			min-height: 2.75rem;
		}
		.time-value {
			font-size: 1.125rem;
		}
	}
</style>
