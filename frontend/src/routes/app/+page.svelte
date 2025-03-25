<script lang="ts">
	import type { PageData } from './$types';

	import Drawer from '$lib/components/Drawer.svelte';
	import NoEventCard from '$lib/components/ui/NoEventCard.svelte';
	import EventCard from '$lib/components/ui/EventCard.svelte';
	import InstallHeader from './InstallHeader.svelte';
	import Header from './Header.svelte';
	import Section from './Section.svelte';
	import ServiceCarousel from './ServiceCarousel.svelte';
	import QRCode from './QRCode.svelte';

	import { NAV_STATES } from '$lib/types';
	import { navigationState } from '$lib/stores/persisted_stores.svelte';
	navigationState.set(NAV_STATES.Accueil); // just to forget the value stored in localstore when reconecting and I had the page to another link.

	interface Props {
		data: PageData;
	}
	let { data }: Props = $props();
	const { name, qrcode } = data;

	let isOpen = $state(false);
</script>

<InstallHeader />
<div class="content flow relative" style="--flow-space: 3rem;">
	<Header {name} toggleDrawer={() => (isOpen = !isOpen)} />
	<Section title="decouvrez nos services" cta="Voir tout">
		<ServiceCarousel />
	</Section>
	<Section title="votre prochain evenement">
		<EventCard />
	</Section>
	<Section title="votre prochain evenement">
		<NoEventCard />
	</Section>
	<div class="next-event container grid" style="--gap: 1rem;">
		<h3 class="fs-h3 subtitle">Tu peux aussi revoir...</h3>
		<p>Jean Dupont, ton dernier prestataire massage</p>
	</div>
</div>
<Drawer bind:isOpen>
	<QRCode {qrcode} />
</Drawer>

<style>
	.content {
		padding-bottom: 7rem;
	}
	.subtitle {
		color: black;
		font-size: var(--fs-1);
		font-weight: 600;
	}
</style>
