<script lang="ts">
    import { reveal } from "$lib/actions/reveal";

    const testimonials = [
        {
            quote:
                "L'approche la plus complète en matière de santé que j'aie jamais connue. Les coachs se soucient véritablement de votre succès à long terme.",
            name: "Emma W.",
            role: "Patiente suivie depuis 2 ans",
            initials: "EW",
        },
        {
            quote:
                "En trois séances d'ostéopathie, mes douleurs chroniques au dos ont quasiment disparu. Un suivi personnalisé et bienveillant que je recommande à 100 %.",
            name: "Thomas R.",
            role: "Consultation en ostéopathie",
            initials: "TR",
        },
        {
            quote:
                "Le coaching en développement personnel m'a permis de reprendre confiance et de structurer mes objectifs. L'équipe est à l'écoute et vraiment professionnelle.",
            name: "Camille D.",
            role: "Coaching mindset & bien-être",
            initials: "CD",
        },
    ];
</script>

<section
    class="py-24 bg-gradient-to-b from-white to-dark-50/50 relative overflow-hidden"
>
    <!-- Background Decor -->
    <div
        class="absolute inset-0 opacity-[0.4] bg-[radial-gradient(ellipse_at_top,_var(--tw-gradient-stops))] from-dark-100 via-transparent to-transparent pointer-events-none"
    ></div>

    <div class="relative max-w-5xl mx-auto px-4 sm:px-6">
        <!-- Section Header -->
        <div class="text-center mb-16" use:reveal={{ preset: "fade-up", delay: 100 }}>
            <span
                class="text-sm font-semibold text-dark-500 uppercase tracking-wider"
                >Témoignages</span
            >
            <h2
                class="mt-3 text-3xl md:text-4xl font-semibold tracking-tight text-dark-900"
            >
                Ce que disent nos patients
            </h2>
            <p
                class="mt-4 text-base md:text-lg text-dark-500 leading-relaxed font-normal max-w-xl mx-auto"
            >
                Des résultats concrets, des relations durables. Voici le retour
                de celles et ceux qui nous font confiance au quotidien.
            </p>
        </div>

        <!-- Testimonials Grid -->
        <div class="grid grid-cols-1 md:grid-cols-3 gap-6 lg:gap-8">
            {#each testimonials as testimonial, i}
                <div
                    class="group h-full p-8 rounded-2xl bg-white border border-dark-100 hover:border-dark-200 shadow-sm hover:shadow-md transition-all duration-200 flex flex-col items-start text-left"
                    use:reveal={{ preset: "fade-up", delay: 150 + i * 80 }}
                >
                    <!-- Rating Stars -->
                    <div class="flex gap-1 mb-6">
                        {#each Array(5) as _}
                            <svg
                                xmlns="http://www.w3.org/2000/svg"
                                aria-hidden="true"
                                role="img"
                                width="16"
                                height="16"
                                viewBox="0 0 24 24"
                                class="iconify text-amber-400 iconify--lucide"
                            >
                                <path
                                    fill="currentColor"
                                    stroke="currentColor"
                                    stroke-linecap="round"
                                    stroke-linejoin="round"
                                    stroke-width="1"
                                    d="M11.525 2.295a.53.53 0 0 1 .95 0l2.31 4.679a2.12 2.12 0 0 0 1.595 1.16l5.166.756a.53.53 0 0 1 .294.904l-3.736 3.638a2.12 2.12 0 0 0-.611 1.878l.882 5.14a.53.53 0 0 1-.771.56l-4.618-2.428a2.12 2.12 0 0 0-1.973 0L6.396 21.01a.53.53 0 0 1-.77-.56l.881-5.139a2.12 2.12 0 0 0-.611-1.879L2.16 9.795a.53.53 0 0 1 .294-.906l5.165-.755a2.12 2.12 0 0 0 1.597-1.16z"
                                />
                            </svg>
                        {/each}
                    </div>

                    <!-- Quote -->
                    <blockquote class="text-base leading-relaxed text-dark-700 mb-8 flex-1">
                        "{testimonial.quote}"
                    </blockquote>

                    <!-- Attribution -->
                    <div class="flex items-center gap-3">
                        <div
                            class="w-10 h-10 rounded-full bg-dark-100 flex items-center justify-center text-xs font-semibold text-dark-600"
                        >
                            {testimonial.initials}
                        </div>
                        <div>
                            <div class="text-sm font-semibold text-dark-900">
                                {testimonial.name}
                            </div>
                            <div class="text-xs text-dark-500 flex items-center gap-1">
                                <span>{testimonial.role}</span>
                                <svg
                                    xmlns="http://www.w3.org/2000/svg"
                                    aria-hidden="true"
                                    role="img"
                                    width="12"
                                    height="12"
                                    viewBox="0 0 24 24"
                                    class="iconify text-blue-500 iconify--lucide"
                                >
                                    <g
                                        fill="none"
                                        stroke="currentColor"
                                        stroke-linecap="round"
                                        stroke-linejoin="round"
                                        stroke-width="2"
                                    >
                                        <path
                                            d="M3.85 8.62a4 4 0 0 1 4.78-4.77a4 4 0 0 1 6.74 0a4 4 0 0 1 4.78 4.78a4 4 0 0 1 0 6.74a4 4 0 0 1-4.77 4.78a4 4 0 0 1-6.75 0a4 4 0 0 1-4.78-4.77a4 4 0 0 1 0-6.76"
                                        ></path>
                                        <path d="m9 12l2 2l4-4"></path>
                                    </g>
                                </svg>
                            </div>
                        </div>
                    </div>
                </div>
            {/each}
        </div>

        <!-- CTA to team -->
        <div class="mt-12 text-center">
            <a
                href="/team"
                class="group inline-flex items-center gap-2 text-sm font-medium text-dark-500 hover:text-dark-900 transition-colors px-4 py-2 rounded-full hover:bg-dark-100"
            >
                Rencontrez l'équipe derrière les résultats
                <svg
                    xmlns="http://www.w3.org/2000/svg"
                    aria-hidden="true"
                    role="img"
                    width="16"
                    height="16"
                    viewBox="0 0 24 24"
                    class="iconify group-hover:translate-x-0.5 transition-transform text-dark-400 group-hover:text-dark-900 iconify--lucide"
                >
                    <path
                        fill="none"
                        stroke="currentColor"
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        stroke-width="2"
                        d="M5 12h14m-7-7l7 7l-7 7"
                    ></path>
                </svg>
            </a>
        </div>
    </div>
</section>
