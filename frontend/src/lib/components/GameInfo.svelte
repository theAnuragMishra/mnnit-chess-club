<script>
	import { formatPostgresTimestamp, formatResultAndMethod, getTimeControl } from '$lib/utils';

	const { result, baseTime, increment, createdAt, method } = $props();
</script>

<div class="game-info panel">
	<p>
		<strong>{baseTime / 60}+{increment}</strong>
		<span>{getTimeControl(baseTime, increment)} · Rated</span>
	</p>
	{#if result !== 0}<p>{formatResultAndMethod(result, method)}</p>{/if}
	<time datetime={createdAt.toISOString()}>{formatPostgresTimestamp(createdAt)}</time>
</div>

<style>
	.game-info {
		display: grid;
		gap: 0.25rem;
		padding: 0.625rem 0.75rem;
		border-radius: 8px;
		font-size: 0.875rem;
	}
	p:first-child {
		display: flex;
		flex-wrap: wrap;
		align-items: baseline;
		gap: 0.5rem;
	}
	span,
	time {
		color: var(--muted);
	}
	time {
		font-size: 0.75rem;
	}
</style>
