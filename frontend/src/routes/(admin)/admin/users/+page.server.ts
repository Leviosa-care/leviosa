import { env } from "$env/dynamic/private";
import type { PageServerLoad, Actions } from "./$types";
import { redirect } from "@sveltejs/kit";
import { mockUsers, mockPartners } from "$lib/data/mockData";

interface BackendUserResponse {
	id: string;
	state: string;
	email: string;
	picture?: string;
	created_at: string;
	logged_in_at: string | null;
	role?: string;
	birthdate?: string;
	last_name?: string;
	first_name?: string;
	gender?: string;
	telephone?: string;
	postal_code?: string;
	city?: string;
	address1?: string;
	address2?: string;
	google_id?: string;
	apple_id?: string;
}

interface FrontendUser {
	id: string;
	status: string;
	email: string;
	picture?: string;
	createdAt: string;
	role: string;
	firstname?: string;
	lastname?: string;
	telephone?: string;
}

interface BackendPartnerResponse {
	id: string;
	user_id: string;
	bio: string;
	experience: string;
	category_ids?: string[];
	product_ids?: string[];
	stripe_account_status: string;
	stripe_onboarding_complete: boolean;
	created_at: string;
	updated_at: string;
}

interface FrontendPartner {
	id: string;
	userId: string;
	bio: string;
	experience: string;
	categoryCount: number;
	productCount: number;
	stripeAccountStatus: string;
	stripeOnboardingComplete: boolean;
	createdAt: string;
	updatedAt: string;
}

function mapBackendUserToFrontend(user: BackendUserResponse): FrontendUser {
	return {
		id: user.id,
		status: user.state === "active" ? "approved" : user.state,
		email: user.email,
		picture: user.picture,
		createdAt: user.created_at,
		role: user.role || "standard",
		firstname: user.first_name,
		lastname: user.last_name,
		telephone: user.telephone
	};
}

function mapBackendPartnerToFrontend(partner: BackendPartnerResponse): FrontendPartner {
	return {
		id: partner.id,
		userId: partner.user_id,
		bio: partner.bio,
		experience: partner.experience,
		categoryCount: partner.category_ids?.length ?? 0,
		productCount: partner.product_ids?.length ?? 0,
		stripeAccountStatus: partner.stripe_account_status,
		stripeOnboardingComplete: partner.stripe_onboarding_complete,
		createdAt: partner.created_at,
		updatedAt: partner.updated_at
	};
}

export const load: PageServerLoad = async ({ fetch }) => {
	if (env.USE_MOCK_DATA === "true") {
		return {
			users: mockUsers.map(mapBackendUserToFrontend),
			pendingUsers: mockUsers.filter((u) => u.state === "pending").map(mapBackendUserToFrontend),
			partners: mockPartners.map(mapBackendPartnerToFrontend),
			partnersError: false
		};
	}

	try {
		const [usersRes, pendingRes, partnersRes] = await Promise.all([
			fetch(`${env.API_URL}/admin/users`),
			fetch(`${env.API_URL}/admin/auth/users/pending`),
			fetch(`${env.API_URL}/admin/partners`)
		]);

		if (!usersRes.ok) {
			throw new Error(`Failed to fetch users: ${usersRes.statusText}`);
		}
		if (!pendingRes.ok) {
			throw new Error(`Failed to fetch pending users: ${pendingRes.statusText}`);
		}

		const [backendUsers, backendPendingUsers] = await Promise.all([
			usersRes.json() as Promise<BackendUserResponse[]>,
			pendingRes.json() as Promise<BackendUserResponse[]>
		]);

		const users = backendUsers.map(mapBackendUserToFrontend);
		const pendingUsers = backendPendingUsers.map(mapBackendUserToFrontend);

		let partners: FrontendPartner[] = [];
		let partnersError = false;
		if (partnersRes.ok) {
			const backendPartners: BackendPartnerResponse[] = await partnersRes.json();
			partners = backendPartners.map(mapBackendPartnerToFrontend);
		} else {
			console.error(`Failed to fetch partners: ${partnersRes.statusText}`);
			partnersError = true;
		}

		return { users, pendingUsers, partners, partnersError };
	} catch (error) {
		console.error("Error loading admin users page:", error);
		return { users: [], pendingUsers: [], partners: [], partnersError: false };
	}
};

export const actions: Actions = {
	deleteUser: async ({ request, fetch }) => {
		const formData = await request.formData();
		const id = formData.get("id") as string;

		if (!id) {
			return { error: "User ID is required" };
		}

		if (env.USE_MOCK_DATA === "true") {
			console.log("Mock delete user:", id);
			return { success: "User deleted successfully" };
		}

		try {
			const res = await fetch(`${env.API_URL}/admin/auth/users/${id}`, {
				method: "DELETE"
			});

			if (!res.ok) {
				const errorText = await res.text();
				throw new Error(`Failed to delete user: ${res.statusText} - ${errorText}`);
			}

			return { success: "User deleted successfully" };
		} catch (error) {
			console.error("Error deleting user:", error);
			return { error: "Failed to delete user" };
		}
	},

	updateRole: async ({ request, fetch }) => {
		const formData = await request.formData();
		const id = formData.get("id") as string;
		const role = formData.get("role") as string;

		if (!id || !role) {
			return { error: "User ID and role are required" };
		}

		if (env.USE_MOCK_DATA === "true") {
			console.log("Mock update role:", id, "to", role);
			return { success: "User role updated successfully" };
		}

		try {
			const res = await fetch(`${env.API_URL}/admin/users/${id}/role`, {
				method: "PATCH",
				headers: { "Content-Type": "application/json" },
				body: JSON.stringify({ role })
			});

			if (!res.ok) {
				const errorText = await res.text();
				throw new Error(`Failed to update role: ${res.statusText} - ${errorText}`);
			}

			return { success: "User role updated successfully" };
		} catch (error) {
			console.error("Error updating role:", error);
			return { error: "Failed to update user role" };
		}
	},

	approveUser: async ({ request, fetch }) => {
		const formData = await request.formData();
		const id = formData.get("id") as string;

		if (!id) {
			return { error: "User ID is required" };
		}

		if (env.USE_MOCK_DATA === "true") {
			console.log("Mock approve user:", id);
			return { success: "User approved successfully" };
		}

		try {
			// Default to "standard" role for approved users if not specified
			const res = await fetch(`${env.API_URL}/admin/users/approve`, {
				method: "PATCH",
				headers: { "Content-Type": "application/json" },
				body: JSON.stringify({ user_id: id, role: "standard" })
			});

			if (!res.ok) {
				const errorText = await res.text();
				throw new Error(`Failed to approve user: ${res.statusText} - ${errorText}`);
			}

			return { success: "User approved successfully" };
		} catch (error) {
			console.error("Error approving user:", error);
			return { error: "Failed to approve user" };
		}
	},

	verifyPartner: async ({ request, fetch }) => {
		const formData = await request.formData();
		const id = formData.get("id") as string;

		if (!id) {
			return { error: "Partner ID is required" };
		}

		if (env.USE_MOCK_DATA === "true") {
			console.log("Mock verify partner:", id);
			return { success: "Partenaire vérifié avec succès" };
		}

		try {
			const res = await fetch(`${env.API_URL}/admin/partners/${id}/verify`, {
				method: "POST"
			});

			if (!res.ok) {
				const errorText = await res.text();
				throw new Error(`Failed to verify partner: ${res.statusText} - ${errorText}`);
			}

			return { success: "Partenaire vérifié avec succès" };
		} catch (error) {
			console.error("Error verifying partner:", error);
			return { error: "Impossible de vérifier le partenaire" };
		}
	},

	deletePartner: async ({ request, fetch }) => {
		const formData = await request.formData();
		const id = formData.get("id") as string;

		if (!id) {
			return { error: "Partner ID is required" };
		}

		if (env.USE_MOCK_DATA === "true") {
			console.log("Mock delete partner:", id);
			return { success: "Partenaire supprimé avec succès" };
		}

		try {
			const res = await fetch(`${env.API_URL}/admin/partners/${id}`, {
				method: "DELETE"
			});

			if (!res.ok) {
				const errorText = await res.text();
				throw new Error(`Failed to delete partner: ${res.statusText} - ${errorText}`);
			}

			return { success: "Partenaire supprimé avec succès" };
		} catch (error) {
			console.error("Error deleting partner:", error);
			return { error: "Impossible de supprimer le partenaire" };
		}
	}
};
