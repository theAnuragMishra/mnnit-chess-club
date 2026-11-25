<script lang="ts">
	import { page } from '$app/state';
	import Icon from '$lib/components/Icon.svelte';
	import {
		formatPostgresTimestamp,
		formatResultAndMethod,
		getBaseURL,
		getTimeControl
	} from '$lib/utils.js';
	import trophyImg from '$lib/assets/icons/trophy.svg';
	let { data } = $props();
	// console.log('mounted');
	let pageNumber = $state(1);
	let hasMore = $state(data.games ? data.games.length == 15 : false);
	let loading = $state(false);
	const items: any = $state(data.games ?? []);

	const memberSince = new Date(data.profile.CreatedAt);

	function intersect(node: any, callback: any) {
		const observer = new IntersectionObserver(
			(entries) => {
				for (const entry of entries) {
					if (entry.isIntersecting) {
						callback();
						break;
					}
				}
			},
			{ threshold: 1 }
		);

		observer.observe(node);

		return {
			destroy() {
				observer.disconnect();
			}
		};
	}

	async function fetchGames() {
		if (!hasMore) return;
		loading = true;
		try {
			const response = await fetch(
				`${getBaseURL()}/games/${page.params.username}?page=${pageNumber}`,
				{
					credentials: 'include'
				}
			);
			const memberData = await response.json();
			if (memberData) items.push(...memberData);
			loading = false;
			hasMore = memberData && memberData.length == 15;
		} catch (e) {
			loading = false;
			console.error(e);
		}
	}
</script>

<svelte:head>
	<title>Profile - {page.params.username}</title>
</svelte:head>
<div class="page-shell profile-page">
	<section class="profile-header panel">
		<span class="profile-avatar">{page.params.username?.slice(0, 1).toUpperCase()}</span>
		<div class="profile-identity">
			<h1>{page.params.username}</h1>
			<p>
				At the board since {memberSince.toLocaleDateString('en-GB', {
					month: 'long',
					year: 'numeric'
				})}
			</p>
		</div>
		<div class="profile-rating">
			<strong>{Math.floor(data.profile.Rating)}{data.profile.Rd > 110 ? '?' : ''}</strong><span
				title="Rating deviation measures rating uncertainty. A rating is provisional when deviation is above 110."
				>Club rating · RD {Math.ceil(data.profile.Rd)}</span
			>
		</div>
	</section>
	<div class="profile-stats">
		<div><strong>{data.profile.GameCount}</strong><span>Games played</span></div>
		<div><strong class="wins">{data.profile.WinCount}</strong><span>Wins</span></div>
		<div><strong>{data.profile.DrawCount}</strong><span>Draws</span></div>
		<div><strong>{data.profile.LossCount}</strong><span>Losses</span></div>
	</div>
	<h2 class="history-heading">Game history</h2>
	{#if items.length === 0}<div class="panel empty-state">
			<span class="empty-icon"><Icon name="grid" size={28} /></span>
			<h2>No games yet</h2>
			<a class="button button-primary" href="/play">Find a game</a>
		</div>{/if}

	<div class="flex w-full flex-col items-center gap-2">
		{#each items as item}
			{@const color =
				(item.WhiteUsername === page.params.username && item.Result === 1) ||
				(item.BlackUsername === page.params.username && item.Result === 2)
					? 'positive'
					: item.Result === 0 || item.Result === 3 || item.Result === 4
						? 'text-gray-300'
						: 'negative'}
			<div class="game-record relative flex w-full flex-col gap-3 px-6 py-5">
				<a
					aria-label="game link"
					href={`/game/${item.ID}`}
					class="absolute top-0 left-0 z-2 h-full w-full"
				></a>
				{#if item.TournamentID}<a
						class="absolute top-[15px] right-[20px] z-3 flex h-fit w-fit items-center gap-[5px]"
						href={`/tournament/${item.TournamentID}`}
					>
						<img
							class="h-[20px] md:h-[30px]"
							src={trophyImg}
							alt="trophy"
						/>{item.TournamentName}</a
					>{/if}
				<div class="flex items-center justify-between">
					{item.BaseTime / 60}+{item.Increment}
					{getTimeControl(item.BaseTime, item.Increment)}
				</div>
				<div class="flex w-full items-center justify-center gap-5">
					<div class="flex flex-col items-center justify-center">
						<span class="md:text-2xl">{item.WhiteUsername}</span><span class="text-[16px]"
							><span>{item.RatingW}</span>&nbsp;&nbsp;<span
								class={`${item.ChangeW > 0 ? 'positive' : item.ChangeW < 0 ? 'negative' : ''}`}
								>{`${item.ChangeW > 0 ? '+' : ''}`}{item.ChangeW}</span
							></span
						>
					</div>
					<div>Vs</div>
					<div class="flex flex-col items-center justify-center">
						<span class="md:text-2xl">{item.BlackUsername}</span><span class="text-[16px]"
							><span>{item.RatingB}</span>&nbsp;&nbsp;<span
								class={`${item.ChangeB > 0 ? 'positive' : item.ChangeB < 0 ? 'negative' : ''}`}
								>{`${item.ChangeB > 0 ? '+' : ''}`}{item.ChangeB}</span
							></span
						>
					</div>
				</div>
				<div class={`text-[14px] ${color} w-full text-center`}>
					{#if item.Result !== 0}
						{formatResultAndMethod(item.Result, item.Method)}
					{:else}
						Playing right now
					{/if}
				</div>
				<div class="flex items-center justify-between">
					<span
						>{Math.ceil(item.GameLength / 2)} move{`${Math.ceil(item.GameLength / 2) > 1 ? 's' : ''}`}</span
					>
					<span class="text-sm">{formatPostgresTimestamp(new Date(item.CreatedAt))}</span>
				</div>
			</div>
		{/each}
		{#if hasMore}
			<div
				class="h-[20px] bg-transparent"
				use:intersect={() => {
					pageNumber += 1;
					fetchGames();
				}}
			>
				{#if loading}Loading...{/if}
			</div>
		{/if}
	</div>
</div>

<style>
	.positive {
		color: var(--accent);
	}
	.negative {
		color: var(--danger);
	}
	.profile-header {
		display: flex;
		align-items: center;
		gap: 22px;
		margin-top: 24px;
		padding: 28px;
	}
	.profile-avatar {
		display: grid;
		width: 75px;
		height: 75px;
		flex-shrink: 0;
		place-items: center;
		border: 1px solid #587550;
		border-radius: 20px;
		background: #344e35;
		color: #cee9be;
		font-family: Georgia, serif;
		font-size: 2.25rem;
	}
	.profile-identity {
		min-width: 0;
	}
	.profile-identity h1 {
		overflow-wrap: anywhere;
		font-family: Georgia, serif;
		font-size: clamp(1.5rem, 4vw, 2.5rem);
		letter-spacing: -0.035em;
	}
	.profile-identity p {
		margin-top: 7px;
		color: #a0aea2;
		font-size: 0.875rem;
	}
	.profile-rating {
		display: flex;
		flex-shrink: 0;
		flex-direction: column;
		align-items: end;
		gap: 5px;
		margin-left: auto;
	}
	.profile-rating strong {
		font-size: 2.5rem;
		font-weight: 550;
		letter-spacing: -0.04em;
	}
	.profile-rating span {
		color: #99ac99;
		font-size: 0.8125rem;
	}
	.profile-stats {
		display: grid;
		grid-template-columns: repeat(4, 1fr);
		gap: 16px;
		margin-top: 20px;
	}
	.profile-stats > div {
		display: flex;
		flex-direction: column;
		gap: 6px;
		border: 1px solid var(--line);
		border-radius: 10px;
		padding: 22px;
		background: var(--surface);
	}
	.profile-stats strong {
		font-size: 1.75rem;
		font-weight: 500;
	}
	.profile-stats strong.wins {
		color: var(--accent);
	}
	.profile-stats span {
		color: #9eaca0;
		font-size: 0.8125rem;
	}
	.history-heading {
		margin: 38px 0 22px;
		font-family: Georgia, serif;
		font-size: 1.875rem;
	}
	.game-record {
		border: 1px solid var(--line);
		border-radius: 12px;
		background: var(--surface);
		transition: border-color 0.15s;
	}
	.game-record:hover {
		border-color: #687f65;
	}
	@media (max-width: 600px) {
		.profile-header {
			flex-wrap: wrap;
			gap: 16px;
			padding: 20px;
		}
		.profile-avatar {
			width: 55px;
			height: 55px;
			font-size: 1.5rem;
		}
		.profile-rating {
			width: 100%;
			flex-direction: row;
			align-items: baseline;
			gap: 15px;
		}
		.profile-rating strong {
			font-size: 2rem;
		}
		.profile-stats {
			gap: 9px;
		}
		.profile-stats > div {
			padding: 15px 10px;
		}
		.profile-stats span {
			font-size: 0.75rem;
		}
		.game-record {
			padding-inline: 16px;
		}
	}
</style>
