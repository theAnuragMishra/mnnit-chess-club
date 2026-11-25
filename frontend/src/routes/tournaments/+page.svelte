<script lang="ts">
	import Icon from '$lib/components/Icon.svelte';
	let { data } = $props();
	const sections = $derived([
		{ title: 'Live arenas', live: true, tournaments: data.live ?? [] },
		{
			title: 'Upcoming arenas',
			live: false,
			tournaments: data.scheduled ?? []
		}
	]);
</script>

<svelte:head><title>Tournaments — MNNIT Chess Club</title></svelte:head>
<div class="page-shell">
	<h1 class="page-heading">Tournaments</h1>
	<div class="tournament-layout">
		<div class="tournament-content">
			{#if data.live === null || data.scheduled === null}<div class="service-note" role="status">
					Some arenas couldn't be loaded. <button
						class="text-link"
						onclick={() => window.location.reload()}>Try again</button
					>
				</div>{/if}
			{#if sections.every((section) => section.tournaments.length === 0)}<div
					class="panel empty-state"
				>
					<span class="empty-icon"><Icon name="trophy" size={31} /></span>
					<h2>
						{data.live === null || data.scheduled === null
							? 'No arenas available'
							: 'No live or upcoming arenas'}
					</h2>
					<a href="/play" class="button button-primary">Play a game</a>
				</div>{:else}
				{#each sections as section}{#if section.tournaments.length}<section class="arena-section">
							<div class="arena-heading">
								<h2>{section.title}</h2>
							</div>
							<div class="arena-list">
								{#each section.tournaments as tournament}<a
										class="arena-card panel"
										href={`/tournament/${tournament.ID}`}
										><span class="arena-icon"><Icon name="trophy" size={24} /></span>
										<div class="arena-details">
											<span class="arena-status" class:live={section.live}
												>{section.live ? 'Live now' : 'Upcoming'}</span
											>
											<h3>{tournament.Name}</h3>
											<p>
												{new Date(tournament.StartTime).toLocaleString('en-IN', {
													month: 'short',
													day: 'numeric',
													hour: '2-digit',
													minute: '2-digit',
													hour12: true
												})}
											</p>
										</div>
										<div class="arena-control">
											<strong>{tournament.BaseTime / 60} + {tournament.Increment}</strong><span
												>{tournament.Duration / 3600} hours</span
											>
										</div>
										<Icon name="chevron" size={19} /></a
									>{/each}
							</div>
						</section>{/if}{/each}{/if}
		</div>
	</div>
</div>

<style>
	.tournament-layout {
		max-width: 52rem;
		margin-top: 1.5rem;
	}
	.service-note {
		margin-bottom: 20px;
		border: 1px solid #676044;
		border-radius: 9px;
		padding: 16px;
		color: #d3cca9;
		font-size: 0.875rem;
		line-height: 1.7;
	}
	.service-note button {
		margin-left: 5px;
	}
	.arena-section + .arena-section {
		margin-top: 32px;
	}
	.arena-heading {
		display: flex;
		align-items: center;
		justify-content: space-between;
		margin-bottom: 17px;
	}
	.arena-heading h2 {
		font-size: 1.25rem;
		font-weight: 550;
	}
	.arena-list {
		display: grid;
		gap: 12px;
	}
	.arena-card {
		display: flex;
		align-items: center;
		gap: 18px;
		padding: 23px;
		transition:
			border-color 0.2s,
			background 0.2s;
	}
	.arena-card:hover {
		border-color: #6d9273;
		background: #253028;
	}
	.arena-icon {
		display: grid;
		width: 50px;
		height: 50px;
		flex-shrink: 0;
		place-items: center;
		border-radius: 12px;
		background: #2d3c2d;
		color: #bfd9a8;
	}
	.arena-details {
		min-width: 0;
	}
	.arena-status {
		color: #a9b695;
		font-size: 0.75rem;
	}
	.arena-status.live {
		color: #b7e4bd;
	}
	.arena-details h3 {
		margin-block: 5px 7px;
		overflow-wrap: anywhere;
		font-size: 1.0625rem;
		font-weight: 550;
	}
	.arena-details p {
		color: #95a298;
		font-size: 0.8125rem;
	}
	.arena-control {
		display: flex;
		flex-shrink: 0;
		flex-direction: column;
		gap: 7px;
		margin-left: auto;
		text-align: right;
	}
	.arena-control strong {
		font-size: 1.25rem;
		font-weight: 550;
	}
	.arena-control span {
		font-size: 0.8125rem;
		color: var(--muted);
	}
	.arena-card > :global(svg) {
		flex-shrink: 0;
		color: #7e9681;
	}
	@media (max-width: 500px) {
		.arena-card {
			gap: 11px;
			padding: 18px 14px;
		}
		.arena-icon {
			width: 38px;
			height: 42px;
		}
		.arena-details h3 {
			font-size: 0.9375rem;
		}
		.arena-control strong {
			font-size: 1rem;
		}
		.arena-card > :global(svg) {
			display: none;
		}
	}
</style>
