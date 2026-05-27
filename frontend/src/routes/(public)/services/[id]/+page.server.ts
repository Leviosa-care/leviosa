import type { PageServerLoad } from "./$types";
import { env } from "$env/dynamic/private";
import { error } from "@sveltejs/kit";
import { mockProducts, mockCategories, mockPrices } from "$lib/data/mockData";

export const load: PageServerLoad = async ({ params, fetch }) => {
    const { id } = params;

    if (env.USE_MOCK_DATA === "true") {
        const product = mockProducts.find((p) => p.id === id);
        if (!product) throw error(404, "Service introuvable");

        const category = mockCategories.find((c) => c.id === product.category) ?? null;
        const price = mockPrices.find(
            (p) => p.productId === id && p.isActive && p.interval === "one_time",
        ) ?? null;

        return {
            product: {
                ...product,
                category: category
                    ? { id: category.id, name: category.name, description: category.description }
                    : null,
            },
            price: price?.amount ?? null,
        };
    }

    const productRes = await fetch(`${env.API_URL}/products/${id}`);
    if (!productRes.ok) {
        if (productRes.status === 404) throw error(404, "Service introuvable");
        throw error(500, "Impossible de charger le service");
    }
    const product = await productRes.json();

    let price: number | null = null;
    try {
        const pricesRes = await fetch(`${env.API_URL}/products/${id}/prices`);
        if (pricesRes.ok) {
            const prices = await pricesRes.json();
            const activePrice =
                prices.find((p: { isActive: boolean; interval: string }) =>
                    p.isActive && p.interval === "one_time",
                ) ?? prices[0];
            if (activePrice) price = activePrice.amount;
        }
    } catch {
        // non-blocking
    }

    return { product, price };
};
