<script>
	import { getBaseURL } from '$lib/utils.js';
	import googleImg from '$lib/assets/icons/google.svg';
</script>

<svelte:head>
	<title>Login</title>
</svelte:head>
<div class="mt-30 flex h-full flex-col items-center justify-center gap-2">
	<h1 class="text-2xl font-bold">Login</h1>
	<button
		onclick={() => (window.location.href = `${getBaseURL()}/auth/login/google`)}
		class="flex cursor-pointer items-center gap-2 rounded-md bg-black px-4 py-2 text-xl text-white"
	>
		Login with Google <img class="w-[24px]" src={googleImg} alt="google logo" />
	</button>
	<p class="text-center">More login providers may (or may not) be added.</p>
</div>
