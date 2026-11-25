<script lang="ts">
	import Icon from '$lib/components/Icon.svelte';
	import TopThree from '$lib/components/TopThree.svelte';
	let { data } = $props();
	let players = $derived(data.leaderboard ?? []);
</script>

<svelte:head><title>Leaderboard — MNNIT Chess Club</title></svelte:head>
<div class="page-shell leaderboard-page">
	<div class="leaderboard-heading">
		<div>
			<h1 class="page-heading">Leaderboard</h1>
		</div>
		<span class="ranking-label"><Icon name="shield" size={16} />Top 10 · Club rating</span>
	</div>
	{#if players.length}
		<TopThree
			players={players.slice(0, 3).map((player: { Username: string; Rating: number }) => ({
				username: player.Username,
				score: Math.floor(player.Rating)
			}))}
			scoreLabel="Rating"
		/>
		<div class="rankings panel">
			<table>
				<thead
					><tr><th scope="col">Rank</th><th scope="col">Player</th><th scope="col">Rating</th></tr
					></thead
				><tbody
					>{#each players as player, index}<tr
							><td><span class:top-rank={index < 3}>{String(index + 1).padStart(2, '0')}</span></td
							><td
								><a href={`/member/${player.Username}`}
									><span class="table-avatar">{player.Username.slice(0, 1).toUpperCase()}</span
									>{player.Username}</a
								></td
							><td>{Math.floor(player.Rating)}</td></tr
						>{/each}</tbody
				>
			</table>
		</div>
	{:else}<div class="panel empty-state">
			<span class="empty-icon"><Icon name="trophy" size={31} /></span>
			<h2>No ranked players yet</h2>
			<a class="button button-primary" href="/play">Find a game</a>
		</div>{/if}
</div>

<style>
	.leaderboard-heading {
		display: flex;
		align-items: end;
		justify-content: space-between;
		gap: 20px;
		margin-bottom: 37px;
	}
	.ranking-label {
		display: flex;
		flex-shrink: 0;
		align-items: center;
		gap: 8px;
		padding: 10px 14px;
		border: 1px solid var(--line);
		border-radius: 7px;
		color: #a9b8aa;
		font-size: 0.8125rem;
	}
	.rankings {
		overflow: hidden;
	}
	table {
		width: 100%;
		border-collapse: collapse;
		text-align: left;
	}
	th {
		padding: 18px 26px;
		background: #242d25;
		color: #a7b5a8;
		font-size: 0.75rem;
		font-weight: 500;
		text-transform: uppercase;
		letter-spacing: 0.08em;
	}
	td {
		border-top: 1px solid var(--line);
		padding: 16px 26px;
		font-size: 0.9375rem;
	}
	td:first-child {
		width: 110px;
		color: #92a095;
		font-variant-numeric: tabular-nums;
	}
	td:last-child,
	th:last-child {
		text-align: right;
	}
	td:last-child {
		font-weight: 600;
	}
	.top-rank {
		color: #c5dba7;
	}
	td a {
		display: inline-flex;
		align-items: center;
		gap: 13px;
		overflow-wrap: anywhere;
	}
	td a:hover {
		color: var(--accent);
	}
	.table-avatar {
		display: grid;
		width: 33px;
		height: 33px;
		flex-shrink: 0;
		place-items: center;
		border-radius: 9px;
		background: #2e3f31;
		color: #b5ceb0;
		font-size: 0.8125rem;
	}
	@media (max-width: 850px) {
		.leaderboard-heading {
			align-items: start;
			flex-direction: column;
		}
	}
	@media (max-width: 600px) {
		td,
		th {
			padding: 15px 13px;
		}
		td:first-child {
			width: 60px;
		}
		td a {
			gap: 9px;
			font-size: 0.875rem;
		}
	}
</style>
