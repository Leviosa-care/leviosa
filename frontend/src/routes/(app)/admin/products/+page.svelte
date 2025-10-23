<script lang="ts">
    import type { PageProps } from "./$types";

    import { Tabs } from "bits-ui";

    import Product from "./Product.svelte";
    import Category from "./Category.svelte";

    let { data }: PageProps = $props();
    let {
        cards,
        statuses,
        categories,
        availabilities,
        deleteProductForm,
        createProductForm,
        updateProductForm,
        createCategoryForm,
    } = data;
    // TODO: add the Tab content brother
</script>

<div class="h-[100vh] flex-1 gap-4 overflow-y-scroll bg-gray-50">
    <!-- TODO: handle the width for that thing brother -->
    <Tabs.Root
        value="category"
        class="border-muted bg-background-alt shadow-card border"
    >
        <div class="p-4">
            <Tabs.List
                class="rounded-9px bg-dark-10 shadow-mini-inset dark:bg-background grid grid-cols-2 gap-1 p-1 text-sm font-semibold leading-[0.01em] dark:border dark:border-neutral-600/30 w-[540px] mx-auto"
            >
                <Tabs.Trigger
                    value="category"
                    class="data-[state=active]:shadow-mini dark:data-[state=active]:bg-muted h-8 rounded-[7px] bg-transparent py-2 data-[state=active]:bg-white"
                    >Categorie</Tabs.Trigger
                >
                <Tabs.Trigger
                    value="product"
                    class="data-[state=active]:shadow-mini dark:data-[state=active]:bg-muted h-8 rounded-[7px] bg-transparent py-2 data-[state=active]:bg-white"
                    >Produit</Tabs.Trigger
                >
            </Tabs.List>
        </div>
        <Tabs.Content value="product" class="select-none pt-3">
            <Product
                {cards}
                {statuses}
                {categories}
                {availabilities}
                {deleteProductForm}
                {createProductForm}
                {updateProductForm}
            />
        </Tabs.Content>
        <Tabs.Content value="category" class="select-none pt-3">
            <Category {createCategoryForm} />
        </Tabs.Content>
    </Tabs.Root>
</div>
