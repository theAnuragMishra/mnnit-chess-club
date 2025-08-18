<script lang="ts">
	import p from '$lib/assets/pieces/captured/P.svg';
	import q from '$lib/assets/pieces/captured/Q.svg';
	import n from '$lib/assets/pieces/captured/N.svg';
	import r from '$lib/assets/pieces/captured/R.svg';
	import b from '$lib/assets/pieces/captured/B.svg';
	const images = {
		p,
		q,
		n,
		r,
		b
	};
	let { pieces }: { pieces: Array<'p' | 'q' | 'n' | 'r' | 'b'> } = $props();
</script>

{#each pieces as piece}
	<div
		style={`background-image: url("${images[piece]}");`}
		class="aspect-square h-[15px] bg-contain md:h-[35px]"
	></div>
{/each}
