import type { PageServerLoad } from "./$types";
import { env } from "$env/dynamic/private";
import { error } from "@sveltejs/kit";
import { mockPartners, mockUsers, mockProducts, mockCategories } from "$lib/data/mockData";

export const load: PageServerLoad = async ({ params, fetch }) => {
    const { id } = params;

    if (env.USE_MOCK_DATA === "true") {
        const partner = mockPartners.find((p) => p.id === id);
        if (!partner) throw error(404, "Praticien introuvable");

        const user = mockUsers.find((u) => u.id === partner.user_id);
        const products = mockProducts.filter((p) => partner.product_ids.includes(p.id));
        const categories = mockCategories.filter((c) => partner.category_ids.includes(c.id));

        return {
            partner: {
                id: partner.id,
                firstname: user?.first_name ?? "—",
                lastname: user?.last_name ?? "—",
                occupation: partner.occupation ?? "",
                quote: partner.quote ?? "",
                tags: partner.tags ?? [],
                picture: partner.picture,
                bio: partner.bio,
                experience: partner.experience,
            },
            products: products.map((p) => ({
                id: p.id,
                name: p.name,
                description: p.description,
                duration: p.duration,
            })),
            categories: categories.map((c) => ({ id: c.id, name: c.name })),
        };
    }

    const partnerRes = await fetch(`${env.API_URL}/partners/${id}`);
    if (!partnerRes.ok) {
        if (partnerRes.status === 404) throw error(404, "Praticien introuvable");
        throw error(500, "Impossible de charger le praticien");
    }
    const partnerData = await partnerRes.json();

    const categoriesRes = await fetch(`${env.API_URL}/categories`);
    const allCategories: any[] = categoriesRes.ok ? await categoriesRes.json() : [];
    const categories = allCategories
        .filter((c: any) => partnerData.category_ids?.includes(c.id))
        .map((c: any) => ({ id: c.id, name: c.name }));

    return {
        partner: {
            id: partnerData.id,
            firstname: partnerData.first_name ?? "—",
            lastname: partnerData.last_name ?? "—",
            occupation: partnerData.occupation ?? "",
            quote: partnerData.quote ?? "",
            tags: partnerData.tags ?? [],
            picture: partnerData.picture,
            bio: partnerData.bio,
            experience: partnerData.experience,
        },
        products: [],
        categories,
    };
};
