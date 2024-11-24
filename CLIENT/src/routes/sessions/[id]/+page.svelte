<script lang="ts">
	import { API_URL } from '$lib/api/constants';
	import PuzzleProgress from '$lib/components/PuzzleProgress.svelte';
	import type { TrainingSession } from '$lib/types/training-session';
	import type { TrainingSet } from '$lib/types/training-set';
	import type { PageProps } from './proxy+page';

	let { data }: { data: PageProps } = $props();

	let trainingSession = $state<TrainingSession | undefined>();
	let trainingSet = $state<TrainingSet | undefined>();

	const getTrainingSession = async () => {
		const trainingSessionResponse = await fetch(
			`${API_URL}/training_sessions/${data.id}`
		);
		trainingSession = (await trainingSessionResponse.json()) as TrainingSession;
	};

	const getTrainingSet = async () => {
		if (!trainingSession) return;

		const trainingSetResponse = await fetch(
			`${API_URL}/training_sets/${trainingSession.trainingSetId}`
		);
		trainingSet = (await trainingSetResponse.json()) as TrainingSet;
	};

	$effect(() => {
		getTrainingSession().then(getTrainingSet);
	});

	let puzzles = $derived([
		{
			link: trainingSet?.p1Link ?? '',
			score: trainingSession?.p1Score ?? 0,
			time: trainingSession?.p1Time ?? 0
		},
		{
			link: trainingSet?.p2Link ?? '',
			score: trainingSession?.p2Score ?? 0,
			time: trainingSession?.p2Time ?? 0
		},
		{
			link: trainingSet?.p3Link ?? '',
			score: trainingSession?.p3Score ?? 0,
			time: trainingSession?.p3Time ?? 0
		},
		{
			link: trainingSet?.p4Link ?? '',
			score: trainingSession?.p4Score ?? 0,
			time: trainingSession?.p4Time ?? 0
		},
		{
			link: trainingSet?.p5Link ?? '',
			score: trainingSession?.p5Score ?? 0,
			time: trainingSession?.p5Time ?? 0
		}
	]);

	let currentPuzzle = $state<number | undefined>(undefined);

	const changeCurrentPuzzle = (index: number) => () => {
		currentPuzzle = currentPuzzle === index ? undefined : index;
	};

	let timer: number | undefined;
	let startTime: number;

	$effect(() => {
		if (trainingSession && currentPuzzle !== undefined) {
			startTime = Date.now();
			timer = setInterval(() => {
				const currentPuzzleTimePropName = `p${currentPuzzle! + 1}Time`;
				// @ts-expect-error I did a stupid thing when it comes to storing data in db
				trainingSession[currentPuzzleTimePropName] = Date.now() - startTime;
			}, 15000);
			return () => clearInterval(timer);
		} else if (trainingSession && timer) {
			return () => clearInterval(timer);
		}
	});
</script>

<h1>Session id #{trainingSession?.id}</h1>
{#if trainingSet}
	<PuzzleProgress {puzzles} {changeCurrentPuzzle} {currentPuzzle} />
{/if}
