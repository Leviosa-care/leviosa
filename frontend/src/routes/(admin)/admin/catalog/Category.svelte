<script lang="ts">
	import { Plus, Tag, Pencil, Trash2, Search, AlertTriangle } from "@lucide/svelte";
	import type { SuperValidated } from "sveltekit-superforms";
	import type { category } from "./schemas";
	import CategoryModal from "./CategoryModal.svelte";
	import type { Category } from "./products";

	type Props = {
		createCategoryForm: SuperValidated<category>;
		categories: Category[];
	};

	let { createCategoryForm, categories }: Props = $props();

	let searchQuery = $state("");
	let statusFilter = $state("all");
	let editDialogOpen = $state(false);
	let deleteDialogOpen = $state(false);
	let selectedCategory = $state<Category | null>(null);

	// Form state for edit
	let editName = $state("");
	let editDescription = $state("");
	let editStatus = $state<"published" | "draft" | "archived">("published");

	let filteredCategories = $derived(
		categories.filter(
			(c) =>
				c.id !== "default" &&
				(statusFilter === "all" || c.status === statusFilter) &&
				c.name.toLowerCase().includes(searchQuery.toLowerCase())
		)
	);

	function getStatusBadge(status: string) {
		switch (status) {
			case "published":
				return "bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-400";
			case "draft":
				return "bg-yellow-100 text-yellow-700 dark:bg-yellow-900/30 dark:text-yellow-400";
			case "archived":
				return "bg-gray-100 text-gray-700 dark:bg-gray-800 dark:text-gray-400";
			default:
				return "bg-gray-100 text-gray-700 dark:bg-gray-800 dark:text-gray-400";
		}
	}

	function openEditDialog(category: Category) {
		selectedCategory = category;
		editName = category.name;
		editDescription = category.description || "";
		editStatus = category.status || "draft";
		editDialogOpen = true;
	}

	function closeEditDialog() {
		editDialogOpen = false;
		selectedCategory = null;
	}

	function openDeleteDialog(category: Category) {
		selectedCategory = category;
		deleteDialogOpen = true;
	}

	function closeDeleteDialog() {
		deleteDialogOpen = false;
		selectedCategory = null;
	}

	function handleSave() {
		// TODO: Connect to backend action
		console.log("Saving category:", {
			id: selectedCategory?.id,
			name: editName,
			description: editDescription,
			status: editStatus
		});
		closeEditDialog();
	}

	function handleDelete() {
		// TODO: Connect to backend action
		console.log("Deleting category:", selectedCategory?.id);
		closeDeleteDialog();
	}
</script>

<div class="flex flex-col gap-6">
	<!-- Header with actions -->
	<div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
		<div>
			<h2 class="text-lg font-semibold text-foreground">Catégories</h2>
			<p class="text-sm text-foreground-alt">
				{categories.filter((c) => c.id !== "default").length} catégorie{categories.filter((c) => c.id !== "default").length > 1
					? 's'
					: ''} au total
			</p>
		</div>

		<CategoryModal modalForm={createCategoryForm}>
			<button
				type="button"
				class="inline-flex items-center gap-2 px-4 py-2 text-sm font-medium text-background bg-foreground rounded-lg hover:opacity-90 transition-colors"
			>
				<Plus size={16} />
				<span>Nouvelle catégorie</span>
			</button>
		</CategoryModal>
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
				bind:value={searchQuery}
				placeholder="Rechercher une catégorie..."
				class="w-full pl-10 pr-4 py-2 text-sm border border-border-card rounded-lg bg-background focus:outline-none focus:ring-2 focus:ring-primary/20 focus:border-primary"
			/>
		</div>

		<select
			bind:value={statusFilter}
			class="px-3 py-2 text-sm border border-border-card rounded-lg bg-background focus:outline-none focus:ring-2 focus:ring-primary/20 focus:border-primary"
		>
			<option value="all">Tous les statuts</option>
			<option value="published">Publié</option>
			<option value="draft">Brouillon</option>
			<option value="archived">Archivé</option>
		</select>
	</div>

	<!-- Categories list -->
	{#if filteredCategories.length === 0}
		<div class="flex flex-col items-center justify-center py-12 text-center">
			<div
				class="w-16 h-16 rounded-full bg-muted flex items-center justify-center mb-4"
			>
				<Tag size={32} class="text-muted-foreground" />
			</div>
			<h3 class="text-lg font-medium text-foreground mb-1">
				Aucune catégorie trouvée
			</h3>
			<p class="text-sm text-foreground-alt">
				{searchQuery || statusFilter !== "all"
					? "Essayez de modifier vos filtres de recherche"
					: "Commencez par créer votre première catégorie"}
			</p>
		</div>
	{:else}
		<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
			{#each filteredCategories as cat (cat.id)}
				<div
					class="bg-card border border-border-card rounded-lg p-5 hover:shadow-md transition-shadow"
				>
					<div class="flex items-start justify-between mb-3">
						<div class="flex items-center gap-3">
							<div
								class="w-10 h-10 rounded-lg bg-primary/10 flex items-center justify-center"
							>
								<Tag size={20} class="text-primary" />
							</div>
							<div>
								<h3 class="font-semibold text-foreground">
									{cat.name}
								</h3>
								<span
									class="inline-flex items-center px-2 py-0.5 text-xs font-medium rounded-full mt-1 {getStatusBadge(
										cat.status || "draft"
									)}"
								>
									{cat.status || "draft"}
								</span>
							</div>
						</div>
					</div>
					{#if cat.description}
						<p class="text-sm text-foreground-alt line-clamp-2 mb-3">
							{cat.description}
						</p>
					{/if}
					<div class="flex items-center justify-between text-sm">
						<span class="text-foreground-alt">
							{cat.productCount || 0} produit
							{cat.productCount === 1 ? '' : 's'}
						</span>
						<div class="flex gap-1">
							<button
								onclick={() => openEditDialog(cat)}
								class="p-1.5 text-muted-foreground hover:text-primary hover:bg-primary/10 rounded-lg transition-colors"
								aria-label="Modifier"
							>
								<Pencil size={14} />
							</button>
							<button
								onclick={() => openDeleteDialog(cat)}
								class="p-1.5 text-muted-foreground hover:text-red-600 hover:bg-red-50 dark:hover:bg-red-950 rounded-lg transition-colors"
								aria-label="Supprimer"
							>
								<Trash2 size={14} />
							</button>
						</div>
					</div>
				</div>
			{/each}
		</div>
	{/if}
</div>

<!-- Edit Category Modal -->
{#if selectedCategory && editDialogOpen}
	<div
		class="fixed inset-0 z-50 flex items-center justify-center bg-black/50"
		onclick={closeEditDialog}
	>
		<div
			class="bg-background border border-border-card rounded-lg shadow-lg w-full max-w-md mx-4 overflow-hidden"
			onclick={(e) => e.stopPropagation()}
		>
			<div class="px-6 py-4 border-b border-border-card">
				<h3 class="text-lg font-semibold text-foreground">
					Modifier la Catégorie
				</h3>
				<p class="text-sm text-foreground-alt mt-1">
					{selectedCategory.name}
				</p>
			</div>
			<div class="p-6 grid gap-4">
				<div>
					<label
						for="edit-name"
						class="block text-sm font-medium text-foreground-alt mb-1"
					>
						Nom
					</label>
					<input
						id="edit-name"
						type="text"
						bind:value={editName}
						class="w-full px-4 py-2.5 border border-border-input rounded-lg focus:outline-none focus:ring-2 focus:ring-foreground focus:border-transparent text-sm"
					/>
				</div>
				<div>
					<label
						for="edit-description"
						class="block text-sm font-medium text-foreground-alt mb-1"
					>
						Description
					</label>
					<textarea
						id="edit-description"
						bind:value={editDescription}
						rows="3"
						class="w-full px-4 py-2.5 border border-border-input rounded-lg focus:outline-none focus:ring-2 focus:ring-foreground focus:border-transparent text-sm resize-none"
					></textarea>
				</div>
				<div>
					<label
						for="edit-status"
						class="block text-sm font-medium text-foreground-alt mb-1"
					>
						Statut
					</label>
					<select
						id="edit-status"
						bind:value={editStatus}
						class="w-full px-4 py-2.5 border border-border-input rounded-lg focus:outline-none focus:ring-2 focus:ring-foreground focus:border-transparent text-sm"
					>
						<option value="published">Publié</option>
						<option value="draft">Brouillon</option>
						<option value="archived">Archivé</option>
					</select>
				</div>
			</div>
			<div class="px-6 py-4 bg-muted/30 border-t border-border-card flex justify-end gap-3">
				<button
					type="button"
					onclick={closeEditDialog}
					class="px-6 py-2.5 border border-border-input text-foreground-alt rounded-lg hover:bg-muted transition-colors font-medium"
				>
					Annuler
				</button>
				<button
					type="button"
					onclick={handleSave}
					class="px-6 py-2.5 bg-foreground text-background rounded-lg hover:opacity-90 transition-colors font-medium"
				>
					Enregistrer
				</button>
			</div>
		</div>
	</div>
{/if}

<!-- Delete Category Modal -->
{#if selectedCategory && deleteDialogOpen}
	<div
		class="fixed inset-0 z-50 flex items-center justify-center bg-black/50"
		onclick={closeDeleteDialog}
	>
		<div
			class="bg-background border border-border-card rounded-lg shadow-lg w-full max-w-md mx-4 overflow-hidden"
			onclick={(e) => e.stopPropagation()}
		>
			<div class="p-6">
				<div class="flex items-center gap-4 mb-4">
					<div
						class="w-12 h-12 rounded-full bg-red-100 flex items-center justify-center flex-shrink-0"
					>
						<AlertTriangle class="text-red-600" size={24} />
					</div>
					<div>
						<h3 class="text-lg font-semibold text-foreground">
							Supprimer la Catégorie
						</h3>
						<p class="text-sm text-foreground-alt mt-1">
							Êtes-vous sûr de vouloir supprimer "{selectedCategory.name}" ? Cette
							action ne peut pas être annulée.
						</p>
					</div>
				</div>
				<div class="flex justify-end gap-3 pt-2">
					<button
						type="button"
						onclick={closeDeleteDialog}
						class="px-6 py-2.5 border border-border-input text-foreground-alt rounded-lg hover:bg-muted transition-colors font-medium"
					>
						Annuler
					</button>
					<button
						type="button"
						onclick={handleDelete}
						class="px-6 py-2.5 bg-red-600 text-white rounded-lg hover:bg-red-700 transition-colors font-medium"
					>
						Supprimer
					</button>
				</div>
			</div>
		</div>
	</div>
{/if}
