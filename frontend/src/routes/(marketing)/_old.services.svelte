<script lang="ts">
    import Button from "$lib/ui/Button.svelte";

    // TODO: I should probably adapt the display based on how many categories/products are available

    interface Card {
        title: string;
        description: string;
        cta: string;
    }
    const serviceCards: Card[] = [
        {
            title: "Massage Therapy",
            description:
                "Retroactive treatment designed for recorvery and stress relied tailored to your body's need.",
            cta: "Voir traitement",
        },
        {
            title: "Mindset Coaching",
            description:
                "Unlock your full potential with proven mental clarity strategies, guided mentality, and coginitive restructuring.",
            cta: "Decouvrir les coachs",
        },
        {
            title: "Physical Training",
            description:
                "Personalized strength and conditionning  programs built for peak performance, longevity, and functional movement.",
            cta: "Start training",
        },
    ];
</script>

<div class="p-24 border-green border-4">
    <div class="flex justify-between items-baseline">
        <div>
            <p class="uppercase">Our expertise</p>
            <h2 class="text-dark-900 text-bold text-3xl">Core Services</h2>
            <p class="text-dark-700">
                We focus on three pillars of wellness to ensure a balanced and
                healthy lifestyle for our clients.
            </p>
        </div>
        <Button class="text-white rounded-3xl">
            <a href="/services"> Explore all services </a>
        </Button>
    </div>
    <div>
        <ul class="flex items-center justify-between list-none">
            {#each serviceCards as item}
                {@render card(item.title, item.description, item.cta)}
            {/each}
        </ul>
    </div>
</div>

{#snippet card(title: string, description: string, cta: string)}
    <li class="grid gap-2 p-4 w-full">
        <div class=" bg-dark-50 w-full aspect-square rounded-input"></div>
        <p class="text-dark-900">{title}</p>
        <p class="text-dark-500">{description}</p>
        {#if cta !== ""}
            <Button class="text-white">{cta}</Button>
        {/if}
    </li>
{/snippet}
