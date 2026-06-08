<script lang="ts">
    import Cards from "./_cards.svelte";
    import Agenda from "./_agenda.svelte";
    import Actions from "./_actions.svelte";
    import Activity from "./_activity.svelte";
    import Volume from "./_volume.svelte";
</script>

<svelte:head>
	<title>Accueil | Staff</title>
</svelte:head>

<div class="container mx-auto px-4 py-8 lg:py-12">
    <div class="mb-8">
        <h1 class="text-3xl lg:text-4xl font-bold mb-1 text-foreground">Accueil</h1>
        <p class="text-muted-foreground">Vue d'ensemble de l'activité de la plateforme</p>
    </div>

    <Cards />
    <div class="mt-8 grid grid-cols-1 lg:grid-cols-[1fr_auto] grid-rows-1 lg:grid-rows-2 gap-6 lg:gap-16">
        <Agenda />
        <Actions />
        <Volume />
        <Activity />
    </div>
</div>
