<script lang="ts">
    import { Tabs as TabsPrimitive } from "bits-ui";
    import { cn } from "$lib/utils/design-system";

    interface Props {
        value: string;
        class?: string;
        children?: any;
    }

    let {
        value,
        class: className = "",
        children,
        ...restProps
    }: Props = $props();
</script>

<TabsPrimitive.Content
    {value}
    class={cn(
        "mt-2 ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2",
        "data-[state=active]:animate-in data-[state=active]:fade-in-0 data-[state=active]:slide-in-2",
        "data-[state=inactive]:animate-out data-[state=inactive]:fade-out-0 data-[state=inactive]:slide-out-2",
        className,
    )}
    {...restProps}
>
    {@render children()}
</TabsPrimitive.Content>

