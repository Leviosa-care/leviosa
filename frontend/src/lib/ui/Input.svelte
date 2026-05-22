<script lang="ts">
    import { goto } from "$app/navigation";
    import { Button } from "bits-ui";
    import { Eye, EyeOff } from "@lucide/svelte";

    import { capitalizeFirstWord } from "$lib/utils/capitalize";

    type Props = {
        name: string;
        label?: string;
        type?: "text" | "email" | "password";
        value: string;
        error: string;
        showForgotPasswordbutton?: boolean;
        autocomplete?: string | null | undefined;
    };
    let {
        name,
        label = name,
        type = "text",
        value = $bindable(""),
        error = $bindable(""),
        showForgotPasswordbutton = false,
        autocomplete,
    }: Props = $props();

    let invalid = $derived(error != "");
    let isFocused = $state(false);
    let visible = $state(false);

    let resolvedInputType = $derived.by(() => {
        if (type === "password") return visible ? "text" : "password";
        else return type;
    });
</script>

<div class="relative w-full">
    <div class="flex justify-right left-auto mb-2">
        {#if type === "password" && showForgotPasswordbutton}
            <Button.Root
                class="ml-auto hover:underline text-sm text-semibold cursor-pointer"
                type="button"
                onclick={() => goto("/auth/forgotten-password")}
            >
                {capitalizeFirstWord("mot de passe oublié ?")}
            </Button.Root>
        {/if}
    </div>
    <div
        class="relative w-full overflow-hidden rounded-xl focus-within:ring-offset-2 focus-within:ring-2 focus-within:ring-foreground"
    >
        <input
            id={name}
            type={resolvedInputType}
            {name}
            autocomplete={autocomplete as any}
            bind:value
            aria-invalid={invalid}
            onfocus={() => (isFocused = true)}
            onblur={() => (isFocused = value.length === 0 ? false : true)}
            class={`w-full px-4 pr-11 pb-3 pt-6 border border-gray-300 rounded-xl peer ${
                invalid ? "border-red-500" : "border-border-input"
            } bg-background hover:border-dark-40 focus:ring-foreground focus:ring-offset-background focus:outline-hidden inline-flex w-full items-center border text-base focus:ring-2 focus:ring-offset-2 sm:text-lg text-black font-base`}
        />
        <label
            for={name}
            class={`
        absolute left-4 top-1/2 transition-all duration-200 ease-in-out origin-top-left
        ${
            isFocused || value.length > 0
                ? "translate-y-[-1.6125rem] scale-80 text-gray-400 "
                : `translate-y-[-50%] scale-100 font-medium text-gray-300`
        }
      `}
        >
            {capitalizeFirstWord(label)}
        </label>
        {#if type === "password"}
            <button
                type="button"
                class="absolute pr-3 top-1/2 right-0 translate-y-[-50%] cursor-pointer"
                onclick={() => (visible = !visible)}
            >
                {#if visible}
                    <EyeOff strokeWidth={1} absoluteStrokeWidth={true} />
                {:else}
                    <Eye strokeWidth={1} absoluteStrokeWidth={true} />
                {/if}
            </button>
        {/if}
    </div>
    {#if invalid}
        <p class="ml-2 mt-1 text-red-500 text-sm">
            {error}
        </p>
    {/if}
</div>
