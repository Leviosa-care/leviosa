<script lang="ts">
    import Button from "$lib/ui/Button.svelte";

    import { type Partner } from "./types";

    let { firstname, lastname, occupation, quote, tags }: Partner = $props();
</script>

<div
    class="flex items-center gap-6 md:gap-8 bg-white p-6 md:p-8 rounded-4xl border border-dark-50"
>
    <div class="grid gap-3 md:gap-4 w-full">
        <div class="flex gap-3 md:gap-4 items-center">
            <div class="bg-dark-200 rounded-3xl aspect-square w-10 md:w-12 flex-shrink-0"></div>
            <div class="">
                <h4 class="text-dark-900 font-bold text-lg md:text-2xl">
                    {firstname}
                    {lastname}
                </h4>
                <p class="text-dark-700 text-sm md:text-base font-medium">{occupation}</p>
            </div>
        </div>
        <div class="flex flex-wrap gap-2">
            {#each tags as item}
                {@render tag(item)}
            {/each}
        </div>
        <p class="max-w-2xl text-sm md:text-base text-dark-700">{quote}</p>
        <div class="flex flex-col sm:flex-row gap-2 md:gap-4 max-w-md">
            <Button class="text-white px-6 py-3 md:px-12 md:py-6 rounded-2xl text-sm md:text-base"
                >Voir le profil</Button
            >
            <Button
                class="border border-dark-900 px-6 py-3 md:px-12 md:py-6 bg-transparent text-dark-900 hover:text-dark-800 rounded-2xl text-sm md:text-base"
            >
                <p>Reserver</p>
            </Button>
        </div>
    </div>
</div>

{#snippet tag(name: string)}
    <div class="bg-dark-50 rounded-3xl px-3 py-1.5 md:px-4 md:py-2">
        <p class="text-xs md:text-sm">{name}</p>
    </div>
{/snippet}
