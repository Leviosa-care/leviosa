<script lang="ts">
    // import { onMount } from "svelte";
    import { goto } from "$app/navigation";

    import type { PageProps } from "./$types";
    import { superForm } from "sveltekit-superforms";
    import { Button } from "bits-ui";

    import ProgressIndicator from "../ProgressIndicator.svelte";
    import PreviousStep from "../PreviousStep.svelte";
    import Input from "$lib/ui/Input.svelte";
    import { getRegistrationContext } from "$lib/context/register.svelte";

    const register = getRegistrationContext();

    let { data }: PageProps = $props();
    const { form, errors, constraints } = superForm(data.form, {
        onUpdated({ form }) {
            if (form.valid) {
                register.clear();
                goto("/auth/pending");
            }
        },
    });

    // onMount(() => {
    //     const values = register.all();
    //     $form.email = values.email;
    //     $form.password = values.password;
    //     $form.address1 = values.address1;
    //     $form.address2 = values.address2;
    //     $form.telephone = values.telephone;
    //     $form.postalCode = values.postalCode;
    //     $form.city = values.city;
    //     $form.firstname = values.firstname;
    //     $form.lastname = values.lastname;
    //     $form.gender = values.gender;
    //     $form.birthdate = values.birthdate;
    // });

    let isSamePassword = $derived(
        $form.password === $form.confirm && $form.confirm != "",
    );

    import type { Step } from "$lib/types/step";
    let steps: Step[] = [
        { id: 1, name: "Email", status: "complete" },
        { id: 2, name: "Identité", status: "complete" },
        { id: 3, name: "Adresse", status: "complete" },
        { id: 4, name: "Mot de passe", status: "current" },
    ];
</script>

<div
    class="h-[100vh] flex flex-col justify-center items-center px-8 max-w-[1080px] mx-auto"
>
    <div class="w-full max-w-[640px]">
        <PreviousStep />
        <div
            class="p-4 md:p-12 grid gap-12 border-gray-200 border-1 rounded-input"
        >
            <div class="w-full flex-none">
                <ProgressIndicator {steps} />
            </div>
            <div>
                <h2 class="text-2xl font-bold">Votre mot de passe</h2>
                <p class="text-gray-600 sm:text-base text-sm">
                    Veuillez fournir votre mot de passe pour compléter votre
                    profil
                </p>
            </div>
            <form class="grid gap-12" method="POST">
                <div class="grid gap-8">
                    <Input
                        bind:value={$form.password}
                        name="password"
                        label="Mot de passe"
                        type="password"
                        error={$errors.password ? $errors.password[0] : ""}
                        {...$constraints.password}
                    />
                    <Input
                        bind:value={$form.confirm}
                        name="confirm"
                        label="Confirmation de mot de passe"
                        type="password"
                        error={$errors.confirm ? $errors.confirm[0] : ""}
                        {...$constraints.confirm}
                    />
                    {#if !isSamePassword}
                        <p class="text-red-500">
                            Les mots de passe ne correspondent pas. Veuillez
                            corriger votre mot de passe.
                        </p>
                    {/if}
                </div>
                <input type="hidden" name="email" bind:value={$form.email} />
                <input
                    type="hidden"
                    name="address1"
                    bind:value={$form.address1}
                />
                <input
                    type="hidden"
                    name="address2"
                    bind:value={$form.address2}
                />
                <input type="hidden" name="city" bind:value={$form.city} />
                <input
                    type="hidden"
                    name="postalCode"
                    bind:value={$form.postalCode}
                />
                <input
                    type="hidden"
                    name="lastname"
                    bind:value={$form.lastname}
                />
                <input
                    type="hidden"
                    name="firstname"
                    bind:value={$form.firstname}
                />
                <input type="hidden" name="gender" bind:value={$form.gender} />
                <input
                    type="hidden"
                    name="birthdate"
                    bind:value={$form.birthdate}
                />
                <Button.Root
                    disabled={!isSamePassword}
                    class="justify-center items-center h-input rounded-input bg-dark text-background shadow-mini hover:bg-dark/95 focus-visible:ring-dark focus-visible:ring-offset-background focus-visible:outline-hidden inline-flex px-4 text-[15px] font-bold focus-visible:ring-2 focus-visible:ring-offset-2 active:scale-[0.98] cursor-pointer"
                    type="submit">Enregistrer votre mot de passe</Button.Root
                >
            </form>
        </div>
    </div>
</div>
