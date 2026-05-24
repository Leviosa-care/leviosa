import { env } from "$env/dynamic/private";
import type { PageServerLoad } from "./$types";

export interface Building {
	id: string;
	name: string;
	address: string;
	city: string;
	postal_code: string;
	country: string;
	description: string;
	phone: string;
	email: string;
	is_active: boolean;
}

export interface Room {
	id: string;
	building_id: string;
	name: string;
	description: string;
	room_number: string;
	capacity: number;
	equipment: string[];
	is_active: boolean;
}

export interface Allocation {
	id: string;
	room_id: string;
	user_id: string;
	allocation_type: "dedicated" | "shared";
	start_date: string | null;
	end_date: string | null;
	is_active: boolean;
}

export interface Partner {
	id: string;
	user_id: string;
	bio: string;
	experience: string;
	category_ids: string[];
	product_ids: string[];
	stripe_account_status: string;
	stripe_onboarding_complete: boolean;
	created_at: string;
	updated_at: string;
}

export const load: PageServerLoad = async ({ fetch }) => {
	async function fetchBuildings(): Promise<Building[]> {
		try {
			const res = await fetch(`${env.API_URL}/buildings`);
			if (!res.ok) return [];
			return await res.json();
		} catch {
			return [];
		}
	}

	async function fetchPartners(): Promise<Partner[]> {
		try {
			const res = await fetch(`${env.API_URL}/admin/partners`);
			if (!res.ok) return [];
			return await res.json();
		} catch {
			return [];
		}
	}

	const [buildings, partners] = await Promise.all([fetchBuildings(), fetchPartners()]);

	return { buildings, partners };
};
