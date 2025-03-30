<script lang="ts">
	import Navbar from '$lib/components/Navbar.svelte';
	import { websocketStore } from '$lib/websocket';
	import '../app.css';

	let { data, children } = $props();

	console.log(websocketStore.socket?.readyState);
</script>

<main class="min-h-screen flex-col items-center bg-gray-900 text-white">
	<Navbar {data} />
	{@render children()}
</main>
