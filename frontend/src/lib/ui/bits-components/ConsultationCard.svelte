<script lang="ts">
    import { cn } from "$lib/utils/design-system";
    import { Calendar, Clock, User, Star } from "@lucide/svelte";

    interface Props {
        id: string;
        title: string;
        description: string;
        duration: string;
        price?: string;
        rating?: number;
        imageUrl?: string;
        onclick?: () => void;
        class?: string;
    }

    let {
        id,
        title,
        description,
        duration,
        price,
        rating,
        imageUrl = "https://images.unsplash.com/photo-1576091160399-112ba8d25d1d?w=400&h=250&fit=crop",
        onclick,
        class: className = "",
    }: Props = $props();
</script>

<div
    {id}
    class={cn(
        "group relative overflow-hidden rounded-lg border border-card bg-background text-foreground shadow-card",
        "transition-all duration-300 hover:shadow-popover hover:scale-[1.02] active:scale-[0.98]",
        "cursor-pointer focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2",
        className,
    )}
    onclick={() => onclick?.()}
    role="button"
    tabindex="0"
>
    <!-- Image -->
    <div class="relative aspect-video overflow-hidden">
        <img
            src={imageUrl}
            alt="Consultation preview"
            class="h-full w-full object-cover transition-transform duration-300 group-hover:scale-105"
        />

        <!-- Rating Badge -->
        {#if rating}
            <div
                class="absolute right-4 top-4 flex items-center gap-1 rounded-full bg-background/80 px-2 py-1 text-xs font-medium backdrop-blur-sm"
            >
                <Star size={12} class="fill-accent text-accent" />
                {rating.toFixed(1)}
            </div>
        {/if}
    </div>

    <!-- Content -->
    <div class="p-6">
        <!-- Title -->
        <h3 class="mb-2 text-lg font-semibold text-foreground line-clamp-1">
            {title}
        </h3>

        <!-- Description -->
        <p class="mb-4 text-sm text-muted-foreground line-clamp-2">
            {description}
        </p>

        <!-- Details -->
        <div class="flex items-center justify-between text-sm">
            <div class="flex items-center gap-4">
                <!-- Duration -->
                <div class="flex items-center gap-1 text-muted-foreground">
                    <Clock size={16} />
                    <span>{duration}</span>
                </div>

                <!-- Type -->
                <div class="flex items-center gap-1 text-muted-foreground">
                    <User size={16} />
                    <span>Online</span>
                </div>
            </div>

            <!-- Price -->
            {#if price}
                <div class="text-right">
                    <p class="text-lg font-semibold text-foreground">{price}</p>
                    <p class="text-xs text-muted-foreground">per session</p>
                </div>
            {/if}
        </div>
    </div>
</div>

<style>
    .line-clamp-1 {
        display: -webkit-box;
        -webkit-line-clamp: 1;
        line-clamp: 1;
        -webkit-box-orient: vertical;
        overflow: hidden;
    }

    .line-clamp-2 {
        display: -webkit-box;
        -webkit-line-clamp: 2;
        line-clamp: 2;
        -webkit-box-orient: vertical;
        overflow: hidden;
    }
</style>
