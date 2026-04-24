<script lang="ts">
	import type { PageProps } from "./$types";
	import { Tabs } from "bits-ui";
	import {
		Package,
		Tag,
		Plus,
		Search,
		Filter,
		MoreVertical,
		Pencil,
		Trash2,
		Eye,
		EyeOff,
	} from "@lucide/svelte";
	import Product from "./Product.svelte";
	import Category from "./Category.svelte";

	let { data }: PageProps = $props();
	let {
		cards,
		statuses,
		categories,
		availabilities,
		deleteProductForm,
		createProductForm,
		updateProductForm,
		createCategoryForm,
	} = data;

	let activeTab = $state("products");

	let searchQuery = $state("");
	let statusFilter = $state("all");
	let categoryFilter = $state("all");
	let availabilityFilter = $state("all");

	let filteredProducts = $derived(
		cards.filter(
			(p) =>
				(statusFilter === "all" || p.published === statusFilter) &&
				(categoryFilter === "all" || p.category === categoryFilter) &&
				(availabilityFilter === "all" || p.availability === availabilityFilter) &&
				p.name.toLowerCase().includes(searchQuery.toLowerCase())
		)
	);

	let filteredCategories = $derived(
		categories.filter(
			(c) =>
				c.id !== "default" &&
				c.name.toLowerCase().includes(searchQuery.toLowerCase())
		)
	);
</script>

<div class="flex flex-col h-full bg-background">
	<!-- Header -->
	<div class="border-b border-border-card px-6 py-5">
		<div class="flex flex-col gap-1">
			<h1 class="text-2xl font-semibold tracking-tight text-foreground">
				Catalogue
			</h1>
			<p class="text-sm text-foreground-alt">
				Gérez vos produits et catégories disponibles à la réservation
			</p>
		</div>
	</div>

	<!-- Tabs Navigation -->
	<div class="border-b border-border-card px-6">
		<Tabs.Root bind:value={activeTab} class="w-full">
			<Tabs.List
				class="inline-flex gap-1 bg-muted/30 rounded-lg p-1 -mb-px"
			>
				<Tabs.Trigger
					value="products"
					class="inline-flex items-center gap-2 px-4 py-2 text-sm font-medium rounded-md transition-all data-[state=active]:bg-background data-[state=active]:text-foreground data-[state=active]:shadow-sm data-[state=inactive]:text-foreground-alt"
				>
					<Package size={16} />
					<span>Produits</span>
					<span class="text-xs bg-muted-foreground/20 px-1.5 py-0.5 rounded">
						{cards.length}
					</span>
				</Tabs.Trigger>
				<Tabs.Trigger
					value="categories"
					class="inline-flex items-center gap-2 px-4 py-2 text-sm font-medium rounded-md transition-all data-[state=active]:bg-background data-[state=active]:text-foreground data-[state=active]:shadow-sm data-[state=inactive]:text-foreground-alt"
				>
					<Tag size={16} />
					<span>Catégories</span>
					<span class="text-xs bg-muted-foreground/20 px-1.5 py-0.5 rounded">
						{categories.length - 1}
					</span>
				</Tabs.Trigger>
			</Tabs.List>

			<!-- Products Tab -->
			<Tabs.Content value="products" class="mt-6">
				<Product
					{cards}
					{statuses}
					{categories}
					{availabilities}
					{deleteProductForm}
					{createProductForm}
					{updateProductForm}
				/>
			</Tabs.Content>

			<!-- Categories Tab -->
			<Tabs.Content value="categories" class="mt-6">
				<Category {createCategoryForm} />
			</Tabs.Content>
		</Tabs.Root>
	</div>
</div>
