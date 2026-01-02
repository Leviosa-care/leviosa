<script lang="ts">
    import { Dialog, Button, Label, Separator, Combobox } from "bits-ui";
    import {
        Plus,
        Pencil,
        Trash2,
        X,
        Image as ImageIcon,
        Clock,
        Filter,
        Check,
        ChevronsUpDown,
    } from "@lucide/svelte";
    import type { Snippet } from "svelte";
    import type { PageData } from "./$types";
    import Drawer from "$lib/ui/Drawer.svelte";

    import { browser } from "$app/environment";
    import { mockProducts, mockCategories } from "./mockData";

    // Detect if we're on mobile
    let isMobile = $state(false);

    if (browser) {
        isMobile = window.innerWidth < 768;
        window.addEventListener("resize", () => {
            isMobile = window.innerWidth < 768;
        });
    }

    type Input = Snippet<[string]>;

    // Product type from API
    type Product = {
        id: string;
        name: string;
        description: string;
        category: string;
        duration: number;
        status: "draft" | "published";
        availability: "online" | "in-person" | "hybrid";
        bufferTime: number;
        cancellationHours: number;
        stripeProductId: string;
        metadata?: Record<string, any>;
        createdAt: string;
        updatedAt: string;
        images?: Array<{
            id: string;
            parent_id: string;
            parent_type: string;
            url: string;
            title: string;
            is_active: boolean;
            created_at: string;
        }>;
    };

    interface Props {
        data: PageData;
    }

    let { data }: Props = $props();

    // Use mock products for now
    let products = $state<Product[]>([...mockProducts]);

    // Category filter state (single select)
    const ALL_CATEGORIES = "";
    let selectedCategoryId = $state<string>(ALL_CATEGORIES);

    // Filtered products based on selected category
    const filteredProducts = $derived(
        selectedCategoryId === ALL_CATEGORIES
            ? products
            : products.filter((p) => p.category === selectedCategoryId),
    );

    // Dialog states
    let createDialogOpen = $state(false);
    let editDialogOpen = $state(false);
    let deleteDialogOpen = $state(false);

    // Currently selected product for edit/delete
    let selectedProduct: Product | null = $state(null);

    // Form states for new/edit product
    let formData = $state<{
        id: string;
        name: string;
        description: string;
        category: string;
        duration: number;
        status: string;
        availability: string;
        bufferTime: number;
        cancellationHours: number;
        stripeProductId: string;
    }>({
        id: "",
        name: "",
        description: "",
        category: "",
        duration: 60,
        status: "draft",
        availability: "in-person",
        bufferTime: 10,
        cancellationHours: 24,
        stripeProductId: "",
    });

    // Reset form to default values
    function resetForm() {
        formData = {
            id: "",
            name: "",
            description: "",
            category: "",
            duration: 60,
            status: "draft",
            availability: "in-person",
            bufferTime: 10,
            cancellationHours: 24,
            stripeProductId: "",
        };
    }

    function openCreateDialog() {
        resetForm();
        createDialogOpen = true;
    }

    function openEditDialog(product: Product) {
        selectedProduct = product;
        formData = {
            id: product.id,
            name: product.name,
            description: product.description,
            category: product.category,
            duration: product.duration,
            status: product.status,
            availability: product.availability,
            bufferTime: product.bufferTime,
            cancellationHours: product.cancellationHours,
            stripeProductId: product.stripeProductId,
        };
        editDialogOpen = true;
    }

    function openDeleteDialog(product: Product) {
        selectedProduct = product;
        deleteDialogOpen = true;
    }

    function getActiveImage(product: Product) {
        return product.images?.find((img) => img.is_active);
    }

    function getCategoryName(categoryId: string): string {
        return (
            mockCategories.find((c) => c.id === categoryId)?.name || "Inconnu"
        );
    }

    function formatDuration(minutes: number): string {
        if (minutes < 60) return `${minutes} min`;
        const hours = Math.floor(minutes / 60);
        const mins = minutes % 60;
        return mins > 0 ? `${hours}h ${mins}min` : `${hours}h`;
    }

    function handleCreateSubmit(e: Event) {
        e.preventDefault();
        const newProduct = {
            id: `prod-${Date.now()}`,
            name: formData.name,
            description: formData.description,
            category: formData.category,
            duration: formData.duration,
            status: formData.status as "draft" | "published",
            availability: formData.availability as
                | "online"
                | "in-person"
                | "hybrid",
            bufferTime: formData.bufferTime,
            cancellationHours: formData.cancellationHours,
            stripeProductId: formData.stripeProductId,
            metadata: {},
            createdAt: new Date().toISOString(),
            updatedAt: new Date().toISOString(),
            images: [],
        } satisfies Product;
        products = [...products, newProduct];
        createDialogOpen = false;
        resetForm();
    }

    function handleEditSubmit(e: Event) {
        e.preventDefault();
        if (!selectedProduct) return;

        const selectedId = selectedProduct.id;
        products = products.map((p) =>
            p.id === selectedId
                ? {
                      ...p,
                      name: formData.name,
                      description: formData.description,
                      category: formData.category,
                      duration: formData.duration,
                      status: formData.status as "draft" | "published",
                      availability: formData.availability as
                          | "online"
                          | "in-person"
                          | "hybrid",
                      bufferTime: formData.bufferTime,
                      cancellationHours: formData.cancellationHours,
                      stripeProductId: formData.stripeProductId,
                      updatedAt: new Date().toISOString(),
                  }
                : p,
        );
        editDialogOpen = false;
        selectedProduct = null;
        resetForm();
    }

    function handleDeleteSubmit(e: Event) {
        e.preventDefault();
        if (!selectedProduct) return;

        const selectedId = selectedProduct.id;
        products = products.filter((p) => p.id !== selectedId);
        deleteDialogOpen = false;
        selectedProduct = null;
    }
</script>

<div class="h-full bg-white">
    <!-- Create button - floating on mobile -->
    <Button.Root
        type="button"
        class="cursor-pointer fixed bottom-20 right-4 md:absolute md:top-6 md:right-8 z-10"
        onclick={openCreateDialog}
    >
        <div
            class="flex gap-2 items-center py-2 bg-dark text-white rounded-input hover:bg-dark/95 transition-all shadow-mini md:px-4 w-12 h-12 md:w-auto md:h-auto justify-center"
        >
            <Plus size={18} />
            <span class="text-sm font-medium hidden md:inline"
                >Nouveau Produit</span
            >
        </div>
    </Button.Root>

    <!-- Category Filter & Products Grid -->
    <div class="p-4 md:p-8">
        <!-- Category Filter Section -->
        <div
            class="flex flex-col md:flex-row md:items-center justify-between gap-4 mb-6"
        >
            <Combobox.Root type="single" bind:value={selectedCategoryId}>
                <div class="relative w-full md:w-auto md:min-w-[280px]">
                    <Filter
                        class="text-muted-foreground absolute start-3 top-1/2 size-4 -translate-y-1/2"
                    />
                    <Combobox.Input
                        class="h-input rounded-card-sm border-border-input bg-background placeholder:text-foreground-alt/50 hover:border-dark-40 focus:ring-foreground focus:ring-offset-background focus:outline-hidden inline-flex w-full items-center border px-9 text-sm focus:ring-2 focus:ring-offset-2 transition-all"
                        placeholder="Filtrer par catégorie"
                    />
                    <Combobox.Trigger
                        class="absolute end-3 top-1/2 size-4 -translate-y-1/2"
                    >
                        <ChevronsUpDown class="text-muted-foreground size-4" />
                    </Combobox.Trigger>
                </div>
                <Combobox.Portal>
                    <Combobox.Content
                        class="bg-background border border-border-input rounded-card-lg shadow-popover z-50 max-h-[300px] overflow-y-auto"
                    >
                        <Combobox.Viewport class="p-1">
                            <Combobox.Item value={ALL_CATEGORIES} label="Toutes les catégories">
                                {#snippet children({ selected })}
                                    <div
                                        class="flex items-center justify-between px-3 py-2 rounded-md hover:bg-dark-04 cursor-pointer"
                                    >
                                        <span class="text-sm"
                                            >Toutes les catégories</span
                                        >
                                        {#if selected}
                                            <Check class="size-4" />
                                        {/if}
                                    </div>
                                {/snippet}
                            </Combobox.Item>
                            <Combobox.Group class="pt-1">
                                {#each mockCategories as category}
                                    <Combobox.Item
                                        value={category.id}
                                        label={category.name}
                                    >
                                        {#snippet children({ selected })}
                                            <div
                                                class="flex items-center justify-between px-3 py-2 rounded-md hover:bg-dark-04 cursor-pointer"
                                            >
                                                <span class="text-sm"
                                                    >{category.name}</span
                                                >
                                                {#if selected}
                                                    <Check class="size-4" />
                                                {/if}
                                            </div>
                                        {/snippet}
                                    </Combobox.Item>
                                {/each}
                            </Combobox.Group>
                        </Combobox.Viewport>
                    </Combobox.Content>
                </Combobox.Portal>
            </Combobox.Root>

            <!-- Product count -->
            <div class="flex items-center gap-2 text-sm text-foreground-alt">
                {#if selectedCategoryId !== ALL_CATEGORIES}
                    <span class="hidden md:inline">•</span>
                    <span
                        >{filteredProducts.length} produit{filteredProducts.length !==
                        1
                            ? "s"
                            : ""}</span
                    >
                    <button
                        type="button"
                        onclick={() => (selectedCategoryId = ALL_CATEGORIES)}
                        class="text-destructive hover:underline"
                    >
                        Réinitialiser
                    </button>
                {:else}
                    <span
                        >{products.length} produit{products.length !== 1
                            ? "s"
                            : ""}</span
                    >
                {/if}
            </div>
        </div>

        <!-- Selected category chip -->
        {#if selectedCategoryId !== ALL_CATEGORIES}
            {@const category = mockCategories.find(
                (c) => c.id === selectedCategoryId,
            )}
            {#if category}
                <div class="flex flex-wrap gap-2 mb-6">
                    <div
                        class="inline-flex items-center gap-1 px-3 py-1 bg-dark-04 border border-border-input rounded-full text-sm"
                    >
                        <span>{category.name}</span>
                        <button
                            type="button"
                            onclick={() => (selectedCategoryId = ALL_CATEGORIES)}
                            class="hover:text-destructive transition-colors"
                        >
                            <X size={14} />
                        </button>
                    </div>
                </div>
            {/if}
        {/if}

        <!-- Products Grid -->
        {#if filteredProducts.length === 0}
            <div
                class="flex flex-col items-center justify-center py-16 text-center"
            >
                <div
                    class="w-16 h-16 rounded-full bg-dark-04 flex items-center justify-center mb-4"
                >
                    <ImageIcon size={32} class="text-dark-400" />
                </div>
                <h3 class="text-lg font-medium mb-2">
                    {products.length === 0 || selectedCategoryId === ALL_CATEGORIES
                        ? "Aucun produit"
                        : "Aucun résultat"}
                </h3>
                <p class="text-sm text-foreground-alt mb-6 max-w-sm">
                    {products.length === 0
                        ? "Commencez par créer votre premier produit pour proposer des services à vos clients."
                        : "Aucun produit ne correspond aux catégories sélectionnées."}
                </p>
                {#if products.length === 0}
                    <Button.Root
                        type="button"
                        class="cursor-pointer"
                        onclick={openCreateDialog}
                    >
                        <div
                            class="flex gap-2 items-center py-2 px-4 bg-dark text-white rounded-input hover:bg-dark/95 transition-all shadow-mini"
                        >
                            <Plus size={18} />
                            <span class="text-sm font-medium"
                                >Créer un produit</span
                            >
                        </div>
                    </Button.Root>
                {:else}
                    <button
                        type="button"
                        onclick={() => (selectedCategoryId = ALL_CATEGORIES)}
                        class="text-destructive hover:underline text-sm font-medium"
                    >
                        Effacer les filtres
                    </button>
                {/if}
            </div>
        {:else}
            <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                {#each filteredProducts as product (product.id)}
                    {@const activeImage = getActiveImage(product)}
                    <div
                        class="border border-border-card rounded-card bg-background shadow-card hover:shadow-popover transition-all overflow-hidden"
                    >
                        <!-- Image -->
                        <div class="aspect-video bg-dark-04 relative">
                            {#if activeImage}
                                <img
                                    src={activeImage.url}
                                    alt={activeImage.title}
                                    class="w-full h-full object-cover"
                                />
                            {:else}
                                <div
                                    class="w-full h-full flex items-center justify-center"
                                >
                                    <ImageIcon
                                        size={48}
                                        class="text-dark-300"
                                    />
                                </div>
                            {/if}
                            <!-- Status Badge -->
                            <div class="absolute top-3 right-3 flex gap-2">
                                <span
                                    class="px-2 py-1 text-xs font-medium rounded-md {product.availability ===
                                    'online'
                                        ? 'bg-blue-100 text-blue-800'
                                        : product.availability === 'hybrid'
                                          ? 'bg-purple-100 text-purple-800'
                                          : 'bg-green-100 text-green-800'}"
                                >
                                    {product.availability === "online"
                                        ? "En ligne"
                                        : product.availability === "hybrid"
                                          ? "Hybride"
                                          : "Présentiel"}
                                </span>
                                <span
                                    class="px-2 py-1 text-xs font-medium rounded-md {product.status ===
                                    'published'
                                        ? 'bg-green-100 text-green-800'
                                        : 'bg-yellow-100 text-yellow-800'}"
                                >
                                    {product.status === "published"
                                        ? "Publié"
                                        : "Brouillon"}
                                </span>
                            </div>
                        </div>

                        <!-- Content -->
                        <div class="p-4">
                            <div class="flex items-center gap-2 mb-2">
                                <h3 class="text-lg font-semibold">
                                    {product.name}
                                </h3>
                            </div>
                            <p
                                class="text-xs text-foreground-alt mb-3 px-2 py-1 bg-dark-04 rounded-md inline-block"
                            >
                                {getCategoryName(product.category)}
                            </p>
                            <p
                                class="text-sm text-foreground-alt line-clamp-2 mb-4"
                            >
                                {product.description}
                            </p>

                            <!-- Metadata -->
                            <div
                                class="flex items-center gap-4 text-xs text-foreground-alt mb-4"
                            >
                                <div class="flex items-center gap-1">
                                    <Clock size={14} />
                                    <span
                                        >{formatDuration(
                                            product.duration,
                                        )}</span
                                    >
                                </div>
                                <div class="flex items-center gap-1">
                                    {#if product.availability === "online"}
                                        <span>En ligne</span>
                                    {:else if product.availability === "hybrid"}
                                        <span>Hybride</span>
                                    {:else}
                                        <span>Présentiel</span>
                                    {/if}
                                </div>
                            </div>

                            <!-- Actions -->
                            <div class="flex gap-2">
                                <Button.Root
                                    type="button"
                                    class="cursor-pointer flex-1"
                                    onclick={() => openEditDialog(product)}
                                >
                                    <div
                                        class="flex gap-2 items-center justify-center py-2 px-3 border border-border-input rounded-input hover:bg-dark-04 transition-all text-sm font-medium"
                                    >
                                        <Pencil size={14} />
                                        <span>Modifier</span>
                                    </div>
                                </Button.Root>
                                <Button.Root
                                    type="button"
                                    class="cursor-pointer"
                                    onclick={() => openDeleteDialog(product)}
                                >
                                    <div
                                        class="flex items-center justify-center py-2 px-3 border border-destructive/20 text-destructive rounded-input hover:bg-destructive/10 transition-all"
                                    >
                                        <Trash2 size={14} />
                                    </div>
                                </Button.Root>
                            </div>
                        </div>
                    </div>
                {/each}
            </div>
        {/if}
    </div>
</div>

<!-- Create Product Dialog/Drawer -->
{#if isMobile}
    <Drawer bind:isOpen={createDialogOpen}>
        <div
            class="sticky top-0 bg-white pb-4 border-b border-border-card -mx-4 px-4 -mt-4 pt-4 z-10"
        >
            <div class="flex items-center justify-between mb-2">
                <h2 class="text-xl font-semibold tracking-tight">
                    Créer un produit
                </h2>
                <button
                    type="button"
                    onclick={() => (createDialogOpen = false)}
                    class="p-2 hover:bg-dark-04 rounded-md transition-all"
                >
                    <X class="text-foreground size-5" />
                </button>
            </div>
            <p class="text-foreground-alt text-sm">
                Remplissez les détails ci-dessous pour créer un nouveau produit.
            </p>
        </div>

        <form onsubmit={handleCreateSubmit} class="grid gap-4 pt-8">
            <div class="grid grid-cols-1 gap-4 w-full pb-4">
                {@render field("name", "Nom", inputCreate, formData.name, null)}
                {@render field(
                    "description",
                    "Description",
                    textareaCreate,
                    formData.description,
                    null,
                )}
                {@render field(
                    "category",
                    "Catégorie",
                    categorySelect,
                    formData.category,
                    null,
                )}
                {@render field(
                    "duration",
                    "Durée (minutes)",
                    numberInput,
                    formData.duration,
                    null,
                )}
                {@render field(
                    "bufferTime",
                    "Temps de buffer (minutes)",
                    numberInput,
                    formData.bufferTime,
                    null,
                )}
                {@render field(
                    "cancellationHours",
                    "Annulation (heures)",
                    numberInput,
                    formData.cancellationHours,
                    null,
                )}
                {@render field(
                    "stripeProductId",
                    "Stripe Product ID",
                    textInput,
                    formData.stripeProductId,
                    null,
                )}
                {@render field(
                    "status",
                    "Statut",
                    statusRadio,
                    formData.status,
                    null,
                )}
                {@render field(
                    "availability",
                    "Modalité",
                    availabilityRadio,
                    formData.availability,
                    null,
                )}
            </div>
            <div class="flex w-full justify-end gap-3">
                <button
                    type="button"
                    onclick={() => (createDialogOpen = false)}
                    class="h-input rounded-input border border-border-input hover:bg-dark-04 focus-visible:ring-dark focus-visible:ring-offset-background focus-visible:outline-hidden inline-flex items-center justify-center px-6 text-sm font-medium focus-visible:ring-2 focus-visible:ring-offset-2 active:scale-[0.98] cursor-pointer transition-all"
                >
                    Annuler
                </button>
                <button
                    type="submit"
                    class="h-input rounded-input bg-dark text-background shadow-mini hover:bg-dark/95 focus-visible:ring-dark focus-visible:ring-offset-background focus-visible:outline-hidden inline-flex items-center justify-center px-8 text-sm font-semibold focus-visible:ring-2 focus-visible:ring-offset-2 active:scale-[0.98] cursor-pointer transition-all"
                >
                    Créer
                </button>
            </div>
        </form>
    </Drawer>
{:else}
    <Dialog.Root bind:open={createDialogOpen}>
        <Dialog.Portal>
            <Dialog.Overlay
                class="data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 fixed inset-0 z-50 bg-black/80"
            />
            <Dialog.Content
                class="rounded-card-lg bg-background shadow-popover data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 data-[state=closed]:zoom-out-95 data-[state=open]:zoom-in-95 outline-hidden fixed left-[50%] top-[50%] z-50 w-full max-w-[calc(100%-2rem)] translate-x-[-50%] translate-y-[-50%] border p-8 sm:max-w-[540px] md:w-full max-h-[90vh] overflow-y-auto"
            >
                <Dialog.Title
                    class="w-full text-xl font-semibold tracking-tight"
                >
                    Créer un produit
                </Dialog.Title>
                <Dialog.Description class="text-foreground-alt !mt-1 text-sm">
                    Remplissez les détails ci-dessous pour créer un nouveau
                    produit.
                </Dialog.Description>

                <Separator.Root class="bg-muted mx-5 !mb-2 !mt-5 block h-px" />

                <form onsubmit={handleCreateSubmit} class="grid gap-4">
                    <div
                        class="grid grid-cols-[max-content_1fr] gap-4 w-full items-center pb-11 pt-7"
                    >
                        {@render field(
                            "name",
                            "Nom",
                            inputCreate,
                            formData.name,
                            null,
                        )}
                        {@render field(
                            "description",
                            "Description",
                            textareaCreate,
                            formData.description,
                            null,
                        )}
                        {@render field(
                            "category",
                            "Catégorie",
                            categorySelect,
                            formData.category,
                            null,
                        )}
                        {@render field(
                            "duration",
                            "Durée (minutes)",
                            numberInput,
                            formData.duration,
                            null,
                        )}
                        {@render field(
                            "bufferTime",
                            "Temps de buffer (minutes)",
                            numberInput,
                            formData.bufferTime,
                            null,
                        )}
                        {@render field(
                            "cancellationHours",
                            "Annulation (heures)",
                            numberInput,
                            formData.cancellationHours,
                            null,
                        )}
                        {@render field(
                            "stripeProductId",
                            "Stripe Product ID",
                            textInput,
                            formData.stripeProductId,
                            null,
                        )}
                        {@render field(
                            "status",
                            "Statut",
                            statusRadio,
                            formData.status,
                            null,
                        )}
                        {@render field(
                            "availability",
                            "Visibilité",
                            availabilityRadio,
                            formData.availability,
                            null,
                        )}
                    </div>
                    <div class="flex w-full justify-end gap-3">
                        <Button.Root type="button" class="cursor-pointer">
                            <Dialog.Close
                                class="h-input rounded-input border border-border-input hover:bg-dark-04 focus-visible:ring-dark focus-visible:ring-offset-background focus-visible:outline-hidden inline-flex items-center justify-center px-6 text-sm font-medium focus-visible:ring-2 focus-visible:ring-offset-2 active:scale-[0.98] cursor-pointer transition-all"
                            >
                                Annuler
                            </Dialog.Close>
                        </Button.Root>
                        <Button.Root type="submit" class="cursor-pointer">
                            <div
                                class="h-input rounded-input bg-dark text-background shadow-mini hover:bg-dark/95 focus-visible:ring-dark focus-visible:ring-offset-background focus-visible:outline-hidden inline-flex items-center justify-center px-8 text-sm font-semibold focus-visible:ring-2 focus-visible:ring-offset-2 active:scale-[0.98] cursor-pointer transition-all"
                            >
                                Créer
                            </div>
                        </Button.Root>
                    </div>
                </form>

                <Button.Root type="button" class="cursor-pointer">
                    <Dialog.Close
                        class="focus-visible:ring-foreground focus-visible:ring-offset-background focus-visible:outline-hidden absolute right-5 top-5 rounded-md focus-visible:ring-2 focus-visible:ring-offset-2 active:scale-[0.98] cursor-pointer"
                    >
                        <X class="text-foreground size-5" />
                        <span class="sr-only">Close</span>
                    </Dialog.Close>
                </Button.Root>
            </Dialog.Content>
        </Dialog.Portal>
    </Dialog.Root>
{/if}

<!-- Edit Product Dialog/Drawer -->
{#if isMobile}
    <Drawer bind:isOpen={editDialogOpen}>
        <div
            class="sticky top-0 bg-white pb-4 border-b border-border-card -mx-4 px-4 -mt-4 pt-4 z-10"
        >
            <div class="flex items-center justify-between mb-2">
                <h2 class="text-xl font-semibold tracking-tight">
                    Modifier "{selectedProduct?.name}"
                </h2>
                <button
                    type="button"
                    onclick={() => (editDialogOpen = false)}
                    class="p-2 hover:bg-dark-04 rounded-md transition-all"
                >
                    <X class="text-foreground size-5" />
                </button>
            </div>
            <p class="text-foreground-alt text-sm">
                Mettez à jour les détails du produit.
            </p>
        </div>

        <form onsubmit={handleEditSubmit} class="grid gap-4 pt-8">
            <div class="grid grid-cols-1 gap-4 w-full pb-4">
                {@render field("name", "Nom", inputEdit, formData.name, null)}
                {@render field(
                    "description",
                    "Description",
                    textareaEdit,
                    formData.description,
                    null,
                )}
                {@render field(
                    "category",
                    "Catégorie",
                    categorySelect,
                    formData.category,
                    null,
                )}
                {@render field(
                    "duration",
                    "Durée (minutes)",
                    numberInput,
                    formData.duration,
                    null,
                )}
                {@render field(
                    "bufferTime",
                    "Temps de buffer (minutes)",
                    numberInput,
                    formData.bufferTime,
                    null,
                )}
                {@render field(
                    "cancellationHours",
                    "Annulation (heures)",
                    numberInput,
                    formData.cancellationHours,
                    null,
                )}
                {@render field(
                    "stripeProductId",
                    "Stripe Product ID",
                    textInput,
                    formData.stripeProductId,
                    null,
                )}
                {@render field(
                    "status",
                    "Statut",
                    statusRadio,
                    formData.status,
                    null,
                )}
                {@render field(
                    "availability",
                    "Modalité",
                    availabilityRadio,
                    formData.availability,
                    null,
                )}
            </div>
            <div class="flex w-full justify-end gap-3">
                <button
                    type="button"
                    onclick={() => (editDialogOpen = false)}
                    class="h-input rounded-input border border-border-input hover:bg-dark-04 focus-visible:ring-dark focus-visible:ring-offset-background focus-visible:outline-hidden inline-flex items-center justify-center px-6 text-sm font-medium focus-visible:ring-2 focus-visible:ring-offset-2 active:scale-[0.98] cursor-pointer transition-all"
                >
                    Annuler
                </button>
                <button
                    type="submit"
                    class="h-input rounded-input bg-dark text-background shadow-mini hover:bg-dark/95 focus-visible:ring-dark focus-visible:ring-offset-background focus-visible:outline-hidden inline-flex items-center justify-center px-8 text-sm font-semibold focus-visible:ring-2 focus-visible:ring-offset-2 active:scale-[0.98] cursor-pointer transition-all"
                >
                    Enregistrer
                </button>
            </div>
        </form>
    </Drawer>
{:else}
    <Dialog.Root bind:open={editDialogOpen}>
        <Dialog.Portal>
            <Dialog.Overlay
                class="data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 fixed inset-0 z-50 bg-black/80"
            />
            <Dialog.Content
                class="rounded-card-lg bg-background shadow-popover data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 data-[state=closed]:zoom-out-95 data-[state=open]:zoom-in-95 outline-hidden fixed left-[50%] top-[50%] z-50 w-full max-w-[calc(100%-2rem)] translate-x-[-50%] translate-y-[-50%] border p-8 sm:max-w-[540px] md:w-full max-h-[90vh] overflow-y-auto"
            >
                <Dialog.Title
                    class="w-full text-xl font-semibold tracking-tight"
                >
                    Modifier {selectedProduct?.name}
                </Dialog.Title>
                <Dialog.Description class="text-foreground-alt !mt-1 text-sm">
                    Mettez à jour les détails du produit.
                </Dialog.Description>

                <Separator.Root class="bg-muted mx-5 !mb-2 !mt-5 block h-px" />

                <form onsubmit={handleEditSubmit} class="grid gap-4">
                    <div
                        class="grid grid-cols-[max-content_1fr] gap-4 w-full items-center pb-11 pt-7"
                    >
                        {@render field(
                            "name",
                            "Nom",
                            inputEdit,
                            formData.name,
                            null,
                        )}
                        {@render field(
                            "description",
                            "Description",
                            textareaEdit,
                            formData.description,
                            null,
                        )}
                        {@render field(
                            "category",
                            "Catégorie",
                            categorySelect,
                            formData.category,
                            null,
                        )}
                        {@render field(
                            "duration",
                            "Durée (minutes)",
                            numberInput,
                            formData.duration,
                            null,
                        )}
                        {@render field(
                            "bufferTime",
                            "Temps de buffer (minutes)",
                            numberInput,
                            formData.bufferTime,
                            null,
                        )}
                        {@render field(
                            "cancellationHours",
                            "Annulation (heures)",
                            numberInput,
                            formData.cancellationHours,
                            null,
                        )}
                        {@render field(
                            "stripeProductId",
                            "Stripe Product ID",
                            textInput,
                            formData.stripeProductId,
                            null,
                        )}
                        {@render field(
                            "status",
                            "Statut",
                            statusRadio,
                            formData.status,
                            null,
                        )}
                        {@render field(
                            "availability",
                            "Visibilité",
                            availabilityRadio,
                            formData.availability,
                            null,
                        )}
                    </div>
                    <div class="flex w-full justify-end gap-3">
                        <Button.Root type="button" class="cursor-pointer">
                            <Dialog.Close
                                class="h-input rounded-input border border-border-input hover:bg-dark-04 focus-visible:ring-dark focus-visible:ring-offset-background focus-visible:outline-hidden inline-flex items-center justify-center px-6 text-sm font-medium focus-visible:ring-2 focus-visible:ring-offset-2 active:scale-[0.98] cursor-pointer transition-all"
                            >
                                Annuler
                            </Dialog.Close>
                        </Button.Root>
                        <Button.Root type="submit" class="cursor-pointer">
                            <div
                                class="h-input rounded-input bg-dark text-background shadow-mini hover:bg-dark/95 focus-visible:ring-dark focus-visible:ring-offset-background focus-visible:outline-hidden inline-flex items-center justify-center px-8 text-sm font-semibold focus-visible:ring-2 focus-visible:ring-offset-2 active:scale-[0.98] cursor-pointer transition-all"
                            >
                                Enregistrer
                            </div>
                        </Button.Root>
                    </div>
                </form>

                <Button.Root type="button" class="cursor-pointer">
                    <Dialog.Close
                        class="focus-visible:ring-foreground focus-visible:ring-offset-background focus-visible:outline-hidden absolute right-5 top-5 rounded-md focus-visible:ring-2 focus-visible:ring-offset-2 active:scale-[0.98] cursor-pointer"
                    >
                        <X class="text-foreground size-5" />
                        <span class="sr-only">Close</span>
                    </Dialog.Close>
                </Button.Root>
            </Dialog.Content>
        </Dialog.Portal>
    </Dialog.Root>
{/if}

<!-- Delete Product Dialog/Drawer -->
{#if isMobile}
    <Drawer bind:isOpen={deleteDialogOpen}>
        <div
            class="sticky top-0 bg-white pb-4 border-b border-border-card -mx-4 px-4 -mt-4 pt-4 z-10"
        >
            <div class="flex items-center justify-between mb-2">
                <h2 class="text-xl font-semibold tracking-tight">
                    Supprimer le produit
                </h2>
                <button
                    type="button"
                    onclick={() => (deleteDialogOpen = false)}
                    class="p-2 hover:bg-dark-04 rounded-md transition-all"
                >
                    <X class="text-foreground size-5" />
                </button>
            </div>
            <p class="text-foreground-alt text-sm">
                Êtes-vous sûr de vouloir supprimer le produit "<span
                    class="font-medium">{selectedProduct?.name}</span
                >" ? Cette action est irréversible.
            </p>
        </div>

        <form onsubmit={handleDeleteSubmit} class="pt-8">
            <div class="flex w-full justify-end gap-3">
                <button
                    type="button"
                    onclick={() => (deleteDialogOpen = false)}
                    class="h-input rounded-input border border-border-input hover:bg-dark-04 focus-visible:ring-dark focus-visible:ring-offset-background focus-visible:outline-hidden inline-flex items-center justify-center px-6 text-sm font-medium focus-visible:ring-2 focus-visible:ring-offset-2 active:scale-[0.98] cursor-pointer transition-all"
                >
                    Annuler
                </button>
                <button
                    type="submit"
                    class="h-input rounded-input bg-destructive text-white shadow-mini hover:bg-destructive/90 focus-visible:ring-destructive focus-visible:ring-offset-background focus-visible:outline-hidden inline-flex items-center justify-center px-8 text-sm font-semibold focus-visible:ring-2 focus-visible:ring-offset-2 active:scale-[0.98] cursor-pointer transition-all"
                >
                    Supprimer
                </button>
            </div>
        </form>
    </Drawer>
{:else}
    <Dialog.Root bind:open={deleteDialogOpen}>
        <Dialog.Portal>
            <Dialog.Overlay
                class="data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 fixed inset-0 z-50 bg-black/80"
            />
            <Dialog.Content
                class="rounded-card-lg bg-background shadow-popover data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 data-[state=closed]:zoom-out-95 data-[state=open]:zoom-in-95 outline-hidden fixed left-[50%] top-[50%] z-50 w-full max-w-[calc(100%-2rem)] translate-x-[-50%] translate-y-[-50%] border p-8 sm:max-w-[440px] md:w-full"
            >
                <Dialog.Title
                    class="w-full text-xl font-semibold tracking-tight"
                >
                    Supprimer le produit
                </Dialog.Title>
                <Dialog.Description class="text-foreground-alt !mt-2 text-sm">
                    Êtes-vous sûr de vouloir supprimer le produit "<span
                        class="font-medium">{selectedProduct?.name}</span
                    >" ? Cette action est irréversible.
                </Dialog.Description>

                <form onsubmit={handleDeleteSubmit} class="mt-8">
                    <div class="flex w-full justify-end gap-3">
                        <Button.Root type="button" class="cursor-pointer">
                            <Dialog.Close
                                class="h-input rounded-input border border-border-input hover:bg-dark-04 focus-visible:ring-dark focus-visible:ring-offset-background focus-visible:outline-hidden inline-flex items-center justify-center px-6 text-sm font-medium focus-visible:ring-2 focus-visible:ring-offset-2 active:scale-[0.98] cursor-pointer transition-all"
                            >
                                Annuler
                            </Dialog.Close>
                        </Button.Root>
                        <Button.Root type="submit" class="cursor-pointer">
                            <div
                                class="h-input rounded-input bg-destructive text-white shadow-mini hover:bg-destructive/90 focus-visible:ring-destructive focus-visible:ring-offset-background focus-visible:outline-hidden inline-flex items-center justify-center px-8 text-sm font-semibold focus-visible:ring-2 focus-visible:ring-offset-2 active:scale-[0.98] cursor-pointer transition-all"
                            >
                                Supprimer
                            </div>
                        </Button.Root>
                    </div>
                </form>

                <Button.Root type="button" class="cursor-pointer">
                    <Dialog.Close
                        class="focus-visible:ring-foreground focus-visible:ring-offset-background focus-visible:outline-hidden absolute right-5 top-5 rounded-md focus-visible:ring-2 focus-visible:ring-offset-2 active:scale-[0.98] cursor-pointer"
                    >
                        <X class="text-foreground size-5" />
                        <span class="sr-only">Close</span>
                    </Dialog.Close>
                </Button.Root>
            </Dialog.Content>
        </Dialog.Portal>
    </Dialog.Root>
{/if}

<!-- Snippets for form fields -->
{#snippet field(
    name: string,
    label: string,
    inputSnippet: Input,
    value: any,
    error: any,
)}
    <Label.Root for={name} class="text-sm font-semibold">
        {label}
    </Label.Root>
    <div class="relative w-full">
        {@render inputSnippet(name)}
        {#if error && error.length > 0}
            <p class="text-xs text-destructive mt-1">{error[0]}</p>
        {/if}
    </div>
{/snippet}

{#snippet inputCreate(fieldName: string)}
    {#if fieldName === "name"}
        <input
            id={fieldName}
            type="text"
            name={fieldName}
            bind:value={formData.name}
            class="h-input rounded-card-sm border-border-input bg-background placeholder:text-foreground-alt/50 hover:border-dark-40 focus:ring-foreground focus:ring-offset-background focus:outline-hidden inline-flex w-full items-center border px-4 text-base focus:ring-2 focus:ring-offset-2 sm:text-sm transition-all"
            required
        />
    {/if}
{/snippet}

{#snippet inputEdit(fieldName: string)}
    {#if fieldName === "name"}
        <input
            id={fieldName}
            type="text"
            name={fieldName}
            bind:value={formData.name}
            class="h-input rounded-card-sm border-border-input bg-background placeholder:text-foreground-alt/50 hover:border-dark-40 focus:ring-foreground focus:ring-offset-background focus:outline-hidden inline-flex w-full items-center border px-4 text-base focus:ring-2 focus:ring-offset-2 sm:text-sm transition-all"
            required
        />
    {/if}
{/snippet}

{#snippet textareaCreate(fieldName: string)}
    {#if fieldName === "description"}
        <textarea
            id={fieldName}
            name={fieldName}
            rows="4"
            bind:value={formData.description}
            class="rounded-card-sm border border-border-input bg-background placeholder:text-foreground-alt/50 hover:border-dark-40 focus:ring-foreground focus:ring-offset-background focus:outline-hidden w-full px-4 py-2 text-base focus:ring-2 focus:ring-offset-2 sm:text-sm transition-all"
            placeholder="Décrivez le produit"
            required
        ></textarea>
    {/if}
{/snippet}

{#snippet textareaEdit(fieldName: string)}
    {#if fieldName === "description"}
        <textarea
            id={fieldName}
            name={fieldName}
            rows="4"
            bind:value={formData.description}
            class="rounded-card-sm border border-border-input bg-background placeholder:text-foreground-alt/50 hover:border-dark-40 focus:ring-foreground focus:ring-offset-background focus:outline-hidden w-full px-4 py-2 text-base focus:ring-2 focus:ring-offset-2 sm:text-sm transition-all"
            placeholder="Décrivez le produit"
            required
        ></textarea>
    {/if}
{/snippet}

{#snippet categorySelect(name: string)}
    <select
        id={name}
        {name}
        bind:value={formData.category}
        class="h-input rounded-card-sm border-border-input bg-background hover:border-dark-40 focus:ring-foreground focus:ring-offset-background focus:outline-hidden inline-flex w-full items-center border px-4 text-base focus:ring-2 focus:ring-offset-2 sm:text-sm transition-all"
        required
    >
        <option value="">Sélectionner une catégorie</option>
        {#each mockCategories as category}
            <option value={category.id}>{category.name}</option>
        {/each}
    </select>
{/snippet}

{#snippet numberInput(name: string)}
    {#if name === "duration" || name === "bufferTime" || name === "cancellationHours"}
        <input
            id={name}
            type="number"
            {name}
            bind:value={formData[name]}
            min="0"
            class="h-input rounded-card-sm border-border-input bg-background placeholder:text-foreground-alt/50 hover:border-dark-40 focus:ring-foreground focus:ring-offset-background focus:outline-hidden inline-flex w-full items-center border px-4 text-base focus:ring-2 focus:ring-offset-2 sm:text-sm transition-all"
            required
        />
    {/if}
{/snippet}

{#snippet textInput(name: string)}
    {#if name === "stripeProductId"}
        <input
            id={name}
            type="text"
            {name}
            bind:value={formData.stripeProductId}
            class="h-input rounded-card-sm border-border-input bg-background placeholder:text-foreground-alt/50 hover:border-dark-40 focus:ring-foreground focus:ring-offset-background focus:outline-hidden inline-flex w-full items-center border px-4 text-base focus:ring-2 focus:ring-offset-2 sm:text-sm transition-all"
        />
    {/if}
{/snippet}

{#snippet statusRadio(name: string)}
    {#if name === "status"}
        <div class="flex gap-4">
            <label class="flex items-center gap-2 cursor-pointer">
                <input
                    type="radio"
                    name="status"
                    value="draft"
                    bind:group={formData.status}
                    class="w-4 h-4 text-dark border-gray-300 focus:ring-dark"
                />
                <span class="text-sm">Brouillon</span>
            </label>
            <label class="flex items-center gap-2 cursor-pointer">
                <input
                    type="radio"
                    name="status"
                    value="published"
                    bind:group={formData.status}
                    class="w-4 h-4 text-dark border-gray-300 focus:ring-dark"
                />
                <span class="text-sm">Publié</span>
            </label>
        </div>
    {/if}
{/snippet}

{#snippet availabilityRadio(name: string)}
    {#if name === "availability"}
        <div class="flex gap-4">
            <label class="flex items-center gap-2 cursor-pointer">
                <input
                    type="radio"
                    name="availability"
                    value="in-person"
                    bind:group={formData.availability}
                    class="w-4 h-4 text-dark border-gray-300 focus:ring-dark"
                />
                <span class="text-sm">Présentiel</span>
            </label>
            <label class="flex items-center gap-2 cursor-pointer">
                <input
                    type="radio"
                    name="availability"
                    value="online"
                    bind:group={formData.availability}
                    class="w-4 h-4 text-dark border-gray-300 focus:ring-dark"
                />
                <span class="text-sm">En ligne</span>
            </label>
            <label class="flex items-center gap-2 cursor-pointer">
                <input
                    type="radio"
                    name="availability"
                    value="hybrid"
                    bind:group={formData.availability}
                    class="w-4 h-4 text-dark border-gray-300 focus:ring-dark"
                />
                <span class="text-sm">Hybride</span>
            </label>
        </div>
    {/if}
{/snippet}
