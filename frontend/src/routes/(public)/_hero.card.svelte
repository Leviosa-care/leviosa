<script lang="ts">
    import { reveal } from "$lib/actions/reveal";
</script>

<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-12 lg:gap-30 items-center">
        <!-- Left Column: Content (Priority 1) -->
        <div class="flex flex-col items-start text-left z-10">
            <!-- Trust Badge -->
            <div
                class="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-surface border border-border-input-hover mb-6"
                use:reveal={{ preset: "fade-down", delay: 100 }}
            >
                <span class="relative flex h-2 w-2">
                    <span
                        class="animate-ping absolute inline-flex h-full w-full rounded-full bg-emerald-400 opacity-75"
                    ></span>
                    <span
                        class="relative inline-flex rounded-full h-2 w-2 bg-emerald-500"
                    ></span>
                </span>
                <span class="text-xs font-medium text-foreground-alt"
                    >Disponible maintenant</span
                >
            </div>

            <!-- Headline -->

            <!-- -->
            <h1
                class="text-4xl sm:text-5xl lg:text-6xl font-medium tracking-tight text-foreground leading-[1.1] mb-6"
                use:reveal={{ preset: "fade-up", delay: 150 }}
            >
                Un expert à vos côtés,<br class="hidden lg:block" /> à vos horaires.
            </h1>

            <!-- Supporting Text -->
            <p
                class="text-lg text-muted-foreground leading-relaxed mb-8 max-w-lg font-normal"
                use:reveal={{ preset: "fade-up", delay: 200 }}
            >
                Réservez instantanément une consultation vérifiée. Fini les
                allers-retours par email, un pro à vos côtés quand vous en avez
                besoin.
            </p>

            <!-- Primary CTA -->
            <div class="w-full sm:w-auto flex flex-col sm:flex-row gap-4" use:reveal={{ preset: "fade-up", delay: 250 }}>
                <button
                    class="group w-full sm:w-auto inline-flex justify-center items-center gap-2 bg-foreground hover:bg-foreground-alt text-white text-sm sm:text-base font-medium px-8 py-3.5 rounded-xl transition-all duration-200 shadow-mini hover:shadow-card ring-offset-2 focus:ring-2 focus:ring-foreground cursor-pointer"
                >
                    Prendre rendez-vous
                    <span
                        class="iconify group-hover:translate-x-0.5 transition-transform"
                        data-icon="lucide:arrow-right"
                        data-width="18"
                        data-stroke-width="1.5"
                    ></span>
                </button>

                <button
                    class="w-full sm:w-auto inline-flex justify-center items-center gap-2 bg-white border border-border-input-hover hover:border-border-input-hover text-foreground-alt hover:text-foreground text-sm sm:text-base font-medium px-8 py-3.5 rounded-xl transition-all duration-200 cursor-pointer"
                >
                    Comment ça fonctionne
                </button>
            </div>

            <!-- Social Proof / Trust Indicators -->
            <div class="mt-8 flex items-center gap-4 text-sm text-muted-foreground" use:reveal={{ preset: "fade-up", delay: 300 }}>
                <div class="flex -space-x-2">
                    <div
                        class="w-8 h-8 rounded-full border-2 border-white bg-border-input-hover flex items-center justify-center text-xs font-medium text-foreground-alt"
                    >
                        A
                    </div>
                    <div
                        class="w-8 h-8 rounded-full border-2 border-white bg-border-input-hover flex items-center justify-center text-xs font-medium text-foreground-alt"
                    >
                        S
                    </div>
                    <div
                        class="w-8 h-8 rounded-full border-2 border-white bg-border-input-hover flex items-center justify-center text-xs font-medium text-foreground-alt"
                    >
                        M
                    </div>
                </div>
                <p class="font-medium">
                    Plus de 2 000 clients nous font confiance
                </p>
            </div>
        </div>

        <!-- Right Column: Visual / Context (Priority 2) -->
        <!-- On Mobile: Appears below CTA. On Desktop: Right column -->
        <div class="relative w-full mt-4 lg:mt-0 select-none" use:reveal={{ preset: "fade-in-scale", delay: 200 }}>
            <!-- Decorative Background Blob -->
            <div
                class="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-[120%] h-[120%] bg-gradient-to-tr from-surface-hover to-white rounded-full blur-3xl opacity-60 -z-10"
            ></div>

            <!-- UI Card Component: Booking Interface Representation -->
            <div
                class="relative bg-white rounded-2xl shadow-[0_8px_30px_rgb(0,0,0,0.04)] border border-border-input p-6 max-w-md mx-auto lg:mx-0 lg:ml-auto"
            >
                <!-- Card Header -->
                <div
                    class="flex items-center justify-between mb-6 border-b border-border-input pb-4"
                >
                    <div class="flex items-center gap-3">
                        <div
                            class="w-10 h-10 rounded-full bg-surface-hover flex items-center justify-center"
                        >
                            <span
                                class="iconify text-foreground-alt"
                                data-icon="lucide:user"
                                data-width="20"
                                data-stroke-width="1.5"
                            ></span>
                        </div>
                        <div>
                            <h3 class="text-sm font-medium text-foreground">
                                Dr. Sarah L.
                            </h3>
                            <p class="text-xs text-muted-foreground">
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
                        class="text-muted-foreground hover:text-foreground-alt cursor-pointer"
                        aria-label="Mois précédent"
                    >
                        <span
                            class="iconify"
                            data-icon="lucide:chevron-left"
                            data-width="20"
                        ></span>
                    </button>
                    <span class="text-sm font-medium text-foreground"
                        >Octobre 2023</span
                    >
                    <button
                        class="text-muted-foreground hover:text-foreground-alt cursor-pointer"
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
                        class="flex flex-col items-center justify-center p-2 rounded-lg border border-border-input text-muted-foreground bg-surface cursor-not-allowed"
                    >
                        <span class="text-xs">Lun</span>
                        <span class="text-sm font-medium">12</span>
                    </div>
                    <div
                        class="flex flex-col items-center justify-center p-2 rounded-lg border border-foreground bg-foreground text-white shadow-card transform scale-105"
                    >
                        <span class="text-xs opacity-80">Mar</span>
                        <span class="text-sm font-medium">13</span>
                    </div>
                    <div
                        class="flex flex-col items-center justify-center p-2 rounded-lg border border-border-input-hover text-foreground-alt hover:border-border-input-hover cursor-pointer"
                    >
                        <span class="text-xs">Mer</span>
                        <span class="text-sm font-medium">14</span>
                    </div>
                    <div
                        class="flex flex-col items-center justify-center p-2 rounded-lg border border-border-input-hover text-foreground-alt hover:border-border-input-hover cursor-pointer"
                    >
                        <span class="text-xs">Jeu</span>
                        <span class="text-sm font-medium">15</span>
                    </div>
                </div>

                <!-- Time Slots -->
                <div class="space-y-2 mb-6">
                    <div
                        class="flex items-center justify-between p-3 rounded-lg border border-border-input hover:border-border-input-hover transition-colors cursor-pointer group"
                    >
                        <div class="flex items-center gap-3">
                            <span
                                class="iconify text-muted-foreground group-hover:text-foreground-alt"
                                data-icon="lucide:clock"
                                data-width="16"
                            ></span>
                            <span class="text-sm text-foreground-alt"
                                >09:00 - 09:45</span
                            >
                        </div>
                        <div
                            class="w-4 h-4 rounded-full border border-border-input-hover group-hover:border-muted-foreground"
                        ></div>
                    </div>
                    <div
                        class="flex items-center justify-between p-3 rounded-lg bg-surface border border-border-input-hover cursor-pointer"
                    >
                        <div class="flex items-center gap-3">
                            <span
                                class="iconify text-foreground"
                                data-icon="lucide:check-circle-2"
                                data-width="16"
                            ></span>
                            <span class="text-sm font-medium text-foreground"
                                >14:00 - 14:45</span
                            >
                        </div>
                    </div>
                </div>

                <!-- Fake Confirm Button -->
                <div
                    class="w-full bg-foreground text-white py-3 rounded-lg text-sm font-medium text-center shadow-card shadow-border-input-hover"
                >
                    Confirmer la réservation
                </div>
            </div>

            <!-- Floating Elements for visual depth -->
            <div
                class="absolute -right-4 top-20 bg-white p-3 rounded-lg shadow-popover border border-border-input flex items-center gap-3 animate-[bounce_3s_infinite]"
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
                    <p class="text-xs font-semibold text-foreground">Visio</p>
                    <p class="text-[10px] text-muted-foreground">Prêt à rejoindre</p>
                </div>
            </div>
        </div>
    </div>
</div>


