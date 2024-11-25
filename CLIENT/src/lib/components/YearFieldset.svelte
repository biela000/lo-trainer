<script lang="ts">
	interface Props {
		minYear: number;
		maxYear: number;
		yearsArray: number[];
	}

	let { minYear = $bindable(), maxYear = $bindable(), yearsArray }: Props = $props();
</script>

<fieldset>
	<legend>Year</legend>
	<label for="minYear">From</label>
	<select id="minYear" bind:value={minYear}>
		{#each yearsArray as year}
			<option value={year}>{year}</option>
		{/each}
	</select>
	<label for="minYear">To</label>
	<select id="maxYear" bind:value={maxYear}>
		{#each yearsArray as year}
			<option value={year}>{year}</option>
		{/each}
	</select>
</fieldset>

<style>
	fieldset {
		padding: 0.5rem;
	}
</style>
