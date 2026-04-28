import { type Role } from '$lib/types/role'

export type Permissions = {
    // consumer
    // canViewServices: boolean
    // canBook: boolean
    // canPay: boolean
    // canAccessPremium: boolean
    // ops
    canAccessOps: boolean
    canAccessApp: boolean
    // canCreateAvailability: boolean
    // canManageBookings: boolean
    // canManageCatalog: boolean
    // canManageUsers: boolean
}

export function computePermissions(role: Role): Permissions {
    switch (role) {
        case "visitor":
            return {

                canAccessOps: false,
                canAccessApp: false,
            }
        case "standard":
            return {

                canAccessOps: false,
                canAccessApp: false,
            }
        case "premium":
            return {

                canAccessOps: false,
                canAccessApp: false,
            }
        case "partner":
            return {

                canAccessOps: true,
                canAccessApp: true,
            }
        case "administrator":
            return {

                canAccessOps: true,
                canAccessApp: true,
            }
        case "guest":
            return {

                canAccessOps: false,
                canAccessApp: false,
            }
    }
}

