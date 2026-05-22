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
        Ticket,
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

    // Coupon type from API
    type Coupon = {
        id: string;
        stripeCouponId: string;
        name: string;
        percentOff?: number;
        amountOff?: number;
        currency?: string;
        duration: "once" | "repeating" | "forever";
        durationInMonths?: number;
        maxRedemptions?: number;
        timesRedeemed: number;
        isValid: boolean;
        redeemBy?: string;
        createdAt: string;
        updatedAt: string;
        metadata?: Record<string, string | undefined>;
    };

    interface Props {
        data: PageData;
    }

    let { data }: Props = $props();

    // Use coupons from server data
    let coupons = $state<Coupon[]>(data.coupons || []);

    // Initialize superforms for create, update, and delete
    const {
        form: createForm,
        errors: createErrors,
        enhance: createEnhance,
    } = superForm(data.createCouponForm, {
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
    } = superForm(data.updateCouponForm, {
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
    } = superForm(data.deleteCouponForm, {
        resetForm: false,
        onUpdated({ form }) {
            if (form.valid) {
                deleteDialogOpen = false;
            }
        },
    });

    // Status filter state (single select)
    const ALL_STATUSES = "all";
    type StatusFilter = "all" | "valid" | "expired";
    let selectedStatus = $state<StatusFilter>(ALL_STATUSES);

    // Sort state
    type SortField = "name" | "value" | "duration" | "status";
    let sortField = $state<SortField | null>(null);
    let sortDirection = $state<"asc" | "desc">("asc");

    // Helper function to check if a coupon is currently valid
    function isCouponCurrentlyValid(coupon: Coupon): boolean {
        if (!coupon.isValid) return false;
        if (coupon.redeemBy) {
            return new Date(coupon.redeemBy) > new Date();
        }
        return true;
    }

    // Filtered coupons based on selected status
    const filteredCoupons = $derived(() => {
        if (selectedStatus === ALL_STATUSES) {
            return coupons;
        } else if (selectedStatus === "valid") {
            return coupons.filter((c) => isCouponCurrentlyValid(c));
        } else {
            return coupons.filter((c) => !isCouponCurrentlyValid(c));
        }
    });

    // Sorted and filtered coupons
    const displayedCoupons = $derived(() => {
        let result = [...filteredCoupons()];

        if (sortField) {
            result.sort((a, b) => {
                let compareValue = 0;

                if (sortField === "name") {
                    compareValue = a.name.localeCompare(b.name);
                } else if (sortField === "value") {
                    const valueA = a.percentOff ?? (a.amountOff ?? 0);
                    const valueB = b.percentOff ?? (b.amountOff ?? 0);
                    compareValue = valueA - valueB;
                } else if (sortField === "duration") {
                    const durationOrder = { once: 0, repeating: 1, forever: 2 };
                    compareValue = durationOrder[a.duration] - durationOrder[b.duration];
                } else if (sortField === "status") {
                    const validA = isCouponCurrentlyValid(a);
                    const validB = isCouponCurrentlyValid(b);
                    compareValue = validA === validB ? 0 : validA ? -1 : 1;
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

    // Currently selected coupon for edit/delete
    let selectedCoupon: Coupon | null = $state(null);

    // Discount type for form (separate for create and edit)
    type DiscountType = "percentage" | "fixed_amount";
    let createDiscountType = $state<DiscountType>("percentage");
    let editDiscountType = $state<DiscountType>("percentage");

    function openCreateDialog() {
        createDiscountType = "percentage";
        $createForm.percentOff = undefined;
        $createForm.amountOff = undefined;
        $createForm.currency = undefined;
        $createForm.durationInMonths = undefined;
        $createForm.maxRedemptions = undefined;
        $createForm.redeemBy = undefined;
        createDialogOpen = true;
    }

    function openEditDialog(coupon: Coupon) {
        selectedCoupon = coupon;
        editDiscountType = coupon.percentOff !== undefined ? "percentage" : "fixed_amount";
        $updateForm.id = coupon.id;
        $updateForm.name = coupon.name;
        editDialogOpen = true;
    }

    function openDeleteDialog(coupon: Coupon) {
        selectedCoupon = coupon;
        $deleteForm.id = coupon.id;
        deleteDialogOpen = true;
    }

    function formatDiscount(coupon: Coupon): string {
        if (coupon.percentOff !== undefined) {
            return `${coupon.percentOff}%`;
        } else if (coupon.amountOff !== undefined && coupon.currency) {
            const value = coupon.amountOff / 100;
            const symbol = coupon.currency === "eur" ? "€" : coupon.currency === "usd" ? "$" : coupon.currency.toUpperCase();
            return `${value.toFixed(2)} ${symbol}`;
        }
        return "-";
    }

    function formatDuration(coupon: Coupon): string {
        if (coupon.duration === "once") {
            return "Une fois";
        } else if (coupon.duration === "repeating") {
            return `Récurrent (${coupon.durationInMonths} mois)`;
        } else {
            return "Toujours";
        }
    }

    function formatUsage(coupon: Coupon): string {
        if (coupon.maxRedemptions) {
            return `${coupon.timesRedeemed} / ${coupon.maxRedemptions}`;
        }
        return `${coupon.timesRedeemed} / Illimité`;
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
    <!-- Filter & Coupons Table -->
    <div class="p-4 md:p-8">
        <!-- Header with Create Button -->
        <div class="flex items-center justify-between mb-6">
            <h2 class="text-lg font-semibold">Liste des coupons</h2>
            <Button.Root
                type="button"
                class="cursor-pointer hidden md:flex"
                onclick={openCreateDialog}
            >
                <div
                    class="flex gap-2 items-center py-2 px-4 bg-dark text-white rounded-input hover:bg-dark/95 transition-all shadow-mini"
                >
                    <Plus size={18} />
                    <span class="text-sm font-medium">Nouveau Coupon</span>
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
                                <Combobox.Item value="valid" label="Valide">
                                    {#snippet children({ selected })}
                                        <div
                                            class="flex items-center justify-between px-3 py-2 rounded-md hover:bg-dark-04 cursor-pointer"
                                        >
                                            <span class="text-sm">Valide</span>
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

            <!-- Coupon count -->
            <div class="flex items-center gap-2 text-sm text-foreground-alt">
                {#if selectedStatus !== ALL_STATUSES}
                    <span class="hidden md:inline">•</span>
                    <span
                        >{displayedCoupons().length} coupons</span
                    >
                    <button
                        type="button"
                        onclick={() => (selectedStatus = ALL_STATUSES)}
                        class="text-destructive hover:underline"
                    >
                        Réinitialiser
                    </button>
                {:else}
                    <span>{coupons.length} coupons</span>
                {/if}
            </div>
        </div>

        <!-- Selected status chip -->
        {#if selectedStatus !== ALL_STATUSES}
            <div class="flex flex-wrap gap-2 mb-6">
                <div
                    class="inline-flex items-center gap-1 px-3 py-1 bg-dark-04 border border-border-input rounded-full text-sm"
                >
                    <span>{selectedStatus === "valid" ? "Valide" : "Expiré"}</span>
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

        <!-- Coupons Table -->
        {#if displayedCoupons().length === 0}
            <div
                class="flex flex-col items-center justify-center py-16 text-center"
            >
                <div
                    class="w-16 h-16 rounded-full bg-dark-04 flex items-center justify-center mb-4"
                >
                    <Ticket size={32} class="text-dark-400" />
                </div>
                <h3 class="text-lg font-medium mb-2">
                    {coupons.length === 0 || selectedStatus === ALL_STATUSES
                        ? "Aucun coupon"
                        : "Aucun résultat"}
                </h3>
                <p class="text-sm text-foreground-alt mb-6 max-w-sm">
                    {coupons.length === 0
                        ? "Commencez par créer votre premier coupon pour offrir des réductions à vos clients."
                        : "Aucun coupon ne correspond au statut sélectionné."}
                </p>
                {#if coupons.length === 0}
                    <Button.Root
                        type="button"
                        class="cursor-pointer"
                        onclick={openCreateDialog}
                    >
                        <div
                            class="flex gap-2 items-center py-2 px-4 bg-dark text-white rounded-input hover:bg-dark/95 transition-all shadow-mini"
                        >
                            <Plus size={18} />
                            <span class="text-sm font-medium">Créer un coupon</span>
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
                                onclick={() => toggleSort("name")}
                            >
                                <div class="flex items-center gap-2">
                                    <span>Nom</span>
                                    {#if sortField === "name"}
                                        <ArrowUpDown size={14} class={sortDirection === "desc" ? "rotate-180" : ""} />
                                    {:else}
                                        <ArrowUpDown size={14} class="opacity-30" />
                                    {/if}
                                </div>
                            </th>
                            <th
                                class="px-4 py-3 text-left text-sm font-semibold text-foreground cursor-pointer hover:bg-dark-10 transition-colors border-b border-r border-border-card"
                                onclick={() => toggleSort("value")}
                            >
                                <div class="flex items-center gap-2">
                                    <span>Remise</span>
                                    {#if sortField === "value"}
                                        <ArrowUpDown size={14} class={sortDirection === "desc" ? "rotate-180" : ""} />
                                    {:else}
                                        <ArrowUpDown size={14} class="opacity-30" />
                                    {/if}
                                </div>
                            </th>
                            <th
                                class="px-4 py-3 text-left text-sm font-semibold text-foreground cursor-pointer hover:bg-dark-10 transition-colors border-b border-r border-border-card"
                                onclick={() => toggleSort("duration")}
                            >
                                <div class="flex items-center gap-2">
                                    <span>Durée</span>
                                    {#if sortField === "duration"}
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
                                class="px-4 py-3 text-left text-sm font-semibold text-foreground border-b border-r border-border-card"
                            >
                                Stripe Coupon ID
                            </th>
                            <th
                                class="px-4 py-3 text-right text-sm font-semibold text-foreground border-b border-border-card"
                            >
                                Actions
                            </th>
                        </tr>
                    </thead>
                    <tbody>
                        {#each displayedCoupons() as coupon (coupon.id)}
                            <tr class="hover:bg-dark-04 transition-colors">
                                <td class="px-4 py-3 text-sm font-medium text-foreground border-b border-r border-border-card">
                                    {coupon.name}
                                </td>
                                <td
                                    class="px-4 py-3 text-sm font-medium text-foreground border-b border-r border-border-card"
                                >
                                    {formatDiscount(coupon)}
                                </td>
                                <td class="px-4 py-3 text-sm text-foreground border-b border-r border-border-card">
                                    {formatDuration(coupon)}
                                </td>
                                <td class="px-4 py-3 text-sm text-foreground border-b border-r border-border-card">
                                    {formatUsage(coupon)}
                                </td>
                                <td class="px-4 py-3 border-b border-r border-border-card">
                                    <span
                                        class="px-2 py-1 text-xs font-medium rounded-md {isCouponCurrentlyValid(coupon)
                                            ? 'bg-green-100 text-green-800'
                                            : 'bg-gray-100 text-gray-800'}"
                                    >
                                        {isCouponCurrentlyValid(coupon) ? "Valide" : "Expiré"}
                                    </span>
                                </td>
                                <td
                                    class="px-4 py-3 text-sm text-foreground-alt font-mono border-b border-r border-border-card"
                                >
                                    {coupon.stripeCouponId}
                                </td>
                                <td class="px-4 py-3 border-b border-border-card">
                                    <div class="flex gap-2 justify-end">
                                        <Button.Root
                                            type="button"
                                            class="cursor-pointer"
                                            onclick={() => openEditDialog(coupon)}
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
                                            onclick={() => openDeleteDialog(coupon)}
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

<!-- Create Coupon Dialog/Drawer -->
{#if isMobile}
    <Drawer bind:isOpen={createDialogOpen}>
        <div
            class="sticky top-0 bg-white pb-4 border-b border-border-card -mx-4 px-4 -mt-4 pt-4 z-10"
        >
            <div class="flex items-center justify-between mb-2">
                <h2 class="text-xl font-semibold tracking-tight">
                    Créer un coupon
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
                Remplissez les détails ci-dessous pour créer un nouveau coupon.
            </p>
        </div>

        <form method="POST" action="?/createCoupon" use:createEnhance class="grid gap-4 pt-8">
            <div class="grid grid-cols-1 gap-4 w-full pb-4">
                {@render createField("name", "Nom du coupon")}
                {@render createField("discountType", "Type de remise")}
                {#if createDiscountType === "percentage"}
                    {@render createField("percentOff", "Pourcentage de réduction")}
                {:else}
                    {@render createField("amountOff", "Montant (en centimes)")}
                    {@render createField("currency", "Devise")}
                {/if}
                {@render createField("duration", "Durée")}
                {#if $createForm.duration === "repeating"}
                    {@render createField("durationInMonths", "Durée en mois")}
                {/if}
                {@render createField("maxRedemptions", "Max. utilisations (0 = illimité)")}
                {@render createField("redeemBy", "Valable jusqu'au")}
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
                    Créer un coupon
                </Dialog.Title>
                <Dialog.Description class="text-foreground-alt !mt-1 text-sm">
                    Remplissez les détails ci-dessous pour créer un nouveau
                    coupon.
                </Dialog.Description>

                <Separator.Root class="bg-muted mx-5 !mb-2 !mt-5 block h-px" />

                <form method="POST" action="?/createCoupon" use:createEnhance class="grid gap-4">
                    <div
                        class="grid grid-cols-[max-content_1fr] gap-4 w-full items-center pb-11 pt-7"
                    >
                        {@render createField("name", "Nom du coupon")}
                        {@render createField("discountType", "Type de remise")}
                        {#if createDiscountType === "percentage"}
                            {@render createField("percentOff", "Pourcentage de réduction")}
                        {:else}
                            {@render createField("amountOff", "Montant (centimes)")}
                            {@render createField("currency", "Devise")}
                        {/if}
                        {@render createField("duration", "Durée")}
                        {#if $createForm.duration === "repeating"}
                            {@render createField("durationInMonths", "Durée en mois")}
                        {/if}
                        {@render createField("maxRedemptions", "Max. utilisations (0 = illimité)")}
                        {@render createField("redeemBy", "Valable jusqu'au")}
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

<!-- Edit Coupon Dialog/Drawer -->
{#if isMobile}
    <Drawer bind:isOpen={editDialogOpen}>
        <div
            class="sticky top-0 bg-white pb-4 border-b border-border-card -mx-4 px-4 -mt-4 pt-4 z-10"
        >
            <div class="flex items-center justify-between mb-2">
                <h2 class="text-xl font-semibold tracking-tight">
                    Modifier le coupon
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
                Mettez à jour les détails du coupon.
            </p>
        </div>

        <form method="POST" action="?/updateCoupon" use:updateEnhance class="grid gap-4 pt-8">
            {#if selectedCoupon}
                <div class="bg-dark-04 p-4 rounded-card mt-4">
                    <p class="text-sm text-foreground-alt mb-2">Coupon sélectionné:</p>
                    <p class="text-sm font-medium">{selectedCoupon.name}</p>
                    <p class="text-xs text-foreground-alt">{formatDiscount(selectedCoupon)} - {formatDuration(selectedCoupon)}</p>
                </div>
            {/if}
            <input type="hidden" name="id" bind:value={$updateForm.id} />
            <div class="grid grid-cols-1 gap-4 w-full pb-4">
                {@render updateField("name", "Nom du coupon")}
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
                    Modifier le coupon
                </Dialog.Title>
                <Dialog.Description class="text-foreground-alt !mt-1 text-sm">
                    Mettez à jour les détails du coupon.
                </Dialog.Description>

                <Separator.Root class="bg-muted mx-5 !mb-2 !mt-5 block h-px" />

                <form method="POST" action="?/updateCoupon" use:updateEnhance class="grid gap-4">
                    {#if selectedCoupon}
                        <div class="bg-dark-04 p-4 rounded-card mt-4">
                            <p class="text-sm text-foreground-alt mb-2">Coupon sélectionné:</p>
                            <p class="text-sm font-medium">{selectedCoupon.name}</p>
                            <p class="text-xs text-foreground-alt">{formatDiscount(selectedCoupon)} - {formatDuration(selectedCoupon)}</p>
                        </div>
                    {/if}
                    <input type="hidden" name="id" bind:value={$updateForm.id} />
                    <div
                        class="grid grid-cols-[max-content_1fr] gap-4 w-full items-center pb-11 pt-7"
                    >
                        {@render updateField("name", "Nom du coupon")}
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

<!-- Delete Coupon Dialog/Drawer -->
{#if isMobile}
    <Drawer bind:isOpen={deleteDialogOpen}>
        <div
            class="sticky top-0 bg-white pb-4 border-b border-border-card -mx-4 px-4 -mt-4 pt-4 z-10"
        >
            <div class="flex items-center justify-between mb-2">
                <h2 class="text-xl font-semibold tracking-tight">
                    Supprimer le coupon
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
                Êtes-vous sûr de vouloir supprimer ce coupon ? Cette action est
                irréversible.
            </p>
        </div>

        <form method="POST" action="?/deleteCoupon" use:deleteEnhance class="pt-8">
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
                    Supprimer le coupon
                </Dialog.Title>
                <Dialog.Description class="text-foreground-alt !mt-2 text-sm">
                    Êtes-vous sûr de vouloir supprimer ce coupon ? Cette action
                    est irréversible.
                </Dialog.Description>

                <form method="POST" action="?/deleteCoupon" use:deleteEnhance class="mt-8">
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
        {#if name === "name"}
            <input
                id={name}
                type="text"
                {name}
                bind:value={$createForm.name}
                class="h-input rounded-card-sm border-border-input bg-background placeholder:text-foreground-alt/50 hover:border-dark-40 focus:ring-foreground focus:ring-offset-background focus:outline-hidden inline-flex w-full items-center border px-4 text-base focus:ring-2 focus:ring-offset-2 sm:text-sm transition-all"
                placeholder="Ex: Bienvenue 2025"
                required
            />
        {:else if name === "discountType"}
            <div class="flex gap-4">
                <label class="flex items-center gap-2 cursor-pointer">
                    <input
                        type="radio"
                        name="discountType"
                        bind:group={createDiscountType}
                        value="percentage"
                        class="w-4 h-4 text-dark border-gray-300 focus:ring-dark"
                    />
                    <span class="text-sm">Pourcentage</span>
                </label>
                <label class="flex items-center gap-2 cursor-pointer">
                    <input
                        type="radio"
                        name="discountType"
                        bind:group={createDiscountType}
                        value="fixed_amount"
                        class="w-4 h-4 text-dark border-gray-300 focus:ring-dark"
                    />
                    <span class="text-sm">Montant fixe</span>
                </label>
            </div>
        {:else if name === "percentOff"}
            <input
                id={name}
                type="number"
                {name}
                bind:value={$createForm.percentOff}
                min="0.1"
                max="100"
                step="0.01"
                class="h-input rounded-card-sm border-border-input bg-background placeholder:text-foreground-alt/50 hover:border-dark-40 focus:ring-foreground focus:ring-offset-background focus:outline-hidden inline-flex w-full items-center border px-4 text-base focus:ring-2 focus:ring-offset-2 sm:text-sm transition-all"
                placeholder="Ex: 25 pour 25%"
            />
        {:else if name === "amountOff"}
            <input
                id={name}
                type="number"
                {name}
                bind:value={$createForm.amountOff}
                min="1"
                step="1"
                class="h-input rounded-card-sm border-border-input bg-background placeholder:text-foreground-alt/50 hover:border-dark-40 focus:ring-foreground focus:ring-offset-background focus:outline-hidden inline-flex w-full items-center border px-4 text-base focus:ring-2 focus:ring-offset-2 sm:text-sm transition-all"
                placeholder="Ex: 1000 pour 10.00€"
            />
        {:else if name === "currency"}
            <select
                id={name}
                {name}
                bind:value={$createForm.currency}
                class="h-input rounded-card-sm border-border-input bg-background hover:border-dark-40 focus:ring-foreground focus:ring-offset-background focus:outline-hidden inline-flex w-full items-center border px-4 text-base focus:ring-2 focus:ring-offset-2 sm:text-sm transition-all"
            >
                <option value="eur">EUR (€)</option>
                <option value="usd">USD ($)</option>
            </select>
        {:else if name === "duration"}
            <select
                id={name}
                {name}
                bind:value={$createForm.duration}
                class="h-input rounded-card-sm border-border-input bg-background hover:border-dark-40 focus:ring-foreground focus:ring-offset-background focus:outline-hidden inline-flex w-full items-center border px-4 text-base focus:ring-2 focus:ring-offset-2 sm:text-sm transition-all"
                required
            >
                <option value="once">Une fois</option>
                <option value="repeating">Récurrent</option>
                <option value="forever">Toujours</option>
            </select>
        {:else if name === "durationInMonths"}
            <input
                id={name}
                type="number"
                {name}
                bind:value={$createForm.durationInMonths}
                min="1"
                step="1"
                class="h-input rounded-card-sm border-border-input bg-background placeholder:text-foreground-alt/50 hover:border-dark-40 focus:ring-foreground focus:ring-offset-background focus:outline-hidden inline-flex w-full items-center border px-4 text-base focus:ring-2 focus:ring-offset-2 sm:text-sm transition-all"
                placeholder="Ex: 3"
            />
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
        {:else if name === "redeemBy"}
            <input
                id={name}
                type="datetime-local"
                {name}
                bind:value={$createForm.redeemBy}
                class="h-input rounded-card-sm border-border-input bg-background placeholder:text-foreground-alt/50 hover:border-dark-40 focus:ring-foreground focus:ring-offset-background focus:outline-hidden inline-flex w-full items-center border px-4 text-base focus:ring-2 focus:ring-offset-2 sm:text-sm transition-all"
            />
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
        {#if name === "name"}
            <input
                id={name}
                type="text"
                {name}
                bind:value={$updateForm.name}
                class="h-input rounded-card-sm border-border-input bg-background placeholder:text-foreground-alt/50 hover:border-dark-40 focus:ring-foreground focus:ring-offset-background focus:outline-hidden inline-flex w-full items-center border px-4 text-base focus:ring-2 focus:ring-offset-2 sm:text-sm transition-all"
                placeholder="Ex: Bienvenue 2025"
            />
        {/if}
        {#if ($updateErrors as Record<string, unknown>)[name]}
            <p class="text-xs text-destructive mt-1">{($updateErrors as Record<string, unknown>)[name]}</p>
        {/if}
    </div>
{/snippet}
