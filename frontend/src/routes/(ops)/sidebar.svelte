<script lang="ts">
    import type { LayoutProps } from "./$types";

    let { data }: LayoutProps = $props();

    let { permissions } = data;
</script>

<div>
    {#if permissions.canAccessOps}
        <div>something that you can display from this</div>
    {/if}
</div>
