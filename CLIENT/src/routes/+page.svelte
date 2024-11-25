<script>
	import { browser } from '$app/environment';
	import { goto } from '$app/navigation';
	import NewTrainingSetForm from '$lib/components/NewTrainingSetForm.svelte';
	import SubjectStats from '$lib/components/SubjectStats.svelte';
	import TrainingSessions from '$lib/components/TrainingSessions.svelte';
	import TrainingSets from '$lib/components/TrainingSets.svelte';
	import Welcome from '$lib/components/Welcome.svelte';

	if (browser) {
		const trainingSessionId = localStorage.getItem('trainingSessionId');
		if (trainingSessionId) {
			goto(`/sessions/${trainingSessionId}`);
		}
	}
</script>

<div class="container">
	<header>
		<Welcome />
	</header>
	<main>
		<NewTrainingSetForm />
		<TrainingSets />
		<TrainingSessions />
		<SubjectStats />
	</main>
</div>

<style>
	.container {
		padding: 3rem;
	}
	main {
		margin-top: 2rem;
	}
</style>
