<script lang="ts">
    import type { LayoutProps } from "./$types";
    let { children }: LayoutProps = $props();
    // NOTE: here we have the sidebar for authenticated users
</script>

{@render children()}
