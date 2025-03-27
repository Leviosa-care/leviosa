<script lang="ts">
	import type { PageData } from './$types';
	interface Props {
		data: PageData;
	}
	let { data }: Props = $props();
	const { messages } = data;

	import { onNavigate } from '$app/navigation';

	onNavigate((navigation) => {
		if (!document.startViewTransition) return;
		return new Promise((resolve) => {
			document.startViewTransition(async () => {
				resolve();
				await navigation.complete;
			});
		});
	});

	import { createHorizontalSwipeHandler } from '$lib/scripts/swipe';
	function swipeBack(direction: 'left' | 'right') {
		const backButton = document.querySelector('.back') as HTMLButtonElement;
		if (direction === 'right') {
			backButton.click();
		}
	}
	const swipeBackAction = createHorizontalSwipeHandler(swipeBack);
	// TODO: make sure that the input in the bottom make the keyboard goes out and but stay on top of the keyboard thing. Need to publish the app to check that
	import Conversations from './Conversations.svelte';
</script>

<div class="content" use:swipeBackAction.action>
	<Conversations {messages} />
</div>

<style>
	.content {
		view-transition-name: pushed;
	}
	/* transitions animations */
	@keyframes back-old {
		to {
			transform: translateX(100%);
		}
	}
	@keyframes back-new {
		from {
			transform: translateX(-30%);
		}
		to {
			transform: translateX(0%);
		}
	}
	/* targetting pages with animations */
	:root::view-transition-old(pushed) {
		animation: 250ms ease-out both back-old;
	}

	:root::view-transition-new(pushing) {
		animation: 250ms ease-out both back-new;
	}
</style>
