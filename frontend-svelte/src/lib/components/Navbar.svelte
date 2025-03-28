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
			<a href="/"> MNNIT Chess Club </a>
		</li>
		<li class="flex items-center justify-center gap-1 text-2xl">
			{#if !data.user}
				<a class="text-xl text-blue-300 underline" href="/login">Login</a>
			{:else}
				<a href={`/member/${data.user.username}`} class="rounded-lg px-1 py-0.5 hover:bg-gray-600"
					>{data.user.username}</a
				>
				<button
					onclick={handleLogout}
					class="cursor-pointer rounded-lg px-1 py-0.5 hover:bg-gray-600"
					aria-label="sign out"
				>
					<svg
						xmlns="http://www.w3.org/2000/svg"
						width="32"
						height="32"
						fill="#fff"
						viewBox="0 0 256 256"
						><path
							d="M120,216a8,8,0,0,1-8,8H48a8,8,0,0,1-8-8V40a8,8,0,0,1,8-8h64a8,8,0,0,1,0,16H56V208h56A8,8,0,0,1,120,216Zm109.66-93.66-40-40a8,8,0,0,0-11.32,11.32L204.69,120H112a8,8,0,0,0,0,16h92.69l-26.35,26.34a8,8,0,0,0,11.32,11.32l40-40A8,8,0,0,0,229.66,122.34Z"
						></path></svg
					>
				</button>
			{/if}
		</li>
	</ul>
</div>
