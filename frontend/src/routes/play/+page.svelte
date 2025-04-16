<script lang="ts">
	import { goto } from '$app/navigation';
	import { websocketStore } from '$lib/websocket.js';

	const { data } = $props();

	if (!data.user) goto('/');

	async function handleInitGame(timeControl: string) {
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
	let baseTime = $state(1);
	let increment = $state(0);
</script>

<div class="box">
	{#each timeControls as timeControl}
		<button
			onclick={() => handleInitGame(timeControl)}
			class="btn cursor-pointer bg-gray-800 text-3xl">{timeControl}</button
		>
	{/each}
	<div class="col-span-3 flex justify-around gap-2">
		<div>
			Minutes per side: <input
				class="rounded-md bg-gray-800 px-2 py-1"
				type="number"
				bind:value={baseTime}
				min="0.5"
				max="180"
			/>
		</div>
		<div>
			Increment: <input
				class="rounded-md bg-gray-800 px-2 py-1"
				type="number"
				bind:value={increment}
				min="0"
				max="180"
			/>
		</div>
	</div>
	<div class="col-span-3">
		<button
			onclick={() =>
				websocketStore.sendMessage({
					type: 'create_challenge',
					payload: { timeControl: `${baseTime}+${increment}` }
				})}
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
