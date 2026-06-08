<script lang="ts">
</script>

<div class="bg-background rounded-2xl border border-border overflow-hidden">
    <div class="px-5 py-4">
        <h2 class="font-display text-sm font-semibold tracking-tight text-foreground">Activité récente</h2>
    </div>
    <ul>
        {@render activity("Nouvelle inscription", "par Paul V.", "10 min", "bg-blue-500")}
        {@render activity("Paiement reçu", "de 60,00 €", "32 min", "bg-green-500")}
        {@render activity("Mise à jour catalogue", "3 items", "2h", "bg-purple-500")}
    </ul>
</div>

{#snippet activity(title: string, subtitle: string, date: string, dotColor: string)}
    <li class="px-5 py-3.5 hover:bg-muted/20 transition-colors duration-300 flex items-start gap-3">
        <div class="w-1.5 h-1.5 rounded-full {dotColor} mt-2 flex-shrink-0"></div>
        <div class="flex-1 min-w-0">
            <div class="flex items-baseline gap-1.5 text-sm">
                <p class="font-medium text-foreground">{title}</p>
                <p class="text-muted-foreground">{subtitle}</p>
            </div>
        </div>
        <span class="text-[11px] text-muted-foreground tabular-nums whitespace-nowrap mt-0.5">il y a {date}</span>
    </li>
{/snippet}
