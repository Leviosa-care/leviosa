<script lang="ts">
	import type { Role, NavigationBarElement } from '$lib/types';

	import { Slack, PanelLeftClose, PanelLeftOpen } from 'lucide-svelte';
	import NavigationBarIcon from '$lib/components/navigation/atoms/NavigationBarIcon.svelte';

	import { navigationState } from '$lib/stores/persisted_stores.svelte';
	import { navigationBarIcons } from '$lib/constructor';

	let hideLabel: boolean = $state(false);

	let { role }: { role: Role } = $props();

	const smallIcons: NavigationBarElement[] = navigationBarIcons[role].small;
	const largeIcons: NavigationBarElement[] = navigationBarIcons[role].large;
</script>

{#snippet navbarIcon(icons: NavigationBarElement[], size: 'small' | 'large')}
	<div class="icons {size}">
		{#each icons as icon}
			<NavigationBarIcon {icon} active={$navigationState === icon.label} {hideLabel} />
		{/each}
		{#if size === 'large'}
			<button onclick={() => (hideLabel = !hideLabel)} class="panel-left-close">
				{#if hideLabel}
					<PanelLeftOpen
						strokeWidth={1.5}
						absoluteStrokeWidth={true}
						style="width: var(--fs-2); height: var(--fs-4);"
					/>
				{:else}
					<PanelLeftClose
						strokeWidth={1.5}
						absoluteStrokeWidth={true}
						style="width: var(--fs-2); height: var(--fs-4);"
					/>
				{/if}
			</button>
		{/if}
	</div>
{/snippet}

<div class="navigation-bar">
	<div class="logo">
		<Slack
			strokeWidth={1.5}
			absoluteStrokeWidth={true}
			style="width: var(--fs-4); height: var(--fs-4);"
		/>
	</div>
	<div class="custom-separator"></div>
	{@render navbarIcon(smallIcons, 'small')}
	{@render navbarIcon(largeIcons, 'large')}
</div>

<style>
	/* TODO: split the display flex from the position, use another class */
	.navigation-bar {
		background-color: hsl(var(--clr-light-primary));
		border-top: 1px solid hsl(var(--clr-stroke));
		min-width: 100vw;
		padding: 0.5rem 2rem 2.5rem;

		z-index: 1000;

		position: fixed;
		bottom: 0;
		left: 0;
		right: 0;
	}
	.icons {
		display: flex;
		align-items: center;
		justify-content: space-between;
	}
	.large {
		display: none;
		visibility: hidden;
	}
	.logo {
		display: none;
		visibility: hidden;
		color: hsl(var(--clr-grey-500));
		padding-inline: 1rem;
	}
	/* for mobile screen */
	@media only screen and (min-width: 500px) {
		.navigation-bar {
			min-width: 0;
			width: fit-content;
			padding: 2rem 1rem;
			border-top: 0;
			border-right: 1px solid hsl(var(--clr-stroke));
			border-right: 1px solid hsl(var(--clr-grey-100));

			position: sticky;
			height: 100vh;
			top: 0;
			flex-shrink: 0;

			display: flex;
			flex-direction: column;
			gap: 3rem;
		}
		.icons {
			flex-direction: column;
			justify-content: flex-start;
			gap: 1rem;
			flex: 1;
		}
		.small {
			display: none;
			visibility: hidden;
		}
		/* NOTE: add the large part when everything is fine brother */
		.large {
			display: flex;
			visibility: visible;
		}
		.logo {
			display: grid;
			visibility: visible;
			place-self: center;
		}
		.custom-separator {
			height: 2px;
			border-radius: 1rem;
			background-color: hsl(var(--clr-grey-100));
			/* background-color: hsl(var(--clr-stroke)); */
			max-width: 8rem;
			width: 100%;
			margin-inline: auto;
		}

		.panel-left-close {
			display: none;
			visibility: hidden;

			color: hsl(var(--clr-grey-500));
			background: transparent;

			position: absolute;
			bottom: 3rem;
			right: 1rem;
		}

		.panel-left-close {
			display: initial;
			visibility: visible;
		}
	}

	@media only screen and (min-width: 1280px) {
		.navigation-bar {
			top: 0;
			bottom: 0;
			left: 0;
		}
		.icons {
			position: relative;
		}
		.panel-left-close {
			display: initial;
			visibility: visible;
		}
	}
</style>
