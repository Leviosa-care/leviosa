<script lang="ts">
	import MessageNavigationBarIcon from './atoms/MessageNavigationBarIcon.svelte';
	import { messageState } from '$lib/stores/persisted_stores.svelte';

	import type { MessageTabs } from '$lib/types';
	import { role } from '$lib/data';
	import { ROLES } from '$lib/types';

	let tabs: MessageTabs = {
		[ROLES.Unknown]: [],
		[ROLES.Anonymous]: [],
		[ROLES.Basic]: [],
		[ROLES.Premium]: [
			{ name: 'Conversations', href: '/app/reservations' },
			{ name: 'Notes de séances', href: '/app/reservations' }
		],
		[ROLES.Freelance]: [],
		[ROLES.Guest]: [],
		[ROLES.Admin]: []
	};
</script>

<nav class="navigation-bar snaps-inline grid" style="--gap: 1.5rem;">
	{#each tabs[role] as { name }}
		<div>
			<MessageNavigationBarIcon {name} active={$messageState === name} />
		</div>
	{/each}
</nav>

<style>
	.navigation-bar {
		gap: 2rem;
		border-bottom: 1px solid hsl(var(--clr-stroke));
		/* text-wrap: nowrap; */
		grid-template-columns: repeat(2, 1fr);

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
