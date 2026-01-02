<script lang="ts">
    import { Tabs as TabsPrimitive } from "bits-ui";
    import { cn } from "$lib/utils/design-system";
    import type { Snippet } from "svelte";

    interface Props {
        value: string;
        disabled?: boolean;
        class?: string;
        children?: Snippet;
    }

    let {
        value,
        disabled = false,
        class: className = "",
        children,
        ...restProps
    }: Props = $props();
</script>

<TabsPrimitive.Trigger
    {value}
    {disabled}
    class={cn(
        "inline-flex items-center justify-center whitespace-nowrap rounded-sm px-3 py-1.5 text-sm font-medium ring-offset-background transition-all",
        "focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2",
        "disabled:pointer-events-none disabled:opacity-50",
        "data-[state=active]:bg-background data-[state=active]:text-foreground data-[state=active]:shadow-sm",
        "data-[state=inactive]:hover:bg-accent hover:text-accent-foreground",
        className,
    )}
    {...restProps}
>
    {#if children}
        {@render children()}
    {/if}
</TabsPrimitive.Trigger>
