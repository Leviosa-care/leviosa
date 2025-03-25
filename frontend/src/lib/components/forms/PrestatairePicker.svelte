<script lang="ts">
	import { User, ChevronUp, ChevronDown, Search } from 'lucide-svelte';
	import Drawer from '$lib/components/Drawer.svelte';

	import type { Freelancer } from '$lib/types';

	let isOpen = $state(false);

	const chevronSize = 12;

	interface Props {
		prestataires: Freelancer[];
		value?: string | undefined;
	}

	let { prestataires, value = $bindable('Selectionne un(e) prestataire') }: Props = $props();

	function handlePrestataireClick(e: MouseEvent) {
		const target = e.currentTarget as HTMLButtonElement;
		const person = prestataires.find((person) => person.id === target.id);
		value = `Prestataire: ${person?.firstname} ${person?.lastname}`;
		isOpen = false;
	}
	// TODO: make the value of the input change what is display on data or a copy of data that I can use after filtering data
</script>

<button type="button" class="toggler stroke" onclick={() => (isOpen = !isOpen)}>
	<div class="flex" style="justify-content: space-between;">
		<div class="icon">
			<User strokeWidth={1.25} aria-label="Icone utilisateur" />
		</div>
		<input type="button" class="toggle-text" {value} />
		<div class="selector flex">
			<ChevronUp size={chevronSize} />
			<ChevronDown size={chevronSize} />
		</div>
	</div>
</button>
<Drawer bind:isOpen>
	<div class="holder flex">
		{#if prestataires.length > 12}
			<div class="drawer-header grid">
				<div class="container">
					<div class="search container flex">
						<Search size={16} color="hsl(var(--clr-dark-primary))" aria-label="Icone recherche" />
						<input type="text" placeholder="Nom, Prenom..." />
					</div>
				</div>
				<div class="separator"></div>
			</div>
		{/if}
		<div class="cards container grid">
			{#each prestataires as { id, avatar, firstname, lastname }}
				<button {id} type="button" class="card" onclick={handlePrestataireClick}>
					<img class="avatar" src={avatar} alt="avatar prestataire" />
					<p class="bg-blured">
						{firstname}
						{lastname}
					</p>
				</button>
			{/each}
		</div>
	</div>
</Drawer>

<style>
	.toggler {
		width: 100%;
		border-radius: 0.5rem;
		position: relative;
		background-color: #f7f7f9;
	}
	.selector {
		position: absolute;
		right: 0;
		top: 0;
		width: 2.5rem;
		height: 100%;
		border-radius: 0.5rem;
		pointer-events: none;
		display: flex;
		align-items: center;
		justify-content: center;
		flex-direction: column;
		gap: 0rem;
	}
	.icon {
		position: absolute;
		left: 0;
		top: 0;
		height: 100%;
		width: 2.5rem;
		display: flex;
		align-items: center;
		justify-content: center;
	}
	.toggle-text {
		border-radius: 0.5rem;
		border: none;
		padding-block: calc(3 * var(--fs-0) / 4);
		padding-inline: 3rem;
		background: transparent;
		appearance: none;
		width: 100%;
	}
	.holder {
		width: calc(100vw - 1.5rem);
		max-height: 85vh;
		position: absolute;
		bottom: 0;
		left: 0;
		right: 0;
		width: 100vw;
		overflow-x: hidden;
		background-color: hsl(var(--clr-light-primary));
		margin-inline: auto;
		border-radius: 0.5rem;
		padding-block: 1rem;
		flex-direction: column;
	}
	.drawer-header {
		flex: none;
	}
	.search {
		background-color: hsl(var(--clr-light-secondary));
		height: 40px;
		width: 100%;
		align-items: center;
		padding: 0.5rem;
		border-radius: 0.5rem;
	}
	.search input {
		background: transparent;
		width: 100%;
	}
	.search input:is(:global(:focus, :hover)) {
		outline: none;
	}
	.cards {
		grid-template-columns: repeat(3, 1fr);
		overflow-y: auto;
		flex: 1;
		padding-bottom: 1rem;
	}
	.card {
		position: relative;
		color: hsl(var(--clr-light-primary));
	}
	.avatar {
		width: 100px;
		aspect-ratio: 1;
		border-radius: 0.5rem;
	}
	.bg-blured {
		position: absolute;
		font-size: 0.8125rem;
		bottom: 0;
		height: 40%;
		width: 100%;
		backdrop-filter: blur(15px);
		display: grid;
		place-content: center;
		padding: 0.5rem;
		text-align: center;
		transition: height 0.3s ease;
	}
	.separator {
		width: 100%;
		height: 1px;
		background-color: hsl(var(--clr-stroke));
	}
</style>
