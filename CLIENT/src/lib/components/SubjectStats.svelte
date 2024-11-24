<script lang="ts">
	import { API_URL } from '$lib/api/constants';

	let subjectStats:
		| {
				compoundingAvg: number;
				morphologyAvg: number;
				numbersAvg: number;
				phonologyAndPhoneticsAvg: number;
				semanticsAvg: number;
				syntaxAvg: number;
				writingSystemAvg: number;
		  }
		| undefined = $state();

	$effect(() => {
		const getSubjectStats = async () => {
			const response = await fetch(`${API_URL}/training_sessions/stats/subjects`);
			subjectStats = await response.json();
		};

		getSubjectStats();
	});
</script>

{#if subjectStats}
	<h2>Your Subject Stats</h2>
	<table>
		<thead>
			<tr>
				<th>Subject</th>
				<th>Average score</th>
			</tr>
		</thead>
		<tbody>
			<tr>
				<td>Compounding</td>
				<td>{subjectStats.compoundingAvg}</td>
			</tr>
			<tr>
				<td>Morphology</td>
				<td>{subjectStats.morphologyAvg}</td>
			</tr>
			<tr>
				<td>Numbers</td>
				<td>{subjectStats.numbersAvg}</td>
			</tr>
			<tr>
				<td>Phonology and Phonetics</td>
				<td>{subjectStats.phonologyAndPhoneticsAvg}</td>
			</tr>
			<tr>
				<td>Semantics</td>
				<td>{subjectStats.semanticsAvg}</td>
			</tr>
			<tr>
				<td>Syntax</td>
				<td>{subjectStats.syntaxAvg}</td>
			</tr>
			<tr>
				<td>Writing System</td>
				<td>{subjectStats.writingSystemAvg}</td>
			</tr>
		</tbody>
	</table>
{/if}
