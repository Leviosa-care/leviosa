<script lang="ts">
    import { cn } from "$lib/utils/design-system";

    interface Props {
        count: number;
        currentIndex: number;
        onIndexChange?: (index: number) => void;
        showDots?: boolean;
        showNumbers?: boolean;
        class?: string;
    }

    let {
        count,
        currentIndex = 0,
        onIndexChange,
        showDots = true,
        showNumbers = false,
        class: className = "",
    }: Props = $props();

    function handleDotClick(index: number) {
        if (onIndexChange) {
            onIndexChange(index);
        }
    }
</script>

<div class={cn("flex items-center gap-2", className)}>
    {#if showNumbers}
        <div class="flex items-center gap-1">
            {#each Array(count) as _, index}
                <button
                    onclick={() => handleDotClick(index)}
                    class={cn(
                        "h-8 w-8 rounded-full text-xs font-medium transition-all duration-200",
                        "focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2",
                        index === currentIndex
                            ? "bg-primary text-primary-foreground"
                            : "bg-muted text-muted-foreground hover:bg-accent hover:text-accent-foreground",
                    )}
                >
                    {index + 1}
                </button>
            {/each}
        </div>
    {/if}

    {#if showDots}
        <div class="flex items-center gap-1">
            {#each Array(count) as _, index}
                <button
                    onclick={() => handleDotClick(index)}
                    class={cn(
                        "transition-all duration-200 rounded-full",
                        "focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2",
                        index === currentIndex
                            ? "bg-dark w-8 h-3"
                            : "bg-muted hover:bg-accent w-3 h-3",
                    )}
                    aria-label={`Go to slide ${index + 1}`}
                >
                </button>
            {/each}
        </div>
    {/if}
</div>
