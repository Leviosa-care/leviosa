import type { Role } from '$lib/types';

type User = {
    firstname: string,
    lastname: string,
    city: string,
    role: Role,
}

// NOTE: just for the time of developping the app
import { ROLES } from '$lib/types';
export const role: Role = ROLES.Basic;
export const mockUser: User = {
    firstname: 'John',
    lastname: 'DOE',
    city: 'Paris',
    role: role,
};
