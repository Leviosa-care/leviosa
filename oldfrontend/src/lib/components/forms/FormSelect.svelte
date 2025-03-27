<script lang="ts">
	const chevronSize = 12;
	const iconSize = 20;

	import { Calendar, ChevronDown, ChevronUp } from 'lucide-svelte';
	interface Props {
		optionName?: string;
		selectName?: string;
		options?: string[];
		icon?: typeof import('lucide-svelte').Icon | null;
	}

	let {
		optionName = 'une option',
		selectName = 'option',
		options = ['Option 1', 'Option 2', 'Option 3'],
		icon: Icon = Calendar
	}: Props = $props();

	// TODO:
	// - style options in the dropdown menu or do a complete custom component with divs instead of select and option tag.
	// - style the icon that goes before the placeholder for the input

	// if options is empty then make the first option that got the text not appear on click
</script>

<div class="container grid">
	<div class="select-container stroke flex" style="--gap: 2rem">
		{#if Icon}
			<div class="icon">
				<Icon size={iconSize} strokeWidth={1} absoluteStrokeWidth={true} />
			</div>
		{/if}
		<select class="selectElement" name={selectName} id={selectName}>
			<option value="" disabled selected>Selectionne {optionName}</option>
			{#each options as option}
				<option value={option}>{option}</option>
			{/each}
		</select>
		<div class="selector flex">
			<ChevronUp size={chevronSize} />
			<ChevronDown size={chevronSize} />
		</div>
	</div>
</div>

<style>
	.icon {
		position: absolute;
		left: 0;
		top: 0;
		height: 100%;
		width: 2.5rem;
		display: flex;
		align-items: center;
		justify-content: center;
	}
	.select-container {
		position: relative;
		background: transparent;
		border-radius: 0.5rem;
		width: 100%;
	}

	select {
		border-radius: 0.5rem;
		border: none;
		padding-block: calc(3 * var(--p) / 4);
		padding-inline: 3rem;
		background: transparent;
		appearance: none;
		width: 100%;
	}
	.selector {
		position: absolute;
		right: 0;
		top: 0;
		width: 2.5rem;
		height: 100%;
		border-radius: 0.5rem;
		pointer-events: none;
		/* NOTE: to center the icon in the container */
		display: flex;
		align-items: center;
		justify-content: center;
		flex-direction: column;
		gap: 0rem;
	}
	.container {
		height: 100vh;
		place-content: center;
		padding-inline: 0.75rem;
	}
</style>
