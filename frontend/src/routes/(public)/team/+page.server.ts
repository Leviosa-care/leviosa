import type { PageServerLoad } from "./$types";
import { env } from "$env/dynamic/private";
import { error } from "@sveltejs/kit";
import { mockPartners, mockUsers, mockCategories } from "$lib/data/mockData";

export const load: PageServerLoad = async ({ fetch }) => {
    if (env.USE_MOCK_DATA === "true") {
        const partners = mockPartners.map((partner) => {
            const user = mockUsers.find((u) => u.id === partner.user_id);
            return {
                id: partner.id,
                firstname: user?.first_name ?? "—",
                lastname: user?.last_name ?? "—",
                occupation: partner.occupation ?? "",
                quote: partner.quote ?? "",
                tags: partner.tags ?? [],
                picture: partner.picture,
            };
        });

        const categories = mockCategories
            .filter((c) => c.status === "published")
            .map((c) => ({
                id: c.id,
                name: c.name,
                description: c.description ?? "",
                partners: partners.filter((p) =>
                    mockPartners.find((mp) => mp.id === p.id)?.category_ids.includes(c.id)
                ),
            }))
            .filter((c) => c.partners.length > 0);

        return { categories };
    }

    const partnersRes = await fetch(`${env.API_URL}/partners`);
    if (!partnersRes.ok) throw error(500, "Impossible de charger l'équipe");
    const partnersData = await partnersRes.json();

    const categoriesRes = await fetch(`${env.API_URL}/categories`);
    const categoriesData = categoriesRes.ok ? await categoriesRes.json() : [];

    const partners = partnersData.map((partner: any) => ({
        id: partner.id,
        firstname: partner.user?.first_name ?? "—",
        lastname: partner.user?.last_name ?? "—",
        occupation: partner.occupation ?? "",
        quote: partner.quote ?? "",
        tags: partner.tags ?? [],
        picture: partner.picture,
    }));

    const categories = categoriesData
        .map((c: any) => ({
            id: c.id,
            name: c.name,
            description: c.description ?? "",
            partners: partners.filter((p: any) =>
                partnersData
                    .find((pd: any) => pd.id === p.id)
                    ?.category_ids?.includes(c.id)
            ),
        }))
        .filter((c: any) => c.partners.length > 0);

    return { categories };
};
