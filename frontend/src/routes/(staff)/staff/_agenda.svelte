<script lang="ts">
</script>

<div class="bg-background rounded-2xl border border-border overflow-hidden">
    <div class="px-5 py-4 flex items-center justify-between">
        <h2 class="font-display text-sm font-semibold tracking-tight text-foreground">Agenda du jour</h2>
        <span class="text-[11px] text-muted-foreground tabular-nums">3 séances</span>
    </div>
    <div class="relative">
        <!-- Timeline rail -->
        <div class="absolute left-[3.25rem] top-0 bottom-0 w-px bg-border"></div>

        {@render time_slot("10:00", "11:00", "Consultation suivi — Dermatologie", "Dr Sarah L", "Marc D.", "En cours", "bg-green-500", "text-green-500")}
        {@render time_slot("11:30", "12:30", "Bilan initial — Nutrition", "Dr Jean M.", "Sophie K.", "À venir", "bg-blue-500", "text-blue-500")}
        {@render time_slot("14:00", "15:00", "Bilan initial — Kiné", "Dr Claire P.", "Lucas M.", "Annulé", "bg-red-400 dark:bg-red-500", "text-red-500 dark:text-red-400")}
    </div>
    <a
        href="/staff/agenda/reservations"
        class="block text-center py-3.5 text-sm font-medium text-muted-foreground hover:text-foreground transition-colors border-t border-border"
    >
        Agenda complet →
    </a>
</div>

{#snippet time_slot(start: string, end: string, title: string, partner: string, client: string, status: string, dotColor: string, dotTextColor: string)}
    <div class="relative px-5 py-4 flex items-center gap-4 hover:bg-muted/30 transition-colors duration-300 group">
        <!-- Timeline dot -->
        <div class="relative z-10 w-2.5 h-2.5 rounded-full {dotColor} ring-4 ring-background group-hover:scale-125 transition-transform duration-300"></div>

        <!-- Time -->
        <div class="flex flex-col min-w-[2.5rem]">
            <span class="text-sm font-semibold text-foreground tabular-nums leading-none">{start}</span>
            <span class="text-[11px] text-muted-foreground tabular-nums">{end}</span>
        </div>

        <!-- Content -->
        <div class="flex-1 min-w-0">
            <p class="text-sm font-medium text-foreground truncate">{title}</p>
            <p class="text-[11px] text-muted-foreground mt-0.5">{partner} · {client}</p>
        </div>

        <!-- Status pill -->
        <span class="text-[11px] font-semibold {dotTextColor} whitespace-nowrap tracking-wide">
            {status}
        </span>
    </div>
{/snippet}
