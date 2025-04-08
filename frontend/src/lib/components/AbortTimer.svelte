<script lang="ts">
	const { time, tb } = $props();
</script>

<div
	class={`text-md text-center rounded-${tb}-md px-2 py-0.5 lg:text-xl ${time > 10 ? 'bg-green-600 text-white' : 'bg-red-500 text-white'}`}
>
	{time} second{time != 1 ? 's' : ''} to play the first move
</div>
