<script lang="ts">
    import { Button } from "bits-ui";
</script>

<div class="bg-background rounded-xl border border-border overflow-hidden">
    <div class="px-5 py-4 border-b border-border">
        <h2 class="text-sm font-semibold tracking-tight text-foreground uppercase">Agenda du jour</h2>
    </div>
    <div class="divide-y divide-border">
        {@render time_slot("10:00", "11:00", "Consultation suivi — Dermatologie", "Dr Sarah L", "Marc D.", "En cours", "bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-400")}
        {@render time_slot("11:30", "12:30", "Bilan initial — Nutrition", "Dr Jean M.", "Sophie K.", "À venir", "bg-blue-100 text-blue-700 dark:bg-blue-900/30 dark:text-blue-400")}
        {@render time_slot("14:00", "15:00", "Bilan initial — Kiné", "Dr Claire P.", "Lucas M.", "Annulé", "bg-red-100 text-red-700 dark:bg-red-900/30 dark:text-red-400")}
    </div>
    <a
        href="/staff/agenda/reservations"
        class="block text-center py-3.5 text-sm font-medium text-muted-foreground hover:text-foreground hover:bg-muted/30 transition-colors border-t border-border"
    >
        Voir l'agenda complet →
    </a>
</div>

{#snippet time_slot(start: string, end: string, title: string, partner: string, client: string, status: string, statusClasses: string)}
    <div class="px-5 py-4 flex items-center gap-4 hover:bg-muted/20 transition-colors duration-200">
        <div class="flex flex-col items-center min-w-[3.5rem]">
            <span class="text-sm font-semibold text-foreground tabular-nums leading-none">{start}</span>
            <span class="text-[11px] text-muted-foreground tabular-nums">{end}</span>
        </div>
        <div class="w-px h-10 bg-border"></div>
        <div class="flex-1 min-w-0">
            <p class="text-sm font-medium text-foreground truncate">{title}</p>
            <p class="text-xs text-muted-foreground mt-0.5">{partner} — Client: {client}</p>
        </div>
        <span class="text-xs font-medium px-2.5 py-1 rounded-full whitespace-nowrap {statusClasses}">
            {status}
        </span>
    </div>
{/snippet}
