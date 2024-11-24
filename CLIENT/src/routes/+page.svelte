<script>
	import NewTrainingSetForm from '$lib/components/NewTrainingSetForm.svelte';
	import TrainingSessions from '$lib/components/TrainingSessions.svelte';
	import TrainingSets from '$lib/components/TrainingSets.svelte';
	import Welcome from '$lib/components/Welcome.svelte';
</script>

<header>
	<Welcome />
</header>
<main>
	<NewTrainingSetForm />
	<TrainingSets />
	<TrainingSessions />
</main>
