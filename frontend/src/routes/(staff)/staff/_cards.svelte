<script lang="ts">
    import { TrendingUp } from "@lucide/svelte";
</script>

<div class="grid grid-cols-2 lg:grid-cols-4 gap-4 lg:gap-6">
    {@render card("Total Revenue", "+12.5%", "$1,250.00", "Trending up this month", "Visitors for the last 6 months")}
    {@render card("Réservations", "+8.2%", "48", "En hausse ce mois", "Réservations mensuelles")}
    {@render card("Taux d'occupation", "94%", "94%", "Performance élevée", "Sur les 7 derniers jours")}
    {@render card("Satisfaction", "4.9/5", "4.9", "Note moyenne clients", "Basée sur 128 avis")}
</div>

{#snippet card(title: string, trend: string, value: string, subtitle: string, description: string)}
    <div class="bg-background rounded-lg border border-border-card p-5">
        <div class="flex items-center justify-between mb-3">
            <span class="text-xs font-medium text-muted-foreground uppercase tracking-wider">{title}</span>
            <div
                class="text-xs font-medium flex items-center gap-1 border border-border-input hover:border-border-input-hover px-2 py-0.5 rounded-full"
            >
                <TrendingUp size={12} />
                <span>{trend}</span>
            </div>
        </div>
        <p class="text-2xl font-bold text-foreground mb-1">{value}</p>
        <p class="text-sm text-muted-foreground">{subtitle}</p>
    </div>
{/snippet}
