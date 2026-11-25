<script>
	import { page } from '$app/state';
	import Icon from '$lib/components/Icon.svelte';
</script>

<svelte:head><title>Page unavailable — MNNIT Chess Club</title></svelte:head>
<div class="page-shell">
	<div class="panel empty-state">
		<span class="empty-icon"><Icon name="knight" size={32} /></span>
		<h1 class="page-heading">
			{page.status === 404 ? 'Page not found' : 'Could not load this page'}
		</h1>
		<p>
			{page.status === 404 ? 'Check the URL or return home.' : 'Try again in a moment.'}
		</p>
		<div class="error-actions">
			<a href="/" class="button button-secondary">Back to home</a>{#if page.status !== 404}<button
					class="button button-primary"
					onclick={() => window.location.reload()}>Try again</button
				>{/if}
		</div>
	</div>
</div>

<style>
	.error-actions {
		display: flex;
		flex-wrap: wrap;
		justify-content: center;
		gap: 12px;
	}
	.page-heading {
		font-size: clamp(2rem, 5vw, 3rem);
	}
</style>
