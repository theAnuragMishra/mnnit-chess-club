<script>
	let { data } = $props();
</script>

<div class="flex h-full w-full flex-col items-center justify-center">
	{#if !data.user}
		<a href="/login" class="text-xl text-blue-500">Login to play</a>
	{:else}
		<a href="/play" class="cursor-pointer rounded-lg bg-gray-500 p-3 text-4xl"> Play! </a>
	{/if}
</div>
