<script lang="ts">
    import { TriangleAlert, ArrowRight } from "@lucide/svelte";
</script>

<div class="h-fit rounded-xl overflow-hidden border border-border-card">
    <div
        class="bg-dark-50 px-4 py-4 flex gap-4 items-center border border-b-dark-50"
    >
        <TriangleAlert />
        <p class="text-dark-900 font-bold">Actions requises</p>
    </div>
    <div>
        {@render required_action(
            "Validation practicien",
            "2h",
            "Dr Thomas R. attend la validation de ses documents.",
            "Examiner",
        )}
        {@render required_action(
            "Profil incomplet",
            "hier",
            "Cabinet KineSport a des informations manquantes.",
            "Relancer",
        )}
        {@render required_action(
            "Litige client",
            "1j",
            "Signalement sur la consultation #4849",
            "Voir le detail",
        )}
    </div>
</div>

{#snippet required_action(
    title: string,
    duration: string,
    description: string,
    action: string,
)}
    <div class="p-4 grid gap-4 bg-background">
        <div class="grid gap-1">
            <div class="flex justify-between items-center">
                <p class="text-dark-800 font-bold">{title}</p>
                <p class="text-dark-300 capitalize">{duration}</p>
            </div>
            <p class="text-dark-500">{description}</p>
        </div>
        <button
            class="flex gap-2 items-center text-dark-700 hover:text-dark-900 cursor-pointer"
        >
            <p>{action}</p>
            <ArrowRight size={12} />
        </button>
    </div>
{/snippet}
