<script lang="ts">
	import { goto } from '$app/navigation';
	import { getBaseURL } from '$lib/utils';

	const { data } = $props();

	if (!data.user) goto('/');

	async function handleInitGame(timerCode: number) {
		await fetch(`${getBaseURL()}/game/init`, {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			credentials: 'include',
			body: JSON.stringify({ username: data.user?.username, timerCode })
		});
	}
</script>

<div class="flex h-[500px] w-full items-center justify-center gap-5">
	<button
		onclick={() => handleInitGame(2)}
		class="h-[100px] w-[150px] cursor-pointer bg-gray-400 text-xl"
	>
		3+2
	</button>
	<button
		onclick={() => handleInitGame(3)}
		class="h-[100px] w-[150px] cursor-pointer bg-gray-400 text-xl"
	>
		10+0
	</button>
	<button
		onclick={() => handleInitGame(1)}
		class="h-[100px] w-[150px] cursor-pointer bg-gray-400 text-xl"
	>
		1+0
	</button>
</div>
