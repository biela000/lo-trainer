<script lang="ts">
	import { goto } from '$app/navigation';
	import { API_URL } from '$lib/api/constants';
	import type { TrainingSession } from '$lib/types/training-session';
	import TrainingSets from './TrainingSets.svelte';

	let trainingSessions: TrainingSession[] = $state([]);

	$effect(() => {
		const getTrainingSessions = async () => {
			const response = await fetch(`${API_URL}/training_sessions`);
			trainingSessions = (await response.json()) as TrainingSession[];
		};

		getTrainingSessions();
	});

	const goToSession = (id: number) => {
		goto(`/sessions/${id}`);
	};
</script>

<table>
	<thead>
		<tr>
			<th>Id</th>
			<th>Training set Id</th>
			<th>P1 Score</th>
			<th>P2 Score</th>
			<th>P3 Score</th>
			<th>P4 Score</th>
			<th>P5 Score</th>
			<th>Full score</th>
			<th>P1 Time</th>
			<th>P2 Time</th>
			<th>P3 Time</th>
			<th>P4 Time</th>
			<th>P5 Time</th>
			<th>Full time</th>
			<th>Finished</th>
			<th></th>
		</tr>
	</thead>
	<tbody>
		{#each trainingSessions as trainingSession}
			<tr>
				<td>{trainingSession.id}</td>
				<td>{trainingSession.trainingSetId}</td>
				<td>{trainingSession.p1Score}</td>
				<td>{trainingSession.p2Score}</td>
				<td>{trainingSession.p3Score}</td>
				<td>{trainingSession.p4Score}</td>
				<td>{trainingSession.p5Score}</td>
				<td>{trainingSession.fullScore}</td>
				<td>{trainingSession.p1Time}</td>
				<td>{trainingSession.p2Time}</td>
				<td>{trainingSession.p3Time}</td>
				<td>{trainingSession.p4Time}</td>
				<td>{trainingSession.p5Time}</td>
				<td>{trainingSession.fullTime}</td>
				<td>{trainingSession.finished}</td>
				{#if !trainingSession.finished}
					<td>
						<button
							onclick={() =>
								goToSession(trainingSession.id)}
						>
							Go to session
						</button>
					</td>
				{/if}
			</tr>
		{/each}
	</tbody>
</table>
