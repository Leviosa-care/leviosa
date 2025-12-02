<script lang="ts">
    import { Button } from "bits-ui";
    import {
        buttonVariants,
        type ButtonVariants,
    } from "$lib/utils/design-system";

    type Props = {
        variant?: ButtonVariants["variant"];
        size?: ButtonVariants["size"];
        disabled?: boolean;
        type?: "button" | "submit" | "reset";
        onclick?: () => void;
        class?: string;
        children?: any;
    };

    let {
        variant = "default",
        size = "default",
        disabled = false,
        type = "button",
        onclick,
        class: className = "",
        children,
        ...restProps
    }: Props = $props();
</script>

<Button.Root
    {type}
    {disabled}
    onclick={() => onclick?.()}
    class={buttonVariants({ variant, size, class: className })}
    {...restProps}
>
    {#if children}
        {@render children()}
    {:else}
        <slot />
    {/if}
</Button.Root>

