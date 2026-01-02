<script lang="ts">
    import { Dialog, Button, Label, Separator } from "bits-ui";
    import {
        Plus,
        Pencil,
        Trash2,
        X,
        Image as ImageIcon,
    } from "@lucide/svelte";
    import type { Snippet } from "svelte";
    import { superForm } from "sveltekit-superforms";
    import type { PageData } from "./$types";
    import Drawer from "$lib/ui/Drawer.svelte";

    import { browser } from "$app/environment";

    // Detect if we're on mobile
    let isMobile = $state(false);

    if (browser) {
        isMobile = window.innerWidth < 768;
        window.addEventListener("resize", () => {
            isMobile = window.innerWidth < 768;
        });
    }

    type Input = Snippet<[string]>;

    // Category type from API
    type Category = {
        id: string;
        name: string;
        description: string;
        status: "draft" | "published" | "archived";
        metadata?: Record<string, any>;
        created_at: string;
        updated_at: string;
        images?: Array<{
            id: string;
            parent_id: string;
            parent_type: string;
            url: string;
            title: string;
            is_active: boolean;
            created_at: string;
        }>;
    };

    interface Props {
        data: PageData;
    }

    let { data }: Props = $props();

    // Extract categories from page data
    let categories: Category[] = data.categories || [];

    // Dialog states
    let createDialogOpen = $state(false);
    let editDialogOpen = $state(false);
    let deleteDialogOpen = $state(false);

    // Currently selected category for edit/delete
    let selectedCategory: Category | null = $state(null);

    // Superforms for create, update, delete
    const {
        form: createForm,
        errors: createErrors,
        enhance: createEnhance,
    } = superForm(data.createCategoryForm, {
        resetForm: true,
        onUpdated({ form }) {
            if (form.valid) {
                createDialogOpen = false;
            }
        },
    });

    const {
        form: updateForm,
        errors: updateErrors,
        enhance: updateEnhance,
    } = superForm(data.updateCategoryForm, {
        resetForm: false,
        onUpdated({ form }) {
            if (form.valid) {
                editDialogOpen = false;
            }
        },
    });

    const {
        form: deleteForm,
        errors: deleteErrors,
        enhance: deleteEnhance,
    } = superForm(data.deleteCategoryForm, {
        resetForm: true,
        onUpdated({ form }) {
            if (form.valid) {
                deleteDialogOpen = false;
            }
        },
    });

    function openCreateDialog() {
        createDialogOpen = true;
    }

    function openEditDialog(category: Category) {
        selectedCategory = category;
        $updateForm.id = category.id;
        $updateForm.name = category.name;
        $updateForm.description = category.description;
        $updateForm.status = category.status;
        editDialogOpen = true;
    }

    function openDeleteDialog(category: Category) {
        selectedCategory = category;
        $deleteForm.id = category.id;
        deleteDialogOpen = true;
    }

    function getActiveImage(category: Category) {
        return category.images?.find((img) => img.is_active);
    }
</script>

<div class="h-full bg-white">
    <!-- Create button - floating on mobile -->
    <Button.Root
        type="button"
        class="cursor-pointer fixed bottom-20 right-4 md:absolute md:top-6 md:right-8 z-10"
        onclick={openCreateDialog}
    >
        <div
            class="flex gap-2 items-center py-2 bg-dark text-white rounded-input hover:bg-dark/95 transition-all shadow-mini md:px-4 w-12 h-12 md:w-auto md:h-auto justify-center"
        >
            <Plus size={18} />
            <span class="text-sm font-medium hidden md:inline"
                >Nouvelle Catégorie</span
            >
        </div>
    </Button.Root>

    <!-- Categories Grid -->
    <div class="p-4 md:p-8">
        {#if categories.length === 0}
            <div
                class="flex flex-col items-center justify-center py-16 text-center"
            >
                <div
                    class="w-16 h-16 rounded-full bg-dark-04 flex items-center justify-center mb-4"
                >
                    <ImageIcon size={32} class="text-dark-400" />
                </div>
                <h3 class="text-lg font-medium mb-2">Aucune catégorie</h3>
                <p class="text-sm text-foreground-alt mb-6 max-w-sm">
                    Commencez par créer votre première catégorie de service pour
                    organiser vos produits.
                </p>
                <Button.Root
                    type="button"
                    class="cursor-pointer"
                    onclick={openCreateDialog}
                >
                    <div
                        class="flex gap-2 items-center py-2 px-4 bg-dark text-white rounded-input hover:bg-dark/95 transition-all shadow-mini"
                    >
                        <Plus size={18} />
                        <span class="text-sm font-medium"
                            >Créer une catégorie</span
                        >
                    </div>
                </Button.Root>
            </div>
        {:else}
            <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                {#each categories as category (category.id)}
                    {@const activeImage = getActiveImage(category)}
                    <div
                        class="border border-border-card rounded-card bg-background shadow-card hover:shadow-popover transition-all overflow-hidden"
                    >
                        <!-- Image -->
                        <div class="aspect-video bg-dark-04 relative">
                            {#if activeImage}
                                <img
                                    src={activeImage.url}
                                    alt={activeImage.title}
                                    class="w-full h-full object-cover"
                                />
                            {:else}
                                <div
                                    class="w-full h-full flex items-center justify-center"
                                >
                                    <ImageIcon
                                        size={48}
                                        class="text-dark-300"
                                    />
                                </div>
                            {/if}
                            <!-- Status Badge -->
                            <div class="absolute top-3 right-3">
                                <span
                                    class="px-2 py-1 text-xs font-medium rounded-md {category.status ===
                                    'published'
                                        ? 'bg-green-100 text-green-800'
                                        : category.status === 'draft'
                                          ? 'bg-yellow-100 text-yellow-800'
                                          : 'bg-gray-100 text-gray-800'}"
                                >
                                    {category.status === "published"
                                        ? "Publié"
                                        : category.status === "draft"
                                          ? "Brouillon"
                                          : "Archivé"}
                                </span>
                            </div>
                        </div>

                        <!-- Content -->
                        <div class="p-4">
                            <h3 class="text-lg font-semibold mb-2">
                                {category.name}
                            </h3>
                            <p
                                class="text-sm text-foreground-alt line-clamp-2 mb-4"
                            >
                                {category.description}
                            </p>

                            <!-- Actions -->
                            <div class="flex gap-2">
                                <Button.Root
                                    type="button"
                                    class="cursor-pointer flex-1"
                                    onclick={() => openEditDialog(category)}
                                >
                                    <div
                                        class="flex gap-2 items-center justify-center py-2 px-3 border border-border-input rounded-input hover:bg-dark-04 transition-all text-sm font-medium"
                                    >
                                        <Pencil size={14} />
                                        <span>Modifier</span>
                                    </div>
                                </Button.Root>
                                <Button.Root
                                    type="button"
                                    class="cursor-pointer"
                                    onclick={() => openDeleteDialog(category)}
                                >
                                    <div
                                        class="flex items-center justify-center py-2 px-3 border border-destructive/20 text-destructive rounded-input hover:bg-destructive/10 transition-all"
                                    >
                                        <Trash2 size={14} />
                                    </div>
                                </Button.Root>
                            </div>
                        </div>
                    </div>
                {/each}
            </div>
        {/if}
    </div>
</div>

<!-- Create Category Dialog/Drawer -->
{#if isMobile}
    <Drawer bind:isOpen={createDialogOpen}>
        <div class="flex items-center justify-between mb-4">
            <h2 class="text-xl font-semibold tracking-tight">
                Créer une catégorie
            </h2>
            <button
                type="button"
                onclick={() => (createDialogOpen = false)}
                class="p-2 hover:bg-dark-04 rounded-md transition-all"
            >
                <X class="text-foreground size-5" />
            </button>
        </div>
        <p class="text-foreground-alt text-sm mb-6">
            Remplissez les détails ci-dessous pour créer une nouvelle catégorie.
        </p>

        <form
            method="POST"
            action="?/createCategory"
            enctype="multipart/form-data"
            use:createEnhance
            class="grid gap-4"
        >
            {#if $createErrors._errors}
                <div class="text-sm text-destructive mt-4">
                    {$createErrors._errors[0]}
                </div>
            {/if}

            <div class="grid grid-cols-1 gap-4 w-full pb-4">
                {@render field(
                    "name",
                    "Nom",
                    input,
                    $createForm.name,
                    $createErrors.name,
                )}
                {@render field(
                    "description",
                    "Description",
                    textarea,
                    $createForm.description,
                    $createErrors.description,
                )}
                {@render field("image", "Image", fileInput, null, null)}
            </div>
            <div class="flex w-full justify-end gap-3">
                <button
                    type="button"
                    onclick={() => (createDialogOpen = false)}
                    class="h-input rounded-input border border-border-input hover:bg-dark-04 focus-visible:ring-dark focus-visible:ring-offset-background focus-visible:outline-hidden inline-flex items-center justify-center px-6 text-sm font-medium focus-visible:ring-2 focus-visible:ring-offset-2 active:scale-[0.98] cursor-pointer transition-all"
                >
                    Annuler
                </button>
                <button
                    type="submit"
                    class="h-input rounded-input bg-dark text-background shadow-mini hover:bg-dark/95 focus-visible:ring-dark focus-visible:ring-offset-background focus-visible:outline-hidden inline-flex items-center justify-center px-8 text-sm font-semibold focus-visible:ring-2 focus-visible:ring-offset-2 active:scale-[0.98] cursor-pointer transition-all"
                >
                    Créer
                </button>
            </div>
        </form>
    </Drawer>
{:else}
    <Dialog.Root bind:open={createDialogOpen}>
        <Dialog.Portal>
            <Dialog.Overlay
                class="data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 fixed inset-0 z-50 bg-black/80"
            />
            <Dialog.Content
                class="rounded-card-lg bg-background shadow-popover data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 data-[state=closed]:zoom-out-95 data-[state=open]:zoom-in-95 outline-hidden fixed left-[50%] top-[50%] z-50 w-full max-w-[calc(100%-2rem)] translate-x-[-50%] translate-y-[-50%] border p-8 sm:max-w-[540px] md:w-full"
            >
                <Dialog.Title
                    class="w-full text-xl font-semibold tracking-tight"
                >
                    Créer une catégorie
                </Dialog.Title>
                <Dialog.Description class="text-foreground-alt !mt-1 text-sm">
                    Remplissez les détails ci-dessous pour créer une nouvelle
                    catégorie.
                </Dialog.Description>

                <Separator.Root class="bg-muted mx-5 !mb-2 !mt-5 block h-px" />

                <form
                    method="POST"
                    action="?/createCategory"
                    enctype="multipart/form-data"
                    use:createEnhance
                    class="grid gap-4"
                >
                    {#if $createErrors._errors}
                        <div class="text-sm text-destructive mt-4">
                            {$createErrors._errors[0]}
                        </div>
                    {/if}

                    <div
                        class="grid grid-cols-[max-content_1fr] gap-4 w-full items-center pb-11 pt-7"
                    >
                        {@render field(
                            "name",
                            "Nom",
                            input,
                            $createForm.name,
                            $createErrors.name,
                        )}
                        {@render field(
                            "description",
                            "Description",
                            textarea,
                            $createForm.description,
                            $createErrors.description,
                        )}
                        {@render field("image", "Image", fileInput, null, null)}
                    </div>
                    <div class="flex w-full justify-end gap-3">
                        <Button.Root type="button" class="cursor-pointer">
                            <Dialog.Close
                                class="h-input rounded-input border border-border-input hover:bg-dark-04 focus-visible:ring-dark focus-visible:ring-offset-background focus-visible:outline-hidden inline-flex items-center justify-center px-6 text-sm font-medium focus-visible:ring-2 focus-visible:ring-offset-2 active:scale-[0.98] cursor-pointer transition-all"
                            >
                                Annuler
                            </Dialog.Close>
                        </Button.Root>
                        <Button.Root type="submit" class="cursor-pointer">
                            <div
                                class="h-input rounded-input bg-dark text-background shadow-mini hover:bg-dark/95 focus-visible:ring-dark focus-visible:ring-offset-background focus-visible:outline-hidden inline-flex items-center justify-center px-8 text-sm font-semibold focus-visible:ring-2 focus-visible:ring-offset-2 active:scale-[0.98] cursor-pointer transition-all"
                            >
                                Créer
                            </div>
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
{/if}

<!-- Edit Category Dialog/Drawer -->
{#if isMobile}
    <Drawer bind:isOpen={editDialogOpen}>
        <div class="flex items-center justify-between mb-4">
            <h2 class="text-xl font-semibold tracking-tight">
                Modifier la catégorie
            </h2>
            <button
                type="button"
                onclick={() => (editDialogOpen = false)}
                class="p-2 hover:bg-dark-04 rounded-md transition-all"
            >
                <X class="text-foreground size-5" />
            </button>
        </div>
        <p class="text-foreground-alt text-sm mb-6">
            Mettez à jour les détails de la catégorie.
        </p>

        <form
            method="POST"
            action="?/updateCategory"
            enctype="multipart/form-data"
            use:updateEnhance
            class="grid gap-4"
        >
            <input type="hidden" name="id" bind:value={$updateForm.id} />

            {#if $updateErrors._errors}
                <div class="text-sm text-destructive mt-4">
                    {$updateErrors._errors[0]}
                </div>
            {/if}

            <div class="grid grid-cols-1 gap-4 w-full pb-4">
                {@render field(
                    "name",
                    "Nom",
                    inputUpdate,
                    $updateForm.name,
                    $updateErrors.name,
                )}
                {@render field(
                    "description",
                    "Description",
                    textareaUpdate,
                    $updateForm.description,
                    $updateErrors.description,
                )}
                {@render field(
                    "status",
                    "Statut",
                    statusSelect,
                    $updateForm.status,
                    $updateErrors.status,
                )}
                {@render field("image", "Image", fileInput, null, null)}
            </div>
            <div class="flex w-full justify-end gap-3">
                <button
                    type="button"
                    onclick={() => (editDialogOpen = false)}
                    class="h-input rounded-input border border-border-input hover:bg-dark-04 focus-visible:ring-dark focus-visible:ring-offset-background focus-visible:outline-hidden inline-flex items-center justify-center px-6 text-sm font-medium focus-visible:ring-2 focus-visible:ring-offset-2 active:scale-[0.98] cursor-pointer transition-all"
                >
                    Annuler
                </button>
                <button
                    type="submit"
                    class="h-input rounded-input bg-dark text-background shadow-mini hover:bg-dark/95 focus-visible:ring-dark focus-visible:ring-offset-background focus-visible:outline-hidden inline-flex items-center justify-center px-8 text-sm font-semibold focus-visible:ring-2 focus-visible:ring-offset-2 active:scale-[0.98] cursor-pointer transition-all"
                >
                    Enregistrer
                </button>
            </div>
        </form>
    </Drawer>
{:else}
    <Dialog.Root bind:open={editDialogOpen}>
        <Dialog.Portal>
            <Dialog.Overlay
                class="data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 fixed inset-0 z-50 bg-black/80"
            />
            <Dialog.Content
                class="rounded-card-lg bg-background shadow-popover data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 data-[state=closed]:zoom-out-95 data-[state=open]:zoom-in-95 outline-hidden fixed left-[50%] top-[50%] z-50 w-full max-w-[calc(100%-2rem)] translate-x-[-50%] translate-y-[-50%] border p-8 sm:max-w-[540px] md:w-full"
            >
                <Dialog.Title
                    class="w-full text-xl font-semibold tracking-tight"
                >
                    Modifier la catégorie
                </Dialog.Title>
                <Dialog.Description class="text-foreground-alt !mt-1 text-sm">
                    Mettez à jour les détails de la catégorie.
                </Dialog.Description>

                <Separator.Root class="bg-muted mx-5 !mb-2 !mt-5 block h-px" />

                <form
                    method="POST"
                    action="?/updateCategory"
                    enctype="multipart/form-data"
                    use:updateEnhance
                    class="grid gap-4"
                >
                    <input
                        type="hidden"
                        name="id"
                        bind:value={$updateForm.id}
                    />

                    {#if $updateErrors._errors}
                        <div class="text-sm text-destructive mt-4">
                            {$updateErrors._errors[0]}
                        </div>
                    {/if}

                    <div
                        class="grid grid-cols-[max-content_1fr] gap-4 w-full items-center pb-11 pt-7"
                    >
                        {@render field(
                            "name",
                            "Nom",
                            inputUpdate,
                            $updateForm.name,
                            $updateErrors.name,
                        )}
                        {@render field(
                            "description",
                            "Description",
                            textareaUpdate,
                            $updateForm.description,
                            $updateErrors.description,
                        )}
                        {@render field(
                            "status",
                            "Statut",
                            statusSelect,
                            $updateForm.status,
                            $updateErrors.status,
                        )}
                        {@render field("image", "Image", fileInput, null, null)}
                    </div>
                    <div class="flex w-full justify-end gap-3">
                        <Button.Root type="button" class="cursor-pointer">
                            <Dialog.Close
                                class="h-input rounded-input border border-border-input hover:bg-dark-04 focus-visible:ring-dark focus-visible:ring-offset-background focus-visible:outline-hidden inline-flex items-center justify-center px-6 text-sm font-medium focus-visible:ring-2 focus-visible:ring-offset-2 active:scale-[0.98] cursor-pointer transition-all"
                            >
                                Annuler
                            </Dialog.Close>
                        </Button.Root>
                        <Button.Root type="submit" class="cursor-pointer">
                            <div
                                class="h-input rounded-input bg-dark text-background shadow-mini hover:bg-dark/95 focus-visible:ring-dark focus-visible:ring-offset-background focus-visible:outline-hidden inline-flex items-center justify-center px-8 text-sm font-semibold focus-visible:ring-2 focus-visible:ring-offset-2 active:scale-[0.98] cursor-pointer transition-all"
                            >
                                Enregistrer
                            </div>
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
{/if}

<!-- Delete Category Dialog/Drawer -->
{#if isMobile}
    <Drawer bind:isOpen={deleteDialogOpen}>
        <div class="flex items-center justify-between mb-4">
            <h2 class="text-xl font-semibold tracking-tight">
                Supprimer la catégorie
            </h2>
            <button
                type="button"
                onclick={() => (deleteDialogOpen = false)}
                class="p-2 hover:bg-dark-04 rounded-md transition-all"
            >
                <X class="text-foreground size-5" />
            </button>
        </div>
        <p class="text-foreground-alt text-sm mb-6">
            Êtes-vous sûr de vouloir supprimer la catégorie "<span
                class="font-medium">{selectedCategory?.name}</span
            >" ? Cette action est irréversible.
        </p>

        {#if $deleteErrors._errors}
            <div class="text-sm text-destructive mb-4">
                {$deleteErrors._errors[0]}
            </div>
        {/if}

        <form method="POST" action="?/deleteCategory" use:deleteEnhance>
            <input type="hidden" name="id" bind:value={$deleteForm.id} />

            <div class="flex w-full justify-end gap-3">
                <button
                    type="button"
                    onclick={() => (deleteDialogOpen = false)}
                    class="h-input rounded-input border border-border-input hover:bg-dark-04 focus-visible:ring-dark focus-visible:ring-offset-background focus-visible:outline-hidden inline-flex items-center justify-center px-6 text-sm font-medium focus-visible:ring-2 focus-visible:ring-offset-2 active:scale-[0.98] cursor-pointer transition-all"
                >
                    Annuler
                </button>
                <button
                    type="submit"
                    class="h-input rounded-input bg-destructive text-white shadow-mini hover:bg-destructive/90 focus-visible:ring-destructive focus-visible:ring-offset-background focus-visible:outline-hidden inline-flex items-center justify-center px-8 text-sm font-semibold focus-visible:ring-2 focus-visible:ring-offset-2 active:scale-[0.98] cursor-pointer transition-all"
                >
                    Supprimer
                </button>
            </div>
        </form>
    </Drawer>
{:else}
    <Dialog.Root bind:open={deleteDialogOpen}>
        <Dialog.Portal>
            <Dialog.Overlay
                class="data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 fixed inset-0 z-50 bg-black/80"
            />
            <Dialog.Content
                class="rounded-card-lg bg-background shadow-popover data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 data-[state-closed]:zoom-out-95 data-[state=open]:zoom-in-95 outline-hidden fixed left-[50%] top-[50%] z-50 w-full max-w-[calc(100%-2rem)] translate-x-[-50%] translate-y-[-50%] border p-8 sm:max-w-[440px] md:w-full"
            >
                <Dialog.Title
                    class="w-full text-xl font-semibold tracking-tight"
                >
                    Supprimer la catégorie
                </Dialog.Title>
                <Dialog.Description class="text-foreground-alt !mt-2 text-sm">
                    Êtes-vous sûr de vouloir supprimer la catégorie "<span
                        class="font-medium">{selectedCategory?.name}</span
                    >" ? Cette action est irréversible.
                </Dialog.Description>

                {#if $deleteErrors._errors}
                    <div class="text-sm text-destructive mt-4">
                        {$deleteErrors._errors[0]}
                    </div>
                {/if}

                <form
                    method="POST"
                    action="?/deleteCategory"
                    use:deleteEnhance
                    class="mt-8"
                >
                    <input
                        type="hidden"
                        name="id"
                        bind:value={$deleteForm.id}
                    />

                    <div class="flex w-full justify-end gap-3">
                        <Button.Root type="button" class="cursor-pointer">
                            <Dialog.Close
                                class="h-input rounded-input border border-border-input hover:bg-dark-04 focus-visible:ring-dark focus-visible:ring-offset-background focus-visible:outline-hidden inline-flex items-center justify-center px-6 text-sm font-medium focus-visible:ring-2 focus-visible:ring-offset-2 active:scale-[0.98] cursor-pointer transition-all"
                            >
                                Annuler
                            </Dialog.Close>
                        </Button.Root>
                        <Button.Root type="submit" class="cursor-pointer">
                            <div
                                class="h-input rounded-input bg-destructive text-white shadow-mini hover:bg-destructive/90 focus-visible:ring-destructive focus-visible:ring-offset-background focus-visible:outline-hidden inline-flex items-center justify-center px-8 text-sm font-semibold focus-visible:ring-2 focus-visible:ring-offset-2 active:scale-[0.98] cursor-pointer transition-all"
                            >
                                Supprimer
                            </div>
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
{/if}

<!-- Snippets for form fields -->
{#snippet field(
    name: string,
    label: string,
    inputSnippet: Input,
    value: any,
    error: any,
)}
    <Label.Root for={name} class="text-sm font-semibold text-right">
        {label}
    </Label.Root>
    <div class="relative w-full">
        {@render inputSnippet(name)}
        {#if error && error.length > 0}
            <p class="text-xs text-destructive mt-1">{error[0]}</p>
        {/if}
    </div>
{/snippet}

{#snippet input(fieldName: string)}
    {#if fieldName === "name"}
        <input
            id={fieldName}
            type="text"
            name={fieldName}
            bind:value={$createForm.name}
            class="h-input rounded-card-sm border-border-input bg-background placeholder:text-foreground-alt/50 hover:border-dark-40 focus:ring-foreground focus:ring-offset-background focus:outline-hidden inline-flex w-full items-center border px-4 text-base focus:ring-2 focus:ring-offset-2 sm:text-sm transition-all"
            required
        />
    {/if}
{/snippet}

{#snippet textarea(fieldName: string)}
    {#if fieldName === "description"}
        <textarea
            id={fieldName}
            name={fieldName}
            rows="5"
            bind:value={$createForm.description}
            class="rounded-card-sm border border-border-input bg-background placeholder:text-foreground-alt/50 hover:border-dark-40 focus:ring-foreground focus:ring-offset-background focus:outline-hidden w-full px-4 py-2 text-base focus:ring-2 focus:ring-offset-2 sm:text-sm transition-all"
            placeholder="Écrivez une description pour la catégorie"
            required
        ></textarea>
    {/if}
{/snippet}

{#snippet inputUpdate(fieldName: string)}
    {#if fieldName === "name"}
        <input
            id={fieldName}
            type="text"
            name={fieldName}
            bind:value={$updateForm.name}
            class="h-input rounded-card-sm border-border-input bg-background placeholder:text-foreground-alt/50 hover:border-dark-40 focus:ring-foreground focus:ring-offset-background focus:outline-hidden inline-flex w-full items-center border px-4 text-base focus:ring-2 focus:ring-offset-2 sm:text-sm transition-all"
            required
        />
    {/if}
{/snippet}

{#snippet textareaUpdate(fieldName: string)}
    {#if fieldName === "description"}
        <textarea
            id={fieldName}
            name={fieldName}
            rows="5"
            bind:value={$updateForm.description}
            class="rounded-card-sm border border-border-input bg-background placeholder:text-foreground-alt/50 hover:border-dark-40 focus:ring-foreground focus:ring-offset-background focus:outline-hidden w-full px-4 py-2 text-base focus:ring-2 focus:ring-offset-2 sm:text-sm transition-all"
            placeholder="Écrivez une description pour la catégorie"
            required
        ></textarea>
    {/if}
{/snippet}

{#snippet statusSelect(name: string)}
    <select
        id={name}
        {name}
        bind:value={$updateForm.status}
        class="h-input rounded-card-sm border-border-input bg-background hover:border-dark-40 focus:ring-foreground focus:ring-offset-background focus:outline-hidden inline-flex w-full items-center border px-4 text-base focus:ring-2 focus:ring-offset-2 sm:text-sm transition-all"
    >
        <option value="draft">Brouillon</option>
        <option value="published">Publié</option>
        <option value="archived">Archivé</option>
    </select>
{/snippet}

{#snippet fileInput(name: string)}
    <input
        type="file"
        id={name}
        {name}
        accept="image/png, image/jpeg, image/jpg, image/webp"
        class="h-input rounded-card-sm border-border-input bg-background hover:border-dark-40 inline-flex w-full items-center border text-base sm:text-sm file:mr-4 file:px-4 file:h-full file:rounded file:border-0 file:text-sm file:font-medium file:bg-dark file:text-background hover:file:bg-dark/90 cursor-pointer transition-all"
    />
{/snippet}
