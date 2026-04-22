<script lang="ts">
    import { reveal } from "$lib/actions/reveal";
</script>

<section id="services" class="py-20 lg:py-24 bg-white border-b border-dark-100">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <!-- Section Header -->
        <div
            class="flex flex-col md:flex-row md:items-end justify-between gap-6 mb-12"
            use:reveal={{ preset: "fade-up", delay: 100 }}
        >
            <div class="max-w-xl">
                <span
                    class="text-sm font-semibold text-dark-500 uppercase tracking-wider"
                    >Notre Expertise</span
                >
                <h2
                    class="mt-3 text-3xl md:text-4xl font-semibold tracking-tight text-dark-900"
                >
                    Conçu sur mesure pour vos besoins
                </h2>
                <p class="mt-4 text-base text-dark-500">
                    Nous proposons actuellement un service spécialisé conçu pour
                    vous apporter une valeur immédiate et des actions concrètes.
                </p>
            </div>
        </div>

        <!-- Service Card: Single Item Layout -->
        <div
            class="group relative bg-white rounded-3xl border border-dark-200 shadow-sm overflow-hidden hover:shadow-md transition-all duration-300"
            use:reveal={{ preset: "fade-up", delay: 150 }}
        >
            <div class="grid grid-cols-1 lg:grid-cols-2">
                <!-- Content Side -->
                <div
                    class="p-8 lg:p-12 flex flex-col justify-center order-2 lg:order-1"
                >
                    <div class="flex items-center gap-3 mb-6">
                        <div
                            class="w-10 h-10 rounded-lg bg-blue-50 flex items-center justify-center text-blue-600"
                        >
                            <svg
                                xmlns="http://www.w3.org/2000/svg"
                                aria-hidden="true"
                                role="img"
                                width="20"
                                height="20"
                                viewBox="0 0 24 24"
                                data-icon="lucide:zap"
                                data-width="20"
                                data-stroke-width="1.5"
                                class="iconify iconify--lucide"
                                ><path
                                    fill="none"
                                    stroke="currentColor"
                                    stroke-linecap="round"
                                    stroke-linejoin="round"
                                    stroke-width="2"
                                    d="M4 14a1 1 0 0 1-.78-1.63l9.9-10.2a.5.5 0 0 1 .86.46l-1.92 6.02A1 1 0 0 0 13 10h7a1 1 0 0 1 .78 1.63l-9.9 10.2a.5.5 0 0 1-.86-.46l1.92-6.02A1 1 0 0 0 11 14z"
                                ></path></svg
                            >
                        </div>
                        <span
                            class="text-sm font-semibold text-blue-600 tracking-wide uppercase"
                            >Le Plus Populaire</span
                        >
                    </div>

                    <h3
                        class="text-2xl md:text-3xl font-semibold text-dark-900 tracking-tight mb-4"
                    >
                        Consultation Stratégique
                    </h3>

                    <p
                        class="text-dark-500 text-base md:text-lg leading-relaxed mb-8"
                    >
                        Une séance de 60 minutes pour analyser vos défis actuels
                        et élaborer une feuille de route concrète, étape par étape.
                        Idéal pour les fondateurs en phase de démarrage.
                    </p>

                    <div class="space-y-4 mb-10">
                        <!-- Benefit 1 -->
                        <div class="flex gap-3 items-start">
                            <div
                                class="mt-1 flex-shrink-0 w-5 h-5 rounded-full bg-dark-100 flex items-center justify-center"
                            >
                                <svg
                                    xmlns="http://www.w3.org/2000/svg"
                                    aria-hidden="true"
                                    role="img"
                                    width="12"
                                    height="12"
                                    viewBox="0 0 24 24"
                                    data-icon="lucide:check"
                                    data-width="12"
                                    data-stroke-width="2"
                                    class="iconify text-dark-900 iconify--lucide"
                                    ><path
                                        fill="none"
                                        stroke="currentColor"
                                        stroke-linecap="round"
                                        stroke-linejoin="round"
                                        stroke-width="2"
                                        d="M20 6L9 17l-5-5"
                                    ></path></svg
                                >
                            </div>
                            <div>
                                <p class="text-sm font-medium text-dark-900">
                                    Audit Approfondi
                                </p>
                                <p class="text-xs text-dark-500 mt-0.5">
                                    Analyse complète de votre stack technique et
                                    de vos indicateurs.
                                </p>
                            </div>
                        </div>
                        <!-- Benefit 2 -->
                        <div class="flex gap-3 items-start">
                            <div
                                class="mt-1 flex-shrink-0 w-5 h-5 rounded-full bg-dark-100 flex items-center justify-center"
                            >
                                <svg
                                    xmlns="http://www.w3.org/2000/svg"
                                    aria-hidden="true"
                                    role="img"
                                    width="12"
                                    height="12"
                                    viewBox="0 0 24 24"
                                    data-icon="lucide:check"
                                    data-width="12"
                                    data-stroke-width="2"
                                    class="iconify text-dark-900 iconify--lucide"
                                    ><path
                                        fill="none"
                                        stroke="currentColor"
                                        stroke-linecap="round"
                                        stroke-linejoin="round"
                                        stroke-width="2"
                                        d="M20 6L9 17l-5-5"
                                    ></path></svg
                                >
                            </div>
                            <div>
                                <p class="text-sm font-medium text-dark-900">
                                    Feuille de Route Actionnable
                                </p>
                                <p class="text-xs text-dark-500 mt-0.5">
                                    Repartez avec un plan clair de 30-60-90 jours.
                                </p>
                            </div>
                        </div>
                        <!-- Benefit 3 -->
                        <div class="flex gap-3 items-start">
                            <div
                                class="mt-1 flex-shrink-0 w-5 h-5 rounded-full bg-dark-100 flex items-center justify-center"
                            >
                                <svg
                                    xmlns="http://www.w3.org/2000/svg"
                                    aria-hidden="true"
                                    role="img"
                                    width="12"
                                    height="12"
                                    viewBox="0 0 24 24"
                                    data-icon="lucide:check"
                                    data-width="12"
                                    data-stroke-width="2"
                                    class="iconify text-dark-900 iconify--lucide"
                                    ><path
                                        fill="none"
                                        stroke="currentColor"
                                        stroke-linecap="round"
                                        stroke-linejoin="round"
                                        stroke-width="2"
                                        d="M20 6L9 17l-5-5"
                                    ></path></svg
                                >
                            </div>
                            <div>
                                <p class="text-sm font-medium text-dark-900">
                                    Enregistrement &amp; Notes
                                </p>
                                <p class="text-xs text-dark-500 mt-0.5">
                                    Transcription complète et liens vers les ressources fournis.
                                </p>
                            </div>
                        </div>
                    </div>

                    <!-- Secondary CTA (Discovery) -->
                    <div class="flex flex-col sm:flex-row gap-4">
                        <button
                            class="inline-flex justify-center items-center gap-2 bg-white text-dark-900 border border-dark-200 hover:border-dark-400 hover:bg-dark-50 text-sm font-medium px-6 py-3 rounded-xl transition-all duration-200 cursor-pointer"
                        >
                            Découvrir le service
                        </button>
                    </div>
                </div>

                <!-- Visual Side -->
                <div
                    class="relative bg-dark-50 order-1 lg:order-2 border-b lg:border-b-0 lg:border-l border-dark-100 min-h-[300px] lg:min-h-full overflow-hidden flex items-center justify-center p-8"
                >
                    <!-- Decorative Pattern -->
                    <div
                        class="absolute inset-0 opacity-[0.03]"
                        style="background-image: radial-gradient(#0f172a 1px, transparent 1px); background-size: 20px 20px;"
                    ></div>

                    <!-- UI Representation: Roadmap Visualization -->
                    <div class="relative w-full max-w-sm">
                        <div
                            class="absolute top-0 bottom-0 left-6 w-0.5 bg-dark-200"
                        ></div>

                        <!-- Item 1 -->
                        <div class="relative flex items-center gap-4 mb-6">
                            <div
                                class="w-12 h-12 rounded-full border-4 border-dark-50 bg-white shadow-sm flex items-center justify-center z-10"
                            >
                                <svg
                                    xmlns="http://www.w3.org/2000/svg"
                                    aria-hidden="true"
                                    role="img"
                                    width="18"
                                    height="18"
                                    viewBox="0 0 24 24"
                                    data-icon="lucide:search"
                                    data-width="18"
                                    class="iconify text-dark-400 iconify--lucide"
                                    ><g
                                        fill="none"
                                        stroke="currentColor"
                                        stroke-linecap="round"
                                        stroke-linejoin="round"
                                        stroke-width="2"
                                        ><path d="m21 21l-4.34-4.34"
                                        ></path><circle cx="11" cy="11" r="8"
                                        ></circle></g
                                    ></svg
                                >
                            </div>
                            <div
                                class="flex-1 bg-white p-3 rounded-lg border border-dark-100 shadow-sm opacity-60"
                            >
                                <div
                                    class="h-2 w-16 bg-dark-100 rounded mb-2"
                                ></div>
                                <div
                                    class="h-1.5 w-full bg-dark-50 rounded"
                                ></div>
                            </div>
                        </div>

                        <!-- Item 2 (Active) -->
                        <div class="relative flex items-center gap-4 mb-6">
                            <div
                                class="w-12 h-12 rounded-full border-4 border-dark-50 bg-dark-900 shadow-lg shadow-dark-200 flex items-center justify-center z-10"
                            >
                                <svg
                                    xmlns="http://www.w3.org/2000/svg"
                                    aria-hidden="true"
                                    role="img"
                                    width="18"
                                    height="18"
                                    viewBox="0 0 24 24"
                                    data-icon="lucide:lightbulb"
                                    data-width="18"
                                    class="iconify text-white iconify--lucide"
                                    ><path
                                        fill="none"
                                        stroke="currentColor"
                                        stroke-linecap="round"
                                        stroke-linejoin="round"
                                        stroke-width="2"
                                        d="M15 14c.2-1 .7-1.7 1.5-2.5c1-.9 1.5-2.2 1.5-3.5A6 6 0 0 0 6 8c0 1 .2 2.2 1.5 3.5c.7.7 1.3 1.5 1.5 2.5m0 4h6m-5 4h4"
                                    ></path></svg
                                >
                            </div>
                            <div
                                class="flex-1 bg-white p-4 rounded-xl border border-dark-200 shadow-md"
                            >
                                <div
                                    class="flex justify-between items-start mb-2"
                                >
                                    <div
                                        class="h-2.5 w-24 bg-dark-900 rounded"
                                    ></div>
                                    <div
                                        class="h-4 w-4 rounded-full bg-emerald-100 border border-emerald-200 flex items-center justify-center"
                                    >
                                        <svg
                                            xmlns="http://www.w3.org/2000/svg"
                                            aria-hidden="true"
                                            role="img"
                                            width="10"
                                            height="10"
                                            viewBox="0 0 24 24"
                                            data-icon="lucide:check"
                                            data-width="10"
                                            class="iconify text-emerald-600 iconify--lucide"
                                            ><path
                                                fill="none"
                                                stroke="currentColor"
                                                stroke-linecap="round"
                                                stroke-linejoin="round"
                                                stroke-width="2"
                                                d="M20 6L9 17l-5-5"
                                            ></path></svg
                                        >
                                    </div>
                                </div>
                                <div class="space-y-1.5">
                                    <div
                                        class="h-2 w-full bg-dark-100 rounded"
                                    ></div>
                                    <div
                                        class="h-2 w-2/3 bg-dark-100 rounded"
                                    ></div>
                                </div>
                            </div>
                        </div>

                        <!-- Item 3 -->
                        <div class="relative flex items-center gap-4">
                            <div
                                class="w-12 h-12 rounded-full border-4 border-dark-50 bg-white shadow-sm flex items-center justify-center z-10"
                            >
                                <svg
                                    xmlns="http://www.w3.org/2000/svg"
                                    aria-hidden="true"
                                    role="img"
                                    width="18"
                                    height="18"
                                    viewBox="0 0 24 24"
                                    data-icon="lucide:trending-up"
                                    data-width="18"
                                    class="iconify text-dark-400 iconify--lucide"
                                    ><g
                                        fill="none"
                                        stroke="currentColor"
                                        stroke-linecap="round"
                                        stroke-linejoin="round"
                                        stroke-width="2"
                                        ><path d="M16 7h6v6"></path><path
                                            d="m22 7l-8.5 8.5l-5-5L2 17"
                                        ></path></g
                                    ></svg
                                >
                            </div>
                            <div
                                class="flex-1 bg-white p-3 rounded-lg border border-dark-100 shadow-sm opacity-60"
                            >
                                <div
                                    class="h-2 w-20 bg-dark-100 rounded mb-2"
                                ></div>
                                <div
                                    class="h-1.5 w-full bg-dark-50 rounded"
                                ></div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</section>
