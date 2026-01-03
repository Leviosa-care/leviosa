<script lang="ts">
    import type { PageData } from "./$types";
    import { superForm } from "sveltekit-superforms";
    import { Dialog, Button, Label, Separator } from "bits-ui";
    import { X } from "@lucide/svelte";

    let { data }: { data: PageData } = $props();

    // Company Name Form
    const {
        form: companyNameForm,
        enhance: companyNameEnhance,
        delayed: companyNameDelayed,
    } = superForm(data.companyNameForm, {
        resetForm: false,
    });

    // Company Email Form
    const {
        form: companyEmailForm,
        enhance: companyEmailEnhance,
        delayed: companyEmailDelayed,
    } = superForm(data.companyEmailForm, {
        resetForm: false,
    });

    // Company Phone Form
    const {
        form: companyPhoneForm,
        enhance: companyPhoneEnhance,
        delayed: companyPhoneDelayed,
    } = superForm(data.companyPhoneForm, {
        resetForm: false,
    });

    // Company Address Form
    const {
        form: companyAddressForm,
        enhance: companyAddressEnhance,
        delayed: companyAddressDelayed,
    } = superForm(data.companyAddressForm, {
        resetForm: false,
    });

    // Company Instagram Form
    const {
        form: companyInstagramForm,
        enhance: companyInstagramEnhance,
        delayed: companyInstagramDelayed,
    } = superForm(data.companyInstagramForm, {
        resetForm: false,
    });

    // OTP Duration Form
    const {
        form: otpDurationForm,
        enhance: otpDurationEnhance,
        delayed: otpDurationDelayed,
    } = superForm(data.otpDurationForm, {
        resetForm: false,
    });

    // OTP Length Form
    const {
        form: otpLengthForm,
        enhance: otpLengthEnhance,
        delayed: otpLengthDelayed,
    } = superForm(data.otpLengthForm, {
        resetForm: false,
    });

    // OTP Max Attempts Form
    const {
        form: otpMaxAttemptsForm,
        enhance: otpMaxAttemptsEnhance,
        delayed: otpMaxAttemptsDelayed,
    } = superForm(data.otpMaxAttemptsForm, {
        resetForm: false,
    });

    // Access Token Duration Form
    const {
        form: accessTokenDurationForm,
        enhance: accessTokenDurationEnhance,
        delayed: accessTokenDurationDelayed,
    } = superForm(data.accessTokenDurationForm, {
        resetForm: false,
    });

    // Refresh Token Duration Form
    const {
        form: refreshTokenDurationForm,
        enhance: refreshTokenDurationEnhance,
        delayed: refreshTokenDurationDelayed,
    } = superForm(data.refreshTokenDurationForm, {
        resetForm: false,
    });

    // Dialog states
    let companyNameDialogOpen = $state(false);
    let companyEmailDialogOpen = $state(false);
    let companyPhoneDialogOpen = $state(false);
    let companyAddressDialogOpen = $state(false);
    let companyInstagramDialogOpen = $state(false);
    let otpDurationDialogOpen = $state(false);
    let otpLengthDialogOpen = $state(false);
    let otpMaxAttemptsDialogOpen = $state(false);
    let accessTokenDurationDialogOpen = $state(false);
    let refreshTokenDurationDialogOpen = $state(false);

    // Current settings from data
    let currentSettings = $derived(data.settings);

    // Helper function to format duration in seconds to human-readable
    function formatSeconds(seconds: number): string {
        if (seconds < 60) return `${seconds}s`;
        const minutes = Math.floor(seconds / 60);
        return `${minutes} min`;
    }

    // Helper function to format duration in minutes to human-readable
    function formatMinutes(minutes: number): string {
        if (minutes < 60) return `${minutes} min`;
        const hours = Math.floor(minutes / 60);
        return `${hours}h`;
    }

    // Helper function to format duration in hours to human-readable
    function formatHours(hours: number): string {
        if (hours < 24) return `${hours}h`;
        const days = Math.floor(hours / 24);
        return `${days} jours`;
    }
</script>

<div class="flex-1 overflow-y-auto p-6">
    <div class="max-w-4xl mx-auto space-y-6">
        <!-- Company Information Section -->
        <div class="bg-white rounded-lg border border-border-card overflow-hidden">
            <div class="px-6 py-4 border-b border-border-card">
                <h2 class="text-lg font-semibold">Informations de l'entreprise</h2>
                <p class="text-sm text-foreground-alt mt-1">
                    Gérez les informations de contact et les réseaux sociaux
                </p>
            </div>
            <div class="p-6 space-y-4">
                <!-- Company Name -->
                <div class="flex items-center justify-between py-3 border-b border-border-card last:border-0">
                    <div class="flex-1">
                        <div class="font-medium text-sm">Nom de l'entreprise</div>
                        <div class="text-sm text-foreground-alt mt-1">
                            {currentSettings?.company?.name || "Non défini"}
                        </div>
                    </div>
                    <Button.Root
                        type="button"
                        class="cursor-pointer"
                        onclick={() => {
                            $companyNameForm.name = currentSettings?.company?.name || "";
                            companyNameDialogOpen = true;
                        }}
                    >
                        <div class="h-input rounded-input border border-border-input hover:bg-dark-04 focus-visible:ring-dark focus-visible:ring-offset-background focus-visible:outline-hidden inline-flex items-center justify-center px-4 text-sm font-medium focus-visible:ring-2 focus-visible:ring-offset-2 active:scale-[0.98] cursor-pointer transition-all">
                            Modifier
                        </div>
                    </Button.Root>
                </div>

                <!-- Company Email -->
                <div class="flex items-center justify-between py-3 border-b border-border-card last:border-0">
                    <div class="flex-1">
                        <div class="font-medium text-sm">Email de contact</div>
                        <div class="text-sm text-foreground-alt mt-1">
                            {currentSettings?.company?.email || "Non défini"}
                        </div>
                    </div>
                    <Button.Root
                        type="button"
                        class="cursor-pointer"
                        onclick={() => {
                            $companyEmailForm.email = currentSettings?.company?.email || "";
                            companyEmailDialogOpen = true;
                        }}
                    >
                        <div class="h-input rounded-input border border-border-input hover:bg-dark-04 focus-visible:ring-dark focus-visible:ring-offset-background focus-visible:outline-hidden inline-flex items-center justify-center px-4 text-sm font-medium focus-visible:ring-2 focus-visible:ring-offset-2 active:scale-[0.98] cursor-pointer transition-all">
                            Modifier
                        </div>
                    </Button.Root>
                </div>

                <!-- Company Phone -->
                <div class="flex items-center justify-between py-3 border-b border-border-card last:border-0">
                    <div class="flex-1">
                        <div class="font-medium text-sm">Téléphone</div>
                        <div class="text-sm text-foreground-alt mt-1">
                            {currentSettings?.company?.telephone || "Non défini"}
                        </div>
                    </div>
                    <Button.Root
                        type="button"
                        class="cursor-pointer"
                        onclick={() => {
                            $companyPhoneForm.telephone = currentSettings?.company?.telephone || "";
                            companyPhoneDialogOpen = true;
                        }}
                    >
                        <div class="h-input rounded-input border border-border-input hover:bg-dark-04 focus-visible:ring-dark focus-visible:ring-offset-background focus-visible:outline-hidden inline-flex items-center justify-center px-4 text-sm font-medium focus-visible:ring-2 focus-visible:ring-offset-2 active:scale-[0.98] cursor-pointer transition-all">
                            Modifier
                        </div>
                    </Button.Root>
                </div>

                <!-- Company Address -->
                <div class="flex items-center justify-between py-3 border-b border-border-card last:border-0">
                    <div class="flex-1">
                        <div class="font-medium text-sm">Adresse</div>
                        <div class="text-sm text-foreground-alt mt-1">
                            {currentSettings?.company?.address || "Non définie"}
                        </div>
                    </div>
                    <Button.Root
                        type="button"
                        class="cursor-pointer"
                        onclick={() => {
                            $companyAddressForm.address = currentSettings?.company?.address || "";
                            companyAddressDialogOpen = true;
                        }}
                    >
                        <div class="h-input rounded-input border border-border-input hover:bg-dark-04 focus-visible:ring-dark focus-visible:ring-offset-background focus-visible:outline-hidden inline-flex items-center justify-center px-4 text-sm font-medium focus-visible:ring-2 focus-visible:ring-offset-2 active:scale-[0.98] cursor-pointer transition-all">
                            Modifier
                        </div>
                    </Button.Root>
                </div>

                <!-- Company Instagram -->
                <div class="flex items-center justify-between py-3 border-b border-border-card last:border-0">
                    <div class="flex-1">
                        <div class="font-medium text-sm">Instagram</div>
                        <div class="text-sm text-foreground-alt mt-1">
                            {currentSettings?.company?.instagram || "Non défini"}
                        </div>
                    </div>
                    <Button.Root
                        type="button"
                        class="cursor-pointer"
                        onclick={() => {
                            $companyInstagramForm.instagram = currentSettings?.company?.instagram || "";
                            companyInstagramDialogOpen = true;
                        }}
                    >
                        <div class="h-input rounded-input border border-border-input hover:bg-dark-04 focus-visible:ring-dark focus-visible:ring-offset-background focus-visible:outline-hidden inline-flex items-center justify-center px-4 text-sm font-medium focus-visible:ring-2 focus-visible:ring-offset-2 active:scale-[0.98] cursor-pointer transition-all">
                            Modifier
                        </div>
                    </Button.Root>
                </div>
            </div>
        </div>

        <!-- OTP Configuration Section -->
        <div class="bg-white rounded-lg border border-border-card overflow-hidden">
            <div class="px-6 py-4 border-b border-border-card">
                <h2 class="text-lg font-semibold">Configuration OTP</h2>
                <p class="text-sm text-foreground-alt mt-1">
                    Paramètres de sécurité pour les codes de vérification
                </p>
            </div>
            <div class="p-6 space-y-4">
                <!-- OTP Duration -->
                <div class="flex items-center justify-between py-3 border-b border-border-card last:border-0">
                    <div class="flex-1">
                        <div class="font-medium text-sm">Durée de validité OTP</div>
                        <div class="text-sm text-foreground-alt mt-1">
                            {formatSeconds(currentSettings?.otp?.duration || 300)} (60-3600 secondes)
                        </div>
                    </div>
                    <Button.Root
                        type="button"
                        class="cursor-pointer"
                        onclick={() => {
                            $otpDurationForm.duration = currentSettings?.otp?.duration || 300;
                            otpDurationDialogOpen = true;
                        }}
                    >
                        <div class="h-input rounded-input border border-border-input hover:bg-dark-04 focus-visible:ring-dark focus-visible:ring-offset-background focus-visible:outline-hidden inline-flex items-center justify-center px-4 text-sm font-medium focus-visible:ring-2 focus-visible:ring-offset-2 active:scale-[0.98] cursor-pointer transition-all">
                            Modifier
                        </div>
                    </Button.Root>
                </div>

                <!-- OTP Length -->
                <div class="flex items-center justify-between py-3 border-b border-border-card last:border-0">
                    <div class="flex-1">
                        <div class="font-medium text-sm">Longueur du code OTP</div>
                        <div class="text-sm text-foreground-alt mt-1">
                            {currentSettings?.otp?.length || 6} chiffres (4-10 chiffres)
                        </div>
                    </div>
                    <Button.Root
                        type="button"
                        class="cursor-pointer"
                        onclick={() => {
                            $otpLengthForm.length = currentSettings?.otp?.length || 6;
                            otpLengthDialogOpen = true;
                        }}
                    >
                        <div class="h-input rounded-input border border-border-input hover:bg-dark-04 focus-visible:ring-dark focus-visible:ring-offset-background focus-visible:outline-hidden inline-flex items-center justify-center px-4 text-sm font-medium focus-visible:ring-2 focus-visible:ring-offset-2 active:scale-[0.98] cursor-pointer transition-all">
                            Modifier
                        </div>
                    </Button.Root>
                </div>

                <!-- OTP Max Attempts -->
                <div class="flex items-center justify-between py-3 border-b border-border-card last:border-0">
                    <div class="flex-1">
                        <div class="font-medium text-sm">Tentatives maximales</div>
                        <div class="text-sm text-foreground-alt mt-1">
                            {currentSettings?.otp?.max_attempts || 5} tentatives (1-10)
                        </div>
                    </div>
                    <Button.Root
                        type="button"
                        class="cursor-pointer"
                        onclick={() => {
                            $otpMaxAttemptsForm.max_attempts = currentSettings?.otp?.max_attempts || 5;
                            otpMaxAttemptsDialogOpen = true;
                        }}
                    >
                        <div class="h-input rounded-input border border-border-input hover:bg-dark-04 focus-visible:ring-dark focus-visible:ring-offset-background focus-visible:outline-hidden inline-flex items-center justify-center px-4 text-sm font-medium focus-visible:ring-2 focus-visible:ring-offset-2 active:scale-[0.98] cursor-pointer transition-all">
                            Modifier
                        </div>
                    </Button.Root>
                </div>
            </div>
        </div>

        <!-- Token Configuration Section -->
        <div class="bg-white rounded-lg border border-border-card overflow-hidden">
            <div class="px-6 py-4 border-b border-border-card">
                <h2 class="text-lg font-semibold">Configuration des tokens</h2>
                <p class="text-sm text-foreground-alt mt-1">
                    Paramètres de durée de vie des tokens d'authentification
                </p>
            </div>
            <div class="p-6 space-y-4">
                <!-- Access Token Duration -->
                <div class="flex items-center justify-between py-3 border-b border-border-card last:border-0">
                    <div class="flex-1">
                        <div class="font-medium text-sm">Durée du token d'accès</div>
                        <div class="text-sm text-foreground-alt mt-1">
                            {formatMinutes(currentSettings?.tokens?.access_duration || 15)} (1-240 minutes)
                        </div>
                    </div>
                    <Button.Root
                        type="button"
                        class="cursor-pointer"
                        onclick={() => {
                            $accessTokenDurationForm.duration = currentSettings?.tokens?.access_duration || 15;
                            accessTokenDurationDialogOpen = true;
                        }}
                    >
                        <div class="h-input rounded-input border border-border-input hover:bg-dark-04 focus-visible:ring-dark focus-visible:ring-offset-background focus-visible:outline-hidden inline-flex items-center justify-center px-4 text-sm font-medium focus-visible:ring-2 focus-visible:ring-offset-2 active:scale-[0.98] cursor-pointer transition-all">
                            Modifier
                        </div>
                    </Button.Root>
                </div>

                <!-- Refresh Token Duration -->
                <div class="flex items-center justify-between py-3 border-b border-border-card last:border-0">
                    <div class="flex-1">
                        <div class="font-medium text-sm">Durée du token de rafraîchissement</div>
                        <div class="text-sm text-foreground-alt mt-1">
                            {formatHours(currentSettings?.tokens?.refresh_duration || 168)} (1-720 heures)
                        </div>
                    </div>
                    <Button.Root
                        type="button"
                        class="cursor-pointer"
                        onclick={() => {
                            $refreshTokenDurationForm.duration = currentSettings?.tokens?.refresh_duration || 168;
                            refreshTokenDurationDialogOpen = true;
                        }}
                    >
                        <div class="h-input rounded-input border border-border-input hover:bg-dark-04 focus-visible:ring-dark focus-visible:ring-offset-background focus-visible:outline-hidden inline-flex items-center justify-center px-4 text-sm font-medium focus-visible:ring-2 focus-visible:ring-offset-2 active:scale-[0.98] cursor-pointer transition-all">
                            Modifier
                        </div>
                    </Button.Root>
                </div>
            </div>
        </div>
    </div>
</div>

<!-- Company Name Dialog -->
<Dialog.Root bind:open={companyNameDialogOpen}>
    <Dialog.Portal>
        <Dialog.Overlay class="data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 fixed inset-0 z-50 bg-black/80" />
        <Dialog.Content class="rounded-card-lg bg-background shadow-popover data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 data-[state=closed]:zoom-out-95 data-[state=open]:zoom-in-95 outline-hidden fixed left-[50%] top-[50%] z-50 w-full max-w-[calc(100%-2rem)] translate-x-[-50%] translate-y-[-50%] border p-8 sm:max-w-[540px] md:w-full">
            <Dialog.Title class="w-full text-xl font-semibold tracking-tight">
                Modifier le nom de l'entreprise
            </Dialog.Title>
            <Dialog.Description class="text-foreground-alt !mt-1 text-sm">
                Définissez le nom officiel de votre entreprise
            </Dialog.Description>

            <Separator.Root class="bg-muted mx-5 !mb-2 !mt-5 block h-px" />

            <form method="POST" action="?/updateCompanyName" use:companyNameEnhance class="grid gap-4">
                <div class="grid gap-2 py-4">
                    <Label.Root for="company-name" class="text-sm font-medium">Nom</Label.Root>
                    <input
                        id="company-name"
                        name="name"
                        type="text"
                        bind:value={$companyNameForm.name}
                        placeholder="Leviosa Spa & Wellness"
                        class="h-input rounded-card-sm border-border-input bg-background placeholder:text-foreground-alt/50 hover:border-dark-40 focus:ring-foreground focus:ring-offset-background focus:outline-hidden inline-flex w-full items-center border px-4 text-base focus:ring-2 focus:ring-offset-2 sm:text-sm transition-all"
                    />
                </div>
                <div class="flex w-full justify-end gap-3">
                    <Button.Root type="button" class="cursor-pointer">
                        <Dialog.Close class="h-input rounded-input border border-border-input hover:bg-dark-04 focus-visible:ring-dark focus-visible:ring-offset-background focus-visible:outline-hidden inline-flex items-center justify-center px-6 text-sm font-medium focus-visible:ring-2 focus-visible:ring-offset-2 active:scale-[0.98] cursor-pointer transition-all">
                            Annuler
                        </Dialog.Close>
                    </Button.Root>
                    <Button.Root type="submit" class="cursor-pointer" disabled={$companyNameDelayed}>
                        <div class="h-input rounded-input bg-dark text-background shadow-mini hover:bg-dark/95 focus-visible:ring-dark focus-visible:ring-offset-background focus-visible:outline-hidden inline-flex items-center justify-center px-8 text-sm font-semibold focus-visible:ring-2 focus-visible:ring-offset-2 active:scale-[0.98] cursor-pointer transition-all">
                            {$companyNameDelayed ? "Enregistrement..." : "Enregistrer"}
                        </div>
                    </Button.Root>
                </div>
            </form>

            <Button.Root type="button" class="cursor-pointer">
                <Dialog.Close class="focus-visible:ring-foreground focus-visible:ring-offset-background focus-visible:outline-hidden absolute right-5 top-5 rounded-md focus-visible:ring-2 focus-visible:ring-offset-2 active:scale-[0.98] cursor-pointer">
                    <X class="text-foreground size-5" />
                    <span class="sr-only">Close</span>
                </Dialog.Close>
            </Button.Root>
        </Dialog.Content>
    </Dialog.Portal>
</Dialog.Root>

<!-- Company Email Dialog -->
<Dialog.Root bind:open={companyEmailDialogOpen}>
    <Dialog.Portal>
        <Dialog.Overlay class="data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 fixed inset-0 z-50 bg-black/80" />
        <Dialog.Content class="rounded-card-lg bg-background shadow-popover data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 data-[state=closed]:zoom-out-95 data-[state=open]:zoom-in-95 outline-hidden fixed left-[50%] top-[50%] z-50 w-full max-w-[calc(100%-2rem)] translate-x-[-50%] translate-y-[-50%] border p-8 sm:max-w-[540px] md:w-full">
            <Dialog.Title class="w-full text-xl font-semibold tracking-tight">
                Modifier l'email de contact
            </Dialog.Title>
            <Dialog.Description class="text-foreground-alt !mt-1 text-sm">
                Définissez l'adresse email de contact de votre entreprise
            </Dialog.Description>

            <Separator.Root class="bg-muted mx-5 !mb-2 !mt-5 block h-px" />

            <form method="POST" action="?/updateCompanyEmail" use:companyEmailEnhance class="grid gap-4">
                <div class="grid gap-2 py-4">
                    <Label.Root for="company-email" class="text-sm font-medium">Email</Label.Root>
                    <input
                        id="company-email"
                        name="email"
                        type="email"
                        bind:value={$companyEmailForm.email}
                        placeholder="contact@leviosa-spa.com"
                        class="h-input rounded-card-sm border-border-input bg-background placeholder:text-foreground-alt/50 hover:border-dark-40 focus:ring-foreground focus:ring-offset-background focus:outline-hidden inline-flex w-full items-center border px-4 text-base focus:ring-2 focus:ring-offset-2 sm:text-sm transition-all"
                    />
                </div>
                <div class="flex w-full justify-end gap-3">
                    <Button.Root type="button" class="cursor-pointer">
                        <Dialog.Close class="h-input rounded-input border border-border-input hover:bg-dark-04 focus-visible:ring-dark focus-visible:ring-offset-background focus-visible:outline-hidden inline-flex items-center justify-center px-6 text-sm font-medium focus-visible:ring-2 focus-visible:ring-offset-2 active:scale-[0.98] cursor-pointer transition-all">
                            Annuler
                        </Dialog.Close>
                    </Button.Root>
                    <Button.Root type="submit" class="cursor-pointer" disabled={$companyEmailDelayed}>
                        <div class="h-input rounded-input bg-dark text-background shadow-mini hover:bg-dark/95 focus-visible:ring-dark focus-visible:ring-offset-background focus-visible:outline-hidden inline-flex items-center justify-center px-8 text-sm font-semibold focus-visible:ring-2 focus-visible:ring-offset-2 active:scale-[0.98] cursor-pointer transition-all">
                            {$companyEmailDelayed ? "Enregistrement..." : "Enregistrer"}
                        </div>
                    </Button.Root>
                </div>
            </form>

            <Button.Root type="button" class="cursor-pointer">
                <Dialog.Close class="focus-visible:ring-foreground focus-visible:ring-offset-background focus-visible:outline-hidden absolute right-5 top-5 rounded-md focus-visible:ring-2 focus-visible:ring-offset-2 active:scale-[0.98] cursor-pointer">
                    <X class="text-foreground size-5" />
                    <span class="sr-only">Close</span>
                </Dialog.Close>
            </Button.Root>
        </Dialog.Content>
    </Dialog.Portal>
</Dialog.Root>

<!-- Company Phone Dialog -->
<Dialog.Root bind:open={companyPhoneDialogOpen}>
    <Dialog.Portal>
        <Dialog.Overlay class="data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 fixed inset-0 z-50 bg-black/80" />
        <Dialog.Content class="rounded-card-lg bg-background shadow-popover data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 data-[state=closed]:zoom-out-95 data-[state=open]:zoom-in-95 outline-hidden fixed left-[50%] top-[50%] z-50 w-full max-w-[calc(100%-2rem)] translate-x-[-50%] translate-y-[-50%] border p-8 sm:max-w-[540px] md:w-full">
            <Dialog.Title class="w-full text-xl font-semibold tracking-tight">
                Modifier le téléphone
            </Dialog.Title>
            <Dialog.Description class="text-foreground-alt !mt-1 text-sm">
                Définissez le numéro de téléphone de votre entreprise
            </Dialog.Description>

            <Separator.Root class="bg-muted mx-5 !mb-2 !mt-5 block h-px" />

            <form method="POST" action="?/updateCompanyPhone" use:companyPhoneEnhance class="grid gap-4">
                <div class="grid gap-2 py-4">
                    <Label.Root for="company-phone" class="text-sm font-medium">Téléphone</Label.Root>
                    <input
                        id="company-phone"
                        name="telephone"
                        type="tel"
                        bind:value={$companyPhoneForm.telephone}
                        placeholder="+33123456789"
                        class="h-input rounded-card-sm border-border-input bg-background placeholder:text-foreground-alt/50 hover:border-dark-40 focus:ring-foreground focus:ring-offset-background focus:outline-hidden inline-flex w-full items-center border px-4 text-base focus:ring-2 focus:ring-offset-2 sm:text-sm transition-all"
                    />
                </div>
                <div class="flex w-full justify-end gap-3">
                    <Button.Root type="button" class="cursor-pointer">
                        <Dialog.Close class="h-input rounded-input border border-border-input hover:bg-dark-04 focus-visible:ring-dark focus-visible:ring-offset-background focus-visible:outline-hidden inline-flex items-center justify-center px-6 text-sm font-medium focus-visible:ring-2 focus-visible:ring-offset-2 active:scale-[0.98] cursor-pointer transition-all">
                            Annuler
                        </Dialog.Close>
                    </Button.Root>
                    <Button.Root type="submit" class="cursor-pointer" disabled={$companyPhoneDelayed}>
                        <div class="h-input rounded-input bg-dark text-background shadow-mini hover:bg-dark/95 focus-visible:ring-dark focus-visible:ring-offset-background focus-visible:outline-hidden inline-flex items-center justify-center px-8 text-sm font-semibold focus-visible:ring-2 focus-visible:ring-offset-2 active:scale-[0.98] cursor-pointer transition-all">
                            {$companyPhoneDelayed ? "Enregistrement..." : "Enregistrer"}
                        </div>
                    </Button.Root>
                </div>
            </form>

            <Button.Root type="button" class="cursor-pointer">
                <Dialog.Close class="focus-visible:ring-foreground focus-visible:ring-offset-background focus-visible:outline-hidden absolute right-5 top-5 rounded-md focus-visible:ring-2 focus-visible:ring-offset-2 active:scale-[0.98] cursor-pointer">
                    <X class="text-foreground size-5" />
                    <span class="sr-only">Close</span>
                </Dialog.Close>
            </Button.Root>
        </Dialog.Content>
    </Dialog.Portal>
</Dialog.Root>

<!-- Company Address Dialog -->
<Dialog.Root bind:open={companyAddressDialogOpen}>
    <Dialog.Portal>
        <Dialog.Overlay class="data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 fixed inset-0 z-50 bg-black/80" />
        <Dialog.Content class="rounded-card-lg bg-background shadow-popover data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 data-[state=closed]:zoom-out-95 data-[state=open]:zoom-in-95 outline-hidden fixed left-[50%] top-[50%] z-50 w-full max-w-[calc(100%-2rem)] translate-x-[-50%] translate-y-[-50%] border p-8 sm:max-w-[540px] md:w-full">
            <Dialog.Title class="w-full text-xl font-semibold tracking-tight">
                Modifier l'adresse
            </Dialog.Title>
            <Dialog.Description class="text-foreground-alt !mt-1 text-sm">
                Définissez l'adresse légale de votre entreprise
            </Dialog.Description>

            <Separator.Root class="bg-muted mx-5 !mb-2 !mt-5 block h-px" />

            <form method="POST" action="?/updateCompanyAddress" use:companyAddressEnhance class="grid gap-4">
                <div class="grid gap-2 py-4">
                    <Label.Root for="company-address" class="text-sm font-medium">Adresse</Label.Root>
                    <input
                        id="company-address"
                        name="address"
                        type="text"
                        bind:value={$companyAddressForm.address}
                        placeholder="123 Wellness Boulevard, 75001 Paris"
                        class="h-input rounded-card-sm border-border-input bg-background placeholder:text-foreground-alt/50 hover:border-dark-40 focus:ring-foreground focus:ring-offset-background focus:outline-hidden inline-flex w-full items-center border px-4 text-base focus:ring-2 focus:ring-offset-2 sm:text-sm transition-all"
                    />
                </div>
                <div class="flex w-full justify-end gap-3">
                    <Button.Root type="button" class="cursor-pointer">
                        <Dialog.Close class="h-input rounded-input border border-border-input hover:bg-dark-04 focus-visible:ring-dark focus-visible:ring-offset-background focus-visible:outline-hidden inline-flex items-center justify-center px-6 text-sm font-medium focus-visible:ring-2 focus-visible:ring-offset-2 active:scale-[0.98] cursor-pointer transition-all">
                            Annuler
                        </Dialog.Close>
                    </Button.Root>
                    <Button.Root type="submit" class="cursor-pointer" disabled={$companyAddressDelayed}>
                        <div class="h-input rounded-input bg-dark text-background shadow-mini hover:bg-dark/95 focus-visible:ring-dark focus-visible:ring-offset-background focus-visible:outline-hidden inline-flex items-center justify-center px-8 text-sm font-semibold focus-visible:ring-2 focus-visible:ring-offset-2 active:scale-[0.98] cursor-pointer transition-all">
                            {$companyAddressDelayed ? "Enregistrement..." : "Enregistrer"}
                        </div>
                    </Button.Root>
                </div>
            </form>

            <Button.Root type="button" class="cursor-pointer">
                <Dialog.Close class="focus-visible:ring-foreground focus-visible:ring-offset-background focus-visible:outline-hidden absolute right-5 top-5 rounded-md focus-visible:ring-2 focus-visible:ring-offset-2 active:scale-[0.98] cursor-pointer">
                    <X class="text-foreground size-5" />
                    <span class="sr-only">Close</span>
                </Dialog.Close>
            </Button.Root>
        </Dialog.Content>
    </Dialog.Portal>
</Dialog.Root>

<!-- Company Instagram Dialog -->
<Dialog.Root bind:open={companyInstagramDialogOpen}>
    <Dialog.Portal>
        <Dialog.Overlay class="data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 fixed inset-0 z-50 bg-black/80" />
        <Dialog.Content class="rounded-card-lg bg-background shadow-popover data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 data-[state=closed]:zoom-out-95 data-[state=open]:zoom-in-95 outline-hidden fixed left-[50%] top-[50%] z-50 w-full max-w-[calc(100%-2rem)] translate-x-[-50%] translate-y-[-50%] border p-8 sm:max-w-[540px] md:w-full">
            <Dialog.Title class="w-full text-xl font-semibold tracking-tight">
                Modifier le lien Instagram
            </Dialog.Title>
            <Dialog.Description class="text-foreground-alt !mt-1 text-sm">
                Définissez l'URL de votre profil Instagram
            </Dialog.Description>

            <Separator.Root class="bg-muted mx-5 !mb-2 !mt-5 block h-px" />

            <form method="POST" action="?/updateCompanyInstagram" use:companyInstagramEnhance class="grid gap-4">
                <div class="grid gap-2 py-4">
                    <Label.Root for="company-instagram" class="text-sm font-medium">URL Instagram</Label.Root>
                    <input
                        id="company-instagram"
                        name="instagram"
                        type="url"
                        bind:value={$companyInstagramForm.instagram}
                        placeholder="https://instagram.com/leviosa_spa"
                        class="h-input rounded-card-sm border-border-input bg-background placeholder:text-foreground-alt/50 hover:border-dark-40 focus:ring-foreground focus:ring-offset-background focus:outline-hidden inline-flex w-full items-center border px-4 text-base focus:ring-2 focus:ring-offset-2 sm:text-sm transition-all"
                    />
                </div>
                <div class="flex w-full justify-end gap-3">
                    <Button.Root type="button" class="cursor-pointer">
                        <Dialog.Close class="h-input rounded-input border border-border-input hover:bg-dark-04 focus-visible:ring-dark focus-visible:ring-offset-background focus-visible:outline-hidden inline-flex items-center justify-center px-6 text-sm font-medium focus-visible:ring-2 focus-visible:ring-offset-2 active:scale-[0.98] cursor-pointer transition-all">
                            Annuler
                        </Dialog.Close>
                    </Button.Root>
                    <Button.Root type="submit" class="cursor-pointer" disabled={$companyInstagramDelayed}>
                        <div class="h-input rounded-input bg-dark text-background shadow-mini hover:bg-dark/95 focus-visible:ring-dark focus-visible:ring-offset-background focus-visible:outline-hidden inline-flex items-center justify-center px-8 text-sm font-semibold focus-visible:ring-2 focus-visible:ring-offset-2 active:scale-[0.98] cursor-pointer transition-all">
                            {$companyInstagramDelayed ? "Enregistrement..." : "Enregistrer"}
                        </div>
                    </Button.Root>
                </div>
            </form>

            <Button.Root type="button" class="cursor-pointer">
                <Dialog.Close class="focus-visible:ring-foreground focus-visible:ring-offset-background focus-visible:outline-hidden absolute right-5 top-5 rounded-md focus-visible:ring-2 focus-visible:ring-offset-2 active:scale-[0.98] cursor-pointer">
                    <X class="text-foreground size-5" />
                    <span class="sr-only">Close</span>
                </Dialog.Close>
            </Button.Root>
        </Dialog.Content>
    </Dialog.Portal>
</Dialog.Root>

<!-- OTP Duration Dialog -->
<Dialog.Root bind:open={otpDurationDialogOpen}>
    <Dialog.Portal>
        <Dialog.Overlay class="data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 fixed inset-0 z-50 bg-black/80" />
        <Dialog.Content class="rounded-card-lg bg-background shadow-popover data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 data-[state=closed]:zoom-out-95 data-[state=open]:zoom-in-95 outline-hidden fixed left-[50%] top-[50%] z-50 w-full max-w-[calc(100%-2rem)] translate-x-[-50%] translate-y-[-50%] border p-8 sm:max-w-[540px] md:w-full">
            <Dialog.Title class="w-full text-xl font-semibold tracking-tight">
                Modifier la durée de validité OTP
            </Dialog.Title>
            <Dialog.Description class="text-foreground-alt !mt-1 text-sm">
                Définissez combien de temps un code OTP reste valide (en secondes)
            </Dialog.Description>

            <Separator.Root class="bg-muted mx-5 !mb-2 !mt-5 block h-px" />

            <form method="POST" action="?/updateOtpDuration" use:otpDurationEnhance class="grid gap-4">
                <div class="grid gap-2 py-4">
                    <Label.Root for="otp-duration" class="text-sm font-medium">Durée (secondes)</Label.Root>
                    <input
                        id="otp-duration"
                        name="duration"
                        type="number"
                        min="60"
                        max="3600"
                        bind:value={$otpDurationForm.duration}
                        placeholder="300"
                        class="h-input rounded-card-sm border-border-input bg-background placeholder:text-foreground-alt/50 hover:border-dark-40 focus:ring-foreground focus:ring-offset-background focus:outline-hidden inline-flex w-full items-center border px-4 text-base focus:ring-2 focus:ring-offset-2 sm:text-sm transition-all"
                    />
                    <p class="text-xs text-foreground-alt">
                        Min: 60 secondes, Max: 3600 secondes
                    </p>
                </div>
                <div class="flex w-full justify-end gap-3">
                    <Button.Root type="button" class="cursor-pointer">
                        <Dialog.Close class="h-input rounded-input border border-border-input hover:bg-dark-04 focus-visible:ring-dark focus-visible:ring-offset-background focus-visible:outline-hidden inline-flex items-center justify-center px-6 text-sm font-medium focus-visible:ring-2 focus-visible:ring-offset-2 active:scale-[0.98] cursor-pointer transition-all">
                            Annuler
                        </Dialog.Close>
                    </Button.Root>
                    <Button.Root type="submit" class="cursor-pointer" disabled={$otpDurationDelayed}>
                        <div class="h-input rounded-input bg-dark text-background shadow-mini hover:bg-dark/95 focus-visible:ring-dark focus-visible:ring-offset-background focus-visible:outline-hidden inline-flex items-center justify-center px-8 text-sm font-semibold focus-visible:ring-2 focus-visible:ring-offset-2 active:scale-[0.98] cursor-pointer transition-all">
                            {$otpDurationDelayed ? "Enregistrement..." : "Enregistrer"}
                        </div>
                    </Button.Root>
                </div>
            </form>

            <Button.Root type="button" class="cursor-pointer">
                <Dialog.Close class="focus-visible:ring-foreground focus-visible:ring-offset-background focus-visible:outline-hidden absolute right-5 top-5 rounded-md focus-visible:ring-2 focus-visible:ring-offset-2 active:scale-[0.98] cursor-pointer">
                    <X class="text-foreground size-5" />
                    <span class="sr-only">Close</span>
                </Dialog.Close>
            </Button.Root>
        </Dialog.Content>
    </Dialog.Portal>
</Dialog.Root>

<!-- OTP Length Dialog -->
<Dialog.Root bind:open={otpLengthDialogOpen}>
    <Dialog.Portal>
        <Dialog.Overlay class="data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 fixed inset-0 z-50 bg-black/80" />
        <Dialog.Content class="rounded-card-lg bg-background shadow-popover data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 data-[state=closed]:zoom-out-95 data-[state=open]:zoom-in-95 outline-hidden fixed left-[50%] top-[50%] z-50 w-full max-w-[calc(100%-2rem)] translate-x-[-50%] translate-y-[-50%] border p-8 sm:max-w-[540px] md:w-full">
            <Dialog.Title class="w-full text-xl font-semibold tracking-tight">
                Modifier la longueur du code OTP
            </Dialog.Title>
            <Dialog.Description class="text-foreground-alt !mt-1 text-sm">
                Définissez le nombre de chiffres dans le code OTP
            </Dialog.Description>

            <Separator.Root class="bg-muted mx-5 !mb-2 !mt-5 block h-px" />

            <form method="POST" action="?/updateOtpLength" use:otpLengthEnhance class="grid gap-4">
                <div class="grid gap-2 py-4">
                    <Label.Root for="otp-length" class="text-sm font-medium">Longueur (chiffres)</Label.Root>
                    <input
                        id="otp-length"
                        name="length"
                        type="number"
                        min="4"
                        max="10"
                        bind:value={$otpLengthForm.length}
                        placeholder="6"
                        class="h-input rounded-card-sm border-border-input bg-background placeholder:text-foreground-alt/50 hover:border-dark-40 focus:ring-foreground focus:ring-offset-background focus:outline-hidden inline-flex w-full items-center border px-4 text-base focus:ring-2 focus:ring-offset-2 sm:text-sm transition-all"
                    />
                    <p class="text-xs text-foreground-alt">
                        Min: 4 chiffres, Max: 10 chiffres
                    </p>
                </div>
                <div class="flex w-full justify-end gap-3">
                    <Button.Root type="button" class="cursor-pointer">
                        <Dialog.Close class="h-input rounded-input border border-border-input hover:bg-dark-04 focus-visible:ring-dark focus-visible:ring-offset-background focus-visible:outline-hidden inline-flex items-center justify-center px-6 text-sm font-medium focus-visible:ring-2 focus-visible:ring-offset-2 active:scale-[0.98] cursor-pointer transition-all">
                            Annuler
                        </Dialog.Close>
                    </Button.Root>
                    <Button.Root type="submit" class="cursor-pointer" disabled={$otpLengthDelayed}>
                        <div class="h-input rounded-input bg-dark text-background shadow-mini hover:bg-dark/95 focus-visible:ring-dark focus-visible:ring-offset-background focus-visible:outline-hidden inline-flex items-center justify-center px-8 text-sm font-semibold focus-visible:ring-2 focus-visible:ring-offset-2 active:scale-[0.98] cursor-pointer transition-all">
                            {$otpLengthDelayed ? "Enregistrement..." : "Enregistrer"}
                        </div>
                    </Button.Root>
                </div>
            </form>

            <Button.Root type="button" class="cursor-pointer">
                <Dialog.Close class="focus-visible:ring-foreground focus-visible:ring-offset-background focus-visible:outline-hidden absolute right-5 top-5 rounded-md focus-visible:ring-2 focus-visible:ring-offset-2 active:scale-[0.98] cursor-pointer">
                    <X class="text-foreground size-5" />
                    <span class="sr-only">Close</span>
                </Dialog.Close>
            </Button.Root>
        </Dialog.Content>
    </Dialog.Portal>
</Dialog.Root>

<!-- OTP Max Attempts Dialog -->
<Dialog.Root bind:open={otpMaxAttemptsDialogOpen}>
    <Dialog.Portal>
        <Dialog.Overlay class="data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 fixed inset-0 z-50 bg-black/80" />
        <Dialog.Content class="rounded-card-lg bg-background shadow-popover data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 data-[state=closed]:zoom-out-95 data-[state=open]:zoom-in-95 outline-hidden fixed left-[50%] top-[50%] z-50 w-full max-w-[calc(100%-2rem)] translate-x-[-50%] translate-y-[-50%] border p-8 sm:max-w-[540px] md:w-full">
            <Dialog.Title class="w-full text-xl font-semibold tracking-tight">
                Modifier les tentatives maximales
            </Dialog.Title>
            <Dialog.Description class="text-foreground-alt !mt-1 text-sm">
                Définissez le nombre de tentatives avant le verrouillage du compte
            </Dialog.Description>

            <Separator.Root class="bg-muted mx-5 !mb-2 !mt-5 block h-px" />

            <form method="POST" action="?/updateOtpMaxAttempts" use:otpMaxAttemptsEnhance class="grid gap-4">
                <div class="grid gap-2 py-4">
                    <Label.Root for="otp-max-attempts" class="text-sm font-medium">Tentatives maximales</Label.Root>
                    <input
                        id="otp-max-attempts"
                        name="max_attempts"
                        type="number"
                        min="1"
                        max="10"
                        bind:value={$otpMaxAttemptsForm.max_attempts}
                        placeholder="5"
                        class="h-input rounded-card-sm border-border-input bg-background placeholder:text-foreground-alt/50 hover:border-dark-40 focus:ring-foreground focus:ring-offset-background focus:outline-hidden inline-flex w-full items-center border px-4 text-base focus:ring-2 focus:ring-offset-2 sm:text-sm transition-all"
                    />
                    <p class="text-xs text-foreground-alt">
                        Min: 1 tentative, Max: 10 tentatives
                    </p>
                </div>
                <div class="flex w-full justify-end gap-3">
                    <Button.Root type="button" class="cursor-pointer">
                        <Dialog.Close class="h-input rounded-input border border-border-input hover:bg-dark-04 focus-visible:ring-dark focus-visible:ring-offset-background focus-visible:outline-hidden inline-flex items-center justify-center px-6 text-sm font-medium focus-visible:ring-2 focus-visible:ring-offset-2 active:scale-[0.98] cursor-pointer transition-all">
                            Annuler
                        </Dialog.Close>
                    </Button.Root>
                    <Button.Root type="submit" class="cursor-pointer" disabled={$otpMaxAttemptsDelayed}>
                        <div class="h-input rounded-input bg-dark text-background shadow-mini hover:bg-dark/95 focus-visible:ring-dark focus-visible:ring-offset-background focus-visible:outline-hidden inline-flex items-center justify-center px-8 text-sm font-semibold focus-visible:ring-2 focus-visible:ring-offset-2 active:scale-[0.98] cursor-pointer transition-all">
                            {$otpMaxAttemptsDelayed ? "Enregistrement..." : "Enregistrer"}
                        </div>
                    </Button.Root>
                </div>
            </form>

            <Button.Root type="button" class="cursor-pointer">
                <Dialog.Close class="focus-visible:ring-foreground focus-visible:ring-offset-background focus-visible:outline-hidden absolute right-5 top-5 rounded-md focus-visible:ring-2 focus-visible:ring-offset-2 active:scale-[0.98] cursor-pointer">
                    <X class="text-foreground size-5" />
                    <span class="sr-only">Close</span>
                </Dialog.Close>
            </Button.Root>
        </Dialog.Content>
    </Dialog.Portal>
</Dialog.Root>

<!-- Access Token Duration Dialog -->
<Dialog.Root bind:open={accessTokenDurationDialogOpen}>
    <Dialog.Portal>
        <Dialog.Overlay class="data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 fixed inset-0 z-50 bg-black/80" />
        <Dialog.Content class="rounded-card-lg bg-background shadow-popover data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 data-[state=closed]:zoom-out-95 data-[state=open]:zoom-in-95 outline-hidden fixed left-[50%] top-[50%] z-50 w-full max-w-[calc(100%-2rem)] translate-x-[-50%] translate-y-[-50%] border p-8 sm:max-w-[540px] md:w-full">
            <Dialog.Title class="w-full text-xl font-semibold tracking-tight">
                Modifier la durée du token d'accès
            </Dialog.Title>
            <Dialog.Description class="text-foreground-alt !mt-1 text-sm">
                Définissez combien de temps un token d'accès reste valide (en minutes)
            </Dialog.Description>

            <Separator.Root class="bg-muted mx-5 !mb-2 !mt-5 block h-px" />

            <form method="POST" action="?/updateAccessTokenDuration" use:accessTokenDurationEnhance class="grid gap-4">
                <div class="grid gap-2 py-4">
                    <Label.Root for="access-token-duration" class="text-sm font-medium">Durée (minutes)</Label.Root>
                    <input
                        id="access-token-duration"
                        name="duration"
                        type="number"
                        min="1"
                        max="240"
                        bind:value={$accessTokenDurationForm.duration}
                        placeholder="15"
                        class="h-input rounded-card-sm border-border-input bg-background placeholder:text-foreground-alt/50 hover:border-dark-40 focus:ring-foreground focus:ring-offset-background focus:outline-hidden inline-flex w-full items-center border px-4 text-base focus:ring-2 focus:ring-offset-2 sm:text-sm transition-all"
                    />
                    <p class="text-xs text-foreground-alt">
                        Min: 1 minute, Max: 240 minutes (4 heures)
                    </p>
                </div>
                <div class="flex w-full justify-end gap-3">
                    <Button.Root type="button" class="cursor-pointer">
                        <Dialog.Close class="h-input rounded-input border border-border-input hover:bg-dark-04 focus-visible:ring-dark focus-visible:ring-offset-background focus-visible:outline-hidden inline-flex items-center justify-center px-6 text-sm font-medium focus-visible:ring-2 focus-visible:ring-offset-2 active:scale-[0.98] cursor-pointer transition-all">
                            Annuler
                        </Dialog.Close>
                    </Button.Root>
                    <Button.Root type="submit" class="cursor-pointer" disabled={$accessTokenDurationDelayed}>
                        <div class="h-input rounded-input bg-dark text-background shadow-mini hover:bg-dark/95 focus-visible:ring-dark focus-visible:ring-offset-background focus-visible:outline-hidden inline-flex items-center justify-center px-8 text-sm font-semibold focus-visible:ring-2 focus-visible:ring-offset-2 active:scale-[0.98] cursor-pointer transition-all">
                            {$accessTokenDurationDelayed ? "Enregistrement..." : "Enregistrer"}
                        </div>
                    </Button.Root>
                </div>
            </form>

            <Button.Root type="button" class="cursor-pointer">
                <Dialog.Close class="focus-visible:ring-foreground focus-visible:ring-offset-background focus-visible:outline-hidden absolute right-5 top-5 rounded-md focus-visible:ring-2 focus-visible:ring-offset-2 active:scale-[0.98] cursor-pointer">
                    <X class="text-foreground size-5" />
                    <span class="sr-only">Close</span>
                </Dialog.Close>
            </Button.Root>
        </Dialog.Content>
    </Dialog.Portal>
</Dialog.Root>

<!-- Refresh Token Duration Dialog -->
<Dialog.Root bind:open={refreshTokenDurationDialogOpen}>
    <Dialog.Portal>
        <Dialog.Overlay class="data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 fixed inset-0 z-50 bg-black/80" />
        <Dialog.Content class="rounded-card-lg bg-background shadow-popover data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 data-[state=closed]:zoom-out-95 data-[state=open]:zoom-in-95 outline-hidden fixed left-[50%] top-[50%] z-50 w-full max-w-[calc(100%-2rem)] translate-x-[-50%] translate-y-[-50%] border p-8 sm:max-w-[540px] md:w-full">
            <Dialog.Title class="w-full text-xl font-semibold tracking-tight">
                Modifier la durée du token de rafraîchissement
            </Dialog.Title>
            <Dialog.Description class="text-foreground-alt !mt-1 text-sm">
                Définissez combien de temps un token de rafraîchissement reste valide (en heures)
            </Dialog.Description>

            <Separator.Root class="bg-muted mx-5 !mb-2 !mt-5 block h-px" />

            <form method="POST" action="?/updateRefreshTokenDuration" use:refreshTokenDurationEnhance class="grid gap-4">
                <div class="grid gap-2 py-4">
                    <Label.Root for="refresh-token-duration" class="text-sm font-medium">Durée (heures)</Label.Root>
                    <input
                        id="refresh-token-duration"
                        name="duration"
                        type="number"
                        min="1"
                        max="720"
                        bind:value={$refreshTokenDurationForm.duration}
                        placeholder="168"
                        class="h-input rounded-card-sm border-border-input bg-background placeholder:text-foreground-alt/50 hover:border-dark-40 focus:ring-foreground focus:ring-offset-background focus:outline-hidden inline-flex w-full items-center border px-4 text-base focus:ring-2 focus:ring-offset-2 sm:text-sm transition-all"
                    />
                    <p class="text-xs text-foreground-alt">
                        Min: 1 heure, Max: 720 heures (30 jours)
                    </p>
                </div>
                <div class="flex w-full justify-end gap-3">
                    <Button.Root type="button" class="cursor-pointer">
                        <Dialog.Close class="h-input rounded-input border border-border-input hover:bg-dark-04 focus-visible:ring-dark focus-visible:ring-offset-background focus-visible:outline-hidden inline-flex items-center justify-center px-6 text-sm font-medium focus-visible:ring-2 focus-visible:ring-offset-2 active:scale-[0.98] cursor-pointer transition-all">
                            Annuler
                        </Dialog.Close>
                    </Button.Root>
                    <Button.Root type="submit" class="cursor-pointer" disabled={$refreshTokenDurationDelayed}>
                        <div class="h-input rounded-input bg-dark text-background shadow-mini hover:bg-dark/95 focus-visible:ring-dark focus-visible:ring-offset-background focus-visible:outline-hidden inline-flex items-center justify-center px-8 text-sm font-semibold focus-visible:ring-2 focus-visible:ring-offset-2 active:scale-[0.98] cursor-pointer transition-all">
                            {$refreshTokenDurationDelayed ? "Enregistrement..." : "Enregistrer"}
                        </div>
                    </Button.Root>
                </div>
            </form>

            <Button.Root type="button" class="cursor-pointer">
                <Dialog.Close class="focus-visible:ring-foreground focus-visible:ring-offset-background focus-visible:outline-hidden absolute right-5 top-5 rounded-md focus-visible:ring-2 focus-visible:ring-offset-2 active:scale-[0.98] cursor-pointer">
                    <X class="text-foreground size-5" />
                    <span class="sr-only">Close</span>
                </Dialog.Close>
            </Button.Root>
        </Dialog.Content>
    </Dialog.Portal>
</Dialog.Root>
