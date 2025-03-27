// See https://kit.svelte.dev/docs/types#app
// for information about these interfaces
import type { Role } from '$lib/types';

declare global {
    namespace App {
        interface User {
            email?: string;
            role: Role;
            lastname?: string;
            firstname?: string;
            gender?: string;
            birthdate?: string;
            telephone?: string;
            picture?: string;
            address1?: string;
            address2?: string;
            city?: string;
            postalCode?: number;
        }
        interface Locals {
            user?: User
        }
        interface PageData {
            user?: User
        }
        // interface Error {}
        // interface PageState {}
        // interface Platform {}
    }
}

export { };
