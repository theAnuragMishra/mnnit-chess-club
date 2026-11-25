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

<div bind:this={boardContainer} class="bcontainer" role="group" aria-label="Chessboard"></div>

<style>
	.bcontainer {
		width: 100%;
		aspect-ratio: 1;
	}
	.bcontainer :global(cg-board) {
		background-image: repeating-conic-gradient(#71916a 0% 25%, #e4e9cf 0% 50%);
		background-size: 25% 25%;
		border-radius: 5px;
	}
	.bcontainer :global(cg-board square.last-move) {
		background: #c6d25b80;
	}
	.bcontainer :global(cg-board square.selected) {
		background: #d1e57390;
	}
	.bcontainer :global(cg-board square.move-dest) {
		background: radial-gradient(#172c2860 19%, transparent 20%);
	}
	.bcontainer :global(cg-board square.move-dest:hover) {
		background: #c6d25b80;
	}
	.bcontainer :global(coords.ranks) {
		top: 0;
		left: 3px;
		color: #d9e5cd;
	}
	.bcontainer :global(coords.ranks coord) {
		transform: none;
		padding-top: 3px;
	}
	.bcontainer :global(coords.files) {
		left: 0;
		bottom: 0;
		height: 14px;
		color: #9cb190;
	}
	.bcontainer :global(coords.files coord) {
		text-align: right;
		padding-right: 3px;
	}
</style>
