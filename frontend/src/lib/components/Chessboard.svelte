<script lang="ts">
	import { Chessground } from '@lichess-org/chessground';
	import '../../../node_modules/@lichess-org/chessground/assets/chessground.base.css';
	import '../../../node_modules/@lichess-org/chessground/assets/chessground.brown.css';
	import '../../../node_modules/@lichess-org/chessground/assets/chessground.cburnett.css';
	import { onMount } from 'svelte';

	let boardContainer: HTMLDivElement;

	let { setGround, boardConfig } = $props();

	onMount(() => {
		if (!boardContainer) return;
		const cg = Chessground(boardContainer, boardConfig);
		setGround(cg);

		return () => cg.destroy();
	});
</script>

<div bind:this={boardContainer} class="bcontainer"></div>

<style>
	.bcontainer {
		width: 100%;
		aspect-ratio: 1;
	}
	@media (width >= 768px) {
		.bcontainer {
			height: 60vw;
			width: 60vw;
		}
		.bcontainer :global(cg-board) {
			border-radius: 6px !important;
		}
	}
	@media (width >= 64rem) {
		.bcontainer {
			height: 650px;
			width: 650px;
		}
	}
</style>
