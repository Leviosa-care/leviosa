<script lang="ts">
    import type { ActivityEvent } from './+page.server';

    let { events }: { events: ActivityEvent[] } = $props();

    const colorMap: Record<string, string> = {
        blue: 'bg-blue-500',
        green: 'bg-green-500',
        purple: 'bg-purple-500',
        amber: 'bg-amber-500',
        red: 'bg-red-500',
    };
</script>

<div class="bg-background rounded-2xl border border-border overflow-hidden">
    <div class="px-5 py-4">
        <h2 class="font-display text-sm font-semibold tracking-tight text-foreground">Activité récente</h2>
    </div>

    {#if events.length === 0}
        <div class="px-5 py-10 text-center">
            <p class="text-sm text-muted-foreground">Aucune activité récente</p>
        </div>
    {:else}
        <ul>
            {#each events as event}
                <li class="px-5 py-3.5 hover:bg-muted/20 transition-colors duration-300 flex items-start gap-3">
                    <div class="w-1.5 h-1.5 rounded-full {colorMap[event.colorKey] ?? 'bg-muted-foreground'} mt-2 flex-shrink-0"></div>
                    <div class="flex-1 min-w-0">
                        <div class="flex items-baseline gap-1.5 text-sm">
                            <p class="font-medium text-foreground">{event.title}</p>
                            <p class="text-muted-foreground">{event.subtitle}</p>
                        </div>
                    </div>
                    <span class="text-[11px] text-muted-foreground tabular-nums whitespace-nowrap mt-0.5">{event.relativeTime}</span>
                </li>
            {/each}
        </ul>
    {/if}
</div>
