<script lang="ts">
    import { ArrowRight } from "@lucide/svelte";
</script>

<div class="bg-background rounded-2xl border border-border overflow-hidden">
    <div class="px-5 py-4 flex items-center gap-3">
        <div class="w-1.5 h-1.5 rounded-full bg-amber-500 animate-pulse"></div>
        <h2 class="font-display text-sm font-semibold tracking-tight text-foreground">Actions requises</h2>
        <span class="ml-auto text-[11px] font-semibold text-amber-600 dark:text-amber-400 tabular-nums">3</span>
    </div>
    <div class="divide-y divide-border">
        {@render required_action("Validation praticien", "2h", "Dr Thomas R. attend la validation de ses documents.", "Examiner")}
        {@render required_action("Profil incomplet", "hier", "Cabinet KineSport a des informations manquantes.", "Relancer")}
        {@render required_action("Litige client", "1j", "Signalement sur la consultation #4849", "Voir le détail")}
    </div>
</div>

{#snippet required_action(title: string, duration: string, description: string, action: string)}
    <div class="px-5 py-4 group hover:bg-muted/20 transition-colors duration-300">
        <div class="flex items-start justify-between gap-3">
            <p class="text-sm font-medium text-foreground leading-snug">{title}</p>
            <span class="text-[11px] text-muted-foreground tabular-nums whitespace-nowrap mt-0.5">{duration}</span>
        </div>
        <p class="text-[11px] text-muted-foreground leading-relaxed mt-1 line-clamp-2">{description}</p>
        <button
            class="inline-flex items-center gap-1 mt-3 text-[11px] font-semibold text-foreground/70 hover:text-foreground cursor-pointer transition-colors group-hover:translate-x-0.5 duration-300"
        >
            {action}
            <ArrowRight size={10} class="transition-transform group-hover:translate-x-0.5 duration-300" />
        </button>
    </div>
{/snippet}
