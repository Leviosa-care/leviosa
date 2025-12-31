<script lang="ts">
    import type { LayoutProps } from "./$types";
    let { children }: LayoutProps = $props();

    // NOTE: Here we render the children + sidebar that is going to be using our defined
    // permissions matrix to display the right elements
</script>

<div class="flex">
    <div id="sidebar"></div>
    <div class="flex-1">
        {@render children()}
    </div>
</div>
