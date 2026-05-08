<script lang="ts">
    import Logo from "./__logo.svelte";
    import { Instagram } from "@lucide/svelte";
</script>

<div class="h-full flex justify-between items-center py-6 px-6 sm:px-12">
    <Logo />
    <div class="flex items-center gap-6">
        <a
            href="https://www.instagram.com/leviosa_care"
            target="_blank"
            rel="noopener noreferrer"
            class="text-dark-900 hover:text-dark-700 transition-colors flex items-center gap-2"
            aria-label="Follow us on Instagram"
        >
            <Instagram class="w-5 h-5" />
        </a>
        <p class="text-dark-900 font-medium text-xs sm:text-base">
            © Leviosa Care. All rights reserved
        </p>
    </div>
</div>
