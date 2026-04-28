export const ROLES = {
    visitor: "visitor",
    standard: "standard",
    premium: "premium",
    guest: "guest",
    partner: "partner",
    admin: "administrator",
} as const;
export type Role = (typeof ROLES)[keyof typeof ROLES]
