<script lang="ts">
    import { Tabs as TabsPrimitive } from "bits-ui";
    import { cn } from "$lib/utils/design-system";

    interface Props {
        value?: string;
        onValueChange?: (value: string) => void;
        orientation?: "horizontal" | "vertical";
        class?: string;
        children?: any;
    }

    let {
        value,
        onValueChange,
        orientation = "horizontal",
        class: className = "",
        children,
        ...restProps
    }: Props = $props();

    // Forward value binding to bits-ui Tabs primitive
    let valueBinding = $bindable(value);
</script>

<TabsPrimitive.Root
    bind:value={valueBinding}
    onValueChange={(v) => onValueChange?.(v)}
    {orientation}
    class={cn(
        "flex",
        orientation === "horizontal" ? "flex-col" : "flex-row",
        className,
    )}
    {...restProps}
>
    {#if children}
        {@render children()}
    {:else}
        <slot />
    {/if}
</TabsPrimitive.Root>

