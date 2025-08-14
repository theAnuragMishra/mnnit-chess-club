<script lang="ts">
	import { dateTimeToDate, getBaseURL } from '$lib/utils';

	const { data } = $props();
	//console.log(data);
	let name = $state('');
	let duration: number | undefined = $state();
	let date: string | undefined = $state();
	let time: string | undefined = $state();
	let startTime = $derived(dateTimeToDate(date, time));

	let startLoading = $state(false);
	let createLoading = $state(false);

	let createError = $state('');
	let startError = $state(false);
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
				startTime: startTime.toISOString()
			}),
			credentials: 'include'
		});

		const response = await res.json();

		if (!res.ok) {
			createError = response.error;
			createLoading = false;
			return;
		}

		window.location.reload();
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

<div class="w-[100%] p-2 lg:p-5">
	<h1 class="mb-5 text-xl">You're an mcc admin. Use this page to create or start a tournament!</h1>

	<div class="m-auto flex w-full flex-col gap-2 border-2 p-[2px] lg:w-1/2 lg:p-2 lg:px-10">
		<h2>Create tournament</h2>
		<label
			>Name:<input
				type="text"
				bind:value={name}
				class={`bg-gray-200 text-black ${nameError ? 'border-2 border-red-500' : ''}`}
			/></label
		>
		<label
			>Duration: <input
				type="text"
				bind:value={duration}
				class={`bg-gray-200 text-black ${durationError ? 'border-2 border-red-500' : ''}`}
			/></label
		>

		<div class="flex flex-col md:flex-row">
			<span>Minutes per side: {baseTimes[baseIndex]}</span>
			<input class="mx-[5px]" type="range" bind:value={baseIndex} min="0" max="37" />
		</div>
		<div class="flex flex-col md:flex-row">
			<span>Increment: {increments[incrementIndex]}</span>
			<input class="mx-[5px]" type="range" bind:value={incrementIndex} min="0" max="30" />
		</div>

		<label>Start Date: <input type="date" bind:value={date} class="bg-gray-200 text-black" /></label
		>
		<label>Start Time: <input type="time" bind:value={time} class="bg-gray-200 text-black" /></label
		>

		<span
			><button
				disabled={createLoading ||
					durationError ||
					nameError ||
					!startTime ||
					startTime <= new Date()}
				class="w-fit cursor-pointer rounded-md bg-green-400 px-4 py-2 text-black disabled:cursor-not-allowed"
				onclick={handleCreateTournament}>Create</button
			></span
		>
		<p class="text-red-500">
			{nameError
				? 'Name length must be in the range [1,90]'
				: durationError
					? 'Duration must be a number in the range [1,24] (hours)'
					: ''}
		</p>
		<p class="text-red-500">{createError}</p>
	</div>

	<div>
		{#if data.tournaments === undefined}<p class="text-xl">
				Error fetching tournaments. Refresh or try again after some time.
			</p>
		{:else if data.tournaments === null}
			<p class="text-xl">No upcoming tournaments</p>
		{:else}
			<p class="text-xl">Upcoming tournaments:</p>
			<p class="mb-[5px] flex justify-around text-xl">
				<span>Name</span><span>Time Control</span><span>Start Time</span><span
					>Duration (hours)</span
				><span>Start</span>
			</p>
			<hr class="mb-4" />
			{#each data.tournaments as tournament}
				<p class="flex justify-around">
					<span><a href={`/tournament/${tournament.ID}`}>{tournament.Name}</a></span><span
						>{tournament.BaseTime / 60}+{tournament.Increment}</span
					><span
						>{new Date(tournament.StartTime).toLocaleString('en-IN', {
							year: 'numeric',
							month: 'short',
							day: '2-digit',
							hour: '2-digit',
							minute: '2-digit',
							hour12: false
						})}</span
					><span>{tournament.Duration / 3600}</span>
					<span
						><button
							class="cursor-pointer rounded-md bg-green-400 px-4 py-2 text-black disabled:cursor-not-allowed"
							disabled={startLoading}
							onclick={async () => {
								startLoading = true;
								const res = await fetch(`${getBaseURL()}/admin/start-tournament`, {
									method: 'POST',
									headers: { 'Content-Type': 'application/json' },
									body: JSON.stringify({
										tournamentID: tournament.ID
									}),
									credentials: 'include'
								});

								if (!res.ok) {
									startError = true;
									startLoading = false;
									return;
								}
								window.location.reload();
							}}>Start</button
						></span
					>
				</p>
			{/each}
		{/if}
	</div>
</div>

<style>
	p > span {
		width: 20%;
		text-align: center;
	}
	label {
		display: flex;
		gap: 5px;
	}
</style>
