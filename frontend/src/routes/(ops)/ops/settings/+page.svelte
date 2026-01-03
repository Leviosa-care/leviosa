<script lang="ts">
    import type { PageData } from "./$types";
    import Settings from "./Settings.svelte";

    let { data }: { data: PageData } = $props();
</script>

<div class="h-[100vh] flex-1 flex flex-col overflow-hidden bg-gray-50">
    <!-- Header with title and description -->
    <div class="bg-white border-b border-border-card px-6 py-6">
        <div class="grid gap-1">
            <h1 class="text-2xl font-semibold tracking-tight">
                Paramètres
            </h1>
            <p class="text-sm text-foreground-alt">
                Configurez les informations de l'entreprise et les paramètres système
            </p>
        </div>
    </div>

    <!-- Main content area -->
    <Settings {data} />
</div>
