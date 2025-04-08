<script lang="ts">
	const { time, active } = $props();

	const formatTime = (time: number): string => {
		const minutes = Math.floor(time / 60000);
		const seconds = Math.floor((time % 60000) / 1000);
		const milliseconds = Math.floor((time % 1000) / 10);
		if (time > 10000) return `${minutes}:${seconds.toString().padStart(2, '0')}`;
		return `${minutes}:${seconds.toString().padStart(2, '0')}.${milliseconds.toString().padStart(2, '0')}`;
	};
</script>

<span
	class={`my-2 w-fit rounded-md px-2 py-1 text-2xl md:text-3xl lg:text-4xl ${
		active ? 'font-bold ' : 'bg-black text-gray-400'
	} ${active && time > 10000 ? 'bg-white text-black' : ''} ${time < 10000 ? 'bg-red-500 text-white' : ''}`}
>
	{formatTime(time)}
</span>
