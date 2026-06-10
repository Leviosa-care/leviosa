<script lang="ts">
    import type { ActionItem } from "./+page.server";
    import { ArrowRight } from "@lucide/svelte";

    let { actions }: { actions: ActionItem[] } = $props();
</script>

{#if actions.length > 0}
<div class="bg-background rounded-2xl border border-border overflow-hidden">
    <div class="px-5 py-4 flex items-center gap-3">
        <div class="w-1.5 h-1.5 rounded-full bg-amber-500 animate-pulse"></div>
        <h2 class="font-display text-sm font-semibold tracking-tight text-foreground">Actions requises</h2>
        <span class="ml-auto text-[11px] font-semibold text-amber-600 dark:text-amber-400 tabular-nums">{actions.length}</span>
    </div>
    <div class="divide-y divide-border">
        {#each actions as action}
            <a
                href={action.href}
                class="block px-5 py-4 group hover:bg-muted/20 transition-colors duration-300"
            >
                <div class="flex items-start justify-between gap-3">
                    <p class="text-sm font-medium text-foreground leading-snug">{action.title}</p>
                    <span class="text-[11px] text-muted-foreground tabular-nums whitespace-nowrap mt-0.5">{action.age}</span>
                </div>
                <p class="text-[11px] text-muted-foreground leading-relaxed mt-1 line-clamp-2">{action.description}</p>
                <span
                    class="inline-flex items-center gap-1 mt-3 text-[11px] font-semibold text-foreground/70 group-hover:text-foreground transition-colors group-hover:translate-x-0.5 duration-300"
                >
                    Voir
                    <ArrowRight size={10} class="transition-transform group-hover:translate-x-0.5 duration-300" />
                </span>
            </a>
        {/each}
    </div>
</div>
{/if}
