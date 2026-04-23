<script lang="ts">
	import { onDestroy } from 'svelte';
	import { browser } from '$app/environment';

	interface Props {
		isOpen?: boolean;
		title?: string;
		description?: string;
		children?: import('svelte').Snippet;
	}

	let { children, isOpen = $bindable(false), title, description }: Props = $props();

	let dialogEl: HTMLElement | null = $state(null);
	let renderChildren = $state(false);
	let closeTimer: ReturnType<typeof setTimeout> | undefined;
	let scrollLocked = false;

	const TRANSITION_MS = 200;

	function handleKeydown(event: KeyboardEvent) {
		if (event.key === 'Escape') {
			isOpen = false;
		}
	}

	function preventBodyTouchMove(e: TouchEvent) {
		if (dialogEl && dialogEl.contains(e.target as Node)) return;
		e.preventDefault();
	}

	function lockScroll() {
		document.body.style.overflow = 'hidden';
		document.addEventListener('touchmove', preventBodyTouchMove, { passive: false });
		scrollLocked = true;
	}

	function unlockScroll() {
		document.body.style.overflow = '';
		document.removeEventListener('touchmove', preventBodyTouchMove);
		void document.body.getBoundingClientRect();
		scrollLocked = false;
	}

	$effect(() => {
		if (!browser) return;
		if (isOpen) {
			clearTimeout(closeTimer);
			renderChildren = true;
			lockScroll();
			if (dialogEl) dialogEl.focus();
		} else if (scrollLocked || renderChildren) {
			closeTimer = setTimeout(() => {
				unlockScroll();
				renderChildren = false;
			}, TRANSITION_MS);
		}
	});

	onDestroy(() => {
		clearTimeout(closeTimer);
		if (!browser) return;
		if (scrollLocked) unlockScroll();
	});
</script>

<div
	class="overlay"
	class:visible={isOpen}
	onclick={() => (isOpen = false)}
	aria-hidden="true"
></div>
<div
	bind:this={dialogEl}
	class="modal"
	class:visible={isOpen}
	role="dialog"
	aria-modal={isOpen}
	aria-hidden={!isOpen}
	aria-labelledby={title ? 'modal-title' : undefined}
	aria-describedby={description ? 'modal-description' : undefined}
	tabindex="-1"
	onclick={(e) => e.stopPropagation()}
	onkeydown={handleKeydown}
>
	{#if renderChildren}
		{#if title || description}
			<div class="modal-header">
				{#if title}
					<h2 id="modal-title" class="modal-title">{title}</h2>
				{/if}
				{#if description}
					<p id="modal-description" class="modal-description">{description}</p>
				{/if}
			</div>
		{/if}
		{@render children?.()}
	{/if}
</div>

<style>
	.overlay {
		position: fixed;
		inset: 0;
		background: rgba(0, 0, 0, 0.5);
		z-index: 9998;
		display: flex;
		align-items: center;
		justify-content: center;
		padding: 1rem;
		opacity: 0;
		pointer-events: none;
		transition: opacity 200ms;
	}
	.overlay.visible {
		opacity: 1;
		pointer-events: auto;
	}

	.modal {
		position: fixed;
		top: 50%;
		left: 50%;
		transform: translate(-50%, -50%) translateY(30px);
		width: 90vw;
		max-width: 540px;
		max-height: 90vh;
		overflow-y: auto;
		background: hsl(var(--color-background));
		border-radius: 0.75rem;
		box-shadow: 0 20px 25px -5px rgb(0 0 0 / 0.1), 0 8px 10px -6px rgb(0 0 0 / 0.1);
		z-index: 9999;
		padding: 1.5rem;
		outline: none;
		opacity: 0;
		pointer-events: none;
		transition: opacity 200ms, transform 250ms;
	}
	.modal.visible {
		opacity: 1;
		pointer-events: auto;
		transform: translate(-50%, -50%) translateY(0);
	}

	.modal-header {
		margin-bottom: 1rem;
		padding-bottom: 1rem;
		border-bottom: 1px solid hsl(var(--color-border-card));
	}

	.modal-title {
		font-size: 1.25rem;
		font-weight: 600;
		line-height: 1.75rem;
		color: hsl(var(--color-foreground));
		margin: 0 0 0.25rem 0;
	}

	.modal-description {
		font-size: 0.875rem;
		line-height: 1.25rem;
		color: hsl(var(--color-foreground-alt));
		margin: 0;
	}
</style>
