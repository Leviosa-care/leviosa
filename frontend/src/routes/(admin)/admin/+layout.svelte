<script lang="ts">
    import type { LayoutProps } from "./$types";
    import { setUserContext } from "$lib/context/user.svelte";
    import NavigationBar from "./NavigationBar.svelte";

    import { Slack } from "@lucide/svelte";

    let { data, children }: LayoutProps = $props();

    setUserContext(data.user);
    // TODO: make that thing work brother
</script>

<div class="flex flex-col md:flex-row-reverse">
    {@render children()}
    <div class="flex flex-col justify-between px-4 py-8 border-r-1">
        <div class="grid gap-8">
            <div class="flex items-center gap-4">
                <Slack />
                <p class="font-bold text-xl">Leviosa</p>
            </div>
            <NavigationBar />
        </div>
        <p>the bottom navigation</p>
    </div>
</div>
