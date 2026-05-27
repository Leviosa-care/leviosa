<script lang="ts">
    import Card from "./_card.svelte";
    import { reveal } from "$lib/actions/reveal";
    import type { PageProps } from "./$types";

    let { data }: PageProps = $props();

    const allCategories = "Toutes les catégories";
    let activeCategory = $state(allCategories);

    let filteredPartners = $derived(
        data.categories.find((c: any) => c.name === activeCategory)?.partners ?? [],
    );
    let activeDescription = $derived(
        data.categories.find((c: any) => c.name === activeCategory)?.description ?? "",
    );

    let totalPartners = $derived(
        data.categories.reduce((acc: number, c: any) => acc + c.partners.length, 0),
    );

    let heroPartners = $derived(data.categories.flatMap((c: any) => c.partners).slice(0, 3));
</script>

<div
    class="min-h-screen bg-white"
    style="background-image: radial-gradient(rgba(15,23,42,0.035) 1px, transparent 1px); background-size: 24px 24px;"
>
    <!-- Hero -->
    <div class="relative">
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
                            >Notre équipe</span
                        >
                    </div>

                    <!-- Headline with weight contrast -->
                    <h1
                        class="text-5xl sm:text-6xl lg:text-7xl font-medium tracking-tight text-dark-900 leading-[1.05] mb-6"
                        use:reveal={{ preset: "fade-up", delay: 150 }}
                    >
                        Une équipe engagée
                        <br class="hidden sm:block" />
                        <span class="text-dark-400 font-light">Pour votre bien-être</span>
                    </h1>

                    <p
                        class="text-lg md:text-xl text-dark-500 leading-relaxed font-normal mb-12"
                        use:reveal={{ preset: "fade-up", delay: 200 }}
                    >
                        Nos praticiens et coachs partagent une vision holistique de la santé.
                        Découvrez des profils aux expertises complémentaires, unis par la volonté
                        de vous accompagner durablement.
                    </p>

                    <!-- Stats row -->
                    <div
                        class="flex flex-wrap items-center gap-6 md:gap-10"
                        use:reveal={{ preset: "fade-up", delay: 250 }}
                    >
                        <div>
                            <p class="text-2xl md:text-3xl font-semibold text-dark-900">{totalPartners}+</p>
                            <p class="text-sm text-dark-500 mt-0.5">praticiens experts</p>
                        </div>
                        <div class="w-px h-10 bg-dark-200 self-stretch"></div>
                        <div>
                            <p class="text-2xl md:text-3xl font-semibold text-dark-900">{data.categories.length}</p>
                            <p class="text-sm text-dark-500 mt-0.5">spécialités</p>
                        </div>
                        <div class="w-px h-10 bg-dark-200 self-stretch"></div>
                        <div>
                            <p class="text-2xl md:text-3xl font-semibold text-dark-900">5 ★</p>
                            <p class="text-sm text-dark-500 mt-0.5">note moyenne</p>
                        </div>
                    </div>
                </div>

                <!-- Right: partner portrait cards -->
                <div class="hidden lg:flex flex-col gap-3" use:reveal={{ preset: "fade-up", delay: 200 }}>
                    {#if heroPartners[0]}
                        <div class="bg-dark-900 text-white rounded-2xl p-6 shadow-md flex items-center gap-4">
                            {#if heroPartners[0].picture}
                                <img
                                    src={heroPartners[0].picture}
                                    alt="{heroPartners[0].firstname} {heroPartners[0].lastname}"
                                    class="w-14 h-14 rounded-xl object-cover flex-shrink-0"
                                />
                            {:else}
                                <div class="w-14 h-14 rounded-xl bg-dark-700 flex items-center justify-center text-sm font-semibold flex-shrink-0">
                                    {heroPartners[0].firstname[0]}{heroPartners[0].lastname[0]}
                                </div>
                            {/if}
                            <div>
                                <p class="text-white font-semibold leading-snug">
                                    {heroPartners[0].firstname} {heroPartners[0].lastname}
                                </p>
                                <p class="text-white/60 text-sm mt-0.5">{heroPartners[0].occupation}</p>
                            </div>
                        </div>
                    {/if}

                    <div class="grid grid-cols-2 gap-3">
                        {#if heroPartners[1]}
                            <div class="bg-white border border-dark-100 rounded-2xl p-5 shadow-sm hover:shadow-md transition-shadow duration-200 flex flex-col gap-3">
                                {#if heroPartners[1].picture}
                                    <img
                                        src={heroPartners[1].picture}
                                        alt="{heroPartners[1].firstname} {heroPartners[1].lastname}"
                                        class="w-10 h-10 rounded-lg object-cover"
                                    />
                                {:else}
                                    <div class="w-10 h-10 rounded-lg bg-dark-100 flex items-center justify-center text-xs font-semibold text-dark-600">
                                        {heroPartners[1].firstname[0]}{heroPartners[1].lastname[0]}
                                    </div>
                                {/if}
                                <div>
                                    <p class="text-dark-900 font-semibold text-sm leading-snug">
                                        {heroPartners[1].firstname} {heroPartners[1].lastname}
                                    </p>
                                    <p class="text-dark-400 text-xs mt-0.5">{heroPartners[1].occupation}</p>
                                </div>
                            </div>
                        {/if}
                        {#if heroPartners[2]}
                            <div class="bg-white border border-dark-100 rounded-2xl p-5 shadow-sm hover:shadow-md transition-shadow duration-200 flex flex-col gap-3">
                                {#if heroPartners[2].picture}
                                    <img
                                        src={heroPartners[2].picture}
                                        alt="{heroPartners[2].firstname} {heroPartners[2].lastname}"
                                        class="w-10 h-10 rounded-lg object-cover"
                                    />
                                {:else}
                                    <div class="w-10 h-10 rounded-lg bg-dark-100 flex items-center justify-center text-xs font-semibold text-dark-600">
                                        {heroPartners[2].firstname[0]}{heroPartners[2].lastname[0]}
                                    </div>
                                {/if}
                                <div>
                                    <p class="text-dark-900 font-semibold text-sm leading-snug">
                                        {heroPartners[2].firstname} {heroPartners[2].lastname}
                                    </p>
                                    <p class="text-dark-400 text-xs mt-0.5">{heroPartners[2].occupation}</p>
                                </div>
                            </div>
                        {/if}
                    </div>

                    {#if totalPartners > 3}
                        <p class="text-xs text-dark-400 text-right">
                            +{totalPartners - 3} autres praticiens disponibles
                        </p>
                    {/if}
                </div>
            </div>
        </div>
    </div>

    <!-- Content: Sticky Sidebar + Partner Cards -->
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
                    {#each data.categories as category}
                        {@render categoryFilter(category.name)}
                    {/each}
                </div>
            </div>

            <!-- Partners List -->
            <div class="w-full">
                <div class="grid gap-6 md:gap-8">
                    {#if activeCategory === allCategories}
                        {#each data.categories as category}
                            {#if category.partners.length > 0}
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
                                <div class="grid grid-cols-1 md:grid-cols-2 gap-6 md:gap-8">
                                    {#each category.partners as partner, index}
                                        <div use:reveal={{ preset: "fade-up", delay: 100 + index * 50 }}>
                                            <Card {...partner} />
                                        </div>
                                    {/each}
                                </div>
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
                        <div class="grid grid-cols-1 md:grid-cols-2 gap-6 md:gap-8">
                            {#each filteredPartners as partner, index}
                                <div use:reveal={{ preset: "fade-up", delay: 100 + index * 50 }}>
                                    <Card {...partner} />
                                </div>
                            {/each}
                        </div>
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
