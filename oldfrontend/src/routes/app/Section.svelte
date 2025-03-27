<script lang="ts">
	import type { Snippet } from 'svelte';
	interface Props {
		title: string;
		cta?: string;
		children: Snippet;
	}

	let { title, cta = '', children }: Props = $props();
</script>

<div class="container grid" style="--gap: 1rem;">
	<div class="section flex">
		<p class="section__title">{title}</p>
		{#if cta !== ''}
			<a href="/app/services" class="cta">{cta}</a>
		{/if}
	</div>
	{@render children()}
</div>

<style>
	.section {
		align-items: center;
		justify-content: space-between;
	}
	.section__title {
		font-size: var(--fs-0);
		font-size: 1.1rem;
		color: hsl(var(--clr-grey-600));
		font-weight: 600;
	}
	.section__title::first-letter {
		text-transform: capitalize;
	}
	.cta {
		font-size: var(--fs--1);
		color: inherit;
		text-decoration: underline;
	}
</style>
