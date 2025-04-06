<script lang="ts">
	const { initialTime, active } = $props();

	let time = $state(initialTime * 1000);

	$effect(() => {
		time = initialTime * 1000;
	});

	let animationFrame: number | null;
	let startTime: DOMHighResTimeStamp | null;

	$effect(() => {
		if (active) {
			startTime = performance.now();
			const tick = (currentTime: number) => {
				if (!startTime) return;
				const elapsed = currentTime - startTime;
				const newTime = initialTime * 1000 - elapsed;

				if (newTime <= 0) {
					time = 0;
					return;
				}

				time = newTime;
				animationFrame = requestAnimationFrame(tick);
			};

			animationFrame = requestAnimationFrame(tick);
		} else {
			if (animationFrame !== null) {
				cancelAnimationFrame(animationFrame);
				animationFrame = null;
				startTime = null;
			}
		}
		return () => {
			if (animationFrame !== null) {
				cancelAnimationFrame(animationFrame);
				animationFrame = null;
				startTime = null;
			}
		};
	});

	const formatTime = (time: number): string => {
		const minutes = Math.floor(time / 60000);
		const seconds = Math.floor((time % 60000) / 1000);
		const milliseconds = Math.floor((time % 1000) / 10);
		if (time > 10000) return `${minutes}:${seconds.toString().padStart(2, '0')}`;
		return `${minutes}:${seconds.toString().padStart(2, '0')}.${milliseconds.toString().padStart(2, '0')}`;
	};
</script>

<span
	class={`my-2 rounded-md px-2 py-1 text-4xl ${
		active ? 'font-bold ' : 'bg-black text-gray-400'
	} ${active && time > 10000 ? 'bg-white text-black' : ''} ${time < 10000 ? 'bg-red-500 text-white' : ''}`}
>
	{formatTime(time)}
</span>
