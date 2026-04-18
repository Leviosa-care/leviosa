<script lang="ts">
    import { type SuperValidated } from "sveltekit-superforms";
    import { superForm } from "sveltekit-superforms";
    import SuperDebug from "sveltekit-superforms";

    import { Plus } from "@lucide/svelte";

    import type { category } from "./schemas";
    import CategoryModal from "./CategoryModal.svelte";

    type Props = {
        createCategoryForm: SuperValidated<category>;
    };

    let { createCategoryForm }: Props = $props();

    const { form, errors } = superForm(createCategoryForm, {
        resetForm: true,
    });
</script>

<SuperDebug
    data={{
        form: $form,
        errors: $errors,
    }}
/>

<div
    class="px-8 py-4 flex justify-between items-center border-b border-gray-200 bg-white"
>
    <div class="grid gap-1">
        <h2 class="text-3xl font-bold">Categories</h2>
        <p class="">
            Créer et gérer les categories de services disponibles à la
            réservation sur la plateforme.
        </p>
    </div>
    <CategoryModal modalForm={createCategoryForm}>
        <div
            class="flex gap-2 items-center py-2 px-4 bg-green-600 text-white rounded-md cursor-pointer"
        >
            <Plus />
            <p>Nouvelle Categorie</p>
        </div>
    </CategoryModal>
</div>
