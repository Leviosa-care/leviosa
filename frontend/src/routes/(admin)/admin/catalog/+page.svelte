<script lang="ts">
	import type { PageProps } from "./$types";
	import { Tabs } from "bits-ui";
	import { Package, Tag } from "@lucide/svelte";
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

	let activeTab = $state("categories");
</script>

<svelte:head>
	<title>Catalogue | Admin</title>
</svelte:head>

<div class="container mx-auto px-4 py-8 lg:py-12">
	<div class="mb-8">
		<h1 class="text-3xl lg:text-4xl font-bold mb-2 text-foreground">
			Catalogue
		</h1>
		<p class="text-muted-foreground">
			Gérez vos produits et catégories disponibles à la réservation
		</p>
	</div>

	<!-- Tabs Navigation -->
	<Tabs.Root bind:value={activeTab} class="space-y-4">
		<Tabs.List
			class="inline-flex items-center w-fit bg-transparent gap-2 text-sm font-semibold border-b border-border-card p-1"
		>
			<Tabs.Trigger
				value="categories"
				class="px-4 py-2 rounded-none bg-transparent border-b-2 data-[state=active]:shadow-none mb-[-2px] data-[state=active]:border-b-foreground data-[state=active]:text-foreground data-[state=inactive]:border-transparent data-[state=inactive]:text-muted-foreground data-[state=inactive]:hover:bg-transparent data-[state=inactive]:hover:text-foreground-alt transition-colors cursor-pointer flex items-center gap-2"
			>
				<Tag size={16} />
				<span>Catégories</span>
				<span class="text-xs bg-muted px-2 py-0.5 rounded-full">
					{categories.length - 1}
				</span>
			</Tabs.Trigger>
			<Tabs.Trigger
				value="products"
				class="px-4 py-2 rounded-none bg-transparent border-b-2 data-[state=active]:shadow-none mb-[-2px] data-[state=active]:border-b-foreground data-[state=active]:text-foreground data-[state=inactive]:border-transparent data-[state=inactive]:text-muted-foreground data-[state=inactive]:hover:bg-transparent data-[state=inactive]:hover:text-foreground-alt transition-colors cursor-pointer flex items-center gap-2"
			>
				<Package size={16} />
				<span>Produits</span>
				<span class="text-xs bg-muted px-2 py-0.5 rounded-full">
					{cards.length}
				</span>
			</Tabs.Trigger>
		</Tabs.List>

		<!-- Categories Tab -->
		<Tabs.Content value="categories" class="p-6">
			<Category {createCategoryForm} />
		</Tabs.Content>

		<!-- Products Tab -->
		<Tabs.Content value="products" class="p-6">
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
	</Tabs.Root>
</div>
