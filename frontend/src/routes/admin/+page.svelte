<script lang="ts">
	import { invalidateAll } from '$app/navigation';
	import { dateTimeToDate, getBaseURL } from '$lib/utils';

	const { data } = $props();
	//console.log(data);
	let name = $state('');
	let duration: number | undefined = $state();
	let date: string | undefined = $state();
	let time: string | undefined = $state();
	let startTime = $derived(dateTimeToDate(date, time));
	let allowBerserk = $state('f');

	let startLoading = $state(false);
	let createLoading = $state(false);
	let deleteLoading = $state(false);

	let createError = $state('');
	let actionError = $state('');
	let nameError = $derived(!name || name.length > 90);
	let durationError = $derived(
		duration === undefined ||
			duration === null ||
			isNaN(Number(duration)) ||
			duration <= 0 ||
			duration > 24
	);

	const handleCreateTournament = async () => {
		if (durationError || nameError || !startTime || startTime <= new Date()) return;
		createLoading = true;
		const res = await fetch(`${getBaseURL()}/admin/create-tournament`, {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify({
				name,
				baseTime: baseTimes[baseIndex] * 60,
				increment: increments[incrementIndex],
				duration: Math.ceil(duration! * 60 * 60),
				startTime: startTime.toISOString(),
				berserkAllowed: allowBerserk === 't'
			}),
			credentials: 'include'
		});

		const response = await res.json();
		createLoading = false;
		if (!res.ok) {
			createError = response.error;
			return;
		}
		invalidateAll();
	};

	const handleDeleteTournament = async (id: string) => {
		actionError = '';
		deleteLoading = true;
		const res = await fetch(`${getBaseURL()}/admin/delete-tournament`, {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify({
				tournamentID: id
			}),
			credentials: 'include'
		});
		deleteLoading = false;
		if (!res.ok) {
			actionError = 'Could not delete the tournament.';
			return;
		}
		invalidateAll();
	};

	const handleStartTournament = async (id: string) => {
		actionError = '';
		startLoading = true;
		try {
			const res = await fetch(`${getBaseURL()}/admin/start-tournament`, {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ tournamentID: id }),
				credentials: 'include'
			});
			if (!res.ok) throw new Error();
			await invalidateAll();
		} catch {
			actionError = 'Could not start the tournament.';
		} finally {
			startLoading = false;
		}
	};

	const baseTimes = [
		0.25, 0.5, 0.75, 1, 1.5, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 25,
		30, 35, 40, 45, 60, 75, 90, 105, 120, 135, 150, 165, 180
	];

	const increments = [
		0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 25, 30, 35, 40, 45,
		60, 90, 120, 150, 180
	];

	let baseIndex = $state(6);
	let incrementIndex = $state(0);
</script>

<svelte:head><title>Manage tournaments — MNNIT Chess Club</title></svelte:head>

<div class="page-shell admin-page">
	<h1 class="page-heading">Manage tournaments</h1>
	<div class="admin-layout">
		<form
			class="admin-form panel"
			onsubmit={(event) => {
				event.preventDefault();
				handleCreateTournament();
			}}
		>
			<h2>Create tournament</h2>
			<label
				>Name <input
					type="text"
					bind:value={name}
					class="form-field"
					maxlength="90"
					required
				/></label
			>
			<label
				>Duration (hours) <input
					type="number"
					bind:value={duration}
					class="form-field"
					min="0"
					max="24"
					step="any"
					required
				/></label
			>
			<label>
				<span>Minutes per side <strong>{baseTimes[baseIndex]}</strong></span>
				<input type="range" bind:value={baseIndex} min="0" max={baseTimes.length - 1} />
			</label>
			<label>
				<span>Increment <strong>{increments[incrementIndex]} sec</strong></span>
				<input type="range" bind:value={incrementIndex} min="0" max={increments.length - 1} />
			</label>
			<fieldset>
				<legend>Allow Berserk</legend>
				<label><input type="radio" name="berserk" value="f" bind:group={allowBerserk} /> No</label>
				<label><input type="radio" name="berserk" value="t" bind:group={allowBerserk} /> Yes</label>
			</fieldset>
			<div class="schedule-inputs">
				<label>Start date <input type="date" bind:value={date} class="form-field" required /></label
				>
				<label>Start time <input type="time" bind:value={time} class="form-field" required /></label
				>
			</div>

			{#if startTime && startTime <= new Date()}<p class="form-error">
					Choose a start time in the future.
				</p>{/if}
			{#if createError}<p class="form-error" role="alert">{createError}</p>{/if}
			<button
				type="submit"
				disabled={createLoading ||
					durationError ||
					nameError ||
					!startTime ||
					startTime <= new Date()}
				class="button button-primary">{createLoading ? 'Creating…' : 'Create tournament'}</button
			>
		</form>

		<section class="scheduled panel" aria-labelledby="scheduled-heading">
			<h2 id="scheduled-heading">Scheduled tournaments</h2>
			{#if actionError}<p class="form-error" role="alert">{actionError}</p>{/if}
			{#if data.tournaments === undefined}
				<p class="form-error">Could not load tournaments.</p>
			{:else if !data.tournaments?.length}
				<p class="empty-schedule">No scheduled tournaments</p>
			{:else}
				<ul>
					{#each data.tournaments as tournament (tournament.ID)}
						<li>
							<div>
								<a class="tournament-name" href={`/tournament/${tournament.ID}`}
									>{tournament.Name}</a
								>
								<p class="tournament-meta">
									{tournament.BaseTime / 60}+{tournament.Increment} · {tournament.Duration / 60} min
								</p>
								<time class="tournament-meta" datetime={tournament.StartTime}>
									{new Date(tournament.StartTime).toLocaleString('en-IN', {
										month: 'short',
										day: 'numeric',
										hour: '2-digit',
										minute: '2-digit'
									})}
								</time>
							</div>
							<div class="tournament-actions">
								<button
									class="button button-primary button-compact"
									disabled={startLoading || deleteLoading}
									onclick={() => handleStartTournament(tournament.ID)}>Start</button
								>
								<button
									class="button button-danger button-compact"
									disabled={deleteLoading || startLoading}
									onclick={() => handleDeleteTournament(tournament.ID)}>Delete</button
								>
							</div>
						</li>
					{/each}
				</ul>
			{/if}
		</section>
	</div>
</div>

<style>
	.admin-layout {
		display: grid;
		grid-template-columns: minmax(0, 1fr) minmax(0, 1.2fr);
		align-items: start;
		gap: 1rem;
		margin-top: 1.25rem;
	}
	.admin-form,
	.scheduled {
		display: grid;
		gap: 1rem;
		min-width: 0;
		padding: 1rem;
	}
	h2 {
		font-size: 1.125rem;
		font-weight: 650;
	}
	label {
		display: grid;
		gap: 0.375rem;
		color: var(--muted);
		font-size: 0.875rem;
	}
	label > span {
		display: flex;
		justify-content: space-between;
	}
	label strong {
		color: var(--accent);
	}
	fieldset {
		display: flex;
		gap: 1rem;
	}
	fieldset label {
		display: flex;
		align-items: center;
	}
	legend {
		margin-bottom: 0.375rem;
		color: var(--muted);
		font-size: 0.875rem;
	}
	.schedule-inputs {
		display: grid;
		grid-template-columns: repeat(auto-fit, minmax(min(8rem, 100%), 1fr));
		gap: 0.625rem;
	}
	.schedule-inputs input {
		min-width: 0;
		color-scheme: dark;
	}
	.empty-schedule,
	.tournament-meta {
		color: var(--muted);
		font-size: 0.875rem;
	}
	li {
		display: flex;
		align-items: center;
		justify-content: space-between;
		flex-wrap: wrap;
		gap: 0.625rem;
		padding-block: 0.75rem;
		border-top: 1px solid var(--line);
	}
	.tournament-name {
		font-weight: 650;
		overflow-wrap: anywhere;
	}
	.tournament-name:hover {
		color: var(--accent);
	}
	.tournament-actions {
		display: flex;
		gap: 0.375rem;
	}
	@media (max-width: 760px) {
		.admin-layout {
			grid-template-columns: minmax(0, 1fr);
		}
	}
</style>
