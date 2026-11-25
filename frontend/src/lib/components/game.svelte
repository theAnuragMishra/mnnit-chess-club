<script lang="ts">
	import { Chess, type Square } from 'chess.js';
	import { page } from '$app/state';
	import Chat from '$lib/components/Chat.svelte';
	import DrawResign from '$lib/components/DrawResign.svelte';
	import GameInfo from '$lib/components/GameInfo.svelte';
	import Chessboard from '$lib/components/Chessboard.svelte';
	import Clock from '$lib/components/Clock.svelte';
	import HistoryTable from '$lib/components/HistoryTable.svelte';
	import { onDestroy, onMount } from 'svelte';
	import type { Api } from '@lichess-org/chessground/api';
	import { getValidMoves, isPromoting } from '$lib/utils.js';
	import AbortTimer from '$lib/components/AbortTimer.svelte';
	import type { Config } from '@lichess-org/chessground/config';
	import Rematch from './Rematch.svelte';
	import { notifyAudio, moveAudio, captureAudio, lowTimeAudio, berserkAudio } from '$lib/audios';
	import { getCapturedAndMaterial } from '$lib/chessUtils';

	import CapturedPieces from './CapturedPieces.svelte';
	import trophyImg from '$lib/assets/icons/trophy.svg';
	import berserkImg from '$lib/assets/icons/kill.svg';
	import { user } from '$lib/user.svelte';
	import { websocketStore } from '$lib/websocket.svelte';

	let { data } = $props();
	//console.log(data);
	let d = data.gameData;
	let ground: Api | null = $state(null);
	const baseTime = d.game.BaseTime;
	const increment = d.game.Increment;
	const createdAt = new Date(d.game.CreatedAt);
	let tournamentName = d.game.TournamentName;
	let tournamentID = d.game.TournamentID;
	let whiteUsername = d.game.WhiteUsername;
	let blackUsername = d.game.BlackUsername;
	let timeBlack = $state(d.game.TimeBlack);
	let timeWhite = $state(d.game.TimeWhite);
	let ratingWhite = d.game.RatingW;
	let ratingBlack = d.game.RatingB;
	let changeWhite = $state(d.game.ChangeW ?? 0);
	let changeBlack = $state(d.game.ChangeB ?? 0);
	let result = $state(d.game.Result);
	let method = $state(d.game.Method);
	let moveHistory = $state(d.moves);
	let activeIndex = $state(d.moves.length - 1);
	let chessLatest = $derived(
		moveHistory.length > 0 ? new Chess(moveHistory[moveHistory.length - 1].MoveFen) : new Chess()
	);
	let chessForView = $derived(
		activeIndex !== -1 ? new Chess(moveHistory[activeIndex].MoveFen) : new Chess()
	);
	const isPlayer = user.username === d.game.WhiteUsername || user.username === d.game.BlackUsername;
	const whiteUp = d.game.WhiteUsername !== user.username;
	let flipped = $state(false);
	const topIsWhite = $derived(flipped ? !whiteUp : whiteUp);
	function flipBoard() {
		if (!ground) return;
		ground.toggleOrientation();
		flipped = !flipped;
	}
	const setActiveIndex = (index: number) => {
		if (index > moveHistory.length - 1 || index < -1 || index == activeIndex) return;
		activeIndex = index;
		ground?.cancelPremove();
		ground?.selectSquare(null);
		if (activeIndex === moveHistory.length - 1) {
			ground?.set({
				movable: { dests: getValidMoves(chessLatest) }
			});
		}
		ground?.set({
			fen: chessForView.fen(),
			turnColor: activeIndex % 2 == 0 ? 'black' : 'white',
			lastMove:
				activeIndex === -1 ? [] : [moveHistory[activeIndex].Orig, moveHistory[activeIndex].Dest],
			viewOnly: !isPlayer || result !== 0 || activeIndex !== moveHistory.length - 1,
			check: chessForView.isCheck()
		});
	};
	const gameID = page.params.gameID;

	const setGround = (g: Api) => {
		ground = g;
	};

	//board config
	const boardConfig: Config = {
		fen: d.moves.length > 0 ? d.moves[d.moves.length - 1].MoveFen : undefined,
		orientation: whiteUp ? 'black' : 'white',
		draggable: { enabled: true },
		turnColor: d.moves.length % 2 == 0 ? 'white' : 'black',
		viewOnly: !isPlayer || d.game.Result !== 0,
		lastMove:
			d.moves.length > 0
				? [d.moves[d.moves.length - 1].Orig, d.moves[d.moves.length - 1].Dest]
				: [],
		check: d.moves.length > 0 ? new Chess(d.moves[d.moves.length - 1].MoveFen).isCheck() : false,
		movable: {
			free: false,
			color: whiteUp ? 'black' : 'white',
			dests: getValidMoves(
				d.moves.length > 0 ? new Chess(d.moves[d.moves.length - 1].MoveFen) : new Chess()
			),
			showDests: true,
			events: {
				after: (orig, dest) => {
					const piece = chessForView.get(orig as Square);
					if (isPromoting(dest, piece!)) {
						const move = chessForView.move({
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
						const move = chessForView.move({ from: orig, to: dest });
						websocketStore.sendMessage({
							type: 'move',
							payload: {
								MoveStr: move.san,
								orig: orig,
								dest: dest,
								GameID: gameID
							}
						});
						//console.log('begin', performance.now());
					}
				}
			}
		},
		highlight: { lastMove: true, check: true }
	};

	const handleGameEnd = (payload: any) => {
		result = payload.Result;
		method = payload.Method;
		changeBlack = payload.changeB;
		changeWhite = payload.changeW;
		timeWhite = payload.timeWhite;
		timeBlack = payload.timeBlack;
		ground?.stop();
		ground?.set({
			viewOnly: true
		});
		notifyAudio.play();
	};

	const handleMoveResponse = (payload: any) => {
		//console.log(payload);
		//console.log('finish', performance.now());
		moveHistory = [...moveHistory, payload.Move];
		timeBlack = payload.TimeBlack;
		timeWhite = payload.TimeWhite;

		if (activeIndex === moveHistory.length - 2) {
			activeIndex = moveHistory.length - 1;
			ground?.set({
				fen: payload.Move.MoveFen,
				turnColor: moveHistory.length % 2 == 0 ? 'white' : 'black',
				movable: { dests: getValidMoves(chessLatest) },
				check: chessForView.isCheck(),
				lastMove: [
					moveHistory[moveHistory.length - 1].Orig,
					moveHistory[moveHistory.length - 1].Dest
				]
			});
		}

		//console.log(ground?.state.movable);
		if (payload.Move.MoveNotation[1] == 'x') captureAudio?.play();
		else moveAudio?.play();
		ground?.playPremove();
	};

	let berserkWhite = $state(d.game.BerserkWhite);
	let berserkBlack = $state(d.game.BerserkBlack);
	const handleBerserk = (payload: any) => {
		if (payload.wb == 0) {
			timeWhite /= 2;
			berserkWhite = true;
		} else if (payload.wb == 1) {
			timeBlack /= 2;
			berserkBlack = true;
		}
		berserkAudio.play();
	};

	//timer setup
	const abortLength = baseTime >= 20 ? 20000 : baseTime >= 10 ? 10000 : baseTime * 1000;
	const lowAbortTime = baseTime >= 20 ? 10000 : baseTime >= 10 ? 5000 : 2000;
	const lowTime =
		(baseTime >= 1800
			? 120
			: baseTime >= 600
				? 60
				: baseTime >= 300
					? 30
					: baseTime >= 30
						? 10
						: baseTime >= 10
							? 3
							: 1) * 1000;
	let btime = $derived(
		result !== 0
			? activeIndex == -1 || activeIndex == 0
				? (baseTime * 1000) / (berserkBlack ? 2 : 1)
				: activeIndex == moveHistory.length - 1
					? timeBlack
					: activeIndex % 2 == 0
						? moveHistory[activeIndex - 1].TimeLeft
						: moveHistory[activeIndex].TimeLeft
			: timeBlack
	);
	let wtime = $derived(
		result !== 0
			? activeIndex == -1
				? (baseTime * 1000) / (berserkWhite ? 2 : 1)
				: activeIndex == moveHistory.length - 1
					? timeWhite
					: activeIndex % 2 == 0
						? moveHistory[activeIndex].TimeLeft
						: moveHistory[activeIndex - 1].TimeLeft
			: timeWhite
	);
	let abortTimeWhite = $state(abortLength);
	let abortTimeBlack = $state(abortLength);
	let animationFrame: number | null;
	let startTime: DOMHighResTimeStamp | null;
	let lowTimePlayed = $state(false);

	$effect(() => {
		if (result === 0) {
			let trn = chessLatest.turn();
			startTime = performance.now();
			const tick = (currentTime: number) => {
				if (!startTime) return;
				const elapsed = currentTime - startTime;
				if (moveHistory.length >= 2) {
					const newTime = (trn == 'w' ? timeWhite : timeBlack) - elapsed;

					if (newTime <= 0) {
						if (trn == 'w') wtime = 0;
						else btime = 0;
						return;
					}
					if (trn == 'w') wtime = newTime;
					else btime = newTime;

					if ((whiteUp && trn == 'b') || (!whiteUp && trn == 'w')) {
						if (!lowTimePlayed && newTime <= lowTime) {
							lowTimePlayed = true;
							lowTimeAudio?.play();
						}
					}
				} else {
					const newTime = abortLength - elapsed;
					if (newTime <= 0) {
						if (trn == 'w') abortTimeWhite = 0;
						else abortTimeBlack = 0;
						return;
					}
					if (trn == 'w') abortTimeWhite = newTime;
					else abortTimeBlack = newTime;
				}
				animationFrame = requestAnimationFrame(tick);
			};

			animationFrame = requestAnimationFrame(tick);
		} else {
			if (animationFrame !== null) {
				cancelAnimationFrame(animationFrame);
				animationFrame = null;
				startTime = null;
			}
		}

		return () => {
			if (animationFrame !== null) {
				cancelAnimationFrame(animationFrame);
				animationFrame = null;
				startTime = null;
			}
		};
	});

	//captured pieces
	let { white: wCap, black: bCap, balance } = $derived(getCapturedAndMaterial(chessForView));
	//$inspect(balance);

	onMount(() => {
		websocketStore.onMessage('game_end', handleGameEnd);
		websocketStore.onMessage('Move_Response', handleMoveResponse);
		websocketStore.onMessage('berserk', handleBerserk);
	});
	onDestroy(() => {
		websocketStore.offMessage('Move_Response', handleMoveResponse);
		websocketStore.offMessage('resignation', handleGameEnd);
	});
</script>

{#snippet playerRow(isWhite: boolean)}
	{@const username = isWhite ? whiteUsername : blackUsername}
	{@const rating = isWhite ? ratingWhite : ratingBlack}
	{@const ratingChange = isWhite ? changeWhite : changeBlack}
	{@const material = isWhite ? balance : -balance}
	{@const captured = isWhite ? bCap : wCap}
	{@const berserk = isWhite ? berserkWhite : berserkBlack}
	<div class="player-row">
		<div class="player-info">
			<div class="player-name">
				<span
					class="color-chip"
					class:white={isWhite}
					aria-label={isWhite ? 'White' : 'Black'}
					role="img"
				></span>
				<a href={`/member/${username}`}>{username}</a>
				<span class="rating">{rating}</span>
				{#if result !== 0 && ratingChange !== 0}
					<span class="rating-change" class:negative={ratingChange < 0}
						>{ratingChange > 0 ? '+' : ''}{ratingChange}</span
					>
				{/if}
			</div>
			{#if captured.length}
				<div class="captured">
					<CapturedPieces pieces={captured} />
					{#if material > 0}<span>+{material}</span>{/if}
				</div>
			{/if}
		</div>
		<div class="player-clock">
			{#if berserk}<img src={berserkImg} alt="Berserk" width="18" height="18" />{/if}
			<Clock
				{lowTime}
				time={isWhite ? wtime : btime}
				active={result === 0 &&
					moveHistory.length >= 2 &&
					chessLatest.turn() === (isWhite ? 'w' : 'b')}
				compact
			/>
		</div>
	</div>
{/snippet}

<svelte:head>
	<title>{whiteUsername} vs {blackUsername} · {baseTime / 60}+{increment} — MNNIT Chess Club</title>
</svelte:head>

<h1 class="sr-only">{whiteUsername} vs {blackUsername}</h1>
<div class="game-layout">
	<div class="board-column">
		{@render playerRow(topIsWhite)}
		<Chessboard {setGround} {boardConfig} />
		{@render playerRow(!topIsWhite)}
	</div>

	<aside class="game-rail" aria-label="Game controls">
		{#if result === 0 && moveHistory.length < 2}
			<AbortTimer
				{lowAbortTime}
				time={moveHistory.length === 0 ? abortTimeWhite : abortTimeBlack}
				player={moveHistory.length === 0 ? whiteUsername : blackUsername}
				yourTurn={user.username === (moveHistory.length === 0 ? whiteUsername : blackUsername)}
			/>
		{/if}

		<HistoryTable
			{moveHistory}
			{setActiveIndex}
			{activeIndex}
			highlightLastArrow={activeIndex !== moveHistory.length - 1 &&
				chessLatest.turn() === (whiteUp ? 'b' : 'w') &&
				result === 0}
		/>

		{#if isPlayer && result === 0}
			<DrawResign
				drawOfferedBy={d.drawOfferedBy}
				handleFlip={flipBoard}
				isDisabled={moveHistory.length < 2}
				{gameID}
				showBerserkButton={d.berserkAllowed &&
					(whiteUp
						? !berserkBlack && moveHistory.length <= 1
						: !berserkWhite && moveHistory.length === 0)}
			/>
		{:else if isPlayer && !tournamentID}
			<Rematch canRematch={d.canRematch} offeredBy={d.rematchOfferedBy ?? 0} />
		{/if}

		<GameInfo {result} {createdAt} {baseTime} {increment} {method} />

		{#if tournamentID}
			<a class="button button-secondary tournament-link" href={`/tournament/${tournamentID}`}>
				<img src={trophyImg} alt="" width="18" height="18" />
				<span>{tournamentName || 'Tournament'}</span>
			</a>
		{/if}

		<details class="game-chat panel">
			<summary>Chat</summary>
			<Chat hei="clamp(5rem, 24dvh, 12rem)" />
		</details>
	</aside>
</div>

<style>
	.game-layout {
		--game-padding: 0.75rem;
		--board-gap: 0.375rem;
		--player-height: 2.75rem;
		--rail-width: clamp(13rem, 28vw, 18rem);
		--board-limit: max(
			10rem,
			calc(
				100dvh - var(--site-header-height) - 2 *
					(var(--game-padding) + var(--board-gap) + var(--player-height))
			)
		);
		display: grid;
		width: min(
			calc(100% - 2 * var(--game-padding)),
			calc(var(--board-limit) + var(--rail-width) + 1rem)
		);
		grid-template-columns: minmax(0, 1fr) var(--rail-width);
		align-items: start;
		gap: 1rem;
		margin-inline: auto;
		padding-block: var(--game-padding);
	}
	.board-column,
	.game-rail {
		display: grid;
		grid-template-columns: minmax(0, 1fr);
		min-width: 0;
		gap: var(--board-gap);
	}
	.game-rail {
		gap: 0.5rem;
	}
	.player-row {
		display: flex;
		height: var(--player-height);
		align-items: center;
		justify-content: space-between;
		gap: 0.5rem;
		padding-inline: 0.625rem 0.25rem;
		border: 1px solid var(--line);
		border-radius: 8px;
		background: var(--surface);
	}
	.player-info,
	.player-name {
		min-width: 0;
	}
	.player-name,
	.player-clock,
	.captured {
		display: flex;
		align-items: center;
		gap: 0.375rem;
	}
	.player-name {
		font-size: 0.875rem;
		line-height: 1.2;
	}
	.player-name a {
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
		font-weight: 650;
	}
	.player-name a:hover {
		color: var(--accent);
	}
	.color-chip {
		width: 0.625rem;
		height: 0.625rem;
		flex: none;
		border: 1px solid #8a968b;
		border-radius: 50%;
		background: #141916;
	}
	.color-chip.white {
		background: #edf1e5;
	}
	.rating {
		color: var(--muted);
	}
	.rating-change {
		color: var(--accent);
	}
	.rating-change.negative {
		color: #f7a59d;
	}
	.captured {
		overflow: hidden;
		gap: 0;
		height: 1rem;
		color: var(--muted);
		font-size: 0.75rem;
	}
	.captured span {
		margin-left: 0.2rem;
	}
	.player-clock {
		flex: none;
	}
	.game-chat {
		border-radius: 8px;
	}
	.game-chat summary {
		padding: 0.625rem 0.875rem;
		font-size: 0.875rem;
	}
	.tournament-link {
		min-height: 2.5rem;
		padding: 0.5rem;
	}
	.tournament-link img {
		flex: none;
	}
	.tournament-link span {
		overflow-wrap: anywhere;
	}
	@media (max-width: 639px) {
		.game-layout {
			/* Leave room for the opening countdown below the board on short phones. */
			width: min(
				calc(100% - 2 * var(--game-padding)),
				max(10rem, calc(var(--board-limit) - 6.5rem))
			);
			grid-template-columns: minmax(0, 1fr);
			gap: 0.625rem;
		}
	}
</style>
