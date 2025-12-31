<script lang="ts">
    // NOTE: A voir si afficher les prix ici est une bonne idea niveau UX.
    import Button from "$lib/ui/Button.svelte";
    import Card from "./_card.svelte";
    import { type Product } from "./types";

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
            name: "Bodywork & Massage",
            description:
                "Restore balance and alleviate tension with our therapeutic bodywork sessions. Tailored to your specific recovery needs.",
            products: [
                {
                    title: "Deep Tissue Therapy",
                    description:
                        "Focus on realigning deep layers of muscle and connectinve tissue. It is especially helpful for chronic aches and pain",
                    duration: 60,
                    price: "120",
                    tags: [],
                },
                {
                    title: "Lymphatic drainage",
                    description:
                        "A gentle massage that encourage the movement of lymph fluids aroung the body. Helps remove waste and toxins.",
                    duration: 60,
                    price: "110",
                    tags: [],
                },
            ],
        },
        {
            name: "Mindset Coaching",
            description:
                "Restore balance and alleviate tension with our therapeutic bodywork sessions. Tailored to your specific recovery needs.",
            products: [
                {
                    title: "Executive Performance",
                    description:
                        "One-on-one coaching for high-performers looking to optimize decision making and leadership presence.",
                    duration: 60,
                    price: "200",
                    tags: [],
                },
            ],
        },
        {
            name: "Physical Training",
            description:
                "Restore balance and alleviate tension with our therapeutic bodywork sessions. Tailored to your specific recovery needs.",
            products: [
                {
                    title: "Strength Foundations",
                    description:
                        "Improve range of motion and joint health through functionnal movement pattern.",
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
    >
        <p class="uppercase">Our philosophy</p>
        <div>
            <!-- <h2 class="text-3xl font-extrabold text-dark-900"> -->
            <h2
                class="text-4xl sm:text-5xl lg:text-6xl font-medium tracking-tight text-dark-900 leading-[1.1]"
            >
                Holistic Wellness For
            </h2>
            <h2
                class="text-4xl sm:text-5xl lg:text-6xl font-medium tracking-tight text-dark-700 leading-[1.1] mb-6"
            >
                Body & Mind
            </h2>
        </div>
        <p
            class="text-lg text-dark-500 leading-relaxed mb-8 max-w-2xl font-normal"
        >
            Explore our comprehensive range of services designed to elevate your
            physical and mental well being through expert care and personalized
            coaching
        </p>
    </div>
    <div
        class="max-w-7xl mx-4 md:mx-auto grid grid-cols-1 lg:grid-cols-[auto_1fr] gap-8 lg:gap-20"
    >
        <div class="w-full lg:w-auto">
            <p class="uppercase text-dark-600 font-medium mb-4 px-2 lg:px-0">
                categories
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
        <div class="w-full px-2 lg:px-0">
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
                        {#each category.products as product}
                            <Card {...product} />
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
                    {#each products as product}
                        <Card {...product} />
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
