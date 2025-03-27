<script lang="ts">
	import BackButton from '$lib/components/navigation/BackButton.svelte';
	import type { PageData } from './$types';

	interface Props {
		data: PageData;
	}

	let { data }: Props = $props();
	const { eventID, eventInformation } = data;
	// TODO: make the content dynamic with the data that I get from the server.
	// Need Livio to know what to add to that card
	// TODO: remove the weird transition that I get from this page
	// TODO: all the field are not compulsory add some variables to handle it. ask Livio for the option that he needs

	import { createHorizontalSwipeHandler } from '$lib/scripts/swipe';
	function swipeBack(direction: 'left' | 'right') {
		const backButton = document.querySelector('.back') as HTMLButtonElement;
		if (direction === 'right') {
			backButton.click();
		}
	}
	const swipeBackAction = createHorizontalSwipeHandler(swipeBack);
</script>

<div class="content" id={eventID} use:swipeBackAction.action>
	<div class="overlay-background"></div>
	<div
		class="background"
		style="background:url({eventInformation.headerImg}); filter:blur(300px) brightness(30%);"
	></div>
	<div class="header container">
		<img src={eventInformation.headerImg} alt="background location illustration" />
		<div class="back-container">
			<BackButton color="white" />
		</div>
	</div>
	<div class="text-content flow">
		<div class="container">
			<p class="title">{eventInformation.name}</p>
		</div>
		<div class="container">
			<div class="card grid">
				<div>
					<p class="subtitle">
						{eventInformation.day}
						{eventInformation.date}
						{eventInformation.month}
					</p>
					<p>{eventInformation.time}</p>
				</div>
				<div class="separator"></div>
				<div>
					<p class="subtitle">{eventInformation.address}</p>
					<p>{eventInformation.city}</p>
					<div
						class="img-container"
						style="background-image: url({eventInformation.mapImg}); background-size: cover;"
					></div>
				</div>
				<div class="card-overlay"></div>
			</div>
		</div>
		<div class="container">
			<div class="card grid">
				<div>
					<p class="subtitle">Dress Code</p>
					<p>Un long paragraphe pour decrire cette partie</p>
				</div>
				<div class="card-overlay"></div>
			</div>
		</div>
		<div class="container">
			<div class="card grid">
				<div>
					<p class="subtitle">Food</p>
					<p>Un long paragraphe pour decrire cette partie</p>
				</div>
				<div class="card-overlay"></div>
			</div>
		</div>
		<div class="container">
			<div class="card grid">
				<div>
					<p class="subtitle">Musique</p>
					<p>Un long paragraphe pour decrire cette partie</p>
				</div>
				<div class="card-overlay"></div>
			</div>
		</div>
		<div class="container">
			<div class="card grid">
				<div>
					<p class="subtitle">Ta reservation massage</p>
					<p>Un long paragraphe pour decrire cette partie</p>
				</div>
				<div class="separator"></div>
				<div>
					<p class="subtitle">Ton practicien</p>
					<p>Met un petit truc avatar + nom + lien vers la fiche prestataire</p>
				</div>
				<div class="card-overlay"></div>
			</div>
		</div>
	</div>
</div>

<style>
	.content {
		padding-block: 2rem 3rem;
		position: relative;
	}
	.overlay-background,
	.background {
		position: absolute;
		left: 0;
		top: 0;
		width: 100vw;
		min-height: 100%;
		z-index: -1;
	}
	.background {
		background-size: cover;
		background-position: center;
	}
	.overlay-background {
		background-color: hsl(var(--clr-dark-ternary));
	}
	.header {
		position: relative;
	}
	.header img {
		border-radius: 0.75rem;
		width: 100vw;
		inline-size: 100%;
		aspect-ratio: 16/9;
	}
	.back-container {
		background: transparent;
		position: absolute;
		z-index: 2;
		--pos: 1rem;
		top: var(--pos);
		left: calc(var(--pos) + 1.5rem);
		background-color: hsl(var(--clr-dark-primary));
		padding: 0.5rem;
		border-radius: 100%;
		display: grid;
		place-content: center;
	}
	.text-content {
		padding-top: 2rem;
	}
	.separator {
		width: 100%;
		height: 1px;
		background-color: hsl(var(--clr-stroke));
		background-color: hsl(var(--clr-light-ternary));
		opacity: 0.125;
	}
	.card {
		--card-border-radius: 1rem;
		--card-padding: 1rem;
		padding: var(--card-padding);
		border-radius: var(--card-border-radius);
		backdrop-filter: blur(15px) saturate(150%);
		color: hsl(var(--clr-light-ternary));
		box-shadow:
			rgba(60, 64, 67, 0.3) 0px 1px 2px 0px,
			rgba(60, 64, 67, 0.15) 0px 1px 3px 1px;
		position: relative;
		overflow: hidden;
	}
	.card-overlay {
		position: absolute;
		top: 0;
		left: 0;
		width: 100%;
		height: 100%;
		background-color: rgba(0, 0, 0, 0.25);
		z-index: -1;
	}
	.title {
		margin: 0;
		color: hsl(var(--clr-light-primary));
		font-weight: 700;
		font-size: var(--fs-2);
		text-wrap: balance;
		line-height: 1.2;
	}
	.subtitle {
		color: hsl(var(--clr-grey-200));
		font-size: var(--fs-1);
		font-weight: 500;
	}
	.img-container {
		margin-top: 0.5rem;
		height: 120px;
		width: 100%;
		border-radius: max(calc(var(--card-border-radius) - var(--card-padding)), 0.25rem);
	}
</style>
