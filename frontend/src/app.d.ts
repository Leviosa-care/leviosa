// See https://svelte.dev/docs/kit/types#app.d.ts
// for information about these interfaces
import { type Role } from "$lib/types/role";

declare global {
    namespace App {
        interface User {
            email: string;
            picture: string;
            role: Role;
            birthdate: string;
            firstname: string;
            lastname: string;
            gender: string;
            telephone: string;
            postalCode: string;
            city: string;
            address1: string;
            address2: string;
        }
        interface Error { }
        interface Locals {
            user: User
            isAdminDomain: boolean
            isStaffDomain: boolean
            sessionCookieName: string
            cookieDomain: string | undefined
        }
        // interface PageData {}
        // interface PageState {}
        // interface Platform {}
    }
}

export { };
