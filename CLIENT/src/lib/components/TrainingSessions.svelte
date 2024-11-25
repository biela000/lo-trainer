<script lang="ts">
	import { goto } from '$app/navigation';
	import { API_URL } from '$lib/api/constants';
	import type { TrainingSession } from '$lib/types/training-session';

	let trainingSessions: TrainingSession[] = $state([]);

	const getTrainingSessions = async () => {
		const response = await fetch(`${API_URL}/training_sessions`);
		trainingSessions = (await response.json()) as TrainingSession[];
	};

	$effect(() => {
		getTrainingSessions();
	});

	const goToSession = (id: number) => {
		goto(`/sessions/${id}`);
	};

	const deleteSession = async (id: number) => {
		await fetch(`${API_URL}/training_sessions/${id}`, {
			method: 'DELETE'
		});

		getTrainingSessions();
	};
</script>

<h2>Your Training Sessions</h2>
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
				<td>{new Date(trainingSession.p1Time).getMinutes()}m</td>
				<td>{new Date(trainingSession.p2Time).getMinutes()}m</td>
				<td>{new Date(trainingSession.p3Time).getMinutes()}m</td>
				<td>{new Date(trainingSession.p4Time).getMinutes()}m</td>
				<td>{new Date(trainingSession.p5Time).getMinutes()}m</td>
				<td>{new Date(trainingSession.fullTime).getMinutes()}m</td>
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
				<td>
					<button onclick={() => deleteSession(trainingSession.id)}>
						Delete
					</button>
				</td>
			</tr>
		{/each}
	</tbody>
</table>

<style>
	table {
		width: 100%;
	}
	h2 {
		margin-top: 2rem;
	}
	tr {
		position: relative;
	}
	tr::before {
		content: '';
		position: absolute;
		width: 100%;
		height: 0.1rem;
		background-color: #000;
		bottom: -0.1rem;
		left: 0;
	}
	button {
		background-color: darkred;
		color: white;
		border: none;
		padding: 0.2rem 1rem;
	}
</style>
