import type { PageServerLoad } from "./$types";
import { env } from "$env/dynamic/private";
import { error } from "@sveltejs/kit";
import { mockCategories, mockProducts, mockPrices } from "$lib/data/mockData";

export const load: PageServerLoad = async ({ fetch }) => {
    if (env.USE_MOCK_DATA === "true") {
        const pricesByProduct: Record<string, number> = {};
        for (const price of mockPrices) {
            if (price.isActive && price.interval === "one_time" && !(price.productId in pricesByProduct)) {
                pricesByProduct[price.productId] = price.amount;
            }
        }
        return {
            categories: mockCategories,
            products: mockProducts.map((p) => ({ ...p, category: { id: p.category, name: "" } })),
            pricesByProduct,
        };
    }

    // Fetch published categories
    const categoriesRes = await fetch(`${env.API_URL}/categories`);
    if (!categoriesRes.ok) {
        throw error(500, "Impossible de charger les catégories");
    }
    const categories = await categoriesRes.json();

    // Fetch published products
    const productsRes = await fetch(`${env.API_URL}/products`);
    if (!productsRes.ok) {
        throw error(500, "Impossible de charger les produits");
    }
    const products = await productsRes.json();

    // Fetch prices for each product
    const pricesByProduct: Record<string, number> = {};
    await Promise.all(
        products.map(async (product: { id: string }) => {
            try {
                const pricesRes = await fetch(
                    `${env.API_URL}/products/${product.id}/prices`
                );
                if (pricesRes.ok) {
                    const prices = await pricesRes.json();
                    const price = prices.find(
                        (p: { interval: string; isActive: boolean }) =>
                            p.isActive && p.interval === "one_time"
                    ) ?? prices[0];
                    if (price) {
                        pricesByProduct[product.id] = price.amount;
                    }
                }
            } catch {
                // Price fetch failure is non-blocking
            }
        })
    );

    return {
        categories,
        products,
        pricesByProduct,
    };
};
