<script lang="ts">
	import { onMount } from 'svelte';
	import { Chessground } from 'chessground';
	import '../../node_modules/chessground/assets/chessground.base.css';
	import '../../node_modules/chessground/assets/chessground.brown.css';
	import '../../node_modules/chessground/assets/chessground.cburnett.css';

	let boardContainer: HTMLDivElement;

	onMount(() => {
		const cg = Chessground(boardContainer, {
			orientation: 'white',
			movable: {
				free: false,
				color: 'white',
				dests: new Map()
			},
			highlight: {
				lastMove: true,
				check: true
			},
			animation: {
				duration: 300
			}
		});

		return () => cg.destroy(); // Cleanup on unmount
	});
</script>

<div bind:this={boardContainer} style="width: 400px; height: 400px;"></div>
