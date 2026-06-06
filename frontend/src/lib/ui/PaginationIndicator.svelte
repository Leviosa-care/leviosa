<script lang="ts">
    interface Props {
        count?: number;
        selectedElement?: number;
    }

    let { count = 1, selectedElement = 0 }: Props = $props();
</script>

<div class="items-center flex gap-2">
    {#each Array(count) as _, index}
        {@const dark = index <= selectedElement}
        <div
            class="w-8 h-2 rounded-sm {dark ? 'bg-dark' : 'bg-muted'}"
        ></div>
    {/each}
</div>
