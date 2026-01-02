<script lang="ts">
    import type { PageData } from "./$types";
    import Tabs from "$lib/ui/bits-components/Tabs.svelte";
    import TabsList from "$lib/ui/bits-components/TabsList.svelte";
    import TabsTrigger from "$lib/ui/bits-components/TabsTrigger.svelte";
    import TabsContent from "$lib/ui/bits-components/TabsContent.svelte";
    import Categories from "./Categories.svelte";
    import Products from "./Products.svelte";
    import Exercices from "./Exercices.svelte";

    let { data }: { data: PageData } = $props();
    interface Trigger {
        value: string;
        name: string;
    }
    let triggers: Trigger[] = [
        { value: "categories", name: "Catégories" },
        { value: "produits", name: "Produits" },
        { value: "prix", name: "Prix" },
        { value: "coupons", name: "Codes promo" },
        { value: "promotion-codes", name: "Codes promotion" },
        { value: "exercices", name: "Exercices" },
    ];
</script>

{#snippet tab_trigger(value: string, name: string)}
    <TabsTrigger
        {value}
        class="px-2 py-1.75 md:px-4 md:py-2 md:h-8 rounded-none md:rounded-[7px] bg-transparent border-b-3 data-[state=active]:shadow-none mb-[-2px] md:mb-0 md:border-b-0 data-[state=active]:border-b-dark-900 data-[state=active]:text-foreground md:data-[state=active]:bg-white md:data-[state=active]:shadow-mini dark:md:data-[state=active]:bg-muted data-[state=inactive]:border-transparent data-[state=inactive]:text-foreground-alt md:data-[state=inactive]:hover:bg-dark-04 hover:text-foreground transition-colors"
    >
        {name}
    </TabsTrigger>
{/snippet}

<div class="h-[100vh] flex-1 flex flex-col overflow-hidden bg-gray-50">
    <Tabs value="categories" class="flex flex-col h-full">
        <!-- Header with title and description -->
        <div class="bg-white border-b border-border-card px-6 py-6">
            <div class="grid gap-1 mb-6">
                <h1 class="text-2xl font-semibold tracking-tight">Catalogue</h1>
                <p class="text-sm text-foreground-alt">
                    Gérez vos catégories, produits, prix et exercices
                </p>
            </div>

            <!-- Scrollable tabs -->
            <div class="overflow-x-auto -mx-6 px-6 scrollbar-hide">
                <TabsList
                    class="inline-flex gap-2 md:gap-1 bg-transparent text-sm font-semibold min-w-max border-b-1 border-border-card md:border-b-0 md:rounded-9px md:bg-dark-10 md:shadow-mini-inset dark:md:bg-background md:p-1 md:leading-[0.01em] dark:md:border dark:md:border-neutral-600/30"
                >
                    {#each triggers as trigger}
                        {@render tab_trigger(trigger.value, trigger.name)}
                    {/each}
                </TabsList>
            </div>
        </div>

        <!-- Tab content area -->
        <div class="flex-1 overflow-y-auto">
            <TabsContent value="categories" class="h-full p-6">
                <Categories {data} />
            </TabsContent>
            <TabsContent value="produits" class="h-full p-6">
                <Products {data} />
            </TabsContent>
            <TabsContent value="prix" class="h-full p-6">
                <!-- Prix content -->
            </TabsContent>
            <TabsContent value="coupons" class="h-full p-6">
                <!-- Coupons content -->
            </TabsContent>
            <TabsContent value="promotion-codes" class="h-full p-6">
                <!-- Promotion codes content -->
            </TabsContent>
            <TabsContent value="exercices" class="h-full p-6">
                <Exercices {data} />
            </TabsContent>
        </div>
    </Tabs>
</div>

<style>
    /* Hide scrollbar but keep functionality */
    .scrollbar-hide::-webkit-scrollbar {
        display: none;
    }
    .scrollbar-hide {
        -ms-overflow-style: none;
        scrollbar-width: none;
    }
</style>
