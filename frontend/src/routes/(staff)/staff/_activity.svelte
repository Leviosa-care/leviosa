<script lang="ts">
</script>

<div class="bg-background rounded-lg border border-border-card overflow-hidden">
    <div class="px-6 py-4 border-b border-border-card">
        <h2 class="text-base font-semibold text-foreground">Activité récente</h2>
        <p class="text-xs text-muted-foreground mt-0.5">Derniers événements</p>
    </div>
    <ul class="divide-y divide-border-card">
        {@render activity("Nouvelle inscription", "par Paul V.", "Il y a 10 min.")}
        {@render activity("Paiement reçu", "de 60,00 €", "Il y a 32 min.")}
        {@render activity("Mise à jour catalogue", "(3 items)", "Il y a 2h.")}
    </ul>
</div>

{#snippet activity(title: string, subtitle: string, date: string)}
    <li class="px-6 py-3 hover:bg-muted/30 transition-colors">
        <div class="flex items-center gap-1 text-sm">
            <p class="font-medium text-foreground">{title}</p>
            <p class="text-muted-foreground">{subtitle}</p>
        </div>
        <p class="text-xs text-muted-foreground mt-0.5">{date}</p>
    </li>
{/snippet}
