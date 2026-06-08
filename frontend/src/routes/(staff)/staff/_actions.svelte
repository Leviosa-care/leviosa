<script lang="ts">
    import { TriangleAlert, ArrowRight } from "@lucide/svelte";
</script>

<div class="bg-background rounded-lg border border-border-card overflow-hidden">
    <div class="px-6 py-4 border-b border-border-card flex items-center gap-2">
        <TriangleAlert size={16} class="text-muted-foreground" />
        <h2 class="text-base font-semibold text-foreground">Actions requises</h2>
    </div>
    <div class="divide-y divide-border-card">
        {@render required_action(
            "Validation praticien",
            "2h",
            "Dr Thomas R. attend la validation de ses documents.",
            "Examiner",
        )}
        {@render required_action(
            "Profil incomplet",
            "hier",
            "Cabinet KineSport a des informations manquantes.",
            "Relancer",
        )}
        {@render required_action(
            "Litige client",
            "1j",
            "Signalement sur la consultation #4849",
            "Voir le détail",
        )}
    </div>
</div>

{#snippet required_action(
    title: string,
    duration: string,
    description: string,
    action: string,
)}
    <div class="px-6 py-4 hover:bg-muted/30 transition-colors">
        <div class="flex items-start justify-between gap-3 mb-1">
            <p class="text-sm font-medium text-foreground">{title}</p>
            <span class="text-xs text-muted-foreground whitespace-nowrap">{duration}</span>
        </div>
        <p class="text-xs text-muted-foreground mb-3">{description}</p>
        <button
            class="flex items-center gap-1.5 text-xs font-medium text-foreground-alt hover:text-foreground cursor-pointer transition-colors"
        >
            <span>{action}</span>
            <ArrowRight size={12} />
        </button>
    </div>
{/snippet}
