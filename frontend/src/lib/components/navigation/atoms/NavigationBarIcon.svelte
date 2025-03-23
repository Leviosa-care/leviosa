<script lang="ts">
	import type { NavState, NavigationBarElement } from '$lib/types';
	import { Icon } from 'lucide-svelte';
	import { navstate } from '$lib/stores/navbar';

	function setState(event: MouseEvent) {
		let targetElement = event.currentTarget as HTMLButtonElement;
		let id = targetElement.id as NavState;
		navstate.set(id);
	}

	// TODO: if I want to fill the icons, use that in the svelte component
	// fill={active ? "hsl(var(--clr-accent))" : "none"}

	interface Props {
		active: boolean;
		icon: NavigationBarElement;
		hideLabel: boolean;
	}

	let { active = false, icon, hideLabel = false }: Props = $props();
</script>

<button
	id={icon.label}
	class:activeButton={active}
	class:hidden={icon.label === 'ghost'}
	class:hideLabel
	class="icon"
	style:display={hideLabel ? 'grid' : 'initial'}
	onclick={setState}
>
	<a class="flex-plus flex" href={icon.href}>
		<Icon
			strokeWidth={1.5}
			absoluteStrokeWidth={true}
			style="width: var(--fs-2); height: var(--fs-2);"
		/>
		<p
			class="label capitalize"
			style:display={hideLabel ? 'none' : 'initial'}
			style:visibility={hideLabel ? 'hidden' : 'visible'}
		>
			{icon.label}
		</p>
	</a>
</button>

<style>
	:root {
		--icon-size: clamp(1.5rem, 1vw, 5rem);
	}
	.icon {
		background: transparent;
		color: hsl(var(--clr-dark-primary));
		display: grid;
		place-content: center;
		border-radius: 0.5rem;
		width: fit-content;
		/* opacity: 0.4; */

		flex-direction: column;
		align-items: center;
		gap: 0.1rem;
	}
	a {
		color: inherit;
		text-decoration: none;
	}
	.flex-plus {
		flex-direction: column;
		align-items: center;
		gap: 0.1rem;
	}
	.activeButton {
		/* TODO: the new color accent for that thing ? */
		/* color: #0c51c4; */
		color: hsl(var(--clr-accent));
		font-weight: 500;
	}
	.hidden {
		visibility: hidden;
	}
	.label {
		font-size: var(--fs--2);
	}

	@media only screen and (min-width: 500px) {
		.icon {
			padding: 0.25rem;
			width: 100%;
			border-radius: 1rem;
		}
		.icon:is(:global(:hover, :focus)) {
			background-color: #f7f7f9;
		}
		.activeButton {
			background-color: #f7f7f9;
			box-shadow: rgba(0, 0, 0, 0.05) 0px 0px 0px 1px;
		}
	}

	@media only screen and (min-width: 1280px) {
		/* do not center the content of the button */
		.icon {
			/* display: initial; */
			padding: 1rem;
		}
		.label {
			font-size: var(--fs-1);
			font-size: var(--fs-0);
		}
		.hideLabel .label {
			display: none;
			visibility: hidden;
		}
		.flex-plus {
			flex-direction: row;
			gap: 1rem;
		}
	}
	/* the next media has the label next to each other brother */
</style>
