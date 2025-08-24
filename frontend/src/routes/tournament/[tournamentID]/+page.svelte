<script lang="ts">
	import { websocketStore } from '$lib/websocket';
	import { page } from '$app/state';
	import { getTimeLeft, getTimeControl } from '$lib/utils';
	import Clock from '$lib/components/Clock.svelte';
	import Chat from '$lib/components/Chat.svelte';
	import TopThree from '$lib/components/TopThree.svelte';

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
	let players = $state(data.tournamentData.players ?? []);
	// $effect(() => {
	// 	players.forEach((player: Player) => console.log(player));
	// });
	$effect(() => {
		players = data.tournamentData.players ?? [];
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
	<svg
		class="h-[50px]"
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
	>{data.tournamentData.name}
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
					><a href={`/member/${player.Username}`}>{player.Username}</a>
					<i>{Math.floor(player.Rating)}</i>
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
