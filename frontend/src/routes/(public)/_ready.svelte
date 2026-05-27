<script lang="ts">
    import { reveal } from "$lib/actions/reveal";
</script>

<!-- Final CTA Section (Conclusion) -->
<section class="py-24 bg-white relative overflow-hidden">
    <!-- Background Glow -->
    <div
        class="absolute top-0 left-1/2 -translate-x-1/2 w-full h-full max-w-3xl opacity-30 bg-[radial-gradient(ellipse_at_top,_var(--tw-gradient-stops))] from-indigo-100 via-white to-white pointer-events-none"
    ></div>

    <div class="relative max-w-3xl mx-auto px-6 text-center z-10">
        <!-- Headline -->
        <h2
            class="text-3xl sm:text-4xl md:text-5xl font-medium tracking-tight text-dark-900 mb-6"
            use:reveal={{ preset: "fade-up", delay: 100 }}
        >
            Prêt·e à prendre soin de votre santé ?
        </h2>

        <!-- Subtext -->
        <p class="text-lg text-dark-500 mb-10 max-w-lg mx-auto leading-relaxed" use:reveal={{ preset: "fade-up", delay: 150 }}>
            Rejoignez les centaines de patients qui ont déjà transformé leur
            quotidien. Accédez à une expertise vérifiée dès aujourd'hui.
        </p>

        <!-- Primary Action -->
        <div class="flex justify-center" use:reveal={{ preset: "fade-up", delay: 200 }}>
            <button
                class="bg-dark-900 text-white hover:bg-dark-800 px-8 py-4 rounded-xl font-semibold text-base transition-all duration-200 shadow-sm hover:shadow-md hover:scale-[1.02] active:scale-[0.98] cursor-pointer"
            >
                Prendre rendez-vous maintenant
            </button>
        </div>
    </div>
</section>
