<script lang="ts">
	interface Props {
		legend: string;
		checkboxes: { value: string; checked: boolean }[];
	}

	let { legend, checkboxes = $bindable([] as { value: string; checked: boolean }[]) }: Props =
		$props();

	const toggleCheckbox = (index: number) => {
		checkboxes[index].checked = !checkboxes[index].checked;
	};
</script>

<fieldset>
	<legend>{legend}</legend>
	{#each checkboxes as { value, checked }, index}
		<div>
			<label for={value}>{value}</label>
			<input
				type="checkbox"
				id={value}
				name={value}
				{value}
				{checked}
				onchange={() => toggleCheckbox(index)}
			/>
		</div>
	{/each}
</fieldset>

<style>
	fieldset {
		display: flex;
		flex-direction: column;
		width: max-content;
		gap: 0.5rem;
		padding: 0.5rem;
	}
</style>
