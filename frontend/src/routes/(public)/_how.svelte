<script lang="ts">
    import { reveal } from "$lib/actions/reveal";

    const steps = [
        {
            number: "01",
            icon: "lucide:search",
            title: "Choisissez votre expert",
            description:
                "Parcourez nos profils vérifiés et filtrez par spécialité pour trouver le praticien adapté à votre besoin.",
        },
        {
            number: "02",
            icon: "lucide:calendar-check",
            title: "Réservez en ligne",
            description:
                "Sélectionnez un créneau disponible et confirmez votre consultation en quelques secondes, sans appel téléphonique.",
        },
        {
            number: "03",
            icon: "lucide:user-check",
            title: "Vivez votre séance",
            description:
                "Rencontrez votre expert en cabinet ou en visio selon votre préférence. Flexible, simple, efficace.",
        },
        {
            number: "04",
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
            <h2
                class="mt-3 text-3xl md:text-4xl font-semibold tracking-tight text-dark-900"
            >
                Votre expert à portée de main en 4 étapes
            </h2>
            <p
                class="mt-4 text-base md:text-lg text-dark-500 leading-relaxed font-normal"
            >
                De la recherche à la consultation, nous avons simplifié chaque
                étape pour que vous puissiez vous concentrer sur l'essentiel.
            </p>
        </div>

        <!-- Steps Grid -->
        <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6 lg:gap-8">
            {#each steps as step, i}
                <div
                    class="relative flex flex-col items-start p-8 rounded-2xl bg-white border border-dark-100 hover:border-dark-200 hover:shadow-sm transition-all duration-200"
                    use:reveal={{ preset: "fade-up", delay: 150 + i * 75 }}
                >
                    <span
                        class="text-5xl font-semibold text-dark-100 leading-none mb-6 select-none"
                    >
                        {step.number}
                    </span>

                    <div
                        class="w-10 h-10 rounded-xl bg-dark-50 border border-dark-100 flex items-center justify-center mb-4"
                    >
                        <span
                            class="iconify text-dark-700"
                            data-icon={step.icon}
                            data-width="20"
                            data-stroke-width="1.5"
                        ></span>
                    </div>

                    <h3
                        class="text-base font-semibold text-dark-900 tracking-tight mb-2"
                    >
                        {step.title}
                    </h3>
                    <p class="text-sm text-dark-500 leading-relaxed">
                        {step.description}
                    </p>
                </div>
            {/each}
        </div>
    </div>
</section>
