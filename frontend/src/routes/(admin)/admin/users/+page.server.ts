import { env } from "$env/dynamic/private";
import type { PageServerLoad, Actions } from "./$types";
import { redirect } from "@sveltejs/kit";
import { mockUsers } from "$lib/data/mockUsers";

export const load: PageServerLoad = async ({ fetch }) => {
	if (env.USE_MOCK_DATA === "true") {
		return {
			users: mockUsers,
			pendingUsers: mockUsers.filter((u) => u.status === "pending")
		};
	}

	try {
		const usersRes = await fetch(`${env.API_URL}/admin/users`);
		if (!usersRes.ok) {
			throw new Error(`Failed to fetch users: ${usersRes.statusText}`);
		}
		const users = await usersRes.json();

		const pendingRes = await fetch(`${env.API_URL}/admin/auth/users/pending`);
		if (!pendingRes.ok) {
			throw new Error(`Failed to fetch pending users: ${pendingRes.statusText}`);
		}
		const pendingUsers = await pendingRes.json();

		return { users, pendingUsers };
	} catch (error) {
		console.error("Error loading users:", error);
		return { users: [], pendingUsers: [] };
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
			// Mock delete - just log it
			console.log("Mock delete user:", id);
			return { success: "User deleted successfully" };
		}

		try {
			const res = await fetch(`${env.API_URL}/admin/auth/users/${id}`, {
				method: "DELETE"
			});

			if (!res.ok) {
				throw new Error(`Failed to delete user: ${res.statusText}`);
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
			// Mock update - just log it
			console.log("Mock update role:", id, "to", role);
			return { success: "User role updated successfully" };
		}

		try {
			const res = await fetch(`${env.API_URL}/admin/users/${id}/role`, {
				method: "PUT",
				headers: { "Content-Type": "application/json" },
				body: JSON.stringify({ role })
			});

			if (!res.ok) {
				throw new Error(`Failed to update role: ${res.statusText}`);
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
			// Mock approve - just log it
			console.log("Mock approve user:", id);
			return { success: "User approved successfully" };
		}

		try {
			const res = await fetch(`${env.API_URL}/admin/users/approve`, {
				method: "POST",
				headers: { "Content-Type": "application/json" },
				body: JSON.stringify({ user_id: id })
			});

			if (!res.ok) {
				throw new Error(`Failed to approve user: ${res.statusText}`);
			}

			return { success: "User approved successfully" };
		} catch (error) {
			console.error("Error approving user:", error);
			return { error: "Failed to approve user" };
		}
	}
};
