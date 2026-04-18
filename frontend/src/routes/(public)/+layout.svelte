<script lang="ts">
    import type { LayoutProps } from "./$types";
    import "../../app.css";
    import Navigation from "./navigation.svelte";
    import Footer from "./footer.svelte";
    let { data, children }: LayoutProps = $props();
    let { permissions } = data;
</script>

<div class="grid">
    <Navigation {permissions} />
    {@render children()}
    <Footer />
</div>
