<script lang="ts">
    import { Button } from "bits-ui";
</script>

<div
    class="bg-white h-fit border border-dark-100/40 rounded-xl overflow-hidden"
>
    <div class="px-8 py-6 flex gap-4 items-center border border-b-dark-50">
        <p class="text-dark-900 font-bold">Agenda du jour</p>
    </div>
    <div class="grid">
        {@render time_slot(
            "10:00",
            "Consultations suivi - Dermatologie",
            "Dr Sarah L",
            "Marc D.",
            "En cours",
        )}
        {@render time_slot(
            "10:00",
            "Bilan initial - Nutrition",
            "Dr Jean M.",
            "Sophie K.",
            "A venir",
        )}
        {@render time_slot(
            "10:00",
            "Bilan initial - Nutrition",
            "Dr Jean M.",
            "Sophie K.",
            "Annule",
        )}
        <Button.Root
            class="bg-dark-50/40 hover:bg-dark-50/60 text-dark-800 font-medium inline-flex justify-center items-center w-full border border-dark-50 py-4 cursor-pointer"
            >Voir l'agenda complet</Button.Root
        >
    </div>
</div>

{#snippet time_slot(
    time: string,
    title: string,
    partner: string,
    client: string,
    status: string,
)}
    <div class="px-8 py-4 flex justify-between items-center">
        <div class="flex items-center gap-8">
            <div
                class="p-3 bg-dark-50 aspect-square inline-flex justify-center items-center rounded-input text-sm font-medium"
            >
                {time}
            </div>
            <div>
                <p class="text-dark-900 font-medium">
                    {title}
                </p>
                <p class="text-dark-600 text-sm">
                    {partner} - Client: {client}
                </p>
            </div>
        </div>
        <p>{status}</p>
    </div>
{/snippet}
