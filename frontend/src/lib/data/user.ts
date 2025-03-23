import type { Role } from '$lib/types';

type User = {
    firstname: string,
    lastname: string,
    city: string,
    role: Role,
}

// NOTE: just for the time of developping the app
export const role: Role = 'admin';
export const mockUser: User = {
    firstname: 'John',
    lastname: 'DOE',
    city: 'Paris',
    role: role,
};
