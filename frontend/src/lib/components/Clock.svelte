<script lang="ts">
	const { time, active, lowTime, compact = false } = $props();

	const formatTime = (time: number): string => {
		const minutes = Math.floor(time / 60000);
		const seconds = Math.floor((time % 60000) / 1000);
		const milliseconds = Math.floor((time % 1000) / 10);
		if (time > lowTime) return `${minutes}:${seconds.toString().padStart(2, '0')}`;
		return `${minutes}:${seconds.toString().padStart(2, '0')}.${milliseconds.toString().padStart(2, '0')}`;
	};
</script>

<span class="clock" class:compact class:active class:low={time < lowTime}>
	{formatTime(time)}
</span>

<style>
	.clock {
		display: inline-block;
		padding: 0.375rem 0.75rem;
		border: 1px solid var(--line);
		border-radius: 8px;
		background: var(--surface-raised);
		color: var(--muted);
		font-size: 1.5rem;
		font-weight: 650;
		font-variant-numeric: tabular-nums;
		line-height: 1;
		white-space: nowrap;
	}
	.clock.compact {
		padding: 0.25rem 0.5rem;
		font-size: 1.125rem;
	}
	.clock.active {
		border-color: var(--accent-line);
		background: var(--accent-surface);
		color: var(--accent);
	}
	.clock.low {
		border-color: var(--danger-line);
		background: var(--danger-surface);
		color: var(--danger);
	}
</style>
