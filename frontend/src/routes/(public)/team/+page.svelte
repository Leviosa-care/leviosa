<script lang="ts">
    // NOTE: A voir si afficher les prix ici est une bonne idea niveau UX.
    import Button from "$lib/ui/Button.svelte";
    import Card from "./_card.svelte";
    import { type Partner } from "./types";
    import { reveal } from "$lib/actions/reveal";

    // import type { PageProps } from "./$types";
    // let { data }: PageProps = $props();
    // let { categories } = data

    interface Category {
        name: string;
        description: string;
        partners: Partner[];
    }

    let categories: Category[] = [
        {
            name: "Bodywork & Massage",
            description:
                "Experts en thérapie manuelle et récupération physique.",
            partners: [
                {
                    firstname: "Sarah",
                    lastname: "Chen",
                    occupation: "Ostéopathe D.O",
                    quote: "Le corps garde en mémoire ce que l'esprit oublie. Mon approche vise à libérer ces tensions profondes.",
                    tags: ["Douleurs Chroniques", "Posturologie"],
                    picture: "https://images.unsplash.com/photo-1594824476967-48c8b964273f?w=200&h=200&auto=format&fit=crop&q=80",
                },
                {
                    firstname: "Marc",
                    lastname: "Dubois",
                    occupation: "Massothérapeute sportif",
                    quote: "Optimiser la récupération pour permettre une performance durable, sans blessure.",
                    tags: ["Récupération", "Deep Tissue"],
                    picture: "https://images.unsplash.com/photo-1570295999919-56ceb5ecca61?w=200&h=200&auto=format&fit=crop&q=80",
                },
                {
                    firstname: "Jean",
                    lastname: "Dupont",
                    occupation: "Massothérapeute sportif",
                    quote: "Optimiser la récupération pour permettre une performance durable, sans blessure.",
                    tags: ["Récupération", "Deep Tissue"],
                    picture: "https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?w=200&h=200&auto=format&fit=crop&q=80",
                },
            ],
        },
        {
            name: "Mindset Coaching",
            description:
                "Psychologues et coachs certifiés pour votre équilibre.",
            partners: [
                {
                    firstname: "Elena",
                    lastname: "Rodriguez",
                    occupation: "Psychologue Clinicienne",
                    quote: "Un espace sécurisant pour comprendre vos mécanismes et retrouver votre sérénité",
                    tags: ["Gestion du stress", "Anxiété"],
                    picture: "https://images.unsplash.com/photo-1438761681033-6461ffad8d80?w=200&h=200&auto=format&fit=crop&q=80",
                },
                {
                    firstname: "David",
                    lastname: "Miller",
                    occupation: "Coach Exécutif",
                    quote: "Clarifier votre vision pour agir avec impact et confiance dans les moments critiques",
                    tags: ["Leadership", "Prise de décision"],
                    picture: "https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=200&h=200&auto=format&fit=crop&q=80",
                },
            ],
        },
        {
            name: "Physical Training",
            description: "Programmation sur mesure et coaching technique",
            partners: [
                {
                    firstname: "Alexandre",
                    lastname: "T.",
                    occupation: "Coach Sportif",
                    quote: "Le mouvement est le meilleur médicament. Construisons un corps fort et capable.",
                    tags: ["Force", "Mobilité"],
                    picture: "https://images.unsplash.com/photo-1568602471122-7832951cc4c5?w=200&h=200&auto=format&fit=crop&q=80",
                },
                {
                    firstname: "Sophie",
                    lastname: "Laurent",
                    occupation: "Instructrice Yoga",
                    quote: "Connecter le souffle et le mouvement pour trouver l'équilibre physique et mental.",
                    tags: ["Vinyasa", "Méditation"],
                    picture: "https://images.unsplash.com/photo-1545389336-cf090694435e?w=200&h=200&auto=format&fit=crop&q=80",
                },
            ],
        },
    ];

    const allCategories = "Toutes les catégories";
    let activeCategory = $state(allCategories);
    let partners = $derived(
        categories.find((c) => c.name === activeCategory)?.partners ?? [],
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
        <div>
            <!-- <h2 class="text-3xl font-extrabold text-dark-900"> -->
            <h2
                class="text-4xl sm:text-5xl lg:text-6xl font-medium tracking-tight text-dark-900 leading-[1.1]"
            >
                Une équipe engagée
            </h2>
            <h2
                class="text-4xl sm:text-5xl lg:text-6xl font-medium tracking-tight text-dark-700 leading-[1.1] mb-6"
            >
                Pour votre bien-être
            </h2>
        </div>
        <p
            class="text-lg text-dark-500 leading-relaxed mb-8 max-w-2xl font-normal"
        >
            Nos praticiens et coachs partagent une vision holistique de la
            santé. Découvrez des profils aux expertises complémentaires, unis
            par la volonté de vous accompagner durablement.
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
                        <div
                            class="grid grid-cols-1 md:grid-cols-2 gap-6 md:gap-8"
                        >
                            {#each category.partners as partner, index}
                                <div use:reveal={{ preset: "fade-up", delay: 200 + index * 60 }}>
                                    <Card {...partner} />
                                </div>
                            {/each}
                        </div>
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

                    <div class="grid grid-cols-1 md:grid-cols-2 gap-6 md:gap-7">
                        {#each partners as partner, index}
                            <div use:reveal={{ preset: "fade-up", delay: 200 + index * 60 }}>
                                <Card {...partner} />
                            </div>
                        {/each}
                    </div>
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
