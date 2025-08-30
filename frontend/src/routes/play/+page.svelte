<script lang="ts">
	import { goto } from '$app/navigation';
	import { websocketStore } from '$lib/websocket';

	const { data } = $props();

	async function handleInitGame(index: number) {
		if (activeIndex === index) activeIndex = -1;
		else activeIndex = index;
		websocketStore.sendMessage({ type: 'init_game', payload: index });
	}

	const timeControls = [
		{ baseTime: 60, increment: 0 },
		{ baseTime: 60, increment: 1 },
		{ baseTime: 120, increment: 1 },
		{ baseTime: 180, increment: 0 },
		{ baseTime: 180, increment: 2 },
		{ baseTime: 300, increment: 0 },
		{ baseTime: 300, increment: 3 },
		{ baseTime: 600, increment: 0 },
		{ baseTime: 600, increment: 5 },
		{ baseTime: 900, increment: 10 },
		{ baseTime: 1800, increment: 0 },
		{ baseTime: 1800, increment: 20 }
	];

	const baseTimes = [
		0.25, 0.5, 0.75, 1, 1.5, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 25,
		30, 35, 40, 45, 60, 75, 90, 105, 120, 135, 150, 165, 180
	];

	const increments = [
		0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 25, 30, 35, 40, 45,
		60, 90, 120, 150, 180
	];
	//console.log(baseTimes.length, increments.length);
	let activeIndex = $state(-1);
	let baseIndex = $state(6);
	let incrementIndex = $state(0);
</script>

<svelte:head>
	<title>Play</title>
</svelte:head>
<div class="box">
	{#each timeControls as timeControl, index}
		<button onclick={() => handleInitGame(index)} class="btn relative cursor-pointer bg-gray-800"
			><span>{timeControl.baseTime / 60}+{timeControl.increment}</span
			>{#if activeIndex === index}<span class="loading-bar">
					<span class="moving-indicator"></span>
				</span>{/if}</button
		>
	{/each}
	<div class="col-span-3 flex flex-col items-center justify-center gap-2">
		<div class="flex flex-col items-center justify-center">
			<span>Minutes per side: {baseTimes[baseIndex]}</span>
			<input type="range" bind:value={baseIndex} min="0" max="37" />
		</div>
		<div class="flex flex-col items-center justify-center">
			<span>Increment: {increments[incrementIndex]}</span>
			<input type="range" bind:value={incrementIndex} min="0" max="30" />
		</div>
	</div>
	<div class="col-span-3">
		<button
			onclick={() => {
				//console.log(baseTime, increment);
				websocketStore.sendMessage({
					type: 'create_challenge',
					payload: { baseTime: baseTimes[baseIndex] * 60, increment: increments[incrementIndex] }
				});
			}}
			class="cursor-pointer rounded-md bg-gray-800 px-3 py-2 text-xl">Create Challenge Link</button
		>
	</div>
</div>

<style>
	.box {
		display: grid;
		height: 300px;
		width: 350px;
		flex: 1;
		grid-template-columns: 1fr 1fr 1fr;
		place-content: center;
		place-items: center;
		place-self: center;
		row-gap: 10px;
	}
	.btn {
		width: 100px;
		height: 70px;
		font-size: 24px;
	}

	.loading-bar {
		position: absolute;
		left: 50%;
		transform: translateX(-50%);
		display: inline-block;
		top: 55px;
		width: 60%;
		height: 4px;
		background-color: #e0e0e0;
		overflow: hidden;
		border-radius: 2px;
	}

	.moving-indicator {
		display: inline-block;
		width: 30%;
		height: 100%;
		background-color: #3b82f6;
		position: absolute;
		animation: move 1.2s ease-in-out infinite alternate;
		border-radius: 2px;
	}

	@keyframes move {
		0% {
			left: 0%;
		}
		100% {
			left: 70%;
		}
	}

	@media (width>=450px) {
		.box {
			height: 400px;
			width: 400px;
		}
		.btn {
			width: 130px;
			height: 90px;
			font-size: 30px;
		}
		.loading-bar {
			top: 70px;
		}
	}

	@media (width>=600px) {
		.box {
			height: 500px;
			width: 500px;
		}
		.btn {
			width: 150px;
			height: 100px;
		}
		.loading-bar {
			top: 80px;
		}
	}

	input[type='range'] {
		appearance: none;
		-webkit-appearance: none;
		width: 300px;
		height: 8px;
		border-radius: 5px;
		background: #ddd;
		outline: none;
		transition: background 0.3s;
		margin: 10px 0;
	}

	input[type='range']:hover {
		background: #ccc;
	}

	/* Thumb styling */
	input[type='range']::-webkit-slider-thumb {
		-webkit-appearance: none;
		height: 20px;
		width: 20px;
		background: #367995;
		border-radius: 50%;
		cursor: pointer;
		border: none;
		transition: background 0.3s;
		margin-top: -6px; /* Align thumb vertically */
	}

	input[type='range']::-moz-range-thumb {
		height: 20px;
		width: 20px;
		background: #367995;
		border: none;
		border-radius: 50%;
		cursor: pointer;
	}

	/* Track styling for Firefox */
	input[type='range']::-moz-range-track {
		height: 8px;
		background: #ddd;
		border-radius: 5px;
	}

	/* IE and Edge */
	input[type='range']::-ms-thumb {
		height: 20px;
		width: 20px;
		background: #367995;
		border: none;
		border-radius: 50%;
		cursor: pointer;
	}

	input[type='range']::-ms-track {
		height: 8px;
		background: transparent;
		border-color: transparent;
		color: transparent;
	}
</style>
