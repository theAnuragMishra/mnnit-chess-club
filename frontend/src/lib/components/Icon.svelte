<script lang="ts">
	export type IconName =
		| 'knight'
		| 'bolt'
		| 'clock'
		| 'trophy'
		| 'people'
		| 'grid'
		| 'chevron'
		| 'menu'
		| 'close'
		| 'logout'
		| 'globe'
		| 'code'
		| 'check'
		| 'link'
		| 'spark'
		| 'shield';
	let { name, size = 22, ...rest }: { name: IconName; size?: number; class?: string } = $props();
</script>

<svg
	width={size}
	height={size}
	viewBox="0 0 24 24"
	fill="none"
	stroke="currentColor"
	stroke-width="1.6"
	stroke-linecap="round"
	stroke-linejoin="round"
	aria-hidden="true"
	{...rest}
>
	{#if name === 'knight'}<path
			d="M5 21h15M6 18h12c0-5-1-7-4-9l-3 3-6-1 2-5 4-2V2l5 3c4 3 4 8 3 13M10 7h.01M7 18v3m11-3v3"
		/>
	{:else if name === 'bolt'}<path d="m13 2-9 12h7l-1 8L20 9h-7l1-7Z" />
	{:else if name === 'clock'}<circle cx="12" cy="12" r="9" /><path d="M12 7v5l3 2" />
	{:else if name === 'trophy'}<path
			d="M8 3h8v6a4 4 0 0 1-8 0V3ZM8 5H4v2a4 4 0 0 0 4 4m8-6h4v2a4 4 0 0 1-4 4m-4 2v5m-4 3h8m-7 0v-3h6v3"
		/>
	{:else if name === 'people'}<path
			d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2m20 0v-2a4 4 0 0 0-3-3.87M16 3a4 4 0 0 1 0 8"
		/><circle cx="9" cy="7" r="4" />
	{:else if name === 'grid'}<rect x="3" y="3" width="7" height="7" rx="1" /><rect
			x="14"
			y="3"
			width="7"
			height="7"
			rx="1"
		/><rect x="3" y="14" width="7" height="7" rx="1" /><rect
			x="14"
			y="14"
			width="7"
			height="7"
			rx="1"
		/>
	{:else if name === 'chevron'}<path d="m9 5 7 7-7 7" />
	{:else if name === 'menu'}<path d="M4 6h16M4 12h16M4 18h16" />
	{:else if name === 'close'}<path d="m6 6 12 12M6 18 18 6" />
	{:else if name === 'logout'}<path d="M9 21H4V3h5m5 5 4 4-4 4m-6-4h13" />
	{:else if name === 'globe'}<circle cx="12" cy="12" r="9" /><ellipse
			cx="12"
			cy="12"
			rx="4"
			ry="9"
		/><path d="M3 12h18" />
	{:else if name === 'code'}<path d="m7 7-5 5 5 5m10-10 5 5-5 5m-4-14-2 18" />
	{:else if name === 'check'}<path d="m5 12 4 4L19 6" />
	{:else if name === 'link'}<path
			d="m10 13 4-4m-6 6-2 2a4 4 0 0 1-6-6l4-4a4 4 0 0 1 6 0m4 2 2-2a4 4 0 0 1 6 6l-4 4a4 4 0 0 1-6 0"
			transform="translate(1 0) scale(.92)"
		/>
	{:else if name === 'spark'}<path
			d="m12 3 2.5 6.5L21 12l-6.5 2.5L12 21l-2.5-6.5L3 12l6.5-2.5L12 3Z"
		/>
	{:else if name === 'shield'}<path d="m12 3 8 3v6c0 5-8 9-8 9s-8-4-8-9V6l8-3Z" /><path
			d="m8 12 3 3 5-6"
		/>{/if}
</svg>
