<script lang="ts">
	import type { Freelancer, Review } from '$lib/types';
	interface Props {
		freelancer: Freelancer | undefined;
	}
	let { freelancer }: Props = $props();
	// TODO: add some function to make the reservation with the specified freelancer

	import { Star, CalendarCheck } from 'lucide-svelte';
	import Button from '$lib/components/Button.svelte';

	type TabState = 'A propos' | 'Specialites' | 'Reviews';
	let tabStates: TabState[] = ['A propos', 'Specialites', 'Reviews'];

	let navState: TabState = $state(tabStates[0]);
	let totalPoints = $derived(
		freelancer?.reviews?.reduce((sum, review) => sum + review.note * review.count, 0)
	);
	let totalCounts = $derived(freelancer?.reviews?.reduce((sum, review) => sum + review.count, 0));
	let avgRating = $derived(Math.round(((totalPoints ?? 0) * 10) / (totalCounts ?? 1)) / 10);

	const widthbar = '5rem';

	function handleTabState(e: MouseEvent) {
		const target = e.currentTarget as HTMLButtonElement;
		navState = target.id as TabState;
	}

	function getBarfillWidth(index: number): string {
		const review: Review | undefined = freelancer?.reviews?.find((review) => review.note === index);
		if (!review || totalCounts === 0) return '0%';
		const percentage = (review.count / (totalCounts ?? 1)) * 100;
		return `${percentage}%`;
	}
</script>

<div class="content flow">
	<div class="header center flex">
		<div class="figure flex">
			<div class="side flex">
				<CalendarCheck color="white" size={20} />
				<p>143</p>
			</div>
			<p>Consulations</p>
		</div>
		<img
			class="avatar"
			src="https://pbs.twimg.com/media/FA9tIzxUUAUy80c.jpg"
			alt="avatar de {freelancer?.firstname} {freelancer?.lastname}"
		/>
		<div class="figure flex">
			<div class="side flex">
				<Star color="white" size={20} />
				<p>{totalCounts}</p>
			</div>
			<p>Reviews</p>
		</div>
	</div>
	<div class="title-container grid">
		<h3 class="center fs-h3">
			{freelancer?.firstname}
			{freelancer?.lastname}
		</h3>
		<p class="center">Leviosa massage, bien-etre</p>
	</div>
	<div class="tabs grid">
		{#each tabStates as tab}
			<button id={tab} class="tab" class:active={navState === tab} onclick={handleTabState}>
				{tab}
			</button>
		{/each}
	</div>
	<div class="text-wrapper">
		{#if navState === tabStates[0]}
			<p>
				Lorem ipsum dolor sit amet consectetur adipisicing elit. Ipsam quaerat blanditiis illum
				fugit. Reprehenderit ab quo quam ipsam cum dolore expedita tempore ratione. Repellat esse a,
				labore minima est quod sit repellendus?
			</p>
		{:else if navState === tabStates[1]}
			<p>Some description for the user</p>
			<p>
				Lorem ipsum dolor sit amet consectetur adipisicing elit. Ipsum aliquam ipsam perferendis
				excepturi cum a corporis laboriosam animi voluptatibus molestiae!
			</p>
		{:else}
			<div class="rating-wrapper container flex">
				<div class="flow">
					<div>
						<h3 class="fs-h3 fw-600">Reviews</h3>
						<div class="grid" style="--gap: -0.8rem;">
							<p class="fs-h1 fw-800">
								{avgRating}
								<span class="fs-paragraph fw-400">/ 5</span>
							</p>
							<p class="fs-paragraph">
								{totalCounts} rating{(totalCounts ?? 0 > 0) ? 's' : ''}
							</p>
						</div>
					</div>
					<div class="stars flex" style="--gap: 0.2rem;">
						{#each Array(5) as _, index}
							<Star fill={index < 4 ? 'black' : 'white'} />
						{/each}
					</div>
				</div>
				<div class="rating-detail-wrapper flex">
					{#each Array(5) as _, index}
						{@const rating = index + 1}
						<div class="rating-detail flex">
							<p class="fs-label">
								{rating} etoile{rating > 1 ? 's' : ''}
							</p>
							<div class="bar" style="width: {widthbar};">
								<div style="width: {getBarfillWidth(rating)}" class="bar-fill"></div>
							</div>
						</div>
					{/each}
				</div>
			</div>
		{/if}
	</div>
	<div class="grid" style="margin-top: 2rem;">
		<Button onClick={() => console.log('make a reservation brother')}
			>Reserve ta place avec {freelancer?.firstname}</Button
		>
	</div>
</div>

<style>
	.header {
		align-items: center;
		justify-content: space-between;
	}
	.figure {
		gap: 0.2rem;
		align-items: center;
		flex-direction: column;
		justify-content: center;
	}
	.side {
		align-items: center;
		color: hsl(var(--clr-light-primary));
		padding-inline: 0.5rem;
		background-color: hsl(var(--clr-dark-primary));
		border-radius: 0.5rem;
		width: 5rem;
		height: 2rem;
		gap: 0.5rem;
		justify-content: center;
	}
	.avatar {
		position: absolute;
		top: 0;
		left: 50%;
		transform: translate(-50%, -25%);

		border: 2px solid hsl(var(--clr-dark-primary));
		width: 80px;
		aspect-ratio: 1;
		background-color: hsl(var(--clr-dark-primary));
		border-radius: 0.5rem;
	}
	.title-container {
		gap: 0.2rem;
		margin-top: 2rem;
	}
	.tabs {
		--border-radius: 0.4rem;
		border-radius: calc(3 * var(--border-radius) / 2);
		margin-top: 2rem;
		background-color: hsl(var(--clr-dark-primary));
		padding: 0.3rem;

		grid-auto-flow: column;
		grid-auto-columns: minmax(max-content, 1fr);
		width: calc(100vw - 2rem);

		/* NOTE: the light colors brother */
		background-color: hsl(var(--clr-light-ternary));
	}
	.tab {
		background: transparent;
		padding: 0.3rem;
		color: hsl(var(--clr-light-primary));
		border-radius: var(--border-radius);
		font-weight: 500;

		/* NOTE: the light colors brother */
		/* color: hsl(var(--clr-dark-primary)); */
		color: hsl(var(--clr-dark-secondary));
	}
	.tab.active {
		color: hsl(var(--clr-dark-primary));
		background-color: hsl(var(--clr-light-primary));
		box-shadow: rgba(0, 0, 0, 0.18) 0px 2px 4px;
	}
	.text-wrapper {
		min-height: 20vh;
	}
	.rating-wrapper {
		margin-top: 2rem;
		align-items: center;
		justify-content: space-between;
	}
	.rating-detail-wrapper {
		flex-direction: column;
		gap: 0.4rem;
	}
	.rating-detail {
		align-items: center;
		justify-content: space-between;
	}
	.bar {
		--heightbar: 12px;
		border: 1px solid hsl(var(--clr-stroke));
		height: var(--heightbar);
		border-radius: 0.2rem;
	}
	.bar-fill {
		background-color: hsl(var(--clr-accent));
		height: var(--heightbar);
	}
</style>
