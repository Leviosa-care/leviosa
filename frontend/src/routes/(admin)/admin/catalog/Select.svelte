<script lang="ts">
    import { Select } from "bits-ui";
    import { Check, ChevronDown } from "@lucide/svelte";

    import { formatItems } from "./formatItems";
    import { capitalizeFirstWord } from "$lib/utils/capitalize";

    type Props = {
        name: string;
        items: Set<string>;
        state: string;
    };

    let { name, items, state = $bindable() }: Props = $props();

    const items_arr = formatItems(items);
</script>

<Select.Root type="single" onValueChange={(v) => (state = v)} items={items_arr}>
    <Select.Trigger
        class="h-input rounded-9px border-border-input bg-background inline-flex justify-between gap-2 touch-none select-none items-center border px-2 transition-colors"
        aria-label="Select {name}"
    >
        {capitalizeFirstWord(state)}
        <ChevronDown class="text-muted-foreground ml-auto size-4" />
    </Select.Trigger>
    <Select.Portal>
        <Select.Content
            class="focus-override border-muted bg-background shadow-popover data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 data-[state=closed]:zoom-out-95 data-[state=open]:zoom-in-95 data-[side=bottom]:slide-in-from-top-2 data-[side=left]:slide-in-from-right-2 data-[side=right]:slide-in-from-left-2 data-[side=top]:slide-in-from-bottom-2 outline-hidden z-50 h-full max-h-[var(--bits-select-content-available-height)] w-[var(--bits-select-anchor-width) + 12px] min-w-[var(--bits-select-anchor-width)] select-none rounded-xl border px-1 py-3 data-[side=bottom]:translate-y-1 data-[side=left]:-translate-x-1 data-[side=right]:translate-x-1 data-[side=top]:-translate-y-1"
            sideOffset={10}
        >
            <Select.Viewport class="p-1">
                {#each items_arr as { value, label }, i (i + value)}
                    <Select.Item
                        class="text-left rounded-button data-highlighted:bg-muted outline-hidden data-disabled:opacity-50 flex gap-2 h-10 w-full select-none items-center py-3 pl-5 pr-1.5 text-sm capitalize"
                        {value}
                        {label}
                    >
                        {#snippet children()}
                            {#if state === value}
                                <Check size={14} aria-label="check" />
                            {/if}
                            {value}
                        {/snippet}
                    </Select.Item>
                {/each}
            </Select.Viewport>
        </Select.Content>
    </Select.Portal>
</Select.Root>
