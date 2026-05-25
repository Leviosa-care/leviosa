import { type Role } from '$lib/types/role'

export type Permissions = {
    // client area
    canAccessClientArea: boolean
    // ops
    canAccessOps: boolean
    canAccessApp: boolean
}

export function computePermissions(role: Role): Permissions {
    switch (role) {
        case "visitor":
            return {
                canAccessClientArea: false,
                canAccessOps: false,
                canAccessApp: false,
            }
        case "standard":
            return {
                canAccessClientArea: true,
                canAccessOps: false,
                canAccessApp: false,
            }
        case "premium":
            return {
                canAccessClientArea: true,
                canAccessOps: false,
                canAccessApp: false,
            }
        case "partner":
            return {
                canAccessClientArea: false,
                canAccessOps: true,
                canAccessApp: true,
            }
        case "administrator":
            return {
                canAccessClientArea: false,
                canAccessOps: true,
                canAccessApp: true,
            }
        case "guest":
            return {
                canAccessClientArea: false,
                canAccessOps: false,
                canAccessApp: false,
            }
    }
}

