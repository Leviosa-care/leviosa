<script lang="ts">
    import Cards from "./_cards.svelte";
    import Agenda from "./_agenda.svelte";
    import Actions from "./_actions.svelte";
    import Activity from "./_activity.svelte";
    import Volume from "./_volume.svelte";

    import { onMount } from 'svelte';
    import { browser } from '$app/environment';

    let visible = $state(false);
    onMount(() => { visible = true; });
</script>

<svelte:head>
	<title>Accueil | Staff</title>
</svelte:head>

<div class="relative min-h-full">
    <!-- Subtle grain texture overlay -->
    <div class="pointer-events-none fixed inset-0 opacity-[0.015] mix-blend-multiply dark:opacity-[0.03] dark:mix-blend-overlay z-0"
        style="background-image: url('data:image/svg+xml,%3Csvg viewBox=%220 0 256 256%22 xmlns=%22http://www.w3.org/2000/svg%22%3E%3Cfilter id=%22noise%22%3E%3CfeTurbulence type=%22fractalNoise%22 baseFrequency=%220.9%22 numOctaves=%224%22 stitchTiles=%22stitch%22/%3E%3C/filter%3E%3Crect width=%22100%25%22 height=%22100%25%22 filter=%22url(%23noise)%22/%3E%3C/svg%3E');"
    ></div>

    <div class="relative z-10 p-6 lg:p-10">
        <!-- Page header with staggered entrance -->
        <div class="mb-12 transition-all duration-700 ease-out {visible ? 'opacity-100 translate-y-0' : 'opacity-0 translate-y-4'}">
            <p class="text-xs font-semibold uppercase tracking-[0.2em] text-muted-foreground mb-3">Tableau de bord</p>
            <h1 class="font-display text-4xl lg:text-5xl font-semibold tracking-tight text-foreground leading-[1.1]">Bonjour{''}</h1>
            <div class="mt-4 h-px w-16 bg-foreground/20"></div>
        </div>

        <!-- KPI row with staggered delay -->
        <div class="transition-all duration-700 delay-150 ease-out {visible ? 'opacity-100 translate-y-0' : 'opacity-0 translate-y-4'}">
            <Cards />
        </div>

        <!-- Two-column layout -->
        <div class="mt-10 grid grid-cols-1 xl:grid-cols-[1fr_360px] gap-8">
            <div class="space-y-6 transition-all duration-700 delay-300 ease-out {visible ? 'opacity-100 translate-y-0' : 'opacity-0 translate-y-4'}">
                <Agenda />
                <Volume />
            </div>
            <div class="space-y-6 transition-all duration-700 delay-[400ms] ease-out {visible ? 'opacity-100 translate-y-0' : 'opacity-0 translate-y-4'}">
                <Actions />
                <Activity />
            </div>
        </div>
    </div>
</div>
