export type Category = {
	id: string;
	name: string;
	description?: string;
	status?: "published" | "draft" | "archived";
	productCount?: number;
};

export type CardType = {
	id: string;
	name: string;
	price: string;
	category: string;
	description: string;
	duration: number;
	image: string;
	updatedAt: string;
	published: "published" | "draft" | "archived";
	availability: "online" | "in-person" | "hybrid";
	bufferTime: number;
	cancellationHours: number;
};

export const categories: Category[] = [
	{
		id: "cat-massage",
		name: "Massage",
		description:
			"Massages thérapeutiques et relaxants pour le bien-être du corps",
		status: "published",
		productCount: 5,
	},
	{
		id: "cat-wellness",
		name: "Wellness",
		description:
			"Soins et traitements pour la santé globale et la relaxation",
		status: "published",
		productCount: 3,
	},
	{
		id: "cat-mental",
		name: "Mental Coaching",
		description: "Séances de coaching mental et de développement personnel",
		status: "published",
		productCount: 2,
	},
	{
		id: "cat-nutrition",
		name: "Nutrition",
		description: "Conseils en nutrition et plans alimentaires personnalisés",
		status: "draft",
		productCount: 0,
	},
];

export let cards: CardType[] = [
	{
		id: "8f275bfa-b7ba-476f-aabf-92d1e9ea5c75",
		name: "Massage Relaxant",
		price: "75.00",
		category: "Massage",
		description:
			"Massage suédois classique pour détendre les muscles et réduire le stress.",
		duration: 60,
		image: "https://placehold.co/360x200",
		updatedAt: "2025-06-03T14:34:00Z",
		published: "published",
		availability: "in-person",
		bufferTime: 15,
		cancellationHours: 24,
	},
	{
		id: "c34021e4-9d95-4f84-903c-8bb2b3a0cd51",
		name: "Séance de Coaching Mental",
		price: "90.00",
		category: "Mental Coaching",
		description:
			"Séance individuelle de coaching pour améliorer la concentration et réduire l'anxiété.",
		duration: 60,
		image: "https://placehold.co/360x200",
		updatedAt: "2025-07-05T13:42:00Z",
		published: "published",
		availability: "online",
		bufferTime: 10,
		cancellationHours: 12,
	},
	{
		id: "ab97d532-476d-49f2-8469-c9060152ca63",
		name: "Massage aux Pierres Chaudes",
		price: "95.00",
		category: "Massage",
		description:
			"Massage thérapeutique utilisant des pierres chauffées pour relaxer les muscles en profondeur.",
		duration: 75,
		image: "https://placehold.co/360x200",
		updatedAt: "2025-06-10T16:20:00Z",
		published: "published",
		availability: "in-person",
		bufferTime: 15,
		cancellationHours: 24,
	},
	{
		id: "fa306985-7080-4ed9-a9e7-39f0c456c9ad",
		name: "Aromathérapie",
		price: "55.00",
		category: "Wellness",
		description:
			"Séance d'aromathérapie avec huiles essentielles pour la relaxation.",
		duration: 45,
		image: "https://placehold.co/360x200",
		updatedAt: "2025-05-15T11:00:00Z",
		published: "published",
		availability: "hybrid",
		bufferTime: 10,
		cancellationHours: 12,
	},
	{
		id: "375feaae-f889-484b-b174-2846546f8c25",
		name: "Massage Tissu Profond",
		price: "85.00",
		category: "Massage",
		description:
			"Massage intense ciblant les couches profondes des muscles pour soulager les tensions chroniques.",
		duration: 60,
		image: "https://placehold.co/360x200",
		updatedAt: "2025-06-20T09:30:00Z",
		published: "published",
		availability: "in-person",
		bufferTime: 15,
		cancellationHours: 24,
	},
	{
		id: "f5f111e2-90f4-46b9-a6ff-793bc9c6f9a1",
		name: "Coaching Performance",
		price: "120.00",
		category: "Mental Coaching",
		description:
			"Programme de coaching pour améliorer les performances mentales et la concentration.",
		duration: 90,
		image: "https://placehold.co/360x200",
		updatedAt: "2025-07-01T14:15:00Z",
		published: "draft",
		availability: "online",
		bufferTime: 10,
		cancellationHours: 24,
	},
	{
		id: "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
		name: "Massage Cranio-Sacral",
		price: "100.00",
		category: "Massage",
		description:
			"Technique douce travaillant sur le crâne et le sacrum pour rééquilibrer le système nerveux.",
		duration: 60,
		image: "https://placehold.co/360x200",
		updatedAt: "2025-04-10T10:00:00Z",
		published: "archived",
		availability: "in-person",
		bufferTime: 15,
		cancellationHours: 48,
	},
	{
		id: "b2c3d4e5-f6a7-8901-bcde-f12345678901",
		name: "Séance de Méditation Guidée",
		price: "45.00",
		category: "Wellness",
		description:
			"Séance de méditation guidée pour apprendre à gérer le stress et trouver la calme intérieur.",
		duration: 30,
		image: "https://placehold.co/360x200",
		updatedAt: "2025-05-20T15:45:00Z",
		published: "published",
		availability: "hybrid",
		bufferTime: 5,
		cancellationHours: 12,
	},
	{
		id: "c3d4e5f6-a7b8-9012-cdef-123456789012",
		name: "Massage Lomi Lomi",
		price: "110.00",
		category: "Massage",
		description:
			"Massage hawaïen traditionnel utilisant des mouvements fluides et rythmés.",
		duration: 90,
		image: "https://placehold.co/360x200",
		updatedAt: "2025-06-25T11:30:00Z",
		published: "published",
		availability: "in-person",
		bufferTime: 20,
		cancellationHours: 24,
	},
	{
		id: "d4e5f6a7-b8c9-0123-def0-234567890123",
		name: "Consultation Nutrition",
		price: "80.00",
		category: "Nutrition",
		description:
			"Première consultation nutritionnelle avec analyse des habitudes alimentaires.",
		duration: 60,
		image: "https://placehold.co/360x200",
		updatedAt: "2025-04-05T14:00:00Z",
		published: "draft",
		availability: "online",
		bufferTime: 10,
		cancellationHours: 24,
	},
];
