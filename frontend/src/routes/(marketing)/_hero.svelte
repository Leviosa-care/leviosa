<script lang="ts">
    import Button from "$lib/ui/Button.svelte";
</script>

<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-12 lg:gap-30 items-center">
        <!-- Left Column: Content (Priority 1) -->
        <div class="flex flex-col items-start text-left z-10">
            <!-- Trust Badge -->
            <div
                class="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-dark-50 border border-dark-200 mb-6"
            >
                <span class="relative flex h-2 w-2">
                    <span
                        class="animate-ping absolute inline-flex h-full w-full rounded-full bg-emerald-400 opacity-75"
                    ></span>
                    <span
                        class="relative inline-flex rounded-full h-2 w-2 bg-emerald-500"
                    ></span>
                </span>
                <span class="text-xs font-medium text-dark-600"
                    >Disponible maintenant</span
                >
            </div>

            <!-- Headline -->

            <!-- -->
            <h1
                class="text-4xl sm:text-5xl lg:text-6xl font-medium tracking-tight text-dark-900 leading-[1.1] mb-6"
            >
                Un expert à vos côtés,<br class="hidden lg:block" /> à vos horaires.
            </h1>

            <!-- Supporting Text -->
            <p
                class="text-lg text-dark-500 leading-relaxed mb-8 max-w-lg font-normal"
            >
                Réservez instantanément une consultation vérifiée. Fini les
                allers-retours par email, un pro à vos côtés quand vous en avez
                besoin.
            </p>

            <!-- Primary CTA -->
            <div class="w-full sm:w-auto flex flex-col sm:flex-row gap-4">
                <button
                    class="group w-full sm:w-auto inline-flex justify-center items-center gap-2 bg-dark-900 hover:bg-dark-800 text-white text-sm sm:text-base font-medium px-8 py-3.5 rounded-xl transition-all duration-200 shadow-sm hover:shadow-md ring-offset-2 focus:ring-2 focus:ring-dark-900"
                >
                    Prendre rendez-vous
                    <span
                        class="iconify group-hover:trandark-x-0.5 transition-transform"
                        data-icon="lucide:arrow-right"
                        data-width="18"
                        data-stroke-width="1.5"
                    ></span>
                </button>

                <button
                    class="w-full sm:w-auto inline-flex justify-center items-center gap-2 bg-white border border-dark-200 hover:border-dark-300 text-dark-600 hover:text-dark-900 text-sm sm:text-base font-medium px-8 py-3.5 rounded-xl transition-all duration-200"
                >
                    Comment ça fonctionne
                </button>
            </div>

            <!-- Social Proof / Trust Indicators -->
            <div class="mt-8 flex items-center gap-4 text-sm text-dark-500">
                <div class="flex -space-x-2">
                    <div
                        class="w-8 h-8 rounded-full border-2 border-white bg-dark-200 flex items-center justify-center text-xs overflow-hidden"
                    >
                        <img
                            src="https://ui-avatars.com/api/?name=Alex&background=random"
                            alt="User"
                        />
                    </div>
                    <div
                        class="w-8 h-8 rounded-full border-2 border-white bg-dark-200 flex items-center justify-center text-xs overflow-hidden"
                    >
                        <img
                            src="https://ui-avatars.com/api/?name=Sarah&background=random"
                            alt="User"
                        />
                    </div>
                    <div
                        class="w-8 h-8 rounded-full border-2 border-white bg-dark-200 flex items-center justify-center text-xs overflow-hidden"
                    >
                        <img
                            src="https://ui-avatars.com/api/?name=Mike&background=random"
                            alt="User"
                        />
                    </div>
                </div>
                <p class="font-medium">
                    Plus de 2 000 clients nous font confiance
                </p>
            </div>
        </div>

        <!-- Right Column: Visual / Context (Priority 2) -->
        <!-- On Mobile: Appears below CTA. On Desktop: Right column -->
        <div class="relative w-full mt-4 lg:mt-0 select-none">
            <!-- Decorative Background Blob -->
            <div
                class="absolute top-1/2 left-1/2 -trandark-x-1/2 -trandark-y-1/2 w-[120%] h-[120%] bg-gradient-to-tr from-dark-100 to-white rounded-full blur-3xl opacity-60 -z-10"
            ></div>

            <!-- UI Card Component: Booking Interface Representation -->
            <div
                class="relative bg-white rounded-2xl shadow-[0_8px_30px_rgb(0,0,0,0.04)] border border-dark-100 p-6 max-w-md mx-auto lg:mx-0 lg:ml-auto"
            >
                <!-- Card Header -->
                <div
                    class="flex items-center justify-between mb-6 border-b border-dark-50 pb-4"
                >
                    <div class="flex items-center gap-3">
                        <div
                            class="w-10 h-10 rounded-full bg-dark-100 flex items-center justify-center"
                        >
                            <span
                                class="iconify text-dark-600"
                                data-icon="lucide:user"
                                data-width="20"
                                data-stroke-width="1.5"
                            ></span>
                        </div>
                        <div>
                            <h3 class="text-sm font-medium text-dark-900">
                                Dr. Sarah L.
                            </h3>
                            <p class="text-xs text-dark-500">
                                Spécialiste en dermatologie
                            </p>
                        </div>
                    </div>
                    <div
                        class="px-2 py-1 rounded bg-green-50 text-green-700 text-xs font-medium"
                    >
                        Vérifié
                    </div>
                </div>

                <!-- Date Picker Strip -->
                <div class="flex justify-between items-center mb-6">
                    <button
                        class="text-dark-400 hover:text-dark-600"
                        aria-label="Mois précédent"
                    >
                        <span
                            class="iconify"
                            data-icon="lucide:chevron-left"
                            data-width="20"
                        ></span>
                    </button>
                    <span class="text-sm font-medium text-dark-900"
                        >Octobre 2023</span
                    >
                    <button
                        class="text-dark-400 hover:text-dark-600"
                        aria-label="Mois suivant"
                    >
                        <span
                            class="iconify"
                            data-icon="lucide:chevron-right"
                            data-width="20"
                        ></span>
                    </button>
                </div>

                <!-- Calendar Grid Visual -->
                <div class="grid grid-cols-4 gap-2 mb-6">
                    <div
                        class="flex flex-col items-center justify-center p-2 rounded-lg border border-dark-100 text-dark-400 bg-dark-50 cursor-not-allowed"
                    >
                        <span class="text-xs">Lun</span>
                        <span class="text-sm font-medium">12</span>
                    </div>
                    <div
                        class="flex flex-col items-center justify-center p-2 rounded-lg border border-dark-900 bg-dark-900 text-white shadow-md transform scale-105"
                    >
                        <span class="text-xs opacity-80">Mar</span>
                        <span class="text-sm font-medium">13</span>
                    </div>
                    <div
                        class="flex flex-col items-center justify-center p-2 rounded-lg border border-dark-200 text-dark-600 hover:border-dark-300 cursor-pointer"
                    >
                        <span class="text-xs">Mer</span>
                        <span class="text-sm font-medium">14</span>
                    </div>
                    <div
                        class="flex flex-col items-center justify-center p-2 rounded-lg border border-dark-200 text-dark-600 hover:border-dark-300 cursor-pointer"
                    >
                        <span class="text-xs">Jeu</span>
                        <span class="text-sm font-medium">15</span>
                    </div>
                </div>

                <!-- Time Slots -->
                <div class="space-y-2 mb-6">
                    <div
                        class="flex items-center justify-between p-3 rounded-lg border border-dark-100 hover:border-dark-300 transition-colors cursor-pointer group"
                    >
                        <div class="flex items-center gap-3">
                            <span
                                class="iconify text-dark-400 group-hover:text-dark-600"
                                data-icon="lucide:clock"
                                data-width="16"
                            ></span>
                            <span class="text-sm text-dark-700"
                                >09:00 - 09:45</span
                            >
                        </div>
                        <div
                            class="w-4 h-4 rounded-full border border-dark-300 group-hover:border-dark-500"
                        ></div>
                    </div>
                    <div
                        class="flex items-center justify-between p-3 rounded-lg bg-dark-50 border border-dark-200 cursor-pointer"
                    >
                        <div class="flex items-center gap-3">
                            <span
                                class="iconify text-dark-900"
                                data-icon="lucide:check-circle-2"
                                data-width="16"
                            ></span>
                            <span class="text-sm font-medium text-dark-900"
                                >14:00 - 14:45</span
                            >
                        </div>
                    </div>
                </div>

                <!-- Fake Confirm Button -->
                <div
                    class="w-full bg-dark-900 text-white py-3 rounded-lg text-sm font-medium text-center shadow-lg shadow-dark-200"
                >
                    Confirmer la réservation
                </div>
            </div>

            <!-- Floating Elements for visual depth -->
            <div
                class="absolute -right-4 top-20 bg-white p-3 rounded-lg shadow-xl border border-dark-100 flex items-center gap-3 animate-[bounce_3s_infinite]"
            >
                <div
                    class="w-8 h-8 rounded-full bg-green-100 flex items-center justify-center text-green-600"
                >
                    <span
                        class="iconify"
                        data-icon="lucide:video"
                        data-width="16"
                    ></span>
                </div>
                <div>
                    <p class="text-xs font-semibold text-dark-900">Visio</p>
                    <p class="text-[10px] text-dark-500">Prêt à rejoindre</p>
                </div>
            </div>
        </div>
    </div>
</div>

{#snippet content()}
    <div
        class="flex flex-col items-center justify-between gap-8 lg:gap-24 p-6 lg:p-24 mt-18 border border-blue-500"
    >
        <div
            class="lg:w-full w-fit grid gap-6 lg:gap-8 border border-red-400 max-w-2xl lg:max-w-none"
        >
            <h1
                class="text-3xl sm:text-4xl lg:text-6xl text-dark-900 font-bold"
            >
                Optimise ton corps <br class="lg:block" /> et apaise ton esprit
            </h1>
            <p
                class="text-foreground-alt font-normal leading-relaxed text-base lg:text-lg max-w-md lg:max-w-160"
            >
                Réserve tes séances instantanément : un accompagnement
                professionnel quand tu en as le plus besoin.
            </p>
            <Button class="text-white justify-self-start"
                >Prendre rendez vous</Button
            >
        </div>
        <div
            class="bg-dark-50 lg:w-full w-fit aspect-6/4 border border-green-400"
        ></div>
    </div>
{/snippet}
