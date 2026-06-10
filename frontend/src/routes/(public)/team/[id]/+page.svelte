<script lang="ts">
    import { ChevronLeft, CalendarCheck, Clock, Briefcase } from "@lucide/svelte";
    import { reveal } from "$lib/actions/reveal";
    import type { PageProps } from "./$types";

    let { data }: PageProps = $props();

    const { partner, products, categories } = data;
</script>

<div
    class="min-h-screen bg-white"
    style="background-image: radial-gradient(rgba(15,23,42,0.035) 1px, transparent 1px); background-size: 24px 24px;"
>
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-10 md:py-14">
        <!-- Back link -->
        <a
            href="/team"
            class="inline-flex items-center gap-1.5 text-sm text-muted-foreground hover:text-foreground transition-colors duration-150 mb-10 group"
            use:reveal={{ preset: "fade-down", delay: 50 }}
        >
            <ChevronLeft size={16} class="group-hover:-translate-x-0.5 transition-transform duration-150" />
            Tous les praticiens
        </a>

        <div class="grid grid-cols-1 lg:grid-cols-[1fr_320px] gap-10 lg:gap-16 items-start">
            <!-- Main content -->
            <div>
                <!-- Profile header -->
                <div
                    class="flex flex-col sm:flex-row gap-6 items-start sm:items-center mb-8"
                    use:reveal={{ preset: "fade-up", delay: 100 }}
                >
                    {#if partner.picture}
                        <img
                            src={partner.picture}
                            alt="{partner.firstname} {partner.lastname}"
                            class="w-24 h-24 md:w-32 md:h-32 rounded-3xl object-cover flex-shrink-0"
                        />
                    {:else}
                        <div
                            class="w-24 h-24 md:w-32 md:h-32 rounded-3xl bg-gradient-to-br from-surface-hover to-surface flex items-center justify-center text-2xl font-semibold text-muted-foreground flex-shrink-0"
                        >
                            {partner.firstname[0]}{partner.lastname[0]}
                        </div>
                    {/if}
                    <div>
                        {#if categories.length > 0}
                            <p class="text-xs font-semibold text-muted-foreground uppercase tracking-wider mb-2">
                                {categories.map((c: any) => c.name).join(" · ")}
                            </p>
                        {/if}
                        <h1 class="text-3xl sm:text-4xl md:text-5xl font-semibold tracking-tight text-foreground">
                            {partner.firstname} {partner.lastname}
                        </h1>
                        <p class="text-lg text-muted-foreground mt-1">{partner.occupation}</p>
                    </div>
                </div>

                <!-- Tags -->
                {#if partner.tags.length > 0}
                    <div
                        class="flex flex-wrap gap-2 mb-8"
                        use:reveal={{ preset: "fade-up", delay: 150 }}
                    >
                        {#each partner.tags as tag}
                            <div
                                class="inline-flex items-center bg-surface border border-border-input rounded-full px-4 py-2"
                            >
                                <span class="text-xs font-medium text-foreground-alt">{tag}</span>
                            </div>
                        {/each}
                    </div>
                {/if}

                <!-- Quote -->
                {#if partner.quote}
                    <blockquote
                        class="border-l-2 border-border-input-hover pl-5 mb-10 italic text-muted-foreground text-base md:text-lg leading-relaxed"
                        use:reveal={{ preset: "fade-up", delay: 200 }}
                    >
                        « {partner.quote} »
                    </blockquote>
                {/if}

                <!-- Bio -->
                {#if partner.bio}
                    <div use:reveal={{ preset: "fade-up", delay: 250 }}>
                        <h2 class="text-lg font-semibold text-foreground mb-3">À propos</h2>
                        <p class="text-muted-foreground leading-relaxed text-base md:text-lg">
                            {partner.bio}
                        </p>
                    </div>
                {/if}

                <!-- Experience -->
                {#if partner.experience}
                    <div class="mt-10 pt-10 border-t border-border-input" use:reveal={{ preset: "fade-up", delay: 300 }}>
                        <h2 class="text-lg font-semibold text-foreground mb-3">Formation & expérience</h2>
                        <p class="text-muted-foreground leading-relaxed">
                            {partner.experience}
                        </p>
                    </div>
                {/if}

                <!-- Associated services -->
                {#if products.length > 0}
                    <div class="mt-10 pt-10 border-t border-border-input" use:reveal={{ preset: "fade-up", delay: 350 }}>
                        <h2 class="text-lg font-semibold text-foreground mb-4">Services proposés</h2>
                        <div class="grid gap-3">
                            {#each products as product}
                                <a
                                    href="/services/{product.id}"
                                    class="group flex items-center justify-between bg-white border border-border-input hover:border-border-input-hover hover:shadow-mini rounded-2xl px-5 py-4 transition-all duration-200"
                                >
                                    <div class="flex items-center gap-3">
                                        <Briefcase size={16} class="text-muted-foreground flex-shrink-0" />
                                        <div>
                                            <p class="text-foreground font-medium text-sm">{product.name}</p>
                                            {#if product.description}
                                                <p class="text-muted-foreground text-xs mt-0.5 line-clamp-1">{product.description}</p>
                                            {/if}
                                        </div>
                                    </div>
                                    <div class="flex items-center gap-3 flex-shrink-0 ml-4">
                                        <span class="inline-flex items-center gap-1 text-xs text-muted-foreground">
                                            <Clock size={12} />
                                            {product.duration} min
                                        </span>
                                        <span
                                            class="iconify text-muted-foreground group-hover:translate-x-0.5 transition-transform"
                                            data-icon="lucide:chevron-right"
                                            data-width="16"
                                        ></span>
                                    </div>
                                </a>
                            {/each}
                        </div>
                    </div>
                {/if}
            </div>

            <!-- Sticky sidebar: booking card -->
            <div class="lg:sticky lg:top-24" use:reveal={{ preset: "fade-up", delay: 200 }}>
                <div class="bg-white border border-border-input rounded-3xl p-6 shadow-mini">
                    <p class="text-sm font-semibold text-foreground mb-1">
                        {partner.firstname} {partner.lastname}
                    </p>
                    <p class="text-sm text-muted-foreground mb-6">{partner.occupation}</p>

                    <!-- Réserver CTA -->
                    <a href="/book?partner={partner.id}" class="block w-full">
                        <button
                            class="group/btn w-full inline-flex justify-center items-center gap-2 bg-foreground hover:bg-foreground-alt text-white text-sm font-medium px-6 py-3.5 rounded-xl transition-all duration-200 shadow-mini hover:shadow-card cursor-pointer"
                        >
                            <CalendarCheck size={16} />
                            Réserver une séance
                            <span
                                class="iconify group-hover/btn:translate-x-0.5 transition-transform"
                                data-icon="lucide:arrow-right"
                                data-width="16"
                                data-stroke-width="2"
                            ></span>
                        </button>
                    </a>

                    <!-- Tags list -->
                    {#if partner.tags.length > 0}
                        <div class="mt-6 pt-6 border-t border-border-input">
                            <p class="text-xs font-semibold text-muted-foreground uppercase tracking-wider mb-3">
                                Spécialités
                            </p>
                            <ul class="space-y-2">
                                {#each partner.tags as tag}
                                    <li class="flex items-center gap-2 text-sm text-muted-foreground">
                                        <span class="iconify text-muted-foreground flex-shrink-0" data-icon="lucide:check" data-width="14"></span>
                                        {tag}
                                    </li>
                                {/each}
                            </ul>
                        </div>
                    {/if}
                </div>
            </div>
        </div>
    </div>
</div>
