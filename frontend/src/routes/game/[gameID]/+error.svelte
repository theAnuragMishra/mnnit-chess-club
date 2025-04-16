<script lang="ts">
	import { page } from '$app/state';
</script>

<div class="abc flex flex-col items-center justify-center">
	<h1 class="text-4xl">{page.error?.message}</h1>
	<p class="sad">:(</p>
</div>

<style>
	.sad {
		font-size: 100px;
	}
	.abc {
		height: 80vh;
	}
</style>
