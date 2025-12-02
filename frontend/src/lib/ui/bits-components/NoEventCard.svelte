<script lang="ts">
    import { cn } from "$lib/utils/design-system";
    import { CalendarPlus, Search } from "@lucide/svelte";

    interface Props {
        title?: string;
        description?: string;
        actionText?: string;
        onAction?: () => void;
        showSearch?: boolean;
        class?: string;
    }

    let {
        title = "No events found",
        description = "There are no events matching your criteria. Try adjusting your filters or check back later.",
        actionText = "Browse all events",
        onAction,
        showSearch = true,
        class: className = "",
    }: Props = $props();
</script>

<div
    class={cn(
        "flex flex-col items-center justify-center rounded-lg border border-dashed border-border bg-muted/20 p-12 text-center",
        className,
    )}
    role="status"
    aria-live="polite"
>
    <!-- Icon -->
    <div
        class="mb-4 flex h-16 w-16 items-center justify-center rounded-full bg-muted text-muted-foreground"
    >
        {#if showSearch}
            <Search size={32} />
        {:else}
            <CalendarPlus size={32} />
        {/if}
    </div>

    <!-- Content -->
    <div class="max-w-md space-y-2">
        <h3 class="text-lg font-semibold text-foreground">
            {title}
        </h3>

        <p class="text-sm text-muted-foreground">
            {description}
        </p>

        {#if onAction && actionText}
            <button
                onclick={onAction}
                class="mt-4 inline-flex items-center gap-2 rounded-button bg-accent px-4 py-2 text-sm font-medium text-accent-foreground transition-colors hover:bg-accent/80 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
            >
                {actionText}
            </button>
        {/if}
    </div>
</div>

