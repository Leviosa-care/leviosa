<script lang="ts">
    import { reveal } from "$lib/actions/reveal";
    import { Users, Hammer, ShieldCheck } from "@lucide/svelte";
</script>

<section id="why-us" class="py-24 bg-white">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <!-- Left-aligned Header -->
        <div class="mb-16" use:reveal={{ preset: "fade-up", delay: 100 }}>
            <h2
                class="text-4xl md:text-5xl font-semibold tracking-tight text-foreground"
            >
                Fondé sur la clarté, pas la complexité.
            </h2>
            <p
                class="mt-4 text-base md:text-lg text-muted-foreground leading-relaxed font-normal"
            >
                Nous avons éliminé le superflu des agences pour nous concentrer
                sur ce qui fait vraiment avancer votre entreprise.
            </p>
        </div>

        <!-- 3 Column Grid -->
        <div class="grid grid-cols-1 md:grid-cols-3 gap-6 lg:gap-8">
            <!-- Card 1: Direct Access -->
            <div
                class="group h-full p-8 rounded-2xl bg-surface border border-border-input hover:border-border-input-hover transition-colors duration-200 flex flex-col items-start text-left"
                use:reveal={{ preset: "fade-up", delay: 150 }}
            >
                <div
                    class="w-12 h-12 rounded-xl bg-white border border-border-input-hover flex items-center justify-center shadow-mini mb-6 group-hover:scale-105 transition-transform duration-200"
                >
                    <Users size={22} class="text-foreground" strokeWidth={1.5} />
                </div>
                <h3
                    class="text-lg font-semibold text-foreground tracking-tight mb-3"
                >
                    Accès direct aux experts
                </h3>
                <p class="text-sm leading-relaxed text-muted-foreground">
                    Vous échangez directement avec les seniors qui réalisent le travail.
                    Pas d'intermédiaires, pas de couches de traduction, juste une
                    collaboration efficace.
                </p>
            </div>

            <!-- Card 2: Pragmatism -->
            <div
                class="group h-full p-8 rounded-2xl bg-surface border border-border-input hover:border-border-input-hover transition-colors duration-200 flex flex-col items-start text-left"
                use:reveal={{ preset: "fade-up", delay: 200 }}
            >
                <div
                    class="w-12 h-12 rounded-xl bg-white border border-border-input-hover flex items-center justify-center shadow-mini mb-6 group-hover:scale-105 transition-transform duration-200"
                >
                    <Hammer size={22} class="text-foreground" strokeWidth={1.5} />
                </div>
                <h3
                    class="text-lg font-semibold text-foreground tracking-tight mb-3"
                >
                    Exécution pragmatique
                </h3>
                <p class="text-sm leading-relaxed text-muted-foreground">
                    Nous évitons les présentations de 50 pages. Notre priorité
                    est de livrer une valeur concrète et de résoudre le vrai problème
                    dès le premier jour.
                </p>
            </div>

            <!-- Card 3: Transparency -->
            <div
                class="group h-full p-8 rounded-2xl bg-surface border border-border-input hover:border-border-input-hover transition-colors duration-200 flex flex-col items-start text-left"
                use:reveal={{ preset: "fade-up", delay: 250 }}
            >
                <div
                    class="w-12 h-12 rounded-xl bg-white border border-border-input-hover flex items-center justify-center shadow-mini mb-6 group-hover:scale-105 transition-transform duration-200"
                >
                    <ShieldCheck size={22} class="text-foreground" strokeWidth={1.5} />
                </div>
                <h3
                    class="text-lg font-semibold text-foreground tracking-tight mb-3"
                >
                    Transparence totale
                </h3>
                <p class="text-sm leading-relaxed text-muted-foreground">
                    Délais clairs, tarification fixe et retours honnêtes. Si nous
                    pensons ne pas pouvoir vous aider, nous vous le dirons
                    immédiatement.
                </p>
            </div>
        </div>
    </div>
</section>
