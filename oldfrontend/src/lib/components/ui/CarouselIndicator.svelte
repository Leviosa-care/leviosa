<script lang="ts">
	import CarouselIndicatorIcon from './atoms/CarouselIndicatorIcon.svelte';
	interface Props {
		count?: number;
		activeIndex?: number;
		isLight?: boolean;
	}

	let { count = 1, activeIndex = 0, isLight = true }: Props = $props();
</script>

<div class="carousel-indicator flex" style="--gap: 0.5rem">
	{#each Array(count) as _, index}
		<CarouselIndicatorIcon {isLight} active={activeIndex === index} width={0.75} />
	{/each}
</div>

<style>
	.carousel-indicator {
		display: fl;
	}
</style>
