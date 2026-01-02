<script lang="ts">
    import type { PageData } from "./$types";
    import Tabs from "$lib/ui/bits-components/Tabs.svelte";
    import TabsList from "$lib/ui/bits-components/TabsList.svelte";
    import TabsTrigger from "$lib/ui/bits-components/TabsTrigger.svelte";
    import TabsContent from "$lib/ui/bits-components/TabsContent.svelte";
    import Categories from "./Categories.svelte";

    let { data }: { data: PageData } = $props();
</script>

<div class="h-[100vh] flex-1 gap-4 overflow-y-scroll bg-gray-50">
    <Tabs
        value="categories"
        class="border-muted bg-background-alt shadow-card border"
    >
        <div class="p-4">
            <TabsList
                class="rounded-9px bg-dark-10 shadow-mini-inset dark:bg-background grid grid-cols-4 gap-1 p-1 text-sm font-semibold leading-[0.01em] dark:border dark:border-neutral-600/30 w-[720px] mx-auto"
            >
                <TabsTrigger
                    value="categories"
                    class="data-[state=active]:shadow-mini dark:data-[state=active]:bg-muted h-8 rounded-[7px] bg-transparent py-2 data-[state=active]:bg-white data-[state=inactive]:hover:bg-dark-04 hover:text-foreground"
                >
                    Categories
                </TabsTrigger>
                <TabsTrigger
                    value="produits"
                    class="data-[state=active]:shadow-mini dark:data-[state=active]:bg-muted h-8 rounded-[7px] bg-transparent py-2 data-[state=active]:bg-white data-[state=inactive]:hover:bg-dark-04 hover:text-foreground"
                >
                    Produits
                </TabsTrigger>
                <TabsTrigger
                    value="prix"
                    class="data-[state=active]:shadow-mini dark:data-[state=active]:bg-muted h-8 rounded-[7px] bg-transparent py-2 data-[state=active]:bg-white data-[state=inactive]:hover:bg-dark-04 hover:text-foreground"
                >
                    Prix
                </TabsTrigger>
                <TabsTrigger
                    value="exercices"
                    class="data-[state=active]:shadow-mini dark:data-[state=active]:bg-muted h-8 rounded-[7px] bg-transparent py-2 data-[state=active]:bg-white data-[state=inactive]:hover:bg-dark-04 hover:text-foreground"
                >
                    Exercices
                </TabsTrigger>
            </TabsList>
        </div>
        <TabsContent value="categories" class="select-none">
            <Categories {data} />
        </TabsContent>
        <TabsContent value="produits" class="select-none pt-3">
            <!-- Produits content -->
        </TabsContent>
        <TabsContent value="prix" class="select-none pt-3">
            <!-- Prix content -->
        </TabsContent>
        <TabsContent value="exercices" class="select-none pt-3">
            <!-- Exercices content -->
        </TabsContent>
    </Tabs>
</div>
