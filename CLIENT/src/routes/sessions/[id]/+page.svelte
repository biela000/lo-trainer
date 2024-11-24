<script lang="ts">
	import { goto } from '$app/navigation';
	import { API_URL } from '$lib/api/constants';
	import PuzzleProgress from '$lib/components/PuzzleProgress.svelte';
	import type { TrainingSession } from '$lib/types/training-session';
	import type { TrainingSet } from '$lib/types/training-set';
	import type { PageProps } from './proxy+page';
	import type { ChangeEventHandler } from 'svelte/elements';

	let { data }: { data: PageProps } = $props();

	let trainingSession = $state<TrainingSession | undefined>();
	let trainingSet = $state<TrainingSet | undefined>();

	$effect(() => {
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

	const changePuzzleScore: ChangeEventHandler<HTMLInputElement> = (event) => {
		// const { name, value } = event.target as HTMLInputElement;
		const { name, value } = event.currentTarget;
		// @ts-expect-error I did a stupid thing when it comes to storing data in db
		trainingSession[name] = +value;
	};

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
				saveSession();
			}, 15000);
			return () => clearInterval(timer);
		} else if (trainingSession && timer) {
			return () => clearInterval(timer);
		}
	});

	const saveSession = async () => {
		if (!trainingSession) return;

		trainingSession.fullTime =
			+trainingSession.p1Time +
			trainingSession.p2Time +
			trainingSession.p3Time +
			trainingSession.p4Time +
			trainingSession.p5Time;

		trainingSession.fullScore = Math.floor(
			(trainingSession.p1Score +
				trainingSession.p2Score +
				trainingSession.p3Score +
				trainingSession.p4Score +
				trainingSession.p5Score) /
				5
		);

		console.log(trainingSession);

		await fetch(`${API_URL}/training_sessions/${data.id}`, {
			method: 'PUT',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify(trainingSession)
		});
	};

	const finishSession = async () => {
		if (!trainingSession) return;

		trainingSession.finished = true;

		await saveSession();

		goto('/');
	};
</script>

<h1>Session id #{trainingSession?.id}</h1>
{#if trainingSet}
	<PuzzleProgress {puzzles} {changePuzzleScore} {changeCurrentPuzzle} {currentPuzzle} />
	<button onclick={finishSession}>Finish session</button>
	<button onclick={saveSession}>Save session</button>
{/if}
