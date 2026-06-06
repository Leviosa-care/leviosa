<script lang="ts">
    import type { PageProps } from "./$types";

    import { Button } from "bits-ui";
    import { superForm } from "sveltekit-superforms";
    import { goto } from "$app/navigation";

    import { ArrowLeft } from "@lucide/svelte";

    import ProgressIndicator from "../ProgressIndicator.svelte";
    import { getRegistrationContext } from "$lib/context/register.svelte";
    import { focusNextTarget, focusPrevTarget } from "./helper";

    let { data }: PageProps = $props();
    const { form, errors, enhance } = superForm(data.form, {
        resetForm: true,
    });

    const register = getRegistrationContext();
    const email = register.value.email;

    export function handleOnPasteOTP(event: ClipboardEvent): void {
        event.preventDefault();
        const clipboardData =
            event.clipboardData || (window as any).clipboardData;
        const pastedData = clipboardData.getData("text");
        // Clean data - keep only numbers
        const cleanedData = pastedData.replace(/[^0-9]/g, "");
        // Update OTP input values
        if (cleanedData) {
            // Limit to 6 digits (or the number of inputs you have)
            const digits = cleanedData.substring(0, 6).split("");

            // Fill the inputs with the digits
            digits.forEach((digit: string, index: number) => {
                ($form as Record<string, string>)[`otp${index}`] = digit;
            });

            // Focus the appropriate input based on how many digits were pasted
            const nextInputIndex = Math.min(digits.length, 5);
            const timeoutID = setTimeout(() => {
                const nextInput = document.getElementById(
                    `${nextInputIndex}`,
                ) as HTMLInputElement;
                nextInput.focus();
            }, 0);
            clearTimeout(timeoutID);
        }
    }
    // import SuperDebug from "sveltekit-superforms";

    import type { Step } from "$lib/types/step";
    let steps: Step[] = [
        { id: 1, name: "Email", status: "current" },
        { id: 2, name: "Identité", status: "upcoming" },
        { id: 3, name: "Adresse", status: "upcoming" },
        { id: 4, name: "Mot de passe", status: "upcoming" },
    ];
</script>

<!-- <SuperDebug -->
<!--     data={{ -->
<!--         values: $form, -->
<!--         err: $errors, -->
<!--     }} -->
<!-- /> -->

<div
    class="h-[100vh] flex flex-col justify-center items-center px-8 max-w-[1080px] mx-auto"
>
    <div class="w-full max-w-[640px]">
        <Button.Root
            class="mb-4 justify-center gap-2 items-center h-input rounded-input hover:bg-muted focus-visible:ring-dark focus-visible:ring-offset-background focus-visible:outline-hidden inline-flex px-4 text-[15px] font-bold focus-visible:ring-2 focus-visible:ring-offset-2 active:scale-[0.98] cursor-pointer"
            type="button"
            onclick={() => goto("/auth")}
        >
            <ArrowLeft size={16} />
            <span class="text-xs sm:text-base">Retour à l'étape précédente</span
            >
        </Button.Root>
        <div
            class="p-4 md:p-12 grid gap-12 border-border-card border-1 rounded-input"
        >
            <div class="grid gap-8">
                <div class="w-full flex-none">
                    <ProgressIndicator {steps} />
                </div>
                <div class="mt-8 grid gap-2">
                    <h3 class="font-bold sm:text-2xl text-lg">
                        Renseignez le code reçu par email
                    </h3>
                    <p class="text-muted-foreground sm:text-base text-sm">
                        Entrez le code à 6 chiffres envoyé à <strong
                            >{email}</strong
                        >.
                    </p>
                </div>
            </div>
            <form method="POST" use:enhance class="grid gap-8">
                <div class="justify-center flex items-center sm:gap-2">
                    {#each Array(6) as _, index}
                        {@const key = `otp${index}` as keyof typeof $form}
                        <input
                            id={`${index}`}
                            oninput={focusNextTarget}
                            onkeyup={focusPrevTarget}
                            onpaste={handleOnPasteOTP}
                            class="h-input w-input sm:h-16 sm:w-16 text-center rounded-card-sm border-border-input bg-background hover:border-dark-40 focus:ring-foreground focus:ring-offset-background focus:outline-hidden inline-flex items-center border text-xl focus:ring-2 focus:ring-offset-2 sm:text-2xl"
                            type="text"
                            name={key}
                            class:border-destructive={$errors[key]}
                            class:focus:ring-destructive={$errors[key]}
                            bind:value={$form[key]}
                        />
                    {/each}
                </div>
                {#if $errors}
                    <span class="text-destructive text-center"
                        >{$errors._errors}</span
                    >
                {/if}
                <div class="grid">
                    <Button.Root
                        class="justify-center items-center h-input rounded-input bg-dark text-background shadow-mini hover:bg-dark/95 focus-visible:ring-dark focus-visible:ring-offset-background focus-visible:outline-hidden inline-flex px-4 text-base font-bold focus-visible:ring-2 focus-visible:ring-offset-2 active:scale-[0.98] cursor-pointer"
                        type="submit"
                    >
                        Vérifie ton code
                    </Button.Root>
                    <!-- TODO: make that thing a button that is going to send an action  -->
                    <p
                        class="text-sm sm:text-md text-center mt-8 text-muted-foreground"
                    >
                        Tu n'as pas reçu de code ? <a
                            class="hover:underline text-muted-foreground"
                            href="/terms-of-service"
                        >
                            Renvoyer le code</a
                        >
                    </p>
                </div>
            </form>
        </div>
    </div>
</div>
