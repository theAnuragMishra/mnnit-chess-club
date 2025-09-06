<script lang="ts">
	import { getBaseURL } from '$lib/utils';

	let username = $state('');
	let loading = $state(false);
	let error = $state('');
	let uiError = $derived(!isValid(username));

	function isValid(str: string) {
		return /^[A-Za-z0-9_]+$/.test(str);
	}

	const handleSubmit = async () => {
		if (!isValid(username)) return;
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

		window.location.href = '/';
	};
</script>

<svelte:head>
	<title>Set your mcc username</title>
</svelte:head>
<div class="mt-10 flex flex-col items-center justify-center gap-4 rounded-lg p-6">
	<input
		type="text"
		placeholder="Enter your username"
		bind:value={username}
		class="rounded-md border border-gray-300 bg-gray-100 px-4 py-2 text-black focus:border-transparent focus:ring-2 focus:ring-blue-400 focus:outline-none"
	/>

	{#if error}
		<p class="text-sm text-red-500">{error}</p>
	{/if}
	{#if uiError}
		<p class="text-sm text-red-500">
			Username may only contain alphanumeric characters and underscores
		</p>
	{/if}
	<button
		class="w-[200px] rounded-md bg-blue-500 py-2 text-white transition-all hover:bg-blue-600 disabled:cursor-not-allowed disabled:bg-gray-400"
		onclick={handleSubmit}
		disabled={loading || uiError}
	>
		{loading ? 'Submitting...' : 'Submit'}
	</button>
</div>
