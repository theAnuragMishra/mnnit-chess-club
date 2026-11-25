<script lang="ts">
	import { user } from '$lib/user.svelte';
	import BoardPreview from '$lib/components/BoardPreview.svelte';
	import Icon from '$lib/components/Icon.svelte';
	const formats = [
		{
			name: 'Bullet',
			time: '1 + 0',
			icon: 'bolt',
			mode: 'bullet'
		},
		{ name: 'Blitz', time: '5 + 0', icon: 'spark', mode: 'blitz' },
		{
			name: 'Rapid',
			time: '15 + 10',
			icon: 'clock',
			mode: 'rapid'
		}
	] as const;
</script>

<svelte:head
	><title>MNNIT Chess Club</title><meta
		name="description"
		content="Play chess online, challenge friends, and join arena tournaments at MNNIT Chess Club."
	/></svelte:head
>
<div class="page-shell home-page">
	<section class="hero" aria-labelledby="hero-title">
		<div class="hero-copy">
			<h1 id="hero-title">MNNIT <em>Chess Club</em></h1>
			<p class="hero-description">
				Play online, challenge friends, and join arena tournaments. Free and open to everyone.
			</p>
			<div class="hero-actions">
				<a class="button button-primary" href={user.id ? '/play' : '/login'}
					><Icon name="knight" size={21} />Play chess</a
				><a class="button button-secondary" href="/tournaments"
					><Icon name="trophy" size={19} />Tournaments</a
				>
			</div>
		</div>
		<div class="hero-board"><BoardPreview /></div>
	</section>
	<section class="play-section" aria-labelledby="play-heading">
		<div class="section-heading">
			<div>
				<h2 id="play-heading">Play options</h2>
			</div>
			<a href="/play" class="text-link">All time controls</a>
		</div>
		<div class="format-grid">
			{#each formats as format}<a class="format-card" href={`/play/?mode=${format.mode}`}
					><div class="format-top">
						<span class="format-icon"><Icon name={format.icon} size={24} /></span><span
							class="format-name">{format.name}</span
						><span class="card-chevron"><Icon name="chevron" size={17} /></span>
					</div>
					<div class="format-bottom">
						<strong>{format.time}</strong>
					</div></a
				>{/each}<a class="format-card friend-card" href="/play/#challenge"
				><div class="format-top">
					<span class="format-icon"><Icon name="people" size={24} /></span><span class="format-name"
						>Play a friend</span
					><span class="card-chevron"><Icon name="chevron" size={17} /></span>
				</div>
				<div class="friend-bottom">
					<span>Create a challenge link</span>
				</div></a
			>
		</div>
	</section>
</div>

<style>
	.home-page {
		padding-top: 0;
		padding-bottom: 50px;
	}
	.hero {
		display: grid;
		grid-template-columns: 1.15fr 1fr;
		align-items: center;
		gap: 65px;
		padding: clamp(2rem, 6dvh, 4rem) 0;
	}
	.hero-copy {
		padding-bottom: 10px;
	}
	h1 {
		margin: 24px 0;
		font-family: Georgia, 'Times New Roman', serif;
		font-size: clamp(3.5rem, 5.6vw, 5.25rem);
		font-weight: 400;
		line-height: 1.06;
		letter-spacing: -0.055em;
	}
	h1 em {
		color: var(--accent);
		font-weight: 400;
	}
	.hero-description {
		max-width: 420px;
		color: #aab5ac;
		font-size: 1rem;
		line-height: 1.85;
	}
	.hero-actions {
		display: flex;
		flex-wrap: wrap;
		gap: 12px;
		margin-top: 31px;
	}
	.hero-actions .button {
		padding-inline: 19px;
	}
	.hero-board {
		isolation: isolate;
		padding: 9px 18px 0 0;
	}
	.play-section {
		border-top: 1px solid var(--line);
		padding-top: 34px;
	}
	.section-heading {
		display: flex;
		align-items: end;
		justify-content: space-between;
		gap: 20px;
		margin-bottom: 23px;
	}
	h2 {
		margin-top: 9px;
		font-family: Georgia, 'Times New Roman', serif;
		font-size: 1.875rem;
		font-weight: 400;
		letter-spacing: -0.025em;
	}
	.section-heading > a {
		margin-bottom: 4px;
		white-space: nowrap;
		font-size: 0.875rem;
	}
	.format-grid {
		display: grid;
		grid-template-columns: repeat(4, 1fr);
		gap: 15px;
	}
	.format-card {
		padding: 21px;
		border: 1px solid #344036;
		border-radius: 12px;
		background: #1c231e;
		transition:
			transform 0.2s,
			background 0.2s,
			border-color 0.2s;
	}
	.format-card:hover {
		border-color: #78967b;
		background: #253129;
		transform: translateY(-4px);
	}
	.format-top {
		display: flex;
		align-items: center;
		gap: 11px;
	}
	.format-icon {
		color: #a9c8a6;
	}
	.format-name {
		color: #dee5d9;
		font-size: 0.875rem;
		font-weight: 550;
	}
	.card-chevron {
		margin-left: auto;
		color: #718575;
	}
	.format-bottom {
		display: flex;
		align-items: baseline;
		gap: 12px;
		margin-top: 24px;
	}
	.format-bottom strong {
		color: #edf1e5;
		font-size: 1.625rem;
		font-weight: 500;
		white-space: nowrap;
		letter-spacing: -0.025em;
	}
	.friend-card {
		border-color: #46573c;
		background: #252c20;
	}
	.friend-bottom {
		display: flex;
		flex-direction: column;
		gap: 4px;
		margin-top: 23px;
	}
	.friend-bottom > span {
		font-size: 0.75rem;
		color: #a3ae94;
	}
	@media (min-width: 1500px) {
		.hero {
			padding-block: 4rem;
		}
	}
	@media (max-width: 1100px) {
		.hero {
			gap: 25px;
		}
		.hero-actions {
			flex-direction: column;
			align-items: start;
		}
		.format-card {
			padding: 18px;
		}
		.format-bottom {
			flex-direction: column;
			gap: 1px;
			margin-top: 18px;
		}
	}
	@media (max-width: 800px) {
		.hero {
			grid-template-columns: 1fr 1fr;
			gap: 20px;
			padding: 2rem 0;
		}
		h1 {
			font-size: 3.25rem;
		}
		.hero-description {
			font-size: 0.875rem;
		}
		.format-grid {
			grid-template-columns: repeat(2, 1fr);
		}
		.format-bottom {
			flex-direction: row;
			gap: 16px;
		}
	}
	@media (max-width: 600px) {
		.hero {
			grid-template-columns: 1fr;
			padding: 1.5rem 0;
			gap: 38px;
		}
		h1 {
			font-size: clamp(3.4rem, 11vw, 4.5rem);
		}
		.hero-actions {
			flex-direction: row;
			gap: 10px;
		}
		.hero-actions .button {
			padding: 12px 14px;
			font-size: 0.8125rem;
		}
		.hero-board {
			padding: 6px 13px 7px;
		}
		.hero-description {
			font-size: 1rem;
		}
		.section-heading {
			align-items: start;
		}
		.section-heading h2 {
			max-width: 240px;
			font-size: 1.625rem;
		}
		.section-heading > a {
			font-size: 0.75rem;
			padding-top: 1px;
		}
		.format-grid {
			gap: 11px;
		}
		.format-card {
			padding: 16px;
		}
		.format-top {
			gap: 8px;
		}
		.format-name {
			font-size: 0.8125rem;
		}
		.format-bottom {
			align-items: start;
			flex-direction: column;
			gap: 2px;
		}
		.format-bottom strong {
			font-size: 1.75rem;
		}
		.friend-bottom {
			gap: 8px;
		}
	}
</style>
