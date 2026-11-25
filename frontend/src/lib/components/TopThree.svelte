<script lang="ts">
	import Icon from './Icon.svelte';

	type PodiumPlayer = { username: string; score: number };
	let { players, scoreLabel }: { players: readonly PodiumPlayer[]; scoreLabel: string } = $props();

	const places = [
		{ rank: 1, label: '2nd' },
		{ rank: 0, label: '1st' },
		{ rank: 2, label: '3rd' }
	] as const;
</script>

<div class="podium" role="group" aria-label="Top three players">
	{#each places as place, column (place.rank)}
		{@const player = players[place.rank]}
		{#if player}
			<a
				href={`/member/${player.username}`}
				class="podium-card panel"
				class:first={place.rank === 0}
				style:grid-column={column + 1}
			>
				<span class="place"><Icon name="trophy" size={15} />{place.label}</span>
				<span class="avatar" aria-hidden="true">{player.username.slice(0, 1).toUpperCase()}</span>
				<strong class="player">{player.username}</strong>
				<span class="score">{player.score} <span>{scoreLabel}</span></span>
			</a>
		{/if}
	{/each}
</div>

<style>
	.podium {
		display: grid;
		grid-template-columns: repeat(3, minmax(0, 1fr));
		align-items: end;
		gap: clamp(0.375rem, 2vw, 1rem);
		margin-bottom: 1.5rem;
	}
	.podium-card {
		display: grid;
		grid-row: 1;
		justify-items: center;
		gap: 0.5rem;
		padding: 0.875rem clamp(0.25rem, 1.5vw, 1rem);
		border-radius: 12px;
		text-align: center;
	}
	.podium-card:hover {
		border-color: #819b7c;
	}
	.podium-card.first {
		padding-block: 1.375rem;
		border-color: #747d53;
		background: #2c3526;
	}
	.place {
		display: flex;
		align-items: center;
		gap: 0.25rem;
		color: var(--muted);
		font-size: 0.8125rem;
	}
	.first .place {
		color: #e2dcaa;
	}
	.avatar {
		display: grid;
		width: 2.25rem;
		height: 2.25rem;
		place-items: center;
		border-radius: 50%;
		background: #344c37;
		color: var(--accent);
		font-size: 1rem;
	}
	.first .avatar {
		width: 2.75rem;
		height: 2.75rem;
	}
	.player {
		max-width: 100%;
		overflow-wrap: anywhere;
		font-size: 0.875rem;
	}
	.score {
		color: var(--accent);
		font-size: 0.9375rem;
		font-variant-numeric: tabular-nums;
		font-weight: 600;
	}
	.score span {
		display: block;
		color: var(--muted);
		font-size: 0.75rem;
		font-weight: 400;
	}
</style>
