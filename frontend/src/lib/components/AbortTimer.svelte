<script lang="ts">
	const { time } = $props();
</script>

<span
	class={`text-md w-full rounded-md px-2 py-1 md:text-xl lg:text-2xl ${time > 10 ? 'bg-green-600 text-white' : 'bg-red-500 text-white'}`}
	>{time} second{time != 1 ? 's' : ''} to play the first move</span
>
