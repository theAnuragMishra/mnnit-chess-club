<script lang="ts">
	import { goto } from '$app/navigation';
	import { navigating, page } from '$app/state';
	import { getBGColorForGameListItem } from '$lib/utils.js';
	import { untrack } from 'svelte';
	let { data } = $props();
	// console.log('mounted');
	const items: any = $state([]);
	$effect.pre(() => {
		if (!data.member) return;
		const i = data.member;
		untrack(() => items.push(...i));
	});
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
</script>

<div class="flex-col rounded-xl bg-black p-4 text-xl">
	<div class="mb-4 text-center text-5xl">{page.params.username}'s Games</div>
	<div class="flex w-full flex-col items-center gap-2">
		{#each items as item}
			<a href={`/game/${item.ID}`} class="flex w-4/5 gap-2 rounded-sm bg-gray-800 px-8 py-4">
				<span class="w-1/3 text-left">{item.WhiteUsername}</span>
				<span
					class={`flex w-1/3 items-center justify-center ${getBGColorForGameListItem(item, page.params.username)}`}
				>
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
		{#if data.hasMore}
			<div
				class="h-[20px] bg-transparent"
				use:intersect={() =>
					!navigating.to &&
					goto(`/member/${page.params.username}?page=${data.page + 1}`, {
						replaceState: true,
						noScroll: true,
						keepFocus: true
					})}
			></div>
		{/if}
	</div>
</div>
