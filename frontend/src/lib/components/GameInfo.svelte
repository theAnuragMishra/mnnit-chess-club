<script>
	import { formatPostgresTimestamp } from '$lib/utils';

	const { whiteUsername, blackUsername, result, baseTime, increment, createdAt, reason } = $props();
	// console.log(createdAt);
	const createdAtString = formatPostgresTimestamp(createdAt);
</script>

<div class="flex h-64 w-1/4 flex-col gap-3 rounded border p-4 shadow-lg">
	<div class="flex items-center justify-around">
		<div class="flex items-center gap-2">
			<div class="inline-block h-[24px] w-[24px] rounded-full border-2 border-black bg-white"></div>
			{whiteUsername}
		</div>
		<div>Vs</div>
		<div class="flex items-center gap-2">
			<div class="inline-block h-[24px] w-[24px] rounded-full border-2 border-white bg-black"></div>
			{blackUsername}
		</div>
	</div>
	<div class="flex items-center justify-around">{createdAtString}</div>
	<div class="flex items-center justify-around">Time Control: {baseTime / 60} + {increment}</div>
	<div class="flex items-center justify-center text-[20px]">
		{reason} | {result === '1-0'
			? 'White is victorious'
			: result === '0-1'
				? 'Black is victorious'
				: 'Game Drawn'}
	</div>
</div>
