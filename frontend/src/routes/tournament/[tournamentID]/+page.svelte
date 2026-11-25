<script lang="ts">
	import { page } from '$app/state';
	import { getTimeLeft, getTimeControl } from '$lib/utils';
	import Clock from '$lib/components/Clock.svelte';
	import Chat from '$lib/components/Chat.svelte';
	import TopThree from '$lib/components/TopThree.svelte';
	import fireImg from '$lib/assets/icons/fire.svg';
	import berserkImg from '$lib/assets/icons/kill.svg';
	import { user } from '$lib/user.svelte.js';
	import { websocketStore } from '$lib/websocket.svelte.js';
	interface Player {
		ID: number;
		Username: string;
		IsActive: boolean;
		Streak: number;
		Score: number;
		Scores: number[];
		Rating: number;
	}

	const { data } = $props();

	let loading = $state(false);
	let players = $state<Player[]>(data.tournamentData.players ?? []);
	// $effect(() => {
	// 	players.forEach((player: Player) => console.log(player));
	// });
	$effect(() => {
		players = data.tournamentData.players ?? [];
	});
	let sortedPlayers = $derived(players.toSorted((a, b) => b.Score - a.Score));

	let joined = $derived(isJoinedActive());
	let status = $derived(data.tournamentData.status);
	let joinAction = $derived(
		!joined[0] ? 'Join' : status === 0 ? 'Leave' : joined[1] ? 'Pause' : 'Resume'
	);

	//$effect(() => console.log(joined));

	function isJoinedActive() {
		const i = players.findIndex((player) => player.Username === user.username);
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
			const i = players.findIndex((player) => player.ID === payload.player.ID);
			if (i === -1) players.push(payload.player);
			else players[i].IsActive = payload.player.IsActive;
		} else if (payload.id) {
			const i = players.findIndex((player) => player.ID === payload.id);
			if (i !== -1) {
				players.splice(i, 1);
			}
		}
		//console.log(players);
	};

	function updateScore(p: Player) {
		const i = players.findIndex((player) => player.ID === p.ID);
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
	<title>{data.tournamentData.name} — MNNIT Chess Club</title>
</svelte:head>

<div class="page-shell tournament-page">
	<a href="/tournaments" class="text-link">All tournaments</a>
	<h1 class="page-heading">{data.tournamentData.name}</h1>

	<div class="tournament-grid">
		<aside class="tournament-sidebar" aria-label="Arena details">
			<section class="arena-details panel" aria-label="Tournament information">
				{#if status !== 0}
					<p class="arena-status" class:live={status === 1}>
						{status === 1 ? 'Live' : 'Finished'}
					</p>
				{/if}
				<p class="arena-format">
					<strong>{data.tournamentData.baseTime / 60}+{data.tournamentData.increment}</strong>
					{getTimeControl(data.tournamentData.baseTime, data.tournamentData.increment)} · {data
						.tournamentData.duration / 60} min
				</p>
				{#if status !== 1}
					<time datetime={data.tournamentData.startTime}>
						{new Date(data.tournamentData.startTime).toLocaleString('en-IN', {
							year: 'numeric',
							month: 'short',
							day: '2-digit',
							hour: '2-digit',
							minute: '2-digit'
						})}
					</time>
				{/if}
				{#if data.tournamentData.berserkAllowed}
					<p class="arena-meta">
						<img src={berserkImg} alt="" width="16" height="16" />Berserk allowed
					</p>
				{/if}
				{#if data.tournamentData.creator}
					<p>
						By <a class="text-link" href={`/member/${data.tournamentData.creator}`}
							>{data.tournamentData.creator}</a
						>
					</p>
				{/if}
				{#if status !== 2}
					<div class="arena-countdown">
						<span>{status === 0 ? 'Starts in' : 'Time remaining'}</span>
						<Clock time={Math.max(0, timeToShow)} active={true} lowTime={0} compact />
					</div>
					<button
						class="button"
						class:button-primary={joinAction === 'Join' || joinAction === 'Resume'}
						class:button-secondary={joinAction === 'Pause' || joinAction === 'Leave'}
						onclick={handleJoinLeave}
						disabled={loading}>{loading ? 'Updating…' : joinAction}</button
					>
				{/if}
			</section>
		</aside>

		<section class="standings panel" aria-labelledby="standings-title">
			<div class="standings-heading">
				<h2 id="standings-title">Standings</h2>
				<span>{sortedPlayers.length} {sortedPlayers.length === 1 ? 'player' : 'players'}</span>
			</div>
			{#if status === 1 && joined[1]}
				<p class="arena-queue" role="status">Waiting for an opponent…</p>
			{/if}
			{#if status === 2 && sortedPlayers.length}
				<TopThree
					players={sortedPlayers
						.slice(0, 3)
						.map((player) => ({ username: player.Username, score: player.Score }))}
					scoreLabel="Points"
				/>
			{/if}
			{#if sortedPlayers.length === 0}
				<p class="empty-standings">No players yet</p>
			{:else}
				<table>
					<thead>
						<tr><th scope="col">Rank</th><th scope="col">Player</th><th scope="col">Points</th></tr>
					</thead>
					<tbody>
						{#each sortedPlayers as player, i (player.ID)}
							<tr>
								<td>{i + 1}</td>
								<td>
									<a href={`/member/${player.Username}`}>{player.Username}</a>
									<span class="player-meta">
										{Math.floor(player.Rating)}
										{#if status === 1 && player.IsActive === false}
											· Paused{/if}
									</span>
									{#if player.Scores?.length}
										<span class="score-strip" aria-label="Game points">
											{#each player.Scores as score}<span>{score}</span>{/each}
										</span>
									{/if}
								</td>
								<td>
									<span class="points">
										{#if player.Streak >= 2}<img
												src={fireImg}
												alt="Winning streak"
												width="16"
												height="16"
											/>{/if}
										<strong>{player.Score}</strong>
									</span>
								</td>
							</tr>
						{/each}
					</tbody>
				</table>
			{/if}
		</section>

		<details class="arena-chat panel">
			<summary>Chat</summary>
			<Chat hei="clamp(6rem, 24dvh, 14rem)" />
		</details>
	</div>
</div>

<style>
	.tournament-page {
		padding-top: clamp(1.25rem, 4dvh, 2.5rem);
	}
	.page-heading {
		overflow-wrap: anywhere;
	}
	.tournament-grid {
		display: grid;
		grid-template-columns: minmax(13rem, 17rem) minmax(0, 1fr);
		grid-template-areas: 'details standings' 'chat standings';
		grid-template-rows: auto 1fr;
		gap: 1rem;
		align-items: start;
		margin-top: 1.25rem;
	}
	.tournament-sidebar {
		grid-area: details;
	}
	.arena-chat {
		grid-area: chat;
	}
	.arena-chat summary {
		padding: 0.75rem 1rem;
		font-size: 0.875rem;
	}
	.arena-details {
		display: grid;
		gap: 0.5rem;
		padding: 1rem;
		color: var(--muted);
		font-size: 0.875rem;
		overflow-wrap: anywhere;
	}
	.arena-status {
		width: fit-content;
		padding: 0.125rem 0.5rem;
		border: 1px solid var(--line);
		border-radius: 5px;
	}
	.arena-status.live,
	.arena-queue {
		border-color: var(--accent-line);
		background: var(--accent-surface);
		color: var(--accent);
	}
	.arena-format,
	.arena-meta {
		display: flex;
		align-items: center;
		gap: 0.5rem;
	}
	.arena-format strong {
		color: #f0f2eb;
		font-size: 1.125rem;
	}
	.arena-countdown {
		display: flex;
		justify-content: space-between;
		align-items: center;
		gap: 0.5rem;
		margin-block: 0.375rem;
		padding-top: 0.625rem;
		border-top: 1px solid var(--line);
	}
	.standings {
		grid-area: standings;
		min-width: 0;
		padding: 1rem;
	}
	.standings-heading {
		display: flex;
		align-items: baseline;
		justify-content: space-between;
		gap: 0.5rem;
		margin-bottom: 0.875rem;
	}
	.standings-heading h2 {
		font-size: 1.125rem;
		font-weight: 650;
	}
	.standings-heading span,
	.empty-standings,
	.player-meta {
		color: var(--muted);
		font-size: 0.8125rem;
	}
	.empty-standings {
		padding: 1.5rem 0;
		text-align: center;
	}
	.arena-queue {
		margin-bottom: 0.875rem;
		padding: 0.625rem;
		border-width: 1px;
		border-radius: 8px;
		font-size: 0.875rem;
	}
	table {
		width: 100%;
		text-align: left;
		font-size: 0.875rem;
	}
	th,
	td {
		padding: 0.5rem;
		border-bottom: 1px solid var(--line);
		vertical-align: top;
	}
	th {
		color: var(--muted);
		font-size: 0.75rem;
		font-weight: 500;
	}
	th:first-child,
	th:last-child {
		width: 1%;
	}
	th:last-child,
	td:last-child {
		text-align: right;
	}
	td a {
		font-weight: 650;
		overflow-wrap: anywhere;
	}
	td a:hover {
		color: var(--accent);
	}
	.player-meta {
		display: block;
	}
	.points,
	.score-strip {
		display: inline-flex;
		align-items: center;
		gap: 0.25rem;
	}
	.score-strip {
		flex-wrap: wrap;
		margin-top: 0.25rem;
	}
	.score-strip > span {
		padding-inline: 0.25rem;
		border-radius: 3px;
		background: var(--surface-raised);
		color: var(--muted);
		font-size: 0.75rem;
	}
	.points {
		font-variant-numeric: tabular-nums;
	}
	@media (max-width: 760px) {
		.tournament-grid {
			grid-template-columns: minmax(0, 1fr);
			grid-template-areas: 'details' 'standings' 'chat';
		}
	}
</style>
