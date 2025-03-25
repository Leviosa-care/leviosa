<script lang="ts">
	import { ROLES, type ConsultationTabs } from '$lib/types';
	import ConsultationNavigationBarIcon from './atoms/ConsultationNavigationBarIcon.svelte';
	import { consultationState } from '$lib/stores/persisted_stores.svelte';

	interface Props {
		role: import('$lib/types').Role;
	}
	let { role }: Props = $props();

	let tabs: ConsultationTabs = {
		[ROLES.Unknown]: [],
		[ROLES.Anonymous]: [],
		[ROLES.Basic]: [],
		[ROLES.Premium]: [
			{ name: 'Consultations a venir', href: '/app/reservations' },
			{
				name: 'Reserve ta consultation',
				href: '/app/reservations/consultations/book'
			}
		],
		[ROLES.Freelance]: [],
		[ROLES.Guest]: [],
		[ROLES.Admin]: []
	};

	// {
	//     name: "Creer un type de consulation",
	//     href: "/app/reservations/events/new",
	// },
</script>

<nav class="navigation-bar snaps-inline container grid" style="--gap: 1.5rem;">
	{#each tabs[role] as { name, href }}
		<div>
			<ConsultationNavigationBarIcon {name} {href} active={$consultationState === name} />
		</div>
	{/each}
</nav>

<style>
	.navigation-bar {
		grid-auto-flow: column;
		/* grid-auto-columns: 16%; */
		gap: 2rem;
		border-bottom: 1px solid hsl(var(--clr-stroke));
		text-wrap: nowrap;

		overflow-x: auto;
		scrollbar-width: none;
		overscroll-behavior-inline: contain;
	}
	.snaps-inline {
		scroll-snap-type: inline mandatory;
		scroll-padding-inline: 1.5rem;
	}
	.snaps-inline > * {
		scroll-snap-align: start;
	}
</style>
