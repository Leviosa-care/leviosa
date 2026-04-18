<script lang="ts">
    import { onMount } from "svelte";

    import { goto } from "$app/navigation";
    import { Button } from "bits-ui";

    let isLargeScreen = $state(false);
    import { icons, type Icon } from "./icons";

    onMount(() => {
        const mediaQuery = window.matchMedia("(min-width: 768px)"); // Adjust breakpoint as needed
        // Set initial value
        isLargeScreen = mediaQuery.matches;

        // Listen for changes
        const handleChange = (e) => {
            isLargeScreen = e.matches;
        };

        mediaQuery.addEventListener("change", handleChange);

        // Cleanup
        return () => {
            mediaQuery.removeEventListener("change", handleChange);
        };
    });
    let currentIcons = $derived(isLargeScreen ? icons.large : icons.small);
    // let selected = $state<Icon>(currentIcons[0]);
    let selected = $state<Icon | null>(null);

    // Update selected when currentIcons changes, but preserve selection if it exists in new array
    $effect(() => {
        if (currentIcons.length > 0) {
            // If no selection or current selection not in new array, select first item
            if (
                !selected ||
                !currentIcons.find((icon) => icon.link === selected?.link)
            ) {
                selected = currentIcons[0];
            }
        } else {
            selected = null;
        }
    });
</script>

<div class="flex flex-col">
    {#each currentIcons as icon}
        {@render navIcon(icon)}
    {/each}
</div>

{#snippet navIcon(icon: Icon)}
    {@const I = icon.icon}
    <Button.Root
        type="button"
        onclick={() => {
            selected = icon;
            goto(icon.link);
        }}
        class="flex rounded-md gap-4 p-3 cursor-pointer font-md {selected?.name ===
        icon.name
            ? 'bg-dark text-white'
            : 'bg-transparent hover:bg-gray-200'}
        "
    >
        <I />
        <p>{icon.name}</p>
    </Button.Root>
{/snippet}
