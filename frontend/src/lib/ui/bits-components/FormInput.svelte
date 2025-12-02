<script lang="ts">
    import { Input } from "bits-ui";
    import { cn } from "$lib/utils/design-system";
    import {
        INPUT_VARIANTS,
        INPUT_SIZES,
        type InputVariantType,
        type InputSizeType,
    } from "./constants";

    interface Props {
        id?: string;
        name?: string;
        type?: "text" | "password" | "email" | "number" | "tel";
        placeholder?: string;
        label?: string;
        value?: string;
        disabled?: boolean;
        required?: boolean;
        error?: string;
        variant?: InputVariantType;
        size?: InputSizeType;
        class?: string;
    }

    let {
        id,
        name,
        type = "text",
        placeholder,
        label,
        value = $bindable(),
        disabled = false,
        required = false,
        error,
        variant = error ? INPUT_VARIANTS.ERROR : INPUT_VARIANTS.DEFAULT,
        size = INPUT_SIZES.DEFAULT,
        class: className = "",
        ...restProps
    }: Props = $props();

    // Generate unique ID if not provided
    id = id || `input-${Math.random().toString(36).substring(2, 9)}`;
</script>

<div class="space-y-2">
    {#if label}
        <label
            for={id}
            class="text-sm font-medium text-foreground"
            class:required
        >
            {label}
            {#if required}
                <span class="text-destructive ml-1">*</span>
            {/if}
        </label>
    {/if}

    <div class="relative">
        <Input.Root
            bind:value
            {id}
            {name}
            {type}
            {placeholder}
            {disabled}
            class={cn(
                "flex w-full rounded-input border border-input bg-background px-3 py-2 text-sm",
                "ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium",
                "placeholder:text-muted-foreground",
                "focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2",
                "disabled:cursor-not-allowed disabled:opacity-50",
                "hover:border-input-hover",
                {
                    "h-10": size === "default",
                    "h-9": size === "sm",
                    "h-11": size === "lg",
                    "border-destructive focus-visible:ring-destructive":
                        variant === "error",
                },
                className,
            )}
            {...restProps}
        />
    </div>

    {#if error}
        <p class="text-sm text-destructive" role="alert">
            {error}
        </p>
    {/if}
</div>

<style>
    .required::after {
        content: "*";
        color: hsl(var(--destructive));
        margin-left: 0.25rem;
    }
</style>
