<script lang="ts">
    import Hero from "./_hero.svelte";
    import How from "./_how.svelte";
    import Stats from "./_stats.svelte";
    import Services from "./_services.svelte";
    import Why from "./_why.svelte";
    import Team from "./_team.svelte";
    import Testimonies from "./_testimonies.svelte";
    import Faq from "./_faq.svelte";
    import Ready from "./_ready.svelte";

    // TODO: After fetching from +page.server.ts page.
    // import type { PageProps } from "./$types";
    // let { data }: PageProps = $props();
    // let { cards } = data
</script>

<main class="relative overflow-hidden">
    <Hero />
    <How />
    <Services />
    <Stats />
    <Why />
    <Team />
    <Testimonies />
    <Faq />
    <Ready />
</main>
