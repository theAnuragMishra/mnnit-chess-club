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
	let { data } = $props();
	// console.log(data);
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
	let chess = $derived(moveHistory ? new Chess(moveHistory[activeIndex].MoveFen) : new Chess());

	const whiteUp = data.gameData.game.WhiteUsername !== data.user.username;
	const setActiveIndex = (index: number) => {
		if (index > moveHistory.length - 1 || index < 0) return;
		activeIndex = index;
	};
	const gameID = Number(page.params.gameID);
	const sendMessage = websocketStore.sendMessage;

	const handleMessage = (e: MessageEvent) => {
		// console.log('event received');
		const data = JSON.parse(e.data);
		if (data.type === 'timeup') {
			result = data.payload.Result;
			reason = data.payload.Reason;
		}
		if (data.type === 'Move_Response') {
			// console.log(data);
			if (moveHistory) moveHistory = [...moveHistory, data.payload.move];
			else moveHistory = [data.payload.move];
			timeBlack = data.payload.timeBlack;
			timeWhite = data.payload.timeWhite;
			activeIndex = activeIndex + 1;

			if (data.payload.Result !== '') {
				result = data.payload.Result;
				reason = data.payload.Reason;
			}
		}
	};

	onMount(() => {
		websocketStore.socket?.addEventListener('message', handleMessage);
	});
	onDestroy(() => websocketStore.socket?.removeEventListener('message', handleMessage));
</script>

<div class="px-5 text-2xl">
	{#if result !== 'ongoing' && result !== ''}
		<ResultModal {result} {reason} />
	{/if}
	<div class="flex w-full items-center justify-around">
		{#if result !== '1-0' && result !== '0-1'}
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
				username={data.user.username}
				{gameID}
				{chess}
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
						sendMessage({
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
							? chess.turn() === 'w'
							: chess.turn() === 'b'}
				/>
			</p>

			<HistoryTable {moveHistory} {setActiveIndex} {activeIndex} />

			<p class="mb-1 flex w-full items-center justify-between">
				{data.user.username}
				<Clock
					onTimeUp={() => {
						sendMessage({
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
							? chess.turn() === 'b'
							: chess.turn() === 'w'}
				/>
			</p>
		</div>
	</div>
</div>
