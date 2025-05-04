<script>
	let { data } = $props();
</script>

<div class="flex h-full w-full flex-col items-center justify-center">
	{#if !data.user}
		<p>Login to play</p>
	{:else}
		<a href="/play" class="cursor-pointer rounded-lg bg-gray-500 p-3 text-4xl"> Play! </a>
	{/if}
</div>
