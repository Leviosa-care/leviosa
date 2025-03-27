import type { LayoutServerLoad } from "./$types"

import type { Role } from "$lib/types";

type PageRes = {
    role: Role
};

export const load: LayoutServerLoad = ({ locals }): PageRes => {
    if (!locals.user?.role) {
        console.error('no role set from the local user, this is embarassing')
    }
    const role = locals.user?.role as Role
    console.log("the role is:", role);
    return {
        role,
    };
}

