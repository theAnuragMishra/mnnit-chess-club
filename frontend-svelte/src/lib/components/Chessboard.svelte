<script lang="ts">
	import { Chessground } from 'chessground';
	import '../../node_modules/chessground/assets/chessground.base.css';
	import '../../node_modules/chessground/assets/chessground.brown.css';
	import '../../node_modules/chessground/assets/chessground.cburnett.css';
	import { getValidMoves, isPromoting } from '$lib/utils';
	import { websocketStore } from '$lib/websocket';
	import type { Square } from 'chess.js';

	const sendMessage = websocketStore.sendMessage;

	let boardContainer: HTMLDivElement;

	let { username, gameID, chess, white, lastMove, viewOnly } = $props();

	$effect(() => {
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
							sendMessage({
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
							sendMessage({
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

		return () => cg.destroy(); // Cleanup on unmount
	});
</script>

<div bind:this={boardContainer} style="width: 400px; height: 400px;"></div>
