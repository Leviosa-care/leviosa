<script lang="ts">
    import { reveal } from "$lib/actions/reveal";
    import { slide } from "svelte/transition";
    import { Plus } from "@lucide/svelte";

    const faqs = [
        {
            question: "Comment fonctionne la prise de rendez-vous ?",
            answer: "Choisissez votre expert, sélectionnez un créneau disponible dans son agenda et confirmez en quelques clics. Vous recevrez une confirmation par email avec tous les détails de votre consultation.",
        },
        {
            question: "Puis-je annuler ou modifier ma réservation ?",
            answer: "Oui, vous pouvez annuler ou reprogrammer votre consultation jusqu'à 24 heures avant le rendez-vous depuis votre espace client. Au-delà de ce délai, des frais d'annulation peuvent s'appliquer selon la politique de l'expert.",
        },
        {
            question: "Les consultations se déroulent-elles en ligne ou en présentiel ?",
            answer: "Cela dépend de l'expert et du service choisi. Certains praticiens proposent uniquement des séances en cabinet, d'autres uniquement en visio, et certains offrent les deux options. Ces informations sont clairement indiquées sur chaque fiche de service.",
        },
        {
            question: "Comment les experts sont-ils vérifiés ?",
            answer: "Chaque praticien passe par un processus de vérification rigoureux : contrôle des diplômes et certifications, vérification des assurances professionnelles, et entretien de qualification. Seuls les experts répondant à nos critères sont acceptés sur la plateforme.",
        },
        {
            question: "Quels modes de paiement sont acceptés ?",
            answer: "Nous acceptons les principales cartes bancaires (Visa, Mastercard, American Express) via notre partenaire de paiement sécurisé Stripe. Le montant est débité au moment de la confirmation de la réservation.",
        },
        {
            question: "Que se passe-t-il après ma consultation ?",
            answer: "Après votre séance, vous recevrez un compte-rendu personnalisé de votre expert. Vous pourrez également retrouver l'historique de toutes vos consultations dans votre espace client et reprendre rendez-vous en un clic.",
        },
    ];

    let openIndex = $state<number | null>(0);

    function toggle(i: number) {
        openIndex = openIndex === i ? null : i;
    }
</script>

<section class="py-20 lg:py-24 bg-white">
    <div class="max-w-3xl mx-auto px-4 sm:px-6 lg:px-8">
        <!-- Section Header -->
        <div
            class="text-center mb-12"
            use:reveal={{ preset: "fade-up", delay: 100 }}
        >
            <span
                class="text-sm font-semibold text-dark-400 uppercase tracking-wider"
                >FAQ</span
            >
            <h2
                class="mt-3 text-3xl md:text-4xl font-semibold tracking-tight text-dark-900"
            >
                Questions fréquentes
            </h2>
            <p
                class="mt-4 text-base text-dark-500 leading-relaxed font-normal"
            >
                Tout ce que vous devez savoir avant de prendre rendez-vous.
            </p>
        </div>

        <!-- Accordion -->
        <div
            class="divide-y divide-dark-100"
            use:reveal={{ preset: "fade-up", delay: 150 }}
        >
            {#each faqs as faq, i}
                <div class="py-5">
                    <button
                        class="w-full flex items-center justify-between gap-4 text-left cursor-pointer group"
                        onclick={() => toggle(i)}
                    >
                        <span
                            class="text-base font-medium text-dark-900 group-hover:text-dark-600 transition-colors"
                        >
                            {faq.question}
                        </span>
                        <span
                            class="flex-shrink-0 w-6 h-6 rounded-full bg-dark-50 border border-dark-100 flex items-center justify-center transition-transform duration-200 {openIndex ===
                            i
                                ? 'rotate-45'
                                : ''}"
                        >
                            <Plus size={14} class="text-dark-500" strokeWidth={2} />
                        </span>
                    </button>
                    {#if openIndex === i}
                        <div
                            transition:slide={{ duration: 200 }}
                            class="mt-3 pr-10"
                        >
                            <p class="text-sm text-dark-500 leading-relaxed">
                                {faq.answer}
                            </p>
                        </div>
                    {/if}
                </div>
            {/each}
        </div>
    </div>
</section>
