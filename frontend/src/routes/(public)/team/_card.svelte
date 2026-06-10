<script lang="ts">
    import { goto } from "$app/navigation";
    import { type Partner } from "./types";

    let { id, firstname, lastname, occupation, quote, tags, picture }: Partner = $props();
</script>

<div
    onclick={() => goto(`/team/${id}`)}
    class="group cursor-pointer flex flex-col bg-white rounded-3xl border border-border-input hover:border-border-input-hover hover:shadow-card transition-all duration-300 overflow-hidden p-6 md:p-8"
>
    <!-- Header: photo + name/occupation -->
    <div class="flex gap-4 items-center mb-4">
        {#if picture}
            <img
                src={picture}
                alt="{firstname} {lastname}"
                class="rounded-2xl aspect-square w-14 md:w-16 flex-shrink-0 object-cover"
            />
        {:else}
            <div
                class="bg-surface-hover rounded-2xl aspect-square w-14 md:w-16 flex-shrink-0 flex items-center justify-center text-sm font-semibold text-muted-foreground"
            >
                {firstname[0]}{lastname[0]}
            </div>
        {/if}
        <div>
            <h4 class="text-foreground font-semibold text-lg md:text-xl tracking-tight">
                {firstname} {lastname}
            </h4>
            <p class="text-muted-foreground text-sm font-medium mt-0.5">{occupation}</p>
        </div>
    </div>

    <!-- Tags -->
    <div class="flex flex-wrap gap-2 mb-4">
        {#each tags as item}
            <div
                class="inline-flex items-center bg-surface border border-border-input rounded-full px-3 py-1.5"
            >
                <span class="text-xs font-medium text-foreground-alt">{item}</span>
            </div>
        {/each}
    </div>

    <!-- Quote -->
    <p class="text-muted-foreground text-sm md:text-base leading-relaxed mb-6 flex-1">
        {quote}
    </p>

    <!-- CTAs -->
    <div class="flex flex-col sm:flex-row gap-3 items-stretch sm:items-center mt-auto">
        <a
            href="/book?partner={id}"
            onclick={(e) => e.stopPropagation()}
            class="inline-flex justify-center items-center gap-2 bg-foreground hover:bg-foreground-alt text-white text-sm font-medium px-6 py-3 rounded-xl transition-all duration-200 shadow-mini hover:shadow-card"
        >
            Réserver
            <span
                class="iconify group-hover:translate-x-0.5 transition-transform"
                data-icon="lucide:arrow-right"
                data-width="16"
                data-stroke-width="2"
            ></span>
        </a>
        <a
            href="/team/{id}"
            onclick={(e) => e.stopPropagation()}
            class="inline-flex justify-center items-center gap-2 bg-transparent hover:bg-surface border border-border-input-hover hover:border-border-input-hover text-foreground-alt hover:text-foreground text-sm font-medium px-6 py-3 rounded-xl transition-all duration-200"
        >
            Voir le profil
        </a>
    </div>
</div>
