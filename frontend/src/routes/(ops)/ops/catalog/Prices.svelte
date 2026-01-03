<script lang="ts">
    import { Dialog, Button, Label, Separator, Combobox } from "bits-ui";
    import {
        Plus,
        Pencil,
        X,
        Filter,
        Check,
        ChevronsUpDown,
        DollarSign,
        ArrowUpDown,
    } from "@lucide/svelte";
    import type { Snippet } from "svelte";
    import type { PageData } from "./$types";
    import Drawer from "$lib/ui/Drawer.svelte";
    import { superForm } from "sveltekit-superforms";

    import { browser } from "$app/environment";

    // Detect if we're on mobile
    let isMobile = $state(false);

    if (browser) {
        isMobile = window.innerWidth < 768;
        window.addEventListener("resize", () => {
            isMobile = window.innerWidth < 768;
        });
    }

    type Input = Snippet<[string]>;

    // Price type from API
    type Price = {
        id: string;
        productId: string;
        stripePriceId: string;
        amount: number; // in cents
        currency: string;
        interval: "one_time" | "month" | "year";
        isActive: boolean;
        createdAt: string;
        updatedAt: string;
    };

    interface Props {
        data: PageData;
    }

    let { data }: Props = $props();

    // Extract prices and products from page data
    let prices = $state<Price[]>(data.prices || []);
    let products = $state(data.products || []);

    // Product filter state (single select)
    const ALL_PRODUCTS = "";
    let selectedProductId = $state<string>(ALL_PRODUCTS);

    // Sort state
    type SortField = "amount" | "product" | "interval" | "isActive";
    let sortField = $state<SortField | null>(null);
    let sortDirection = $state<"asc" | "desc">("asc");

    // Filtered prices based on selected product
    const filteredPrices = $derived(
        selectedProductId === ALL_PRODUCTS
            ? prices
            : prices.filter((p) => p.productId === selectedProductId),
    );

    // Sorted and filtered prices
    const displayedPrices = $derived(() => {
        let result = [...filteredPrices];

        if (sortField) {
            result.sort((a, b) => {
                let compareValue = 0;

                if (sortField === "amount") {
                    compareValue = a.amount - b.amount;
                } else if (sortField === "product") {
                    const productA = getProductName(a.productId);
                    const productB = getProductName(b.productId);
                    compareValue = productA.localeCompare(productB);
                } else if (sortField === "interval") {
                    const intervalOrder = { one_time: 0, month: 1, year: 2 };
                    compareValue = intervalOrder[a.interval] - intervalOrder[b.interval];
                } else if (sortField === "isActive") {
                    compareValue = (a.isActive === b.isActive) ? 0 : a.isActive ? -1 : 1;
                }

                return sortDirection === "asc" ? compareValue : -compareValue;
            });
        }

        return result;
    });

    // Dialog states
    let createDialogOpen = $state(false);
    let editDialogOpen = $state(false);

    // Currently selected price for edit
    let selectedPrice: Price | null = $state(null);

    // Superforms for create and update
    const {
        form: createForm,
        errors: createErrors,
        enhance: createEnhance,
    } = superForm(data.createPriceForm, {
        resetForm: true,
        onUpdated({ form }) {
            if (form.valid) {
                createDialogOpen = false;
            }
        },
    });

    const {
        form: updateForm,
        errors: updateErrors,
        enhance: updateEnhance,
    } = superForm(data.updatePriceForm, {
        resetForm: false,
        onUpdated({ form }) {
            if (form.valid) {
                editDialogOpen = false;
            }
        },
    });

    function openCreateDialog() {
        createDialogOpen = true;
    }

    function openEditDialog(price: Price) {
        selectedPrice = price;
        $updateForm.id = price.id;
        $updateForm.active = price.isActive;
        $updateForm.nickname = undefined;
        editDialogOpen = true;
    }

    function getProductName(productId: string): string {
        return products.find((p) => p.id === productId)?.name || "Inconnu";
    }

    function formatAmount(amount: number, currency: string): string {
        const value = amount / 100;
        const symbol = currency === "eur" ? "€" : currency === "usd" ? "$" : currency.toUpperCase();
        return `${value.toFixed(2)} ${symbol}`;
    }

    function formatInterval(interval: string): string {
        const labels: Record<string, string> = {
            one_time: "Paiement unique",
            month: "Mensuel",
            year: "Annuel",
        };
        return labels[interval] || interval;
    }

    function toggleSort(field: SortField) {
        if (sortField === field) {
            sortDirection = sortDirection === "asc" ? "desc" : "asc";
        } else {
            sortField = field;
            sortDirection = "asc";
        }
    }
</script>

<div class="h-full bg-white relative">
    <!-- Filter & Prices Table -->
    <div class="p-4 md:p-8">
        <!-- Header with Create Button -->
        <div class="flex items-center justify-between mb-6">
            <h2 class="text-lg font-semibold">Liste des prix</h2>
            <Button.Root
                type="button"
                class="cursor-pointer hidden md:flex"
                onclick={openCreateDialog}
            >
                <div
                    class="flex gap-2 items-center py-2 px-4 bg-dark text-white rounded-input hover:bg-dark/95 transition-all shadow-mini"
                >
                    <Plus size={18} />
                    <span class="text-sm font-medium">Nouveau Prix</span>
                </div>
            </Button.Root>
        </div>

        <!-- Mobile floating button -->
        <Button.Root
            type="button"
            class="cursor-pointer fixed bottom-20 right-4 z-10 md:hidden"
            onclick={openCreateDialog}
        >
            <div
                class="flex gap-2 items-center py-2 bg-dark text-white rounded-input hover:bg-dark/95 transition-all shadow-mini w-12 h-12 justify-center"
            >
                <Plus size={18} />
            </div>
        </Button.Root>
        <!-- Filter Section -->
        <div
            class="flex flex-col md:flex-row md:items-center justify-between gap-4 mb-6"
        >
            <Combobox.Root type="single" bind:value={selectedProductId}>
                <div class="relative w-full md:w-auto md:min-w-[280px]">
                    <Filter
                        class="text-muted-foreground absolute start-3 top-1/2 size-4 -translate-y-1/2"
                    />
                    <Combobox.Input
                        class="h-input rounded-card-sm border-border-input bg-background placeholder:text-foreground-alt/50 hover:border-dark-40 focus:ring-foreground focus:ring-offset-background focus:outline-hidden inline-flex w-full items-center border px-9 text-sm focus:ring-2 focus:ring-offset-2 transition-all"
                        placeholder="Filtrer par produit"
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
                            <Combobox.Item value={ALL_PRODUCTS} label="Tous les produits">
                                {#snippet children({ selected })}
                                    <div
                                        class="flex items-center justify-between px-3 py-2 rounded-md hover:bg-dark-04 cursor-pointer"
                                    >
                                        <span class="text-sm"
                                            >Tous les produits</span
                                        >
                                        {#if selected}
                                            <Check class="size-4" />
                                        {/if}
                                    </div>
                                {/snippet}
                            </Combobox.Item>
                            <Combobox.Group class="pt-1">
                                {#each mockProducts as product}
                                    <Combobox.Item
                                        value={product.id}
                                        label={product.name}
                                    >
                                        {#snippet children({ selected })}
                                            <div
                                                class="flex items-center justify-between px-3 py-2 rounded-md hover:bg-dark-04 cursor-pointer"
                                            >
                                                <span class="text-sm"
                                                    >{product.name}</span
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

            <!-- Price count -->
            <div class="flex items-center gap-2 text-sm text-foreground-alt">
                {#if selectedProductId !== ALL_PRODUCTS}
                    <span class="hidden md:inline">•</span>
                    <span
                        >{displayedPrices().length} prix</span
                    >
                    <button
                        type="button"
                        onclick={() => (selectedProductId = ALL_PRODUCTS)}
                        class="text-destructive hover:underline"
                    >
                        Réinitialiser
                    </button>
                {:else}
                    <span>{prices.length} prix</span>
                {/if}
            </div>
        </div>

        <!-- Selected product chip -->
        {#if selectedProductId !== ALL_PRODUCTS}
            {@const product = mockProducts.find(
                (p) => p.id === selectedProductId,
            )}
            {#if product}
                <div class="flex flex-wrap gap-2 mb-6">
                    <div
                        class="inline-flex items-center gap-1 px-3 py-1 bg-dark-04 border border-border-input rounded-full text-sm"
                    >
                        <span>{product.name}</span>
                        <button
                            type="button"
                            onclick={() => (selectedProductId = ALL_PRODUCTS)}
                            class="hover:text-destructive transition-colors"
                        >
                            <X size={14} />
                        </button>
                    </div>
                </div>
            {/if}
        {/if}

        <!-- Prices Table -->
        {#if displayedPrices().length === 0}
            <div
                class="flex flex-col items-center justify-center py-16 text-center"
            >
                <div
                    class="w-16 h-16 rounded-full bg-dark-04 flex items-center justify-center mb-4"
                >
                    <DollarSign size={32} class="text-dark-400" />
                </div>
                <h3 class="text-lg font-medium mb-2">
                    {prices.length === 0 || selectedProductId === ALL_PRODUCTS
                        ? "Aucun prix"
                        : "Aucun résultat"}
                </h3>
                <p class="text-sm text-foreground-alt mb-6 max-w-sm">
                    {prices.length === 0
                        ? "Commencez par créer votre premier prix pour vos produits."
                        : "Aucun prix ne correspond au produit sélectionné."}
                </p>
                {#if prices.length === 0}
                    <Button.Root
                        type="button"
                        class="cursor-pointer"
                        onclick={openCreateDialog}
                    >
                        <div
                            class="flex gap-2 items-center py-2 px-4 bg-dark text-white rounded-input hover:bg-dark/95 transition-all shadow-mini"
                        >
                            <Plus size={18} />
                            <span class="text-sm font-medium">Créer un prix</span>
                        </div>
                    </Button.Root>
                {:else}
                    <button
                        type="button"
                        onclick={() => (selectedProductId = ALL_PRODUCTS)}
                        class="text-destructive hover:underline text-sm font-medium"
                    >
                        Effacer les filtres
                    </button>
                {/if}
            </div>
        {:else}
            <div class="overflow-x-auto border border-border-card rounded-card">
                <table class="w-full">
                    <thead class="bg-dark-04">
                        <tr>
                            <th
                                class="px-4 py-3 text-left text-sm font-semibold text-foreground cursor-pointer hover:bg-dark-10 transition-colors border-b border-r border-border-card"
                                onclick={() => toggleSort("product")}
                            >
                                <div class="flex items-center gap-2">
                                    <span>Produit</span>
                                    {#if sortField === "product"}
                                        <ArrowUpDown size={14} class={sortDirection === "desc" ? "rotate-180" : ""} />
                                    {:else}
                                        <ArrowUpDown size={14} class="opacity-30" />
                                    {/if}
                                </div>
                            </th>
                            <th
                                class="px-4 py-3 text-left text-sm font-semibold text-foreground cursor-pointer hover:bg-dark-10 transition-colors border-b border-r border-border-card"
                                onclick={() => toggleSort("amount")}
                            >
                                <div class="flex items-center gap-2">
                                    <span>Montant</span>
                                    {#if sortField === "amount"}
                                        <ArrowUpDown size={14} class={sortDirection === "desc" ? "rotate-180" : ""} />
                                    {:else}
                                        <ArrowUpDown size={14} class="opacity-30" />
                                    {/if}
                                </div>
                            </th>
                            <th
                                class="px-4 py-3 text-left text-sm font-semibold text-foreground cursor-pointer hover:bg-dark-10 transition-colors border-b border-r border-border-card"
                                onclick={() => toggleSort("interval")}
                            >
                                <div class="flex items-center gap-2">
                                    <span>Type</span>
                                    {#if sortField === "interval"}
                                        <ArrowUpDown size={14} class={sortDirection === "desc" ? "rotate-180" : ""} />
                                    {:else}
                                        <ArrowUpDown size={14} class="opacity-30" />
                                    {/if}
                                </div>
                            </th>
                            <th
                                class="px-4 py-3 text-left text-sm font-semibold text-foreground cursor-pointer hover:bg-dark-10 transition-colors border-b border-r border-border-card"
                                onclick={() => toggleSort("isActive")}
                            >
                                <div class="flex items-center gap-2">
                                    <span>Statut</span>
                                    {#if sortField === "isActive"}
                                        <ArrowUpDown size={14} class={sortDirection === "desc" ? "rotate-180" : ""} />
                                    {:else}
                                        <ArrowUpDown size={14} class="opacity-30" />
                                    {/if}
                                </div>
                            </th>
                            <th
                                class="px-4 py-3 text-left text-sm font-semibold text-foreground border-b border-r border-border-card"
                            >
                                Stripe Price ID
                            </th>
                            <th
                                class="px-4 py-3 text-right text-sm font-semibold text-foreground border-b border-border-card"
                            >
                                Actions
                            </th>
                        </tr>
                    </thead>
                    <tbody>
                        {#each displayedPrices() as price (price.id)}
                            <tr class="hover:bg-dark-04 transition-colors">
                                <td class="px-4 py-3 text-sm text-foreground border-b border-r border-border-card">
                                    {getProductName(price.productId)}
                                </td>
                                <td
                                    class="px-4 py-3 text-sm font-medium text-foreground border-b border-r border-border-card"
                                >
                                    {formatAmount(price.amount, price.currency)}
                                </td>
                                <td class="px-4 py-3 text-sm text-foreground border-b border-r border-border-card">
                                    {formatInterval(price.interval)}
                                </td>
                                <td class="px-4 py-3 border-b border-r border-border-card">
                                    <span
                                        class="px-2 py-1 text-xs font-medium rounded-md {price.isActive
                                            ? 'bg-green-100 text-green-800'
                                            : 'bg-gray-100 text-gray-800'}"
                                    >
                                        {price.isActive ? "Actif" : "Inactif"}
                                    </span>
                                </td>
                                <td
                                    class="px-4 py-3 text-sm text-foreground-alt font-mono border-b border-r border-border-card"
                                >
                                    {price.stripePriceId}
                                </td>
                                <td class="px-4 py-3 border-b border-border-card">
                                    <div class="flex gap-2 justify-end">
                                        <Button.Root
                                            type="button"
                                            class="cursor-pointer"
                                            onclick={() => openEditDialog(price)}
                                        >
                                            <div
                                                class="flex gap-2 items-center p-2 px-3 border border-border-input rounded-md hover:bg-dark-04 transition-all"
                                            >
                                                <Pencil size={14} />
                                                <span class="text-sm font-medium">Modifier</span>
                                            </div>
                                        </Button.Root>
                                    </div>
                                </td>
                            </tr>
                        {/each}
                    </tbody>
                </table>
            </div>
        {/if}
    </div>
</div>

<!-- Create Price Dialog/Drawer -->
{#if isMobile}
    <Drawer bind:isOpen={createDialogOpen}>
        <div
            class="sticky top-0 bg-white pb-4 border-b border-border-card -mx-4 px-4 -mt-4 pt-4 z-10"
        >
            <div class="flex items-center justify-between mb-2">
                <h2 class="text-xl font-semibold tracking-tight">
                    Créer un prix
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
                Remplissez les détails ci-dessous pour créer un nouveau prix.
            </p>
        </div>

        <form onsubmit={handleCreateSubmit} class="grid gap-4 pt-8">
            <div class="grid grid-cols-1 gap-4 w-full pb-4">
                {@render field("productId", "Produit", productSelect, formData.productId, null)}
                {@render field("amount", "Montant (en centimes)", numberInput, formData.amount, null)}
                {@render field("currency", "Devise", currencySelect, formData.currency, null)}
                {@render field("interval", "Type", intervalSelect, formData.interval, null)}
                {@render field("stripePriceId", "Stripe Price ID", textInput, formData.stripePriceId, null)}
                {@render field("isActive", "Statut", statusCheckbox, formData.isActive, null)}
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
                    Créer un prix
                </Dialog.Title>
                <Dialog.Description class="text-foreground-alt !mt-1 text-sm">
                    Remplissez les détails ci-dessous pour créer un nouveau
                    prix.
                </Dialog.Description>

                <Separator.Root class="bg-muted mx-5 !mb-2 !mt-5 block h-px" />

                <form onsubmit={handleCreateSubmit} class="grid gap-4">
                    <div
                        class="grid grid-cols-[max-content_1fr] gap-4 w-full items-center pb-11 pt-7"
                    >
                        {@render field(
                            "productId",
                            "Produit",
                            productSelect,
                            formData.productId,
                            null,
                        )}
                        {@render field(
                            "amount",
                            "Montant (centimes)",
                            numberInput,
                            formData.amount,
                            null,
                        )}
                        {@render field(
                            "currency",
                            "Devise",
                            currencySelect,
                            formData.currency,
                            null,
                        )}
                        {@render field(
                            "interval",
                            "Type",
                            intervalSelect,
                            formData.interval,
                            null,
                        )}
                        {@render field(
                            "stripePriceId",
                            "Stripe Price ID",
                            textInput,
                            formData.stripePriceId,
                            null,
                        )}
                        {@render field(
                            "isActive",
                            "Statut",
                            statusCheckbox,
                            formData.isActive,
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

<!-- Edit Price Dialog/Drawer -->
{#if isMobile}
    <Drawer bind:isOpen={editDialogOpen}>
        <div
            class="sticky top-0 bg-white pb-4 border-b border-border-card -mx-4 px-4 -mt-4 pt-4 z-10"
        >
            <div class="flex items-center justify-between mb-2">
                <h2 class="text-xl font-semibold tracking-tight">
                    Modifier le prix
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
                Mettez à jour les détails du prix.
            </p>
        </div>

        <form onsubmit={handleEditSubmit} class="grid gap-4 pt-8">
            <div class="grid grid-cols-1 gap-4 w-full pb-4">
                {@render field("productId", "Produit", productSelect, formData.productId, null)}
                {@render field("amount", "Montant (en centimes)", numberInput, formData.amount, null)}
                {@render field("currency", "Devise", currencySelect, formData.currency, null)}
                {@render field("interval", "Type", intervalSelect, formData.interval, null)}
                {@render field("stripePriceId", "Stripe Price ID", textInput, formData.stripePriceId, null)}
                {@render field("isActive", "Statut", statusCheckbox, formData.isActive, null)}
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
                    Modifier le prix
                </Dialog.Title>
                <Dialog.Description class="text-foreground-alt !mt-1 text-sm">
                    Mettez à jour les détails du prix.
                </Dialog.Description>

                <Separator.Root class="bg-muted mx-5 !mb-2 !mt-5 block h-px" />

                <form onsubmit={handleEditSubmit} class="grid gap-4">
                    <div
                        class="grid grid-cols-[max-content_1fr] gap-4 w-full items-center pb-11 pt-7"
                    >
                        {@render field(
                            "productId",
                            "Produit",
                            productSelect,
                            formData.productId,
                            null,
                        )}
                        {@render field(
                            "amount",
                            "Montant (centimes)",
                            numberInput,
                            formData.amount,
                            null,
                        )}
                        {@render field(
                            "currency",
                            "Devise",
                            currencySelect,
                            formData.currency,
                            null,
                        )}
                        {@render field(
                            "interval",
                            "Type",
                            intervalSelect,
                            formData.interval,
                            null,
                        )}
                        {@render field(
                            "stripePriceId",
                            "Stripe Price ID",
                            textInput,
                            formData.stripePriceId,
                            null,
                        )}
                        {@render field(
                            "isActive",
                            "Statut",
                            statusCheckbox,
                            formData.isActive,
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

{#snippet productSelect(name: string)}
    <select
        id={name}
        {name}
        bind:value={formData.productId}
        class="h-input rounded-card-sm border-border-input bg-background hover:border-dark-40 focus:ring-foreground focus:ring-offset-background focus:outline-hidden inline-flex w-full items-center border px-4 text-base focus:ring-2 focus:ring-offset-2 sm:text-sm transition-all"
        required
    >
        <option value="">Sélectionner un produit</option>
        {#each mockProducts as product}
            <option value={product.id}>{product.name}</option>
        {/each}
    </select>
{/snippet}

{#snippet numberInput(name: string)}
    {#if name === "amount"}
        <input
            id={name}
            type="number"
            {name}
            bind:value={formData.amount}
            min="0"
            step="1"
            class="h-input rounded-card-sm border-border-input bg-background placeholder:text-foreground-alt/50 hover:border-dark-40 focus:ring-foreground focus:ring-offset-background focus:outline-hidden inline-flex w-full items-center border px-4 text-base focus:ring-2 focus:ring-offset-2 sm:text-sm transition-all"
            placeholder="Ex: 5000 pour 50.00€"
            required
        />
    {/if}
{/snippet}

{#snippet currencySelect(name: string)}
    <select
        id={name}
        {name}
        bind:value={formData.currency}
        class="h-input rounded-card-sm border-border-input bg-background hover:border-dark-40 focus:ring-foreground focus:ring-offset-background focus:outline-hidden inline-flex w-full items-center border px-4 text-base focus:ring-2 focus:ring-offset-2 sm:text-sm transition-all"
        required
    >
        <option value="eur">EUR (€)</option>
        <option value="usd">USD ($)</option>
    </select>
{/snippet}

{#snippet intervalSelect(name: string)}
    <select
        id={name}
        {name}
        bind:value={formData.interval}
        class="h-input rounded-card-sm border-border-input bg-background hover:border-dark-40 focus:ring-foreground focus:ring-offset-background focus:outline-hidden inline-flex w-full items-center border px-4 text-base focus:ring-2 focus:ring-offset-2 sm:text-sm transition-all"
        required
    >
        <option value="one_time">Paiement unique</option>
        <option value="month">Mensuel</option>
        <option value="year">Annuel</option>
    </select>
{/snippet}

{#snippet textInput(name: string)}
    {#if name === "stripePriceId"}
        <input
            id={name}
            type="text"
            {name}
            bind:value={formData.stripePriceId}
            class="h-input rounded-card-sm border-border-input bg-background placeholder:text-foreground-alt/50 hover:border-dark-40 focus:ring-foreground focus:ring-offset-background focus:outline-hidden inline-flex w-full items-center border px-4 text-base focus:ring-2 focus:ring-offset-2 sm:text-sm transition-all"
            placeholder="price_..."
            required
        />
    {/if}
{/snippet}

{#snippet statusCheckbox(name: string)}
    {#if name === "isActive"}
        <label class="flex items-center gap-2 cursor-pointer">
            <input
                type="checkbox"
                id={name}
                {name}
                bind:checked={formData.isActive}
                class="w-4 h-4 text-dark border-gray-300 rounded focus:ring-dark"
            />
            <span class="text-sm">Actif</span>
        </label>
    {/if}
{/snippet}
