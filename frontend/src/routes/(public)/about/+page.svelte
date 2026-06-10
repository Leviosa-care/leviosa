<script lang="ts">
    import { reveal } from "$lib/actions/reveal";

    const domains = [
        {
            number: "01",
            label: "Massage",
            description: "Détente profonde et soulagement des tensions musculaires.",
        },
        {
            number: "02",
            label: "Psychologie",
            description: "Un espace d'écoute pour traverser vos moments difficiles.",
        },
        {
            number: "03",
            label: "Coaching",
            description: "Clarté et impact pour votre vie personnelle et professionnelle.",
        },
        {
            number: "04",
            label: "Méditation & pleine conscience",
            description: "Des pratiques pour retrouver calme et ancrage au quotidien.",
        },
    ];

    const spheres = [
        {
            label: "Bien-être physique",
            description: "Le corps libéré de ses tensions.",
        },
        {
            label: "Santé mentale",
            description: "L'esprit apaisé, ancré dans le présent.",
        },
        {
            label: "Développement personnel",
            description: "La trajectoire que vous méritez.",
        },
    ];
</script>

<!-- Hero -->
<section class="py-28 sm:py-36 bg-background-section">
    <div
        class="max-w-3xl mx-auto px-4 sm:px-6 text-center"
        use:reveal={{ preset: "fade-up", delay: 100 }}
    >
        <h1
            class="text-4xl sm:text-5xl lg:text-6xl font-semibold tracking-tight text-foreground leading-[1.1] mb-6"
        >
            Prendre soin de soi,<br />simplement.
        </h1>
        <p class="text-lg text-muted-foreground leading-relaxed max-w-xl mx-auto">
            Un accès direct aux experts qui veillent sur votre corps et votre
            esprit — où que vous en soyez dans la vie.
        </p>
    </div>
</section>

<!-- Mission -->
<section class="py-20 bg-white">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div
            class="max-w-2xl mx-auto text-center"
            use:reveal={{ preset: "fade-up", delay: 100 }}
        >
            <span
                class="text-xs font-medium text-muted-foreground uppercase tracking-widest bg-surface-hover px-3 py-1 rounded-full"
            >
                Notre Mission
            </span>
            <h2
                class="mt-6 text-3xl md:text-4xl font-semibold tracking-tight text-foreground"
            >
                Né d'un constat simple.
            </h2>
            <p class="mt-4 text-base md:text-lg text-muted-foreground leading-relaxed">
                Trouver le bon expert prend trop de temps. Les échanges
                s'étirent, les disponibilités restent floues. Leviosa change ça
                — une réservation directe, des praticiens vérifiés, un
                soulagement à portée de main.
            </p>
        </div>
    </div>
</section>

<!-- Ce qui nous rassemble -->
<section class="py-20 bg-foreground relative overflow-hidden">
    <div
        class="absolute inset-0 opacity-20 bg-[radial-gradient(ellipse_at_center,_var(--tw-gradient-stops))] from-white via-transparent to-transparent pointer-events-none"
    ></div>

    <div class="relative max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div
            class="max-w-2xl mx-auto text-center mb-12"
            use:reveal={{ preset: "fade-up", delay: 100 }}
        >
            <span
                class="text-xs font-medium text-white/50 uppercase tracking-widest bg-white/10 border border-white/20 px-3 py-1 rounded-full"
            >
                Ce qui nous rassemble
            </span>
            <h2
                class="mt-6 text-3xl md:text-4xl font-semibold tracking-tight text-white"
            >
                Votre bien-être, sous toutes ses formes.
            </h2>
            <p class="mt-4 text-base md:text-lg text-white/60 leading-relaxed">
                Peu importe votre âge ou vos besoins, Leviosa réunit les
                experts adaptés à votre situation.
            </p>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
            {#each spheres as sphere, i}
                <div
                    class="p-8 rounded-2xl bg-white/10 border border-white/20 backdrop-blur-sm"
                    use:reveal={{ preset: "fade-up", delay: 150 + i * 75 }}
                >
                    <h3 class="text-lg font-semibold text-white">
                        {sphere.label}
                    </h3>
                    <p class="mt-2 text-sm text-white/60 leading-relaxed">
                        {sphere.description}
                    </p>
                </div>
            {/each}
        </div>
    </div>
</section>

<!-- Domaines couverts -->
<section class="py-20 bg-white">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div
            class="max-w-2xl mx-auto text-center mb-12"
            use:reveal={{ preset: "fade-up", delay: 100 }}
        >
            <span
                class="text-xs font-medium text-muted-foreground uppercase tracking-widest bg-surface-hover px-3 py-1 rounded-full"
            >
                Nos spécialités
            </span>
            <h2
                class="mt-6 text-3xl md:text-4xl font-semibold tracking-tight text-foreground"
            >
                Ce que vous trouverez ici.
            </h2>
        </div>

        <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
            {#each domains as domain, i}
                <div
                    class="group p-8 rounded-2xl bg-surface border border-border-input hover:border-border-input-hover transition-colors duration-200 flex flex-col"
                    use:reveal={{ preset: "fade-up", delay: 150 + i * 75 }}
                >
                    <span
                        class="text-xs font-medium text-muted-foreground font-mono tracking-wider"
                        >{domain.number}</span
                    >
                    <h3
                        class="mt-4 text-lg font-semibold text-foreground tracking-tight"
                    >
                        {domain.label}
                    </h3>
                    <p class="mt-2 text-sm text-muted-foreground leading-relaxed">
                        {domain.description}
                    </p>
                </div>
            {/each}
        </div>
    </div>
</section>

<!-- Nos valeurs -->
<section class="py-20 bg-background-section">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div
            class="max-w-2xl mx-auto text-center mb-12"
            use:reveal={{ preset: "fade-up", delay: 100 }}
        >
            <span
                class="text-xs font-medium text-muted-foreground uppercase tracking-widest bg-surface-hover px-3 py-1 rounded-full"
            >
                Nos valeurs
            </span>
            <h2
                class="mt-6 text-3xl md:text-4xl font-semibold tracking-tight text-foreground"
            >
                Fondé sur la clarté, pas la complexité.
            </h2>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-3 gap-6 lg:gap-8">
            <div
                class="group h-full p-8 rounded-2xl bg-white border border-border-input hover:border-border-input-hover transition-colors duration-200 flex flex-col items-start text-left"
                use:reveal={{ preset: "fade-up", delay: 150 }}
            >
                <div
                    class="w-12 h-12 rounded-xl bg-surface border border-border-input-hover flex items-center justify-center shadow-mini mb-6 group-hover:scale-105 transition-transform duration-200"
                >
                    <svg
                        xmlns="http://www.w3.org/2000/svg"
                        width="22"
                        height="22"
                        viewBox="0 0 24 24"
                        class="text-foreground"
                    >
                        <g
                            fill="none"
                            stroke="currentColor"
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            stroke-width="2"
                        >
                            <path
                                d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2M16 3.128a4 4 0 0 1 0 7.744M22 21v-2a4 4 0 0 0-3-3.87"
                            ></path>
                            <circle cx="9" cy="7" r="4"></circle>
                        </g>
                    </svg>
                </div>
                <h3 class="text-lg font-semibold text-foreground tracking-tight mb-3">
                    Accès direct aux experts
                </h3>
                <p class="text-sm leading-relaxed text-muted-foreground">
                    Vous échangez directement avec les praticiens. Pas
                    d'intermédiaires, pas de couches de traduction, juste une
                    collaboration efficace.
                </p>
            </div>

            <div
                class="group h-full p-8 rounded-2xl bg-white border border-border-input hover:border-border-input-hover transition-colors duration-200 flex flex-col items-start text-left"
                use:reveal={{ preset: "fade-up", delay: 225 }}
            >
                <div
                    class="w-12 h-12 rounded-xl bg-surface border border-border-input-hover flex items-center justify-center shadow-mini mb-6 group-hover:scale-105 transition-transform duration-200"
                >
                    <svg
                        xmlns="http://www.w3.org/2000/svg"
                        width="22"
                        height="22"
                        viewBox="0 0 24 24"
                        class="text-foreground"
                    >
                        <g
                            fill="none"
                            stroke="currentColor"
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            stroke-width="2"
                        >
                            <path
                                d="m15 12l-9.373 9.373a1 1 0 0 1-3.001-3L12 9m6 6l4-4"
                            ></path>
                            <path
                                d="m21.5 11.5l-1.914-1.914A2 2 0 0 1 19 8.172v-.344a2 2 0 0 0-.586-1.414l-1.657-1.657A6 6 0 0 0 12.516 3H9l1.243 1.243A6 6 0 0 1 12 8.485V10l2 2h1.172a2 2 0 0 1 1.414.586L18.5 14.5"
                            ></path>
                        </g>
                    </svg>
                </div>
                <h3 class="text-lg font-semibold text-foreground tracking-tight mb-3">
                    Exécution pragmatique
                </h3>
                <p class="text-sm leading-relaxed text-muted-foreground">
                    Notre priorité est de vous apporter une valeur concrète dès
                    la première séance. Pas de promesses creuses, des résultats
                    réels.
                </p>
            </div>

            <div
                class="group h-full p-8 rounded-2xl bg-white border border-border-input hover:border-border-input-hover transition-colors duration-200 flex flex-col items-start text-left"
                use:reveal={{ preset: "fade-up", delay: 300 }}
            >
                <div
                    class="w-12 h-12 rounded-xl bg-surface border border-border-input-hover flex items-center justify-center shadow-mini mb-6 group-hover:scale-105 transition-transform duration-200"
                >
                    <svg
                        xmlns="http://www.w3.org/2000/svg"
                        width="22"
                        height="22"
                        viewBox="0 0 24 24"
                        class="text-foreground"
                    >
                        <g
                            fill="none"
                            stroke="currentColor"
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            stroke-width="2"
                        >
                            <path
                                d="M20 13c0 5-3.5 7.5-7.66 8.95a1 1 0 0 1-.67-.01C7.5 20.5 4 18 4 13V6a1 1 0 0 1 1-1c2 0 4.5-1.2 6.24-2.72a1.17 1.17 0 0 1 1.52 0C14.51 3.81 17 5 19 5a1 1 0 0 1 1 1z"
                            ></path>
                            <path d="m9 12l2 2l4-4"></path>
                        </g>
                    </svg>
                </div>
                <h3 class="text-lg font-semibold text-foreground tracking-tight mb-3">
                    Transparence totale
                </h3>
                <p class="text-sm leading-relaxed text-muted-foreground">
                    Délais clairs, tarification fixe et retours honnêtes. Si un
                    expert n'est pas adapté à votre besoin, nous vous le dirons
                    immédiatement.
                </p>
            </div>
        </div>
    </div>
</section>

<!-- CTA final -->
<section class="py-24 bg-background-section relative overflow-hidden">
    <div
        class="absolute top-0 left-1/2 -translate-x-1/2 w-full h-full max-w-3xl opacity-40 bg-[radial-gradient(ellipse_at_top,_var(--tw-gradient-stops))] from-indigo-100 via-transparent to-transparent pointer-events-none"
    ></div>

    <div class="relative max-w-3xl mx-auto px-6 text-center z-10">
        <h2
            class="text-3xl sm:text-4xl md:text-5xl font-semibold tracking-tight text-foreground mb-6"
            use:reveal={{ preset: "fade-up", delay: 100 }}
        >
            Prêt·e à prendre soin de vous ?
        </h2>
        <p
            class="text-lg text-muted-foreground mb-10 max-w-lg mx-auto leading-relaxed"
            use:reveal={{ preset: "fade-up", delay: 150 }}
        >
            Rejoignez les milliers de personnes qui ont déjà trouvé leur
            expert. Ou rejoignez notre réseau en tant que praticien.
        </p>
        <div
            class="flex flex-col sm:flex-row justify-center gap-4"
            use:reveal={{ preset: "fade-up", delay: 200 }}
        >
            <a
                href="/book"
                class="inline-flex justify-center items-center gap-2 bg-foreground hover:bg-foreground-alt text-white text-sm font-medium px-8 py-3.5 rounded-xl transition-all duration-200 shadow-mini hover:shadow-card"
            >
                Prendre rendez-vous
                <span
                    class="iconify group-hover:translate-x-0.5 transition-transform"
                    data-icon="lucide:arrow-right"
                    data-width="16"
                    data-stroke-width="2"
                ></span>
            </a>
            <a
                href="/auth"
                class="inline-flex justify-center items-center gap-2 bg-transparent hover:bg-surface-hover border border-border-input-hover hover:border-border-input-hover text-foreground-alt hover:text-foreground text-sm font-medium px-8 py-3.5 rounded-xl transition-all duration-200"
            >
                Rejoindre en tant qu'expert
            </a>
        </div>
    </div>
</section>
