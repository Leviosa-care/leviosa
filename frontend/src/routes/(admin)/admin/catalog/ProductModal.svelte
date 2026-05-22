<script lang="ts">
    import type { Snippet } from "svelte";

    type Input = Snippet<[string]>;
    type Select = Snippet<
        [
            string,
            Set<string>,
            { current: string; previous: string; next: string },
        ]
    >;

    import { Dialog, Label, Separator, Select, Button } from "bits-ui";
    import { type SuperValidated } from "sveltekit-superforms";
    import { superForm } from "sveltekit-superforms";

    import { X, ChevronDown, Check } from "@lucide/svelte";

    import { capitalizeFirstWord } from "$lib/utils/capitalize";
    import type { CardType, Category } from "$lib/data/mockData";
    import {
        defaultStatus,
        defaultCategory,
        defaultAvailability,
    } from "./default";
    import { formatItems } from "./formatItems";
    import { type product } from "./schemas";
    type Props = {
        children: import("svelte").Snippet;
        statuses: Set<string>;
        categories: Category[];
        availabilities: Set<string>;
        modalForm: SuperValidated<product>;
        card?: CardType;
    };
    let {
        children,
        statuses,
        categories,
        availabilities,
        modalForm = $bindable(),
        card,
    }: Props = $props();

    const {
        form,
        errors,
        enhance,
        tainted,
        // constraints: deleteConstraints,
    } = superForm(modalForm, {
        resetForm: true,
        // NOTE: just for debug, I want to check if something happens
        onUpdate: ({ form: updatedForm }) => {
            modalForm = updatedForm;
        },
        // onSubmit: ({ formData, cancel }) => {
        onSubmit: ({ formData }) => {
            if (card) {
                const modifiedFields = getModifiedFields();

                if (Object.keys(modifiedFields).length === 0) {
                    // No changes made, might want to cancel or handle this case
                    return;
                }

                // Delete all existing entries
                for (const key of formData.keys()) {
                    formData.delete(key);
                }

                Object.entries(modifiedFields).forEach(([key, value]) => {
                    formData.append(key, value);
                });

                if (card?.id) {
                    formData.append("id", card.id);
                }
            }
            console.log("FormData contents:", Object.fromEntries(formData));
        },
        onResult: ({ result }) => {
            console.log("Form submission result:", result);
        },
        onError: ({ result }) => {
            console.log("Form submission error:", result);
        },
    });

    async function urlToFile(
        url: string,
        filename: string,
        mimeType: string,
    ): Promise<File> {
        const response = await fetch(url);
        const blob = await response.blob();
        return new File([blob], filename, { type: mimeType });
    }

    $effect(() => {
        if (card) {
            (async () => {
                const file = await urlToFile(
                    card.image,
                    "image.jpg",
                    "image/jpeg",
                );

                $form = { ...card };
            })();
        }
    });

    function getModifiedFields() {
        const modifiedData: Record<string, any> = {};

        const taintedFields = $tainted;
        const formData = $form;

        if (taintedFields) {
            (
                Object.keys(taintedFields) as Array<keyof typeof formData>
            ).forEach((key) => {
                if (taintedFields[key]) {
                    modifiedData[key] = formData[key];
                }
            });
        }
        return modifiedData;
    }

    // remove values from select lists that are not used
    let modalCategories = new Set(categories.map((category) => category.name));
    modalCategories.delete(defaultCategory);
    const addType = "Ajouter un type...";
    modalCategories.add(addType);

    let modalAvailabilities = new Set(availabilities);
    modalAvailabilities.delete(defaultAvailability);

    let modalStatuses = new Set(statuses);
    modalStatuses.delete(defaultStatus);

    let statusState = $state({
        // NOTE: this is the thing that I use in the production, but here for testing it helps to have default values
        // current: card ? card.published : "Choisis un statut...",
        current: card ? card.published : $form.published,
        previous: card ? card.published : $form.published,
        next: "",
    });

    let categoryState = $state({
        current: card ? card.category : $form.category,
        previous: card ? card.category : $form.category,
        next: "",
    });

    // TODO: the thing that I need to do, once the category thing is well defined in cards
    // because I should use the category instead of the value
    // let categoryState = $state({
    //     current: card ? card.categoryID : $form.categoryID,
    //     previous: card ? card.categoryID : $form.categoryID,
    //     next: "",
    // });

    let availabilityState = $state({
        current: card ? card.availability : $form.availability,
        previous: card ? card.availability : $form.availability,
        next: "",
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
                {card ? "Modifier" : "Creer"} un produit
            </Dialog.Title>
            <Dialog.Description class="text-foreground-alt !mt-1 text-sm"
                >Remplissez les détails ci-dessous pour {card
                    ? "modifier un produit."
                    : "créer un nouveau produit."}</Dialog.Description
            >

            <Separator.Root class="bg-muted mx-5 !mb-2 !mt-5 block h-px" />
            <form
                method="POST"
                action="?/{card ? 'updateProduct' : 'createProduct'}"
                class="grid gap-4"
                use:enhance
                enctype="multipart/form-data"
            >
                <div
                    class="grid grid-cols-[max-content_1fr] gap-4 w-full items-center pb-11 pt-7"
                >
                    {@render field("image", "Image", file)}
                    {@render field("name", "Nom", input)}
                    <input
                        type="hidden"
                        name="category"
                        id="category"
                        value={categoryState.current}
                    />
                    {@render field(
                        "category",
                        "Categorie",
                        selectWithAddOption,
                    )}
                    {@render field("description", "Description", textarea)}
                    {@render field("price", "Prix", input)}
                    <input
                        type="hidden"
                        name="published"
                        id="published"
                        value={statusState.current}
                    />
                    {@render field(
                        "published",
                        "Status",
                        undefined,
                        select,
                        modalStatuses,
                        statusState,
                    )}
                    {@render field("duration", "Duree", input)}
                    <input
                        type="hidden"
                        name="availability"
                        id="availability"
                        value={availabilityState.current}
                    />
                    {@render field(
                        "availability",
                        "Disponibilite",
                        undefined,
                        select,
                        modalAvailabilities,
                        availabilityState,
                    )}
                    {@render field("bufferTime", "Temps de battement", input)}
                    {@render field(
                        "cancellationHours",
                        "Heure limite d’annulation",
                        input,
                    )}
                </div>
                <div class="flex w-full justify-end">
                    <Button.Root type="submit" class="cursor-pointer">
                        <Dialog.Close
                            class="h-input rounded-input bg-dark text-background shadow-mini hover:bg-dark/95 focus-visible:ring-dark focus-visible:ring-offset-background focus-visible:outline-hidden inline-flex items-center justify-center px-[50px] text-[15px] font-semibold focus-visible:ring-2 focus-visible:ring-offset-2 active:scale-[0.98] cursor-pointer"
                        >
                            Sauvegarder le produit
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

{#snippet field(
    name: string,
    label: string,
    input?: Input,
    select?: Select,
    options?: Set<string>,
    state?: { current: string; previous: string; next: string },
)}
    <Label.Root for={name} class="text-sm font-semibold text-right"
        >{label}</Label.Root
    >

    <div class="relative w-full">
        {#if input}
            {@render input(name)}
        {:else if select && options && state}
            {@render select(name, options, state)}
        {/if}
    </div>
{/snippet}

{#snippet input(name: string)}
    <input
        id={name}
        class="h-input rounded-card-sm border-border-input bg-background placeholder:text-foreground-alt/50 hover:border-dark-40 focus:ring-foreground focus:ring-offset-background focus:outline-hidden inline-flex w-full items-center border px-4 text-base focus:ring-2 focus:ring-offset-2 sm:text-sm"
        {name}
        bind:value={$form[name as keyof product]}
    />
{/snippet}

{#snippet selectWithAddOption(name: string)}
    {#if categoryState.current === addType}
        <!-- Input mode -->
        <div class="flex items-center gap-2">
            <!-- TODO: add the name to the input or it wont be sent properly -->
            <!-- TODO: add something for the description of that category -->
            <input
                name="newCategoryName"
                type="text"
                bind:value={categoryState.next}
                placeholder="Nouveau type..."
                class="h-input rounded-9px border-border-input bg-background inline-flex w-full items-center border px-4 text-base focus:ring-2 focus:ring-foreground focus:ring-offset-2 focus:ring-offset-background focus:outline-hidden sm:text-sm"
            />
            <button
                type="button"
                onclick={() => {
                    categoryState.current = categoryState.previous;
                    categoryState.next = "";
                }}
                class="h-input rounded-9px border-border-input bg-background hover:bg-muted inline-flex items-center justify-center border px-3 text-sm transition-colors"
                aria-label="Annuler"
            >
                Annuler
            </button>
        </div>
    {:else}
        <!-- Select mode -->
        {@render select(name, modalCategories, categoryState)}
    {/if}
{/snippet}
{#snippet select(
    name: string,
    options: Set<string>,
    state: { current: string; previous: string; next: string },
)}
    {@const items = formatItems(options)}
    <Select.Root
        type="single"
        onValueChange={(v) => {
            if (v && v !== addType) {
                state.previous = state.current;
                // TODO: should I add something for the form syncing or not ?
                // $form[name] = state.current;
            }
            state.current = v;
        }}
        {items}
    >
        <Select.Trigger
            id={name}
            class="h-input rounded-9px border-border-input bg-background inline-flex justify-between gap-4 touch-none select-none items-center border px-4 w-full transition-colors"
            aria-label="Select a product"
        >
            {capitalizeFirstWord(state.current)}
            <ChevronDown class="text-muted-foreground ml-auto size-6" />
        </Select.Trigger>
        <Select.Portal>
            <Select.Content
                class="focus-override border-muted bg-background shadow-popover data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 data-[state=closed]:zoom-out-95 data-[state=open]:zoom-in-95 data-[side=bottom]:slide-in-from-top-2 data-[side=left]:slide-in-from-right-2 data-[side=right]:slide-in-from-left-2 data-[side=top]:slide-in-from-bottom-2 outline-hidden z-50 h-full max-h-[var(--bits-select-content-available-height)] w-[var(--bits-select-anchor-width)] min-w-[var(--bits-select-anchor-width)] select-none rounded-xl border px-1 py-3 data-[side=bottom]:translate-y-1 data-[side=left]:-translate-x-1 data-[side=right]:translate-x-1 data-[side=top]:-translate-y-1"
                sideOffset={10}
            >
                <Select.Viewport class="p-1">
                    {#each items as { value, label }, i (i + value)}
                        <Select.Item
                            class="rounded-button data-highlighted:bg-muted outline-hidden data-disabled:opacity-50 flex h-10 w-full select-none items-center py-3 pl-5 pr-1.5 text-sm capitalize"
                            {value}
                            {label}
                        >
                            {#snippet children({ selected })}
                                {value}
                                {#if selected}
                                    <div class="ml-auto">
                                        <Check aria-label="check" />
                                    </div>
                                {/if}
                            {/snippet}
                        </Select.Item>
                    {/each}
                </Select.Viewport>
            </Select.Content>
        </Select.Portal>
    </Select.Root>
{/snippet}

{#snippet textarea(name: string)}
    <textarea
        id={name}
        {name}
        rows="5"
        class="rounded-card-sm border border-border-input bg-background placeholder:text-foreground-alt/50 hover:border-dark-40 focus:ring-foreground focus:ring-offset-background focus:outline-hidden w-full px-4 py-2 text-base focus:ring-2 focus:ring-offset-2 sm:text-sm"
        placeholder="Ecrivez une description pour le produit"
        bind:value={$form[name as keyof product]}
    ></textarea>
{/snippet}

{#snippet file(name: string)}
    <input
        type="file"
        id={name}
        class="h-input rounded-card-sm border-border-input bg-background hover:border-dark-40 inline-flex w-full items-center border text-base sm:text-sm file:mr-4 file:px-4 file:h-full file:rounded file:border-0 file:text-sm file:font-medium file:bg-foreground file:text-background hover:file:bg-foreground/90 cursor-pointer"
        {name}
        accept="image/png, image/jpeg"
        bind:value={$form[name as keyof product]}
    />
{/snippet}
