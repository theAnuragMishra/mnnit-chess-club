<script>
	let { data } = $props();
</script>

<div class="flex flex-col items-center justify-center gap-5 px-2 py-4 md:w-1/2">
	<p>
		Welcome! This is mnnit chess club's official website. You can play against your friends or
		participate in weekly arenas. We will decorate the home page better soon!
	</p>
	<p>
		In the mean time, you can take a look at some of these:
		<a class="text-blue-700" href="https://github.com/theAnuragMishra/mnnit-chess-club"
			>source code</a
		>,
		<a class="text-blue-700" href="/dev">a high level design of the server</a>.
	</p>
	<div class="flex w-full flex-col gap-2">
		<p>Or you can:</p>
		{#if !data.user}
			<a href="/login" class="w-fit text-xl text-blue-500">Login to play</a>
		{:else}
			<a href="/play" class="w-fit cursor-pointer rounded-lg bg-gray-500 p-3 text-4xl"> Play! </a>
		{/if}
		<p>
			Heads up: matchmaking may not be instant because we don't have users always trying to match on
			the site. You can create and share challenge links to your friends though!
		</p>
	</div>
</div>
