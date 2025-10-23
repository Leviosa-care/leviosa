<script lang="ts">
    import { goto } from "$app/navigation";
    import type { PageProps } from "./$types";

    import { superForm } from "sveltekit-superforms";
    import { Button } from "bits-ui";
    import { DateField } from "bits-ui";
    import { CalendarDate } from "@internationalized/date";
    import { ArrowLeft } from "@lucide/svelte";

    import Input from "$lib/ui/Input.svelte";
    import ProgressIndicator from "../ProgressIndicator.svelte";

    let { data }: PageProps = $props();
    const { form, enhance, errors } = superForm(data.form, {
        resetForm: true,
        onUpdated({ form }) {
            if (form.valid) {
                register.setGeneral({
                    firstname: form.data.firstname,
                    lastname: form.data.lastname,
                    gender: form.data.gender,
                    telephone: form.data.telephone,
                    birthdate: form.data.birthdate,
                });
                goto("/auth/address");
            }
        },
    });
    let birthdate = $state(new CalendarDate(2000, 1, 1));
    // TODO: have the maxValue for the date of tody - 18 years
    // TODO: find a way to change the format of the date in the component using the french format for dates, not the english one
    import { getRegistrationContext } from "$lib/context/register.svelte";
    const register = getRegistrationContext();

    import type { Step } from "$lib/types/step";
    let steps: Step[] = [
        { id: 1, name: "Email", status: "complete" },
        { id: 2, name: "Identité", status: "current" },
        { id: 3, name: "Adresse", status: "upcoming" },
        { id: 4, name: "Mot de passe", status: "upcoming" },
    ];
</script>

<div
    class="h-[100vh] flex flex-col justify-center items-center px-8 max-w-[1080px] mx-auto"
>
    <div class="w-full max-w-[640px]">
        <Button.Root
            class="mb-4 justify-center gap-2 items-center h-input rounded-input hover:bg-gray-50 focus-visible:ring-dark focus-visible:ring-offset-background focus-visible:outline-hidden inline-flex px-4 text-[15px] font-bold focus-visible:ring-2 focus-visible:ring-offset-2 active:scale-[0.98] cursor-pointer"
            type="button"
        >
            <ArrowLeft size={16} />
            <span class="text-xs sm:text-base">Retour à l'étape précédente</span
            >
        </Button.Root>
        <div
            class="p-4 md:p-12 grid gap-12 border-gray-200 border-1 rounded-input"
        >
            <div class="w-full flex-none">
                <ProgressIndicator {steps} />
            </div>
            <div class="mt-8">
                <h3 class="font-bold sm:text-2xl text-lg">
                    Informations personnelles
                </h3>
                <p class="text-gray-600 sm:text-base text-sm">
                    Veuillez fournir vos coordonnées personnelles pour compléter
                    votre profil
                </p>
            </div>
            <form method="POST" use:enhance class="grid gap-6">
                <div class="flex items-center gap-4">
                    <Input
                        label="nom"
                        name="lastname"
                        bind:value={$form.lastname}
                        error={$errors.lastname ? $errors.lastname[0] : ""}
                    />
                    <Input
                        label="prenom"
                        name="firstname"
                        bind:value={$form.firstname}
                        error={$errors.firstname ? $errors.firstname[0] : ""}
                    />
                </div>
                <Input
                    label="telephone"
                    name="telephone"
                    error={$errors.telephone ? $errors.telephone[0] : ""}
                    bind:value={$form.telephone}
                />
                <div class="flex gap-4 items-center">
                    <div class="grid gap-1.5 w-full">
                        <label class="font-semibold" for="gender">Genre</label>
                        <select
                            name="gender"
                            id="gender"
                            class="border-border-input border px-2 py-3 rounded-lg"
                        >
                            <option value=""
                                >Sélectionnez votre genre (optionnel)</option
                            >
                            <option value="woman">Femme</option>
                            <option value="man">Homme</option>
                            <option value="non_binary">Non binaire</option>
                            <option value="prefer_not_to_say"
                                >Je préfère ne pas le dire</option
                            >
                            <option value="custom"
                                >Je préfère décrire mon genre</option
                            >
                        </select>
                    </div>
                    <DateField.Root bind:value={birthdate} locale="fr">
                        <div class="flex w-full flex-col gap-1.5">
                            <DateField.Label
                                class="block select-none font-semibold"
                                >Date de naissance</DateField.Label
                            >
                            <DateField.Input
                                name="birthdate"
                                class="h-input rounded-input border-border-input bg-background text-foreground focus-within:border-border-input-hover focus-within:shadow-date-field-focus hover:border-border-input-hover data-invalid:border-destructive flex w-full select-none items-center border px-2 py-3 text-sm tracking-[0.01em] "
                            >
                                {#snippet children({ segments })}
                                    {#each segments as { part, value }}
                                        <div class="inline-block select-none">
                                            {#if part === "literal"}
                                                <DateField.Segment
                                                    {part}
                                                    class="text-muted-foreground p-1"
                                                >
                                                    {value}
                                                </DateField.Segment>
                                            {:else}
                                                <DateField.Segment
                                                    {part}
                                                    class="rounded-5px hover:bg-muted focus:bg-muted focus:text-foreground aria-[valuetext=Empty]:text-muted-foreground data-invalid:text-destructive focus-visible:ring-0! focus-visible:ring-offset-0! px-1 py-1"
                                                >
                                                    {value}
                                                </DateField.Segment>
                                            {/if}
                                        </div>
                                    {/each}
                                {/snippet}
                            </DateField.Input>
                        </div>
                    </DateField.Root>
                </div>
                <Button.Root
                    class="mt-10 justify-center gap-4 items-center h-input rounded-input bg-dark text-background shadow-mini hover:bg-dark/95 focus-visible:ring-dark focus-visible:ring-offset-background focus-visible:outline-hidden inline-flex px-4 text-[15px] font-bold focus-visible:ring-2 focus-visible:ring-offset-2 active:scale-[0.98] cursor-pointer"
                    type="submit"
                >
                    Poursuivre l'inscription
                </Button.Root>
                <p class="text-sm text-center mt-4 text-muted-foreground">
                    Vos informations sont stockées et protégées en toute
                    sécurité.
                </p>
            </form>
        </div>
    </div>
</div>
