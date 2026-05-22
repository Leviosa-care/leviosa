import { ROLES } from "$lib/types/role";

export const homeByRole = {
    [ROLES.visitor]: "/portal",
    [ROLES.standard]: "/app",
    [ROLES.premium]: "/premium",
    [ROLES.guest]: "/guest",
    [ROLES.partner]: "/partners",
    [ROLES.admin]: "/admin",
}
