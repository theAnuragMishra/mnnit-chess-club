<script lang="ts">
	import Navbar from '$lib/components/Navbar.svelte';
	import { websocketStore } from '$lib/websocket';
	import { onMount } from 'svelte';
	import '../app.css';

	let { data, children } = $props();
	onMount(async () => {
		await websocketStore.connect();
	});
</script>

<main class="flex min-h-screen flex-col items-center bg-gray-900 p-2 text-white">
	<Navbar {data} />
	{@render children()}
</main>
