<script lang="ts">
    import type { PageProps } from "./$types";
    import { goto } from "$app/navigation";

    import { Button } from "bits-ui";
    import { superForm } from "sveltekit-superforms";

    import { getRegistrationContext } from "$lib/context/register.svelte";
    import Input from "$lib/ui/Input.svelte";
    import ProgressIndicator from "../ProgressIndicator.svelte";

    let { data }: PageProps = $props();
    const { form, enhance, errors, constraints } = superForm(data.form, {
        resetForm: true,
        onUpdated({ form }) {
            if (form.valid) {
                const register = getRegistrationContext();
                register.setAddress({
                    address1: form.data.address1,
                    address2: form.data.address2,
                    postalCode: form.data.postalCode,
                    city: form.data.city,
                });
                goto("/auth/password");
            }
        },
    });

    import type { Step } from "$lib/types/step";
    let steps: Step[] = [
        { id: 1, name: "Email", status: "complete" },
        { id: 2, name: "Identité", status: "complete" },
        { id: 3, name: "Adresse", status: "current" },
        { id: 4, name: "Mot de passe", status: "upcoming" },
    ];
    import { ArrowLeft } from "@lucide/svelte";
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
                    Vos informations de domiciliation
                </h3>
                <p class="text-gray-600 sm:text-base text-sm">
                    Veuillez fournir vos informations de domiciliation pour
                    compléter votre profil
                </p>
            </div>
            <form class="grid gap-12" method="POST" use:enhance>
                <div class="grid gap-4">
                    <Input
                        label="Adresse"
                        name="address1"
                        bind:value={$form.address1}
                        error={$errors.address1 ? $errors.address1[0] : ""}
                        {...$constraints.address1}
                    />
                    <Input
                        label="Complément d'adresse"
                        name="address2"
                        {...$constraints.address2}
                        bind:value={$form.address2}
                        error={$errors.address2 ? $errors.address2[0] : ""}
                    />
                    <div class="flex gap-4 items-center">
                        <div class="grid gap-1.5 w-full">
                            <Input
                                label="Code postal"
                                name="postalCode"
                                {...$constraints.postalCode}
                                bind:value={$form.postalCode}
                                error={$errors.postalCode
                                    ? $errors.postalCode[0]
                                    : ""}
                            />
                        </div>
                        <div class="grid gap-1.5 w-full">
                            <Input
                                label="Ville"
                                name="city"
                                {...$constraints.city}
                                bind:value={$form.city}
                                error={$errors.city ? $errors.city[0] : ""}
                            />
                        </div>
                    </div>
                </div>
                <Button.Root
                    class="mt-2 justify-center items-center h-input rounded-input bg-dark text-background shadow-mini hover:bg-dark/95 focus-visible:ring-dark focus-visible:ring-offset-background focus-visible:outline-hidden inline-flex px-4 text-[15px] font-bold focus-visible:ring-2 focus-visible:ring-offset-2 active:scale-[0.98] cursor-pointer"
                    type="submit"
                >
                    Poursuivre l'inscription
                </Button.Root>
            </form>
        </div>
    </div>
</div>
