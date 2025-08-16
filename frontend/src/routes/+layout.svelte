<script lang="ts">
	import Navbar from '$lib/components/Navbar.svelte';
	import '../app.css';

	let { data, children } = $props();
</script>

<main class="flex min-h-screen flex-col items-center overflow-x-hidden bg-gray-900 text-white">
	<Navbar {data} />
	{@render children()}
</main>
