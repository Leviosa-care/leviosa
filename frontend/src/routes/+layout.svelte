<script lang="ts">
    import type { LayoutProps } from "./$types";
    import "../app.css";
    import ToastProvider from "$lib/components/toast/ToastProvider.svelte";
    let { children }: LayoutProps = $props();
</script>

<ToastProvider />
{@render children()}
