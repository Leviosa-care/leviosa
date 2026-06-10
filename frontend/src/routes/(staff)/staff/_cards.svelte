<script lang="ts">
    import type { KpiData } from './+page.server';

    import {
        TrendingUp,
        CalendarCheck,
        Clock,
        BarChart3,
    } from "@lucide/svelte";

    let { kpi }: { kpi: KpiData } = $props();

    const revenueValue = $derived(
        kpi.revenueCents > 0
            ? (kpi.revenueCents / 100).toLocaleString('fr-FR', { minimumFractionDigits: 0, maximumFractionDigits: 0 })
            : '0'
    );
</script>

<div class="grid grid-cols-2 lg:grid-cols-4 gap-3 lg:gap-4">
    {@render card(TrendingUp, "Revenus", kpi.revenueGrowthPct, `${revenueValue} €`, "bg-green-100 dark:bg-green-900/30 text-green-600 dark:text-green-400", "7 derniers jours")}
    {@render card(CalendarCheck, "Réservations", kpi.bookingsGrowthPct, String(kpi.bookingsCount), "bg-blue-100 dark:bg-blue-900/30 text-blue-600 dark:text-blue-400", "7 derniers jours")}
    {@render card(Clock, "Occupation", "—", kpi.occupationPct, "bg-amber-100 dark:bg-amber-900/30 text-amber-600 dark:text-amber-400", "Cette semaine")}
    {@render card(BarChart3, "Satisfaction", "—", "—", "bg-purple-100 dark:bg-purple-900/30 text-purple-600 dark:text-purple-400", "Sur 0 avis")}
</div>

{#snippet card(Icon: typeof import('@lucide/svelte').TrendingUp, title: string, trend: string, value: string, iconColor: string, subtitle: string)}
    <div class="group relative bg-background rounded-2xl border border-border p-5 lg:p-6 hover:shadow-card hover:shadow-foreground/[0.03] hover:-translate-y-0.5 transition-all duration-500">
        <div class="flex items-start justify-between mb-5">
            <div class="w-10 h-10 rounded-xl {iconColor} flex items-center justify-center group-hover:scale-110 transition-transform duration-500">
                <Icon size={18} strokeWidth={2} />
            </div>
            <span class="text-[11px] font-semibold tracking-wider text-muted-foreground tabular-nums mt-1">{trend}</span>
        </div>
        <p class="text-3xl font-bold text-foreground tracking-tight tabular-nums leading-none">{value}</p>
        <div class="mt-3 flex items-center justify-between">
            <p class="text-sm font-medium text-foreground">{title}</p>
        </div>
        <p class="text-[11px] text-muted-foreground mt-0.5">{subtitle}</p>
    </div>
{/snippet}
