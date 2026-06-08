<script lang="ts">
</script>

<div class="bg-background rounded-xl border border-border overflow-hidden">
    <div class="px-5 py-4 border-b border-border">
        <h2 class="text-sm font-semibold tracking-tight text-foreground uppercase">Volume de consultations</h2>
    </div>
    <div class="p-5">
        <div class="flex items-end gap-2 h-40">
            {#each [35, 50, 42, 68, 55, 78, 62] as height, i}
                <div class="flex-1 flex flex-col items-center gap-1.5">
                    <div class="w-full rounded-md relative overflow-hidden bg-muted" style="height: 100%">
                        <div
                            class="absolute bottom-0 left-0 right-0 rounded-md transition-all duration-500 {i === 5
                                ? 'bg-foreground'
                                : 'bg-foreground/20'}"
                            style="height: {height}%"
                        ></div>
                    </div>
                    <span class="text-[10px] font-medium text-muted-foreground tabular-nums">
                        {['Lun', 'Mar', 'Mer', 'Jeu', 'Ven', 'Sam', 'Dim'][i]}
                    </span>
                </div>
            {/each}
        </div>
    </div>
</div>
