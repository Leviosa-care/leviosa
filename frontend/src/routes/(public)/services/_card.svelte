<script lang="ts">
    import { ChevronDown, Clock } from "@lucide/svelte";
    import Button from "$lib/ui/Button.svelte";

    let { id, title, description, duration, price, tags }: {
        id: string;
        title: string;
        description: string;
        duration: number;
        price: string;
        tags: string[];
    } = $props();
</script>

<div
    class="flex flex-col md:flex-row md:items-center gap-6 md:gap-8 bg-white p-6 md:p-8 rounded-4xl"
>
    <div
        class="bg-dark-200 rounded-3xl w-full h-48 md:aspect-square md:w-64 md:h-auto flex-shrink-0"
    ></div>
    <div class="grid gap-4 w-full">
        <div
            class="flex flex-col sm:flex-row sm:justify-between sm:items-center gap-2"
        >
            <h4 class="text-dark-900 font-bold text-xl md:text-2xl">{title}</h4>
            {#if price}<p class="text-dark-900 font-bold text-xl md:text-2xl">{price}€</p>{/if}
        </div>
        <p class="max-w-2xl text-dark-700">{description}</p>
        <div class="flex flex-wrap gap-2">
            {@render tag(`${duration} min.`)}
            {#each tags as item}
                {@render tag(item)}
            {/each}
        </div>
        <div
            class="flex flex-col sm:flex-row gap-2 items-stretch sm:items-center"
        >
            <a href="/book?product={id}">
                <Button class="text-white px-8 py-4 md:px-12 md:py-6 rounded-2xl cursor-pointer"
                    >Réserver</Button
                >
            </a>
            <Button
                class="bg-transparent hover:bg-dark-50/30 px-8 py-4 md:px-4 md:py-2"
            >
                <p>Plus de détails</p>
                <ChevronDown />
            </Button>
        </div>
    </div>
</div>

{#snippet tag(name: string)}
    <div class="flex gap-3 bg-dark-50 rounded-3xl px-4 py-2">
        <Clock class="text-dark-600" />
        <p>{name}</p>
    </div>
{/snippet}
