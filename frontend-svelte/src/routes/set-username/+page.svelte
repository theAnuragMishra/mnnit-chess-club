<script>
	import { goto } from '$app/navigation';
	import { getBaseURL } from '$lib/utils';

	let { data } = $props();

	//Todo: check if this works
	if (data.user.username) goto('/');

	let username = $state('');
	let loading = $state(false);
	let error = $state('');

	const handleSubmit = async () => {
		loading = true;

		const res = await fetch(`${getBaseURL()}/set-username`, {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify({ username }),
			credentials: 'include'
		});

		const response = await res.json();

		if (!res.ok) {
			loading = false;
			error = response.error;
			return;
		}

		goto('/');
	};
</script>

<div
	class="mt-10 flex flex-col items-center justify-center gap-4 rounded-lg bg-gray-100 p-6 shadow-md"
>
	<input
		type="text"
		placeholder="Enter your username"
		bind:value={username}
		class="rounded-md border border-gray-300 px-4 py-2 focus:border-transparent focus:ring-2 focus:ring-blue-400 focus:outline-none"
	/>

	{#if error}
		<p class="text-sm text-red-500">{error}</p>
	{/if}
	<button
		class="w-[200px] rounded-md bg-blue-500 py-2 text-white transition-all hover:bg-blue-600 disabled:cursor-not-allowed disabled:bg-gray-400"
		onclick={handleSubmit}
		disabled={loading}
	>
		{loading ? 'Submitting...' : 'Submit'}
	</button>
</div>
