<script lang="ts">
    import type { LayoutProps } from "./$types";
    import { setUserContext } from "$lib/context/user.svelte";
    import "../../app.css";

    let { data, children }: LayoutProps = $props();

    setUserContext(data.user);
</script>

<div class="flex flex-col">
    {@render children()}
</div>
