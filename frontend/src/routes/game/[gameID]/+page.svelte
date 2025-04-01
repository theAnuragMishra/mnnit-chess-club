<script lang="ts">
	import { Chess } from 'chess.js';
	import { page } from '$app/state';
	import ResultModal from '$lib/components/ResultModal.svelte';
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
	let { data } = $props();
	// console.log(data);
	let ground: Api | null = $state(null);
	let whiteUsername = data.gameData.game.WhiteUsername;
	let blackUsername = data.gameData.game.BlackUsername;
	let whiteID = data.gameData.game.WhiteID;
	let blackID = data.gameData.game.BlackID;
	let timeBlack = $state(data.gameData.timeBlack);
	let timeWhite = $state(data.gameData.timeWhite);
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

	const whiteUp = data.gameData.game.WhiteUsername !== data.user.username;
	const setActiveIndex = (index: number) => {
		if (index > moveHistory.length - 1 || index < 0) return;
		activeIndex = index;
	};
	const gameID = Number(page.params.gameID);

	const setGround = (g: Api) => {
		ground = g;
	};

	const handleTimeUp = (payload: any) => {
		if (payload.gameID != gameID) return;
		result = payload.Result;
		reason = payload.Reason;
	};

	const handleResignation = (payload: any) => {
		if (payload.gameID != gameID) return;
		result = payload.Result;
		reason = payload.Reason;
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
		}
		ground?.set({
			fen: payload.move.MoveFen,
			turnColor: chessLatest.turn() === 'w' ? 'white' : 'black',
			movable: { dests: getValidMoves(chessLatest) }
		});
		// console.log('finished movehandler');
		ground?.playPremove();
		// console.log(x);
	};

	onMount(() => {
		websocketStore.onMessage('timeup', handleTimeUp);
		websocketStore.onMessage('Move_Response', handleMoveResponse);
		websocketStore.onMessage('resignation', handleResignation);
	});
	onDestroy(() => {
		websocketStore.offMessage('timeup', handleTimeUp);
		websocketStore.offMessage('Move_Response', handleMoveResponse);
		websocketStore.offMessage('resignation', handleResignation);
	});
</script>

<div class="px-5 text-2xl">
	{#if result !== 'ongoing' && result !== ''}
		<ResultModal {result} {reason} />
	{/if}
	<div class="flex w-full items-center justify-around">
		{#if result === '' || result === 'ongoing'}
			<div class="flex w-1/4 flex-col gap-10">
				<Chat
					username={data.user.username}
					userID={data.user.userID}
					{blackID}
					{whiteID}
					{gameID}
					{whiteUsername}
					{blackUsername}
				/>
				<DrawResign
					{gameID}
					userID={data.user.userID}
					setResultReason={(res: string, rea: string) => {
						result = res;
						reason = rea;
					}}
				/>
			</div>
		{:else}
			<GameInfo />
		{/if}
		<div class="flex items-center justify-center">
			<Chessboard
				{setGround}
				username={data.user.username}
				{gameID}
				chess={chessForView}
				white={whiteUsername}
				lastMove={moveHistory ? [moveHistory[activeIndex].Orig, moveHistory[activeIndex].Dest] : []}
				viewOnly={(result != 'ongoing' && result != '') ||
					(moveHistory && activeIndex !== moveHistory.length - 1)}
			/>
		</div>
		<div class="flex h-full w-1/4 flex-col gap-2">
			<p class="mb-1 flex w-full items-center justify-between">
				{whiteUp ? whiteUsername : blackUsername}
				<Clock
					onTimeUp={() => {
						websocketStore.sendMessage({
							type: 'timeup',
							payload: {
								gameID: gameID
							}
						});
					}}
					initialTime={whiteUp ? timeWhite : timeBlack}
					active={result !== 'ongoing' && result !== ''
						? false
						: whiteUp
							? chessLatest.turn() === 'w'
							: chessLatest.turn() === 'b'}
				/>
			</p>

			<HistoryTable
				{moveHistory}
				{setActiveIndex}
				{activeIndex}
				highlightLastArrow={activeIndex !== moveHistory.length - 1 &&
					(whiteUp ? chessLatest.turn() === 'b' : chessLatest.turn() === 'w') &&
					(result === 'ongoing' || result === '')}
			/>

			<p class="mb-1 flex w-full items-center justify-between">
				{data.user.username}
				<Clock
					onTimeUp={() => {
						websocketStore.sendMessage({
							type: 'timeup',
							payload: {
								gameID: gameID
							}
						});
					}}
					initialTime={whiteUp ? timeBlack : timeWhite}
					active={result !== 'ongoing' && result !== ''
						? false
						: whiteUp
							? chessLatest.turn() === 'b'
							: chessLatest.turn() === 'w'}
				/>
			</p>
		</div>
	</div>
</div>
