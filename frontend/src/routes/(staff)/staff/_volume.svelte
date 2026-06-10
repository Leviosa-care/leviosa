<script lang="ts">
    import type { VolumeDay } from './+page.server';

    let { days }: { days: VolumeDay[] } = $props();

    // Highlight the day with the highest value (or today if all equal)
    const todayLabel = (() => {
        const d = new Date();
        const dayIdx = d.getDay(); // 0=Sun … 6=Sat
        const labelIdx = dayIdx === 0 ? 6 : dayIdx - 1; // Mon=0 … Sun=6
        return ['L', 'M', 'M', 'J', 'V', 'S', 'D'][labelIdx];
    })();

    const maxPct = $derived(Math.max(...days.map((d) => d.pct), 0));
    const highlightIdx = $derived(
        maxPct > 0
            ? days.reduce((best, d, i) => (d.pct > days[best].pct ? i : best), 0)
            : days.findIndex((d) => d.label === todayLabel)
    );
</script>

<div class="bg-background rounded-2xl border border-border overflow-hidden">
    <div class="px-5 py-4 flex items-center justify-between">
        <h2 class="font-display text-sm font-semibold tracking-tight text-foreground">Volume hebdomadaire</h2>
        <span class="text-[11px] text-muted-foreground">7 derniers jours</span>
    </div>
    <div class="px-5 pb-5">
        <div class="flex items-end gap-3 h-36">
            {#each days as day, i}
                <div class="flex-1 flex flex-col items-center gap-2">
                    <!-- Hover area with bar -->
                    <div class="w-full relative group cursor-default" style="height: 128px">
                        <!-- Background track -->
                        <div class="absolute inset-0 rounded-lg bg-muted/50"></div>
                        <!-- Filled bar -->
                        <div
                            class="absolute bottom-0 left-0 right-0 rounded-lg transition-all duration-700 ease-out {i === highlightIdx
                                ? 'bg-foreground'
                                : 'bg-foreground/15 group-hover:bg-foreground/25'}"
                            style="height: {day.pct}%"
                        ></div>
                        <!-- Hover value tooltip -->
                        <div class="absolute -top-6 left-1/2 -translate-x-1/2 text-[10px] font-semibold text-foreground tabular-nums opacity-0 group-hover:opacity-100 transition-opacity duration-200">
                            {day.pct}%
                        </div>
                    </div>
                    <span class="text-[10px] font-medium text-muted-foreground tracking-wide {i === highlightIdx ? '!text-foreground' : ''}">
                        {day.label}
                    </span>
                </div>
            {/each}
        </div>
    </div>
</div>
