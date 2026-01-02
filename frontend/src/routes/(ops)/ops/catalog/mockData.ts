// Mock data for categories and products in development mode

export const mockCategories = [
    {
        id: "cat-1",
        name: "Bodywork & Massage",
        description: "Restore balance and alleviate tension with our therapeutic bodywork sessions. Tailored to your specific recovery needs.",
        status: "published" as const,
        metadata: {},
        createdAt: "2024-12-01T10:00:00Z",
        updatedAt: "2025-01-15T14:30:00Z",
        images: [
            {
                id: "img-cat-1",
                parent_id: "cat-1",
                parent_type: "category",
                url: "https://images.unsplash.com/photo-1544161515-4ab6ce6db874?w=800&auto=format&fit=crop",
                title: "Bodywork & Massage",
                is_active: true,
                created_at: "2024-12-01T10:00:00Z",
            }
        ]
    },
    {
        id: "cat-2",
        name: "Mindset Coaching",
        description: "Transform your mindset and unlock your potential with personalized coaching sessions designed for high-performers.",
        status: "published" as const,
        metadata: {},
        createdAt: "2024-12-02T11:00:00Z",
        updatedAt: "2025-01-16T09:00:00Z",
        images: [
            {
                id: "img-cat-2",
                parent_id: "cat-2",
                parent_type: "category",
                url: "https://images.unsplash.com/photo-1573496359142-b8d87734a5a2?w=800&auto=format&fit=crop",
                title: "Mindset Coaching",
                is_active: true,
                created_at: "2024-12-02T11:00:00Z",
            }
        ]
    },
    {
        id: "cat-3",
        name: "Physical Training",
        description: "Build strength, improve mobility, and optimize your physical performance with expert-guided training programs.",
        status: "published" as const,
        metadata: {},
        createdAt: "2024-12-03T09:30:00Z",
        updatedAt: "2025-01-17T16:45:00Z",
        images: [
            {
                id: "img-cat-3",
                parent_id: "cat-3",
                parent_type: "category",
                url: "https://images.unsplash.com/photo-1571019614242-c5c5dee9f50b?w=800&auto=format&fit=crop",
                title: "Physical Training",
                is_active: true,
                created_at: "2024-12-03T09:30:00Z",
            }
        ]
    },
    {
        id: "cat-4",
        name: "Wellness & Recovery",
        description: "Comprehensive wellness solutions focused on recovery, rejuvenation, and maintaining optimal health.",
        status: "draft" as const,
        metadata: {},
        createdAt: "2024-12-04T14:00:00Z",
        updatedAt: "2025-01-18T11:20:00Z",
        images: []
    },
];

export const mockProducts = [
    {
        id: "prod-1",
        name: "Deep Tissue Therapy",
        description: "Focus on realigning deep layers of muscle and connective tissue. It is especially helpful for chronic aches and pain in contracted areas.",
        category: "cat-1",
        duration: 60,
        status: "published" as const,
        availability: "in-person" as const,
        bufferTime: 10,
        cancellationHours: 24,
        stripeProductId: "prod_stripe_1",
        metadata: { intensity: "high", targetAreas: ["back", "shoulders", "neck"] },
        createdAt: "2025-01-01T10:00:00Z",
        updatedAt: "2025-01-15T14:30:00Z",
        images: [
            {
                id: "img-prod-1",
                parent_id: "prod-1",
                parent_type: "product",
                url: "https://images.unsplash.com/photo-1544161515-4ab6ce6db874?w=800&auto=format&fit=crop",
                title: "Deep Tissue Therapy",
                is_active: true,
                created_at: "2025-01-01T10:00:00Z",
            }
        ]
    },
    {
        id: "prod-2",
        name: "Lymphatic Drainage",
        description: "A gentle massage that encourages the movement of lymph fluids around the body. Helps remove waste and toxins.",
        category: "cat-1",
        duration: 60,
        status: "published" as const,
        availability: "in-person" as const,
        bufferTime: 10,
        cancellationHours: 24,
        stripeProductId: "prod_stripe_2",
        metadata: { intensity: "low", technique: "gentle" },
        createdAt: "2025-01-02T11:00:00Z",
        updatedAt: "2025-01-16T09:00:00Z",
        images: [
            {
                id: "img-prod-2",
                parent_id: "prod-2",
                parent_type: "product",
                url: "https://images.unsplash.com/photo-1540555700478-4be289fbecef?w=800&auto=format&fit=crop",
                title: "Lymphatic Drainage",
                is_active: true,
                created_at: "2025-01-02T11:00:00Z",
            }
        ]
    },
    {
        id: "prod-3",
        name: "Massage Sportif",
        description: "Massage profond ciblant les muscles tendus. Idéal pour les sportifs et personnes actives souhaitant améliorer leurs performances.",
        category: "cat-1",
        duration: 90,
        status: "published" as const,
        availability: "in-person" as const,
        bufferTime: 15,
        cancellationHours: 48,
        stripeProductId: "prod_stripe_3",
        metadata: { intensity: "high", targetAudience: "athletes" },
        createdAt: "2025-01-03T09:30:00Z",
        updatedAt: "2025-01-17T16:45:00Z",
        images: [
            {
                id: "img-prod-3",
                parent_id: "prod-3",
                parent_type: "product",
                url: "https://images.unsplash.com/photo-1519823551278-64ac92734fb1?w=800&auto=format&fit=crop",
                title: "Massage Sportif",
                is_active: true,
                created_at: "2025-01-03T09:30:00Z",
            }
        ]
    },
    {
        id: "prod-4",
        name: "Executive Performance Coaching",
        description: "One-on-one coaching for high-performers looking to optimize decision making and leadership presence.",
        category: "cat-2",
        duration: 60,
        status: "published" as const,
        availability: "online" as const,
        bufferTime: 5,
        cancellationHours: 48,
        stripeProductId: "prod_stripe_4",
        metadata: { sessionFormat: "1-on-1", focus: ["leadership", "decision-making"] },
        createdAt: "2025-01-04T14:00:00Z",
        updatedAt: "2025-01-18T11:20:00Z",
        images: [
            {
                id: "img-prod-4",
                parent_id: "prod-4",
                parent_type: "product",
                url: "https://images.unsplash.com/photo-1573496359142-b8d87734a5a2?w=800&auto=format&fit=crop",
                title: "Executive Performance Coaching",
                is_active: true,
                created_at: "2025-01-04T14:00:00Z",
            }
        ]
    },
    {
        id: "prod-5",
        name: "Strength Foundations",
        description: "Improve range of motion and joint health through functional movement patterns. Build a solid foundation for long-term fitness.",
        category: "cat-3",
        duration: 90,
        status: "published" as const,
        availability: "hybrid" as const,
        bufferTime: 15,
        cancellationHours: 24,
        stripeProductId: "prod_stripe_5",
        metadata: { skillLevel: "beginner-intermediate", equipment: ["bodyweight", "resistance-bands"] },
        createdAt: "2025-01-05T10:15:00Z",
        updatedAt: "2025-01-19T13:00:00Z",
        images: [
            {
                id: "img-prod-5",
                parent_id: "prod-5",
                parent_type: "product",
                url: "https://images.unsplash.com/photo-1571019614242-c5c5dee9f50b?w=800&auto=format&fit=crop",
                title: "Strength Foundations",
                is_active: true,
                created_at: "2025-01-05T10:15:00Z",
            }
        ]
    },
    {
        id: "prod-6",
        name: "Aromathérapie",
        description: "Massage utilisant des huiles essentielles pour un bien-être physique et mental optimal. Personnalisé selon vos besoins.",
        category: "cat-4",
        duration: 75,
        status: "draft" as const,
        availability: "in-person" as const,
        bufferTime: 10,
        cancellationHours: 24,
        stripeProductId: "",
        metadata: { essentialOils: true, customizable: true },
        createdAt: "2025-01-06T15:30:00Z",
        updatedAt: "2025-01-20T10:00:00Z",
        images: [
            {
                id: "img-prod-6",
                parent_id: "prod-6",
                parent_type: "product",
                url: "https://images.unsplash.com/photo-1600334129128-685c5582fd35?w=800&auto=format&fit=crop",
                title: "Aromathérapie",
                is_active: true,
                created_at: "2025-01-06T15:30:00Z",
            }
        ]
    },
];
