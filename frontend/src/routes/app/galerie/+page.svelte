<script lang="ts">
	import Tabs from '$lib/components/Tabs.svelte';
	import SubOffer from './SubOffer.svelte';
	import BackButton from '$lib/components/navigation/BackButton.svelte';

	import { Search } from 'lucide-svelte';

	import type { PageData } from './$types';

	function displayPhoto(e: MouseEvent) {
		const target = e.currentTarget as HTMLButtonElement;
		console.log('the ID for that thing is :', target.id);
	}
	type Content = 'Photos' | 'Videos';
	let content: Content = $state('Photos');
	let subcontent: Content | 'Albums';
	$effect(() => {
		let subcontent: Content | 'Albums' = $state(content);
	});
	function changeContent(newContent: Content) {
		content = newContent;
	}
	let mainTabs = [
		{ name: 'Photos', action: () => changeContent('Photos') },
		{ name: 'Videos', action: () => changeContent('Videos') }
	];
	// TODO: add the part when there is no photos, we need something to indicate that
	// TODO: do the href thing to get to the right page
	import { createHorizontalSwipeHandler } from '$lib/scripts/swipe';
	interface Props {
		data: PageData;
	}

	let { data }: Props = $props();
	function swipeBack(direction: 'left' | 'right') {
		const backButton = document.querySelector('.back') as HTMLButtonElement;
		if (direction === 'right') {
			backButton.click();
		}
	}
	const swipeBackAction = createHorizontalSwipeHandler(swipeBack);
	let { eventsPhotos, eventsVideos } = $derived(data);
	let subTabs = $derived([
		{ name: content, action: () => (subcontent = content) },
		{ name: 'Albums', action: () => (subcontent = 'Albums') }
	]);
</script>

<div class="content flow" use:swipeBackAction.action>
	<div class="backbutton">
		<BackButton pathname="/app/reservations" />
	</div>
	<div class="container">
		<Tabs offers={mainTabs} isSecondary={true} />
	</div>
	<div style="margin-inline: auto; max-width: 200px;">
		<SubOffer offers={subTabs} isSecondary={true} />
	</div>
	<div class="container">
		<div class="search flex">
			<Search size={16} color="hsl(var(--clr-dark-primary))" />
			<input type="text" placeholder="Photos, Date, Lieux..." />
		</div>
	</div>
	<div class="photos grid" style="--gap:1px;">
		{#if content === 'Photos'}
			{#each eventsPhotos as eventPhoto}
				{#each eventPhoto.photos as content, index}
					<button aria-label="photo" onclick={displayPhoto} id="photo_{index}" class="photo-button">
						<div
							class="photo"
							style="background-image: url({content}); background-size: cover;;"
						></div>
					</button>
				{/each}
			{/each}
		{:else}
			{#each eventsVideos as eventVideo}
				{#each eventVideo.videos as video, index}
					<button aria-label="video" onclick={displayPhoto} id="photo_{index}" class="video-button">
						<div
							class="video"
							style="background-image: url({video.thumbnail}); background-size: cover;;"
						></div>
						<p class="timer">{video.duration}</p>
					</button>
				{/each}
			{/each}
		{/if}
	</div>
</div>

<style>
	.content {
		padding-top: 1rem;
		padding-bottom: 7rem;
	}
	.backbutton {
		background-color: #f7f7f9;
		width: fit-content;
		height: fit-content;
		margin-inline: 0.75rem;
		padding: 0.5rem;
		border-radius: 100%;
		display: grid;
		place-content: center;
	}
	.search {
		background-color: hsl(var(--clr-light-secondary));
		height: 40px;
		width: 100%;
		align-items: center;
		padding: 0.5rem;
		border-radius: 0.5rem;
	}
	.search input {
		background: transparent;
		width: 100%;
	}
	.search input:is(:global(:focus, :hover)) {
		outline: none;
	}
	.photos {
		margin-inline: auto;
		grid-template-columns: repeat(3, 1fr);
	}
	.photo-button,
	.video-button {
		background-color: #d9d9d9;
		aspect-ratio: 4/3;
		max-height: 300px;
	}
	.photo,
	.video {
		height: 100%;
		aspect-ratio: 4/3;
	}
	.video-button {
		position: relative;
	}
	.timer {
		--distance: 0.25rem;
		position: absolute;
		bottom: var(--distance);
		right: var(--distance);
		font-size: 0.75rem;
		font-weight: 600;
		color: hsl(var(--clr-light-primary));
	}
</style>
