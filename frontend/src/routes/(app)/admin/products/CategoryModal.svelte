<script lang="ts">
    import type { Snippet } from "svelte";
    import { Dialog, Label, Separator, Button } from "bits-ui";

    type Input = Snippet<[string]>;
    import { type SuperValidated } from "sveltekit-superforms";
    import { superForm } from "sveltekit-superforms";

    import { X } from "@lucide/svelte";

    // import { capitalizeFirstWord } from "$lib/utils/capitalize";
    // import type { Category } from "./products";
    import { type category } from "./schemas";
    type Props = {
        children: import("svelte").Snippet;
        modalForm: SuperValidated<category>;
        // categories: Category[];
        // category?: Category;
    };
    let {
        children,
        // categories,
        modalForm = $bindable(),
        // category,
    }: Props = $props();

    const {
        form,
        errors,
        enhance,
        // constraints: deleteConstraints,
    } = superForm(modalForm, {
        resetForm: true,
    });

    // TODO: here just use the superdebug to understand what is going on now ?
    // import SuperDebug from "sveltekit-superforms";

    // TODO: something to take into account is that if there is no image, maybe we should not publish the product ?
</script>

<!-- <SuperDebug -->
<!--     data={{ -->
<!--         form: $form, -->
<!--         errors: $errors, -->
<!--     }} -->
<!-- /> -->

<Dialog.Root>
    <Dialog.Trigger>
        {@render children()}
    </Dialog.Trigger>
    <Dialog.Portal>
        <Dialog.Overlay
            class="data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 fixed inset-0 z-50 bg-black/80"
        />
        <Dialog.Content
            class="rounded-card-lg bg-background shadow-popover data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 data-[state=closed]:zoom-out-95 data-[state=open]:zoom-in-95 outline-hidden fixed left-[50%] top-[50%] z-50 w-full max-w-[calc(100%-2rem)] translate-x-[-50%] translate-y-[-50%] border p-8 sm:max-w-[540px] md:w-full"
        >
            <Dialog.Title class="w-full text-xl font-semibold tracking-tight">
                Creer une categorie
            </Dialog.Title>
            <Dialog.Description class="text-foreground-alt !mt-1 text-sm"
                >Remplissez les détails ci-dessous pour créer une nouvelle
                categorie.</Dialog.Description
            >

            <Separator.Root class="bg-muted mx-5 !mb-2 !mt-5 block h-px" />
            <form
                method="POST"
                action="?/createCategory"
                class="grid gap-4"
                use:enhance
                enctype="multipart/form-data"
            >
                <div
                    class="grid grid-cols-[max-content_1fr] gap-4 w-full items-center pb-11 pt-7"
                >
                    <!-- {@render field("image", "Image", file)} -->
                    {@render field("name", "Nom", input)}
                    {@render field("description", "Description", textarea)}
                </div>
                <div class="flex w-full justify-end">
                    <Button.Root type="submit" class="cursor-pointer">
                        <Dialog.Close
                            class="h-input rounded-input bg-dark text-background shadow-mini hover:bg-dark/95 focus-visible:ring-dark focus-visible:ring-offset-background focus-visible:outline-hidden inline-flex items-center justify-center px-[50px] text-[15px] font-semibold focus-visible:ring-2 focus-visible:ring-offset-2 active:scale-[0.98] cursor-pointer"
                        >
                            Sauvegarder la categorie
                        </Dialog.Close>
                    </Button.Root>
                </div>
            </form>
            <Button.Root type="button" class="cursor-pointer">
                <Dialog.Close
                    class="focus-visible:ring-foreground focus-visible:ring-offset-background focus-visible:outline-hidden absolute right-5 top-5 rounded-md focus-visible:ring-2 focus-visible:ring-offset-2 active:scale-[0.98] cursor-pointer"
                >
                    <X class="text-foreground size-5" />
                    <span class="sr-only">Close</span>
                </Dialog.Close>
            </Button.Root>
        </Dialog.Content>
    </Dialog.Portal>
</Dialog.Root>

{#snippet field(name: string, label: string, input: Input)}
    <Label.Root for={name} class="text-sm font-semibold text-right"
        >{label}</Label.Root
    >

    <div class="relative w-full">
        {@render input(name)}
    </div>
{/snippet}

{#snippet input(name: string)}
    <input
        id={name}
        class="h-input rounded-card-sm border-border-input bg-background placeholder:text-foreground-alt/50 hover:border-dark-40 focus:ring-foreground focus:ring-offset-background focus:outline-hidden inline-flex w-full items-center border px-4 text-base focus:ring-2 focus:ring-offset-2 sm:text-sm"
        {name}
        bind:value={$form[name as keyof category]}
    />
{/snippet}

{#snippet textarea(name: string)}
    <textarea
        id={name}
        {name}
        rows="5"
        class="rounded-card-sm border border-border-input bg-background placeholder:text-foreground-alt/50 hover:border-dark-40 focus:ring-foreground focus:ring-offset-background focus:outline-hidden w-full px-4 py-2 text-base focus:ring-2 focus:ring-offset-2 sm:text-sm"
        placeholder="Ecrivez une description pour le produit"
        bind:value={$form[name as keyof category]}
    ></textarea>
{/snippet}

<!-- {#snippet file(name: string)} -->
<!--     <input -->
<!--         type="file" -->
<!--         id={name} -->
<!--         class="h-input rounded-card-sm border-border-input bg-background hover:border-dark-40 inline-flex w-full items-center border text-base sm:text-sm file:mr-4 file:px-4 file:h-full file:rounded file:border-0 file:text-sm file:font-medium file:bg-foreground file:text-background hover:file:bg-foreground/90 cursor-pointer" -->
<!--         {name} -->
<!--         accept="image/png, image/jpeg" -->
<!--         bind:value={$form[name as keyof product]} -->
<!--     /> -->
<!-- {/snippet} -->
