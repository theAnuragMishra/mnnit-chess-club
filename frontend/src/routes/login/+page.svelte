<script>
	import { getBaseURL } from '$lib/utils.js';
	import googleImg from '$lib/assets/icons/google.svg';
	import BoardPreview from '$lib/components/BoardPreview.svelte';
</script>

<svelte:head><title>Join the club — MNNIT Chess Club</title></svelte:head>
<div class="page-shell login-layout">
	<section class="login-copy">
		<h1 class="page-heading">Log in</h1>
		<div class="login-panel panel">
			<button
				class="google-button"
				onclick={() => (window.location.href = `${getBaseURL()}/auth/login/google`)}
				><img src={googleImg} alt="" width="21" height="21" />Continue with Google</button
			>
		</div>
		<a href="/" class="text-link">Back to home</a>
	</section>
	<div class="login-board">
		<BoardPreview />
	</div>
</div>

<style>
	.login-layout {
		display: grid;
		grid-template-columns: 1.1fr 1fr;
		align-items: center;
		gap: 100px;
		padding-block: 2.5rem;
	}
	.login-copy .page-heading {
		margin-block: 22px;
		line-height: 1.15;
	}
	.login-panel {
		margin-block: 30px 24px;
		padding: 28px;
	}
	.google-button {
		display: flex;
		width: 100%;
		min-height: 50px;
		align-items: center;
		justify-content: center;
		gap: 12px;
		border: 1px solid #d8e0d4;
		border-radius: 8px;
		background: #eef2e9;
		color: #202a21;
		font-size: 0.875rem;
		font-weight: 600;
		transition: background 0.15s;
	}
	.google-button:hover {
		background: #fff;
	}
	.login-board {
		isolation: isolate;
		padding-right: 15px;
	}
	@media (max-width: 1000px) {
		.login-layout {
			gap: 45px;
		}
	}
	@media (max-width: 750px) {
		.login-layout {
			grid-template-columns: 1fr;
			max-width: 520px;
			padding-block: 45px;
		}
		.login-board {
			display: none;
		}
	}
</style>
