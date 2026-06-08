<script lang="ts">
    import type { PageData } from "./$types";
    import Settings from "./Settings.svelte";

    let { data }: { data: PageData } = $props();
</script>

<svelte:head>
	<title>Paramètres | Staff</title>
</svelte:head>

<div class="p-6 lg:p-10 max-w-4xl">
    <div class="mb-10">
        <p class="text-[11px] font-semibold uppercase tracking-[0.2em] text-muted-foreground mb-3">Configuration</p>
        <h1 class="font-display text-4xl lg:text-5xl font-semibold tracking-tight text-foreground leading-[1.1]">Paramètres</h1>
        <p class="text-muted-foreground mt-3 text-sm">Configurez les informations de l'entreprise et les paramètres système</p>
        <div class="mt-4 h-px w-16 bg-foreground/20"></div>
    </div>
    <Settings {data} />
</div>
