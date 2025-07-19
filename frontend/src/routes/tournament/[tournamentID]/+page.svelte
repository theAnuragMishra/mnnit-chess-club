<script lang="ts">
	import { websocketStore } from '$lib/websocket';
	import { onDestroy, onMount } from 'svelte';
	import { page } from '$app/state';
	import { getTimeLeft, getTimeControl } from '$lib/utils';
	import Clock from '$lib/components/Clock.svelte';
	const { data } = $props();
	const tournamentID = page.params.tournamentID;

	let loading = $state(false);
	let players = $state(data.tournamentData.players ? data.tournamentData.players : []);
	// $effect(() => {
	// 	players.forEach((player) => console.log(player));
	// });
	let sortedPlayers = $derived(players.toSorted((a: any, b: any) => b.Score - a.Score));

	let joined = $derived(isJoined());

	function isJoined() {
		const i = players.findIndex((player: any) => player.Username === data.user.username);
		// console.log(i);
		return i !== -1 && (!data.tournamentData.ongoing || players[i].IsActive);
	}

	const handleJoinLeave = async () => {
		loading = true;
		websocketStore.sendMessage({ type: 'join_leave', payload: { tournamentID } });
	};

	const handleJLResponse = (payload: any) => {
		loading = false;

		if (payload.player) {
			const i = players.findIndex((player: any) => player.Username === payload.player.username);
			if (i === -1) players.push(payload.player);
		} else if (payload.id) {
			const i = players.findIndex((player: any) => player.ID === payload.id);
			if (i !== -1) {
				players.splice(i, 1);
			}
		}
	};

	function updateScore(id: number, score: number) {
		const i = players.findIndex((player: any) => player.ID === id);
		players[i].Score = score;
	}

	const handleScoreUpdate = (payload: any) => {
		updateScore(payload.p1ID, payload.p1Score);
		updateScore(payload.p2ID, payload.p2Score);
	};

	//timers
	let totalTime = data.tournamentData.ongoing
		? getTimeLeft(data.tournamentData.startTime, data.tournamentData.duration)
		: new Date(data.tournamentData.startTime).getTime() - new Date().getTime();
	let animationFrame: number | null;
	let startTime: DOMHighResTimeStamp | null;
	let timeToShow = $state(totalTime);

	$effect(() => {
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

		return () => {
			if (animationFrame !== null) {
				cancelAnimationFrame(animationFrame);
				animationFrame = null;
				startTime = null;
			}
		};
	});

	onMount(() => {
		websocketStore.sendMessage({ type: 'room_change', payload: { room: tournamentID } });
		websocketStore.onMessage('jl_response', handleJLResponse);
		websocketStore.onMessage('update_score', handleScoreUpdate);
	});
	onDestroy(() => {
		websocketStore.sendMessage({ type: 'leave_room', payload: { room: tournamentID } });
		websocketStore.offMessage('jl_response', handleJLResponse);
		websocketStore.offMessage('update_score', handleScoreUpdate);
	});
</script>

<div class="flex w-full justify-between gap-10 p-10">
	<div class="w-1/4">
		<h1 class="text-4xl">
			{data.tournamentData.name}
		</h1>
		<p>By <a href={`/member/${data.tournamentData.creator}`}>{data.tournamentData.creator}</a></p>
		<p>
			Starts {new Date(data.tournamentData.startTime).toLocaleString('en-IN', {
				year: 'numeric',
				month: 'short',
				day: '2-digit',
				hour: '2-digit',
				minute: '2-digit',
				hour12: false
			})}
		</p>
		<p>Duration: {data.tournamentData.duration / 3600} hours</p>
		<p>
			{data.tournamentData.baseTime / 60}+{data.tournamentData.increment} Rated {getTimeControl(
				data.tournamentData.baseTime,
				data.tournamentData.increment
			)}
		</p>
		<div>
			{#if !data.tournamentData.ongoing}
				<p>Starting in</p>
			{/if}
			<Clock time={timeToShow} active={true} lowTime={0} />min:sec
		</div>
		<button
			onclick={handleJoinLeave}
			class={`${joined ? 'bg-red-600' : 'bg-green-500'} my-2 cursor-pointer rounded-lg px-3 py-1 text-white disabled:cursor-not-allowed`}
			disabled={loading}
			>{data.tournamentData.ongoing
				? joined
					? 'Pause'
					: 'Resume'
				: joined
					? 'Leave'
					: 'Join'}</button
		>
	</div>
	<div class="grid w-2/3 grid-cols-[50px_1fr_2fr] justify-start gap-[10px] text-xl">
		{#each sortedPlayers as player, i}
			<span>{i + 1}.</span><span>{player.Username} <i>{Math.floor(player.Rating)}</i></span><span
				>{player.Score}</span
			>
		{/each}
	</div>
</div>
