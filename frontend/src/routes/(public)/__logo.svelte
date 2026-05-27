<script lang="ts">
    interface Props {
        scrolled?: boolean;
    }
    let { scrolled = true }: Props = $props();
</script>

<a href="/" class="text-lg sm:text-2xl font-bold transition-colors duration-300 {scrolled ? 'text-dark-900' : 'text-white'}">Leviosa</a>
