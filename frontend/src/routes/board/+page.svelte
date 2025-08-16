<script lang="ts">
	import Chessboard from '$lib/components/Chessboard.svelte';
	import type { Api } from '@lichess-org/chessground/api';

	let g: Api | null = $state(null);
</script>

<div class="flex">
	<Chessboard boardConfig={{ viewOnly: true }} setGround={(x) => (g = x)} />
	<div class="flex flex-col">
		<button
			class="cursor-pointer bg-blue-500 p-4"
			onclick={() => {
				g?.set({ viewOnly: false });
			}}>toggle viewonly</button
		>
		<button
			class="cursor-pointer bg-blue-500 p-4"
			onclick={() => {
				g?.set({ drawable: { visible: !g.state.drawable.visible } });
			}}>toggle drawable.visible</button
		>
	</div>
</div>
