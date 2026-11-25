<script lang="ts">
	import hldImage from '$lib/assets/mcc_hld.svg';
</script>

<svelte:head><title>Server architecture — MNNIT Chess Club</title></svelte:head>
<div class="page-shell">
	<h1 class="page-heading">Server architecture.</h1>
	<div class="architecture panel">
		<img src={hldImage} alt="High level design of the MNNIT Chess Club server architecture" />
	</div>
</div>

<style>
	.architecture {
		margin-top: 30px;
		padding: 24px;
		overflow-x: auto;
		background: #f3f4ed;
	}
	img {
		display: block;
		width: 100%;
		min-width: 620px;
		max-height: 75vh;
		object-fit: contain;
	}
</style>
