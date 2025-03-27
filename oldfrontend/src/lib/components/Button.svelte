<script lang="ts">
	type buttonStyle = 'primary' | 'secondary' | 'ternary';
	type ButtonType = 'submit' | 'reset' | 'button';

	import type { Component } from 'svelte';

	interface Props {
		style?: buttonStyle;
		leftIcon?: typeof import('lucide-svelte').Icon | Component;
		rightIcon?: typeof import('lucide-svelte').Icon | Component;
		buttonType?: ButtonType;
		onClick: () => void;
		children?: import('svelte').Snippet;
	}

	let {
		style = 'primary',
		leftIcon: LeftIcon,
		rightIcon: RightIcon,
		buttonType = 'button',
		onClick,
		children
	}: Props = $props();
</script>

<button
	type={buttonType}
	class="{style} flex"
	onclick={onClick}
	style="--gap: {LeftIcon !== null || RightIcon !== null ? '1rem' : '0rem'};"
>
	{@render buttonIcon(LeftIcon)}
	{@render children?.()}
	{@render buttonIcon(RightIcon)}
</button>

{#snippet buttonIcon(Icon: any)}
	{#if Icon}
		<div class="icon" style:display={Icon ? 'visible' : 'none'}>
			<Icon />
		</div>
	{/if}
{/snippet}

<style>
	button {
		border-radius: 0.5rem;
		padding: 0.75rem 1.5rem;
		align-items: center;
		justify-content: center;
		background: transparent;
		font-weight: 500;
	}
	.icon {
		margin-top: -0.1rem;
		width: 1.5rem;
		height: 1.5rem;
	}
	.primary {
		background-color: hsl(var(--clr-dark-primary));
		color: hsl(var(--clr-light-primary));
	}
	.secondary {
		outline: 1px solid hsl(var(--clr-grey-200));
		color: hsl(var(--clr-grey-700));
		background: #f7f7f9;
	}
	.ternary {
		padding: 0;
	}
</style>
