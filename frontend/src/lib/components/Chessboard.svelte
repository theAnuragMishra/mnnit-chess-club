<script lang="ts">
	import { Chessground } from 'chessground';
	import { getValidMoves, isPromoting } from '$lib/utils';
	import { websocketStore } from '$lib/websocket';
	import type { Square } from 'chess.js';
	import '../../../node_modules/chessground/assets/chessground.base.css';
	import '../../../node_modules/chessground/assets/chessground.brown.css';
	import '../../../node_modules/chessground/assets/chessground.cburnett.css';

	let boardContainer: HTMLDivElement;

	let { orientation, gameID, chess, lastMove, viewOnly, setGround } = $props();

	$effect(() => {
		// console.log('effect ran');
		// console.log(username, chess.ascii(), lastMove, viewOnly);
		if (!boardContainer) return;
		const cg = Chessground(boardContainer, {
			fen: chess.fen(),
			orientation,
			draggable: { enabled: true },
			turnColor: chess.turn() == 'w' ? 'white' : 'black',
			viewOnly: viewOnly,
			lastMove: lastMove,
			check: chess.isCheck(),
			movable: {
				free: false,
				color: orientation,
				dests: getValidMoves(chess),
				showDests: true,
				events: {
					after: (orig, dest) => {
						const piece = chess.get(orig as Square);
						if (isPromoting(dest, piece!)) {
							const move = chess.move({
								from: orig,
								to: dest,
								promotion: 'q'
							});
							websocketStore.sendMessage({
								type: 'move',
								payload: {
									MoveStr: move.san,
									orig: orig,
									dest: dest,
									GameID: gameID
								}
							});
						} else {
							const move = chess.move({ from: orig, to: dest });
							websocketStore.sendMessage({
								type: 'move',
								payload: {
									MoveStr: move.san,
									orig: orig,
									dest: dest,
									GameID: gameID
								}
							});
						}
					}
				}
			},
			highlight: { lastMove: true, check: true }
		});
		setGround(cg);

		return () => cg.destroy();
	});
</script>

<div bind:this={boardContainer} class="bcontainer"></div>

<style>
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
