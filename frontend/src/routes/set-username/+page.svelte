<script lang="ts">
	import { getBaseURL } from '$lib/utils';
	import Icon from '$lib/components/Icon.svelte';

	let username = $state('');
	let loading = $state(false);
	let error = $state('');
	let uiError = $derived(!isValid(username));

	function isValid(str: string) {
		return /^[A-Za-z0-9_]+$/.test(str);
	}

	const handleSubmit = async () => {
		if (!isValid(username)) return;
		loading = true;

		const res = await fetch(`${getBaseURL()}/set-username`, {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify({ username }),
			credentials: 'include'
		});

		const response = await res.json();

		if (!res.ok) {
			loading = false;
			error = response.error;
			return;
		}

		window.location.href = '/';
	};
</script>

<svelte:head><title>Your club username — MNNIT Chess Club</title></svelte:head>
<div class="page-shell username-page">
	<form
		class="panel username-panel"
		onsubmit={(event) => {
			event.preventDefault();
			handleSubmit();
		}}
	>
		<span class="empty-icon"><Icon name="knight" size={30} /></span>
		<h1 class="page-heading">Choose a username</h1>
		<label for="username">Your club username</label><input
			id="username"
			name="username"
			class="form-field"
			type="text"
			placeholder="e.g. knight_owl"
			autocomplete="username"
			bind:value={username}
			aria-describedby="username-hint username-error"
			aria-invalid={username.length > 0 && uiError}
		/>
		<p id="username-hint" class="hint">Use letters, numbers, and underscores.</p>
		<p id="username-error" class="error-message" aria-live="polite">
			{error ||
				(username.length > 0 && uiError ? 'Use only letters, numbers, and underscores.' : '')}
		</p>
		<button class="button button-primary" type="submit" disabled={loading || uiError}
			>{loading ? 'Saving your username…' : 'Join the club'}</button
		>
	</form>
</div>

<style>
	.username-page {
		max-width: 600px;
		padding-top: 65px;
	}
	.username-panel {
		padding: 38px;
	}
	.empty-icon {
		width: 58px;
		height: 58px;
		margin-bottom: 25px;
	}
	.page-heading {
		font-size: 2.5rem;
		line-height: 1.15;
		margin-block: 15px;
	}
	label {
		display: block;
		margin: 30px 0 10px;
		color: #bfccbf;
		font-size: 0.875rem;
	}
	.hint {
		color: #97a599;
		margin-top: 12px;
		font-size: 0.8125rem;
		line-height: 1.7;
	}
	.error-message {
		min-height: 24px;
		margin-block: 10px;
		color: #f1a49d;
		font-size: 0.875rem;
	}
	.button {
		width: 100%;
	}
	@media (max-width: 500px) {
		.username-panel {
			padding: 26px 22px;
		}
	}
</style>
