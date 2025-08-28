<script lang="ts">
	const { time, tb, lowAbortTime } = $props();
</script>

<div
	class={`text-md text-center md:w-[300px] rounded-${tb}-md px-2 py-0.5 ${time > lowAbortTime ? 'bg-green-600 text-white' : 'bg-red-500 text-white'}`}
>
	{Math.floor(time / 1000)} seconds to play the first move
</div>
