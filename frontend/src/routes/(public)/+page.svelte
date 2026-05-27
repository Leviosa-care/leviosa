<script lang="ts">
    import Hero from "./_hero.svelte";
    import Services from "./_services.svelte";
    import Why from "./_why.svelte";
    import Testimonies from "./_testimonies.svelte";
    import Ready from "./_ready.svelte";

    // TODO: After fetching from +page.server.ts page.
    // import type { PageProps } from "./$types";
    // let { data }: PageProps = $props();
    // let { cards } = data
</script>

<main class="relative overflow-hidden">
    <Hero />
    <Services />
    <Why />
    <Testimonies />
    <Ready />
</main>
