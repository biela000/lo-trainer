<script lang="ts">
	import type { ChangeEventHandler } from 'svelte/elements';

	let {
		puzzles,
		changePuzzleScore,
		changeCurrentPuzzle,
		currentPuzzle
	}: {
		puzzles: {
			link: string;
			score: number;
			time: number;
		}[];
		changePuzzleScore: ChangeEventHandler<HTMLInputElement>;
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
				<td>
					<input
						type="number"
						min="0"
						max="100"
						name={`p${i + 1}Score`}
						onchange={changePuzzleScore}
						value={score}
					/>
				</td>
				<td>{new Date(time).getHours()}h {new Date(time).getMinutes()}m</td>
				<td>
					<button onclick={changeCurrentPuzzle(i)}>
						{currentPuzzle === i ? 'Stop' : 'Start'}
					</button>
				</td>
			</tr>
		{/each}
	</tbody>
</table>

<style>
	table {
		margin-top: 1rem;
	}
	tr {
		position: relative;
	}
	td {
		padding: 1rem;
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
		padding: 0.5rem 2rem;
		border: none;
		background-color: #000;
		color: #fff;
	}
</style>
