import { type } from 'arktype'
import type { Infer } from 'sveltekit-superforms'

// Approve user schema
export const approveUserSchema = type({
    user_id: "string",
    role: "'visitor' | 'standard' | 'partner' | 'administrator'"
})

// Update user role schema
export const updateUserRoleSchema = type({
    user_id: "string",
    role: "'visitor' | 'standard' | 'partner' | 'administrator'"
})

// Types
export type ApproveUser = Infer<typeof approveUserSchema>
export type UpdateUserRole = Infer<typeof updateUserRoleSchema>

// Defaults
export const approveUserDefaults: ApproveUser = {
    user_id: "",
    role: "standard",
}

export const updateUserRoleDefaults: UpdateUserRole = {
    user_id: "",
    role: "standard",
}
