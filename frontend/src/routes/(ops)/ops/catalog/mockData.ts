// Mock data for categories in development mode

export const mockCategories = [
    {
        id: "cat-1",
        name: "Massage Relaxant",
        description: "Massage doux pour la relaxation et le bien-être général. Parfait pour décompresser après une longue journée.",
        status: "published" as const,
        metadata: {},
        created_at: "2025-01-01T10:00:00Z",
        updated_at: "2025-01-15T14:30:00Z",
        images: [
            {
                id: "img-1",
                parent_id: "cat-1",
                parent_type: "category",
                url: "https://images.unsplash.com/photo-1544161515-4ab6ce6db874?w=800&auto=format&fit=crop",
                title: "Massage Relaxant",
                is_active: true,
                created_at: "2025-01-01T10:00:00Z",
            }
        ]
    },
    {
        id: "cat-2",
        name: "Massage Sportif",
        description: "Massage profond ciblant les muscles tendus. Idéal pour les sportifs et personnes actives.",
        status: "published" as const,
        metadata: {},
        created_at: "2025-01-02T11:00:00Z",
        updated_at: "2025-01-16T09:00:00Z",
        images: [
            {
                id: "img-2",
                parent_id: "cat-2",
                parent_type: "category",
                url: "https://images.unsplash.com/photo-1540555700478-4be289fbecef?w=800&auto=format&fit=crop",
                title: "Massage Sportif",
                is_active: true,
                created_at: "2025-01-02T11:00:00Z",
            }
        ]
    },
    {
        id: "cat-3",
        name: "Réflexologie",
        description: "Stimulation des points de pression des pieds pour améliorer la circulation et réduire le stress.",
        status: "published" as const,
        metadata: {},
        created_at: "2025-01-03T09:30:00Z",
        updated_at: "2025-01-17T16:45:00Z",
        images: [
            {
                id: "img-3",
                parent_id: "cat-3",
                parent_type: "category",
                url: "https://images.unsplash.com/photo-1519823551278-64ac92734fb1?w=800&auto=format&fit=crop",
                title: "Réflexologie",
                is_active: true,
                created_at: "2025-01-03T09:30:00Z",
            }
        ]
    },
    {
        id: "cat-4",
        name: "Massage Prénatal",
        description: "Massage adapté aux femmes enceintes pour soulager les tensions et favoriser la détente.",
        status: "draft" as const,
        metadata: {},
        created_at: "2025-01-04T14:00:00Z",
        updated_at: "2025-01-18T11:20:00Z",
        images: []
    },
    {
        id: "cat-5",
        name: "Aromathérapie",
        description: "Massage utilisant des huiles essentielles pour un bien-être physique et mental optimal.",
        status: "published" as const,
        metadata: {},
        created_at: "2025-01-05T10:15:00Z",
        updated_at: "2025-01-19T13:00:00Z",
        images: [
            {
                id: "img-5",
                parent_id: "cat-5",
                parent_type: "category",
                url: "https://images.unsplash.com/photo-1600334129128-685c5582fd35?w=800&auto=format&fit=crop",
                title: "Aromathérapie",
                is_active: true,
                created_at: "2025-01-05T10:15:00Z",
            }
        ]
    },
    {
        id: "cat-6",
        name: "Shiatsu",
        description: "Technique japonaise de massage par pression des doigts pour rééquilibrer l'énergie du corps.",
        status: "archived" as const,
        metadata: {},
        created_at: "2025-01-06T15:30:00Z",
        updated_at: "2025-01-20T10:00:00Z",
        images: [
            {
                id: "img-6",
                parent_id: "cat-6",
                parent_type: "category",
                url: "https://images.unsplash.com/photo-1590736969955-71cc94901144?w=800&auto=format&fit=crop",
                title: "Shiatsu",
                is_active: true,
                created_at: "2025-01-06T15:30:00Z",
            }
        ]
    },
]
