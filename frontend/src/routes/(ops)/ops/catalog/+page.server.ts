import { redirect } from '@sveltejs/kit'

export const load = async ({ parent }) => {
    // ⬅️ pulls data from (ops)/+layout.server.ts
    const { permissions } = await parent()

    if (!permissions.canAccessOps) {
        throw redirect(302, '/app')
    }

    return {}
}
