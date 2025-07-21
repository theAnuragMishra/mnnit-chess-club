<script lang="ts">
	import { Chessground } from 'chessground';
	import '../../../node_modules/chessground/assets/chessground.base.css';
	import '../../../node_modules/chessground/assets/chessground.brown.css';
	import '../../../node_modules/chessground/assets/chessground.cburnett.css';

	let boardContainer: HTMLDivElement;

	let { setGround, boardConfig } = $props();

	$effect(() => {
		console.log('effect ran');

		if (!boardContainer) return;
		const cg = Chessground(boardContainer, boardConfig);
		setGround(cg);

		return () => cg.destroy();
	});
</script>

<div bind:this={boardContainer} class="bcontainer"></div>

<style>
	.bcontainer :global(cg-board) {
		border-radius: 6px !important;
	}
	.bcontainer {
		height: 80vw;
		width: 80vw;
	}
	@media (width >= 768px) {
		.bcontainer {
			height: 60vw;
			width: 60vw;
		}
	}
	@media (width >= 64rem) {
		.bcontainer {
			height: 650px;
			width: 650px;
		}
	}
</style>
