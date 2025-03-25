<script lang="ts">
	import type { PageData } from './$types';

	import { ChevronLeft, SlidersHorizontal } from 'lucide-svelte';

	import Partenaire from './Partenaire.svelte';
	import Slide from './Slide.svelte';
	import Drawer from '$lib/components/Drawer.svelte';

	import { createVerticalSwipeHandler } from '$lib/scripts/swipe';

	let isOpen = $state(false);

	interface Props {
		data: PageData;
	}

	let { data }: Props = $props();

	let VerticalSlidePosition: number = $state(0);
	function swipeCarousel(direction: 'top' | 'bottom') {
		const topCondition = direction === 'top' && VerticalSlidePosition < events.length - 1;
		const bottomCondition = direction === 'bottom' && VerticalSlidePosition > 0;
		if (topCondition) VerticalSlidePosition++;
		else if (bottomCondition) VerticalSlidePosition--;
	}

	const swipeAction = createVerticalSwipeHandler(swipeCarousel);

	// TODO: handle the content on the drawer when clicking on the share and download icon
	// TODO: make that thing responsive too.
	let { events, eventID } = data;
</script>

<div class="content">
	<div class="header">
		<div class="header-content container flex">
			<ChevronLeft />
			<p class="title">
				Event du {events[VerticalSlidePosition].event.date}
			</p>
			<SlidersHorizontal />
		</div>
		<div class="separator"></div>
	</div>
	<div class="carousel">
		<div
			class="slides"
			style="transform: translateY({`calc(-100% * ${VerticalSlidePosition})`});"
			use:swipeAction.action
		>
			{#each events as event}
				<Slide images={event.images} openDrawer={() => (isOpen = true)} />
			{/each}
		</div>
	</div>
</div>
<Drawer bind:isOpen>
	<div class="flex" style="justify-content: center;">
		<div class="swipe-down"></div>
	</div>
	<div class="partners flow">
		{#each events[VerticalSlidePosition].partners as partner}
			<Partenaire {partner} />
		{/each}
	</div>
</Drawer>

<style>
	.carousel {
		overflow-y: hidden;
		height: 100vh;
	}
	.slides {
		transition: transform 0.3s ease;
		height: 100%;
	}
	.swipe-down {
		background-color: hsl(var(--clr-dark-ternary));
		height: 6px;
		width: 20%;
		border-radius: 0.5rem;
	}
	.header {
		position: fixed;
		top: 0;
		left: 0;
		right: 0;
		z-index: 2;
		/* if the background does not exist */
		background: transparent;
		color: hsl(var(--clr-light-primary));
	}
	.header-content {
		padding-block: 1rem;
		justify-content: space-between;
		align-items: center;
		/* TODO: make the header fixed */
	}
	.separator {
		width: 100%;
		height: 1px;
		background-color: hsl(var(--clr-stroke));
	}
	.title {
		font-size: 1rem;
		font-weight: 500;
	}
	.partners {
		margin-top: 2rem;
	}
</style>
