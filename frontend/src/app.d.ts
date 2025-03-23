// See https://kit.svelte.dev/docs/types#app
// for information about these interfaces
declare global {
    namespace App {
        // interface Error {}
        interface Locals {
            // TODO: remove the "?" when done with the code
            user: {
                email?: string;
                role: string;
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
            };
        }
        // interface PageData {}
        // interface PageState {}
        // interface Platform {}
    }
}

export { };
