<script lang="ts">
	import { Chess } from 'chess.js';
	import Icon from './Icon.svelte';
	const pieces = import.meta.glob('$lib/assets/pieces/*.svg', {
		eager: true,
		query: '?url',
		import: 'default'
	}) as Record<string, string>;
	const position = new Chess();
	for (const move of ['e4', 'e5', 'Nf3', 'Nc6', 'Bc4', 'Bc5']) position.move(move);
	const board = position.board();
	function pieceImage(color: string, type: string) {
		return Object.entries(pieces).find(([path]) =>
			path.endsWith(`/${color}${type.toUpperCase()}.svg`)
		)?.[1];
	}
</script>

<div class="board-scene">
	<div class="board-card">
		<div class="board-header">
			<span class="player-avatar"><Icon name="knight" size={21} /></span>
			<strong>Italian Game</strong>
			<span class="board-time">05:00</span>
		</div>
		<div
			class="preview-board"
			role="img"
			aria-label="Chessboard showing the Italian opening after three moves by each player, with white to move"
		>
			{#each board as row, rank}{#each row as piece, file}<div
						class="square"
						class:dark={(rank + file) % 2 === 1}
						class:last-move={(rank === 0 && file === 5) || (rank === 3 && file === 2)}
					>
						{#if file === 0}<span class="rank">{8 - rank}</span>{/if}{#if rank === 7}<span
								class="file">{'abcdefgh'[file]}</span
							>{/if}{#if piece}<img
								src={pieceImage(piece.color, piece.type)}
								alt=""
								draggable="false"
							/>{/if}
					</div>{/each}{/each}
		</div>
		<div class="board-bottom">
			<span class="white-piece"></span><span>White to move</span>
		</div>
	</div>
</div>

<style>
	.board-scene {
		position: relative;
		width: 100%;
		max-width: 454px;
		margin: 0 auto;
		padding: 0 0 17px;
	}
	.board-scene::before {
		position: absolute;
		z-index: -1;
		inset: 14% -3% -5%;
		content: '';
		background: radial-gradient(ellipse, #5a8d4733 0, transparent 68%);
	}
	.board-card {
		padding: 18px;
		border: 1px solid #405045;
		border-radius: 15px;
		background: #202923;
		box-shadow: 0 24px 60px #0004;
		transform: rotate(-3deg);
	}
	.board-header {
		display: flex;
		align-items: center;
		gap: 9px;
		padding: 0 1px 15px;
	}
	.player-avatar {
		display: grid;
		width: 35px;
		height: 35px;
		place-items: center;
		border: 1px solid #4b6050;
		border-radius: 50%;
		color: var(--accent);
		background: #2c3b2e;
	}
	.board-header strong {
		display: block;
		font-size: 0.875rem;
		font-weight: 550;
	}
	.board-time {
		margin-left: auto;
		padding: 5px 9px;
		border-radius: 5px;
		background: #111b14;
		color: #dce9d9;
		font-size: 1.125rem;
		font-variant-numeric: tabular-nums;
		font-weight: 600;
	}
	.preview-board {
		display: grid;
		grid-template-columns: repeat(8, 1fr);
		overflow: hidden;
		border-radius: 4px;
		aspect-ratio: 1;
	}
	.square {
		position: relative;
		aspect-ratio: 1;
		background: #e4e9cf;
		color: #5a7854;
	}
	.square.dark {
		background: #71916a;
		color: #edf1d7;
	}
	.square.last-move {
		background: #b6cc7b;
	}
	.square.dark.last-move {
		background: #a1b966;
	}
	.square img {
		position: relative;
		z-index: 1;
		width: 100%;
		height: 100%;
		padding: 2px;
		filter: drop-shadow(0 2px 1px #0003);
	}
	.rank,
	.file {
		position: absolute;
		z-index: 2;
		font-size: 0.5625rem;
		font-weight: 650;
	}
	.rank {
		top: 2px;
		left: 3px;
	}
	.file {
		right: 3px;
		bottom: 1px;
	}
	.board-bottom {
		display: flex;
		align-items: center;
		gap: 7px;
		padding-top: 13px;
		color: #b2bfb4;
		font-size: 0.75rem;
	}
	.white-piece {
		width: 9px;
		height: 9px;
		border-radius: 50%;
		background: #e4e9cf;
	}
	@media (max-width: 700px) {
		.board-scene {
			max-width: 370px;
		}
		.board-header strong {
			font-size: 0.75rem;
		}
	}
</style>
