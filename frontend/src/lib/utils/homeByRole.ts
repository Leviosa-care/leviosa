import { ROLES } from "$lib/types/role";

export const homeByRole = {
    [ROLES.visitor]: "/portal",
    [ROLES.basic]: "/app",
    [ROLES.premium]: "/premium",
    [ROLES.guest]: "/guest",
    [ROLES.partners]: "/partners",
    [ROLES.admin]: "/admin",
}
