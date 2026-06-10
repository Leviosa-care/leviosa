<script lang="ts">
    import { reveal } from "$lib/actions/reveal";
    import Button from "$lib/ui/Button.svelte";
    import type { PageProps } from "./$types";
    import type { Category, Product, Price, Partner, Availability } from "$lib/types/booking-flow.ts";

    let { data, form }: PageProps = $props();

    // --- State ---
    let loadingProducts = $state(false);
    let loadingPartners = $state(false);
    let loadingSlots = $state(false);
    let submitError = $state("");

    // Step 1 — Category
    let categories: Category[] = $state(data.categories ?? []);
    let selectedCategoryId = $state<string | null>(
        data.preselectedCategory?.id ?? null
    );
    let selectedCategoryName = $state<string | null>(
        data.preselectedCategory?.name ?? null
    );

    // Step 2 — Product
    let products = $state<Product[]>([]);
    let productPrices = $state<Record<string, number>>({});
    let selectedProductId = $state<string | null>(
        data.preselectedProduct?.id ?? null
    );
    let selectedProductName = $state<string | null>(
        data.preselectedProduct?.name ?? null
    );
    let selectedProductDuration = $state<number | null>(
        data.preselectedProduct?.duration ?? null
    );
    let selectedProductPrice = $state<string | null>(null);

    // Step 3 — Partner
    let partners = $state<Partner[]>([]);
    let selectedPartnerId = $state<string | null>(null);
    let selectedPartnerName = $state<string | null>(null);

    // Step 4 — Availability
    let availabilities = $state<Availability[]>([]);
    let selectedAvailabilityId = $state<string | null>(null);
    let selectedSlotStartTime = $state<string | null>(null);
    let selectedSlotDisplay = $state<string | null>(null);

    // Step 5 — Guest info
    let guestFirstName = $state("");
    let guestLastName = $state("");
    let guestEmail = $state("");
    let guestPhone = $state("");
    let guestFormErrors = $state<Record<string, string>>({});

    // --- Derived: which step to show ---
    let step = $derived.by(() => {
        if (!selectedCategoryId) return 1;
        if (!selectedProductId) return 2;
        if (!selectedPartnerId) return 3;
        if (!selectedAvailabilityId) return 4;
        if (!data.user && (!guestFirstName || !guestLastName || (!guestEmail && !guestPhone))) return 5;
        return 6;
    });

    // --- Pre-load on mount if product is preselected ---
    let initialized = $state(false);
    $effect(() => {
        if (initialized) return;
        initialized = true;

        if (data.preselectedProduct) {
            selectedCategoryId = data.preselectedCategory?.id ?? null;
            selectedCategoryName = data.preselectedCategory?.name ?? null;
            selectedProductId = data.preselectedProduct.id;
            selectedProductName = data.preselectedProduct.name;
            selectedProductDuration = data.preselectedProduct.duration;

            fetchProducts(data.preselectedCategory?.id);
            fetchPrice(data.preselectedProduct.id);

            if (data.preselectedPartners?.length > 0) {
                partners = data.preselectedPartners;
            } else {
                fetchPartners(data.preselectedProduct.id);
            }
        } else if (data.preselectedCategory) {
            fetchProducts(data.preselectedCategory.id);
        }
    });

    // --- Data fetching — all through SvelteKit API proxy routes ---
    async function fetchProducts(categoryId: string | null | undefined) {
        if (!categoryId) return;
        loadingProducts = true;
        try {
            const res = await fetch('/api/products');
            if (res.ok) {
                const all: Product[] = await res.json();
                products = all.filter((p) => p.category?.id === categoryId);
                for (const p of products) {
                    fetchPrice(p.id);
                }
            }
        } catch (e) {
            console.error("Failed to fetch products:", e);
        }
        loadingProducts = false;
    }

    async function fetchPrice(productId: string) {
        try {
            const res = await fetch(`/api/products/${productId}/prices`);
            if (res.ok) {
                const prices: Price[] = await res.json();
                const price = prices.find((p) => p.interval === "one_time") ?? prices[0];
                if (price) {
                    productPrices[productId] = price.amount;
                }
            }
        } catch {
            // Non-blocking
        }
    }

    async function fetchPartners(productId: string) {
        loadingPartners = true;
        try {
            const res = await fetch(`/api/partners/products/${productId}`);
            if (res.ok) {
                partners = await res.json();
            }
        } catch (e) {
            console.error("Failed to fetch partners:", e);
        }
        loadingPartners = false;
    }

    async function fetchAvailabilities(partnerId: string) {
        loadingSlots = true;
        try {
            const now = new Date().toISOString();
            const res = await fetch(
                `/api/partners/${partnerId}/availabilities?start_time=${encodeURIComponent(now)}&status=available`
            );
            if (res.ok) {
                const all: Availability[] = await res.json();
                availabilities = all.filter(
                    (a) => a.status === "available" && new Date(a.start_time) > new Date()
                );
            }
        } catch (e) {
            console.error("Failed to fetch availabilities:", e);
        }
        loadingSlots = false;
    }

    // --- Group availabilities by day ---
    let availabilitiesByDay = $derived.by(() => {
        const grouped: Record<string, Availability[]> = {};
        for (const a of availabilities) {
            const date = new Date(a.start_time);
            const key = date.toLocaleDateString("fr-FR", {
                weekday: "long",
                day: "numeric",
                month: "long",
            });
            if (!grouped[key]) grouped[key] = [];
            grouped[key].push(a);
        }
        for (const key of Object.keys(grouped)) {
            grouped[key].sort(
                (a, b) => new Date(a.start_time).getTime() - new Date(b.start_time).getTime()
            );
        }
        return grouped;
    });

    // --- Selection handlers ---
    function selectCategory(cat: Category) {
        selectedCategoryId = cat.id;
        selectedCategoryName = cat.name;
        selectedProductId = null;
        selectedProductName = null;
        selectedProductDuration = null;
        selectedProductPrice = null;
        selectedPartnerId = null;
        selectedPartnerName = null;
        selectedAvailabilityId = null;
        selectedSlotStartTime = null;
        selectedSlotDisplay = null;
        partners = [];
        availabilities = [];
        fetchProducts(cat.id);
    }

    function selectProduct(product: Product) {
        selectedProductId = product.id;
        selectedProductName = product.name;
        selectedProductDuration = product.duration;
        selectedProductPrice = productPrices[product.id]
            ? (productPrices[product.id] / 100).toFixed(2).replace(/\.00$/, "")
            : null;
        selectedPartnerId = null;
        selectedPartnerName = null;
        selectedAvailabilityId = null;
        selectedSlotStartTime = null;
        selectedSlotDisplay = null;
        availabilities = [];
        fetchPartners(product.id);
    }

    function partnerDisplayName(partner: Partner): string {
        return partner.bio?.split(/[.(]/)[0]?.trim() || "Praticien";
    }

    function selectPartner(partner: Partner) {
        selectedPartnerId = partner.id;
        selectedPartnerName = partnerDisplayName(partner);
        selectedAvailabilityId = null;
        selectedSlotStartTime = null;
        selectedSlotDisplay = null;
        availabilities = [];
        fetchAvailabilities(partner.id);
    }

    function selectSlot(availability: Availability) {
        selectedAvailabilityId = availability.id;
        selectedSlotStartTime = availability.start_time;
        const start = new Date(availability.start_time);
        const end = new Date(availability.end_time);
        selectedSlotDisplay = `${start.toLocaleDateString("fr-FR", {
            weekday: "long",
            day: "numeric",
            month: "long",
        })} — ${start.toLocaleTimeString("fr-FR", { hour: "2-digit", minute: "2-digit" })} à ${end.toLocaleTimeString("fr-FR", { hour: "2-digit", minute: "2-digit" })}`;
    }

    function formatTime(iso: string): string {
        return new Date(iso).toLocaleTimeString("fr-FR", {
            hour: "2-digit",
            minute: "2-digit",
        });
    }

    // --- Submit validation — on form submit event so Enter-key is also covered ---
    function handleSubmit(event: SubmitEvent) {
        if (!data.user) {
            guestFormErrors = {};
            if (!guestFirstName.trim()) guestFormErrors.guest_first_name = "Le prénom est requis";
            if (!guestLastName.trim()) guestFormErrors.guest_last_name = "Le nom est requis";
            if (!guestEmail.trim() && !guestPhone.trim()) {
                guestFormErrors.guest_email = "Email ou téléphone requis";
                guestFormErrors.guest_phone = "Email ou téléphone requis";
            }
            if (Object.keys(guestFormErrors).length > 0) {
                event.preventDefault();
                return;
            }
        }
    }
</script>

<div class="bg-surface min-h-screen py-24 md:py-32 px-4 lg:px-8">
    <div class="max-w-4xl mx-auto">
        <!-- Header -->
        <div class="text-center mb-12" use:reveal={{ preset: "fade-up", delay: 100 }}>
            <h1 class="text-4xl md:text-5xl font-bold text-foreground mb-4">
                Réserver une Séance
            </h1>
            <p class="text-lg text-foreground-alt">
                Choisissez votre soin, votre praticien et votre créneau en quelques clics
            </p>
        </div>

        <!-- Error banner -->
        {#if form?.errors?._form}
            <div class="bg-red-50 border border-red-200 text-red-700 rounded-2xl p-4 mb-8">
                {form.errors._form}
            </div>
        {/if}

        <form method="POST" action="/book" class="grid gap-8" onsubmit={handleSubmit}>
            <input type="hidden" name="availability_id" value={selectedAvailabilityId ?? ""} />
            <input type="hidden" name="product_id" value={selectedProductId ?? ""} />
            <input type="hidden" name="slot_start_time" value={selectedSlotStartTime ?? ""} />

            <!-- Step 1 — Catégorie -->
            <section class="bg-white rounded-3xl p-6 md:p-8" use:reveal={{ preset: "fade-up", delay: 150 }}>
                <div class="flex items-center gap-3 mb-6">
                    <span class="w-8 h-8 rounded-full bg-foreground text-white flex items-center justify-center text-sm font-bold">1</span>
                    <h2 class="text-xl md:text-2xl font-bold text-foreground">Catégorie</h2>
                </div>
                <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
                    {#each categories as category}
                        <button
                            type="button"
                            onclick={() => selectCategory(category)}
                            class="text-left p-5 rounded-2xl border-2 transition-all cursor-pointer {selectedCategoryId === category.id
                                ? 'border-foreground bg-surface'
                                : 'border-border-input hover:border-border-input-hover bg-white'}"
                        >
                            <h3 class="font-semibold text-foreground mb-1">{category.name}</h3>
                            <p class="text-sm text-foreground-alt line-clamp-2">{category.description}</p>
                        </button>
                    {/each}
                </div>
            </section>

            <!-- Step 2 — Soin -->
            {#if selectedCategoryId}
                <section class="bg-white rounded-3xl p-6 md:p-8" use:reveal={{ preset: "fade-up", delay: 100 }}>
                    <div class="flex items-center gap-3 mb-6">
                        <span class="w-8 h-8 rounded-full bg-foreground text-white flex items-center justify-center text-sm font-bold">2</span>
                        <h2 class="text-xl md:text-2xl font-bold text-foreground">Soin</h2>
                        {#if selectedCategoryName}
                            <span class="text-sm text-muted-foreground ml-2">— {selectedCategoryName}</span>
                        {/if}
                    </div>
                    {#if loadingProducts}
                        <p class="text-muted-foreground">Chargement des soins...</p>
                    {:else if products.length === 0}
                        <p class="text-muted-foreground">Aucun soin disponible dans cette catégorie</p>
                    {:else}
                        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
                            {#each products as product}
                                <button
                                    type="button"
                                    onclick={() => selectProduct(product)}
                                    class="text-left p-5 rounded-2xl border-2 transition-all cursor-pointer {selectedProductId === product.id
                                        ? 'border-foreground bg-surface'
                                        : 'border-border-input hover:border-border-input-hover bg-white'}"
                                >
                                    <div class="flex justify-between items-start mb-2">
                                        <h3 class="font-semibold text-foreground">{product.name}</h3>
                                        {#if productPrices[product.id]}
                                            <span class="font-bold text-foreground whitespace-nowrap ml-2">
                                                {(productPrices[product.id] / 100).toFixed(2).replace(/\.00$/, "")}€
                                            </span>
                                        {/if}
                                    </div>
                                    <p class="text-sm text-foreground-alt line-clamp-2 mb-2">{product.description ?? ""}</p>
                                    <span class="text-xs text-muted-foreground">{product.duration} min.</span>
                                </button>
                            {/each}
                        </div>
                    {/if}
                </section>
            {/if}

            <!-- Step 3 — Praticien -->
            {#if selectedProductId}
                <section class="bg-white rounded-3xl p-6 md:p-8" use:reveal={{ preset: "fade-up", delay: 100 }}>
                    <div class="flex items-center gap-3 mb-6">
                        <span class="w-8 h-8 rounded-full bg-foreground text-white flex items-center justify-center text-sm font-bold">3</span>
                        <h2 class="text-xl md:text-2xl font-bold text-foreground">Praticien</h2>
                    </div>
                    {#if loadingPartners}
                        <p class="text-muted-foreground">Chargement des praticiens...</p>
                    {:else if partners.length === 0}
                        <p class="text-muted-foreground">Aucun praticien disponible pour ce soin</p>
                    {:else}
                        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
                            {#each partners as partner}
                                <button
                                    type="button"
                                    onclick={() => selectPartner(partner)}
                                    class="text-left p-5 rounded-2xl border-2 transition-all cursor-pointer {selectedPartnerId === partner.id
                                        ? 'border-foreground bg-surface'
                                        : 'border-border-input hover:border-border-input-hover bg-white'}"
                                >
                                    <div class="flex items-center gap-4 mb-3">
                                        <div class="w-12 h-12 rounded-full bg-border-input-hover flex-shrink-0"></div>
                                        <div>
                                            <h3 class="font-semibold text-foreground">{partnerDisplayName(partner)}</h3>
                                        </div>
                                    </div>
                                    {#if partner.experience}
                                        <p class="text-sm text-foreground-alt line-clamp-3">{partner.experience}</p>
                                    {/if}
                                </button>
                            {/each}
                        </div>
                    {/if}
                </section>
            {/if}

            <!-- Step 4 — Créneau -->
            {#if selectedPartnerId}
                <section class="bg-white rounded-3xl p-6 md:p-8" use:reveal={{ preset: "fade-up", delay: 100 }}>
                    <div class="flex items-center gap-3 mb-6">
                        <span class="w-8 h-8 rounded-full bg-foreground text-white flex items-center justify-center text-sm font-bold">4</span>
                        <h2 class="text-xl md:text-2xl font-bold text-foreground">Créneau</h2>
                    </div>
                    {#if loadingSlots}
                        <p class="text-muted-foreground">Chargement des disponibilités...</p>
                    {:else if availabilities.length === 0}
                        <p class="text-muted-foreground">Aucun créneau disponible pour ce praticien</p>
                    {:else}
                        <div class="grid gap-6">
                            {#each Object.entries(availabilitiesByDay) as [day, slots]}
                                <div>
                                    <h4 class="font-semibold text-foreground capitalize mb-3">{day}</h4>
                                    <div class="flex flex-wrap gap-2">
                                        {#each slots as slot}
                                            <button
                                                type="button"
                                                onclick={() => selectSlot(slot)}
                                                class="px-4 py-2 rounded-xl border-2 transition-all cursor-pointer text-sm {selectedAvailabilityId === slot.id
                                                    ? 'border-foreground bg-foreground text-white'
                                                    : 'border-border-input-hover hover:border-border-input-hover bg-white text-foreground-alt'}"
                                            >
                                                {formatTime(slot.start_time)}
                                            </button>
                                        {/each}
                                    </div>
                                </div>
                            {/each}
                        </div>
                    {/if}
                </section>
            {/if}

            <!-- Step 5 — Vos coordonnées (guest only) -->
            {#if selectedAvailabilityId && !data.user}
                <section class="bg-white rounded-3xl p-6 md:p-8" use:reveal={{ preset: "fade-up", delay: 100 }}>
                    <div class="flex items-center gap-3 mb-6">
                        <span class="w-8 h-8 rounded-full bg-foreground text-white flex items-center justify-center text-sm font-bold">5</span>
                        <h2 class="text-xl md:text-2xl font-bold text-foreground">Vos coordonnées</h2>
                    </div>
                    <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
                        <div>
                            <label for="guest_first_name" class="block text-sm font-medium text-foreground-alt mb-1">Prénom *</label>
                            <input
                                id="guest_first_name"
                                name="guest_first_name"
                                type="text"
                                bind:value={guestFirstName}
                                required
                                class="w-full px-4 py-3 border border-border-input-hover rounded-xl focus:ring-2 focus:ring-foreground focus:ring-offset-2 focus:outline-none"
                                placeholder="Jean"
                            />
                            {#if guestFormErrors.guest_first_name}
                                <p class="text-red-500 text-sm mt-1">{guestFormErrors.guest_first_name}</p>
                            {/if}
                        </div>
                        <div>
                            <label for="guest_last_name" class="block text-sm font-medium text-foreground-alt mb-1">Nom *</label>
                            <input
                                id="guest_last_name"
                                name="guest_last_name"
                                type="text"
                                bind:value={guestLastName}
                                required
                                class="w-full px-4 py-3 border border-border-input-hover rounded-xl focus:ring-2 focus:ring-foreground focus:ring-offset-2 focus:outline-none"
                                placeholder="Dupont"
                            />
                            {#if guestFormErrors.guest_last_name}
                                <p class="text-red-500 text-sm mt-1">{guestFormErrors.guest_last_name}</p>
                            {/if}
                        </div>
                        <div>
                            <label for="guest_email" class="block text-sm font-medium text-foreground-alt mb-1">Email</label>
                            <input
                                id="guest_email"
                                name="guest_email"
                                type="email"
                                bind:value={guestEmail}
                                class="w-full px-4 py-3 border border-border-input-hover rounded-xl focus:ring-2 focus:ring-foreground focus:ring-offset-2 focus:outline-none"
                                placeholder="jean@exemple.com"
                            />
                            {#if guestFormErrors.guest_email}
                                <p class="text-red-500 text-sm mt-1">{guestFormErrors.guest_email}</p>
                            {/if}
                        </div>
                        <div>
                            <label for="guest_phone" class="block text-sm font-medium text-foreground-alt mb-1">Téléphone</label>
                            <input
                                id="guest_phone"
                                name="guest_phone"
                                type="tel"
                                bind:value={guestPhone}
                                class="w-full px-4 py-3 border border-border-input-hover rounded-xl focus:ring-2 focus:ring-foreground focus:ring-offset-2 focus:outline-none"
                                placeholder="+33 6 12 34 56 78"
                            />
                            {#if guestFormErrors.guest_phone}
                                <p class="text-red-500 text-sm mt-1">{guestFormErrors.guest_phone}</p>
                            {/if}
                        </div>
                    </div>
                    <p class="text-sm text-muted-foreground mt-4">* Au moins un moyen de contact (email ou téléphone) est requis</p>
                </section>
            {/if}

            <!-- Step 5b — Authenticated user: show name as confirmation -->
            {#if selectedAvailabilityId && data.user}
                <section class="bg-white rounded-3xl p-6 md:p-8" use:reveal={{ preset: "fade-up", delay: 100 }}>
                    <div class="flex items-center gap-3 mb-4">
                        <span class="w-8 h-8 rounded-full bg-foreground text-white flex items-center justify-center text-sm font-bold">5</span>
                        <h2 class="text-xl md:text-2xl font-bold text-foreground">Réservation</h2>
                    </div>
                    <p class="text-foreground-alt">
                        Réservation au nom de <span class="font-semibold">{data.user.firstname} {data.user.lastname}</span>
                    </p>
                </section>
            {/if}

            <!-- Step 6 — Récapitulatif et confirmation -->
            {#if step === 6}
                <section class="bg-white rounded-3xl p-6 md:p-8" use:reveal={{ preset: "fade-up", delay: 100 }}>
                    <div class="flex items-center gap-3 mb-6">
                        <span class="w-8 h-8 rounded-full bg-foreground text-white flex items-center justify-center text-sm font-bold">6</span>
                        <h2 class="text-xl md:text-2xl font-bold text-foreground">Récapitulatif</h2>
                    </div>
                    <div class="grid gap-4 mb-8">
                        {#if selectedProductName}
                            <div class="flex justify-between py-3 border-b border-border-input">
                                <span class="text-foreground-alt">Soin</span>
                                <span class="font-semibold text-foreground">{selectedProductName}</span>
                            </div>
                        {/if}
                        {#if selectedProductDuration}
                            <div class="flex justify-between py-3 border-b border-border-input">
                                <span class="text-foreground-alt">Durée</span>
                                <span class="font-semibold text-foreground">{selectedProductDuration} min.</span>
                            </div>
                        {/if}
                        {#if selectedProductPrice}
                            <div class="flex justify-between py-3 border-b border-border-input">
                                <span class="text-foreground-alt">Prix</span>
                                <span class="font-semibold text-foreground">{selectedProductPrice}€</span>
                            </div>
                        {/if}
                        {#if selectedSlotDisplay}
                            <div class="flex justify-between py-3 border-b border-border-input">
                                <span class="text-foreground-alt">Créneau</span>
                                <span class="font-semibold text-foreground capitalize">{selectedSlotDisplay}</span>
                            </div>
                        {/if}
                        {#if !data.user && guestFirstName}
                            <div class="flex justify-between py-3 border-b border-border-input">
                                <span class="text-foreground-alt">Nom</span>
                                <span class="font-semibold text-foreground">{guestFirstName} {guestLastName}</span>
                            </div>
                        {/if}
                    </div>
                    <Button
                        type="submit"
                        class="w-full text-white px-8 py-5 rounded-2xl text-lg font-semibold cursor-pointer"
                    >
                        Confirmer la réservation
                    </Button>
                </section>
            {/if}
        </form>
    </div>
</div>
