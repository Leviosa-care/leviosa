<script lang="ts">
    import Card from "./_card.svelte";
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
            <div class="max-w-3xl">
                <!-- Badge -->
                <div
                    class="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-dark-50 border border-dark-200 mb-8"
                    use:reveal={{ preset: "fade-down", delay: 100 }}
                >
                    <span class="text-xs font-semibold text-dark-500 uppercase tracking-wider"
                        >Notre philosophie</span
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
                    class="text-lg md:text-xl text-dark-500 leading-relaxed max-w-2xl font-normal mb-12"
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
        </div>
    </div>

    <!-- Content: Sticky Sidebar + Products -->
    <div class="relative max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-16 md:py-20">
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
