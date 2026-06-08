<script lang="ts"></script>

<div class="bg-background rounded-lg border border-border-card overflow-hidden">
    <div class="px-6 py-4 border-b border-border-card">
        <h2 class="text-base font-semibold text-foreground">Volume de consultations</h2>
        <p class="text-xs text-muted-foreground mt-0.5">Les 7 derniers jours</p>
    </div>
    <div class="w-full h-48 p-6 flex items-end gap-2">
        {#each [35, 50, 42, 68, 55, 78, 62] as height, i}
            <div class="flex-1 bg-muted rounded-md relative overflow-hidden">
                <div
                    class="absolute bottom-0 left-0 right-0 bg-foreground/80 rounded-md transition-all"
                    style="height: {height}%"
                ></div>
            </div>
        {/each}
    </div>
</div>
