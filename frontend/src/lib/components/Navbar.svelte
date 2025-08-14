<script>
	import { invalidateAll } from '$app/navigation';
	import { getBaseURL } from '$lib/utils';
	import logo from '$lib/assets/mcc-logo.webp';

	let { data } = $props();

	const handleLogout = async () => {
		await fetch(`${getBaseURL()}/logout`, {
			method: 'POST',
			credentials: 'include'
		});
		invalidateAll();
	};

	let expanded = $state(false);
</script>

<nav class="relative z-50 mt-2 w-full bg-gray-900 text-xl text-white">
	<div class="flex items-center justify-between px-4">
		<div class="flex flex-shrink-0 items-center">
			<a href="/">
				<img src={logo} alt="mcc-logo" class="invert-100 h-10" />
			</a>
		</div>

		<div class="hidden flex-1 justify-center space-x-6 sm:flex">
			<a href="/play">Play</a>
			<a href="/tournaments">Tournaments</a>
			<a href="/achievements">Achievements</a>
			<a href="/leaderboard">Leaderboard</a>
		</div>

		<div class="flex items-center gap-4">
			<div class="hidden items-center gap-1 sm:flex md:text-2xl">
				{#if !data?.user}
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

			<button
				class="cursor-pointer sm:hidden"
				onclick={() => (expanded = !expanded)}
				aria-label="menu"
			>
				<svg
					xmlns="http://www.w3.org/2000/svg"
					class="h-6 w-6"
					fill="none"
					viewBox="0 0 24 24"
					stroke="currentColor"
				>
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						stroke-width="2"
						d="M4 6h16M4 12h16M4 18h16"
					/>
				</svg>
			</button>
		</div>
	</div>

	{#if expanded}
		<button
			class="fixed inset-0 z-40 bg-black bg-opacity-50"
			onclick={() => (expanded = false)}
			aria-label="x"
		></button>
	{/if}

	<div
		class={`fixed left-0 top-0 z-50 h-full w-64 transform bg-gray-800 transition-transform duration-200 ease-in-out ${
			expanded ? 'translate-x-0' : '-translate-x-full'
		}`}
	>
		<div class="space-y-4 p-4">
			<button
				class="mb-6 cursor-pointer"
				onclick={() => (expanded = false)}
				aria-label="close menu"
			>
				<svg
					xmlns="http://www.w3.org/2000/svg"
					class="h-6 w-6"
					fill="none"
					viewBox="0 0 24 24"
					stroke="currentColor"
				>
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						stroke-width="2"
						d="M6 18L18 6M6 6l12 12"
					/>
				</svg>
			</button>

			<a href="/play" class="block" onclick={() => (expanded = false)}>Play</a>
			<a href="/tournaments" class="block" onclick={() => (expanded = false)}>Tournaments</a>
			<a href="/achievements" class="block" onclick={() => (expanded = false)}>Achievements</a>
			<a href="/leaderboard" class="block" onclick={() => (expanded = false)}>Leaderboard</a>
			<hr class="border-gray-700" />
			{#if !data?.user}
				<a href="/login" class="block" onclick={() => (expanded = false)}>Login</a>
			{:else}
				<a href={`/member/${data.user.username}`} class="block" onclick={() => (expanded = false)}
					>{data.user.username}</a
				>
				<button
					onclick={() => {
						handleLogout();
						expanded = false;
					}}
					class="flex cursor-pointer items-center justify-center gap-2 px-[0.5rem] py-[0.25rem]"
				>
					Logout <svg
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
</nav>

<style>
	a {
		border-radius: 0.5rem;
		padding: 0.25rem 0.5rem;
	}
	a:hover {
		background-color: #1e293b;
	}
</style>
