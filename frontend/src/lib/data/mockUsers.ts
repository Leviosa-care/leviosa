import type { Role } from "$lib/types/role";

export interface MockUser {
	id: string;
	email: string;
	role: Role;
	status: "approved" | "pending";
	createdAt: string;
	firstname?: string;
	lastname?: string;
	telephone?: string;
}

export const mockUsers: MockUser[] = [
	{
		id: "1",
		email: "sophie.martin@example.com",
		role: "admin",
		status: "approved",
		createdAt: "2024-01-15T10:30:00Z",
		firstname: "Sophie",
		lastname: "Martin",
		telephone: "+33612345678"
	},
	{
		id: "2",
		email: "pierre.leroy@example.com",
		role: "standard",
		status: "approved",
		createdAt: "2024-02-10T14:20:00Z",
		firstname: "Pierre",
		lastname: "Leroy",
		telephone: "+33623456789"
	},
	{
		id: "3",
		email: "marie.dubois@example.com",
		role: "standard",
		status: "approved",
		createdAt: "2024-02-15T09:45:00Z",
		firstname: "Marie",
		lastname: "Dubois",
		telephone: "+33634567890"
	},
	{
		id: "4",
		email: "jean.dupont@example.com",
		role: "partner",
		status: "approved",
		createdAt: "2024-03-01T16:00:00Z",
		firstname: "Jean",
		lastname: "Dupont",
		telephone: "+33645678901"
	},
	{
		id: "5",
		email: "claire.bernard@example.com",
		role: "standard",
		status: "pending",
		createdAt: "2024-03-10T11:30:00Z",
		firstname: "Claire",
		lastname: "Bernard",
		telephone: "+33656789012"
	},
	{
		id: "6",
		email: "lucas.petit@example.com",
		role: "standard",
		status: "pending",
		createdAt: "2024-03-12T08:15:00Z",
		firstname: "Lucas",
		lastname: "Petit",
		telephone: "+33678901234"
	},
	{
		id: "7",
		email: "emma.moreau@example.com",
		role: "premium",
		status: "approved",
		createdAt: "2024-03-05T13:20:00Z",
		firstname: "Emma",
		lastname: "Moreau",
		telephone: "+33689012345"
	},
	{
		id: "8",
		email: "thomas.richard@example.com",
		role: "standard",
		status: "approved",
		createdAt: "2024-03-08T10:00:00Z",
		firstname: "Thomas",
		lastname: "Richard",
		telephone: "+33690123456"
	}
];
