<script lang="ts">
    import { run, createBubbler } from "svelte/legacy";
    import { onDestroy } from "svelte";
    import { browser } from "$app/environment";
    import { slide, fade } from "svelte/transition";

    const bubble = function (e: MouseEvent) {
        e.stopPropagation();
        createBubbler()("click");
    };

    let scrollPosition: number = 0;

    // =======================
    // Store and State Imports
    // =======================
    import { createVerticalSwipeHandler } from "$lib/scripts/swipe";
    interface Props {
        isOpen?: boolean;
        children?: import("svelte").Snippet;
    }

    let { children, isOpen = $bindable(false) }: Props = $props();

    // =======================
    // Helper Functions
    // =======================
    const onSwipe = (direction: "top" | "bottom") => {
        if (direction === "bottom") isOpen = false;
    };

    const closeSwipeAction = createVerticalSwipeHandler(onSwipe);
    function handleKeydown(event: KeyboardEvent) {
        event.stopPropagation();
        if (event.key === "Enter" || event.key === " ") {
            event.preventDefault();
            isOpen = false;
        }
    }

    function toggleBodyScroll(isOpen: boolean) {
        if (!browser) return;
        if (isOpen) {
            // save the current position
            scrollPosition = window.scrollY;
            document.body.style.position = "fixed";
            document.body.style.top = `-${scrollPosition}px`;
            document.body.style.width = "100%";
        } else {
            document.body.style.position = "";
            document.body.style.top = "";
            document.body.style.width = "";
            window.scrollTo(0, scrollPosition);
        }
    }

    $effect(() => {
        toggleBodyScroll(isOpen);
    });

    onDestroy(() => {
        if (!browser) return;
        document.body.style.position = "";
        document.body.style.top = "";
        window.scrollTo(0, scrollPosition);
    });
</script>

{#if isOpen}
    <div
        transition:fade={{ duration: 300 }}
        class="overlay"
        onclick={() => (isOpen = false)}
        onkeydown={handleKeydown}
        class:visible={isOpen}
        tabindex="0"
        role="button"
    ></div>
    <div
        transition:slide={{ duration: 300 }}
        class="drawer bg-white"
        class:visible={isOpen}
        onclick={bubble}
        onkeydown={handleKeydown}
        use:closeSwipeAction.action
        tabindex="0"
        role="button"
    >
        {@render children?.()}
    </div>
{/if}

<style>
    .overlay {
        position: fixed;
        top: 0;
        left: 0;
        height: 100%;
        width: 100vw;
        background: rgba(0, 0, 0, 0.2);
        opacity: 0;
        visibility: hidden;
    }
    .overlay.visible {
        opacity: 1;
        visibility: visible;
    }

    .drawer {
        --border-top-radius: 1.2rem;
        position: fixed;
        bottom: -100%;
        left: 0;
        width: 100%;
        max-height: 85vh;
        overflow-y: auto;
        padding-inline: 1rem;
        padding-bottom: 1rem;
        box-shadow: 0 -1px 10px rgba(0, 0, 0, 0.2);
        border-top-left-radius: var(--border-top-radius);
        border-top-right-radius: var(--border-top-radius);
        /* TODO: should the drawer be on top of the navigation, it makes sense but it is weird right? */
        z-index: 9999;
        color: hsl(var(--dark-900));
    }
    .drawer.visible {
        bottom: 0;
    }
</style>
