<script lang="ts">
    import { Dialog, Button, Label, Separator, Combobox } from "bits-ui";
    import {
        Plus,
        Pencil,
        Trash2,
        X,
        Filter,
        Check,
        ChevronsUpDown,
        Tag,
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

    // PromotionCode type from API
    type PromotionCode = {
        id: string;
        stripePromotionId: string;
        couponId: string;
        code: string;
        active: boolean;
        maxRedemptions?: number;
        timesRedeemed: number;
        expiresAt?: string;
        firstTimeTransaction: boolean;
        minimumAmount?: number;
        minimumAmountCurrency?: string;
        createdAt: string;
        updatedAt: string;
        metadata?: Record<string, string | undefined>;
    };

    // Coupon type for lookups
    type Coupon = {
        id: string;
        name: string;
        percentOff?: number;
        amountOff?: number;
        currency?: string;
    };

    interface Props {
        data: PageData;
    }

    let { data }: Props = $props();

    // Use promotion codes from server data
    let promotionCodes = $state<PromotionCode[]>(data.promotionCodes || []);
    let coupons = $state<Coupon[]>(data.coupons || []);

    // Initialize superforms for create, update, and delete
    const {
        form: createForm,
        errors: createErrors,
        enhance: createEnhance,
    } = superForm(data.createPromotionCodeForm, {
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
    } = superForm(data.updatePromotionCodeForm, {
        resetForm: false,
        onUpdated({ form }) {
            if (form.valid) {
                editDialogOpen = false;
            }
        },
    });

    const {
        form: deleteForm,
        errors: deleteErrors,
        enhance: deleteEnhance,
    } = superForm(data.deletePromotionCodeForm, {
        resetForm: false,
        onUpdated({ form }) {
            if (form.valid) {
                deleteDialogOpen = false;
            }
        },
    });

    // Status filter state (single select)
    const ALL_STATUSES = "all";
    type StatusFilter = "all" | "active" | "expired";
    let selectedStatus = $state<StatusFilter>(ALL_STATUSES);

    // Sort state
    type SortField = "code" | "coupon" | "status" | "expiresAt";
    let sortField = $state<SortField | null>(null);
    let sortDirection = $state<"asc" | "desc">("asc");

    // Helper function to check if a promotion code is currently active
    function isPromotionCodeActive(promo: PromotionCode): boolean {
        if (!promo.active) return false;
        if (promo.expiresAt) {
            return new Date(promo.expiresAt) > new Date();
        }
        return true;
    }

    // Helper function to get coupon name by ID
    function getCouponName(couponId: string): string {
        return coupons.find((c) => c.id === couponId)?.name || "Inconnu";
    }

    // Helper function to get coupon discount info
    function getCouponDiscount(couponId: string): string {
        const coupon = coupons.find((c) => c.id === couponId);
        if (!coupon) return "-";
        if (coupon.percentOff !== undefined) {
            return `${coupon.percentOff}%`;
        } else if (coupon.amountOff !== undefined && coupon.currency) {
            const value = coupon.amountOff / 100;
            const symbol = coupon.currency === "eur" ? "€" : coupon.currency === "usd" ? "$" : coupon.currency.toUpperCase();
            return `${value.toFixed(2)} ${symbol}`;
        }
        return "-";
    }

    // Filtered promotion codes based on selected status
    const filteredPromotionCodes = $derived(() => {
        if (selectedStatus === ALL_STATUSES) {
            return promotionCodes;
        } else if (selectedStatus === "active") {
            return promotionCodes.filter((p) => isPromotionCodeActive(p));
        } else {
            return promotionCodes.filter((p) => !isPromotionCodeActive(p));
        }
    });

    // Sorted and filtered promotion codes
    const displayedPromotionCodes = $derived(() => {
        let result = [...filteredPromotionCodes()];

        if (sortField) {
            result.sort((a, b) => {
                let compareValue = 0;

                if (sortField === "code") {
                    compareValue = a.code.localeCompare(b.code);
                } else if (sortField === "coupon") {
                    const couponA = getCouponName(a.couponId);
                    const couponB = getCouponName(b.couponId);
                    compareValue = couponA.localeCompare(couponB);
                } else if (sortField === "status") {
                    const activeA = isPromotionCodeActive(a);
                    const activeB = isPromotionCodeActive(b);
                    compareValue = activeA === activeB ? 0 : activeA ? -1 : 1;
                } else if (sortField === "expiresAt") {
                    if (!a.expiresAt && !b.expiresAt) {
                        compareValue = 0;
                    } else if (!a.expiresAt) {
                        compareValue = 1;
                    } else if (!b.expiresAt) {
                        compareValue = -1;
                    } else {
                        compareValue = new Date(a.expiresAt).getTime() - new Date(b.expiresAt).getTime();
                    }
                }

                return sortDirection === "asc" ? compareValue : -compareValue;
            });
        }

        return result;
    });

    // Dialog states
    let createDialogOpen = $state(false);
    let editDialogOpen = $state(false);
    let deleteDialogOpen = $state(false);

    // Currently selected promotion code for edit/delete
    let selectedPromo: PromotionCode | null = $state(null);

    function openCreateDialog() {
        $createForm.couponId = "";
        $createForm.code = "";
        $createForm.maxRedemptions = undefined;
        $createForm.expiresAt = undefined;
        $createForm.firstTimeTransaction = false;
        $createForm.minimumAmount = undefined;
        $createForm.minimumAmountCurrency = undefined;
        createDialogOpen = true;
    }

    function openEditDialog(promo: PromotionCode) {
        selectedPromo = promo;
        $updateForm.id = promo.id;
        $updateForm.active = promo.active;
        editDialogOpen = true;
    }

    function openDeleteDialog(promo: PromotionCode) {
        selectedPromo = promo;
        $deleteForm.id = promo.id;
        deleteDialogOpen = true;
    }

    function formatUsage(promo: PromotionCode): string {
        if (promo.maxRedemptions) {
            return `${promo.timesRedeemed} / ${promo.maxRedemptions}`;
        }
        return `${promo.timesRedeemed} / Illimité`;
    }

    function formatExpiration(promo: PromotionCode): string {
        if (!promo.expiresAt) return "-";
        const date = new Date(promo.expiresAt);
        return date.toLocaleDateString("fr-FR", {
            day: "2-digit",
            month: "short",
            year: "numeric",
        });
    }

    function formatMinimumAmount(promo: PromotionCode): string {
        if (!promo.minimumAmount) return "-";
        const value = promo.minimumAmount / 100;
        const currency = promo.minimumAmountCurrency || "eur";
        const symbol = currency === "eur" ? "€" : currency === "usd" ? "$" : currency.toUpperCase();
        return `${value.toFixed(2)} ${symbol}`;
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
    <!-- Filter & Promotion Codes Table -->
    <div class="p-4 md:p-8">
        <!-- Header with Create Button -->
        <div class="flex items-center justify-between mb-6">
            <h2 class="text-lg font-semibold">Liste des codes promotion</h2>
            <Button.Root
                type="button"
                class="cursor-pointer hidden md:flex"
                onclick={openCreateDialog}
            >
                <div
                    class="flex gap-2 items-center py-2 px-4 bg-dark text-white rounded-input hover:bg-dark/95 transition-all shadow-mini"
                >
                    <Plus size={18} />
                    <span class="text-sm font-medium">Nouveau Code</span>
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
            <Combobox.Root type="single" bind:value={selectedStatus}>
                <div class="relative w-full md:w-auto md:min-w-[280px]">
                    <Filter
                        class="text-muted-foreground absolute start-3 top-1/2 size-4 -translate-y-1/2"
                    />
                    <Combobox.Input
                        class="h-input rounded-card-sm border-border-input bg-background placeholder:text-foreground-alt/50 hover:border-dark-40 focus:ring-foreground focus:ring-offset-background focus:outline-hidden inline-flex w-full items-center border px-9 text-sm focus:ring-2 focus:ring-offset-2 transition-all"
                        placeholder="Filtrer par statut"
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
                            <Combobox.Item value={ALL_STATUSES} label="Tous les statuts">
                                {#snippet children({ selected })}
                                    <div
                                        class="flex items-center justify-between px-3 py-2 rounded-md hover:bg-dark-04 cursor-pointer"
                                    >
                                        <span class="text-sm"
                                            >Tous les statuts</span
                                        >
                                        {#if selected}
                                            <Check class="size-4" />
                                        {/if}
                                    </div>
                                {/snippet}
                            </Combobox.Item>
                            <Combobox.Group class="pt-1">
                                <Combobox.Item value="active" label="Actif">
                                    {#snippet children({ selected })}
                                        <div
                                            class="flex items-center justify-between px-3 py-2 rounded-md hover:bg-dark-04 cursor-pointer"
                                        >
                                            <span class="text-sm">Actif</span>
                                            {#if selected}
                                                <Check class="size-4" />
                                            {/if}
                                        </div>
                                    {/snippet}
                                </Combobox.Item>
                                <Combobox.Item value="expired" label="Expiré">
                                    {#snippet children({ selected })}
                                        <div
                                            class="flex items-center justify-between px-3 py-2 rounded-md hover:bg-dark-04 cursor-pointer"
                                        >
                                            <span class="text-sm">Expiré</span>
                                            {#if selected}
                                                <Check class="size-4" />
                                            {/if}
                                        </div>
                                    {/snippet}
                                </Combobox.Item>
                            </Combobox.Group>
                        </Combobox.Viewport>
                    </Combobox.Content>
                </Combobox.Portal>
            </Combobox.Root>

            <!-- Promotion code count -->
            <div class="flex items-center gap-2 text-sm text-foreground-alt">
                {#if selectedStatus !== ALL_STATUSES}
                    <span class="hidden md:inline">•</span>
                    <span
                        >{displayedPromotionCodes().length} codes</span
                    >
                    <button
                        type="button"
                        onclick={() => (selectedStatus = ALL_STATUSES)}
                        class="text-destructive hover:underline"
                    >
                        Réinitialiser
                    </button>
                {:else}
                    <span>{promotionCodes.length} codes</span>
                {/if}
            </div>
        </div>

        <!-- Selected status chip -->
        {#if selectedStatus !== ALL_STATUSES}
            <div class="flex flex-wrap gap-2 mb-6">
                <div
                    class="inline-flex items-center gap-1 px-3 py-1 bg-dark-04 border border-border-input rounded-full text-sm"
                >
                    <span>{selectedStatus === "active" ? "Actif" : "Expiré"}</span>
                    <button
                        type="button"
                        onclick={() => (selectedStatus = ALL_STATUSES)}
                        class="hover:text-destructive transition-colors"
                    >
                        <X size={14} />
                    </button>
                </div>
            </div>
        {/if}

        <!-- Promotion Codes Table -->
        {#if displayedPromotionCodes().length === 0}
            <div
                class="flex flex-col items-center justify-center py-16 text-center"
            >
                <div
                    class="w-16 h-16 rounded-full bg-dark-04 flex items-center justify-center mb-4"
                >
                    <Tag size={32} class="text-dark-400" />
                </div>
                <h3 class="text-lg font-medium mb-2">
                    {promotionCodes.length === 0 || selectedStatus === ALL_STATUSES
                        ? "Aucun code promotion"
                        : "Aucun résultat"}
                </h3>
                <p class="text-sm text-foreground-alt mb-6 max-w-sm">
                    {promotionCodes.length === 0
                        ? "Commencez par créer votre premier code promotion pour offrir des réductions à vos clients."
                        : "Aucun code promotion ne correspond au statut sélectionné."}
                </p>
                {#if promotionCodes.length === 0}
                    <Button.Root
                        type="button"
                        class="cursor-pointer"
                        onclick={openCreateDialog}
                    >
                        <div
                            class="flex gap-2 items-center py-2 px-4 bg-dark text-white rounded-input hover:bg-dark/95 transition-all shadow-mini"
                        >
                            <Plus size={18} />
                            <span class="text-sm font-medium">Créer un code</span>
                        </div>
                    </Button.Root>
                {:else}
                    <button
                        type="button"
                        onclick={() => (selectedStatus = ALL_STATUSES)}
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
                                onclick={() => toggleSort("code")}
                            >
                                <div class="flex items-center gap-2">
                                    <span>Code</span>
                                    {#if sortField === "code"}
                                        <ArrowUpDown size={14} class={sortDirection === "desc" ? "rotate-180" : ""} />
                                    {:else}
                                        <ArrowUpDown size={14} class="opacity-30" />
                                    {/if}
                                </div>
                            </th>
                            <th
                                class="px-4 py-3 text-left text-sm font-semibold text-foreground cursor-pointer hover:bg-dark-10 transition-colors border-b border-r border-border-card"
                                onclick={() => toggleSort("coupon")}
                            >
                                <div class="flex items-center gap-2">
                                    <span>Coupon associé</span>
                                    {#if sortField === "coupon"}
                                        <ArrowUpDown size={14} class={sortDirection === "desc" ? "rotate-180" : ""} />
                                    {:else}
                                        <ArrowUpDown size={14} class="opacity-30" />
                                    {/if}
                                </div>
                            </th>
                            <th
                                class="px-4 py-3 text-left text-sm font-semibold text-foreground border-b border-r border-border-card"
                            >
                                Utilisations
                            </th>
                            <th
                                class="px-4 py-3 text-left text-sm font-semibold text-foreground cursor-pointer hover:bg-dark-10 transition-colors border-b border-r border-border-card"
                                onclick={() => toggleSort("expiresAt")}
                            >
                                <div class="flex items-center gap-2">
                                    <span>Expiration</span>
                                    {#if sortField === "expiresAt"}
                                        <ArrowUpDown size={14} class={sortDirection === "desc" ? "rotate-180" : ""} />
                                    {:else}
                                        <ArrowUpDown size={14} class="opacity-30" />
                                    {/if}
                                </div>
                            </th>
                            <th
                                class="px-4 py-3 text-left text-sm font-semibold text-foreground cursor-pointer hover:bg-dark-10 transition-colors border-b border-r border-border-card"
                                onclick={() => toggleSort("status")}
                            >
                                <div class="flex items-center gap-2">
                                    <span>Statut</span>
                                    {#if sortField === "status"}
                                        <ArrowUpDown size={14} class={sortDirection === "desc" ? "rotate-180" : ""} />
                                    {:else}
                                        <ArrowUpDown size={14} class="opacity-30" />
                                    {/if}
                                </div>
                            </th>
                            <th
                                class="px-4 py-3 text-right text-sm font-semibold text-foreground border-b border-border-card"
                            >
                                Actions
                            </th>
                        </tr>
                    </thead>
                    <tbody>
                        {#each displayedPromotionCodes() as promo (promo.id)}
                            <tr class="hover:bg-dark-04 transition-colors">
                                <td class="px-4 py-3 border-b border-r border-border-card">
                                    <div class="flex flex-col">
                                        <span class="text-sm font-mono font-medium text-foreground">{promo.code}</span>
                                        {#if promo.firstTimeTransaction}
                                            <span class="text-xs text-foreground-alt mt-1">Nouveaux clients uniquement</span>
                                        {/if}
                                        {#if promo.minimumAmount}
                                            <span class="text-xs text-foreground-alt">Min. {formatMinimumAmount(promo)}</span>
                                        {/if}
                                    </div>
                                </td>
                                <td
                                    class="px-4 py-3 text-sm text-foreground border-b border-r border-border-card"
                                >
                                    <div class="flex flex-col">
                                        <span class="font-medium">{getCouponName(promo.couponId)}</span>
                                        <span class="text-xs text-foreground-alt">{getCouponDiscount(promo.couponId)}</span>
                                    </div>
                                </td>
                                <td class="px-4 py-3 text-sm text-foreground border-b border-r border-border-card">
                                    {formatUsage(promo)}
                                </td>
                                <td class="px-4 py-3 text-sm text-foreground border-b border-r border-border-card">
                                    {formatExpiration(promo)}
                                </td>
                                <td class="px-4 py-3 border-b border-r border-border-card">
                                    <span
                                        class="px-2 py-1 text-xs font-medium rounded-md {isPromotionCodeActive(promo)
                                            ? 'bg-green-100 text-green-800'
                                            : 'bg-gray-100 text-gray-800'}"
                                    >
                                        {isPromotionCodeActive(promo) ? "Actif" : "Expiré"}
                                    </span>
                                </td>
                                <td class="px-4 py-3 border-b border-border-card">
                                    <div class="flex gap-2 justify-end">
                                        <Button.Root
                                            type="button"
                                            class="cursor-pointer"
                                            onclick={() => openEditDialog(promo)}
                                        >
                                            <div
                                                class="flex items-center justify-center p-2 border border-border-input rounded-md hover:bg-dark-04 transition-all"
                                            >
                                                <Pencil size={14} />
                                            </div>
                                        </Button.Root>
                                        <Button.Root
                                            type="button"
                                            class="cursor-pointer"
                                            onclick={() => openDeleteDialog(promo)}
                                        >
                                            <div
                                                class="flex items-center justify-center p-2 border border-destructive/20 text-destructive rounded-md hover:bg-destructive/10 transition-all"
                                            >
                                                <Trash2 size={14} />
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

<!-- Create Promotion Code Dialog/Drawer -->
{#if isMobile}
    <Drawer bind:isOpen={createDialogOpen}>
        <div
            class="sticky top-0 bg-white pb-4 border-b border-border-card -mx-4 px-4 -mt-4 pt-4 z-10"
        >
            <div class="flex items-center justify-between mb-2">
                <h2 class="text-xl font-semibold tracking-tight">
                    Créer un code promotion
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
                Remplissez les détails ci-dessous pour créer un nouveau code promotion.
            </p>
        </div>

        <form method="POST" action="?/createPromotionCode" use:createEnhance class="grid gap-4 pt-8">
            <div class="grid grid-cols-1 gap-4 w-full pb-4">
                {@render createField("code", "Code promotion")}
                {@render createField("couponId", "Coupon associé")}
                {@render createField("maxRedemptions", "Max. utilisations (0 = illimité)")}
                {@render createField("expiresAt", "Date d'expiration")}
                {@render createField("minimumAmount", "Montant minimum (centimes)")}
                {#if $createForm.minimumAmount && $createForm.minimumAmount > 0}
                    {@render createField("minimumAmountCurrency", "Devise")}
                {/if}
                {@render createField("firstTimeTransaction", "Restrictions")}
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
                    Créer un code promotion
                </Dialog.Title>
                <Dialog.Description class="text-foreground-alt !mt-1 text-sm">
                    Remplissez les détails ci-dessous pour créer un nouveau
                    code promotion.
                </Dialog.Description>

                <Separator.Root class="bg-muted mx-5 !mb-2 !mt-5 block h-px" />

                <form method="POST" action="?/createPromotionCode" use:createEnhance class="grid gap-4">
                    <div
                        class="grid grid-cols-[max-content_1fr] gap-4 w-full items-center pb-11 pt-7"
                    >
                        {@render createField("code", "Code promotion")}
                        {@render createField("couponId", "Coupon associé")}
                        {@render createField("maxRedemptions", "Max. utilisations (0 = illimité)")}
                        {@render createField("expiresAt", "Date d'expiration")}
                        {@render createField("minimumAmount", "Montant minimum (centimes)")}
                        {#if $createForm.minimumAmount && $createForm.minimumAmount > 0}
                            {@render createField("minimumAmountCurrency", "Devise")}
                        {/if}
                        {@render createField("firstTimeTransaction", "Restrictions")}
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

<!-- Edit Promotion Code Dialog/Drawer -->
{#if isMobile}
    <Drawer bind:isOpen={editDialogOpen}>
        <div
            class="sticky top-0 bg-white pb-4 border-b border-border-card -mx-4 px-4 -mt-4 pt-4 z-10"
        >
            <div class="flex items-center justify-between mb-2">
                <h2 class="text-xl font-semibold tracking-tight">
                    Modifier le code promotion
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
                Mettez à jour les détails du code promotion.
            </p>
        </div>

        <form method="POST" action="?/updatePromotionCode" use:updateEnhance class="grid gap-4 pt-8">
            {#if selectedPromo}
                <div class="bg-dark-04 p-4 rounded-card mt-4">
                    <p class="text-sm text-foreground-alt mb-2">Code promotion sélectionné:</p>
                    <p class="text-sm font-medium font-mono">{selectedPromo.code}</p>
                    <p class="text-xs text-foreground-alt">{getCouponName(selectedPromo.couponId)} - {getCouponDiscount(selectedPromo.couponId)}</p>
                </div>
            {/if}
            <input type="hidden" name="id" bind:value={$updateForm.id} />
            <div class="grid grid-cols-1 gap-4 w-full pb-4">
                {@render updateField("active", "Statut actif")}
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
                    Modifier le code promotion
                </Dialog.Title>
                <Dialog.Description class="text-foreground-alt !mt-1 text-sm">
                    Mettez à jour les détails du code promotion.
                </Dialog.Description>

                <Separator.Root class="bg-muted mx-5 !mb-2 !mt-5 block h-px" />

                <form method="POST" action="?/updatePromotionCode" use:updateEnhance class="grid gap-4">
                    {#if selectedPromo}
                        <div class="bg-dark-04 p-4 rounded-card mt-4">
                            <p class="text-sm text-foreground-alt mb-2">Code promotion sélectionné:</p>
                            <p class="text-sm font-medium font-mono">{selectedPromo.code}</p>
                            <p class="text-xs text-foreground-alt">{getCouponName(selectedPromo.couponId)} - {getCouponDiscount(selectedPromo.couponId)}</p>
                        </div>
                    {/if}
                    <input type="hidden" name="id" bind:value={$updateForm.id} />
                    <div
                        class="grid grid-cols-[max-content_1fr] gap-4 w-full items-center pb-11 pt-7"
                    >
                        {@render updateField("active", "Statut actif")}
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

<!-- Delete Promotion Code Dialog/Drawer -->
{#if isMobile}
    <Drawer bind:isOpen={deleteDialogOpen}>
        <div
            class="sticky top-0 bg-white pb-4 border-b border-border-card -mx-4 px-4 -mt-4 pt-4 z-10"
        >
            <div class="flex items-center justify-between mb-2">
                <h2 class="text-xl font-semibold tracking-tight">
                    Supprimer le code promotion
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
                Êtes-vous sûr de vouloir supprimer ce code promotion ? Cette action est
                irréversible.
            </p>
        </div>

        <form method="POST" action="?/deletePromotionCode" use:deleteEnhance class="pt-8">
            <input type="hidden" name="id" bind:value={$deleteForm.id} />
            {#if $deleteErrors._errors}
                <p class="text-sm text-destructive mb-4">{$deleteErrors._errors}</p>
            {/if}
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
                    Supprimer le code promotion
                </Dialog.Title>
                <Dialog.Description class="text-foreground-alt !mt-2 text-sm">
                    Êtes-vous sûr de vouloir supprimer ce code promotion ? Cette action
                    est irréversible.
                </Dialog.Description>

                <form method="POST" action="?/deletePromotionCode" use:deleteEnhance class="mt-8">
                    <input type="hidden" name="id" bind:value={$deleteForm.id} />
                    {#if $deleteErrors._errors}
                        <p class="text-sm text-destructive mb-4">{$deleteErrors._errors}</p>
                    {/if}
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

<!-- Snippets for create form fields -->
{#snippet createField(name: string, label: string)}
    <Label.Root for={name} class="text-sm font-semibold">{label}</Label.Root>
    <div class="relative w-full">
        {#if name === "code"}
            <input
                id={name}
                type="text"
                {name}
                bind:value={$createForm.code}
                class="h-input rounded-card-sm border-border-input bg-background placeholder:text-foreground-alt/50 hover:border-dark-40 focus:ring-foreground focus:ring-offset-background focus:outline-hidden inline-flex w-full items-center border px-4 text-base font-mono uppercase focus:ring-2 focus:ring-offset-2 sm:text-sm transition-all"
                placeholder="Ex: BIENVENUE2025"
                required
            />
        {:else if name === "couponId"}
            <select
                id={name}
                {name}
                bind:value={$createForm.couponId}
                class="h-input rounded-card-sm border-border-input bg-background hover:border-dark-40 focus:ring-foreground focus:ring-offset-background focus:outline-hidden inline-flex w-full items-center border px-4 text-base focus:ring-2 focus:ring-offset-2 sm:text-sm transition-all"
                required
            >
                <option value="">Sélectionner un coupon</option>
                {#each coupons as coupon}
                    <option value={coupon.id}>{coupon.name}</option>
                {/each}
            </select>
        {:else if name === "maxRedemptions"}
            <input
                id={name}
                type="number"
                {name}
                bind:value={$createForm.maxRedemptions}
                min="0"
                step="1"
                class="h-input rounded-card-sm border-border-input bg-background placeholder:text-foreground-alt/50 hover:border-dark-40 focus:ring-foreground focus:ring-offset-background focus:outline-hidden inline-flex w-full items-center border px-4 text-base focus:ring-2 focus:ring-offset-2 sm:text-sm transition-all"
                placeholder="0 pour illimité"
            />
        {:else if name === "expiresAt"}
            <input
                id={name}
                type="datetime-local"
                {name}
                bind:value={$createForm.expiresAt}
                class="h-input rounded-card-sm border-border-input bg-background placeholder:text-foreground-alt/50 hover:border-dark-40 focus:ring-foreground focus:ring-offset-background focus:outline-hidden inline-flex w-full items-center border px-4 text-base focus:ring-2 focus:ring-offset-2 sm:text-sm transition-all"
            />
        {:else if name === "minimumAmount"}
            <input
                id={name}
                type="number"
                {name}
                bind:value={$createForm.minimumAmount}
                min="0"
                step="1"
                class="h-input rounded-card-sm border-border-input bg-background placeholder:text-foreground-alt/50 hover:border-dark-40 focus:ring-foreground focus:ring-offset-background focus:outline-hidden inline-flex w-full items-center border px-4 text-base focus:ring-2 focus:ring-offset-2 sm:text-sm transition-all"
                placeholder="0 pour aucun minimum"
            />
        {:else if name === "minimumAmountCurrency"}
            <select
                id={name}
                {name}
                bind:value={$createForm.minimumAmountCurrency}
                class="h-input rounded-card-sm border-border-input bg-background hover:border-dark-40 focus:ring-foreground focus:ring-offset-background focus:outline-hidden inline-flex w-full items-center border px-4 text-base focus:ring-2 focus:ring-offset-2 sm:text-sm transition-all"
            >
                <option value="eur">EUR (€)</option>
                <option value="usd">USD ($)</option>
            </select>
        {:else if name === "firstTimeTransaction"}
            <label class="flex items-center gap-2 cursor-pointer">
                <input
                    type="checkbox"
                    id={name}
                    {name}
                    bind:checked={$createForm.firstTimeTransaction}
                    class="w-4 h-4 text-dark border-gray-300 rounded focus:ring-dark"
                />
                <span class="text-sm">Nouveaux clients uniquement</span>
            </label>
        {/if}
        {#if ($createErrors as Record<string, unknown>)[name]}
            <p class="text-xs text-destructive mt-1">{($createErrors as Record<string, unknown>)[name]}</p>
        {/if}
    </div>
{/snippet}

<!-- Snippet for update form fields -->
{#snippet updateField(name: string, label: string)}
    <Label.Root for={name} class="text-sm font-semibold">{label}</Label.Root>
    <div class="relative w-full">
        {#if name === "active"}
            <label class="flex items-center gap-2 cursor-pointer">
                <input
                    type="checkbox"
                    id={name}
                    {name}
                    bind:checked={$updateForm.active}
                    class="w-4 h-4 text-dark border-gray-300 rounded focus:ring-dark"
                />
                <span class="text-sm">Actif</span>
            </label>
        {/if}
        {#if ($updateErrors as Record<string, unknown>)[name]}
            <p class="text-xs text-destructive mt-1">{($updateErrors as Record<string, unknown>)[name]}</p>
        {/if}
    </div>
{/snippet}
