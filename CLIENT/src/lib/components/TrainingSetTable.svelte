<script lang="ts">
	import type { TrainingSet } from '$lib/types/training-set';

	interface Props {
		trainingSets: TrainingSet[];
	}

	let { trainingSets }: Props = $props();
</script>

<table>
	<thead>
		<tr>
			<th>Id</th>
			<th>P1 Link</th>
			<th>P2 Link</th>
			<th>P3 Link</th>
			<th>P4 Link</th>
			<th>P5 Link</th>
		</tr>
	</thead>
	<tbody>
		{#each trainingSets as { id, p1Link, p2Link, p3Link, p4Link, p5Link }}
			<tr>
				<td>{id}</td>
				<td><a href={p1Link} target="_blank">{p1Link}</a></td>
				<td><a href={p2Link} target="_blank">{p2Link}</a></td>
				<td><a href={p3Link} target="_blank">{p3Link}</a></td>
				<td><a href={p4Link} target="_blank">{p4Link}</a></td>
				<td><a href={p5Link} target="_blank">{p5Link}</a></td>
			</tr>
		{/each}
	</tbody>
</table>
