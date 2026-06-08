<script lang="ts">
    import { Button } from "bits-ui";
</script>

<div class="bg-background rounded-lg border border-border-card overflow-hidden">
    <div class="px-6 py-4 border-b border-border-card">
        <h2 class="text-base font-semibold text-foreground">Agenda du jour</h2>
        <p class="text-xs text-muted-foreground mt-0.5">Réservations à venir</p>
    </div>
    <div class="divide-y divide-border-card">
        {@render time_slot(
            "10:00",
            "Consultation suivi - Dermatologie",
            "Dr Sarah L",
            "Marc D.",
            "En cours",
            true,
        )}
        {@render time_slot(
            "11:30",
            "Bilan initial - Nutrition",
            "Dr Jean M.",
            "Sophie K.",
            "À venir",
            false,
        )}
        {@render time_slot(
            "14:00",
            "Bilan initial - Kiné",
            "Dr Claire P.",
            "Lucas M.",
            "Annulé",
            false,
        )}
    </div>
    <a
        href="/staff/agenda/reservations"
        class="block text-center py-3 text-sm font-medium text-muted-foreground hover:text-foreground hover:bg-muted/40 transition-colors border-t border-border-card"
    >
        Voir l'agenda complet
    </a>
</div>

{#snippet time_slot(
    time: string,
    title: string,
    partner: string,
    client: string,
    status: string,
    active: boolean,
)}
    <div class="px-6 py-4 flex items-center gap-4 hover:bg-muted/30 transition-colors">
        <div
            class="p-3 bg-muted rounded-lg text-sm font-medium text-foreground tabular-nums"
        >
            {time}
        </div>
        <div class="flex-1 min-w-0">
            <p class="text-sm font-medium text-foreground truncate">
                {title}
            </p>
            <p class="text-xs text-muted-foreground">
                {partner} — Client: {client}
            </p>
        </div>
        <span class="text-xs font-medium px-2.5 py-0.5 rounded-full whitespace-nowrap {status === 'En cours'
            ? 'bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-400'
            : status === 'Annulé'
            ? 'bg-red-100 text-red-700 dark:bg-red-900/30 dark:text-red-400'
            : 'bg-blue-100 text-blue-700 dark:bg-blue-900/30 dark:text-blue-400'}">
            {status}
        </span>
    </div>
{/snippet}
