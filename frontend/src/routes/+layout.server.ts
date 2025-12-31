import { computePermissions } from '$lib/security/permissions'
import { ROLES } from '$lib/types/role'

export const load = async ({ locals }) => {
    const user = locals.user

    const role = user?.role ?? ROLES.visitor;
    const permissions = computePermissions(role);

    return {
        user,
        permissions
    }
}

