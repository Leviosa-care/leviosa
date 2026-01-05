<script lang="ts">
</script>

<div class="h-fit rounded-xl overflow-hidden border border-border-dark">
    <div
        class="bg-dark-50 px-4 py-4 flex gap-4 items-center border border-b-dark-50"
    >
        <p class="text-dark-900 font-bold">Activite recente</p>
    </div>
    <ul class="px-4 py-4">
        {@render activity(
            "Nouvelle inscription",
            "par Paul V.",
            "Il y a 10 min.",
        )}
        {@render activity("Paiement recu", "de $60.00", "Il y a 32 min.")}
        {@render activity("Mise a jour catalogue", "(3 items)", "Il y a 2h.")}
    </ul>
</div>

{#snippet activity(title: string, subtitle: string, date: string)}
    <li class="py-2">
        <div class="flex items-center gap-1 text-sm">
            <p class="text-dark-900 font-medium">{title}</p>
            <p class="text-dark-600">{subtitle}</p>
        </div>
        <p class="text-dark-500 text-xs">{date}</p>
    </li>
{/snippet}
