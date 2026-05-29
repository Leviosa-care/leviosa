<script lang="ts">
    import { onMount } from "svelte";
    import { goto } from "$app/navigation";

    onMount(() => {
        if (window.innerWidth >= 768) {
            goto("/admin/bookings/consultations");
        }
    });
</script>

here I have both events and consultations
