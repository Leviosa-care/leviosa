<script lang="ts">
    import { reveal } from "$lib/actions/reveal";

    const experts = [
        {
            id: "p-001",
            name: "Jane Smith",
            occupation: "Ostéopathe D.O",
            tags: ["Douleurs Chroniques", "Posturologie"],
            quote: "Le corps garde en mémoire ce que l'esprit oublie. Mon approche vise à libérer ces tensions profondes.",
            picture:
                "https://images.unsplash.com/photo-1594824476967-48c8b964273f?w=200&h=200&auto=format&fit=crop&q=80",
        },
        {
            id: "p-002",
            name: "Emma Davis",
            occupation: "Psychologue Clinicienne",
            tags: ["Gestion du stress", "Anxiété"],
            quote: "Un espace sécurisant pour comprendre vos mécanismes et retrouver votre sérénité.",
            picture:
                "https://images.unsplash.com/photo-1534528741775-53994a69daeb?w=200&h=200&auto=format&fit=crop&q=80",
        },
        {
            id: "p-003",
            name: "Alice Martin",
            occupation: "Coach Exécutif",
            tags: ["Leadership", "Prise de décision"],
            quote: "Clarifier votre vision pour agir avec impact et confiance dans les moments critiques.",
            picture:
                "https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=200&h=200&auto=format&fit=crop&q=80",
        },
    ];
</script>

<section class="py-20 lg:py-24 bg-dark-50">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <!-- Section Header -->
        <div
            class="flex flex-col md:flex-row md:items-end justify-between gap-6 mb-12"
            use:reveal={{ preset: "fade-up", delay: 100 }}
        >
            <div class="max-w-xl">
                <span
                    class="text-sm font-semibold text-dark-500 uppercase tracking-wider"
                    >Notre Équipe</span
                >
                <h2
                    class="mt-3 text-3xl md:text-4xl font-semibold tracking-tight text-dark-900"
                >
                    Des experts à votre service
                </h2>
                <p class="mt-4 text-base text-dark-500">
                    Chaque praticien est soigneusement sélectionné et vérifié
                    pour garantir la qualité de votre suivi.
                </p>
            </div>
            <a
                href="/team"
                class="group inline-flex items-center gap-2 text-sm font-medium text-dark-500 hover:text-dark-900 transition-colors px-4 py-2 rounded-full hover:bg-dark-100 whitespace-nowrap"
            >
                Voir toute l'équipe
                <span
                    class="iconify group-hover:translate-x-0.5 transition-transform text-dark-400 group-hover:text-dark-900"
                    data-icon="lucide:arrow-right"
                    data-width="16"
                    data-stroke-width="2"
                ></span>
            </a>
        </div>

        <!-- Expert Cards -->
        <div class="grid grid-cols-1 md:grid-cols-3 gap-6 lg:gap-8">
            {#each experts as expert, i}
                <div
                    class="group flex flex-col bg-white rounded-3xl border border-dark-100 hover:border-dark-200 hover:shadow-lg transition-all duration-300 overflow-hidden p-6 md:p-8"
                    use:reveal={{ preset: "fade-up", delay: 150 + i * 75 }}
                >
                    <!-- Header: photo + name/occupation -->
                    <div class="flex gap-4 items-center mb-4">
                        <img
                            src={expert.picture}
                            alt={expert.name}
                            class="rounded-2xl aspect-square w-14 md:w-16 flex-shrink-0 object-cover"
                        />
                        <div>
                            <h4
                                class="text-dark-900 font-semibold text-lg tracking-tight"
                            >
                                {expert.name}
                            </h4>
                            <p class="text-dark-500 text-sm font-medium mt-0.5">
                                {expert.occupation}
                            </p>
                        </div>
                    </div>

                    <!-- Tags -->
                    <div class="flex flex-wrap gap-2 mb-4">
                        {#each expert.tags as tag}
                            <div
                                class="inline-flex items-center bg-dark-50 border border-dark-100 rounded-full px-3 py-1.5"
                            >
                                <span class="text-xs font-medium text-dark-600"
                                    >{tag}</span
                                >
                            </div>
                        {/each}
                    </div>

                    <!-- Quote -->
                    <p
                        class="text-dark-500 text-sm leading-relaxed mb-6 flex-1"
                    >
                        "{expert.quote}"
                    </p>

                    <!-- CTAs -->
                    <div class="flex flex-col sm:flex-row gap-3 mt-auto">
                        <a
                            href="/book?partner={expert.id}"
                            class="inline-flex justify-center items-center gap-2 bg-dark-900 hover:bg-dark-800 text-white text-sm font-medium px-6 py-3 rounded-xl transition-all duration-200 shadow-sm hover:shadow-md"
                        >
                            Réserver
                            <span
                                class="iconify group-hover:translate-x-0.5 transition-transform"
                                data-icon="lucide:arrow-right"
                                data-width="16"
                                data-stroke-width="2"
                            ></span>
                        </a>
                        <a
                            href="/team/{expert.id}"
                            class="inline-flex justify-center items-center gap-2 bg-transparent hover:bg-dark-50 border border-dark-200 hover:border-dark-300 text-dark-600 hover:text-dark-900 text-sm font-medium px-6 py-3 rounded-xl transition-all duration-200"
                        >
                            Voir le profil
                        </a>
                    </div>
                </div>
            {/each}
        </div>
    </div>
</section>
