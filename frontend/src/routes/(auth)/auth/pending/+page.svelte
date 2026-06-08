<script lang="ts">
    import { goto } from "$app/navigation";
    import { Button } from "bits-ui";
    import { Slack, ArrowRight } from "@lucide/svelte";
</script>

<div class="grid items-center h-[100vh] px-6">
    <div class="grid gap-8">
        <div class="w-10 mx-auto aspect-square grid items-center rounded-lg">
            <Slack size={24} />
        </div>
        <div class="grid gap-4">
            <h3 class="text-foreground font-semibold text-lg">
                Merci pour votre inscription
            </h3>
            <div>
                <p class="font-medium">
                    Votre inscription est en attente de validation par
                    un membre de l'équipe Leviosa. Vous recevrez un e-mail
                    dès que votre compte sera activé.
                </p>
            </div>
        </div>
    </div>
    <Button.Root class="mt-12" onclick={() => goto("/")}>
        Accueil
        <ArrowRight />
    </Button.Root>
</div>

<style>
</style>
