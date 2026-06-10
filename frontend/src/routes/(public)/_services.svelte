<script lang="ts">
    import { reveal } from "$lib/actions/reveal";
    import { Zap, Check, Search, Lightbulb, TrendingUp } from "@lucide/svelte";
</script>

<section id="services" class="py-20 lg:py-24 bg-white">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <!-- Section Header -->
        <div
            class="flex flex-col md:flex-row md:items-end justify-between gap-6 mb-12"
            use:reveal={{ preset: "fade-up", delay: 100 }}
        >
            <div class="max-w-xl">
                <span
                    class="text-sm font-semibold text-muted-foreground uppercase tracking-wider"
                    >Notre Expertise</span
                >
                <h2
                    class="mt-3 text-3xl md:text-4xl font-semibold tracking-tight text-foreground"
                >
                    Conçu sur mesure pour vos besoins
                </h2>
                <p class="mt-4 text-base text-muted-foreground">
                    Nous proposons actuellement un service spécialisé conçu pour
                    vous apporter une valeur immédiate et des actions concrètes.
                </p>
            </div>
        </div>

        <!-- Service Card: Single Item Layout -->
        <div
            class="group relative bg-white rounded-3xl border border-border-input-hover shadow-mini overflow-hidden hover:shadow-card transition-all duration-300"
            use:reveal={{ preset: "fade-up", delay: 150 }}
        >
            <div class="grid grid-cols-1 lg:grid-cols-2">
                <!-- Content Side -->
                <div
                    class="p-8 lg:p-12 flex flex-col justify-center order-2 lg:order-1"
                >
                    <div class="flex items-center gap-3 mb-6">
                        <div
                            class="w-10 h-10 rounded-lg bg-tertiary/10 flex items-center justify-center text-foreground"
                        >
                            <Zap size={20} strokeWidth={1.5} />
                        </div>
                        <span
                            class="text-sm font-semibold text-foreground tracking-wide uppercase"
                            >Le Plus Populaire</span
                        >
                    </div>

                    <h3
                        class="text-2xl md:text-3xl font-semibold text-foreground tracking-tight mb-4"
                    >
                        Consultation Stratégique
                    </h3>

                    <p
                        class="text-muted-foreground text-base md:text-lg leading-relaxed mb-8"
                    >
                        Une séance de 60 minutes pour analyser vos défis actuels
                        et élaborer une feuille de route concrète, étape par étape.
                        Idéal pour les fondateurs en phase de démarrage.
                    </p>

                    <div class="space-y-4 mb-10">
                        <!-- Benefit 1 -->
                        <div class="flex gap-3 items-start">
                            <div
                                class="mt-1 flex-shrink-0 w-5 h-5 rounded-full bg-surface-hover flex items-center justify-center"
                            >
                                <Check size={12} class="text-foreground" strokeWidth={2} />
                            </div>
                            <div>
                                <p class="text-sm font-medium text-foreground">
                                    Audit Approfondi
                                </p>
                                <p class="text-xs text-muted-foreground mt-0.5">
                                    Analyse complète de votre stack technique et
                                    de vos indicateurs.
                                </p>
                            </div>
                        </div>
                        <!-- Benefit 2 -->
                        <div class="flex gap-3 items-start">
                            <div
                                class="mt-1 flex-shrink-0 w-5 h-5 rounded-full bg-surface-hover flex items-center justify-center"
                            >
                                <Check size={12} class="text-foreground" strokeWidth={2} />
                            </div>
                            <div>
                                <p class="text-sm font-medium text-foreground">
                                    Feuille de Route Actionnable
                                </p>
                                <p class="text-xs text-muted-foreground mt-0.5">
                                    Repartez avec un plan clair de 30-60-90 jours.
                                </p>
                            </div>
                        </div>
                        <!-- Benefit 3 -->
                        <div class="flex gap-3 items-start">
                            <div
                                class="mt-1 flex-shrink-0 w-5 h-5 rounded-full bg-surface-hover flex items-center justify-center"
                            >
                                <Check size={12} class="text-foreground" strokeWidth={2} />
                            </div>
                            <div>
                                <p class="text-sm font-medium text-foreground">
                                    Enregistrement &amp; Notes
                                </p>
                                <p class="text-xs text-muted-foreground mt-0.5">
                                    Transcription complète et liens vers les ressources fournis.
                                </p>
                            </div>
                        </div>
                    </div>

                    <!-- Secondary CTA (Discovery) -->
                    <div class="flex flex-col sm:flex-row gap-4">
                        <a
                            href="/services"
                            class="inline-flex justify-center items-center gap-2 bg-white text-foreground border border-border-input-hover hover:border-border-input-hover hover:bg-surface text-sm font-medium px-6 py-3 rounded-xl transition-all duration-200 cursor-pointer"
                        >
                            Découvrir le service
                        </a>
                    </div>
                </div>

                <!-- Visual Side -->
                <div
                    class="relative bg-surface order-1 lg:order-2 border-b lg:border-b-0 lg:border-l border-border-input min-h-[300px] lg:min-h-full overflow-hidden flex items-center justify-center p-8"
                >
                    <!-- Decorative Pattern -->
                    <div
                        class="absolute inset-0 opacity-[0.03]"
                        style="background-image: radial-gradient(#0f172a 1px, transparent 1px); background-size: 20px 20px;"
                    ></div>

                    <!-- UI Representation: Roadmap Visualization -->
                    <div class="relative w-full max-w-sm">
                        <div
                            class="absolute top-0 bottom-0 left-6 w-0.5 bg-border-input-hover"
                        ></div>

                        <!-- Item 1 -->
                        <div class="relative flex items-center gap-4 mb-6">
                            <div
                                class="w-12 h-12 rounded-full border-4 border-border-input bg-white shadow-mini flex items-center justify-center z-10"
                            >
                                <Search size={18} class="text-muted-foreground" />
                            </div>
                            <div
                                class="flex-1 bg-white p-3 rounded-lg border border-border-input shadow-mini opacity-60"
                            >
                                <div
                                    class="h-2 w-16 bg-surface-hover rounded mb-2"
                                ></div>
                                <div
                                    class="h-1.5 w-full bg-surface rounded"
                                ></div>
                            </div>
                        </div>

                        <!-- Item 2 (Active) -->
                        <div class="relative flex items-center gap-4 mb-6">
                            <div
                                class="w-12 h-12 rounded-full border-4 border-border-input bg-foreground shadow-card shadow-border-input-hover flex items-center justify-center z-10"
                            >
                                <Lightbulb size={18} class="text-white" />
                            </div>
                            <div
                                class="flex-1 bg-white p-4 rounded-xl border border-border-input-hover shadow-card"
                            >
                                <div
                                    class="flex justify-between items-start mb-2"
                                >
                                    <div
                                        class="h-2.5 w-24 bg-foreground rounded"
                                    ></div>
                                    <div
                                        class="h-4 w-4 rounded-full bg-emerald-100 border border-emerald-200 flex items-center justify-center"
                                    >
                                        <Check size={10} class="text-emerald-600" strokeWidth={2} />
                                    </div>
                                </div>
                                <div class="space-y-1.5">
                                    <div
                                        class="h-2 w-full bg-surface-hover rounded"
                                    ></div>
                                    <div
                                        class="h-2 w-2/3 bg-surface-hover rounded"
                                    ></div>
                                </div>
                            </div>
                        </div>

                        <!-- Item 3 -->
                        <div class="relative flex items-center gap-4">
                            <div
                                class="w-12 h-12 rounded-full border-4 border-border-input bg-white shadow-mini flex items-center justify-center z-10"
                            >
                                <TrendingUp size={18} class="text-muted-foreground" />
                            </div>
                            <div
                                class="flex-1 bg-white p-3 rounded-lg border border-border-input shadow-mini opacity-60"
                            >
                                <div
                                    class="h-2 w-20 bg-surface-hover rounded mb-2"
                                ></div>
                                <div
                                    class="h-1.5 w-full bg-surface rounded"
                                ></div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</section>
