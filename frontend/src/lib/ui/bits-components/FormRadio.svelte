<script lang="ts">
    import { RadioGroup } from "bits-ui";
    import { cn } from "$lib/utils/design-system";

    interface Option {
        value: string;
        label: string;
        disabled?: boolean;
    }

    interface Props {
        name?: string;
        label?: string;
        value?: string;
        options: Option[];
        disabled?: boolean;
        required?: boolean;
        error?: string;
        variant?: "default" | "error";
        orientation?: "horizontal" | "vertical";
        class?: string;
    }

    let {
        name,
        label,
        value = $bindable(),
        options = [],
        disabled = false,
        required = false,
        error,
        variant = error ? "error" : "default",
        orientation = "vertical",
        class: className = "",
        ...restProps
    }: Props = $props();

    // Generate unique ID if not provided
    name = name || `radio-${Math.random().toString(36).substring(2, 9)}`;

    let valueBinding = $state(value);

    // Sync valueBinding with the bindable value
    $effect(() => {
        value = valueBinding;
    });
</script>

<div class="space-y-3">
    {#if label}
        <label class="text-sm font-medium text-foreground" class:required>
            {label}
            {#if required}
                <span class="text-destructive ml-1">*</span>
            {/if}
        </label>
    {/if}

    <RadioGroup.Root
        bind:value={valueBinding}
        {name}
        {disabled}
        {orientation}
        class={cn(
            "space-y-3",
            orientation === "horizontal" &&
                "space-x-6 space-y-0 flex items-center",
            className,
        )}
        {...restProps}
    >
        {#each options as option}
            <div class="flex items-center space-x-2">
                <RadioGroup.Item
                    value={option.value}
                    disabled={option.disabled}
                    id={`${name}-${option.value}`}
                    class="aspect-square h-4 w-4 rounded-full border border-input text-primary ring-offset-background focus:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50 data-[state=checked]:bg-primary data-[state=checked]:border-primary relative"
                >
                    <div class="absolute inset-0 flex items-center justify-center">
                        <div class="h-2 w-2 rounded-full bg-foreground opacity-0 data-[state=checked]:opacity-100" />
                    </div>
                </RadioGroup.Item>
                <label
                    for={`${name}-${option.value}`}
                    class={cn(
                        "text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70",
                        option.disabled && "cursor-not-allowed opacity-50",
                        variant === "error" && "text-destructive",
                    )}
                >
                    {option.label}
                </label>
            </div>
        {/each}
    </RadioGroup.Root>

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
