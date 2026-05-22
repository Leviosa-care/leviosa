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
    {disabled}
    onclick={() => onclick?.()}
    class={cn(
        "w-full justify-start text-left font-normal border border-input bg-background hover:bg-accent hover:text-accent-foreground rounded-input px-3 py-2 text-sm",
        "focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2",
        "disabled:cursor-not-allowed disabled:opacity-50",
        className,
    )}
    {...restProps}
>
    <Clock class="mr-2 h-4 w-4 text-muted-foreground" />
    {selectedTime || placeholder}
</Button.Root>
