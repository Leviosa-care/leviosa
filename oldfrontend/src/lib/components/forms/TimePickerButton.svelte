<script lang="ts">
	import { Calendar, ChevronUp, ChevronDown } from 'lucide-svelte';

	// TODO: make that responsive
	const calendarSize = 20;
	const chevronSize = 12;

	interface Props {
		onClick: () => void;
		value?: string | undefined;
	}

	let { onClick, value = 'Selectionne un creneau' }: Props = $props();
</script>

<button type="button" class="toggle-button stroke flex" onclick={onClick}>
	<div class="icon">
		<Calendar size={calendarSize} strokeWidth={1.25} absoluteStrokeWidth={true} />
	</div>
	<input type="button" class="toggle-text" {value} />
	<div class="selector flex">
		<ChevronUp size={chevronSize} />
		<ChevronDown size={chevronSize} />
	</div>
</button>

<style>
	.toggle-button {
		position: relative;
		background: transparent;
		border-radius: 0.5rem;
		width: 100%;
		background-color: hsl(var(--clr-light-secondary));
		background-color: #f7f7f9;
	}
	.selector {
		position: absolute;
		right: 0;
		top: 0;
		width: 2.5rem;
		height: 100%;
		border-radius: 0.5rem;
		pointer-events: none;
		display: flex;
		align-items: center;
		justify-content: center;
		flex-direction: column;
		gap: 0rem;
	}
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
	.toggle-text {
		border-radius: 0.5rem;
		border: none;
		padding-block: calc(3 * var(--fs-0) / 4);
		padding-inline: 3rem;
		background: transparent;
		appearance: none;
		width: 100%;
	}
</style>
