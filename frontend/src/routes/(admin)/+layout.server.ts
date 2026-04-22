import { redirect, error } from '@sveltejs/kit';
import { computePermissions } from '$lib/security/permissions';

export const load = async ({ locals, parent }) => {
	const { user, permissions } = await parent();

	// Check if user is authenticated
	if (!user) {
		throw redirect(302, '/auth?redirect=' + encodeURIComponent('/admin'));
	}

	// Check if user has admin role
	if (user.role !== 'admin') {
		throw error(403, 'Accès non autorisé');
	}

	return {
		user,
		permissions,
	};
};
