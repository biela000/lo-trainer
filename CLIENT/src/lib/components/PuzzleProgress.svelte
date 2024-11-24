<script lang="ts">
	let {
		puzzles,
		changeCurrentPuzzle,
		currentPuzzle
	}: {
		puzzles: {
			link: string;
			score: number;
			time: number;
		}[];
		changeCurrentPuzzle: (index: number) => () => void;
		currentPuzzle: number | undefined;
	} = $props();
</script>

<table>
	<thead>
		<tr>
			<th>Id</th>
			<th>Puzzle</th>
			<th>Score</th>
			<th>Time</th>
		</tr>
	</thead>
	<tbody>
		{#each puzzles as { link, score, time }, i}
			<tr>
				<td>{i + 1}</td>
				<td>
					<a href={link} target="_blank" rel="noopener noreferrer">
						Puzzle {i + 1}
					</a>
				</td>
				<td>{score}</td>
				<td>{new Date(time).getMinutes()}m</td>
				<td>
					<button onclick={changeCurrentPuzzle(i)}>
						{currentPuzzle === i ? 'Stop' : 'Start'}
					</button>
				</td>
			</tr>
		{/each}
	</tbody>
</table>
