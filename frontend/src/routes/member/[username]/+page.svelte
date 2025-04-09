<script lang="ts">
	import { page } from '$app/state';
	import { getBaseURL } from '$lib/utils.js';
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

<div class="flex w-full flex-col rounded-xl bg-black p-4 text-xl">
	<div class="mb-4 text-center text-5xl">{page.params.username}'s Games</div>
	<div class="flex w-full flex-col items-center gap-2">
		{#each items as item}
			{@const color =
				(item.WhiteUsername === page.params.username && item.Result === '1-0') ||
				(item.BlackUsername === page.params.username && item.Result === '0-1')
					? 'bg-green-700'
					: item.Result === 'ongoing' || item.Result === '1/2-1/2'
						? 'bg-gray-600'
						: 'bg-red-500'}
			<a href={`/game/${item.ID}`} class="flex w-4/5 gap-2 rounded-sm bg-gray-800 px-8 py-4">
				<span class="w-1/3 text-left">{item.WhiteUsername}</span>
				<span class={`flex w-1/3 items-center justify-center ${color}`}>
					{#if item.Result !== 'ongoing'}
						{item.Result}
					{:else}
						<svg
							xmlns="http://www.w3.org/2000/svg"
							width="24"
							height="24"
							viewBox="0 0 24 24"
							fill="none"
							stroke="currentColor"
							stroke-width="2"
							stroke-linecap="round"
							stroke-linejoin="round"
							class="lucide lucide-asterisk-icon lucide-asterisk"
							><path d="M12 6v12" /><path d="M17.196 9 6.804 15" /><path
								d="m6.804 9 10.392 6"
							/></svg
						>
					{/if}
				</span>
				<span class="w-1/3 text-right">{item.BlackUsername}</span></a
			>
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
