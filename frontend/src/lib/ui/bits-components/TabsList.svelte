<script lang="ts">
    import { Tabs as TabsPrimitive } from "bits-ui";
    import { cn } from "$lib/utils/design-system";

    interface Props {
        class?: string;
        children?: any;
    }

    let { class: className = "", children, ...restProps }: Props = $props();
</script>

<TabsPrimitive.List
    class={cn(
        "inline-flex h-10 items-center justify-center rounded-md bg-muted p-1 text-muted-foreground",
        "focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2",
        "disabled:pointer-events-none disabled:opacity-50",
        className,
    )}
    {...restProps}
>
    {@render children()}
</TabsPrimitive.List>

