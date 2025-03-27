<script lang="ts">
	import { run } from 'svelte/legacy';

	import { Share2, Download, ChevronLeft, ChevronRight } from 'lucide-svelte';
	import { createHorizontalSwipeHandler } from '$lib/scripts/swipe';

	interface Props {
		images: string[];
		openDrawer: () => void;
	}

	let { images, openDrawer }: Props = $props();

	let slidePosition: number = $state(0);
	const leftArrowClass = 'left-arrow';
	const rightArrowClass = 'right-arrow';

	function swipeImage(direction: 'left' | 'right') {
		console.log('the count brother in the swipe function', images.length - 1);
		const leftCondition = direction === 'left' && slidePosition < images.length - 1;
		const rightCondition = direction === 'right' && slidePosition > 0;
		if (leftCondition) slidePosition++;
		else if (rightCondition) slidePosition--;
	}
	const swipeAction = createHorizontalSwipeHandler(swipeImage);

	console.log('the count that I should have is:', images.length - 1);
	function handleClick(e: MouseEvent) {
		const target = e.currentTarget as HTMLButtonElement;
		const leftSwipeCondition = target.classList.contains(leftArrowClass) && slidePosition > 0;
		const rightSwipeCondition =
			target.classList.contains(rightArrowClass) && slidePosition < images.length - 1;

		if (leftSwipeCondition) slidePosition--;
		else if (rightSwipeCondition) {
			slidePosition++;
		}
	}

	run(() => {
		console.log('the new value of slidePosition:', slidePosition);
	});
	let transformStyle = $derived(`translateX(calc(-100% * ${slidePosition}))`);
	const arrowPosition = '5%';
</script>

<div class="content" use:swipeAction.action>
	<div class="flow">
		<div class="counter container">
			<p style="color: hsl(var(--clr-light-ternary));text-align: right;">
				{slidePosition + 1}/{images.length}
			</p>
		</div>
		<div class="images-container">
			<div class="images" style="transform: {transformStyle}; --flow-space: 1rem;">
				<button
					onclick={handleClick}
					class="arrows {leftArrowClass}"
					style="left: calc({arrowPosition} + ({slidePosition} * 100%));"
				>
					<ChevronLeft />
				</button>
				{#each images as image}
					<div class="img">
						<img src={image} alt="party" />
					</div>
				{/each}
				<button
					onclick={handleClick}
					class="arrows {rightArrowClass}"
					style="right: calc({arrowPosition} - ({slidePosition} * 100%));"
				>
					<ChevronRight />
				</button>
			</div>
		</div>
		<div class="icons container flex" style="margin-top: 2rem;">
			<Share2 color="hsl(var(--clr-light-ternary))" />
			<button onclick={openDrawer} class="dark-ternary-content toggler">
				Voir nos partenaires
			</button>
			<Download color="hsl(var(--clr-light-ternary))" />
		</div>
	</div>
</div>

<style>
	.content {
		height: 100vh;
		display: grid;
		place-content: center;
		background-color: black;
		--max-width: 600px;
	}
	.images-container {
		overflow-x: hidden;
		width: 100vw;
	}
	.images {
		display: flex;
		align-items: center;
		transition: transform 0.3s ease;
		width: 100%;
		/* NOTE: test for the buttons placement */
		position: relative;
	}
	.img {
		min-width: 100vw;
		border-radius: 0.5rem;
		display: flex;
		align-items: center;
		gap: 2rem;
	}
	/* TODO: use some media queries to make the thing go bigger brother */
	img {
		width: min(var(--max-width), 100%);
		margin-inline: auto;
		aspect-ratio: initial;
	}
	.arrows {
		display: none;
		background-color: hsl(var(--clr-light-primary));
		padding: 0.5rem;
		border-radius: 100%;
		position: absolute;
		top: 50%;
		transform: translateY(-50%);

		@media (width > 768px) {
			display: initial;
		}
	}
	.left-arrow {
		transform: translateX(-50%);
	}
	.right-arrow {
		transform: translateX(50%);
	}

	.icons {
		text-align: center;
		justify-content: space-between;
		align-items: center;
	}
	.counter,
	.icons {
		max-width: var(--max-width);
		margin-inline: auto;
	}
	.toggler {
		background: transparent;
		padding: 0.5rem 1rem;
		/* background-color: hsl(var(--clr-dark-primary)); */
		background-color: hsl(var(--clr-light-secondary));
		color: hsl(var(--clr-dark-primary));
		border-radius: 0.5rem;
		font-weight: 500;
	}
</style>
