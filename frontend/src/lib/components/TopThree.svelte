<script lang="ts">
	import goldImg from '$lib/assets/icons/gold.svg';
	import silverImg from '$lib/assets/icons/silver.svg';
	import bronzeImg from '$lib/assets/icons/bronze.svg';
	let { u1, u2, u3 } = $props();
</script>

<div class="mb-5 flex w-full justify-center">
	{#if u2}
		<div class="flex flex-col items-center justify-center">
			<img src={silverImg} alt="second place" class="h-[110px]" />{u2}
		</div>{/if}
	{#if u1}
		<div class="flex flex-col items-center justify-center">
			<img src={goldImg} alt="first place" class="h-[140px]" />{u1}
		</div>{/if}
	{#if u3}
		<div class="flex flex-col items-center justify-center">
			<img src={bronzeImg} alt="third place" class="h-[100px]" />{u3}
		</div>{/if}
</div>
