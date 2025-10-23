<script lang="ts">
    import BackButton from "$lib/navigation/BackButton.svelte";
    import PaginationIndicator from "$lib/ui/PaginationIndicator.svelte";

    interface Props {
        selectedElement: number;
        pathname?: string;
    }

    let { selectedElement, pathname = "" }: Props = $props();
</script>

<div class="flex justify-between items-center py-4">
    <BackButton pathname={pathname !== "" ? pathname : "/"} />
    <PaginationIndicator count={4} {selectedElement} />
    <div class="w-8 aspect-square invisible"></div>
</div>

<style>
</style>
