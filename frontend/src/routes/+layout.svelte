<script lang="ts">
	import Navbar from '$lib/components/Navbar.svelte';
	import Footer from '$lib/components/Footer.svelte';
	import { page } from '$app/state';
	import '../app.css';
	let { children } = $props();
</script>

<div class="app-shell">
	<a href="#main-content" class="skip-link">Skip to content</a>
	<Navbar />
	<main id="main-content" class="site-content" tabindex="-1">{@render children()}</main>
	{#if !page.url.pathname.startsWith('/game/') && page.url.pathname !== '/board/'}<Footer />{/if}
</div>
