<script lang="ts">
	import { getBaseURL } from '../../utils/urlUtils';

	let username = $state('');
	let password = $state('');

	const handleSubmit = async (e: SubmitEvent) => {
		e.preventDefault();
		const res = await fetch(`${getBaseURL()}/login`, {
			method: 'POST',
			body: JSON.stringify({ username, password }),
			headers: {
				'Content-Type': 'application/json'
			},
			credentials: 'include'
		});
	};
</script>

<div class="flex flex-col items-center justify-center h-full mt-30">
	<h2 class="text-2xl font-bold mb-4">Login</h2>
	<form onsubmit={handleSubmit} class="bg-gray-800 p-6 rounded shadow-md w-full max-w-sm">
		<div class="mb-4">
			<label class="block font-medium mb-2"
				>Username<input
					class="w-full px-3 py-2 border rounded"
					type="text"
					name="username"
					bind:value={username}
				/></label
			>
		</div>
		<div class="mb-4">
			<label class="block font-medium mb-2"
				>Password<input
					class="w-full px-3 py-2 border rounded"
					type="password"
					name="password"
					bind:value={password}
				/></label
			>
		</div>

		<button
			type="submit"
			class="w-full bg-blue-500 text-white py-2 rounded hover:bg-blue-600 cursor-pointer"
			>Login</button
		>
	</form>
</div>
