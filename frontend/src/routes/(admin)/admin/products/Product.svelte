<script lang="ts">
    import { Plus } from "@lucide/svelte";

    import ProductCard from "./ProductCard.svelte";
    import SearchBar from "./SearchBar.svelte";
    import Select from "./Select.svelte";
    import ProductModal from "./ProductModal.svelte";

    import { type CardType, type Category } from "./products";

    import { type SuperValidated } from "sveltekit-superforms";

    import type { DeleteProduct, product } from "./schemas";

    type Props = {
        cards: CardType[];
        statuses: Set<string>;
        categories: Category[];
        availabilities: Set<string>;
        deleteProductForm: SuperValidated<DeleteProduct>;
        createProductForm: SuperValidated<product>;
        updateProductForm: SuperValidated<product>;
    };

    let {
        cards,
        statuses,
        categories,
        availabilities,
        deleteProductForm,
        createProductForm,
        updateProductForm,
    }: Props = $props();

    import {
        defaultStatus,
        defaultCategory,
        defaultAvailability,
    } from "./default";

    // filters
    let status = $state(defaultStatus);
    let category = $state(defaultCategory);
    let availability = $state(defaultAvailability);

    let searchValue = $state("");

    let filteredCards = $derived(
        cards
            .filter(
                (card) => status === defaultStatus || card.published === status,
            )
            .filter(
                (card) =>
                    // category === defaultCategory || card.category === category,
                    category === defaultCategory || card.category === category,
            )
            .filter(
                (card) =>
                    availability === defaultAvailability ||
                    card.published === availability,
            )
            // TODO: how to have the selection of an option defining the displayed card ?
            .filter((card) =>
                card.name.toLowerCase().includes(searchValue.toLowerCase()),
            ),
    );

    // NOTE: here is the debug part that I need to remove later when the form works
    import SuperDebug from "sveltekit-superforms";
    import { superForm } from "sveltekit-superforms";

    const { form, errors } = superForm(createProductForm, {
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
        <h2 class="text-3xl font-bold">Products</h2>
        <p class="">
            Créer et gérer les services disponibles à la réservation sur la
            plateforme.
        </p>
    </div>
    <ProductModal
        {statuses}
        {categories}
        {availabilities}
        modalForm={createProductForm}
    >
        <div
            class="flex gap-2 items-center py-2 px-4 bg-green-600 text-white rounded-md cursor-pointer"
        >
            <Plus />
            <p>Nouveau Produit</p>
        </div>
    </ProductModal>
</div>
<div
    class="px-8 border-b rounded-lg py-4 flex justify-between items-center gap-4 bg-white"
>
    <div class="flex-1">
        <SearchBar {cards} bind:searchValue />
    </div>
    <Select name="status" items={statuses} bind:state={status} />
    <Select
        name="category"
        items={new Set(categories.map((category) => category.name))}
        bind:state={category}
    />
    <Select
        name="availability"
        items={availabilities}
        bind:state={availability}
    />
</div>
<div
    class="px-8 grid mt-8 grid-cols-[repeat(auto-fit,minmax(320px,1fr))] gap-4 max-sm:justify-self-center"
>
    {#each filteredCards as card (card.id)}
        <ProductCard
            {card}
            {statuses}
            {categories}
            {availabilities}
            {deleteProductForm}
            {updateProductForm}
        />
    {/each}
</div>
