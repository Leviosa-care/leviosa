<script lang="ts">
</script>

<div class="bg-background rounded-xl border border-border overflow-hidden">
    <div class="px-5 py-4 border-b border-border">
        <h2 class="text-sm font-semibold tracking-tight text-foreground uppercase">Activité récente</h2>
    </div>
    <ul class="divide-y divide-border">
        {@render activity("Nouvelle inscription", "par Paul V.", "Il y a 10 min")}
        {@render activity("Paiement reçu", "de 60,00 €", "Il y a 32 min")}
        {@render activity("Mise à jour catalogue", "(3 items)", "Il y a 2h")}
    </ul>
</div>

{#snippet activity(title: string, subtitle: string, date: string)}
    <li class="px-5 py-3.5 hover:bg-muted/20 transition-colors duration-200">
        <div class="flex items-baseline gap-1.5 text-sm">
            <p class="font-medium text-foreground">{title}</p>
            <p class="text-muted-foreground">{subtitle}</p>
        </div>
        <p class="text-[11px] text-muted-foreground mt-0.5 tabular-nums">{date}</p>
    </li>
{/snippet}
