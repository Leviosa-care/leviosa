// TODO: make a role type and use it here as I should
import { ROLES } from "$lib/types/role";

export const mockUser: App.User = {
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
};
