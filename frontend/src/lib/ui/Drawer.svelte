<script lang="ts">
	import { onDestroy, getContext } from 'svelte';
	import { browser } from '$app/environment';
	import { createVerticalSwipeHandler } from '$lib/scripts/swipe';

	interface Props {
		isOpen?: boolean;
		children?: import('svelte').Snippet;
	}

	let { children, isOpen = $bindable(false) }: Props = $props();

	// Layouts that use a scroll-jail (overflow on a container, not body) should set
	// this context to false so we don't apply scroll locking at all.
	const lockBodyScroll = (getContext<boolean>('drawer-lock-body-scroll') ?? true);

	let drawerEl: HTMLElement | null = null;
	let renderChildren = $state(false);
	let closeTimer: ReturnType<typeof setTimeout> | undefined;
	let scrollLocked = false;

	const TRANSITION_MS = 300;

	function handleKeydown(event: KeyboardEvent) {
		event.stopPropagation();
		if (event.key === 'Escape') {
			isOpen = false;
		}
	}

	function preventBodyTouchMove(e: TouchEvent) {
		if (drawerEl && drawerEl.contains(e.target as Node)) return;
		e.preventDefault();
	}

	function lockScroll() {
		document.body.style.overflow = 'hidden'; // desktop fallback
		document.addEventListener('touchmove', preventBodyTouchMove, { passive: false });
		scrollLocked = true;
	}

	function unlockScroll() {
		document.body.style.overflow = '';
		document.removeEventListener('touchmove', preventBodyTouchMove);
		// Flush iOS Safari's touch hit-test cache after unlocking.
		void document.body.getBoundingClientRect();
		scrollLocked = false;
	}

	$effect(() => {
		if (!browser) return;
		if (isOpen) {
			clearTimeout(closeTimer);
			renderChildren = true;
			if (lockBodyScroll) lockScroll();
		} else if (scrollLocked || renderChildren) {
			closeTimer = setTimeout(() => {
				if (lockBodyScroll) unlockScroll();
				renderChildren = false;
			}, TRANSITION_MS);
		}
	});

	onDestroy(() => {
		clearTimeout(closeTimer);
		if (!browser) return;
		if (scrollLocked) unlockScroll();
	});

	const onSwipe = (direction: 'top' | 'bottom') => {
		if (direction === 'bottom') isOpen = false;
	};

	const closeSwipeAction = createVerticalSwipeHandler(onSwipe);
</script>

<!--
	Overlay and drawer are always in the DOM (not inside {#if}).
	CSS transitions handle animation. pointer-events:none when closed ensures
	touches are never dispatched to the hidden elements — removing a DOM node
	mid-touch-sequence on iOS Safari corrupts the entire touch-dispatch state
	until the page is reloaded.
-->
<div
	class="overlay"
	class:visible={isOpen}
	onclick={() => (isOpen = false)}
	aria-hidden="true"
></div>
<div
	class="drawer"
	bind:this={drawerEl}
	class:visible={isOpen}
	use:closeSwipeAction.action
	role="dialog"
	aria-modal={isOpen}
	aria-hidden={!isOpen}
	tabindex="-1"
	onkeydown={handleKeydown}
>
	{#if renderChildren}
		{@render children?.()}
	{/if}
</div>

<style>
	.overlay {
		position: fixed;
		top: 0;
		left: 0;
		height: 100%;
		width: 100vw;
		background: rgba(0, 0, 0, 0.5);
		opacity: 0;
		pointer-events: none;
		z-index: 9998;
		transition: opacity 300ms;
	}
	.overlay.visible {
		opacity: 1;
		pointer-events: auto;
	}

	.drawer {
		--border-top-radius: 1rem;
		position: fixed;
		bottom: 0;
		left: 0;
		width: 100%;
		max-height: 90vh;
		overflow-y: auto;
		padding: 1rem;
		box-shadow: 0 -4px 20px rgba(0, 0, 0, 0.15);
		border-top-left-radius: var(--border-top-radius);
		border-top-right-radius: var(--border-top-radius);
		z-index: 9999;
		background: hsl(var(--color-background));
		outline: none;
		pointer-events: none;
		transform: translateY(100%);
		transition: transform 300ms;
	}
	.drawer.visible {
		pointer-events: auto;
		transform: translateY(0);
	}
</style>
