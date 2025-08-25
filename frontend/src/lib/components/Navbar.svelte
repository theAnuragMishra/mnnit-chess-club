<script>
	import { invalidateAll } from '$app/navigation';
	import { getBaseURL } from '$lib/utils';
	import logo from '$lib/assets/mcc-logo.webp';
	import signoutImg from '$lib/assets/icons/logout.svg';
	import menuImg from '$lib/assets/icons/menu.svg';
	import crossImg from '$lib/assets/icons/cross.svg';

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
				<img src={logo} alt="mcc-logo" class="h-10 invert-100" />
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
						<img alt="signout" src={signoutImg} />
					</button>
				{/if}
			</div>

			<button
				class="cursor-pointer sm:hidden"
				onclick={() => (expanded = !expanded)}
				aria-label="menu"
			>
				<img alt="menu" src={menuImg} />
			</button>
		</div>
	</div>

	{#if expanded}
		<button
			class="bg-opacity-50 fixed inset-0 z-40 bg-black"
			onclick={() => (expanded = false)}
			aria-label="x"
		></button>
	{/if}

	<div
		class={`fixed top-0 left-0 z-50 h-full w-64 transform bg-gray-800 transition-transform duration-200 ease-in-out ${
			expanded ? 'translate-x-0' : '-translate-x-full'
		}`}
	>
		<div class="space-y-4 p-4">
			<button
				class="mb-6 cursor-pointer"
				onclick={() => (expanded = false)}
				aria-label="close menu"
			>
				<img src={crossImg} alt="close menu" />
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
					Logout <img src={signoutImg} alt="logout" />
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
