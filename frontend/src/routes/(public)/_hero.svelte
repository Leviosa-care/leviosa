<script lang="ts">
    import { reveal } from "$lib/actions/reveal";
    import { ArrowRight } from "@lucide/svelte";
</script>

<div class="relative min-h-screen flex items-center overflow-hidden">
    <!-- Background Video -->
    <video
        class="absolute inset-0 w-full h-full object-cover"
        src="/hero.mp4"
        autoplay
        muted
        loop
        playsinline
    ></video>

    <!-- Dark overlay for text legibility -->
    <div class="absolute inset-0 bg-black/50"></div>

    <!-- Content -->
    <div class="relative z-10 w-full max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="max-w-2xl flex flex-col items-start text-left">
            <!-- Trust Badge -->
            <div
                class="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-white/10 border border-white/20 backdrop-blur-sm mb-6"
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
                <span class="text-xs font-medium text-white/90">Disponible maintenant</span>
            </div>

            <!-- Headline -->
            <h1
                class="text-4xl sm:text-5xl lg:text-6xl font-display font-medium tracking-tight text-white leading-[1.1] mb-6"
                use:reveal={{ preset: "fade-up", delay: 150 }}
            >
                Un expert à vos côtés,<br class="hidden lg:block" /> à vos horaires.
            </h1>

            <!-- Supporting Text -->
            <p
                class="text-lg text-white/75 leading-relaxed mb-8 max-w-lg font-normal"
                use:reveal={{ preset: "fade-up", delay: 200 }}
            >
                Réservez instantanément une consultation vérifiée. Fini les
                allers-retours par email, un pro à vos côtés quand vous en avez
                besoin.
            </p>

            <!-- Primary CTA -->
            <div class="w-full sm:w-auto flex flex-col sm:flex-row gap-4" use:reveal={{ preset: "fade-up", delay: 250 }}>
                <a
                    href="/book"
                    class="group w-full sm:w-auto inline-flex justify-center items-center gap-2 bg-tertiary hover:bg-tertiary/90 text-dark-900 text-sm sm:text-base font-medium px-8 py-3.5 rounded-xl transition-all duration-200 shadow-sm hover:shadow-md cursor-pointer"
                >
                    Prendre rendez-vous
                    <ArrowRight size={18} class="group-hover:translate-x-0.5 transition-transform" strokeWidth={1.5} />
                </a>

                <a
                    href="#how"
                    class="w-full sm:w-auto inline-flex justify-center items-center gap-2 bg-white/10 backdrop-blur-sm border border-white/25 hover:bg-white/20 text-white text-sm sm:text-base font-medium px-8 py-3.5 rounded-xl transition-all duration-200"
                >
                    Comment ça fonctionne
                </a>
            </div>

            <!-- Social Proof -->
            <div class="mt-8 flex items-center gap-4 text-sm text-white/70" use:reveal={{ preset: "fade-up", delay: 300 }}>
                <div class="flex -space-x-2">
                    <img
                        src="https://images.unsplash.com/photo-1487412720507-e7ab37603c6f?w=80&h=80&auto=format&fit=crop&q=80"
                        alt="Emma W."
                        class="w-8 h-8 rounded-full border-2 border-white/30 object-cover"
                    />
                    <img
                        src="https://images.unsplash.com/photo-1500648767791-00dcc994a43e?w=80&h=80&auto=format&fit=crop&q=80"
                        alt="Thomas R."
                        class="w-8 h-8 rounded-full border-2 border-white/30 object-cover"
                    />
                    <img
                        src="https://images.unsplash.com/photo-1534528741775-53994a69daeb?w=80&h=80&auto=format&fit=crop&q=80"
                        alt="Camille D."
                        class="w-8 h-8 rounded-full border-2 border-white/30 object-cover"
                    />
                </div>
                <p class="font-medium">Plus de 2 000 clients nous font confiance</p>
            </div>
        </div>
    </div>
</div>
