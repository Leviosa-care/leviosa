<script lang="ts">
    import Card from "./_card.svelte";
    import { Clock } from "@lucide/svelte";
    import { reveal } from "$lib/actions/reveal";
    import type { PageProps } from "./$types";

    let { data }: PageProps = $props();

    const allCategories = "Toutes les categories";
    let activeCategory = $state(allCategories);

    let categories = $derived.by(() => {
        const cats = data.categories.map((cat: any) => ({
            id: cat.id,
            name: cat.name,
            description: cat.description,
            products: data.products
                .filter((p: any) => p.category?.id === cat.id)
                .map((p: any) => ({
                    id: p.id,
                    title: p.name,
                    description: p.description ?? "",
                    duration: p.duration,
                    price: data.pricesByProduct[p.id]
                        ? (data.pricesByProduct[p.id] / 100).toFixed(2).replace(/\.00$/, "")
                        : "",
                    tags: [],
                    image: p.images?.find((img: any) => img.is_active)?.url ?? p.images?.[0]?.url ?? "",
                })),
        }));
        return cats;
    });

    let products = $derived(
        categories.find((c: any) => c.name === activeCategory)?.products ?? [],
    );
    let activeDescription = $derived(
        categories.find((c: any) => c.name === activeCategory)?.description ?? "",
    );

    let heroProducts = $derived(
        data.products.slice(0, 3).map((p: any) => ({
            name: p.name,
            category: p.category?.name ?? "",
            duration: p.duration,
        })),
    );
</script>

<div
    class="min-h-screen bg-white"
    style="background-image: radial-gradient(rgba(15,23,42,0.035) 1px, transparent 1px); background-size: 24px 24px;"
>
    <!-- Hero -->
    <div class="relative">
        <!-- Soft gradient overlay -->
        <div
            class="absolute inset-0 bg-gradient-to-br from-dark-50/60 via-white/50 to-transparent pointer-events-none"
        ></div>

        <div class="relative max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-20 md:py-28">
            <div class="grid grid-cols-1 lg:grid-cols-[3fr_2fr] gap-12 lg:gap-16 items-center">
                <!-- Left: text content -->
                <div class="max-w-2xl">
                    <!-- Badge -->
                    <div
                        class="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-dark-50 border border-dark-200 mb-8"
                        use:reveal={{ preset: "fade-down", delay: 100 }}
                    >
                        <span class="text-xs font-semibold text-dark-500 uppercase tracking-wider"
                            >Nos services</span
                        >
                    </div>

                    <!-- Headline with weight contrast -->
                    <h1
                        class="text-5xl sm:text-6xl lg:text-7xl font-medium tracking-tight text-dark-900 leading-[1.05] mb-6"
                        use:reveal={{ preset: "fade-up", delay: 150 }}
                    >
                        Bien-être Holistique
                        <br class="hidden sm:block" />
                        <span class="text-dark-400 font-light">Pour Corps & Esprit</span>
                    </h1>

                    <p
                        class="text-lg md:text-xl text-dark-500 leading-relaxed font-normal mb-12"
                        use:reveal={{ preset: "fade-up", delay: 200 }}
                    >
                        Découvrez notre gamme complète de services conçus pour améliorer votre
                        bien-être physique et mental grâce à des soins d'experts et un
                        coaching personnalisé.
                    </p>

                    <!-- Stats row -->
                    <div
                        class="flex flex-wrap items-center gap-6 md:gap-10"
                        use:reveal={{ preset: "fade-up", delay: 250 }}
                    >
                        <div>
                            <p class="text-2xl md:text-3xl font-semibold text-dark-900">2 000+</p>
                            <p class="text-sm text-dark-500 mt-0.5">clients satisfaits</p>
                        </div>
                        <div class="w-px h-10 bg-dark-200 self-stretch"></div>
                        <div>
                            <p class="text-2xl md:text-3xl font-semibold text-dark-900">10+</p>
                            <p class="text-sm text-dark-500 mt-0.5">services disponibles</p>
                        </div>
                        <div class="w-px h-10 bg-dark-200 self-stretch"></div>
                        <div>
                            <p class="text-2xl md:text-3xl font-semibold text-dark-900">5 ★</p>
                            <p class="text-sm text-dark-500 mt-0.5">note moyenne</p>
                        </div>
                    </div>
                </div>

                <!-- Right: highlight cards grid -->
                <div class="hidden lg:flex flex-col gap-3" use:reveal={{ preset: "fade-up", delay: 200 }}>
                    {#if heroProducts[0]}
                        <div
                            class="bg-dark-900 text-white rounded-2xl p-6 shadow-md"
                        >
                            <div class="flex items-center justify-between mb-4">
                                <span class="text-xs font-semibold text-white/50 uppercase tracking-wider">
                                    {heroProducts[0].category}
                                </span>
                                <span class="inline-flex items-center gap-1.5 text-xs text-white/60">
                                    <Clock size={11} />
                                    {heroProducts[0].duration} min
                                </span>
                            </div>
                            <p class="text-white font-semibold text-lg leading-snug">
                                {heroProducts[0].name}
                            </p>
                        </div>
                    {/if}

                    <div class="grid grid-cols-2 gap-3">
                        {#if heroProducts[1]}
                            <div
                                class="bg-white border border-dark-100 rounded-2xl p-5 shadow-sm hover:shadow-md transition-shadow duration-200"
                            >
                                <div class="flex items-start justify-between gap-2 mb-3">
                                    <span class="text-xs font-semibold text-dark-400 uppercase tracking-wider leading-tight">
                                        {heroProducts[1].category}
                                    </span>
                                    <span class="inline-flex items-center gap-1 text-xs text-dark-500 flex-shrink-0">
                                        <Clock size={11} />
                                        {heroProducts[1].duration} min
                                    </span>
                                </div>
                                <p class="text-dark-900 font-semibold text-sm leading-snug">
                                    {heroProducts[1].name}
                                </p>
                            </div>
                        {/if}
                        {#if heroProducts[2]}
                            <div
                                class="bg-white border border-dark-100 rounded-2xl p-5 shadow-sm hover:shadow-md transition-shadow duration-200"
                            >
                                <div class="flex items-start justify-between gap-2 mb-3">
                                    <span class="text-xs font-semibold text-dark-400 uppercase tracking-wider leading-tight">
                                        {heroProducts[2].category}
                                    </span>
                                    <span class="inline-flex items-center gap-1 text-xs text-dark-500 flex-shrink-0">
                                        <Clock size={11} />
                                        {heroProducts[2].duration} min
                                    </span>
                                </div>
                                <p class="text-dark-900 font-semibold text-sm leading-snug">
                                    {heroProducts[2].name}
                                </p>
                            </div>
                        {/if}
                    </div>

                    {#if data.products.length > 3}
                        <p class="text-xs text-dark-400 text-right">
                            +{data.products.length - 3} autres services disponibles
                        </p>
                    {/if}
                </div>
            </div>
        </div>
    </div>

    <!-- Content: Sticky Sidebar + Products -->
    <div class="relative max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-10 md:py-12">
        <div class="grid grid-cols-1 lg:grid-cols-[260px_1fr] gap-10 lg:gap-16 items-start">
            <!-- Sticky Category Sidebar -->
            <div class="lg:sticky lg:top-24" use:reveal={{ preset: "fade-up", delay: 150 }}>
                <p class="uppercase text-dark-500 text-xs font-semibold tracking-wider mb-4">
                    catégories
                </p>
                <div
                    class="flex flex-row lg:flex-col gap-2 overflow-x-auto -mx-4 px-4 lg:mx-0 lg:px-0 pb-2 lg:pb-0"
                >
                    {@render categoryFilter(allCategories)}
                    {#each categories as category}
                        {@render categoryFilter(category.name)}
                    {/each}
                </div>
            </div>

            <!-- Products List -->
            <div class="w-full">
                <div class="grid gap-6 md:gap-8">
                    {#if activeCategory === allCategories}
                        {#each categories as category}
                            {#if category.products.length > 0}
                                <div
                                    class="mb-2"
                                    use:reveal={{ preset: "fade-up", delay: 100 }}
                                >
                                    <h2
                                        class="text-2xl md:text-3xl font-semibold text-dark-900 tracking-tight"
                                    >
                                        {category.name}
                                    </h2>
                                    <p class="text-sm md:text-base text-dark-500 mt-1">
                                        {category.description}
                                    </p>
                                </div>
                                {#each category.products as product, index}
                                    <div use:reveal={{ preset: "fade-up", delay: 100 + index * 50 }}>
                                        <Card {...product} />
                                    </div>
                                {/each}
                            {/if}
                        {/each}
                    {:else}
                        <div class="mb-2" use:reveal={{ preset: "fade-up", delay: 100 }}>
                            <h2
                                class="text-2xl md:text-3xl font-semibold text-dark-900 tracking-tight"
                            >
                                {activeCategory}
                            </h2>
                            <p class="text-sm md:text-base text-dark-500 mt-1">
                                {activeDescription}
                            </p>
                        </div>
                        {#each products as product, index}
                            <div use:reveal={{ preset: "fade-up", delay: 100 + index * 50 }}>
                                <Card {...product} />
                            </div>
                        {/each}
                    {/if}
                </div>
            </div>
        </div>
    </div>
</div>

{#snippet categoryFilter(name: string)}
    <button
        onclick={() => (activeCategory = name)}
        class="{activeCategory === name
            ? 'bg-dark-900 text-white border-dark-900 shadow-sm font-medium'
            : 'bg-white text-dark-600 border-dark-200 hover:border-dark-300 hover:text-dark-900 font-normal'} border text-sm px-5 py-3 rounded-xl whitespace-nowrap transition-all duration-200 cursor-pointer lg:w-full lg:text-left"
    >
        {name}
    </button>
{/snippet}
