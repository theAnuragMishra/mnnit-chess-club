<script lang="ts">
	import { goto } from '$app/navigation';
	import { websocketStore } from '$lib/websocket.js';

	const { data } = $props();

	if (!data.user) goto('/');

	async function handleInitGame(timeControl: string, index: number) {
		activeState[index] = !activeState[index];
		websocketStore.sendMessage({ type: 'init_game', payload: { timeControl } });
	}

	const timeControls = [
		'1+0',
		'1+1',
		'2+1',
		'3+0',
		'3+2',
		'5+0',
		'5+3',
		'10+0',
		'10+5',
		'15+10',
		'30+0',
		'30+20'
	];
	let baseTime = $state(3);
	let increment = $state(0);
	let activeState = $state(Array(12).fill(false));
	let btError = $derived(
		baseTime === undefined || baseTime === null || baseTime <= 0 || baseTime > 180
	);
	let iError = $derived(
		increment === undefined ||
			increment === null ||
			increment < 0 ||
			increment > 180 ||
			!Number.isInteger(increment)
	);
</script>

<div class="box">
	{#each timeControls as timeControl, index}
		<button
			onclick={() => handleInitGame(timeControl, index)}
			class="btn relative cursor-pointer bg-gray-800 text-3xl"
			><span>{timeControl}</span>{#if activeState[index]}<span class="loading-bar">
					<span class="moving-indicator"></span>
				</span>{/if}</button
		>
	{/each}
	<div class="col-span-3 flex justify-around gap-2">
		<div>
			Minutes per side: <input
				class={`rounded-md bg-gray-800 px-2 py-1 ${btError ? 'border-2 border-red-500' : ''}`}
				type="number"
				bind:value={baseTime}
				min="0.5"
				max="180"
			/>
		</div>
		<div>
			Increment: <input
				class={`rounded-md bg-gray-800 px-2 py-1 ${iError ? 'border-2 border-red-500' : ''}`}
				type="number"
				bind:value={increment}
				min="0"
				max="180"
			/>
		</div>
	</div>
	<div class="col-span-3">
		<button
			onclick={() => {
				if (btError || iError) {
					return;
				}
				//console.log(baseTime, increment);
				websocketStore.sendMessage({
					type: 'create_challenge',
					payload: { timeControl: `${baseTime}+${increment}` }
				});
			}}
			class="cursor-pointer rounded-md bg-gray-800 px-3 py-2 text-xl">Create Challenge Link</button
		>
	</div>
</div>

<style>
	input::-webkit-outer-spin-button,
	input::-webkit-inner-spin-button {
		-webkit-appearance: none;
		margin: 0;
	}

	input[type='number'] {
		appearance: none;
		-moz-appearance: textfield;
	}
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
	}

	.loading-bar {
		position: absolute;
		top: 80px;
		left: 50%;
		transform: translateX(-50%);
		display: inline-block;
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
	}
</style>
