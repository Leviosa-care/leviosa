export type Category = {
    id: string;
    name: string;
}

// NOTE: this is what I get from the database to print the Products
export type CardType = {
    id: string;
    name: string;
    price: string;
    // TODO: change the following to have and ID and a name, so a custom type
    // category: Category;
    category: string;
    description: string;
    duration: number;
    image: string;
    updatedAt: string;
    published: "published" | "draft" | "archived";
    availability: "online" | "in-person" | "hybrid";
    bufferTime: number;
    cancellationHours: number;
}

export let cards: CardType[] = [
    {
        id: "8f275bfa-b7ba-476f-aabf-92d1e9ea5c75",
        name: "Relaxation package",
        price: "150.00",
        category: "massage",
        description:
            "Complete wellness package including massage, aromatherapy, and meditation session.",
        duration: 60,
        image: "https://placehold.co/360x200",
        updatedAt: "Updated: Jun 3, 2025, 02:34 PM",
        published: "published",
        availability: "online",
        bufferTime: 12,
        cancellationHours: 24,
    },
    {
        id: "c34021e4-9d95-4f84-903c-8bb2b3a0cd51",
        name: "Mental Clarity Coaching Session",
        price: "110.00",
        category: "mental coaching",
        description:
            "One-on-one coaching to help you refocus, reduce stress, and boost mental clarity through guided techniques.",
        duration: 60,
        image: "https://placehold.co/360x200",
        updatedAt: "Updated: Jul 5, 2025, 01:42 PM",
        published: "draft",
        availability: "online",
        bufferTime: 12,
        cancellationHours: 24,
    },
    {
        id: "ab97d532-476d-49f2-8469-c9060152ca63",
        name: "Hot stone package",
        price: "110.00",
        category: "massage",
        description:
            "Therapeutic massage using heated stones to relax muscles and improve circulation.",
        duration: 60,
        image: "https://placehold.co/360x200",
        updatedAt: "Updated: Jun 3, 2025, 02:34 PM",
        published: "draft",
        availability: "online",
        bufferTime: 12,
        cancellationHours: 24,
    },
    {
        id: "fa306985-7080-4ed9-a9e7-39f0c456c9ad",
        name: "Aromatherapy Oil",
        price: "25.00",
        category: "wellness",
        description:
            "Premium essential oil blend for relaxation and stress relief. ",
        duration: 60,
        image: "https://placehold.co/360x200",
        updatedAt: "Updated: Jun 3, 2025, 02:34 PM",
        published: "archived",
        availability: "in-person",
        bufferTime: 12,
        cancellationHours: 24,
    },
    {
        id: "375feaae-f889-484b-b174-2846546f8c25",
        name: "Deep Tissue massage",
        price: "90.00",
        category: "massage",
        description:
            "A more intense massage technique focused on the deeper layers of muscle tissue.",
        duration: 60,
        image: "https://placehold.co/360x200",
        updatedAt: "Updated: Jun 3, 2025, 02:34 PM",
        published: "published",
        availability: "hybrid",
        bufferTime: 12,
        cancellationHours: 24,
    },
    {
        id: "f5f111e2-90f4-46b9-a6ff-793bc9c6f9a1",
        name: "Swedish masage",
        price: "75.00",
        category: "massage",
        description:
            "A more intense massage technique focused on the deeper layers of muscle tissue.",
        duration: 60,
        image: "https://placehold.co/360x200",
        updatedAt: "Updated: Jun 3, 2025, 02:34 PM",
        published: "published",
        availability: "hybrid",
        bufferTime: 12,
        cancellationHours: 24,
    },]
