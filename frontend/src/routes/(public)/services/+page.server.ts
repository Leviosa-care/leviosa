import type { PageServerLoad } from "./$types";
import { env } from "$env/dynamic/private";
import { error } from "@sveltejs/kit";

export const load: PageServerLoad = async ({ fetch }) => {
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
                    // Use the first active one_time price as the display price
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
