import type { Role } from '$lib/types';

export const role: Role = 'admin';
interface User {
    firstname: string,
    lastname: string,
    city: string,
    role: Role,
}
export const mockUser: User = {
    firstname: 'John',
    lastname: 'DOE',
    city: 'Paris',
    // role: 'userPremium',
    role: 'admin',
    // role: 'user',
    // role: 'freelance',
};
