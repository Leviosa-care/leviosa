<script lang="ts">
    import { Button } from "bits-ui";
    import { cn } from "$lib/utils/design-system";
    import { ArrowLeft, ChevronLeft } from "@lucide/svelte";
    import {
        type ButtonSizeType,
        type ButtonVariantType,
        BUTTON_SIZES,
        BUTTON_VARIANTS,
    } from "./constants";

    interface Props {
        href?: string;
        onclick?: () => void;
        label?: string;
        variant?: ButtonVariantType;
        size?: ButtonSizeType;
        disabled?: boolean;
        class?: string;
    }

    let {
        href,
        onclick,
        label = "Go back",
        variant = BUTTON_VARIANTS.CHEVRON,
        size = BUTTON_SIZES.DEFAULT,
        disabled = false,
        class: className = "",
        ...restProps
    }: Props = $props();

    function handleClick() {
        if (href) {
            window.history.back();
        } else if (onclick) {
            onclick();
        }
    }
</script>

<Button.Root
    {disabled}
    onclick={handleClick}
    aria-label={label}
    class={cn("transition-all duration-200 p-3 cursor-pointer", className)}
    {...restProps}
>
    {#if variant === BUTTON_VARIANTS.CHEVRON}
        <ChevronLeft />
    {:else if variant === BUTTON_VARIANTS.ARROW}
        <ArrowLeft />
    {/if}
    {#if size !== "icon"}
        <span class="ml-2">{label}</span>
    {/if}
</Button.Root>
