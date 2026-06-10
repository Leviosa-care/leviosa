<script lang="ts">
    import type { PageData } from "./$types";
    import Tabs from "$lib/ui/bits-components/Tabs.svelte";
    import TabsList from "$lib/ui/bits-components/TabsList.svelte";
    import TabsTrigger from "$lib/ui/bits-components/TabsTrigger.svelte";
    import TabsContent from "$lib/ui/bits-components/TabsContent.svelte";
    import Categories from "./Categories.svelte";
    import Products from "./Products.svelte";
    import Prices from "./Prices.svelte";
    import Coupons from "./Coupons.svelte";
    import PromotionCodes from "./PromotionCodes.svelte";
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
        { value: "coupons", name: "Coupons" },
        { value: "promotion-codes", name: "Codes de promotion" },
        { value: "exercices", name: "Exercices" },
    ];
</script>

<div class="h-[100vh] flex-1 flex flex-col overflow-hidden bg-background">
    <Tabs value="categories" class="flex flex-col h-full">
        <!-- Header -->
        <div class="border-b border-border px-6 py-6">
            <div class="mb-6">
                <p class="text-[11px] font-semibold uppercase tracking-[0.2em] text-muted-foreground mb-3">Administration</p>
                <h1 class="font-display text-3xl font-semibold tracking-tight text-foreground">
                    Catalogue
                </h1>
                <p class="text-sm text-muted-foreground mt-2">
                    Configurez et organisez votre offre commerciale
                </p>
                <div class="mt-4 h-px w-16 bg-foreground/20"></div>
            </div>

            <!-- Scrollable tabs -->
            <div class="overflow-x-auto -mx-6 px-6 scrollbar-hide">
                <!-- NOTE: here is the old navigation with the classic tab that we removed since there are some many triggers -->
                <!-- <TabsList -->
                <!--     class="inline-flex gap-2 md:gap-1 bg-transparent text-sm font-semibold min-w-max border-b-1 border-border-card md:border-b-0 md:rounded-9px md:bg-dark-10 md:shadow-mini-inset dark:md:bg-background md:p-1 md:leading-[0.01em] dark:md:border dark:md:border-neutral-600/30" -->
                <!-- > -->
                <TabsList
                    class="inline-flex gap-2 md:gap-1 bg-transparent text-sm font-semibold min-w-max border-b-1 border-border-card md:p-1 md:leading-[0.01em]"
                >
                    {#each triggers as trigger}
                        <TabsTrigger
                            value={trigger.value}
                            class="px-2 py-1.75 md:px-4 md:py-2 md:h-8 rounded-none bg-transparent border-b-3 data-[state=active]:shadow-none mb-[-2px] data-[state=active]:border-b-foreground data-[state=active]:text-foreground data-[state=inactive]:border-transparent data-[state=inactive]:text-muted-foreground data-[state=inactive]:hover:bg-transparent  data-[state=inactive]:hover:text-foreground-alt transition-colors cursor-pointer"
                        >
                            {trigger.name}
                        </TabsTrigger>
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
                <Prices {data} />
            </TabsContent>
            <TabsContent value="coupons" class="h-full p-6">
                <Coupons {data} />
            </TabsContent>
            <TabsContent value="promotion-codes" class="h-full p-6">
                <PromotionCodes {data} />
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
