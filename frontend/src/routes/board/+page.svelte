<script lang="ts">
	import Chessboard from '$lib/components/Chessboard.svelte';
	import type { Api } from '@lichess-org/chessground/api';

	let g: Api | null = $state(null);
</script>

<svelte:head>
	<title>Practice board — MNNIT Chess Club</title>
</svelte:head>
<h1 class="sr-only">Practice board</h1>
<div class="practice-board">
	<Chessboard boardConfig={{ viewOnly: true }} setGround={(x: Api) => (g = x)} />
	<div class="practice-actions">
		<button
			class="button button-secondary"
			onclick={() => {
				g?.set({ viewOnly: !g.state.viewOnly });
			}}>Toggle piece movement</button
		>
		<button
			class="button button-secondary"
			onclick={() => {
				g?.set({ drawable: { visible: !g.state.drawable.visible } });
			}}>Toggle board drawings</button
		>
	</div>
</div>

<style>
	.practice-board {
		--board-size: clamp(10rem, calc(100dvh - var(--site-header-height) - 2rem), 44rem);
		--actions-width: 13rem;
		display: grid;
		width: min(calc(100% - 2rem), calc(var(--board-size) + var(--actions-width) + 1rem));
		grid-template-columns: minmax(0, 1fr) var(--actions-width);
		align-items: start;
		gap: 1rem;
		margin: 1rem auto;
	}
	.practice-actions {
		display: grid;
		gap: 0.5rem;
	}
	@media (max-width: 639px) {
		.practice-board {
			width: min(calc(100% - 2rem), var(--board-size));
			grid-template-columns: minmax(0, 1fr);
		}
	}
</style>
