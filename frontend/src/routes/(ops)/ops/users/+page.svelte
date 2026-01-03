<script lang="ts">
    import type { PageData } from "./$types";
    import Users from "./Users.svelte";

    let { data }: { data: PageData } = $props();
</script>

<div class="h-[100vh] flex-1 flex flex-col overflow-hidden bg-gray-50">
    <!-- Header with title and description -->
    <div class="bg-white border-b border-border-card px-6 py-6">
        <div class="grid gap-1">
            <h1 class="text-2xl font-semibold tracking-tight">
                Gestion des utilisateurs
            </h1>
            <p class="text-sm text-foreground-alt">
                Gérez les utilisateurs, approuvez les demandes et modifiez les
                rôles
            </p>
        </div>
    </div>

    <!-- Main content area -->
    <div class="flex-1 overflow-y-auto">
        <Users {data} />
    </div>
</div>
