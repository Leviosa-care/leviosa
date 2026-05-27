<script lang="ts">
    import { goto } from "$app/navigation";
    import { type Partner } from "./types";

    let { id, firstname, lastname, occupation, quote, tags, picture }: Partner = $props();
</script>

<div
    onclick={() => goto(`/team/${id}`)}
    class="group cursor-pointer flex flex-col bg-white rounded-3xl border border-dark-100 hover:border-dark-200 hover:shadow-lg transition-all duration-300 overflow-hidden p-6 md:p-8"
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
                class="bg-dark-100 rounded-2xl aspect-square w-14 md:w-16 flex-shrink-0 flex items-center justify-center text-sm font-semibold text-dark-500"
            >
                {firstname[0]}{lastname[0]}
            </div>
        {/if}
        <div>
            <h4 class="text-dark-900 font-semibold text-lg md:text-xl tracking-tight">
                {firstname} {lastname}
            </h4>
            <p class="text-dark-500 text-sm font-medium mt-0.5">{occupation}</p>
        </div>
    </div>

    <!-- Tags -->
    <div class="flex flex-wrap gap-2 mb-4">
        {#each tags as item}
            <div
                class="inline-flex items-center bg-dark-50 border border-dark-100 rounded-full px-3 py-1.5"
            >
                <span class="text-xs font-medium text-dark-600">{item}</span>
            </div>
        {/each}
    </div>

    <!-- Quote -->
    <p class="text-dark-500 text-sm md:text-base leading-relaxed mb-6 flex-1">
        {quote}
    </p>

    <!-- CTAs -->
    <div class="flex flex-col sm:flex-row gap-3 items-stretch sm:items-center mt-auto">
        <a
            href="/book?partner={id}"
            onclick={(e) => e.stopPropagation()}
            class="inline-flex justify-center items-center gap-2 bg-dark-900 hover:bg-dark-800 text-white text-sm font-medium px-6 py-3 rounded-xl transition-all duration-200 shadow-sm hover:shadow-md"
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
            class="inline-flex justify-center items-center gap-2 bg-transparent hover:bg-dark-50 border border-dark-200 hover:border-dark-300 text-dark-600 hover:text-dark-900 text-sm font-medium px-6 py-3 rounded-xl transition-all duration-200"
        >
            Voir le profil
        </a>
    </div>
</div>
