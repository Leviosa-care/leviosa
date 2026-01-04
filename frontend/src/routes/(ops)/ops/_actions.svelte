<script lang="ts">
    import { TriangleAlert, ArrowRight } from "@lucide/svelte";
</script>

<div class="bg-white h-fit rounded-xl overflow-hidden">
    <div class="px-4 py-6 flex gap-4 items-center border border-b-dark-50">
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
            "Validation practicien",
            "hier",
            "Dr Thomas R. attend la validation de ses documents.",
            "Relancer",
        )}
        {@render required_action(
            "Validation practicien",
            "1j",
            "Dr Thomas R. attend la validation de ses documents.",
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
    <div class="p-4">
        <div class="flex justify-between items-center">
            <p class="text-dark-800 font-bold">{title}</p>
            <p class="text-dark-300 capitalize">{duration}</p>
        </div>
        <p class="text-dark-600">{description}</p>
        <button class="flex gap-2 items-center text-blue-900">
            <p>{action}</p>
            <ArrowRight size={12} />
        </button>
    </div>
{/snippet}
