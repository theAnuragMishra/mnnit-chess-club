<script>
	import { page } from '$app/state';
	let { data } = $props();

	const getX = (item) => {
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
					{item.Result !== 'ongoing' ? item.Result : '*'}
				</span>
				<span class="w-1/3 text-right">{item.BlackUsername}</span></a
			>
		{/each}
	</div>
</div>
