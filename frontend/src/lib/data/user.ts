import { ROLES } from "$lib/types/role";

export const mockUser: App.User = {
    id: "123e4567-e89b-12d3-a456-426614174000",
    email: "john.doe@example.com",
    picture: "",
    role: ROLES.admin,
    birthdate: "",
    firstname: "John",
    lastname: "DOE",
    gender: "male",
    telephone: "0123456789",
    postalCode: "75000",
    city: "Paris",
    address1: "01 Rue du capitaine framework",
    address2: "Chez monsieur Truc",
    google_id: "mock-google-id-12345",
    apple_id: "",
    has_password: true,
};
