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
        <h1 class="font-display text-3xl lg:text-4xl font-semibold tracking-tight text-foreground">Paramètres</h1>
        <p class="text-muted-foreground mt-1.5 text-base">Configurez les informations de l'entreprise et les paramètres système</p>
    </div>
    <Settings {data} />
</div>
