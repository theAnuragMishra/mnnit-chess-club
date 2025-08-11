<script lang="ts">
	import { page } from '$app/state';
	import {
		formatPostgresTimestamp,
		formatResultAndReason,
		getBaseURL,
		getTimeControl
	} from '$lib/utils.js';
	let { data } = $props();
	// console.log('mounted');
	let pageNumber = $state(1);
	let hasMore = $state(data.games ? data.games.length == 15 : false);
	let loading = $state(false);
	const items: any = $state(data.games ? data.games : []);

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

<div class="flex w-full flex-col rounded-xl bg-black p-4 text-xl text-gray-300">
	<div class="mb-4">
		<div class="text-3xl">{page.params.username}</div>
		<div class="text-[1rem]">
			<p>
				Member since: {memberSince.toLocaleDateString('en-GB', {
					day: 'numeric',
					month: 'long',
					year: 'numeric'
				})}
			</p>
			<p>Rating: {Math.floor(data.profile.Rating)}{`${data.profile.Rd > 110 ? '?' : ''}`}</p>
			<p
				class="w-fit cursor-help"
				title="Rating deviation shows how stable your rating is, lower is more stable. Your rating is provisional if it is above 110."
			>
				Rating Deviation: <span>{Math.ceil(data.profile.Rd)}</span>
			</p>
			<p>
				{data.profile.GameCount} Game{`${data.profile.GameCount > 1 ? 's' : ''}`}, {data.profile
					.WinCount} win{`${data.profile.WinCount > 1 ? 's' : ''}`},
				{data.profile.LossCount} loss{`${data.profile.LossCount > 1 ? 'es' : ''}`},
				{data.profile.DrawCount} draw{`${data.profile.DrawCount > 1 ? 's' : ''}`}
			</p>
		</div>
	</div>
	<div class="flex w-full flex-col items-center gap-2">
		{#each items as item}
			{@const color =
				(item.WhiteUsername === page.params.username && item.Result === 1) ||
				(item.BlackUsername === page.params.username && item.Result === 2)
					? 'text-green-500'
					: item.Result === 0 || item.Result === 3 || item.Result === 4
						? 'text-gray-300'
						: 'text-red-500'}
			<div class="relative flex w-full flex-col gap-2 rounded-sm bg-gray-800 px-8 py-4 md:w-4/5">
				<a
					aria-label="game link"
					href={`/game/${item.ID}`}
					class="z-2 absolute left-0 top-0 h-full w-full"
				></a>

				<div class="flex items-center justify-between">
					<span
						>{item.BaseTime / 60}+{item.Increment}
						{getTimeControl(item.BaseTime, item.Increment)}</span
					><span class="relative"
						>{#if item.TournamentID}<a
								class="z-3 absolute right-0 top-0 flex h-fit w-fit items-center gap-[5px]"
								href={`/tournament/${item.TournamentID}`}
							>
								<svg
									class="h-[30px]"
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
								>{item.TournamentName}</a
							>{/if}</span
					>
				</div>
				<div class="flex w-full items-center justify-center gap-5">
					<div class="flex flex-col items-center justify-center">
						<span class="text-2xl">{item.WhiteUsername}</span><span class="text-[16px]"
							><span>{item.RatingW}</span>&nbsp;&nbsp;<span
								class={`${item.ChangeW > 0 ? 'text-green-500' : item.ChangeW < 0 ? 'text-red-500' : ''}`}
								>{`${item.ChangeW > 0 ? '+' : ''}`}{item.ChangeW}</span
							></span
						>
					</div>
					<div>Vs</div>
					<div class="flex flex-col items-center justify-center">
						<span class="text-2xl">{item.BlackUsername}</span><span class="text-[16px]"
							><span>{item.RatingB}</span>&nbsp;&nbsp;<span
								class={`${item.ChangeB > 0 ? 'text-green-500' : item.ChangeB < 0 ? 'text-red-500' : ''}`}
								>{`${item.ChangeB > 0 ? '+' : ''}`}{item.ChangeB}</span
							></span
						>
					</div>
				</div>
				<div class={`text-[14px] ${color} w-full text-center`}>
					{#if item.Result !== 0}
						{formatResultAndReason(item.Result, item.ResultReason)}
					{:else}
						Playing right now
					{/if}
				</div>
				<div class="flex items-center justify-between text-lg">
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
