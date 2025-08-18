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
	<img
		src={images[piece]}
		alt="piece"
		class={`h-[15px] md:h-[35px] ${piece == 'p' && i < pieces.length - 1 && pieces[i + 1] == 'p' ? '-mr-[7px] md:-mr-[17px]' : ''} ${piece == 'r' && i < pieces.length - 1 && pieces[i + 1] == 'r' ? '-mr-[6px] md:-mr-[14px]' : ''} ${piece == 'n' && i < pieces.length - 1 && pieces[i + 1] == 'n' ? '-mr-[5px] md:-mr-[12px]' : ''} ${piece == 'b' && i < pieces.length - 1 && pieces[i + 1] == 'b' ? '-mr-[6px] md:-mr-[10px]' : ''} ${piece == 'q' && i < pieces.length - 1 && pieces[i + 1] == 'q' ? '-mr-[4px]' : ''}`}
	/>
{/each}
