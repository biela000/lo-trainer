<script lang="ts">
	import { goto } from '$app/navigation';
	import { API_URL } from '$lib/api/constants';
	import type { TrainingSession } from '$lib/types/training-session';
	import type { TrainingSet } from '$lib/types/training-set';
	import CheckboxFieldset from './CheckboxFieldset.svelte';
	import TrainingSetTable from './TrainingSetTable.svelte';
	import YearFieldset from './YearFieldset.svelte';

	let levels = $state([
		{
			value: 'Breakthrough',
			checked: false
		},
		{
			value: 'Foundation',
			checked: false
		},
		{
			value: 'Intermediate',
			checked: false
		},
		{
			value: 'Advanced',
			checked: false
		},
		{
			value: 'Round 2',
			checked: false
		}
	]);

	let subjects = $state([
		{
			value: 'Compounding',
			checked: false
		},
		{
			value: 'Morphology',
			checked: false
		},
		{
			value: 'Numbers',
			checked: false
		},
		{
			value: 'Phonology and Phonetics',
			checked: false
		},
		{
			value: 'Semantics',
			checked: false
		},
		{
			value: 'Syntax',
			checked: false
		},
		{
			value: 'Writing System',
			checked: false
		}
	]);

	let formats = $state([
		{
			value: 'Rosetta',
			checked: false
		},
		{
			value: 'Match-up',
			checked: false
		},
		{
			value: 'Monolingual',
			checked: false
		},
		{
			value: 'Pattern',
			checked: false
		},
		{
			value: 'Computational',
			checked: false
		},
		{
			value: 'Text',
			checked: false
		}
	]);

	const startYear = 2010;
	const endYear = new Date().getFullYear();
	const yearsArray = Array.from(
		{ length: endYear - startYear + 1 },
		(_, index) => startYear + index
	);

	let minYear = $state(startYear);
	let maxYear = $state(endYear);

	let minScore = $state(0);
	let maxScore = $state(100);

	const prepareCheckboxValues = (
		name: string,
		checkboxes: { value: string; checked: boolean }[]
	) =>
		checkboxes
			.filter((checkbox) => checkbox.checked)
			.map((checkbox) => checkbox.value)
			.join(`&${name}=`);

	let newTrainingSet = $state<TrainingSet | undefined>(undefined);

	const buildQueryString = () => {
		const selectedLevels = prepareCheckboxValues('levels', levels);
		const selectedSubjects = prepareCheckboxValues('subjects', subjects);
		const selectedFormats = prepareCheckboxValues('formats', formats);

		let queryString = '';
		queryString += selectedLevels ? `levels=${selectedLevels}` : '';
		queryString += selectedSubjects ? `&subjects=${selectedSubjects}` : '';
		queryString += selectedFormats ? `&formats=${selectedFormats}` : '';
		queryString += `&min_year=${minYear}`;
		queryString += `&max_year=${maxYear}`;
		queryString += `&min_score=${minScore}`;
		queryString += `&max_score=${maxScore}`;

		return queryString;
	};

	const handleSubmit = async () => {
		const queryString = buildQueryString();

		const response = await fetch(`${API_URL}/training_sets?${queryString}`, {
			method: 'POST'
		});

		newTrainingSet = (await response.json()) as TrainingSet;
	};

	const startTrainingSession = async () => {
		if (!newTrainingSet) return;

		const response = await fetch(`${API_URL}/training_sessions`, {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify({
				trainingSetId: newTrainingSet.id
			})
		});

		const { id: trainingSessionId } = (await response.json()) as TrainingSession;

		goto(`/sessions/${trainingSessionId}`);
	};
</script>

<form onsubmit={handleSubmit}>
	<div class="checkbox-fieldsets">
		<CheckboxFieldset
			legend="Choose the levels you want to include"
			bind:checkboxes={levels}
		/>
		<CheckboxFieldset
			legend="Choose the subjects you want your levels to be about"
			bind:checkboxes={subjects}
		/>
		<CheckboxFieldset
			legend="Choose the formats of your levels"
			bind:checkboxes={formats}
		/>
	</div>
	<div class="number-inputs">
		<YearFieldset {yearsArray} bind:minYear bind:maxYear />
		<fieldset>
			<legend>Score</legend>
			<label for="minScore">Min</label>
			<input
				type="number"
				id="minScore"
				bind:value={minScore}
				min={0}
				max={100}
			/>
			<label for="maxScore">Max</label>
			<input
				type="number"
				id="maxScore"
				bind:value={maxScore}
				min={0}
				max={100}
			/>
		</fieldset>
	</div>
	<button type="submit">Create training set</button>
</form>
{#if newTrainingSet}
	<TrainingSetTable trainingSets={[newTrainingSet]} />
	<button onclick={startTrainingSession}>Start training session</button>
{/if}

<style>
	fieldset {
		padding: 0.5em;
	}
	.checkbox-fieldsets,
	.number-inputs {
		display: flex;
		gap: 2rem;
	}
	button {
		margin-top: 0.5rem;
		background-color: #000;
		color: #eee;
		border: none;
		padding: 0.5rem 2rem;
	}
</style>
