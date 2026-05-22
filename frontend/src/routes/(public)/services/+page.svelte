<script lang="ts">
    // NOTE: A voir si afficher les prix ici est une bonne idea niveau UX.
    import Button from "$lib/ui/Button.svelte";
    import Card from "./_card.svelte";
    import { type Product } from "./types";
    import { reveal } from "$lib/actions/reveal";

    // import type { PageProps } from "./$types";
    // let { data }: PageProps = $props();
    // let { services } = data

    interface Category {
        name: string;
        description: string;
        products: Product[];
    }

    let categories: Category[] = [
        {
            name: "Massothérapie & Corps",
            description:
                "Rétablissez l'équilibre et soulagez les tensions grâce à nos séances de thérapie corporelle. Adapté à vos besoins spécifiques de récupération.",
            products: [
                {
                    title: "Thérapie des Tissus Profonds",
                    description:
                        "Concentration sur le réalignement des couches profondes des muscles et des tissus conjonctifs. Particulièrement utile pour les douleurs chroniques.",
                    duration: 60,
                    price: "120",
                    tags: [],
                },
                {
                    title: "Drainage Lymphatique",
                    description:
                        "Un massage doux qui encourage le mouvement des fluides lymphatiques dans tout le corps. Aide à éliminer les déchets et les toxines.",
                    duration: 60,
                    price: "110",
                    tags: [],
                },
            ],
        },
        {
            name: "Coaching Mental",
            description:
                "Développez votre état d'esprit et renforcez votre résilience grâce à nos séances de coaching personnalisé. Adapté à vos objectifs de croissance personnelle.",
            products: [
                {
                    title: "Performance Exécutive",
                    description:
                        "Coaching individuel pour les performeurs cherchant à optimiser leur prise de décision et leur présence leadership.",
                    duration: 60,
                    price: "200",
                    tags: [],
                },
            ],
        },
        {
            name: "Entraînement Physique",
            description:
                "Améliorez votre condition physique grâce à des programmes d'entraînement personnalisés. Adapté à votre niveau et vos objectifs de santé.",
            products: [
                {
                    title: "Fondations en Force",
                    description:
                        "Améliorez l'amplitude des mouvements et la santé articulaire à travers des modèles de mouvements fonctionnels.",
                    duration: 90,
                    price: "80",
                    tags: [],
                },
            ],
        },
    ];

    const allCategories = "Toutes les categories";
    let activeCategory = $state(allCategories);
    let products = $derived(
        categories.find((c) => c.name === activeCategory)?.products ?? [],
    );
    let activeDescription = $derived(
        categories.find((c) => c.name === activeCategory)?.description ?? "",
    );
</script>

<div class="bg-dark-50 py-12 md:py-24 px-4 lg:px-8">
    <!-- Hero Section -->
    <div
        class="bg-white rounded-3xl px-6 py-8 md:px-12 lg:px-16 md:py-12 my-8 md:my-16 max-w-7xl mx-4 md:mx-auto"
        use:reveal={{ preset: "fade-up", delay: 100 }}
    >
        <p class="uppercase">Notre philosophie</p>
        <div>
            <!-- <h2 class="text-3xl font-extrabold text-dark-900"> -->
            <h2
                class="text-4xl sm:text-5xl lg:text-6xl font-medium tracking-tight text-dark-900 leading-[1.1]"
            >
                Bien-être Holistique Pour
            </h2>
            <h2
                class="text-4xl sm:text-5xl lg:text-6xl font-medium tracking-tight text-dark-700 leading-[1.1] mb-6"
            >
                Corps & Esprit
            </h2>
        </div>
        <p
            class="text-lg text-dark-500 leading-relaxed mb-8 max-w-2xl font-normal"
        >
            Découvrez notre gamme complète de services conçus pour améliorer votre
            bien-être physique et mental grâce à des soins d'experts et un
            coaching personnalisé
        </p>
    </div>
    <div
        class="max-w-7xl mx-4 md:mx-auto grid grid-cols-1 lg:grid-cols-[auto_1fr] gap-8 lg:gap-20"
    >
        <div class="w-full lg:w-auto" use:reveal={{ preset: "fade-up", delay: 150 }}>
            <p class="uppercase text-dark-600 font-medium mb-4 px-2 lg:px-0">
                catégories
            </p>
            <div
                class="w-full lg:max-w-xs flex flex-row lg:flex-col gap-2 overflow-x-auto pb-2 px-2 xl:px-0"
            >
                {@render categoryFilter(allCategories)}
                {#each categories as category}
                    {@render categoryFilter(category.name)}
                {/each}
            </div>
        </div>
        <div class="w-full px-2 lg:px-0" use:reveal={{ preset: "fade-up", delay: 150 }}>
            <div class="grid gap-8 md:gap-12">
                {#if activeCategory === allCategories}
                    {#each categories as category}
                        <div>
                            <p
                                class="text-2xl md:text-3xl font-bold text-dark-900"
                            >
                                {category.name}
                            </p>
                            <p class="text-sm md:text-base text-dark-700">
                                {category.description}
                            </p>
                        </div>
                        {#each category.products as product, index}
                            <div use:reveal={{ preset: "fade-up", delay: 200 + index * 60 }}>
                                <Card {...product} />
                            </div>
                        {/each}
                    {/each}
                {:else}
                    <div>
                        <p class="text-2xl md:text-3xl font-bold text-dark-900">
                            {activeCategory}
                        </p>
                        <p class="text-sm md:text-base text-dark-600">
                            {activeDescription}
                        </p>
                    </div>
                    {#each products as product, index}
                        <div use:reveal={{ preset: "fade-up", delay: 200 + index * 60 }}>
                            <Card {...product} />
                        </div>
                    {/each}
                {/if}
            </div>
        </div>
    </div>
</div>

{#snippet categoryFilter(name: string)}
    <Button
        onclick={() => (activeCategory = name)}
        class="{activeCategory === name
            ? `bg-white font-bold text-dark-900 hover:bg-white`
            : `bg-transparent text-dark-700 font-medium hover:bg-white/60`} justify-start px-6 py-4 lg:px-12 lg:py-6 cursor-pointer rounded-3xl whitespace-nowrap"
    >
        {name}
    </Button>
{/snippet}
