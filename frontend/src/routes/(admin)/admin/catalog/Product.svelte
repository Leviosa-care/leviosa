<script lang="ts">
	import { Plus, Search, Package } from "@lucide/svelte";
	import ProductCard from "./ProductCard.svelte";
	import { type CardType, type Category } from "$lib/data/mockData";
	import { type SuperValidated } from "sveltekit-superforms";
	import type { DeleteProduct, product } from "./schemas";
	import {
		defaultStatus,
		defaultCategory,
		defaultAvailability,
	} from "./default";

	type Props = {
		cards: CardType[];
		statuses: Set<string>;
		categories: Category[];
		availabilities: Set<string>;
		deleteProductForm: SuperValidated<DeleteProduct>;
		updateProductForm: SuperValidated<product>;
	};

	let {
		cards,
		statuses,
		categories,
		availabilities,
		deleteProductForm,
		updateProductForm,
	}: Props = $props();

	// filters
	let status = $state(defaultStatus);
	let category = $state(defaultCategory);
	let availability = $state(defaultAvailability);
	let searchValue = $state("");

	let filteredCards = $derived(
		cards
			.filter((card) => status === defaultStatus || card.published === status)
			.filter(
				(card) =>
					category === defaultCategory || card.category === category
			)
			.filter(
				(card) =>
					availability === defaultAvailability ||
					card.published === availability
			)
			.filter((card) =>
				card.name.toLowerCase().includes(searchValue.toLowerCase())
			)
	);
</script>

<div class="flex flex-col gap-6">
	<!-- Header with actions -->
	<div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
		<div>
			<h2 class="text-lg font-semibold text-foreground">Produits</h2>
			<p class="text-sm text-foreground-alt">
				{cards.length} produit{cards.length > 1 ? 's' : ''} au total
			</p>
		</div>

		<a
			href="/admin/catalog/products/new"
			class="inline-flex items-center gap-2 px-4 py-2 text-sm font-medium text-background bg-foreground rounded-lg hover:opacity-90 transition-colors"
		>
			<Plus size={16} />
			<span>Nouveau produit</span>
		</a>
	</div>

	<!-- Filters and search -->
	<div class="flex flex-col sm:flex-row gap-3">
		<div class="relative flex-1">
			<Search
				class="absolute left-3 top-1/2 -translate-y-1/2 text-muted-foreground"
				size={16}
			/>
			<input
				type="text"
				bind:value={searchValue}
				placeholder="Rechercher un produit..."
				class="w-full pl-10 pr-4 py-2 text-sm border border-border-card rounded-lg bg-background focus:outline-none focus:ring-2 focus:ring-primary/20 focus:border-primary"
			/>
		</div>

		<select
			bind:value={status}
			class="px-3 py-2 text-sm border border-border-card rounded-lg bg-background focus:outline-none focus:ring-2 focus:ring-primary/20 focus:border-primary"
		>
			<option value={defaultStatus}>Tous les statuts</option>
			{#each Array.from(statuses).filter(s => s !== defaultStatus) as s}
				<option value={s}>{s}</option>
			{/each}
		</select>

		<select
			bind:value={category}
			class="px-3 py-2 text-sm border border-border-card rounded-lg bg-background focus:outline-none focus:ring-2 focus:ring-primary/20 focus:border-primary"
		>
			<option value={defaultCategory}>Toutes les catégories</option>
			{#each categories.filter(c => c.id !== 'default') as c}
				<option value={c.name}>{c.name}</option>
			{/each}
		</select>

		<select
			bind:value={availability}
			class="px-3 py-2 text-sm border border-border-card rounded-lg bg-background focus:outline-none focus:ring-2 focus:ring-primary/20 focus:border-primary"
		>
			<option value={defaultAvailability}>Tous les types</option>
			{#each Array.from(availabilities).filter(a => a !== defaultAvailability) as a}
				<option value={a}>{a}</option>
			{/each}
		</select>
	</div>

	<!-- Products grid -->
	{#if filteredCards.length === 0}
		<div class="flex flex-col items-center justify-center py-12 text-center">
			<div
				class="w-16 h-16 rounded-full bg-muted flex items-center justify-center mb-4"
			>
				<Package size={32} class="text-muted-foreground" />
			</div>
			<h3 class="text-lg font-medium text-foreground mb-1">
				Aucun produit trouvé
			</h3>
			<p class="text-sm text-foreground-alt">
				{searchValue ||
				status !== defaultStatus ||
				category !== defaultCategory ||
				availability !== defaultAvailability
					? "Essayez de modifier vos filtres de recherche"
					: "Commencez par créer votre premier produit"}
			</p>
		</div>
	{:else}
		<div
			class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4"
		>
			{#each filteredCards as card (card.id)}
				<ProductCard
					{card}
					{statuses}
					{categories}
					{availabilities}
					{deleteProductForm}
					{updateProductForm}
				/>
			{/each}
		</div>
	{/if}
</div>
