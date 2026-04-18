<script lang="ts">
    import { Combobox } from "bits-ui";
    import {
        Search,
        ChevronsDown,
        ChevronsUp,
        ChevronsUpDown,
        Check,
    } from "@lucide/svelte";

    import { capitalizeFirstWord } from "$lib/utils/capitalize";
    import { type Card } from "./products";
    interface Props {
        cards: Card[];
        searchValue: string;
    }
    let { cards, searchValue = $bindable() }: Props = $props();
    const products = cards.map((card) => {
        return {
            value: card.name,
            label: capitalizeFirstWord(card.name),
        };
    });

    const filteredProducts = $derived(
        searchValue === ""
            ? products
            : products.filter((product) =>
                  product.label
                      .toLowerCase()
                      .includes(searchValue.toLowerCase()),
              ),
    );
</script>

<Combobox.Root
    type="multiple"
    name="affaire"
    onOpenChange={(o) => {
        if (!o) searchValue = "";
    }}
>
    <div class="relative">
        <Search
            class="text-muted-foreground absolute start-3 top-1/2 size-6 -translate-y-1/2"
        />
        <Combobox.Input
            oninput={(e) => (searchValue = e.currentTarget.value)}
            class="h-input rounded-9px border-border-input bg-background placeholder:text-foreground-alt/50 focus:ring-foreground focus:ring-offset-background focus:outline-hidden inline-flex w-full truncate border px-11 text-base transition-colors focus:ring-2 focus:ring-offset-2 sm:text-sm"
            placeholder="Chercher un produit par son nom ou sa description"
            aria-label="Chercher un produit par son nom ou sa description"
        />
        <Combobox.Trigger
            class="absolute end-3 top-1/2 size-6 -translate-y-1/2"
        >
            <ChevronsUpDown class="text-muted-foreground size-6" />
        </Combobox.Trigger>
    </div>
    <Combobox.Portal>
        <Combobox.Content
            class="focus-override border-muted bg-background shadow-popover data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 data-[state=closed]:zoom-out-95 data-[state=open]:zoom-in-95 data-[side=bottom]:slide-in-from-top-2 data-[side=left]:slide-in-from-right-2 data-[side=right]:slide-in-from-left-2 data-[side=top]:slide-in-from-bottom-2 outline-hidden z-50 h-96 max-h-[var(--bits-combobox-content-available-height)] w-[var(--bits-combobox-anchor-width)] min-w-[var(--bits-combobox-anchor-width)] select-none rounded-xl border px-1 py-3 data-[side=bottom]:translate-y-1 data-[side=left]:-translate-x-1 data-[side=right]:translate-x-1 data-[side=top]:-translate-y-1"
            sideOffset={10}
        >
            <Combobox.ScrollUpButton
                class="flex w-full items-center justify-center py-1"
            >
                <ChevronsUp class="size-3" />
            </Combobox.ScrollUpButton>
            <Combobox.Viewport class="p-1">
                {#each filteredProducts as fruit, i (i + fruit.value)}
                    <Combobox.Item
                        class="rounded-button data-highlighted:bg-muted outline-hidden flex h-10 w-full select-none items-center py-3 pl-5 pr-1.5 text-sm capitalize"
                        value={fruit.value}
                        label={fruit.label}
                    >
                        {#snippet children({ selected })}
                            {fruit.label}
                            {#if selected}
                                <div class="ml-auto">
                                    <Check />
                                </div>
                            {/if}
                        {/snippet}
                    </Combobox.Item>
                {:else}
                    <span class="block px-5 py-2 text-sm text-muted-foreground">
                        No results found, try again.
                    </span>
                {/each}
            </Combobox.Viewport>
            <Combobox.ScrollDownButton
                class="flex w-full items-center justify-center py-1"
            >
                <ChevronsDown class="size-3" />
            </Combobox.ScrollDownButton>
        </Combobox.Content>
    </Combobox.Portal>
</Combobox.Root>
