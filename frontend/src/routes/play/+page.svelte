<script lang="ts">
	import { goto } from '$app/navigation';
	import { getBaseURL } from '$lib/utils';
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
</script>

<div class="box">
	{#each timeControls as timeControl}
		<button
			onclick={() => handleInitGame(timeControl)}
			class="btn cursor-pointer bg-gray-800 text-3xl"
		>
			{timeControl}
		</button>
	{/each}
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
