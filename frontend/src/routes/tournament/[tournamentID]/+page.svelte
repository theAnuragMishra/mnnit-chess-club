<script lang="ts">
	import { websocketStore } from '$lib/websocket';
	import { page } from '$app/state';
	import { getTimeLeft, getTimeControl } from '$lib/utils';
	import Clock from '$lib/components/Clock.svelte';
	import Chat from '$lib/components/Chat.svelte';
	import TopThree from '$lib/components/TopThree.svelte';
	import trophyImg from '$lib/assets/icons/trophy.svg';
	import pauseImg from '$lib/assets/icons/pause.svg';
	import fireImg from '$lib/assets/icons/fire.svg';
	import berserkImg from '$lib/assets/icons/kill.svg';
	interface Player {
		ID: number;
		IsActive: boolean;
		Streak: number;
		Score: number;
		Scores: number[];
		Rating: number;
	}

	const { data } = $props();

	let loading = $state(false);
	let players = $state(data.tournamentData.players);
	// $effect(() => {
	// 	players.forEach((player: Player) => console.log(player));
	// });
	$effect(() => {
		players = data.tournamentData.players;
	});
	let sortedPlayers = $derived(players.toSorted((a: any, b: any) => b.Score - a.Score));

	let joined = $derived(isJoinedActive());

	//$effect(() => console.log(joined));

	function isJoinedActive() {
		const i = players.findIndex((player: any) => player.Username === data.user.username);
		// console.log(i);
		return [i !== -1, i !== -1 && players[i].IsActive];
	}

	const handleJoinLeave = async () => {
		loading = true;
		websocketStore.sendMessage({
			type: 'join_leave',
			payload: { tournamentID: page.params.tournamentID }
		});
	};

	const handleJLResponse = (payload: any) => {
		loading = false;

		if (payload.player) {
			//console.log(payload.player);
			const i = players.findIndex((player: any) => player.ID === payload.player.ID);
			if (i === -1) players.push(payload.player);
			else players[i].IsActive = payload.player.IsActive;
		} else if (payload.id) {
			const i = players.findIndex((player: any) => player.ID === payload.id);
			if (i !== -1) {
				players.splice(i, 1);
			}
		}
		//console.log(players);
	};

	function updateScore(p: Player) {
		const i = players.findIndex((player: any) => player.ID === p.ID);
		if (i === -1) return;
		Object.assign(players[i], p);
		// console.log(players);
	}

	const handleScoreUpdate = (payload: any) => {
		// console.log(payload);
		updateScore(payload.p1);
		updateScore(payload.p2);
	};

	//timers
	let totalTime = $derived(
		data.tournamentData.status === 1
			? getTimeLeft(data.tournamentData.startTime, data.tournamentData.duration)
			: new Date(data.tournamentData.startTime).getTime() - new Date().getTime()
	);
	let animationFrame: number | null;
	let startTime: DOMHighResTimeStamp | null;
	let timeToShow = $derived(totalTime);

	$effect(() => {
		if (data.tournamentData.status !== 2) {
			//console.log('effect working');
			startTime = performance.now();
			const tick = (currentTime: number) => {
				if (!startTime) return;
				const elapsed = currentTime - startTime;
				const newTime = totalTime - elapsed;

				if (newTime <= 0) {
					timeToShow = 0;
					return;
				}
				timeToShow = newTime;

				animationFrame = requestAnimationFrame(tick);
			};

			animationFrame = requestAnimationFrame(tick);
		}

		return () => {
			if (animationFrame !== null) {
				cancelAnimationFrame(animationFrame);
				animationFrame = null;
				startTime = null;
			}
		};
	});

	$effect(() => {
		websocketStore.sendMessage({
			type: 'room_change',
			payload: { room: page.params.tournamentID }
		});
		websocketStore.onMessage('jl_response', handleJLResponse);
		websocketStore.onMessage('update_score', handleScoreUpdate);
		return () => {
			websocketStore.sendMessage({ type: 'leave_room' });
			websocketStore.offMessage('jl_response', handleJLResponse);
			websocketStore.offMessage('update_score', handleScoreUpdate);
		};
	});
</script>

<svelte:head>
	<title>{data.tournamentData.name} - MCC Arena</title>
</svelte:head>
<h1 class="flex items-center gap-[5px] text-3xl">
	<img class="h-[30px]" src={trophyImg} alt="trophy" />{data.tournamentData.name}
</h1>
<div class="flex w-full flex-col-reverse gap-10 p-2 md:flex-row md:gap-2 md:p-5">
	<div class="flex w-full flex-col gap-2 md:w-1/4">
		<div>
			{#if data.tournamentData.status !== 1}
				<p>
					{#if data.tournamentData.status === 0}
						Starts{/if}
					{new Date(data.tournamentData.startTime).toLocaleString('en-IN', {
						year: 'numeric',
						month: 'short',
						day: '2-digit',
						hour: '2-digit',
						minute: '2-digit',
						hour12: false
					})}
				</p>
			{/if}
			<p>
				{data.tournamentData.baseTime / 60}+{data.tournamentData.increment} &bull; {getTimeControl(
					data.tournamentData.baseTime,
					data.tournamentData.increment
				)} &bull; {data.tournamentData.duration / 60}m
			</p>
			{#if data.tournamentData.berserkAllowed}
				<p class="flex items-center gap-[5px]">
					<img class="h-[16px]" src={berserkImg} alt="berserk allowed" />Berserk Allowed
				</p>
			{/if}
			<p>By <a href={`/member/${data.tournamentData.creator}`}>{data.tournamentData.creator}</a></p>
			{#if data.tournamentData.status !== 2}
				<div>
					{#if data.tournamentData.status === 0}
						<p>Starting in</p>
					{/if}
					<Clock time={timeToShow} active={true} lowTime={0} />min:sec
				</div>
				<button
					onclick={handleJoinLeave}
					class={`${!joined[0] || (joined[0] && joined[1] === false) ? 'bg-green-500' : 'bg-red-600'} my-2 cursor-pointer rounded-lg px-3 py-1 text-white disabled:cursor-not-allowed`}
					disabled={loading}
					>{data.tournamentData.status === 1
						? joined[0]
							? joined[1]
								? 'Pause'
								: 'Resume'
							: 'Join'
						: joined[0]
							? 'Leave'
							: 'Join'}</button
				>
			{:else}
				<p>Tournament Over</p>
			{/if}
		</div>
		<div>
			<h1 class="mb-2 text-xl">Chat Room</h1>
			<Chat hei="250px" />
		</div>
	</div>
	<div class="mr-10 w-full md:w-1/2">
		{#if data.tournamentData.status === 2}
			<TopThree
				u1={sortedPlayers[0]?.Username}
				u2={sortedPlayers[1]?.Username}
				u3={sortedPlayers[2]?.Username}
			/>
		{/if}
		{#if joined[1]}
			<h2
				class="shiny relative mb-5 overflow-hidden bg-green-600 py-0.5 text-center text-[16px] md:px-2 md:text-xl"
			>
				Standby {data.user.username}, pairing players. Get ready!
			</h2>
		{/if}
		<div class="grid grid-cols-[auto_auto] content-start gap-[10px] text-[16px] md:text-xl">
			{#each sortedPlayers as player, i}
				<div>
					<div class="inline-flex items-center gap-[5px]">
						<span class="inline-flex h-[24px] w-[24px] items-center justify-center">
							{#if player.IsActive === false}
								<img src={pauseImg} alt="played paused" />
							{:else}
								{i + 1}
							{/if}</span
						><a href={`/member/${player.Username}`}>{player.Username}</a>
						<i>{Math.floor(player.Rating)}</i>
					</div>
					<span class="text-[16px] text-gray-300">
						{#each player.Scores as s}
							{s}
						{/each}
					</span>
				</div>
				<div class="inline-flex justify-end gap-1 justify-self-end">
					<img
						src={fireImg}
						alt="player on streak"
						class={`h-[24px] w-[24px] ${player.Streak >= 2 ? '' : 'hidden'}`}
					/>
					{player.Score}
				</div>
			{/each}
		</div>
	</div>
</div>

<style>
	.shiny::before {
		content: '';
		position: absolute;
		top: 0;
		left: -100%;
		height: 100%;
		width: 100%;
		background: linear-gradient(
			120deg,
			transparent 0%,
			rgba(255, 255, 255, 0.8) 30%,
			transparent 100%
		);
		animation: shine 3s linear infinite;
	}
	@keyframes shine {
		0% {
			left: -100%;
		}
		100% {
			left: 100%;
		}
	}
</style>
