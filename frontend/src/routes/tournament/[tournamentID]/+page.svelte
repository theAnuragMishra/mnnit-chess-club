<script lang="ts">
	import { websocketStore } from '$lib/websocket';
	import { onDestroy, onMount } from 'svelte';
	import { page } from '$app/state';
	import { getTimeLeft, getTimeControl } from '$lib/utils';
	import Clock from '$lib/components/Clock.svelte';
	import Chat from '$lib/components/Chat.svelte';

	interface Player {
		ID: number;
		IsActive: boolean;
		Streak: number;
		Score: number;
		Scores: number[];
		Rating: number;
	}

	const { data } = $props();
	const tournamentID = page.params.tournamentID;

	let loading = $state(false);
	let players = $state(data.tournamentData.players ? data.tournamentData.players : []);
	// $effect(() => {
	// 	players.forEach((player: Player) => console.log(player));
	// });
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
		websocketStore.sendMessage({ type: 'join_leave', payload: { tournamentID } });
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
	};

	function updateScore(p: Player) {
		const i = players.findIndex((player: any) => player.ID === p.ID);
		if (i === -1) return;
		Object.assign(players[i], p);
	}

	const handleScoreUpdate = (payload: any) => {
		updateScore(payload.p1);
		updateScore(payload.p2);
	};

	//timers
	let totalTime = data.tournamentData.ongoing
		? getTimeLeft(data.tournamentData.startTime, data.tournamentData.duration)
		: new Date(data.tournamentData.startTime).getTime() - new Date().getTime();
	let animationFrame: number | null;
	let startTime: DOMHighResTimeStamp | null;
	let timeToShow = $state(totalTime);

	$effect(() => {
		if (totalTime > 0) {
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

<h1 class="text-4xl">
	{data.tournamentData.name}
</h1>
<div class="flex w-full justify-between gap-2 p-5">
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
		{#if totalTime >= 0}
			<div>
				{#if !data.tournamentData.ongoing}
					<p>Starting in</p>
				{/if}
				<Clock time={timeToShow} active={true} lowTime={0} />min:sec
			</div>
			<button
				onclick={handleJoinLeave}
				class={`${!joined[0] || (joined[0] && joined[1] === false) ? 'bg-green-500' : 'bg-red-600'} my-2 cursor-pointer rounded-lg px-3 py-1 text-white disabled:cursor-not-allowed`}
				disabled={loading}
				>{data.tournamentData.ongoing
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
	<div class="mr-10 w-1/2">
		{#if joined[1]}
			<h2 class="shiny relative mb-5 overflow-hidden bg-green-600 px-2 py-0.5 text-center text-xl">
				Wait {data.user.username}, pairing players. Get ready!
			</h2>
		{/if}
		<div class="grid grid-cols-[140px_auto_50px] content-start gap-[10px] text-xl">
			{#each sortedPlayers as player, i}
				<span class="flex items-center gap-2">
					<span class="flex w-[30px] items-center justify-center">
						{#if player.IsActive === false}
							<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24"
								><path fill="currentColor" d="M14 19V5h4v14zm-8 0V5h4v14z" /></svg
							>
						{:else}
							{i + 1}
						{/if}</span
					>{player.Username} <i>{Math.floor(player.Rating)}</i>
				</span><span class="flex items-center text-[16px] text-gray-300">
					{#each player.Scores as s}
						{s}
					{/each}
				</span><span class="inline-flex items-center justify-end gap-1 justify-self-end"
					><svg
						class={`h-[24px] w-[24px] ${player.Streak >= 2 ? '' : 'hidden'}`}
						xmlns="http://www.w3.org/2000/svg"
						width="512"
						height="512"
						viewBox="0 0 512 512"
						><path
							fill="#ff8f1f"
							d="M266.91 500.44c-168.738 0-213.822-175.898-193.443-291.147c.887-5.016 7.462-6.461 10.327-2.249c8.872 13.04 16.767 31.875 29.848 30.24c19.661-2.458 33.282-175.946 149.807-224.761c3.698-1.549 7.567 1.39 7.161 5.378c-5.762 56.533 28.181 137.468 88.316 137.468c34.472 0 58.058-27.512 69.844-55.142c3.58-8.393 15.843-7.335 17.896 1.556c21.031 91.082 77.25 398.657-179.756 398.657"
						/><path
							fill="#ffb636"
							d="M207.756 330.827c3.968-3.334 9.992-1.046 10.893 4.058c2.108 11.943 9.04 32.468 31.778 32.468c27.352 0 45.914-75.264 50.782-97.399c.801-3.642 4.35-6.115 8.004-5.372c68.355 13.898 101.59 235.858-48.703 235.858c-109.412 0-84.625-142.839-52.754-169.613M394.537 90.454c2.409-18.842-31.987 32.693-31.987 32.693s26.223 12.386 31.987-32.693M47.963 371.456c.725-8.021-9.594-29.497-11.421-20.994c-4.373 20.344 10.696 29.016 11.421 20.994"
						/><path
							fill="#ffd469"
							d="M323.176 348.596c-2.563-10.69-11.755 14.14-10.6 24.254c1.155 10.113 16.731 1.322 10.6-24.254"
						/></svg
					>
					{player.Score}</span
				>
			{/each}
		</div>
	</div>
	<div class="w-1/4">
		<h1 class="text-xl">Chat Room</h1>
		<Chat hei="400" />
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
