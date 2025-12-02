<script lang="ts">
    import { cn } from "$lib/utils/design-system";
    // import type { PageProps } from "./$types";
    import {
        CARD_VARIANTS,
        CARD_SIZES,
        type CardVariantType,
        type CardSizeType,
    } from "./constants";

    interface Props {
        variant?: CardVariantType;
        size?: CardSizeType;
        interactive?: boolean;
        class?: string;
        // children?: PageProps;
        onclick?: () => void;
    }

    let {
        variant = CARD_VARIANTS.DEFAULT,
        size = CARD_SIZES.DEFAULT,
        interactive = false,
        class: className = "",
        // children,
        onclick,
        ...restProps
    }: Props = $props();

    const baseClasses = [
        "rounded-card border border-card bg-background text-foreground",
        "transition-all duration-200",
        "focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2",
    ];

    const variantClasses = {
        [CARD_VARIANTS.DEFAULT]: "shadow-card",
        [CARD_VARIANTS.ELEVATED]: "shadow-popover",
        [CARD_VARIANTS.OUTLINE]: "border-border",
    };

    const sizeClasses = {
        [CARD_SIZES.DEFAULT]: "p-6",
        [CARD_SIZES.SM]: "p-4",
        [CARD_SIZES.LG]: "p-8",
    };

    const interactiveClasses = interactive
        ? [
              "cursor-pointer hover:shadow-popover hover:scale-[1.02]",
              "active:scale-[0.98]",
          ]
        : [];
</script>

<button
    class={cn(
        baseClasses,
        variantClasses[variant],
        sizeClasses[size],
        interactiveClasses,
        className,
    )}
    onclick={() => onclick?.()}
    role={interactive ? "button" : undefined}
    disabled={!interactive}
    tabindex={interactive ? 0 : undefined}
    {...restProps}
>
    <!-- {@render children()} -->
    <div>Here is the content of the card that needs to be displayed</div>
</button>
