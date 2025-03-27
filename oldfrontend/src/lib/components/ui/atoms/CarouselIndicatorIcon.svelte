<script lang="ts">
	interface Props {
		active?: boolean;
		isLight?: boolean;
		// the unit is in rem for now
		width?: number;
	}

	let { active = false, isLight = true, width = 0.75 }: Props = $props();
</script>

<div
	class="dot"
	class:active
	class:isLight
	style="width: {active ? `calc(3 * ${width}rem)` : `${width}rem`};"
></div>

<style>
	.dot {
		aspect-ratio: 1;
		background-color: hsl(var(--clr-light-secondary));
		background-color: hsl(var(--clr-light-ternary));
		border-radius: 100%;
	}
	.dot.active {
		aspect-ratio: unset;
		border-radius: 0.5rem;
		background-color: hsl(var(--clr-light-primary));
		/* TODO:find the right value for that transition */
		transition: width 0.5s ease;
	}
	.dot.isLight {
		background-color: hsl(var(--clr-dark-secondary));
	}
	.dot.isLight.active {
		background-color: hsl(var(--clr-dark-primary));
	}
</style>
