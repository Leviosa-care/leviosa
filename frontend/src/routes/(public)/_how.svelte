<script lang="ts">
    import { reveal } from "$lib/actions/reveal";
</script>

<section id="how" class="relative py-24 sm:py-32 bg-dark-25 overflow-hidden">
    <div class="max-w-7xl mx-auto px-6 lg:px-8">
        <!-- Section Header -->
        <div class="max-w-2xl mx-auto text-center" use:reveal={{ preset: "fade-up", delay: 100 }}>
            <span class="text-xs font-medium text-dark-500 uppercase tracking-widest bg-dark-100 px-3 py-1 rounded-full">
                Comment ça fonctionne
            </span>
            <h2 class="mt-6 text-3xl sm:text-4xl font-semibold tracking-tight text-dark-900">
                Votre expert à portée de main en 4 étapes
            </h2>
            <p class="mt-4 text-base sm:text-lg text-dark-500 leading-relaxed">
                De la recherche à la consultation, nous avons simplifié chaque étape pour
                que vous puissiez vous concentrer sur l'essentiel.
            </p>
        </div>

        <!-- Steps Grid -->
        <div class="mt-16 sm:mt-20 lg:mt-24 grid grid-cols-1 gap-8 lg:grid-cols-4 lg:gap-6 xl:gap-8">

            <!-- Step 1: Choisissez votre expert -->
            <div
                class="flex flex-col rounded-2xl border border-dark-200 bg-white p-8 shadow-sm transition-all hover:shadow-md hover:border-dark-300 group"
                use:reveal={{ preset: "fade-up", delay: 150 }}
            >
                <div class="h-48 rounded-xl bg-dark-50/80 border border-dark-100 mb-8 p-5 flex flex-col gap-4 relative overflow-hidden group-hover:bg-dark-50 transition-colors">
                    <!-- Search bar -->
                    <div class="h-10 bg-white border border-dark-200 rounded-lg flex items-center px-3 shadow-sm relative z-10">
                        <span data-icon="solar:magnifer-linear" data-width="18" class="text-dark-400"></span>
                        <div class="w-24 h-1.5 bg-dark-200 rounded-full ml-3"></div>
                        <div class="ml-auto w-5 h-5 bg-dark-100 rounded flex items-center justify-center border border-dark-200">
                            <span data-icon="solar:tuning-linear" data-width="12" class="text-dark-400"></span>
                        </div>
                    </div>
                    <!-- Profile rows -->
                    <div class="flex flex-col gap-3 relative z-10">
                        <div class="bg-white p-3 rounded-lg border border-dark-100 flex items-center gap-4 shadow-sm w-11/12 transform transition-transform group-hover:translate-x-1">
                            <div class="w-8 h-8 rounded-full bg-dark-200 flex-shrink-0"></div>
                            <div class="flex flex-col gap-1.5 w-full">
                                <div class="w-16 h-1.5 bg-dark-300 rounded-full"></div>
                                <div class="w-10 h-1 bg-dark-200 rounded-full"></div>
                            </div>
                            <div class="w-4 h-4 rounded-full border border-dark-200 ml-auto"></div>
                        </div>
                        <div class="bg-white p-3 rounded-lg border border-dark-100 flex items-center gap-4 shadow-sm w-full opacity-60 ml-4">
                            <div class="w-8 h-8 rounded-full bg-dark-200 flex-shrink-0"></div>
                            <div class="flex flex-col gap-1.5 w-full">
                                <div class="w-20 h-1.5 bg-dark-300 rounded-full"></div>
                                <div class="w-12 h-1 bg-dark-200 rounded-full"></div>
                            </div>
                        </div>
                    </div>
                    <div class="absolute bottom-0 left-0 right-0 h-16 bg-gradient-to-t from-dark-50 to-transparent z-20"></div>
                </div>
                <div class="flex items-center gap-3 mb-3">
                    <span class="text-xs font-medium text-dark-400 font-mono tracking-wider">01</span>
                    <h3 class="text-lg font-semibold tracking-tight text-dark-900">Choisissez votre expert</h3>
                </div>
                <p class="text-sm text-dark-500 leading-relaxed flex-grow">
                    Parcourez nos profils vérifiés et filtrez par spécialité pour trouver le praticien adapté à votre besoin.
                </p>
            </div>

            <!-- Step 2: Réservez en ligne -->
            <div
                class="flex flex-col rounded-2xl border border-dark-200 bg-white p-8 shadow-sm transition-all hover:shadow-md hover:border-dark-300 group"
                use:reveal={{ preset: "fade-up", delay: 225 }}
            >
                <div class="h-48 rounded-xl bg-dark-50/80 border border-dark-100 mb-8 p-5 flex flex-col items-center justify-center relative overflow-hidden group-hover:bg-dark-50 transition-colors">
                    <div class="bg-white border border-dark-200 rounded-xl p-4 shadow-sm w-full max-w-[180px] relative z-10 transform transition-transform group-hover:scale-105">
                        <div class="flex justify-between items-center mb-4">
                            <div class="w-12 h-1.5 bg-dark-300 rounded-full"></div>
                            <div class="flex gap-1">
                                <span data-icon="solar:alt-arrow-left-linear" data-width="12" class="text-dark-300"></span>
                                <span data-icon="solar:alt-arrow-right-linear" data-width="12" class="text-dark-400"></span>
                            </div>
                        </div>
                        <div class="grid grid-cols-4 gap-2">
                            <div class="aspect-square rounded bg-dark-50 border border-dark-100"></div>
                            <div class="aspect-square rounded bg-dark-50 border border-dark-100"></div>
                            <div class="aspect-square rounded bg-dark-900 shadow-sm ring-2 ring-dark-900/20 flex items-center justify-center">
                                <div class="w-1 h-1 bg-white rounded-full"></div>
                            </div>
                            <div class="aspect-square rounded bg-dark-50 border border-dark-100"></div>
                            <div class="aspect-square rounded bg-dark-50 border border-dark-100"></div>
                            <div class="aspect-square rounded bg-dark-50 border border-dark-100"></div>
                            <div class="aspect-square rounded bg-dark-50 border border-dark-100"></div>
                            <div class="aspect-square rounded bg-dark-50 border border-dark-100"></div>
                        </div>
                        <div class="mt-4 w-full h-8 bg-dark-900 rounded-lg text-white flex items-center justify-center text-xs font-medium shadow-sm">
                            Réserver
                        </div>
                    </div>
                </div>
                <div class="flex items-center gap-3 mb-3">
                    <span class="text-xs font-medium text-dark-400 font-mono tracking-wider">02</span>
                    <h3 class="text-lg font-semibold tracking-tight text-dark-900">Réservez en ligne</h3>
                </div>
                <p class="text-sm text-dark-500 leading-relaxed flex-grow">
                    Sélectionnez un créneau disponible et confirmez votre consultation en quelques secondes, sans appel téléphonique.
                </p>
            </div>

            <!-- Step 3: Vivez votre séance -->
            <div
                class="flex flex-col rounded-2xl border border-dark-200 bg-white p-8 shadow-sm transition-all hover:shadow-md hover:border-dark-300 group"
                use:reveal={{ preset: "fade-up", delay: 300 }}
            >
                <div class="h-48 rounded-xl bg-dark-50/80 border border-dark-100 mb-8 p-5 flex flex-col items-center justify-center relative overflow-hidden group-hover:bg-dark-50 transition-colors">
                    <div class="relative w-full max-w-[200px] aspect-[16/10] bg-dark-900 rounded-xl shadow-md overflow-hidden z-10">
                        <div class="absolute inset-0 flex items-center justify-center opacity-50">
                            <span data-icon="solar:user-linear" data-width="40" class="text-dark-400"></span>
                        </div>
                        <div class="absolute bottom-2 right-2 w-12 h-16 bg-dark-800 border border-dark-700 rounded-lg shadow-sm"></div>
                        <div class="absolute inset-0 bg-white/95 backdrop-blur-sm flex flex-col items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity duration-300 ease-out">
                            <div class="w-8 h-8 rounded-full bg-green-100 flex items-center justify-center mb-2">
                                <span data-icon="solar:check-read-linear" data-width="16" class="text-green-600"></span>
                            </div>
                            <div class="w-16 h-1.5 bg-dark-200 rounded-full mb-1"></div>
                            <div class="w-10 h-1 bg-dark-200 rounded-full"></div>
                        </div>
                    </div>
                    <div class="absolute -right-2 top-4 w-16 h-16 bg-dark-200 rounded-full blur-3xl opacity-50"></div>
                    <div class="absolute -left-4 bottom-4 w-20 h-20 bg-dark-200 rounded-full blur-3xl opacity-50"></div>
                </div>
                <div class="flex items-center gap-3 mb-3">
                    <span class="text-xs font-medium text-dark-400 font-mono tracking-wider">03</span>
                    <h3 class="text-lg font-semibold tracking-tight text-dark-900">Vivez votre séance</h3>
                </div>
                <p class="text-sm text-dark-500 leading-relaxed flex-grow">
                    Rencontrez votre expert en cabinet ou en visio selon votre préférence. Flexible, simple, efficace.
                </p>
            </div>

            <!-- Step 4: Progressez -->
            <div
                class="flex flex-col rounded-2xl border border-dark-200 bg-white p-8 shadow-sm transition-all hover:shadow-md hover:border-dark-300 group"
                use:reveal={{ preset: "fade-up", delay: 375 }}
            >
                <div class="h-48 rounded-xl bg-dark-50/80 border border-dark-100 mb-8 p-5 flex flex-col items-center justify-center relative overflow-hidden group-hover:bg-dark-50 transition-colors">
                    <div class="bg-white border border-dark-200 rounded-xl p-4 shadow-sm w-full max-w-[160px] relative z-10 flex flex-col gap-2 transform transition-transform group-hover:scale-105">
                        <div class="flex justify-between items-center mb-2">
                            <div class="w-12 h-1.5 bg-dark-300 rounded-full"></div>
                            <span data-icon="solar:chart-square-linear" data-width="12" class="text-dark-400"></span>
                        </div>
                        <div class="flex items-end gap-1.5 h-16">
                            <div class="w-full bg-dark-100 rounded-t h-1/3"></div>
                            <div class="w-full bg-dark-200 rounded-t h-1/2"></div>
                            <div class="w-full bg-dark-300 rounded-t h-2/3"></div>
                            <div class="w-full bg-dark-900 rounded-t h-full shadow-sm"></div>
                        </div>
                    </div>
                    <div class="absolute -right-4 bottom-4 w-20 h-20 bg-dark-200 rounded-full blur-3xl opacity-50"></div>
                </div>
                <div class="flex items-center gap-3 mb-3">
                    <span class="text-xs font-medium text-dark-400 font-mono tracking-wider">04</span>
                    <h3 class="text-lg font-semibold tracking-tight text-dark-900">Progressez</h3>
                </div>
                <p class="text-sm text-dark-500 leading-relaxed flex-grow">
                    Recevez un compte-rendu personnalisé après chaque séance et revenez dès que vous en avez besoin.
                </p>
            </div>

        </div>
    </div>
</section>
