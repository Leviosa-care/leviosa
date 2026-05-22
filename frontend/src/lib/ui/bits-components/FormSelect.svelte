<script lang="ts">
    import { Select } from "bits-ui";
    import { cn } from "$lib/utils/design-system";
    import { ChevronDown } from "@lucide/svelte";
    import {
        INPUT_VARIANTS,
        INPUT_SIZES,
        type InputVariantType,
        type InputSizeType,
    } from "./constants";

    interface Option {
        value: string;
        label: string;
        disabled?: boolean;
    }

    interface Props {
        id?: string;
        name?: string;
        label?: string;
        value?: string;
        options: Option[];
        placeholder?: string;
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
        label,
        value = $bindable(),
        options = [],
        placeholder = "Select an option",
        disabled = false,
        required = false,
        error,
        variant = error ? INPUT_VARIANTS.ERROR : INPUT_VARIANTS.DEFAULT,
        size = INPUT_SIZES.DEFAULT,
        class: className = "",
        ...restProps
    }: Props = $props();

    // Generate unique ID if not provided
    id = id || `select-${Math.random().toString(36).substring(2, 9)}`;

    let valueBinding = $state(value);

    // Sync valueBinding with the bindable value
    $effect(() => {
        value = valueBinding;
    });

    // Get the selected option label for display
    let selectedLabel = $derived(options.find(opt => opt.value === valueBinding)?.label ?? "");
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

    <Select.Root type="single" value={valueBinding} onValueChange={(v: string) => { valueBinding = v; }} {disabled}>
        <Select.Trigger
            {id}
            {name}
            class={cn(
                "flex w-full items-center justify-between rounded-input border border-input bg-background px-3 py-2 text-sm",
                "ring-offset-background placeholder:text-muted-foreground",
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
        >
            <span class="flex-1 truncate">
                {#if selectedLabel}
                    {selectedLabel}
                {:else}
                    <span class="text-muted-foreground">{placeholder}</span>
                {/if}
            </span>
            <ChevronDown class="h-4 w-4 opacity-50 pointer-events-none" />
        </Select.Trigger>
        <Select.Portal>
            <Select.Content
                class="relative z-50 min-w-[8rem] overflow-hidden rounded-card border border-card bg-background text-foreground shadow-popover data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 data-[state=closed]:zoom-out-95 data-[state=open]:zoom-in-95 data-[side=bottom]:slide-in-from-top-2 data-[side=left]:slide-in-from-right-2 data-[side=right]:slide-in-from-left-2 data-[side=top]:slide-in-from-bottom-2"
                sideOffset={4}
            >
                <Select.Viewport class="p-1">
                    <Select.Group>
                        {#each options as option}
                            <Select.Item
                                value={option.value}
                                disabled={option.disabled}
                                class="relative flex w-full cursor-default select-none items-center rounded-sm py-1.5 pl-2 pr-8 text-sm outline-none focus:bg-accent focus:text-accent-foreground data-[disabled]:pointer-events-none data-[disabled]:opacity-50"
                            >
                                <span class="flex-1 truncate">
                                    {option.label}
                                </span>
                            </Select.Item>
                        {/each}
                    </Select.Group>
                </Select.Viewport>
            </Select.Content>
        </Select.Portal>
    </Select.Root>

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
