<script lang="ts">
	import { Chess } from 'chess.js';
	import { page } from '$app/state';
	import Chat from '$lib/components/Chat.svelte';
	import DrawResign from '$lib/components/DrawResign.svelte';
	import GameInfo from '$lib/components/GameInfo.svelte';
	import Chessboard from '$lib/components/Chessboard.svelte';
	import Clock from '$lib/components/Clock.svelte';
	import { websocketStore } from '$lib/websocket.js';
	import HistoryTable from '$lib/components/HistoryTable.svelte';
	import { onDestroy, onMount } from 'svelte';
	import type { Api } from 'chessground/api';
	import { getValidMoves } from '$lib/utils.js';
	import AbortTimer from '$lib/components/AbortTimer.svelte';
	let { data } = $props();
	//console.log(data);

	//audios
	let moveAudio: HTMLAudioElement;
	let captureAudio: HTMLAudioElement;
	let notifyAudio: HTMLAudioElement;
	let lowTimeAudio: HTMLAudioElement;

	let ground: Api | null = $state(null);
	const baseTime = data.gameData.game.BaseTime;
	const increment = data.gameData.game.Increment;
	const createdAt = new Date(data.gameData.game.CreatedAt);
	let whiteUsername = data.gameData.game.WhiteUsername;
	let blackUsername = data.gameData.game.BlackUsername;
	let timeBlack = $state(data.gameData.timeBlack);
	let timeWhite = $state(data.gameData.timeWhite);
	let ratingWhite = data.gameData.game.RatingW;
	let ratingBlack = data.gameData.game.RatingB;
	let changeWhite = $state(data.gameData.game.ChangeW ?? 0);
	let changeBlack = $state(data.gameData.game.ChangeB ?? 0);
	let result = $state(data.gameData.game.Result);
	let reason = $state(
		data.gameData.game.Result != 'ongoing' && data.gameData.game.Result != ''
			? data.gameData.game.ResultReason
			: ''
	);
	let moveHistory = $state(data.gameData.moves);
	let activeIndex = $state(data.gameData.moves ? data.gameData.moves.length - 1 : -1);
	let chessLatest = $derived(
		moveHistory ? new Chess(moveHistory[moveHistory.length - 1].MoveFen) : new Chess()
	);
	let chessForView = $derived(
		moveHistory ? new Chess(moveHistory[activeIndex].MoveFen) : new Chess()
	);
	const isPlayer =
		data.user.username === data.gameData.game.WhiteUsername ||
		data.user.username === data.gameData.game.BlackUsername;
	const whiteUp = data.gameData.game.WhiteUsername !== data.user.username;
	const setActiveIndex = (index: number) => {
		if (index > moveHistory.length - 1 || index < 0) return;
		activeIndex = index;
	};
	const gameID = page.params.gameID;

	const setGround = (g: Api) => {
		ground = g;
	};

	const handleTimeUp = (payload: any) => {
		if (payload.gameID != gameID) return;
		result = payload.Result;
		reason = payload.Reason;
		changeBlack = payload.changeB;
		changeWhite = payload.changeW;
	};

	const handleResignation = (payload: any) => {
		if (payload.gameID != gameID) return;
		result = payload.Result;
		reason = payload.Reason;
		changeBlack = payload.changeB;
		changeWhite = payload.changeW;
	};

	const handleMoveResponse = (payload: any) => {
		if (payload.gameID != gameID) return;
		if (moveHistory) moveHistory = [...moveHistory, payload.move];
		else moveHistory = [payload.move];
		timeBlack = payload.timeBlack;
		timeWhite = payload.timeWhite;
		if (activeIndex === moveHistory.length - 2) activeIndex = moveHistory.length - 1;

		if (payload.Result !== '') {
			result = payload.Result;
			reason = payload.message;
			changeBlack = payload.changeB;
			changeWhite = payload.changeW;
		}
		ground?.set({
			fen: payload.move.MoveFen,
			turnColor: chessLatest.turn() === 'w' ? 'white' : 'black',
			movable: { dests: getValidMoves(chessLatest) }
		});
		// console.log('finished movehandler');
		if (payload.move.MoveNotation[1] == 'x') captureAudio?.play();
		else moveAudio?.play();
		ground?.playPremove();
		// console.log(x);
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
	let btime = $derived(timeBlack);
	let wtime = $derived(timeWhite);
	let animationFrame: number | null;
	let startTime: DOMHighResTimeStamp | null;
	let lowTimePlayed = $state(false);

	$effect(() => {
		if (result === 'ongoing' || result === '') {
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
		websocketStore.onMessage('timeup', handleTimeUp);
		websocketStore.onMessage('game_abort', handleTimeUp);
		websocketStore.onMessage('Move_Response', handleMoveResponse);
		websocketStore.onMessage('resignation', handleResignation);
		moveAudio = new Audio('/Move.mp3');
		captureAudio = new Audio('/Capture.mp3');
		notifyAudio = new Audio('/GenericNotify.mp3');
		lowTimeAudio = new Audio('/LowTime.mp3');
	});
	onDestroy(() => {
		websocketStore.offMessage('timeup', handleTimeUp);
		websocketStore.offMessage('game_abort', handleTimeUp);
		websocketStore.offMessage('Move_Response', handleMoveResponse);
		websocketStore.offMessage('resignation', handleResignation);
	});
</script>

<div class="flex flex-col-reverse items-center justify-center gap-5 xl:flex-row">
	<div class="flex w-[95%] flex-col gap-5 md:flex-row xl:w-[300px] xl:flex-col">
		<GameInfo
			{whiteUsername}
			{blackUsername}
			{result}
			{createdAt}
			{baseTime}
			{increment}
			{reason}
		/>
		<Chat {gameID} />
	</div>
	<div class="acontainer xl:w-3/4">
		<div class="abortt">
			{#if (result === 'ongoing' || result === '') && (whiteUp ? !moveHistory || moveHistory.length == 0 : moveHistory && moveHistory.length == 1)}
				<AbortTimer
					{lowAbortTime}
					time={abortLength - (baseTime - Math.floor((whiteUp ? wtime : btime) / 1000))}
					tb="t"
				/>
			{/if}
		</div>
		<div class="board flex flex-col justify-center">
			<Chessboard
				{setGround}
				orientation={whiteUp ? 'black' : 'white'}
				{gameID}
				chess={chessForView}
				lastMove={moveHistory ? [moveHistory[activeIndex].Orig, moveHistory[activeIndex].Dest] : []}
				viewOnly={!isPlayer ||
					(result != 'ongoing' && result != '') ||
					(moveHistory && activeIndex !== moveHistory.length - 1)}
			/>
		</div>
		<div class="abortb">
			{#if (result === 'ongoing' || result === '') && (whiteUp ? moveHistory && moveHistory.length == 1 : !moveHistory || moveHistory.length == 0)}
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
				active={result !== 'ongoing' && result !== ''
					? false
					: whiteUp
						? chessLatest.turn() === 'w'
						: chessLatest.turn() === 'b'}
			/>
		</div>

		<div class="namet flex h-fit justify-between md:w-[300px]">
			<a href={`/member/${whiteUp ? whiteUsername : blackUsername}`}
				>{whiteUp ? whiteUsername : blackUsername}&nbsp;&nbsp;</a
			><span
				>{whiteUp ? ratingWhite : ratingBlack}&nbsp;&nbsp;<span
					class={`${whiteUp ? (changeWhite > 0 ? 'text-green-500' : changeWhite == 0 ? 'text-gray-500' : 'text-red-500') : changeBlack > 0 ? 'text-green-500' : changeBlack == 0 ? 'text-gray-500' : 'text-red-500'}`}
					>{result != '' && result != 'ongoing'
						? whiteUp
							? `${changeWhite > 0 ? '+' : ''}${changeWhite}`
							: `${changeBlack > 0 ? '+' : ''}${changeBlack}`
						: ''}</span
				></span
			>
		</div>
		<div class="draw-resign h-fit w-full">
			{#if isPlayer && (result == '' || result == 'ongoing')}
				<DrawResign
					isDisabled={!moveHistory || moveHistory.length < 2}
					{gameID}
					userID={data.user.userID}
					setResultReason={(res: string, rea: string, cw: number, cb: number) => {
						result = res;
						reason = rea;
						changeBlack = cb;
						changeWhite = cw;
					}}
				/>
			{/if}
		</div>
		<div class="history w-full">
			<HistoryTable
				{moveHistory}
				{setActiveIndex}
				{activeIndex}
				highlightLastArrow={activeIndex !== moveHistory.length - 1 &&
					(whiteUp ? chessLatest.turn() === 'b' : chessLatest.turn() === 'w') &&
					(result === 'ongoing' || result === '')}
			/>
		</div>

		<div class="nameb flex h-fit justify-between md:w-[300px]">
			<a href={`/member/${whiteUp ? blackUsername : whiteUsername}`}
				>{whiteUp ? blackUsername : whiteUsername}&nbsp;&nbsp;</a
			>
			<span
				>{whiteUp ? ratingBlack : ratingWhite}&nbsp;&nbsp;<span
					class={`${!whiteUp ? (changeWhite > 0 ? 'text-green-500' : changeWhite == 0 ? 'text-gray-500' : 'text-red-500') : changeBlack > 0 ? 'text-green-500' : changeBlack == 0 ? 'text-gray-500' : 'text-red-500'}`}
					>{result != '' && result != 'ongoing'
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
				active={result !== 'ongoing' && result !== ''
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
			'draw-resign draw-resign'
			'history history';
	}
	.namet {
		grid-area: namet;
		justify-self: start;
	}
	.nameb {
		grid-area: nameb;
		justify-self: start;
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
	.abortt {
		grid-area: abortt;
	}
	.abortb {
		grid-area: abortb;
	}
	@media (width>= 768px) {
		.acontainer {
			row-gap: 3px;
			column-gap: 10px;
			place-content: start;
			place-items: start;
			grid-template-areas:
				'board .'
				'board clockt'
				'board namet'
				'board abortt'
				'board history'
				'board draw-resign'
				'board abortb'
				'board nameb'
				'board clockb'
				'board .';
		}
		.nameb {
			justify-self: auto;
		}
		.namet {
			justify-self: auto;
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
		.board {
			justify-self: auto;
		}
	}
</style>
