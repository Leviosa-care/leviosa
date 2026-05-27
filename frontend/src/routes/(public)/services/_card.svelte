<script lang="ts">
    import { Clock } from "@lucide/svelte";

    let { id, title, description, duration, price, tags, image }: {
        id: string;
        title: string;
        description: string;
        duration: number;
        price: string;
        tags: string[];
        image?: string;
    } = $props();
</script>

<div
    class="group flex flex-col md:flex-row md:items-stretch bg-white rounded-3xl border border-dark-100 hover:border-dark-200 hover:shadow-lg transition-all duration-300 overflow-hidden"
>
    <!-- Image or placeholder -->
    {#if image}
        <img
            src={image}
            alt={title}
            class="w-full h-56 md:w-72 md:h-auto flex-shrink-0 object-cover"
        />
    {:else}
        <div
            class="bg-gradient-to-br from-dark-100 to-dark-50 w-full h-56 md:w-72 md:h-auto flex-shrink-0 flex items-center justify-center"
        >
            <span
                class="iconify text-dark-300"
                data-icon="lucide:image"
                data-width="48"
            ></span>
        </div>
    {/if}

    <!-- Content -->
    <div class="flex flex-col gap-4 w-full p-6 md:p-8">
        <!-- Title + Price -->
        <div class="flex flex-col sm:flex-row sm:justify-between sm:items-start gap-3">
            <h4 class="text-dark-900 font-semibold text-xl md:text-2xl tracking-tight">
                {title}
            </h4>
            {#if price}
                <div
                    class="flex-shrink-0 inline-flex items-baseline gap-0.5 bg-dark-900 text-white px-3 py-1.5 rounded-xl self-start"
                >
                    <span class="text-xl font-bold">{price}</span>
                    <span class="text-sm font-medium opacity-70">€</span>
                </div>
            {/if}
        </div>

        <!-- Description -->
        <p class="text-dark-500 text-sm md:text-base leading-relaxed max-w-2xl">
            {description}
        </p>

        <!-- Tags -->
        <div class="flex flex-wrap gap-2">
            <div
                class="inline-flex items-center gap-1.5 bg-dark-50 border border-dark-100 rounded-full px-3 py-1.5"
            >
                <Clock size={13} class="text-dark-400 flex-shrink-0" />
                <span class="text-xs font-medium text-dark-600">{duration} min.</span>
            </div>
            {#each tags as item}
                <div
                    class="inline-flex items-center bg-dark-50 border border-dark-100 rounded-full px-3 py-1.5"
                >
                    <span class="text-xs font-medium text-dark-600">{item}</span>
                </div>
            {/each}
        </div>

        <!-- CTAs -->
        <div class="flex flex-col sm:flex-row gap-3 items-stretch sm:items-center mt-auto pt-2">
            <a href="/book?product={id}">
                <button
                    class="group/btn inline-flex justify-center items-center gap-2 bg-dark-900 hover:bg-dark-800 text-white text-sm font-medium px-6 py-3 rounded-xl transition-all duration-200 shadow-sm hover:shadow-md cursor-pointer w-full sm:w-auto"
                >
                    Réserver
                    <span
                        class="iconify group-hover/btn:translate-x-0.5 transition-transform"
                        data-icon="lucide:arrow-right"
                        data-width="16"
                        data-stroke-width="2"
                    ></span>
                </button>
            </a>
            <button
                class="inline-flex justify-center items-center gap-2 bg-transparent hover:bg-dark-50 border border-dark-200 hover:border-dark-300 text-dark-600 hover:text-dark-900 text-sm font-medium px-6 py-3 rounded-xl transition-all duration-200 cursor-pointer"
            >
                Plus de détails
            </button>
        </div>
    </div>
</div>
