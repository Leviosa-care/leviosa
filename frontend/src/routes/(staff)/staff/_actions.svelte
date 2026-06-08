<script lang="ts">
    import { TriangleAlert, ArrowRight } from "@lucide/svelte";
</script>

<div class="bg-background rounded-xl border border-border overflow-hidden">
    <div class="px-5 py-4 border-b border-border flex items-center gap-2">
        <TriangleAlert size={14} class="text-amber-500" />
        <h2 class="text-sm font-semibold tracking-tight text-foreground uppercase">Actions requises</h2>
    </div>
    <div class="divide-y divide-border">
        {@render required_action("Validation praticien", "2h", "Dr Thomas R. attend la validation de ses documents.", "Examiner")}
        {@render required_action("Profil incomplet", "hier", "Cabinet KineSport a des informations manquantes.", "Relancer")}
        {@render required_action("Litige client", "1j", "Signalement sur la consultation #4849", "Voir le détail")}
    </div>
</div>

{#snippet required_action(title: string, duration: string, description: string, action: string)}
    <div class="px-5 py-4 hover:bg-muted/20 transition-colors duration-200">
        <div class="flex items-center justify-between mb-1.5">
            <p class="text-sm font-medium text-foreground">{title}</p>
            <span class="text-[11px] font-medium text-muted-foreground tabular-nums">{duration}</span>
        </div>
        <p class="text-xs text-muted-foreground leading-relaxed mb-3">{description}</p>
        <button
            class="inline-flex items-center gap-1.5 text-xs font-semibold text-foreground hover:text-foreground/70 cursor-pointer transition-colors"
        >
            {action}
            <ArrowRight size={11} />
        </button>
    </div>
{/snippet}
