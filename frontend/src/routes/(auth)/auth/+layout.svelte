<script lang="ts">
    import type { LayoutProps } from "./$types";
    import "../../../app.css";
    let { children }: LayoutProps = $props();

    import { setRegistrationContext } from "$lib/context/register.svelte";
    setRegistrationContext();
</script>

{@render children()}
