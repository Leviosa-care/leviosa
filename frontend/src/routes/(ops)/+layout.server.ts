import { redirect } from '@sveltejs/kit'
import { computePermissions } from '$lib/security/permissions'


export const load = async ({ locals }) => {
    const user = locals.user

    if (!user) {
        // throw redirect(302, '/auth')
        // TODO: is the way to do this ?
        throw redirect(302, 'redirectFrom?/ops/auth')
    }

    const permissions = computePermissions(user.role)

    return {
        user,
        permissions
    }
}

