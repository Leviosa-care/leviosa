<script lang="ts">
    import type { PageProps } from "./$types";

    import { ArrowRight, Slack } from "@lucide/svelte";
    import { Tabs, Button } from "bits-ui";
    import { superForm } from "sveltekit-superforms";

    import Input from "$lib/ui/Input.svelte";
    import Alert from "$lib/ui/Alert.svelte";
    import { getRegistrationContext } from "$lib/context/register.svelte";
    import { buttons, type OAuthButton } from "./oauth_buttons";

    let { data }: PageProps = $props();

    import { goto } from "$app/navigation";
    import { page } from "$app/state";
    const redirectFrom = page.url.searchParams.get("redirectFrom");
    import { MESSAGES } from "$lib/utils/redirect";
    import { capitalizeFirstWord } from "$lib/utils/capitalize";


    const { form, errors, enhance, constraints } = superForm(data.loginForm, {
        resetForm: true,
    });

    const {
        form: registerForm,
        errors: registerErrors,
        enhance: registerEnhance,
        constraints: registerConstraints,
    } = superForm(data.registerForm, {
        resetForm: true,
        onUpdated({ form }) {
            if (form.valid) {
                const register = getRegistrationContext();
                register.setEmail(form.data.registerEmail);
                goto("/auth/verify-email");
            }
        },
    });
    // import SuperDebug from "sveltekit-superforms";

    type AuthState = "login" | "register";
    let authState = $state<AuthState>("login");
    type OauthState = "google" | "apple";
    let oauthState = $state<OauthState>("google");
</script>

<!-- <SuperDebug -->
<!--     data={{ -->
<!--         form: $registerForm.registerEmail, -->
<!--         errors: $registerErrors.registerEmail, -->
<!--         constraints: $registerConstraints, -->
<!--         state: authState, -->
<!--     }} -->
<!-- /> -->

<div class="grid place-content-center h-[100vh]">
    <div class="pt-6">
        <Tabs.Root
            bind:value={authState}
            class="rounded-card border-muted bg-background-alt shadow-card md:w-[540px] w-[360px] border p-8"
        >
            <Tabs.Content value="login" class="select-none pt-8">
                {@render header("Connectez-vous à votre compte pour continuer")}
                {#if redirectFrom && !$errors._errors}
                    <Alert message={MESSAGES.expiredSession} />
                {/if}
                {#if $errors._errors}
                    <Alert message={$errors._errors[0]} />
                {/if}
            </Tabs.Content>
            <Tabs.Content value="register" class="select-none pt-8">
                {@render header(
                    "Inscrivez-vous pour profiter de tous nos services",
                )}
                {#if $registerErrors._errors}
                    <Alert type="error" message={$registerErrors._errors[0]} />
                {/if}
            </Tabs.Content>
            <Tabs.List
                class="rounded-9px bg-dark-10 shadow-mini-inset dark:bg-background grid w-full grid-cols-2 gap-1 p-1 text-sm font-semibold leading-[0.01em] dark:border dark:border-border-input"
            >
                <Tabs.Trigger
                    value="register"
                    class="data-[state=active]:shadow-mini dark:data-[state=active]:bg-muted h-8 rounded-[7px] bg-transparent py-2 data-[state=active]:bg-background cursor-pointer"
                    >S'enregistrer</Tabs.Trigger
                >
                <Tabs.Trigger
                    value="login"
                    class="data-[state=active]:shadow-mini dark:data-[state=active]:bg-muted h-8 rounded-[7px] bg-transparent py-2 data-[state=active]:bg-background cursor-pointer"
                    >Se connecter</Tabs.Trigger
                >
            </Tabs.List>
            <Tabs.Content value="login" class="select-none pt-8">
                <form
                    method="POST"
                    action="?/login"
                    class="grid gap-4"
                    use:enhance
                >
                    <div class="grid gap-6">
                        <div>
                            <Input
                                label="email"
                                name="email"
                                type="email"
                                autocomplete="email"
                                bind:value={$form.email}
                                error={$errors.email ? $errors.email[0] : ""}
                                {...$constraints.email}
                            />
                        </div>
                        <div>
                            <Input
                                name="password"
                                label="Mot de passe"
                                type="password"
                                autocomplete="current-password"
                                bind:value={$form.password}
                                error={$errors.password
                                    ? $errors.password[0]
                                    : ""}
                                showForgotPasswordbutton={true}
                                {...$constraints.password}
                            />
                        </div>
                    </div>
                    {@render button("Se connecter", ArrowRight)}
                </form>
            </Tabs.Content>
            <Tabs.Content value="register" class="select-none pt-6">
                <form
                    method="POST"
                    action="?/register"
                    class="grid grid-rows-2 gap-4"
                    use:registerEnhance
                >
                    <div class="grid gap-4">
                        <Input
                            name="registerEmail"
                            label="email"
                            type="email"
                            autocomplete="email"
                            bind:value={$registerForm.registerEmail}
                            {...$registerConstraints.registerEmail}
                            error={$registerForm.registerEmail
                                ? $registerForm.registerEmail[0]
                                : ""}
                        />
                    </div>
                    {@render button("Continuer votre inscription", ArrowRight)}
                </form>
            </Tabs.Content>
            {@render oauthButtons(buttons)}
        </Tabs.Root>
    </div>
</div>

{#snippet header(subtitle: string)}
    <div class="text-center mb-8 grid gap-4">
        <div class="flex gap-2 justify-center items-center">
            <Slack />
            <h2 class="font-extrabold text-2xl">Bienvenue à Leviosa</h2>
        </div>
        <p class="text-muted-foreground">
            {subtitle}
        </p>
    </div>
{/snippet}

{#snippet button(content: string, icon: typeof import("@lucide/svelte").Icon)}
    {@const Icon = icon}
    <Button.Root
        class="mt-6 justify-center gap-4 items-center h-input rounded-input bg-dark text-background shadow-mini hover:bg-dark/95 focus-visible:ring-dark focus-visible:ring-offset-background focus-visible:outline-hidden inline-flex px-4 text-[15px] font-bold focus-visible:ring-2 focus-visible:ring-offset-2 active:scale-[0.98] cursor-pointer"
        type="submit"
    >
        {content}
        <Icon size={20} />
    </Button.Root>
{/snippet}

{#snippet oauthButtons(buttons: OAuthButton[])}
    {#snippet oauthButton({ icon, name }: OAuthButton)}
        {@const Icon = icon}
        {@const prefix =
            authState === "login" ? "Se connecter" : "S'enregistrer"}
        <Button.Root
            type="submit"
            class="gap-4 items-center justify-center h-input rounded-input bg-transparent text-dark hover:bg-[#fafafa] border-border-input border-1 focus-visible:ring-dark focus-visible:ring-offset-background focus-visible:outline-hidden inline-flex px-4 text-[15px] font-semibold focus-visible:ring-2 focus-visible:ring-offset-2 active:scale-[0.98]"
        >
            <div class="size-5">
                <Icon />
            </div>
            {prefix} avec {capitalizeFirstWord(name)}
        </Button.Root>
    {/snippet}
    <form
        method="POST"
        action="?/oauth"
        class="grid grid-rows-{buttons.length} gap-2 mt-8"
    >
        <input type="hidden" name="provider" bind:value={oauthState} />
        {#each buttons as button}
            {@render oauthButton(button)}
        {/each}
    </form>
    {#if authState === "register"}
        <p class="text-sm text-center mt-2 text-muted-foreground">
            En continuant, vous acceptez nos <a
                class="underline hover:text-foreground"
                href="/terms-of-service">Conditions d'utilisation</a
            >
            et notre
            <a class="underline hover:text-foreground" href="/privacy-policy"
                >Politique de confidentialité</a
            >.
        </p>
    {/if}
{/snippet}
