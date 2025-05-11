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

<div class="flex w-full justify-between px-4 py-1 text-xl">
	<div class="flex w-[200px] items-center justify-start text-xl md:text-2xl">
		<a href="/"> MNNIT Chess Club </a>
	</div>
	<div class="flex items-center justify-center gap-5 text-sm md:text-xl">
		<a href="/play">Play</a>
		<a href="/about">About</a>
		<a href="/tournaments">Tournaments</a>
		<a href="/achievements">Achievements</a>
		<a href="/leaderboard">Leaderboard</a>
	</div>
	<div class="flex w-[200px] items-center justify-end gap-1 text-xl md:text-2xl">
		{#if !data.user}
			<a href="/login">Login</a>
		{:else}
			<a href={`/member/${data.user.username}`}>{data.user.username}</a>
			<button
				onclick={handleLogout}
				class="cursor-pointer rounded-lg hover:bg-gray-800"
				aria-label="sign out"
			>
				<svg
					xmlns="http://www.w3.org/2000/svg"
					viewBox="0 0 24 24"
					fill="none"
					stroke="currentColor"
					stroke-width="2"
					stroke-linecap="round"
					stroke-linejoin="round"
					class="h-[20px] w-[20px] md:h-[20px] md:w-[20px]"
					><path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4" /><polyline
						points="16 17 21 12 16 7"
					/><line x1="21" x2="9" y1="12" y2="12" /></svg
				>
			</button>
		{/if}
	</div>
</div>

<style>
	a {
		border-radius: 0.5rem;
	}
	a:hover {
		background-color: #1e2939;
	}
</style>
