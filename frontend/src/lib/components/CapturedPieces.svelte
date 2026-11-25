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

{#each pieces as piece, i}
	<img src={images[piece]} alt={piece} class:stacked={i > 0 && pieces[i - 1] === piece} />
{/each}

<style>
	img {
		width: 1rem;
		height: 1rem;
		flex: none;
	}
	img.stacked {
		margin-left: -0.5rem;
	}
</style>
