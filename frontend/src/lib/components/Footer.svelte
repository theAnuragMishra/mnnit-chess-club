<script>
	import Icon from './Icon.svelte';
</script>

<footer>
	<div class="footer-inner">
		<a href="/" class="footer-brand"
			><Icon name="knight" size={21} /><span>MNNIT Chess Club</span></a
		>
		<div>
			<a href="/dev">Server architecture</a><a
				href="https://github.com/theAnuragMishra/mnnit-chess-club"
				target="_blank"
				rel="noreferrer"
				aria-label="Source code on GitHub"><Icon name="code" size={18} /><span>Source code</span></a
			>
		</div>
	</div>
</footer>

<style>
	footer {
		border-top: 1px solid var(--line);
	}
	.footer-inner {
		display: flex;
		width: min(1184px, calc(100% - 80px));
		min-height: 3.75rem;
		align-items: center;
		justify-content: space-between;
		gap: 24px;
		margin: auto;
		color: #91a094;
		font-size: 0.75rem;
	}
	.footer-brand {
		display: flex;
		align-items: center;
		gap: 9px;
		color: #c5d1c6;
		font-weight: 550;
	}
	.footer-inner > div,
	.footer-inner > div a {
		display: flex;
		align-items: center;
		gap: 8px;
	}
	.footer-inner > div {
		gap: 23px;
	}
	a:hover {
		color: var(--accent);
	}
	@media (max-width: 900px) {
		.footer-inner {
			flex-wrap: wrap;
			justify-content: center;
			padding: 0.75rem 0;
			gap: 15px 28px;
		}
	}
	@media (max-width: 480px) {
		.footer-inner {
			width: calc(100% - 40px);
		}
		.footer-inner > div {
			gap: 16px;
		}
		.footer-inner > div a span {
			display: none;
		}
	}
</style>
