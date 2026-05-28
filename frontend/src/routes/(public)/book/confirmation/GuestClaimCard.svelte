<script lang="ts">
    import { superForm } from "sveltekit-superforms";
    import { fade } from "svelte/transition";

    import Input from "$lib/ui/Input.svelte";
    import type { SuperValidated } from "sveltekit-superforms";

    type ClaimForm = {
        email: string;
        phone: string;
        password: string;
    };

    type VerifyForm = {
        otp0: string;
        otp1: string;
        otp2: string;
        otp3: string;
        otp4: string;
        otp5: string;
    };

    type GuestInfo = {
        guest_first_name: string;
        guest_last_name: string;
        guest_email: string;
        guest_phone: string;
    };

    type Props = {
        guestInfo: GuestInfo | null;
        guestClaimForm: SuperValidated<ClaimForm>;
        guestClaimVerifyForm: SuperValidated<VerifyForm>;
    };

    let { guestInfo, guestClaimForm: initialClaimForm, guestClaimVerifyForm: initialVerifyForm }: Props = $props();

    const hasEmail = !!guestInfo?.guest_email;
    const hasPhone = !!guestInfo?.guest_phone;

    // Phase 1 form
    const {
        form: claimForm,
        errors: claimErrors,
        enhance: claimEnhance,
        constraints: claimConstraints,
    } = superForm(initialClaimForm, {
        resetForm: false,
        onSubmit: () => {
            claimSubmitting = true;
            claimError = "";
        },
        onResult: ({ result }) => {
            claimSubmitting = false;
            if (result.type === "success" && result.data?.claimPhase === "otp") {
                phase = "otp";
            }
            if (result.type === "error") {
                claimError = result.error?.message ?? "Une erreur est survenue.";
            }
        },
        onUpdated: ({ form }) => {
            if (!form.valid) {
                const errors = form.errors._errors;
                if (Array.isArray(errors) && errors.length > 0) {
                    claimError = errors[0];
                }
            }
        },
    });

    // Phase 2 form
    const {
        form: verifyForm,
        errors: verifyErrors,
        enhance: verifyEnhance,
    } = superForm(initialVerifyForm, {
        resetForm: false,
        onSubmit: () => {
            verifySubmitting = true;
            verifyError = "";
        },
        onResult: ({ result }) => {
            verifySubmitting = false;
            if (result.type === "error") {
                verifyError = result.error?.message ?? "Une erreur est survenue.";
            }
        },
        onUpdated: ({ form }) => {
            if (!form.valid) {
                const errors = form.errors._errors;
                if (Array.isArray(errors) && errors.length > 0) {
                    verifyError = errors[0];
                }
            }
        },
    });

    type Phase = "credentials" | "otp";
    let phase = $state<Phase>("credentials");
    let claimSubmitting = $state(false);
    let claimError = $state("");
    let verifySubmitting = $state(false);
    let verifyError = $state("");

    // The email for the OTP helper text
    const otpEmail = $derived($claimForm.email || guestInfo?.guest_email || "");

    function handleOnPasteOTP(event: ClipboardEvent): void {
        event.preventDefault();
        const clipboardData = event.clipboardData || (window as any).clipboardData;
        const pastedData = clipboardData.getData("text");
        const cleanedData = pastedData.replace(/[^0-9]/g, "");
        if (cleanedData) {
            const digits = cleanedData.substring(0, 6).split("");
            digits.forEach((digit: string, index: number) => {
                ($verifyForm as Record<string, string>)[`otp${index}`] = digit;
            });
            const nextInputIndex = Math.min(digits.length, 5);
            setTimeout(() => {
                const nextInput = document.getElementById(`otp-${nextInputIndex}`) as HTMLInputElement;
                nextInput?.focus();
            }, 0);
        }
    }

    function handleOTPInput(event: Event): void {
        const input = event.target as HTMLInputElement;
        const currentID = parseInt(input.dataset.index ?? "0");
        const nextInput = document.getElementById(`otp-${currentID + 1}`) as HTMLInputElement | null;

        if (!/^\d$/.test(input.value)) {
            input.value = "";
            ($verifyForm as Record<string, string>)[`otp${currentID}`] = "";
            return;
        }

        if (input.value && nextInput) {
            nextInput.focus();
        }
    }

    function handleOTPKeyup(event: KeyboardEvent): void {
        const key = event.key.toLowerCase();
        if (key === "backspace" || key === "delete") {
            const input = event.target as HTMLInputElement;
            const currentID = parseInt(input.dataset.index ?? "0");
            if (input.value.length === 0) {
                const prevInput = document.getElementById(`otp-${currentID - 1}`) as HTMLInputElement | null;
                prevInput?.focus();
            }
        }
    }
</script>

<div class="bg-white rounded-3xl p-6 md:p-8 shadow-sm">
    {#if phase === "credentials"}
        <div transition:fade={{ duration: 200 }}>
            <div class="text-center mb-6">
                <h3 class="text-lg font-semibold text-dark-900 mb-2">
                    Créer votre compte en un clic
                </h3>
                <p class="text-dark-600 text-sm">
                    Renseignez un mot de passe pour accéder à vos réservations.
                </p>
            </div>

            {#if claimError}
                <div class="mb-4 p-3 bg-red-50 border border-red-200 rounded-xl text-red-700 text-sm text-center">
                    {claimError}
                </div>
            {/if}

            <form method="POST" action="?/guestClaim" use:claimEnhance class="grid gap-4">
                <!-- Pre-filled name (read-only display) -->
                <div class="grid grid-cols-2 gap-3">
                    <div>
                        <label class="block text-sm font-medium text-dark-600 mb-1">Prénom</label>
                        <div class="w-full px-4 pb-3 pt-3 border border-dark-100 rounded-xl bg-dark-50 text-dark-400 text-base">
                            {guestInfo?.guest_first_name ?? ""}
                        </div>
                    </div>
                    <div>
                        <label class="block text-sm font-medium text-dark-600 mb-1">Nom</label>
                        <div class="w-full px-4 pb-3 pt-3 border border-dark-100 rounded-xl bg-dark-50 text-dark-400 text-base">
                            {guestInfo?.guest_last_name ?? ""}
                        </div>
                    </div>
                </div>

                <!-- Email field: shown if guest didn't provide one -->
                {#if !hasEmail}
                    <div>
                        <Input
                            name="email"
                            label="Email"
                            type="email"
                            autocomplete="email"
                            bind:value={$claimForm.email}
                            error={$claimErrors.email ? $claimErrors.email[0] : ""}
                        />
                    </div>
                {:else}
                    <div>
                        <label class="block text-sm font-medium text-dark-600 mb-1">Email</label>
                        <div class="w-full px-4 pb-3 pt-3 border border-dark-100 rounded-xl bg-dark-50 text-dark-400 text-base truncate">
                            {guestInfo.guest_email}
                        </div>
                        <input type="hidden" name="email" value={guestInfo.guest_email} />
                    </div>
                {/if}

                <!-- Phone field: shown if guest didn't provide one -->
                {#if !hasPhone}
                    <div>
                        <Input
                            name="phone"
                            label="Téléphone"
                            type="text"
                            autocomplete="tel"
                            bind:value={$claimForm.phone}
                            error={$claimErrors.phone ? $claimErrors.phone[0] : ""}
                        />
                    </div>
                {:else}
                    <input type="hidden" name="phone" value={guestInfo!.guest_phone} />
                {/if}

                <!-- Password -->
                <div>
                    <Input
                        name="password"
                        label="Mot de passe"
                        type="password"
                        autocomplete="new-password"
                        bind:value={$claimForm.password}
                        error={$claimErrors.password ? $claimErrors.password[0] : ""}
                    />
                </div>

                <button
                    type="submit"
                    disabled={claimSubmitting}
                    class="mt-2 justify-center items-center py-3.5 rounded-2xl bg-dark-900 text-white font-semibold text-base
                           hover:bg-dark-800 focus-visible:ring-2 focus-visible:ring-dark-500 focus-visible:ring-offset-2
                           focus-visible:outline-hidden transition-colors active:scale-[0.98] cursor-pointer
                           disabled:opacity-50 disabled:cursor-not-allowed"
                >
                    {#if claimSubmitting}
                        <span class="inline-flex items-center gap-2">
                            <svg class="animate-spin h-4 w-4" viewBox="0 0 24 24">
                                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" fill="none" />
                                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z" />
                            </svg>
                            Envoi en cours…
                        </span>
                    {:else}
                        Créer mon compte
                    {/if}
                </button>
            </form>

            <!-- Fallback link for no-JS -->
            <noscript>
                <div class="mt-4 text-center">
                    <a href="/auth" class="text-dark-600 hover:text-dark-900 underline text-sm">
                        Créer un compte via le formulaire d'inscription
                    </a>
                </div>
            </noscript>
        </div>
    {:else if phase === "otp"}
        <div transition:fade={{ duration: 200 }}>
            <div class="text-center mb-6">
                <h3 class="text-lg font-semibold text-dark-900 mb-2">
                    Vérifiez votre email
                </h3>
                <p class="text-dark-600 text-sm">
                    Entrez le code à 6 chiffres envoyé à <strong>{otpEmail}</strong>.
                </p>
            </div>

            {#if verifyError}
                <div class="mb-4 p-3 bg-red-50 border border-red-200 rounded-xl text-red-700 text-sm text-center">
                    {verifyError}
                </div>
            {/if}

            <form method="POST" action="?/guestClaimVerify" use:verifyEnhance class="grid gap-6">
                <div class="flex justify-center items-center gap-2">
                    {#each Array(6) as _, index}
                        {@const key = `otp${index}` as keyof typeof $verifyForm}
                        <input
                            id="otp-{index}"
                            data-index={index}
                            oninput={handleOTPInput}
                            onkeyup={handleOTPKeyup}
                            onpaste={handleOnPasteOTP}
                            class="h-14 w-14 text-center rounded-xl border border-dark-200 bg-white hover:border-dark-300
                                   focus:ring-2 focus:ring-dark-500 focus:ring-offset-2 focus:outline-hidden
                                   text-xl font-mono text-dark-900 transition-colors"
                            type="text"
                            inputmode="numeric"
                            maxlength="1"
                            name={key}
                            bind:value={$verifyForm[key]}
                        />
                    {/each}
                </div>

                <button
                    type="submit"
                    disabled={verifySubmitting}
                    class="justify-center items-center py-3.5 rounded-2xl bg-dark-900 text-white font-semibold text-base
                           hover:bg-dark-800 focus-visible:ring-2 focus-visible:ring-dark-500 focus-visible:ring-offset-2
                           focus-visible:outline-hidden transition-colors active:scale-[0.98] cursor-pointer
                           disabled:opacity-50 disabled:cursor-not-allowed"
                >
                    {#if verifySubmitting}
                        <span class="inline-flex items-center gap-2">
                            <svg class="animate-spin h-4 w-4" viewBox="0 0 24 24">
                                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" fill="none" />
                                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z" />
                            </svg>
                            Vérification…
                        </span>
                    {:else}
                        Vérifier le code
                    {/if}
                </button>
            </form>

            <button
                type="button"
                onclick={() => { phase = "credentials"; verifyError = ""; }}
                class="mt-4 w-full text-center text-sm text-dark-500 hover:text-dark-700 cursor-pointer"
            >
                ← Modifier les informations
            </button>
        </div>
    {/if}
</div>
