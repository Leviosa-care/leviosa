<script lang="ts">
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    const httpEquivCacheControl: any = "Cache-Control";
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    const httpEquivPragma: any = "Pragma";
</script>

<svelte:head>
    <title>Leviosa — À Venir</title>
    <meta name="robots" content="noindex" />
    <meta http-equiv={httpEquivCacheControl} content="no-store, no-cache, must-revalidate" />
    <meta http-equiv={httpEquivPragma} content="no-cache" />
</svelte:head>

<div class="min-h-screen flex items-center justify-center px-4" style="background-color: var(--foreground);">
    <div class="text-center" style="color: var(--background);">
        <h1 class="text-5xl md:text-6xl font-bold tracking-widest uppercase mb-6" style="font-family: 'Cal Sans', sans-serif;">
            Leviosa
        </h1>
        <p class="text-sm tracking-widest uppercase" style="color: var(--muted-foreground);">
            Quelque chose est en préparation &mdash; à venir
        </p>
    </div>
</div>
