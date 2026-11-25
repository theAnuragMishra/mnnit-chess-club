<script lang="ts">
	import { page } from '$app/state';
	import { afterNavigate } from '$app/navigation';
	import { logout, user } from '$lib/user.svelte';
	import Icon from './Icon.svelte';
	import logo from '$lib/assets/mcc-logo.webp';
	let expanded = $state(false);
	const links = [
		{ href: '/play', label: 'Play', icon: 'grid' },
		{ href: '/tournaments', label: 'Tournaments', icon: 'trophy' },
		{ href: '/leaderboard', label: 'Leaderboard', icon: 'people' },
		{ href: '/achievements', label: 'About', icon: 'spark' }
	] as const;
	afterNavigate(() => (expanded = false));
</script>

<svelte:window
	onkeydown={(event) => {
		if (event.key === 'Escape') expanded = false;
	}}
/>
<header class="site-header">
	<nav aria-label="Main navigation" class="nav-inner">
		<a href="/" class="brand" aria-label="MNNIT Chess Club home"
			><span class="brand-mark"
				><img class="club-logo" src={logo} alt="" width="36" height="30" /></span
			><span class="brand-text">MNNIT<span>CHESS CLUB</span></span></a
		>
		<div class="desktop-links">
			{#each links as link}<a
					href={link.href}
					class:active={page.url.pathname.startsWith(link.href)}
					aria-current={page.url.pathname.startsWith(link.href) ? 'page' : undefined}
					>{link.label}</a
				>{/each}
		</div>
		<div class="nav-actions">
			{#if user.id}<a
					class="profile-link"
					href={`/member/${user.username}`}
					aria-label={`${user.username}'s profile`}
					><span class="avatar">{user.username?.slice(0, 1).toUpperCase()}</span><span
						class="username">{user.username}</span
					></a
				><button class="icon-button signout" onclick={logout} aria-label="Sign out"
					><Icon name="logout" size={19} /></button
				>
			{:else}<a class="join-link" href="/login">Log in</a>{/if}
			<button
				class="icon-button menu-button"
				onclick={() => (expanded = !expanded)}
				aria-label={expanded ? 'Close navigation' : 'Open navigation'}
				aria-expanded={expanded}
				aria-controls="mobile-navigation"><Icon name={expanded ? 'close' : 'menu'} /></button
			>
		</div>
	</nav>
	{#if expanded}<nav id="mobile-navigation" class="mobile-links" aria-label="Mobile navigation">
			{#each links as link}<a
					href={link.href}
					class:active={page.url.pathname.startsWith(link.href)}
					aria-current={page.url.pathname.startsWith(link.href) ? 'page' : undefined}
					><Icon name={link.icon} size={19} />{link.label}</a
				>{/each}
		</nav>{/if}
</header>

<style>
	.club-logo {
		width: 36px;
		height: 30px;
		object-fit: contain;
		filter: invert(89%) sepia(14%) saturate(649%) hue-rotate(74deg) brightness(98%) contrast(89%);
	}
	.site-header {
		position: relative;
		z-index: 50;
		border-bottom: 1px solid var(--line);
		background: #141916;
	}
	.nav-inner {
		display: flex;
		width: min(1320px, calc(100% - 80px));
		min-height: calc(var(--site-header-height) - 1px);
		align-items: center;
		justify-content: space-between;
		gap: 28px;
		margin: auto;
	}
	.brand {
		display: inline-flex;
		flex-shrink: 0;
		align-items: center;
		gap: 10px;
	}
	.brand-mark {
		display: grid;
		height: 46px;
		width: 46px;
		place-items: center;
		border: 1px solid #485e4b;
		border-radius: 12px;
		color: var(--accent);
	}
	.brand-text {
		font-size: 1.125rem;
		font-weight: 750;
		letter-spacing: 0.085em;
		line-height: 1.2;
	}
	.brand-text > span {
		display: block;
		margin-top: 4px;
		color: var(--muted);
		font-size: 0.625rem;
		letter-spacing: 0.19em;
		font-weight: 500;
	}
	.desktop-links {
		display: flex;
		align-items: center;
		gap: 29px;
	}
	.desktop-links a {
		position: relative;
		display: grid;
		height: calc(var(--site-header-height) - 1px);
		align-items: center;
		font-size: 0.875rem;
		color: #b8c1b9;
		transition: color 0.15s;
	}
	.desktop-links a:hover,
	.desktop-links a.active {
		color: var(--accent);
	}
	.desktop-links a.active::after {
		position: absolute;
		bottom: -1px;
		left: 0;
		width: 100%;
		height: 2px;
		content: '';
		background: var(--accent);
	}
	.nav-actions {
		display: flex;
		align-items: center;
		gap: 20px;
	}
	.join-link {
		border: 1px solid #536b55;
		border-radius: 7px;
		padding: 10px 17px;
		font-size: 0.8125rem;
		font-weight: 600;
		color: var(--accent);
		transition: background 0.15s;
	}
	.join-link:hover {
		background: #26372a;
	}
	.profile-link {
		display: flex;
		align-items: center;
		gap: 9px;
		font-size: 0.875rem;
	}
	.avatar {
		display: grid;
		width: 34px;
		height: 34px;
		place-items: center;
		border-radius: 50%;
		background: #324c38;
		color: var(--accent);
	}
	.icon-button {
		display: grid;
		min-width: 40px;
		min-height: 40px;
		place-items: center;
		border-radius: 6px;
		color: var(--muted);
	}
	.icon-button:hover {
		background: var(--surface-raised);
		color: var(--accent);
	}
	.menu-button {
		display: none;
	}
	.mobile-links {
		display: none;
		position: absolute;
		inset: 100% 0 auto;
		max-height: calc(100dvh - var(--site-header-height));
		overflow-y: auto;
		background: #141916;
		border-bottom: 1px solid var(--line);
		flex-direction: column;
		gap: 6px;
		padding: 12px 20px 20px;
		border-top: 1px solid var(--line);
	}
	.mobile-links a {
		display: flex;
		align-items: center;
		gap: 12px;
		padding: 14px;
		border-radius: 8px;
		color: #b8c1b9;
	}
	.mobile-links a:hover,
	.mobile-links a.active {
		background: #26372a;
		color: var(--accent);
	}
	@media (max-width: 1050px) {
		.desktop-links {
			gap: 20px;
		}
		.nav-inner {
			gap: 20px;
			width: calc(100% - 48px);
		}
		.nav-actions {
			gap: 12px;
		}
	}
	@media (max-width: 860px) {
		.mobile-links {
			display: flex;
		}
		.desktop-links {
			display: none;
		}
		.menu-button {
			display: grid;
		}
	}
	@media (max-width: 480px) {
		.nav-inner {
			width: calc(100% - 32px);
			gap: 10px;
		}
		.brand-mark {
			width: 40px;
			height: 42px;
		}
		.brand-text {
			font-size: 1rem;
		}
		.username {
			display: none;
		}
		.join-link {
			font-size: 0.75rem;
			padding: 9px 12px;
		}
		.nav-actions {
			gap: 6px;
		}
	}
</style>
