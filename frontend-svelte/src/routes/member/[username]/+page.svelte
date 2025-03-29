<script lang="ts">
	import { page } from '$app/state';
	let { data } = $props();

	const getX = (item: { WhiteUsername: String; BlackUsername: string; Result: string }) => {
		let x = 'bg-red-500';
		if (
			(item.WhiteUsername === page.params.username && item.Result === '1-0') ||
			(item.BlackUsername === page.params.username && item.Result === '0-1')
		) {
			x = 'bg-green-700';
		} else if (item.Result === 'ongoing') {
			x = 'bg-gray-600';
		}
	};
</script>

<div class="flex-col rounded-xl bg-black p-4 text-xl">
	<div class="mb-4 text-center text-5xl">{page.params.username}'s Games</div>
	<div class="flex w-full flex-col items-center gap-2">
		{#each data.member as item}
			<a href={`/game/${item.ID}`} class="flex w-4/5 gap-2 rounded-sm bg-gray-800 px-8 py-4">
				<span class="w-1/3 text-left">{item.WhiteUsername}</span>
				<span class={`flex w-1/3 items-center justify-center ${getX(item)}`}>
					{#if item.Result !== 'ongoing'}
						item.Result
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
	</div>
</div>
