<script lang="ts">
    import { Button } from "bits-ui";
    import { cn } from "$lib/utils/design-system";
    import { Clock } from "@lucide/svelte";

    interface Props {
        selectedTime?: string;
        placeholder?: string;
        disabled?: boolean;
        onclick?: () => void;
        class?: string;
    }

    let {
        selectedTime = "",
        placeholder = "Select time",
        disabled = false,
        onclick,
        class: className = "",
        ...restProps
    }: Props = $props();
</script>

<Button.Root
    variant="outline"
    {disabled}
    onclick={() => onclick?.()}
    class={cn(
        "w-full justify-start text-left font-normal",
        "hover:bg-accent hover:text-accent-foreground",
        className,
    )}
    {...restProps}
>
    <Clock class="mr-2 h-4 w-4 text-muted-foreground" />
    {selectedTime || placeholder}
</Button.Root>

