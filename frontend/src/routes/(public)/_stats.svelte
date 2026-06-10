<script lang="ts">
    import { reveal } from "$lib/actions/reveal";

    const stats = [
        { value: "2 000+", label: "patients accompagnés" },
        { value: "15+", label: "experts vérifiés" },
        { value: "8", label: "spécialités" },
        { value: "4.9/5", label: "satisfaction client" },
    ];
</script>

<section class="py-14 bg-background-section border-t border-b border-border-input">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="grid grid-cols-2 md:grid-cols-4 gap-8">
            {#each stats as stat, i}
                <div
                    class="flex flex-col items-center text-center {i < 3 ? 'md:border-r md:border-border-input' : ''}"
                    use:reveal={{ preset: "fade-up", delay: 100 + i * 75 }}
                >
                    <span
                        class="text-3xl sm:text-4xl font-semibold tracking-tight text-foreground"
                    >
                        {stat.value}
                    </span>
                    <span class="mt-1.5 text-sm text-muted-foreground">{stat.label}</span>
                </div>
            {/each}
        </div>
    </div>
</section>
