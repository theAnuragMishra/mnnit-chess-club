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
	import { notifyAudio, moveAudio, captureAudio, lowTimeAudio } from '$lib/audios';

	let { data } = $props();
	// console.log(data);
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
	let reason = $state(d.game.ResultReason);
	let moveHistory = $state(d.moves);
	let activeIndex = $state(d.moves ? d.moves.length - 1 : -1);
	let chessLatest = $derived(
		moveHistory ? new Chess(moveHistory[moveHistory.length - 1].MoveFen) : new Chess()
	);
	let chessForView = $derived(
		moveHistory && activeIndex !== -1 ? new Chess(moveHistory[activeIndex].MoveFen) : new Chess()
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
			viewOnly: moveHistory && activeIndex !== moveHistory.length - 1,
			check: chessForView.isCheck()
		});
	};
	const gameID = page.params.gameID;

	const setGround = (g: Api) => {
		ground = g;
	};

	//board config
	const boardConfig: Config = {
		fen: d.moves ? d.moves[d.moves.length - 1].MoveFen : undefined,
		orientation: whiteUp ? 'black' : 'white',
		draggable: { enabled: true },
		turnColor: d.moves ? (d.moves.length % 2 == 0 ? 'white' : 'black') : 'white',
		viewOnly: !isPlayer || d.game.Result !== 0,
		lastMove: d.moves ? [d.moves[d.moves.length - 1].Orig, d.moves[d.moves.length - 1].Dest] : [],
		check: d.moves ? new Chess(d.moves[d.moves.length - 1].MoveFen).isCheck() : false,
		movable: {
			free: false,
			color: whiteUp ? 'black' : 'white',
			dests: getValidMoves(d.moves ? new Chess(d.moves[d.moves.length - 1].MoveFen) : new Chess()),
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
					}
				}
			}
		},
		highlight: { lastMove: true, check: true }
	};

	const handleGameEnd = (payload: any) => {
		result = payload.Result;
		reason = payload.Reason;
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
		if (moveHistory) moveHistory = [...moveHistory, payload.move];
		else moveHistory = [payload.move];
		timeBlack = payload.timeBlack;
		timeWhite = payload.timeWhite;
		ground?.set({
			turnColor: moveHistory.length % 2 == 0 ? 'white' : 'black',
			movable: { dests: getValidMoves(chessLatest) }
		});
		if (activeIndex === moveHistory.length - 2) {
			activeIndex = moveHistory.length - 1;
			ground?.set({
				fen: payload.move.MoveFen,
				check: chessForView.isCheck(),
				lastMove: [
					moveHistory[moveHistory.length - 1].Orig,
					moveHistory[moveHistory.length - 1].Dest
				]
			});
		}

		//console.log(ground?.state.movable);
		if (payload.Result !== 0) {
			result = payload.Result;
			reason = payload.reason;
			changeBlack = payload.changeB;
			changeWhite = payload.changeW;
			ground?.set({
				viewOnly: true
			});
			notifyAudio.play();
			return;
		}

		if (payload.move.MoveNotation[1] == 'x') captureAudio?.play();
		else moveAudio?.play();
		ground?.playPremove();
	};

	//timer setup
	const abortLength = baseTime >= 20 ? 20 : baseTime >= 10 ? 10 : baseTime;
	const lowAbortTime = baseTime >= 20 ? 10 : baseTime >= 10 ? 5 : 2;
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
				? baseTime * 1000
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
				? baseTime * 1000
				: activeIndex == moveHistory.length - 1
					? timeWhite
					: activeIndex % 2 == 0
						? moveHistory[activeIndex].TimeLeft
						: moveHistory[activeIndex - 1].TimeLeft
			: timeWhite
	);
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

	onMount(() => {
		websocketStore.onMessage('game_end', handleGameEnd);
		websocketStore.onMessage('Move_Response', handleMoveResponse);
	});
	onDestroy(() => {
		websocketStore.offMessage('Move_Response', handleMoveResponse);
		websocketStore.offMessage('resignation', handleGameEnd);
	});
</script>

<div class="flex w-full flex-col-reverse items-center justify-center gap-5 xl:flex-row">
	<div class="flex w-[95%] flex-col gap-5 md:flex-row xl:w-[300px] xl:flex-col">
		{#if tournamentID}
			<a
				class="flex h-fit w-fit items-center gap-[5px] text-2xl"
				href={`/tournament/${tournamentID}`}
			>
				<svg
					class="h-[30px]"
					xmlns="http://www.w3.org/2000/svg"
					width="36"
					height="36"
					viewBox="0 0 36 36"
					><path
						fill="#ffac33"
						d="M5.123 5h6C12.227 5 13 4.896 13 6V4c0-1.104-.773-2-1.877-2h-8c-2 0-3.583 2.125-3 5c0 0 1.791 9.375 1.917 9.958C2.373 18.5 4.164 20 6.081 20h6.958c1.105 0-.039-1.896-.039-3v-2c0 1.104-.773 2-1.877 2h-4c-1.104 0-1.833-1.042-2-2S3.539 7.667 3.539 7.667C3.206 5.75 4.018 5 5.123 5m25.812 0h-6C23.831 5 22 4.896 22 6V4c0-1.104 1.831-2 2.935-2h8c2 0 3.584 2.125 3 5c0 0-1.633 9.419-1.771 10c-.354 1.5-2.042 3-4 3h-7.146C21.914 20 22 18.104 22 17v-2c0 1.104 1.831 2 2.935 2h4c1.104 0 1.834-1.042 2-2s1.584-7.333 1.584-7.333C32.851 5.75 32.04 5 30.935 5M20.832 22c0-6.958-2.709 0-2.709 0s-3-6.958-3 0s-3.291 10-3.291 10h12.292c-.001 0-3.292-3.042-3.292-10"
					/><path
						fill="#ffcc4d"
						d="M29.123 6.577c0 6.775-6.77 18.192-11 18.192s-11-11.417-11-18.192c0-5.195 1-6.319 3-6.319c1.374 0 6.025-.027 8-.027l7-.001c2.917-.001 4 .684 4 6.347"
					/><path
						fill="#c1694f"
						d="M27 33c0 1.104.227 2-.877 2h-16C9.018 35 9 34.104 9 33v-1c0-1.104 1.164-2 2.206-2h13.917c1.042 0 1.877.896 1.877 2z"
					/><path
						fill="#c1694f"
						d="M29 34.625c0 .76.165 1.375-1.252 1.375H8.498C7.206 36 7 35.385 7 34.625v-.25C7 33.615 7.738 33 8.498 33h19.25c.759 0 1.252.615 1.252 1.375z"
					/></svg
				>{tournamentName}</a
			>
		{/if}
		<GameInfo
			{whiteUsername}
			{blackUsername}
			{result}
			{createdAt}
			{baseTime}
			{increment}
			{reason}
		/>
		<Chat hei="256px" />
	</div>
	<div class="acontainer w-full xl:w-3/4">
		<div class="abortt">
			{#if result === 0 && (whiteUp ? !moveHistory || moveHistory.length == 0 : moveHistory && moveHistory.length == 1)}
				<AbortTimer
					{lowAbortTime}
					time={abortLength - (baseTime - Math.floor((whiteUp ? wtime : btime) / 1000))}
					tb="t"
				/>
			{/if}
		</div>
		<div class="board flex w-full flex-col justify-center">
			<Chessboard {setGround} {boardConfig} />
		</div>
		<div class="abortb">
			{#if result === 0 && (whiteUp ? moveHistory && moveHistory.length == 1 : !moveHistory || moveHistory.length == 0)}
				<AbortTimer
					{lowAbortTime}
					time={abortLength - (baseTime - Math.floor((whiteUp ? btime : wtime) / 1000))}
					tb="b"
				/>
			{/if}
		</div>
		<div class="clockt h-fit">
			<Clock
				{lowTime}
				time={whiteUp ? wtime : btime}
				active={result !== 0
					? false
					: whiteUp
						? chessLatest.turn() === 'w'
						: chessLatest.turn() === 'b'}
			/>
		</div>
		<div class="rematch flex w-full justify-center">
			{#if result !== 0 && isPlayer && !tournamentID}
				<Rematch {baseTime} {increment} opponentID={whiteUp ? d.game.WhiteID : d.game.BlackID} />
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
				<DrawResign isDisabled={!moveHistory || moveHistory.length < 2} {gameID} />
			{/if}
		</div>
		<div class="back_to_tournament w-full text-center">
			{#if isPlayer && tournamentID && result !== 0}
				<a
					class="flex w-full items-center justify-center gap-[5px] px-3 py-2"
					href={`/tournament/${tournamentID}`}
					><svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24"
						><path fill="currentColor" d="M8 5.14v14l11-7z" /></svg
					>BACK TO TOURNAMENT</a
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
		<div class="clockb h-fit">
			<Clock
				{lowTime}
				time={whiteUp ? btime : wtime}
				active={result !== 0
					? false
					: whiteUp
						? chessLatest.turn() === 'b'
						: chessLatest.turn() === 'w'}
			/>
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
			'abortt abortt'
			'board board'
			'abortb abortb'
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
	}
	.clockb {
		grid-area: clockb;
		justify-self: end;
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
	@media (width>= 768px) {
		.acontainer {
			row-gap: 3px;
			column-gap: 10px;
			place-content: center;
			place-items: start;
			grid-template-areas:
				'board .'
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
		}
		.clockb {
			justify-self: auto;
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
	}
</style>
