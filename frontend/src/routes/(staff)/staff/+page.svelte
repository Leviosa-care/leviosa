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

<div class="p-6 lg:p-10">
    <!-- Page header -->
    <div class="mb-10">
        <h1 class="font-display text-3xl lg:text-4xl font-semibold tracking-tight text-foreground">Accueil</h1>
        <p class="text-muted-foreground mt-1.5 text-base">Vue d'ensemble de l'activité de la plateforme</p>
    </div>

    <!-- KPI row -->
    <Cards />

    <!-- Two-column layout: agenda + sidebar -->
    <div class="mt-8 grid grid-cols-1 xl:grid-cols-[1fr_340px] gap-6">
        <div class="space-y-6">
            <Agenda />
            <Volume />
        </div>
        <div class="space-y-6">
            <Actions />
            <Activity />
        </div>
    </div>
</div>
