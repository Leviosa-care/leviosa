<script lang="ts">
    import type { LayoutProps } from "./$types";
    import Sidebar from "./sidebar.svelte";

    let { children, data }: LayoutProps = $props();
</script>

<!-- Main layout container with sidebar and content area -->
<div class="flex min-h-screen">
    <!-- Sidebar component handles both desktop sidebar and mobile bottom nav -->
    <Sidebar {data} />

    <!-- Main content area -->
    <main class="flex-1 overflow-x-hidden">
        <!-- Add bottom padding on mobile to prevent content from being hidden behind bottom nav -->
        <div class="pb-20 lg:pb-0">
            {@render children()}
        </div>
    </main>
</div>
