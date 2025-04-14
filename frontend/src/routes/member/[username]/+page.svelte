<script lang="ts">
	import { page } from '$app/state';
	import { formatResultAndReason, getBaseURL, getTimeControl } from '$lib/utils.js';
	let { data } = $props();
	// console.log('mounted');
	let pageNumber = $state(1);
	let hasMore = $state(data.member ? data.member.length == 15 : false);
	let loading = $state(false);
	const items: any = $state(data.member ? data.member : []);

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
				`${getBaseURL()}/profile/${page.params.username}?page=${pageNumber}`,
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
	<div class="mb-4 text-center text-5xl">{page.params.username}'s Games</div>
	<div class="flex w-full flex-col items-center gap-2">
		{#each items as item}
			{@const color =
				(item.WhiteUsername === page.params.username && item.Result === '1-0') ||
				(item.BlackUsername === page.params.username && item.Result === '0-1')
					? 'text-green-700'
					: item.Result === 'ongoing' || item.Result === '1/2-1/2' || item.Result === 'aborted'
						? 'text-gray-300'
						: 'text-red-500'}
			<div class="relative flex w-4/5 flex-col gap-2 rounded-sm bg-gray-800 px-8 py-4">
				<a
					aria-label="game link"
					href={`/game/${item.ID}`}
					class="absolute top-0 left-0 z-2 h-full w-full"
				></a>

				<div>
					{item.BaseTime / 60}+{item.Increment}
					{getTimeControl(item.BaseTime, item.Increment)}
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
					{#if item.Result !== 'ongoing'}
						{formatResultAndReason(item.Result, item.ResultReason)}
					{:else}
						Playing right now
					{/if}
				</div>
				<div class="text-lg">
					{Math.ceil(item.GameLength / 2)} move{`${Math.ceil(item.GameLength / 2) > 1 ? 's' : ''}`}
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
