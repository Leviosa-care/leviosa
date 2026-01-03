<script lang="ts">
    import { Dialog, Button, Label, Separator } from "bits-ui";
    import { UserCheck, X, Shield, Users as UsersIcon } from "@lucide/svelte";
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

    // User type from API
    type User = {
        id: string;
        state: "pending" | "active" | "inactive" | "locked";
        email: string;
        picture?: string;
        created_at: string;
        logged_in_at: string | null;
        role: "visitor" | "standard" | "partner" | "administrator";
        birthdate: string;
        last_name: string;
        first_name: string;
        gender: "male" | "female" | "other" | "prefer_not_to_say";
        telephone: string;
        postal_code: string;
        city: string;
        address1: string;
        address2?: string;
    };

    interface Props {
        data: PageData;
    }

    let { data }: Props = $props();

    // Extract users from page data
    let users: User[] = data.users || [];
    let pendingUsers: User[] = data.pendingUsers || [];

    // Show all users or only pending
    let showPending = $state(false);
    let displayedUsers = $derived(showPending ? pendingUsers : users);

    // Dialog states
    let approveDialogOpen = $state(false);
    let updateRoleDialogOpen = $state(false);

    // Currently selected user
    let selectedUser: User | null = $state(null);

    // Superforms for approve and update role
    const {
        form: approveForm,
        errors: approveErrors,
        enhance: approveEnhance,
    } = superForm(data.approveUserForm, {
        resetForm: true,
        onUpdated({ form }) {
            if (form.valid) {
                approveDialogOpen = false;
            }
        },
    });

    const {
        form: updateRoleForm,
        errors: updateRoleErrors,
        enhance: updateRoleEnhance,
    } = superForm(data.updateUserRoleForm, {
        resetForm: false,
        onUpdated({ form }) {
            if (form.valid) {
                updateRoleDialogOpen = false;
            }
        },
    });

    function openApproveDialog(user: User) {
        selectedUser = user;
        $approveForm.user_id = user.id;
        $approveForm.role = "standard";
        approveDialogOpen = true;
    }

    function openUpdateRoleDialog(user: User) {
        selectedUser = user;
        $updateRoleForm.user_id = user.id;
        $updateRoleForm.role = user.role;
        updateRoleDialogOpen = true;
    }

    function formatDate(dateString: string | null) {
        if (!dateString) return "Jamais";
        return new Date(dateString).toLocaleDateString("fr-FR", {
            year: "numeric",
            month: "short",
            day: "numeric",
        });
    }

    function getRoleBadgeClass(role: string) {
        switch (role) {
            case "administrator":
                return "bg-purple-100 text-purple-800 dark:bg-purple-900 dark:text-purple-200";
            case "partner":
                return "bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-200";
            case "standard":
                return "bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200";
            case "visitor":
                return "bg-gray-100 text-gray-800 dark:bg-gray-700 dark:text-gray-200";
            default:
                return "bg-gray-100 text-gray-800";
        }
    }

    function getStateBadgeClass(state: string) {
        switch (state) {
            case "active":
                return "bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200";
            case "pending":
                return "bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-200";
            case "inactive":
                return "bg-gray-100 text-gray-800 dark:bg-gray-700 dark:text-gray-200";
            case "locked":
                return "bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200";
            default:
                return "bg-gray-100 text-gray-800";
        }
    }

    function getRoleLabel(role: string) {
        switch (role) {
            case "administrator":
                return "Administrateur";
            case "partner":
                return "Partenaire";
            case "standard":
                return "Standard";
            case "visitor":
                return "Visiteur";
            default:
                return role;
        }
    }

    function getStateLabel(state: string) {
        switch (state) {
            case "active":
                return "Actif";
            case "pending":
                return "En attente";
            case "inactive":
                return "Inactif";
            case "locked":
                return "Verrouillé";
            default:
                return state;
        }
    }
</script>

<div class="h-full bg-white">
    <!-- Filter Toggle -->
    <div class="p-4 md:p-6 border-b border-border-card">
        <div class="flex gap-2">
            <button
                type="button"
                onclick={() => (showPending = false)}
                class="px-4 py-2 rounded-input text-sm font-medium transition-all {!showPending
                    ? 'bg-dark text-white'
                    : 'bg-transparent border border-border-input hover:bg-dark-04'}"
            >
                Tous les utilisateurs ({users.length})
            </button>
            <button
                type="button"
                onclick={() => (showPending = true)}
                class="px-4 py-2 rounded-input text-sm font-medium transition-all {showPending
                    ? 'bg-dark text-white'
                    : 'bg-transparent border border-border-input hover:bg-dark-04'}"
            >
                En attente d'approbation ({pendingUsers.length})
            </button>
        </div>
    </div>

    <!-- Users Table -->
    <div class="p-4 md:p-6 overflow-x-auto">
        {#if displayedUsers.length === 0}
            <div
                class="flex flex-col items-center justify-center py-16 text-center"
            >
                <div
                    class="w-16 h-16 rounded-full bg-dark-04 flex items-center justify-center mb-4"
                >
                    <UsersIcon size={32} class="text-dark-400" />
                </div>
                <h3 class="text-lg font-medium mb-2">
                    {showPending
                        ? "Aucun utilisateur en attente"
                        : "Aucun utilisateur"}
                </h3>
                <p class="text-sm text-foreground-alt max-w-sm">
                    {showPending
                        ? "Il n'y a actuellement aucun utilisateur en attente d'approbation."
                        : "Aucun utilisateur enregistré dans le système."}
                </p>
            </div>
        {:else}
            <div class="overflow-x-auto">
                <table class="w-full">
                    <thead>
                        <tr
                            class="border-b border-border-card text-left text-sm font-semibold text-foreground-alt"
                        >
                            <th class="pb-3 pr-4">Nom</th>
                            <th class="pb-3 pr-4">Email</th>
                            <th class="pb-3 pr-4">Rôle</th>
                            <th class="pb-3 pr-4">Statut</th>
                            <th class="pb-3 pr-4">Ville</th>
                            <th class="pb-3 pr-4">Créé le</th>
                            <th class="pb-3 pr-4">Dernière connexion</th>
                            <th class="pb-3 text-right">Actions</th>
                        </tr>
                    </thead>
                    <tbody>
                        {#each displayedUsers as user (user.id)}
                            <tr
                                class="border-b border-border-card hover:bg-dark-04 transition-colors"
                            >
                                <td class="py-4 pr-4">
                                    <div class="flex items-center gap-3">
                                        <div
                                            class="w-10 h-10 rounded-full bg-dark-10 flex items-center justify-center text-sm font-semibold"
                                        >
                                            {user.first_name[0]}{user
                                                .last_name[0]}
                                        </div>
                                        <div>
                                            <div class="font-medium">
                                                {user.first_name}
                                                {user.last_name}
                                            </div>
                                        </div>
                                    </div>
                                </td>
                                <td
                                    class="py-4 pr-4 text-sm text-foreground-alt"
                                >
                                    {user.email}
                                </td>
                                <td class="py-4 pr-4">
                                    <span
                                        class="px-2 py-1 text-xs font-medium rounded-md {getRoleBadgeClass(
                                            user.role,
                                        )}"
                                    >
                                        {getRoleLabel(user.role)}
                                    </span>
                                </td>
                                <td class="py-4 pr-4">
                                    <span
                                        class="px-2 py-1 text-xs font-medium rounded-md {getStateBadgeClass(
                                            user.state,
                                        )}"
                                    >
                                        {getStateLabel(user.state)}
                                    </span>
                                </td>
                                <td class="py-4 pr-4 text-sm">
                                    {user.city}
                                </td>
                                <td
                                    class="py-4 pr-4 text-sm text-foreground-alt"
                                >
                                    {formatDate(user.created_at)}
                                </td>
                                <td
                                    class="py-4 pr-4 text-sm text-foreground-alt"
                                >
                                    {formatDate(user.logged_in_at)}
                                </td>
                                <td class="py-4 text-right">
                                    <div class="flex gap-2 justify-end">
                                        {#if user.state === "pending"}
                                            <Button.Root
                                                type="button"
                                                class="cursor-pointer"
                                                onclick={() =>
                                                    openApproveDialog(user)}
                                            >
                                                <div
                                                    class="flex items-center gap-2 py-2 px-3 border border-green-500/20 text-green-700 bg-green-50 rounded-input hover:bg-green-100 transition-all text-sm font-medium"
                                                >
                                                    <UserCheck size={14} />
                                                    <span>Approuver</span>
                                                </div>
                                            </Button.Root>
                                        {:else}
                                            <Button.Root
                                                type="button"
                                                class="cursor-pointer"
                                                onclick={() =>
                                                    openUpdateRoleDialog(user)}
                                            >
                                                <div
                                                    class="flex items-center gap-2 py-2 px-3 border border-border-input rounded-input hover:bg-dark-04 transition-all text-sm font-medium"
                                                >
                                                    <Shield size={14} />
                                                    <span>Rôle</span>
                                                </div>
                                            </Button.Root>
                                        {/if}
                                    </div>
                                </td>
                            </tr>
                        {/each}
                    </tbody>
                </table>
            </div>
        {/if}
    </div>
</div>

<!-- Approve User Dialog/Drawer -->
{#if isMobile}
    <Drawer bind:isOpen={approveDialogOpen}>
        <div
            class="sticky top-0 bg-white pb-4 border-b border-border-card -mx-4 px-4 -mt-4 pt-4 z-10"
        >
            <div class="flex items-center justify-between mb-2">
                <h2 class="text-xl font-semibold tracking-tight">
                    Approuver l'utilisateur
                </h2>
                <button
                    type="button"
                    onclick={() => (approveDialogOpen = false)}
                    class="p-2 hover:bg-dark-04 rounded-md transition-all"
                >
                    <X class="text-foreground size-5" />
                </button>
            </div>
            <p class="text-foreground-alt text-sm">
                Approuvez l'utilisateur "<span class="font-medium"
                    >{selectedUser?.first_name}
                    {selectedUser?.last_name}</span
                >" et attribuez-lui un rôle.
            </p>
        </div>

        <form
            method="POST"
            action="?/approveUser"
            use:approveEnhance
            class="grid gap-4 pt-8"
        >
            <input
                type="hidden"
                name="user_id"
                bind:value={$approveForm.user_id}
            />

            {#if $approveErrors._errors}
                <div class="text-sm text-destructive mt-4">
                    {$approveErrors._errors[0]}
                </div>
            {/if}

            <div class="grid grid-cols-1 gap-4 w-full pb-4">
                {@render field(
                    "role",
                    "Rôle",
                    roleSelect,
                    $approveForm.role,
                    $approveErrors.role,
                )}
            </div>
            <div class="flex w-full justify-end gap-3">
                <button
                    type="button"
                    onclick={() => (approveDialogOpen = false)}
                    class="h-input rounded-input border border-border-input hover:bg-dark-04 focus-visible:ring-dark focus-visible:ring-offset-background focus-visible:outline-hidden inline-flex items-center justify-center px-6 text-sm font-medium focus-visible:ring-2 focus-visible:ring-offset-2 active:scale-[0.98] cursor-pointer transition-all"
                >
                    Annuler
                </button>
                <button
                    type="submit"
                    class="h-input rounded-input bg-dark text-background shadow-mini hover:bg-dark/95 focus-visible:ring-dark focus-visible:ring-offset-background focus-visible:outline-hidden inline-flex items-center justify-center px-8 text-sm font-semibold focus-visible:ring-2 focus-visible:ring-offset-2 active:scale-[0.98] cursor-pointer transition-all"
                >
                    Approuver
                </button>
            </div>
        </form>
    </Drawer>
{:else}
    <Dialog.Root bind:open={approveDialogOpen}>
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
                    Approuver l'utilisateur
                </Dialog.Title>
                <Dialog.Description class="text-foreground-alt !mt-2 text-sm">
                    Approuvez l'utilisateur "<span class="font-medium"
                        >{selectedUser?.first_name}
                        {selectedUser?.last_name}</span
                    >" et attribuez-lui un rôle.
                </Dialog.Description>

                <Separator.Root class="bg-muted mx-5 !mb-2 !mt-5 block h-px" />

                <form
                    method="POST"
                    action="?/approveUser"
                    use:approveEnhance
                    class="grid gap-4"
                >
                    <input
                        type="hidden"
                        name="user_id"
                        bind:value={$approveForm.user_id}
                    />

                    {#if $approveErrors._errors}
                        <div class="text-sm text-destructive mt-4">
                            {$approveErrors._errors[0]}
                        </div>
                    {/if}

                    <div
                        class="grid grid-cols-[max-content_1fr] gap-4 w-full items-center pb-11 pt-7"
                    >
                        {@render field(
                            "role",
                            "Rôle",
                            roleSelect,
                            $approveForm.role,
                            $approveErrors.role,
                        )}
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
                                Approuver
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

<!-- Update Role Dialog/Drawer -->
{#if isMobile}
    <Drawer bind:isOpen={updateRoleDialogOpen}>
        <div
            class="sticky top-0 bg-white pb-4 border-b border-border-card -mx-4 px-4 -mt-4 pt-4 z-10"
        >
            <div class="flex items-center justify-between mb-2">
                <h2 class="text-xl font-semibold tracking-tight">
                    Modifier le rôle
                </h2>
                <button
                    type="button"
                    onclick={() => (updateRoleDialogOpen = false)}
                    class="p-2 hover:bg-dark-04 rounded-md transition-all"
                >
                    <X class="text-foreground size-5" />
                </button>
            </div>
            <p class="text-foreground-alt text-sm">
                Modifiez le rôle de "<span class="font-medium"
                    >{selectedUser?.first_name}
                    {selectedUser?.last_name}</span
                >".
            </p>
        </div>

        <form
            method="POST"
            action="?/updateUserRole"
            use:updateRoleEnhance
            class="grid gap-4 pt-8"
        >
            <input
                type="hidden"
                name="user_id"
                bind:value={$updateRoleForm.user_id}
            />

            {#if $updateRoleErrors._errors}
                <div class="text-sm text-destructive mt-4">
                    {$updateRoleErrors._errors[0]}
                </div>
            {/if}

            <div class="grid grid-cols-1 gap-4 w-full pb-4">
                {@render field(
                    "role",
                    "Rôle",
                    roleSelectUpdate,
                    $updateRoleForm.role,
                    $updateRoleErrors.role,
                )}
            </div>
            <div class="flex w-full justify-end gap-3">
                <button
                    type="button"
                    onclick={() => (updateRoleDialogOpen = false)}
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
    <Dialog.Root bind:open={updateRoleDialogOpen}>
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
                    Modifier le rôle
                </Dialog.Title>
                <Dialog.Description class="text-foreground-alt !mt-2 text-sm">
                    Modifiez le rôle de "<span class="font-medium"
                        >{selectedUser?.first_name}
                        {selectedUser?.last_name}</span
                    >".
                </Dialog.Description>

                <Separator.Root class="bg-muted mx-5 !mb-2 !mt-5 block h-px" />

                <form
                    method="POST"
                    action="?/updateUserRole"
                    use:updateRoleEnhance
                    class="grid gap-4"
                >
                    <input
                        type="hidden"
                        name="user_id"
                        bind:value={$updateRoleForm.user_id}
                    />

                    {#if $updateRoleErrors._errors}
                        <div class="text-sm text-destructive mt-4">
                            {$updateRoleErrors._errors[0]}
                        </div>
                    {/if}

                    <div
                        class="grid grid-cols-[max-content_1fr] gap-4 w-full items-center pb-11 pt-7"
                    >
                        {@render field(
                            "role",
                            "Rôle",
                            roleSelectUpdate,
                            $updateRoleForm.role,
                            $updateRoleErrors.role,
                        )}
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

<!-- Snippets for form fields -->
{#snippet field(
    name: string,
    label: string,
    inputSnippet: Input,
    value: any,
    error: any,
)}
    <Label.Root for={name} class="text-sm font-semibold">
        {label}
    </Label.Root>
    <div class="relative w-full">
        {@render inputSnippet(name)}
        {#if error && error.length > 0}
            <p class="text-xs text-destructive mt-1">{error[0]}</p>
        {/if}
    </div>
{/snippet}

{#snippet roleSelect(name: string)}
    <select
        id={name}
        {name}
        bind:value={$approveForm.role}
        class="h-input rounded-card-sm border-border-input bg-background hover:border-dark-40 focus:ring-foreground focus:ring-offset-background focus:outline-hidden inline-flex w-full items-center border px-4 text-base focus:ring-2 focus:ring-offset-2 sm:text-sm transition-all"
    >
        <option value="visitor">Visiteur</option>
        <option value="standard">Standard</option>
        <option value="partner">Partenaire</option>
        <option value="administrator">Administrateur</option>
    </select>
{/snippet}

{#snippet roleSelectUpdate(name: string)}
    <select
        id={name}
        {name}
        bind:value={$updateRoleForm.role}
        class="h-input rounded-card-sm border-border-input bg-background hover:border-dark-40 focus:ring-foreground focus:ring-offset-background focus:outline-hidden inline-flex w-full items-center border px-4 text-base focus:ring-2 focus:ring-offset-2 sm:text-sm transition-all"
    >
        <option value="visitor">Visiteur</option>
        <option value="standard">Standard</option>
        <option value="partner">Partenaire</option>
        <option value="administrator">Administrateur</option>
    </select>
{/snippet}
