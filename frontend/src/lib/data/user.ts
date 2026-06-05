import { ROLES, type Role } from "$lib/types/role";

const base = {
    birthdate: "",
    gender: "male" as const,
    telephone: "0123456789",
    postalCode: "75000",
    city: "Paris",
    address1: "01 Rue du capitaine framework",
    address2: "",
    google_id: "",
    apple_id: "",
    has_password: true,
    profile_incomplete: false,
};

export const mockAdminUser: App.User = {
    ...base,
    id: "123e4567-e89b-12d3-a456-426614174000",
    email: "admin@demo.leviosa.fr",
    picture: "",
    role: ROLES.admin,
    firstname: "Admin",
    lastname: "Demo",
};

export const mockPartnerUser: App.User = {
    ...base,
    id: "223e4567-e89b-12d3-a456-426614174001",
    email: "partner@demo.leviosa.fr",
    picture: "",
    role: ROLES.partner,
    firstname: "Partner",
    lastname: "Demo",
};

export const mockClientUser: App.User = {
    ...base,
    id: "323e4567-e89b-12d3-a456-426614174002",
    email: "client@demo.leviosa.fr",
    picture: "",
    role: ROLES.standard,
    firstname: "Client",
    lastname: "Demo",
};

export const mockUser = mockAdminUser;

export function getMockUserByRole(role: string): App.User | null {
    switch (role) {
        case ROLES.admin: return mockAdminUser;
        case ROLES.partner: return mockPartnerUser;
        case ROLES.standard: return mockClientUser;
        default: return null;
    }
}
