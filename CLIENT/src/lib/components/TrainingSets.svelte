<script lang="ts">
	import { goto } from '$app/navigation';
	import { API_URL } from '$lib/api/constants';
	import type { TrainingSession } from '$lib/types/training-session';
	import type { TrainingSet } from '$lib/types/training-set';

	let trainingSets: TrainingSet[] = $state([]);

	const fetchTrainingSets = async () => {
		const response = await fetch(`${API_URL}/training_sets`);

		trainingSets = await response.json();
	};

	$effect(() => {
		fetchTrainingSets();
	});

	const startTrainingSession = async (trainingSetId: number) => {
		const response = await fetch(`${API_URL}/training_sessions`, {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify({ trainingSetId })
		});

		const { id: trainingSessionId } = (await response.json()) as TrainingSession;

		goto(`/sessions/${trainingSessionId}`);
	};

	const deleteTrainingSet = async (trainingSetId: number) => {
		await fetch(`${API_URL}/training_sets/${trainingSetId}`, {
			method: 'DELETE'
		});

		fetchTrainingSets();
	};
</script>

<table>
	<thead>
		<tr>
			<th>Id</th>
			<th>P1</th>
			<th>P2</th>
			<th>P3</th>
			<th>P4</th>
			<th>P5</th>
			<th></th>
			<th></th>
		</tr>
	</thead>
	<tbody>
		{#each trainingSets as { id, p1Link, p2Link, p3Link, p4Link, p5Link }}
			<tr>
				<td>{id}</td>
				<td>{p1Link}</td>
				<td>{p2Link}</td>
				<td>{p3Link}</td>
				<td>{p4Link}</td>
				<td>{p5Link}</td>
				<td>
					<button onclick={() => startTrainingSession(id)}>
						Start session
					</button>
				</td>
				<td>
					<button onclick={() => deleteTrainingSet(id)}>
						Delete
					</button>
				</td>
			</tr>
		{/each}
	</tbody>
</table>
