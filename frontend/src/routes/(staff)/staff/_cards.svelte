<script lang="ts">
    import { TrendingUp } from "@lucide/svelte";
    // TODO: add  parameter to card snippet
</script>

<div class="grid grid-cols-4 gap-8">
    {@render card()}
    {@render card()}
    {@render card()}
    {@render card()}
</div>

{#snippet card()}
    <div
        class="bg-dark-50 border border-border-card rounded-lg p-6 w-full max-w-90 grid gap-6 shadow-mini"
    >
        <div>
            <div class="flex items-center justify-between">
                <p class="text-dark-500 text-md font-medium capitalize">
                    Total Revenue
                </p>
                <div
                    class="text-sm font-medium flex items-center gap-1 border border-border-input hover:border-border-input-hover px-2 py-1 rounded-3xl"
                >
                    <TrendingUp size={16} />
                    <p>+12.5%</p>
                </div>
            </div>
            <p class="text-dark-900 text-3xl font-bold">$1,250.00</p>
        </div>
        <div class="grid gap-2">
            <div class="flex items-center gap-2">
                <p class="text-sm text-dark-900">Trending up this month</p>
                <TrendingUp size={16} />
            </div>
            <p class="text-sm text-dark-500">Visitors for the last 6 months</p>
        </div>
    </div>
{/snippet}
