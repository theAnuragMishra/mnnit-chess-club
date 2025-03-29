<script lang="ts">
	let result = '1-0';
	let reason = 'Checkmate';
	let open = $state(true);
</script>

{#if open}
	<div class="fixed inset-0 z-50 flex items-center justify-center">
		<div
			class="relative w-80 rounded-lg border border-gray-300 bg-gray-800 p-6 text-center shadow-lg"
		>
			<h2 class="mb-2 text-2xl font-bold text-white">Game Over</h2>
			<p class="mb-4 text-4xl text-white">{result}</p>
			<p class="text-xl text-white">{reason}!</p>
			<p>
				{result === '1-0' ? 'White is victorious' : 'Black is victorious'}.
			</p>
			<button
				class="absolute top-2 right-2 text-gray-500 hover:text-gray-800"
				onclick={() => (open = false)}
			>
				✖
			</button>
		</div>
	</div>
{/if}
