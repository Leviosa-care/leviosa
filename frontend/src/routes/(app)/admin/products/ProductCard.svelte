<script lang="ts">
    import { Button } from "bits-ui";
    import { SquarePen, Trash2, Clock, Archive } from "@lucide/svelte";
    import { type SuperValidated } from "sveltekit-superforms";
    import { superForm } from "sveltekit-superforms";

    import { capitalizeFirstWord } from "$lib/utils/capitalize";
    import { type Card, type Category } from "./products";
    import ProductModal from "./ProductModal.svelte";
    import AlertDialog from "./AlertDialog.svelte";
    import { type DeleteProduct, type product } from "./schemas";

    type Props = {
        card: Card;
        statuses: Set<string>;
        categories: Category[];
        availabilities: Set<string>;
        deleteProductForm: SuperValidated<DeleteProduct>;
        updateProductForm: SuperValidated<product>;
    };

    let {
        card,
        statuses,
        categories,
        availabilities,
        deleteProductForm,
        updateProductForm,
    }: Props = $props();
    let {
        id,
        name,
        price,
        category,
        description,
        duration,
        image,
        updatedAt,
        published,
        // availability, // TODO: use that one later to explain where I can access that service
    } = card;

    const {
        form: deleteForm,
        // errors: deleteErrors,
        enhance: deleteEnhance,
        // constraints: deleteConstraints,
        // TODO: problem here supposedly
    } = superForm(deleteProductForm, {
        resetForm: true,
    });

    function formatMinutes(totalMinutes: number) {
        const hours = Math.floor(totalMinutes / 60);
        const minutes = totalMinutes % 60;

        let result = "";
        if (hours > 0) {
            result += `${hours}h`;
        }
        if (minutes > 0 || hours === 0) {
            result += `${minutes}min`;
        }

        return result;
    }

    const statusClassMap = {
        published: "bg-dark text-contrast",
        draft: "bg-tertiary/10 text-tertiary",
        archived: "bg-destructive/10 text-destructive",
    };

    $effect(() => {
        $deleteForm.id = card.id;
    });

    // import SuperDebug from "sveltekit-superforms";
</script>

<!-- <SuperDebug -->
<!--     data={{ -->
<!--         form: $deleteForm.id, -->
<!--     }} -->
<!-- /> -->

<div
    {id}
    class="relative max-w-[360px] border border-input rounded-md grid gap-4 bg-background"
>
    <div class="cursor-pointer grid gap-1">
        <img src={image} alt="service illustration" />
        <div class="grid gap-2 px-4 pb-4">
            <p class="font-bold text-xl">{name}</p>
            <div class="flex justify-between items-center bg-grey">
                <p class="bg-muted px-4 py-1 rounded-lg font-black text-sm">
                    {capitalizeFirstWord(category)}
                </p>
                <p
                    class="text-sm px-4 py-1 font-black rounded-lg {statusClassMap[
                        published
                    ]}"
                >
                    {capitalizeFirstWord(published)}
                </p>
            </div>
            <div class="flex items-center justify-between">
                <p
                    class="font-black text-xl {published === 'published'
                        ? 'text-accent-foreground'
                        : 'text-foreground'}"
                >
                    {price}€
                </p>
                <div class="flex items-center gap-2 text-muted-foreground">
                    <Clock size={16} />
                    <p class="text-md">{formatMinutes(duration)}</p>
                </div>
            </div>
            <p class="text-foreground">{description}</p>
            <div class="text-muted-foreground">
                <p>{updatedAt}</p>
            </div>
            <div class="border-b mt-2"></div>
            <div class="flex items-center justify-between mt-2">
                <ProductModal
                    {statuses}
                    {categories}
                    {availabilities}
                    {card}
                    modalForm={updateProductForm}
                >
                    <div
                        class="cursor-pointer flex items-center gap-2 p-2 hover:bg-muted"
                    >
                        <SquarePen size={20} />
                        <p>Editer</p>
                    </div>
                </ProductModal>
                <div class="flex items-center gap-2">
                    <Button.Root
                        class="cursor-pointer p-2 hover:bg-muted hover:text-tertiary"
                    >
                        <Archive size={20} />
                    </Button.Root>
                    <AlertDialog>
                        <form
                            method="POST"
                            action="?/delete"
                            class="grid gap-4"
                            use:deleteEnhance
                        >
                            <input
                                id="id"
                                name="id"
                                type="hidden"
                                value={$deleteForm.id}
                            />
                            <Button.Root
                                type="button"
                                class="cursor-pointer p-2 hover:bg-muted hover:text-destructive"
                            >
                                <Trash2 size={20} />
                            </Button.Root>
                        </form>
                    </AlertDialog>
                </div>
            </div>
        </div>
    </div>
</div>
