import { redirect, error } from '@sveltejs/kit';
import { computePermissions } from '$lib/security/permissions';

export const load = async ({ locals }) => {
	const user = locals.user;

	// Check if user is authenticated
	if (!user) {
		throw redirect(302, '/auth?redirect=' + encodeURIComponent('/staff'));
	}

	// Compute permissions to check access
	const permissions = computePermissions(user.role);

	// Check if user has partner or admin role
	if (!permissions.canAccessOps) {
		throw error(403, 'Accès non autorisé');
	}

	return {
		user,
		permissions,
	};
};
