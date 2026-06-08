<script lang="ts">
    import {
        TrendingUp,
        CalendarCheck,
        Clock,
        BarChart3,
    } from "@lucide/svelte";
</script>

<div class="grid grid-cols-2 lg:grid-cols-4 gap-4 lg:gap-5">
    {@render card(TrendingUp, "Revenus (7j)", "+12,5%", "1 250,00 €", "bg-green-100 dark:bg-green-900/30 text-green-600 dark:text-green-400")}
    {@render card(CalendarCheck, "Réservations", "+8,2%", "48", "bg-blue-100 dark:bg-blue-900/30 text-blue-600 dark:text-blue-400")}
    {@render card(Clock, "Taux d'occupation", "94%", "94%", "bg-amber-100 dark:bg-amber-900/30 text-amber-600 dark:text-amber-400")}
    {@render card(BarChart3, "Satisfaction", "4.9/5", "4,9", "bg-purple-100 dark:bg-purple-900/30 text-purple-600 dark:text-purple-400")}
</div>

{#snippet card(Icon: typeof import('@lucide/svelte').TrendingUp, title: string, trend: string, value: string, iconColor: string)}
    <div class="bg-background rounded-xl border border-border p-5 group hover:border-foreground/10 transition-colors duration-300">
        <div class="flex items-center justify-between mb-4">
            <div class="w-9 h-9 rounded-lg {iconColor} flex items-center justify-center">
                <Icon size={18} strokeWidth={2} />
            </div>
            <span class="text-xs font-medium tracking-wide text-muted-foreground">{trend}</span>
        </div>
        <p class="text-2xl font-bold text-foreground tracking-tight tabular-nums">{value}</p>
        <p class="text-sm text-muted-foreground mt-1">{title}</p>
    </div>
{/snippet}
