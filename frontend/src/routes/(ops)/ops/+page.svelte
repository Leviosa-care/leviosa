<script lang="ts">
    import Cards from "./_cards.svelte";
    import Agenda from "./_agenda.svelte";
    import Actions from "./_actions.svelte";
    import Activity from "./_activity.svelte";
    import Volume from "./_volume.svelte";
</script>

<div class="flex flex-col h-screen">
    <div class="grid gap-1 p-8 border border-dark-50">
        <h1 class="text-2xl font-semibold tracking-tight">Accueil</h1>
        <p class="text-sm text-foreground-alt">
            Vue d'ensemble de l'activite de la plateforme
        </p>
    </div>
    <div class="p-12 bg-dark-50/60 flex-1">
        <Cards />
        <div class="mt-8 grid grid-cols-[1fr_auto] grid-rows-2 gap-16">
            <Agenda />
            <Actions />
            <Volume />
            <Activity />
        </div>
    </div>
</div>
