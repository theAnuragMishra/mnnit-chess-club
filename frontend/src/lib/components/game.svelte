<script lang="ts">
	import { Chess, type Square } from 'chess.js';
	import { page } from '$app/state';
	import Chat from '$lib/components/Chat.svelte';
	import DrawResign from '$lib/components/DrawResign.svelte';
	import GameInfo from '$lib/components/GameInfo.svelte';
	import Chessboard from '$lib/components/Chessboard.svelte';
	import Clock from '$lib/components/Clock.svelte';
	import { websocketStore } from '$lib/websocket';
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
	import playImg from '$lib/assets/icons/play.svg';
	import berserkImg from '$lib/assets/icons/kill.svg';

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
	let timeBlack = $state(d.timeBlack);
	let timeWhite = $state(d.timeWhite);
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
	const isPlayer =
		data.user.username === d.game.WhiteUsername || data.user.username === d.game.BlackUsername;
	const whiteUp = d.game.WhiteUsername !== data.user.username;
	const setActiveIndex = (index: number) => {
		if (index > moveHistory.length - 1 || index < -1 || index == activeIndex) return;
		activeIndex = index;
		ground?.cancelPremove();
		ground?.selectSquare(null);
		ground?.set({
			fen: chessForView.fen(),
			lastMove:
				activeIndex === -1 ? [] : [moveHistory[activeIndex].Orig, moveHistory[activeIndex].Dest],
			viewOnly: activeIndex !== moveHistory.length - 1,
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
						console.log('begin', performance.now());
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
		ground?.cancelPremove();
		ground?.selectSquare(null);
		ground?.set({
			viewOnly: true
		});
		notifyAudio.play();
	};

	const handleMoveResponse = (payload: any) => {
		//console.log(payload);
		console.log('finish', performance.now());
		moveHistory = [...moveHistory, payload.Move];
		timeBlack = payload.TimeBlack;
		timeWhite = payload.TimeWhite;
		ground?.set({
			turnColor: moveHistory.length % 2 == 0 ? 'white' : 'black',
			movable: { dests: getValidMoves(chessLatest) }
		});
		if (activeIndex === moveHistory.length - 2) {
			activeIndex = moveHistory.length - 1;
			ground?.set({
				fen: payload.Move.MoveFen,
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

<svelte:head>
	<title
		>{result === 0 ? 'Play' : 'Game Over'} vs
		{whiteUp ? whiteUsername : blackUsername} - {baseTime / 60} + {increment}</title
	>
</svelte:head>
<div class="flex w-full flex-col-reverse items-center justify-center gap-5 xl:flex-row">
	<div class="flex w-[95%] flex-col gap-5 md:flex-row xl:w-[300px] xl:flex-col">
		{#if tournamentID}
			<a
				class="flex h-fit w-fit items-center gap-[5px] text-2xl"
				href={`/tournament/${tournamentID}`}
			>
				<img class="h-[30px]" src={trophyImg} alt="trophy" />{tournamentName}</a
			>
		{/if}
		<GameInfo
			{whiteUsername}
			{blackUsername}
			{result}
			{createdAt}
			{baseTime}
			{increment}
			{method}
		/>
		<Chat hei="256px" />
	</div>
	<div class="acontainer w-full xl:w-3/4">
		<div class="abortt">
			{#if result === 0 && (whiteUp ? moveHistory.length == 0 : moveHistory.length == 1)}
				<AbortTimer {lowAbortTime} time={whiteUp ? abortTimeWhite : abortTimeBlack} tb="t" />
			{/if}
		</div>
		<div class="board flex w-full flex-col justify-center">
			<Chessboard {setGround} {boardConfig} />
		</div>
		<div class="abortb">
			{#if result === 0 && (whiteUp ? moveHistory.length == 1 : moveHistory.length == 0)}
				<AbortTimer {lowAbortTime} time={whiteUp ? abortTimeBlack : abortTimeWhite} tb="b" />
			{/if}
		</div>
		<div
			class="piecest flex h-[15px] items-center text-[10px] text-gray-400 md:h-[35px] md:w-[300px] md:flex-wrap-reverse md:text-[20px]"
		>
			<CapturedPieces pieces={whiteUp ? bCap : wCap} /><span class="ml-[6px]"
				>{whiteUp
					? balance > 0
						? `+${Math.abs(balance)}`
						: ''
					: balance < 0
						? `+${Math.abs(balance)}`
						: ''}
			</span>
		</div>
		<div class="clockt flex h-fit items-center gap-2">
			<Clock
				{lowTime}
				time={whiteUp ? wtime : btime}
				active={result !== 0 || moveHistory.length < 2
					? false
					: whiteUp
						? chessLatest.turn() === 'w'
						: chessLatest.turn() === 'b'}
			/>
			{#if whiteUp ? berserkWhite : berserkBlack}
				<img src={berserkImg} alt="berserk icon" class="h-[32px]" />
			{/if}
		</div>
		<div class="rematch flex w-full justify-center">
			{#if result !== 0 && isPlayer && !tournamentID}
				<Rematch canRematch={d.canRematch} />
			{/if}
		</div>
		<div class="namet flex h-fit justify-between md:w-[300px]">
			<a href={`/member/${whiteUp ? whiteUsername : blackUsername}`}
				>{whiteUp ? whiteUsername : blackUsername}&nbsp;&nbsp;</a
			><span
				>{whiteUp ? ratingWhite : ratingBlack}&nbsp;&nbsp;<span
					class={`${whiteUp ? (changeWhite > 0 ? 'text-green-500' : changeWhite == 0 ? 'text-gray-500' : 'text-red-500') : changeBlack > 0 ? 'text-green-500' : changeBlack == 0 ? 'text-gray-500' : 'text-red-500'}`}
					>{result !== 0
						? whiteUp
							? `${changeWhite > 0 ? '+' : ''}${changeWhite}`
							: `${changeBlack > 0 ? '+' : ''}${changeBlack}`
						: ''}</span
				></span
			>
		</div>
		<div class="draw-resign h-fit w-full">
			{#if isPlayer && result === 0}
				<DrawResign
					isDisabled={moveHistory.length < 2}
					{gameID}
					showBerserkButton={result === 0 &&
						d.berserkAllowed &&
						(whiteUp
							? !berserkBlack && moveHistory.length <= 1
							: !berserkWhite && moveHistory.length == 0)}
				/>
			{/if}
		</div>
		<div class="back_to_tournament w-full text-center">
			{#if isPlayer && tournamentID && result !== 0}
				<a
					class="flex w-full items-center justify-center gap-[5px] px-3 py-2"
					href={`/tournament/${tournamentID}`}
				>
					<img src={playImg} alt="back to tournament" /> BACK TO TOURNAMENT</a
				>
			{/if}
		</div>
		<div class="history w-full">
			<HistoryTable
				{moveHistory}
				{setActiveIndex}
				{activeIndex}
				highlightLastArrow={activeIndex !== moveHistory.length - 1 &&
					(whiteUp ? chessLatest.turn() === 'b' : chessLatest.turn() === 'w') &&
					result === 0}
			/>
		</div>

		<div class="nameb flex h-fit justify-between md:w-[300px]">
			<a href={`/member/${whiteUp ? blackUsername : whiteUsername}`}
				>{whiteUp ? blackUsername : whiteUsername}&nbsp;&nbsp;</a
			>
			<span
				>{whiteUp ? ratingBlack : ratingWhite}&nbsp;&nbsp;<span
					class={`${!whiteUp ? (changeWhite > 0 ? 'text-green-500' : changeWhite == 0 ? 'text-gray-500' : 'text-red-500') : changeBlack > 0 ? 'text-green-500' : changeBlack == 0 ? 'text-gray-500' : 'text-red-500'}`}
					>{result !== 0
						? whiteUp
							? `${changeBlack > 0 ? '+' : ''}${changeBlack}`
							: `${changeWhite > 0 ? '+' : ''}${changeWhite}`
						: ''}</span
				></span
			>
		</div>
		<div class="clockb flex h-fit items-center gap-2">
			<Clock
				{lowTime}
				time={whiteUp ? btime : wtime}
				active={result !== 0 || moveHistory.length < 2
					? false
					: whiteUp
						? chessLatest.turn() === 'b'
						: chessLatest.turn() === 'w'}
			/>
			{#if whiteUp ? berserkBlack : berserkWhite}
				<img src={berserkImg} alt="berserk icon" class="h-[32px]" />
			{/if}
		</div>
		<div
			class="piecesb flex h-[15px] items-center text-[10px] text-gray-400 md:h-[35px] md:w-[300px] md:flex-wrap md:text-[20px]"
		>
			<CapturedPieces pieces={whiteUp ? wCap : bCap} />
			<span class="ml-[6px]"
				>{whiteUp
					? balance < 0
						? `+${Math.abs(balance)}`
						: ''
					: balance > 0
						? `+${Math.abs(balance)}`
						: ''}
			</span>
		</div>
	</div>
</div>

<style>
	.acontainer {
		display: grid;
		row-gap: 1px;
		/* grid-template-columns: auto auto;
		grid-template-rows: auto auto auto auto auto; */
		grid-template-areas:
			'namet clockt'
			'piecest clockt'
			'abortt abortt'
			'board board'
			'abortb abortb'
			'piecesb clockb'
			'nameb clockb'
			'btt btt'
			'draw-resign draw-resign'
			'rematch rematch'
			'history history';
	}
	.namet {
		grid-area: namet;
		justify-self: start;
		align-self: center;
		margin-left: 10px;
	}
	.nameb {
		grid-area: nameb;
		justify-self: start;
		align-self: center;
		margin-left: 10px;
	}
	.board {
		grid-area: board;
		justify-self: center;
	}
	.draw-resign {
		grid-area: draw-resign;
	}
	.clockt {
		grid-area: clockt;
		justify-self: end;
		flex-direction: row-reverse;
	}
	.clockb {
		grid-area: clockb;
		justify-self: end;
		flex-direction: row-reverse;
	}
	.history {
		grid-area: history;
		justify-self: center;
	}
	.back_to_tournament {
		grid-area: btt;
		justify-self: center;
	}
	.abortt {
		grid-area: abortt;
		display: none;
	}
	.abortb {
		grid-area: abortb;
	}
	.rematch {
		grid-area: rematch;
	}

	.piecest {
		grid-area: piecest;
		align-self: center;
		margin-left: 8px;
	}
	.piecesb {
		grid-area: piecesb;
		align-self: center;
		margin-left: 8px;
	}
	@media (width>= 768px) {
		.acontainer {
			row-gap: 3px;
			column-gap: 10px;
			place-content: center;
			place-items: start;
			grid-template-areas:
				'board .'
				'board piecest'
				'board clockt'
				'board namet'
				'board abortt'
				'board btt'
				'board history'
				'board draw-resign'
				'board rematch'
				'board abortb'
				'board nameb'
				'board clockb'
				'board piecesb'
				'board .';
		}
		.nameb {
			justify-self: auto;
			margin: 0;
		}
		.namet {
			justify-self: auto;
			margin: 0;
		}
		.clockt {
			justify-self: auto;
			/* margin-top: auto; */
			flex-direction: row;
		}
		.clockb {
			justify-self: auto;
			flex-direction: row;
			/* margin-bottom: auto; */
		}
		.history {
			justify-self: auto;
		}
		.back_to_tournament {
			justify-self: auto;
		}
		.board {
			justify-self: auto;
		}
		.abortt {
			display: unset;
		}
		.piecesb {
			margin: 0;
		}
		.piecest {
			margin: 0;
		}
	}
</style>
