export const ROLES = {
    visitor: "visitor",
    standard: "standard",
    premium: "premium",
    guest: "guest",
    partner: "partner",
    admin: "admin",
} as const;
export type Role = keyof typeof ROLES
