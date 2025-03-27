<script lang="ts">
	import { goto } from '$app/navigation';
	import type { PageData } from './$types';

	import { Contrast, LogOut, Instagram } from 'lucide-svelte';
	import Drawer from '$lib/components/Drawer.svelte';
	import Button from '$lib/components/Button.svelte';

	// TODO: add the part for when a user access that part of the page without having an account
	import ProfileCard from './ProfileCard.svelte';
	import Field from './Field.svelte';
	import AppearanceDrawer from './AppearanceDrawer.svelte';
	import Socials from './Socials.svelte';
	import Google from '../../../assets/google.svelte';

	interface Props {
		data: PageData;
	}

	let { data }: Props = $props();
	let { perso, juridique, parameters } = data;

	let isOpen = $state(false);

	function disconnect() {
		console.log('just delete the token that allows authentication');
	}
</script>

<div class="content grid" style="--gap: 2rem;">
	<div class="container">
		<div class="field container flex">
			<ProfileCard />
		</div>
	</div>
	<div class="container grid" style="--gap: 3rem;">
		<div class="grid" style="--gap: 0.75rem;">
			<p class="title">Informations personnelles</p>
			<div>
				{#each perso as { name, icon, pathname }}
					<Field {name} {icon} action={() => goto(pathname)} />
				{/each}
			</div>
		</div>
		<div class="grid" style="--gap: 0.75rem;">
			<p class="title">Parametres</p>
			<div>
				<Field name="Apparence" icon={Contrast} action={() => (isOpen = !isOpen)} />
				{#each parameters as { name, icon, pathname }}
					<Field {name} {icon} action={() => console.log('hello sir')} />
				{/each}
			</div>
		</div>
		<div class="grid" style="--gap: 0.75rem;">
			<p class="title">Documents legaux</p>
			<div>
				{#each juridique as { name, icon, pathname }}
					<Field {name} {icon} action={() => console.log('hello sir')} />
				{/each}
			</div>
		</div>
		<div class="grid" style="--gap: 1rem;">
			<Button leftIcon={LogOut} style="secondary" onClick={disconnect}>Se deconnecter</Button>
			<Socials />
		</div>
	</div>
</div>
<Drawer bind:isOpen>
	<AppearanceDrawer />
</Drawer>

<style>
	.content {
		position: relative;
	}
	.title {
		font-size: var(--fs-2);
		font-weight: 600;
		color: hsl(var(--clr-grey-700));
	}
	.field {
		padding: 1rem 0.5rem;
		justify-content: space-between;
		align-items: center;
		color: hsl(var(--clr-grey-600));
		background-color: hsl(var(--clr-light-primary));
		border-bottom: 1px solid hsl(var(--clr-stroke));
	}
</style>
