<script lang="ts">
    import { run, createBubbler } from "svelte/legacy";
    import { onDestroy } from "svelte";
    import { browser } from "$app/environment";
    import { fly, fade } from "svelte/transition";
    import { createHorizontalSwipeHandler } from "$lib/scripts/swipe";

    interface Props {
        isOpen?: boolean;
        children?: import("svelte").Snippet;
    }

    let { children, isOpen = $bindable(false) }: Props = $props();

    let scrollPosition: number = 0;

    const bubble = function (e: MouseEvent) {
        e.stopPropagation();
        createBubbler()("click");
    };

    const onSwipe = (direction: "left" | "right") => {
        if (direction === "left") isOpen = false;
    };

    const closeSwipeAction = createHorizontalSwipeHandler(onSwipe);

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
    <!-- Render overlay first, then drawer, so drawer appears on top in DOM -->
    <div
        transition:fade={{ duration: 300 }}
        class="overlay bg-white"
        onclick={() => (isOpen = false)}
        onkeydown={handleKeydown}
        class:visible={isOpen}
        tabindex="0"
        role="button"
    ></div>
    <div
        transition:fly={{ x: -300, duration: 300 }}
        class="side-drawer"
        class:visible={isOpen}
        onclick={bubble}
        onkeydown={handleKeydown}
        use:closeSwipeAction.action
        tabindex="0"
        role="dialog"
        aria-modal="true"
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
        z-index: 100;
        pointer-events: auto;
    }
    .overlay.visible {
        opacity: 1;
        visibility: visible;
    }

    .side-drawer {
        position: fixed;
        top: 0;
        left: 0;
        height: 100vh;
        width: 280px;
        max-width: 85vw;
        background: var(--background);
        padding: 1.5rem;
        box-shadow: 2px 0 10px rgba(0, 0, 0, 0.2);
        z-index: 200;
        color: hsl(var(--foreground));
        overflow-y: auto;
        pointer-events: auto;
    }
</style>
