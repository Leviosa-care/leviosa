<script lang="ts">
    import { cn } from "$lib/utils/design-system";
    import { Calendar, MapPin, Clock, Sparkles } from "@lucide/svelte";

    interface Props {
        id: string;
        eventType: string;
        title: string;
        date: string;
        eventBegin: string;
        city: string;
        prestationType?: string;
        prestationStart?: string;
        imageUrl?: string;
        onclick?: () => void;
        class?: string;
    }

    let {
        id,
        eventType,
        title,
        date,
        eventBegin,
        city,
        prestationType = "",
        prestationStart = "",
        imageUrl = "https://images.unsplash.com/photo-1541339907198-e08756dedf3f?w=400&h=225&fit=crop",
        onclick,
        class: className = "",
    }: Props = $props();
</script>

<div
    {id}
    class={cn(
        "group relative overflow-hidden rounded-lg border border-card bg-background text-foreground shadow-card",
        "transition-all duration-200 hover:shadow-popover hover:scale-[1.02] active:scale-[0.98]",
        "cursor-pointer focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2",
        className,
    )}
    onclick={() => onclick?.()}
    role="button"
    tabindex="0"
>
    <!-- Image Container -->
    <div class="relative aspect-video overflow-hidden">
        <img
            src={imageUrl}
            alt="Event location"
            class="h-full w-full object-cover transition-transform duration-300 group-hover:scale-105"
        />

        <!-- Prestation Tag -->
        {#if prestationType && prestationStart}
            <div
                class="absolute left-4 top-4 flex items-center gap-2 rounded-full bg-accent px-3 py-1.5 text-sm font-medium text-accent-foreground shadow-mini"
            >
                <Sparkles size={16} class="text-accent-foreground" />
                <span class="capitalize">
                    {prestationType} :
                    <span class="font-semibold">{prestationStart}</span>
                </span>
            </div>
        {/if}
    </div>

    <!-- Content -->
    <div class="p-6">
        <!-- Title -->
        <div class="mb-4">
            <h3 class="text-lg font-semibold text-foreground line-clamp-2">
                <span class="font-bold">{eventType}</span> - {title}
            </h3>
        </div>

        <!-- Event Details -->
        <div class="space-y-2">
            <!-- Date and Time -->
            <div class="flex items-center gap-2 text-sm text-muted-foreground">
                <Calendar size={16} class="text-accent" />
                <span class="font-medium">{date} • {eventBegin}</span>
            </div>

            <!-- Location -->
            <div class="flex items-center gap-2 text-sm text-muted-foreground">
                <MapPin size={16} class="text-accent" />
                <span>{city}</span>
            </div>
        </div>
    </div>
</div>

<style>
    .line-clamp-2 {
        display: -webkit-box;
        -webkit-line-clamp: 2;
        line-clamp: 2;
        -webkit-box-orient: vertical;
        overflow: hidden;
    }
</style>
