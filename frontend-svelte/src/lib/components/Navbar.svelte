<script>
	import { invalidateAll } from '$app/navigation';
	import { getBaseURL } from '$lib/utils';

	let { data } = $props();

	const handleLogout = async () => {
		await fetch(`${getBaseURL()}/logout`, {
			method: 'POST',
			credentials: 'include'
		});
		invalidateAll();
	};
</script>

<div class="flex w-full px-4 py-4 text-xl">
	<ul class="flex w-full justify-between">
		<li class="text-2xl">
			<a href="/" class="rounded-lg p-1 hover:bg-gray-600"> MNNIT Chess Club </a>
		</li>
		<li class="flex items-center justify-center gap-1 text-2xl">
			{#if !data.user}
				<a class="text-xl text-blue-300 underline" href="/login">Login</a>
			{:else}
				<a href={`/member/${data.user.username}`} class="rounded-lg p-1 hover:bg-gray-600"
					>{data.user.username}</a
				>
				<button
					onclick={handleLogout}
					class="cursor-pointer rounded-lg p-2 hover:bg-gray-600"
					aria-label="sign out"
				>
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
						class="lucide lucide-log-out-icon lucide-log-out"
						><path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4" /><polyline
							points="16 17 21 12 16 7"
						/><line x1="21" x2="9" y1="12" y2="12" /></svg
					>
				</button>
			{/if}
		</li>
	</ul>
</div>
