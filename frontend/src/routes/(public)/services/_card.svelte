<script lang="ts">
    import { Clock } from "@lucide/svelte";
    import { goto } from "$app/navigation";

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
    onclick={() => goto(`/services/${id}`)}
    class="group cursor-pointer flex flex-col md:flex-row md:items-stretch bg-white rounded-3xl border border-border-input hover:border-border-input-hover hover:shadow-card transition-all duration-300 overflow-hidden"
>
    <!-- Image or placeholder -->
    {#if image}
        <div class="relative w-full h-56 md:w-72 md:h-auto flex-shrink-0 overflow-hidden">
            <img
                src={image}
                alt={title}
                class="absolute inset-0 w-full h-full object-cover"
            />
        </div>
    {:else}
        <div
            class="bg-gradient-to-br from-surface-hover to-surface w-full h-56 md:w-72 md:h-auto flex-shrink-0 flex items-center justify-center"
        >
            <span
                class="iconify text-muted-foreground"
                data-icon="lucide:image"
                data-width="48"
            ></span>
        </div>
    {/if}

    <!-- Content -->
    <div class="flex flex-col gap-4 w-full p-6 md:p-8">
        <!-- Title + Price -->
        <div class="flex flex-col sm:flex-row sm:justify-between sm:items-start gap-3">
            <h4 class="text-foreground font-semibold text-xl md:text-2xl tracking-tight">
                {title}
            </h4>
            {#if price}
                <div
                    class="flex-shrink-0 inline-flex items-baseline gap-0.5 bg-foreground text-white px-3 py-1.5 rounded-xl self-start"
                >
                    <span class="text-xl font-bold">{price}</span>
                    <span class="text-sm font-medium opacity-70">€</span>
                </div>
            {/if}
        </div>

        <!-- Description -->
        <p class="text-muted-foreground text-sm md:text-base leading-relaxed max-w-2xl">
            {description}
        </p>

        <!-- Tags -->
        <div class="flex flex-wrap gap-2">
            <div
                class="inline-flex items-center gap-1.5 bg-surface border border-border-input rounded-full px-3 py-1.5"
            >
                <Clock size={13} class="text-muted-foreground flex-shrink-0" />
                <span class="text-xs font-medium text-foreground-alt">{duration} min.</span>
            </div>
            {#each tags as item}
                <div
                    class="inline-flex items-center bg-surface border border-border-input rounded-full px-3 py-1.5"
                >
                    <span class="text-xs font-medium text-foreground-alt">{item}</span>
                </div>
            {/each}
        </div>

        <!-- CTAs -->
        <div class="flex flex-col sm:flex-row gap-3 items-stretch sm:items-center mt-auto pt-2">
            <a href="/book?product={id}" onclick={(e) => e.stopPropagation()}>
                <button
                    class="group/btn inline-flex justify-center items-center gap-2 bg-foreground hover:bg-foreground-alt text-white text-sm font-medium px-6 py-3 rounded-xl transition-all duration-200 shadow-mini hover:shadow-card cursor-pointer w-full sm:w-auto"
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
            <a
                href="/services/{id}"
                onclick={(e) => e.stopPropagation()}
                class="inline-flex justify-center items-center gap-2 bg-transparent hover:bg-surface border border-border-input-hover hover:border-border-input-hover text-foreground-alt hover:text-foreground text-sm font-medium px-6 py-3 rounded-xl transition-all duration-200"
            >
                Plus de détails
            </a>
        </div>
    </div>
</div>
