<script lang="ts">
	import { Chessground } from 'chessground';
	import { getValidMoves, isPromoting } from '$lib/utils';
	import { websocketStore } from '$lib/websocket';
	import type { Square } from 'chess.js';
	import '../../../node_modules/chessground/assets/chessground.base.css';
	import '../../../node_modules/chessground/assets/chessground.brown.css';
	import '../../../node_modules/chessground/assets/chessground.cburnett.css';

	let boardContainer: HTMLDivElement;

	let { username, gameID, chess, white, lastMove, viewOnly, setGround } = $props();

	$effect(() => {
		// console.log('effect ran');
		// console.log(username, chess.ascii(), lastMove, viewOnly);
		if (!boardContainer) return;
		const cg = Chessground(boardContainer, {
			fen: chess.fen(),
			orientation: white == username ? 'white' : 'black',
			draggable: { enabled: true },
			turnColor: chess.turn() == 'w' ? 'white' : 'black',
			viewOnly: viewOnly,
			lastMove: lastMove,
			check: chess.isCheck(),
			movable: {
				free: false,
				color: white == username ? 'white' : 'black',
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
									GameID: Number(gameID)
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
									GameID: Number(gameID)
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
		height: 90vw;
		width: 90vw;
	}
	@media (width >= 768px) {
		.bcontainer {
			height: 60vw;
			width: 60vw;
		}
	}
	@media (width >= 64rem) {
		.bcontainer {
			height: 630px;
			width: 630px;
		}
	}
</style>
