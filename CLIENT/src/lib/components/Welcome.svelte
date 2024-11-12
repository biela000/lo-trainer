<script lang="ts">
	import { API_URL } from '$lib/api/constants';

	let streak = $state(0);

	$effect(() => {
		const getStreak = async () => {
			const response = await fetch(`${API_URL}/training_sessions/streak`);
			const { streak } = (await response.json()) as { streak: number };

			return streak;
		};

		getStreak().then((newStreak) => {
			streak = newStreak;
		});
	});
</script>

<h1>Welcome to LO Trainer</h1>
<h2>Your current streak is {streak} 🔥</h2>
