<script lang="ts">
	let {
		time,
		lowAbortTime,
		player,
		yourTurn
	}: {
		time: number;
		lowAbortTime: number;
		player: string;
		yourTurn: boolean;
	} = $props();
</script>

<div class="abort-timer" class:urgent={time <= lowAbortTime}>
	<span>{yourTurn ? 'Make your first move' : `Waiting for ${player}`}</span>
	<strong aria-label={`${Math.ceil(time / 1000)} seconds left`}>{Math.ceil(time / 1000)}s</strong>
</div>

<style>
	.abort-timer {
		display: flex;
		width: 100%;
		min-height: 2.5rem;
		align-items: center;
		justify-content: space-between;
		gap: 0.5rem;
		padding: 0.5rem 0.75rem;
		border: 1px solid var(--accent-line);
		border-radius: 8px;
		background: var(--accent-surface);
		color: #d9eed9;
		font-size: 0.875rem;
		overflow-wrap: anywhere;
	}
	.abort-timer.urgent {
		border-color: var(--danger-line);
		background: var(--danger-surface);
		color: var(--danger);
	}
	strong {
		flex: none;
		font-variant-numeric: tabular-nums;
	}
</style>
