<script lang="ts">
    import { onMount } from "svelte";
    import { goto } from "$app/navigation";
    onMount(() => {
        if (window.innerWidth >= 768) {
            // Use any logic to choose default, e.g. based on previous visit, preference, or default to events
            goto("/bookings");
        }
    });
</script>

here I am in the events handling page
