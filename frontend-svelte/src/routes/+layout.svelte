<script>
	let { children } = $props();
	import '../app.css';
</script>

<div class="bg-gray-700 text-white min-h-screen flex flex-col items-center">
	{@render children()}
</div>
