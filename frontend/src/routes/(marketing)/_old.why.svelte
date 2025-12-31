<script lang="ts">
    interface Card {
        title: string;
        description: string;
    }
    const serviceCards: Card[] = [
        {
            title: "Massage Therapy",
            description:
                "Retroactive treatment designed for recorvery and stress relied tailored to your body's need.",
        },
        {
            title: "Mindset Coaching",
            description:
                "Unlock your full potential with proven mental clarity strategies, guided mentality, and coginitive restructuring.",
        },
        {
            title: "Physical Training",
            description:
                "Personalized strength and conditionning  programs built for peak performance, longevity, and functional movement.",
        },
    ];
</script>

<div class="bg-dark-50 p-24">
    <div>
        <h2 class="text-3xl font-bold">Why Choose Leviosa</h2>
        <p class="text-dark-700">
            We believe in sustainable, holistic results driven by expers
            guidance and a supportive community
        </p>
    </div>

    <div>
        <ul class="flex items-center justify-between list-none">
            {#each serviceCards as item}
                {@render card(item.title, item.description)}
            {/each}
        </ul>
    </div>
</div>

{#snippet card(title: string, description: string)}
    <li class="grid gap-2 p-4 w-full">
        <div class=" bg-dark-100 w-full aspect-square rounded-input"></div>
        <p class="text-dark-900">{title}</p>
        <p class="text-dark-500">{description}</p>
    </li>
{/snippet}
