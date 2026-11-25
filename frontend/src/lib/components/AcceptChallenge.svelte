<script lang="ts">
	import { page } from '$app/state';
	import { PUBLIC_FRONTEND_URL } from '$env/static/public';
	import { user } from '$lib/user.svelte';
	import { websocketStore } from '$lib/websocket.svelte';
	import Icon from './Icon.svelte';
	const { data } = $props();
	const gameLink = `${PUBLIC_FRONTEND_URL}/game/${page.params.gameID}/`;
	let copied = $state(false);
	let challengeInput = $state<HTMLInputElement>();
</script>

<svelte:head>
	<title
		>Challenge - {data.gameData.TimeControl.baseTime / 60}+{data.gameData.TimeControl
			.increment}</title
	>
</svelte:head>
<div class="page-shell challenge-page">
	<section class="panel challenge-card">
		<span class="empty-icon"><Icon name="people" size={29} /></span>
		<h1 class="page-heading">
			{user.username === data.gameData.CreatorUsername ? 'Challenge created' : 'Chess challenge'}
		</h1>
		<div class="time-badge">
			<Icon name="clock" size={19} />{data.gameData.TimeControl.baseTime / 60} + {data.gameData
				.TimeControl.increment}<span>minutes + seconds</span>
		</div>
		{#if user.username === data.gameData.CreatorUsername}<p>Share this link with your opponent.</p>
			<label for="challenge-url">Your challenge link</label><input
				id="challenge-url"
				bind:this={challengeInput}
				class="form-field"
				value={gameLink}
				readonly
				onclick={(event) => event.currentTarget.select()}
			/><button
				class="button button-primary"
				onclick={async () => {
					try {
						await navigator.clipboard.writeText(gameLink);
						copied = true;
						setTimeout(() => (copied = false), 2000);
					} catch {
						challengeInput?.select();
					}
				}}
				><Icon name={copied ? 'check' : 'link'} size={18} />{copied
					? 'Link copied!'
					: 'Copy challenge link'}</button
			><span class="waiting" role="status">{'Waiting for your opponent to accept…'}</span>
		{:else}<p>
				<strong>{data.gameData.CreatorUsername}</strong> sent you a challenge.
			</p>
			<button
				class="button button-primary"
				onclick={() =>
					websocketStore.sendMessage({
						type: 'accept_challenge',
						payload: { GameID: page.params.gameID }
					})}><Icon name="knight" size={20} />Accept challenge</button
			>{/if}
	</section>
</div>

<style>
	.challenge-page {
		max-width: 650px;
		padding-top: 65px;
	}
	.challenge-card {
		display: flex;
		flex-direction: column;
		align-items: center;
		padding: 38px;
		text-align: center;
	}
	.page-heading {
		font-size: clamp(2rem, 5vw, 2.75rem);
		margin: 18px 0 25px;
	}
	.time-badge {
		display: flex;
		align-items: center;
		gap: 9px;
		padding: 12px 18px;
		border: 1px solid #4b614a;
		border-radius: 9px;
		background: #29392b;
		color: #d9ebd1;
		font-size: 1.25rem;
	}
	.time-badge span {
		font-size: 0.75rem;
		color: #a5b69c;
	}
	.challenge-card > p:not(.eyebrow) {
		margin: 25px 0;
		color: var(--muted);
		font-size: 0.9375rem;
		line-height: 1.85;
	}
	label {
		width: 100%;
		text-align: left;
		margin-bottom: 10px;
		color: #afbeb0;
		font-size: 0.8125rem;
	}
	input {
		font-size: 0.875rem;
	}
	.button {
		width: 100%;
		margin-top: 18px;
	}
	.waiting {
		margin-top: 19px;
		color: #97aa98;
		font-size: 0.8125rem;
	}
	@media (max-width: 500px) {
		.challenge-card {
			padding: 28px 20px;
		}
	}
</style>
