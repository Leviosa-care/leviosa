<script lang="ts">
	import { ChevronRight } from 'lucide-svelte';

	interface Props {
		action: () => void;
		name: string;
		icon: typeof import('lucide-svelte').Icon;
	}

	let { action, name, icon: Icon }: Props = $props();
</script>

<button onclick={action} class="field flex">
	<div class="field-left flex">
		<Icon />
		<p>{name}</p>
	</div>
	<ChevronRight color="black" />
</button>

<style>
	.field {
		justify-content: space-between;
		align-items: center;
		padding: 0.5rem;
		width: 100%;
		/* border-bottom: 1px solid hsl(var(--clr-stroke)); */
		/* TODO: test that one to see if that thing is intersting. */
		/* background-color: hsl(var(--clr-light-secondary)); */
		background-color: hsl(var(--clr-light-primary));
		padding-block: 0.75rem;
		padding-block: 1rem;
		border-bottom: 1px solid hsl(var(--clr-stroke));
		color: hsl(var(--clr-dark-secondary));
		color: hsl(var(--clr-grey-600));
	}
</style>
