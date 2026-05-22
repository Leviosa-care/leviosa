<script lang="ts">
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    const httpEquivCacheControl: any = "Cache-Control";
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    const httpEquivPragma: any = "Pragma";
</script>

<svelte:head>
    <title>Leviosa — Mes réservations</title>
    <meta name="robots" content="noindex" />
    <meta http-equiv={httpEquivCacheControl} content="no-store, no-cache, must-revalidate" />
    <meta http-equiv={httpEquivPragma} content="no-cache" />
</svelte:head>

<div class="min-h-screen flex items-center justify-center px-4" style="background-color: var(--dark-900);">
    <div class="text-center max-w-md" style="color: var(--background);">
        <h1 class="text-5xl md:text-6xl font-bold tracking-widest uppercase mb-6" style="font-family: 'Cal Sans', sans-serif;">
            Leviosa
        </h1>
        <p class="text-xl font-semibold mb-4">
            Mes réservations
        </p>
        <p class="text-sm tracking-wide" style="color: var(--dark-300);">
            Cette fonctionnalité est en cours de développement.
        </p>
        <p class="text-sm tracking-wide" style="color: var(--dark-300);">
            À venir très bientôt.
        </p>
    </div>
</div>
