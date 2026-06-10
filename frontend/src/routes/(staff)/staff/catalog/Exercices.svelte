<script lang="ts">
    import type { PageData } from "./$types";
    import { CalendarClock } from "@lucide/svelte";

    let { data }: { data: PageData } = $props();
</script>

<div class="flex flex-col items-center justify-center h-full min-h-[400px]">
    <div class="text-center max-w-md">
        <div
            class="w-20 h-20 bg-dark-10 rounded-full flex items-center justify-center mx-auto mb-6"
        >
            <CalendarClock class="text-foreground-alt size-10" />
        </div>
        <h2 class="text-2xl font-semibold tracking-tight mb-3">
            Bientôt disponible
        </h2>
        <p class="text-foreground-alt text-base mb-6">
            La gestion des exercices sera bientôt implémentée. Cette
            fonctionnalité vous permettra de créer et gérer des exercices pour
            vos produits.
        </p>
        <div
            class="bg-surface border border-border-card rounded-lg px-4 py-3 text-sm text-foreground-alt"
        >
            <p class="font-medium mb-1">Statut du développement</p>
            <p class="text-xs">
                Cette fonctionnalité est actuellement en cours de développement.
                Veuillez revenir plus tard pour les mises à jour.
            </p>
        </div>
    </div>
</div>
