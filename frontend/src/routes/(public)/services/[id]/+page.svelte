<script lang="ts">
    import { Clock, MapPin, Monitor, Shuffle, ChevronLeft, CalendarCheck } from "@lucide/svelte";
    import { reveal } from "$lib/actions/reveal";
    import type { PageProps } from "./$types";

    let { data }: PageProps = $props();

    const { product, price } = data;

    const formattedPrice = price
        ? (price / 100).toFixed(2).replace(/\.00$/, "")
        : null;

    const activeImage =
        product.images?.find((img: any) => img.is_active)?.url ??
        product.images?.[0]?.url ??
        null;

    const availabilityLabel: Record<string, string> = {
        "in-person": "En présentiel",
        online: "En ligne",
        hybrid: "Hybride",
    };

    const availabilityIcon: Record<string, typeof MapPin> = {
        "in-person": MapPin,
        online: Monitor,
        hybrid: Shuffle,
    };

    const AvailIcon = availabilityIcon[product.availability] ?? MapPin;
</script>

<div
    class="min-h-screen bg-white"
    style="background-image: radial-gradient(rgba(15,23,42,0.035) 1px, transparent 1px); background-size: 24px 24px;"
>
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-10 md:py-14">
        <!-- Back link -->
        <a
            href="/services"
            class="inline-flex items-center gap-1.5 text-sm text-dark-500 hover:text-dark-900 transition-colors duration-150 mb-10 group"
            use:reveal={{ preset: "fade-down", delay: 50 }}
        >
            <ChevronLeft size={16} class="group-hover:-translate-x-0.5 transition-transform duration-150" />
            Tous les services
        </a>

        <div class="grid grid-cols-1 lg:grid-cols-[1fr_320px] gap-10 lg:gap-16 items-start">
            <!-- Main content -->
            <div>
                <!-- Image -->
                {#if activeImage}
                    <div
                        class="w-full aspect-[16/7] rounded-3xl overflow-hidden mb-8"
                        use:reveal={{ preset: "fade-up", delay: 100 }}
                    >
                        <img
                            src={activeImage}
                            alt={product.name}
                            class="w-full h-full object-cover"
                        />
                    </div>
                {:else}
                    <div
                        class="w-full aspect-[16/7] rounded-3xl bg-gradient-to-br from-dark-100 to-dark-50 flex items-center justify-center mb-8"
                        use:reveal={{ preset: "fade-up", delay: 100 }}
                    >
                        <span class="iconify text-dark-300" data-icon="lucide:image" data-width="56"></span>
                    </div>
                {/if}

                <!-- Category + title -->
                <div use:reveal={{ preset: "fade-up", delay: 150 }}>
                    {#if product.category}
                        <p class="text-xs font-semibold text-dark-400 uppercase tracking-wider mb-3">
                            {product.category.name}
                        </p>
                    {/if}
                    <h1 class="text-3xl sm:text-4xl md:text-5xl font-semibold tracking-tight text-dark-900 mb-6">
                        {product.name}
                    </h1>
                </div>

                <!-- Meta pills -->
                <div
                    class="flex flex-wrap gap-2 mb-8"
                    use:reveal={{ preset: "fade-up", delay: 200 }}
                >
                    <div class="inline-flex items-center gap-1.5 bg-dark-50 border border-dark-100 rounded-full px-4 py-2">
                        <Clock size={13} class="text-dark-400" />
                        <span class="text-xs font-medium text-dark-600">{product.duration} min</span>
                    </div>
                    <div class="inline-flex items-center gap-1.5 bg-dark-50 border border-dark-100 rounded-full px-4 py-2">
                        <AvailIcon size={13} class="text-dark-400" />
                        <span class="text-xs font-medium text-dark-600">
                            {availabilityLabel[product.availability] ?? product.availability}
                        </span>
                    </div>
                </div>

                <!-- Description -->
                <div use:reveal={{ preset: "fade-up", delay: 250 }}>
                    <h2 class="text-lg font-semibold text-dark-900 mb-3">À propos de ce service</h2>
                    <p class="text-dark-500 leading-relaxed text-base md:text-lg">
                        {product.description}
                    </p>
                </div>

                <!-- Category description -->
                {#if product.category?.description}
                    <div class="mt-10 pt-10 border-t border-dark-100" use:reveal={{ preset: "fade-up", delay: 300 }}>
                        <h2 class="text-lg font-semibold text-dark-900 mb-3">{product.category.name}</h2>
                        <p class="text-dark-500 leading-relaxed">
                            {product.category.description}
                        </p>
                    </div>
                {/if}
            </div>

            <!-- Sticky sidebar: booking card -->
            <div class="lg:sticky lg:top-24" use:reveal={{ preset: "fade-up", delay: 200 }}>
                <div class="bg-white border border-dark-100 rounded-3xl p-6 shadow-sm">
                    <!-- Price -->
                    {#if formattedPrice}
                        <div class="flex items-baseline gap-1 mb-6">
                            <span class="text-4xl font-bold text-dark-900">{formattedPrice}</span>
                            <span class="text-lg text-dark-400">€</span>
                        </div>
                    {:else}
                        <p class="text-sm text-dark-400 mb-6">Prix sur demande</p>
                    {/if}

                    <!-- Réserver CTA -->
                    <a href="/book?product={product.id}" class="block w-full">
                        <button
                            class="group/btn w-full inline-flex justify-center items-center gap-2 bg-dark-900 hover:bg-dark-800 text-white text-sm font-medium px-6 py-3.5 rounded-xl transition-all duration-200 shadow-sm hover:shadow-md cursor-pointer"
                        >
                            <CalendarCheck size={16} />
                            Réserver ce service
                            <span
                                class="iconify group-hover/btn:translate-x-0.5 transition-transform"
                                data-icon="lucide:arrow-right"
                                data-width="16"
                                data-stroke-width="2"
                            ></span>
                        </button>
                    </a>

                    <!-- Details list -->
                    <ul class="mt-6 space-y-3 text-sm text-dark-500">
                        <li class="flex items-center gap-2">
                            <Clock size={14} class="text-dark-400 flex-shrink-0" />
                            Durée : {product.duration} min
                        </li>
                        <li class="flex items-center gap-2">
                            <AvailIcon size={14} class="text-dark-400 flex-shrink-0" />
                            {availabilityLabel[product.availability] ?? product.availability}
                        </li>
                        {#if product.cancellationHours}
                            <li class="flex items-center gap-2">
                                <span class="iconify text-dark-400 flex-shrink-0" data-icon="lucide:shield-check" data-width="14"></span>
                                Annulation gratuite sous {product.cancellationHours}h
                            </li>
                        {/if}
                    </ul>
                </div>
            </div>
        </div>
    </div>
</div>
