<script lang="ts">
    import type { PageProps } from "./$types";
    import { superForm } from "sveltekit-superforms";
    import { Button } from "bits-ui";
    import { ArrowLeft } from "@lucide/svelte";

    import Input from "$lib/ui/Input.svelte";
    import ProgressIndicator from "../ProgressIndicator.svelte";

    let { data }: PageProps = $props();
    const { form, enhance, errors } = superForm(data.form);

    // Category / product data from load
    let categories = $state(data.categories ?? []);
    let products = $state(data.products ?? []);

    // Track selected IDs (synced with hidden form fields)
    let selectedCategoryIds = $state<string[]>(data.form.data.category_ids ?? []);
    let selectedProductIds = $state<string[]>(data.form.data.product_ids ?? []);

    // Password match check for client-side feedback
    let isSamePassword = $derived(
        $form.password === $form.confirm && $form.confirm !== "",
    );

    // Toggle a value in a multi-select array
    function toggleItem(list: string[], id: string): string[] {
        return list.includes(id) ? list.filter((x) => x !== id) : [...list, id];
    }

    // Filter products by selected categories (show all when no category selected)
    let filteredProducts = $derived(
        selectedCategoryIds.length === 0
            ? products
            : products.filter((p: any) => selectedCategoryIds.includes(p.category ?? p.category_id ?? "")),
    );

    import type { Step } from "$lib/types/step";
    let steps: Step[] = [
        { id: 1, name: "Email", status: "complete" },
        { id: 2, name: "Identité", status: "complete" },
        { id: 3, name: "Adresse", status: "complete" },
        { id: 4, name: "Profil partenaire", status: "current" },
    ];
</script>

<div
    class="min-h-[100vh] flex flex-col justify-center items-center px-8 max-w-[1080px] mx-auto py-12"
>
    <div class="w-full max-w-[640px]">
        <Button.Root
            class="mb-4 justify-center gap-2 items-center h-input rounded-input hover:bg-gray-50 focus-visible:ring-dark focus-visible:ring-offset-background focus-visible:outline-hidden inline-flex px-4 text-[15px] font-bold focus-visible:ring-2 focus-visible:ring-offset-2 active:scale-[0.98] cursor-pointer"
            type="button"
        >
            <ArrowLeft size={16} />
            <span class="text-xs sm:text-base">Retour à l'étape précédente</span>
        </Button.Root>

        <div
            class="p-4 md:p-12 grid gap-12 border-gray-200 border-1 rounded-input"
        >
            <div class="w-full flex-none">
                <ProgressIndicator {steps} />
            </div>
            <div class="mt-8">
                <h3 class="font-bold sm:text-2xl text-lg">
                    Profil partenaire
                </h3>
                <p class="text-gray-600 sm:text-base text-sm">
                    Complétez votre profil pour finaliser votre inscription en tant que partenaire
                </p>
            </div>

            <form method="POST" use:enhance class="grid gap-8">
                <!-- Hidden fields from previous steps -->
                <input type="hidden" name="firstname" bind:value={$form.firstname} />
                <input type="hidden" name="lastname" bind:value={$form.lastname} />
                <input type="hidden" name="gender" bind:value={$form.gender} />
                <input type="hidden" name="birthdate" bind:value={$form.birthdate} />
                <input type="hidden" name="telephone" bind:value={$form.telephone} />
                <input type="hidden" name="address1" bind:value={$form.address1} />
                <input type="hidden" name="address2" bind:value={$form.address2} />
                <input type="hidden" name="city" bind:value={$form.city} />
                <input type="hidden" name="postalCode" bind:value={$form.postalCode} />

                <!-- Global form error -->
                {#if $errors._errors && $errors._errors.length > 0}
                    <div class="rounded-input border border-red-300 bg-red-50 p-4 text-sm text-red-700">
                        {$errors._errors[0]}
                    </div>
                {/if}

                <!-- Password section -->
                <div class="grid gap-6">
                    <h4 class="font-semibold text-lg border-b border-gray-200 pb-2">
                        Mot de passe
                    </h4>
                    <div class="grid gap-4">
                        <Input
                            bind:value={$form.password}
                            name="password"
                            label="Mot de passe"
                            type="password"
                            error={$errors.password ? $errors.password[0] : ""}
                        />
                        <Input
                            bind:value={$form.confirm}
                            name="confirm"
                            label="Confirmation du mot de passe"
                            type="password"
                            error={$errors.confirm ? $errors.confirm[0] : ""}
                        />
                        {#if $form.confirm !== "" && !isSamePassword}
                            <p class="text-red-500 text-sm">
                                Les mots de passe ne correspondent pas.
                            </p>
                        {/if}
                    </div>
                </div>

                <!-- Bio & Experience section -->
                <div class="grid gap-6">
                    <h4 class="font-semibold text-lg border-b border-gray-200 pb-2">
                        À propos de vous
                    </h4>
                    <div class="grid gap-4">
                        <div class="grid gap-1.5">
                            <label class="font-semibold" for="bio">Biographie</label>
                            <textarea
                                id="bio"
                                name="bio"
                                rows="4"
                                bind:value={$form.bio}
                                class="border-border-input border px-3 py-2 rounded-lg bg-background hover:border-dark-40 focus:ring-foreground focus:ring-offset-background focus:outline-hidden w-full text-sm focus:ring-2 focus:ring-offset-2 transition-all resize-none"
                                placeholder="Décrivez-vous en quelques mots (max. 1000 caractères)"
                                maxlength="1000"
                            ></textarea>
                            {#if $errors.bio}
                                <p class="text-red-500 text-sm">{$errors.bio[0]}</p>
                            {/if}
                            <p class="text-xs text-gray-400 text-right">{$form.bio.length}/1000</p>
                        </div>
                        <div class="grid gap-1.5">
                            <label class="font-semibold" for="experience">Expérience</label>
                            <textarea
                                id="experience"
                                name="experience"
                                rows="4"
                                bind:value={$form.experience}
                                class="border-border-input border px-3 py-2 rounded-lg bg-background hover:border-dark-40 focus:ring-foreground focus:ring-offset-background focus:outline-hidden w-full text-sm focus:ring-2 focus:ring-offset-2 transition-all resize-none"
                                placeholder="Décrivez votre expérience professionnelle (max. 2000 caractères)"
                                maxlength="2000"
                            ></textarea>
                            {#if $errors.experience}
                                <p class="text-red-500 text-sm">{$errors.experience[0]}</p>
                            {/if}
                            <p class="text-xs text-gray-400 text-right">{$form.experience.length}/2000</p>
                        </div>
                    </div>
                </div>

                <!-- Categories multi-select -->
                <div class="grid gap-6">
                    <h4 class="font-semibold text-lg border-b border-gray-200 pb-2">
                        Catégories
                    </h4>
                    <p class="text-sm text-gray-600">Sélectionnez les catégories dans lesquelles vous exercez</p>
                    {#if categories.length === 0}
                        <p class="text-sm text-gray-400 italic">Aucune catégorie disponible pour le moment.</p>
                    {:else}
                        <div class="flex flex-wrap gap-2">
                            {#each categories as category (category.id)}
                                {@const isSelected = selectedCategoryIds.includes(category.id)}
                                <button
                                    type="button"
                                    class="px-4 py-2 rounded-full border text-sm font-medium transition-all cursor-pointer
                                        {isSelected
                                            ? 'bg-dark text-white border-dark'
                                            : 'bg-white text-gray-700 border-gray-300 hover:border-dark hover:bg-gray-50'}"
                                    onclick={() => {
                                        selectedCategoryIds = toggleItem(selectedCategoryIds, category.id);
                                    }}
                                >
                                    {category.name}
                                </button>
                            {/each}
                        </div>
                    {/if}
                    <!-- Hidden inputs for category_ids -->
                    {#each selectedCategoryIds as catId}
                        <input type="hidden" name="category_ids" value={catId} />
                    {/each}
                </div>

                <!-- Products multi-select -->
                <div class="grid gap-6">
                    <h4 class="font-semibold text-lg border-b border-gray-200 pb-2">
                        Produits & Services
                    </h4>
                    <p class="text-sm text-gray-600">Sélectionnez les produits et services que vous proposez</p>
                    {#if filteredProducts.length === 0}
                        <p class="text-sm text-gray-400 italic">
                            {selectedCategoryIds.length > 0
                                ? "Aucun produit dans les catégories sélectionnées."
                                : "Aucun produit disponible pour le moment."}
                        </p>
                    {:else}
                        <div class="grid gap-2">
                            {#each filteredProducts as product (product.id)}
                                {@const isSelected = selectedProductIds.includes(product.id)}
                                <button
                                    type="button"
                                    class="flex items-center gap-3 px-4 py-3 rounded-lg border text-sm text-left transition-all cursor-pointer
                                        {isSelected
                                            ? 'bg-dark/5 border-dark ring-1 ring-dark'
                                            : 'bg-white border-gray-300 hover:border-gray-400 hover:bg-gray-50'}"
                                    onclick={() => {
                                        selectedProductIds = toggleItem(selectedProductIds, product.id);
                                    }}
                                >
                                    <div
                                        class="flex-shrink-0 w-5 h-5 rounded border-2 flex items-center justify-center transition-all
                                            {isSelected ? 'bg-dark border-dark' : 'border-gray-300'}"
                                    >
                                        {#if isSelected}
                                            <svg class="w-3 h-3 text-white" viewBox="0 0 12 12" fill="none">
                                                <path d="M2 6l3 3 5-5" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                                            </svg>
                                        {/if}
                                    </div>
                                    <div>
                                        <span class="font-medium">{product.name}</span>
                                        {#if product.description}
                                            <span class="text-gray-500 ml-2">{product.description}</span>
                                        {/if}
                                    </div>
                                </button>
                            {/each}
                        </div>
                    {/if}
                    <!-- Hidden inputs for product_ids -->
                    {#each selectedProductIds as prodId}
                        <input type="hidden" name="product_ids" value={prodId} />
                    {/each}
                </div>

                <!-- Submit -->
                <Button.Root
                    disabled={!isSamePassword}
                    class="mt-4 justify-center items-center h-input rounded-input bg-dark text-background shadow-mini hover:bg-dark/95 focus-visible:ring-dark focus-visible:ring-offset-background focus-visible:outline-hidden inline-flex px-4 text-[15px] font-bold focus-visible:ring-2 focus-visible:ring-offset-2 active:scale-[0.98] cursor-pointer disabled:opacity-50 disabled:cursor-not-allowed"
                    type="submit"
                >
                    Finaliser mon inscription
                </Button.Root>
                <p class="text-sm text-center mt-4 text-muted-foreground">
                    Vos informations sont stockées et protégées en toute sécurité.
                </p>
            </form>
        </div>
    </div>
</div>
