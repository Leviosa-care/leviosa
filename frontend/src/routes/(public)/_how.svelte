<script lang="ts">
    import { reveal } from "$lib/actions/reveal";

    const steps = [
        {
            icon: "lucide:search",
            title: "Choisissez votre expert",
            description:
                "Parcourez nos profils vérifiés et filtrez par spécialité pour trouver le praticien adapté à votre besoin.",
        },
        {
            icon: "lucide:calendar-check",
            title: "Réservez en ligne",
            description:
                "Sélectionnez un créneau disponible et confirmez votre consultation en quelques secondes, sans appel téléphonique.",
        },
        {
            icon: "lucide:user-check",
            title: "Vivez votre séance",
            description:
                "Rencontrez votre expert en cabinet ou en visio selon votre préférence. Flexible, simple, efficace.",
        },
        {
            icon: "lucide:trending-up",
            title: "Progressez",
            description:
                "Recevez un compte-rendu personnalisé après chaque séance et revenez dès que vous en avez besoin.",
        },
    ];
</script>

<section id="how" class="py-20 lg:py-24 bg-dark-50">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <!-- Section Header -->
        <div
            class="max-w-2xl mx-auto text-center mb-16"
            use:reveal={{ preset: "fade-up", delay: 100 }}
        >
            <span class="text-sm font-semibold text-dark-500 uppercase tracking-wider">
                Comment ça fonctionne
            </span>
            <h2 class="mt-3 text-3xl md:text-4xl font-semibold tracking-tight text-dark-900">
                Votre expert à portée de main en 4 étapes
            </h2>
            <p class="mt-4 text-base md:text-lg text-dark-500 leading-relaxed font-normal">
                De la recherche à la consultation, nous avons simplifié chaque étape pour
                que vous puissiez vous concentrer sur l'essentiel.
            </p>
        </div>

        <!-- Steps Grid -->
        <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6 lg:gap-8">
            {#each steps as step, i}
                <div
                    class="group h-full p-8 rounded-2xl bg-white border border-dark-100 hover:border-dark-200 transition-colors duration-200 flex flex-col items-start text-left"
                    use:reveal={{ preset: "fade-up", delay: 150 + i * 75 }}
                >
                    <span class="text-xs font-medium text-dark-400 mb-4 tracking-widest">
                        0{i + 1}
                    </span>

                    <div
                        class="w-12 h-12 rounded-xl bg-dark-50 border border-dark-100 flex items-center justify-center shadow-sm mb-6 group-hover:scale-105 transition-transform duration-200"
                    >
                        <span
                            class="iconify text-dark-900"
                            data-icon={step.icon}
                            data-width="22"
                            data-stroke-width="1.5"
                        ></span>
                    </div>

                    <h3 class="text-base font-semibold text-dark-900 tracking-tight mb-2">
                        {step.title}
                    </h3>
                    <p class="text-sm leading-relaxed text-dark-500">
                        {step.description}
                    </p>
                </div>
            {/each}
        </div>
    </div>
</section>
