import type { Actions, PageServerLoad } from "./$types"
import { fail, redirect } from "@sveltejs/kit"
import { API_URL } from "$env/static/private"

import { setError, superValidate } from 'sveltekit-superforms';
import { arktype } from 'sveltekit-superforms/adapters';

import {
    approveUserSchema,
    approveUserDefaults,
    updateUserRoleSchema,
    updateUserRoleDefaults,
} from "./schemas"
import { mockUsers, mockPendingUsers } from "./mockData"

// Enable mock mode for development
const MOCK_MODE = true

export const load: PageServerLoad = async ({ parent, fetch, cookies }) => {
    // ⬅️ pulls data from (ops)/+layout.server.ts
    const { permissions } = await parent()

    if (!permissions.canAccessOps) {
        throw redirect(302, '/app')
    }

    // Initialize forms
    const approveUserForm = await superValidate(arktype(approveUserSchema, { defaults: approveUserDefaults }))
    const updateUserRoleForm = await superValidate(arktype(updateUserRoleSchema, { defaults: updateUserRoleDefaults }))

    let users = []
    let pendingUsers = []

    if (MOCK_MODE) {
        // Use mock data in development
        users = mockUsers
        pendingUsers = mockPendingUsers
    } else {
        const sessionCookie = cookies.get('session');

        // Fetch all users (admin endpoint)
        const usersRes = await fetch(`${API_URL}/admin/users`, {
            headers: {
                'Authorization': `Bearer ${sessionCookie}`,
            }
        })

        if (usersRes.ok) {
            users = await usersRes.json()
        }

        // Fetch pending users (admin endpoint)
        const pendingUsersRes = await fetch(`${API_URL}/admin/auth/admin/users/pending`, {
            headers: {
                'Authorization': `Bearer ${sessionCookie}`,
            }
        })

        if (pendingUsersRes.ok) {
            pendingUsers = await pendingUsersRes.json()
        }
    }

    return {
        approveUserForm,
        updateUserRoleForm,
        users,
        pendingUsers,
    }
}

export const actions = {
    approveUser: async ({ request, fetch, cookies }) => {
        const form = await superValidate(request, arktype(approveUserSchema, { defaults: approveUserDefaults }))

        if (!form.valid) {
            return fail(400, { approveUserForm: form })
        }

        if (MOCK_MODE) {
            // Mock mode: just log and return success
            console.log('✅ [MOCK] Approving user:', {
                user_id: form.data.user_id,
                role: form.data.role
            })

            // Simulate success
            return { approveUserForm: form }
        }

        const sessionCookie = cookies.get('session');

        // Approve user
        const res = await fetch(`${API_URL}/admin/users/approve`, {
            method: "PATCH",
            headers: {
                'Content-Type': "application/json",
                'Authorization': `Bearer ${sessionCookie}`,
            },
            body: JSON.stringify({
                user_id: form.data.user_id,
                role: form.data.role,
            })
        })

        if (!res.ok) {
            switch (res.status) {
                case 400:
                    return setError(form, "Rôle invalide.");
                case 401:
                    return setError(form, "Non autorisé. Veuillez vous reconnecter.");
                case 403:
                    return setError(form, "Accès refusé. Vous devez être administrateur.");
                case 404:
                    return setError(form, "Utilisateur introuvable.");
                case 415:
                    return setError(form, "Type de contenu non supporté.");
                case 500:
                    return setError(form, "Erreur serveur. Veuillez réessayer.");
                case 503:
                    return setError(form, "Service temporairement indisponible.");
                default:
                    return setError(form, "Une erreur inattendue est survenue.");
            }
        }

        return { approveUserForm: form }
    },

    updateUserRole: async ({ request, fetch, cookies }) => {
        const form = await superValidate(request, arktype(updateUserRoleSchema, { defaults: updateUserRoleDefaults }))

        if (!form.valid) {
            return fail(400, { updateUserRoleForm: form })
        }

        if (MOCK_MODE) {
            // Mock mode: just log and return success
            console.log('✏️ [MOCK] Updating user role:', {
                user_id: form.data.user_id,
                role: form.data.role
            })

            // Simulate success
            return { updateUserRoleForm: form }
        }

        const sessionCookie = cookies.get('session');

        // Update user role
        const res = await fetch(`${API_URL}/admin/users/${form.data.user_id}/role`, {
            method: "PATCH",
            headers: {
                'Content-Type': "application/json",
                'Authorization': `Bearer ${sessionCookie}`,
            },
            body: JSON.stringify({
                user_id: form.data.user_id,
                role: form.data.role,
            })
        })

        if (!res.ok) {
            switch (res.status) {
                case 400:
                    return setError(form, "Rôle invalide.");
                case 401:
                    return setError(form, "Non autorisé. Veuillez vous reconnecter.");
                case 403:
                    return setError(form, "Accès refusé. Vous devez être administrateur.");
                case 404:
                    return setError(form, "Utilisateur introuvable.");
                case 415:
                    return setError(form, "Type de contenu non supporté.");
                case 500:
                    return setError(form, "Erreur serveur. Veuillez réessayer.");
                case 503:
                    return setError(form, "Service temporairement indisponible.");
                default:
                    return setError(form, "Une erreur inattendue est survenue.");
            }
        }

        return { updateUserRoleForm: form }
    },
} satisfies Actions
