<script>
	import { formatPostgresTimestamp, formatResultAndReason } from '$lib/utils';

	const { whiteUsername, blackUsername, result, baseTime, increment, createdAt, reason } = $props();
	// console.log(createdAt);
	const createdAtString = formatPostgresTimestamp(createdAt);
	const totalTime = baseTime / 60 + (increment * 2) / 3;
	const format =
		totalTime < 3 ? 'Bullet' : totalTime < 15 ? 'Blitz' : totalTime < 60 ? 'Rapid' : 'Classical';
</script>

<div
	class="flex h-30 w-full flex-col justify-around gap-1 rounded bg-[#1c1d1e] p-4 shadow-lg md:h-40"
>
	<div class="text-[16px] leading-tight">
		<div class="flex items-center">
			{baseTime / 60} + {increment} Rated {format}
		</div>
		<div class="flex items-center">{createdAtString}</div>
	</div>
	<div class="hidden flex-col md:flex">
		<div class="flex items-center gap-2">
			<div class="inline-block h-[16px] w-[16px] rounded-full border-2 bg-white"></div>
			{whiteUsername}
		</div>

		<div class="flex items-center gap-2">
			<div class="inline-block h-[16px] w-[16px] rounded-full border-2 bg-black"></div>
			{blackUsername}
		</div>
	</div>

	{#if result != '' && result != 'ongoing'}
		<div class="h-[1px] w-full bg-gray-500 opacity-50"></div>
		<div class=" flex items-center justify-center">
			{formatResultAndReason(result, reason)}
		</div>
	{/if}
</div>
