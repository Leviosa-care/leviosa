// types
import type { Role } from "$lib/types/role"

export const NAVIGATION_BAR_SIZES = {
    Small: 'small',
    Large: 'large',
} as const;
export type NavigationBarSize = typeof NAVIGATION_BAR_SIZES[keyof typeof NAVIGATION_BAR_SIZES];
export const NAVIGATION_BAR_SIZES_ARRAY = Object.values(NAVIGATION_BAR_SIZES);
export type NavigationBarElement = {
    label: string;
    href: string;
    icon: typeof import('@lucide/svelte').Icon;
};

export type IconsByRole = Record<NavigationBarSize, NavigationBarElement[]>
export type NavigationBarIcons = Record<Role, IconsByRole>;

import {
    Home,
    Slack,
    MessageSquare,
    Calendar,
    User,
    NotebookPen,
    Ticket,
    Users
} from '@lucide/svelte';

const visitor = {
    small: [],
    large: [],
}
export const basic: IconsByRole = {
    small: [
        { label: 'accueil', icon: Home, href: '/portal/' },
        { label: 'messages', icon: MessageSquare, href: '/portal/messages' },
        { label: 'services', icon: Slack, href: '/portal/services' },
        { label: 'reservations', icon: Calendar, href: '/portal/reservations' },
        { label: 'profil', icon: User, href: '/portal/profile' }
    ],
    large: [
        { label: 'accueil', icon: Home, href: '/portal/' },
        { label: 'conversations', icon: MessageSquare, href: '/portal/messages' },
        { label: 'Note de seance', icon: NotebookPen, href: '/portal/messages' },
        { label: 'services', icon: Slack, href: '/portal/services' },
        { label: 'reservations', icon: Calendar, href: '/portal/reservations' },
        { label: 'profil', icon: User, href: '/portal/profile' }
    ]
}
export const premium: IconsByRole = {
    small: [
        { label: 'accueil', icon: Home, href: '/premium/' },
        { label: 'messages', icon: MessageSquare, href: '/premium/messages' },
        { label: 'services', icon: Slack, href: '/premium/services' },
        { label: 'reservations', icon: Calendar, href: '/premium/reservations' },
        { label: 'profil', icon: User, href: '/premium/profile' }
    ],
    large: [
        { label: 'accueil', icon: Home, href: '/premium/' },
        { label: 'conversations', icon: MessageSquare, href: '/premium/messages' },
        { label: 'notes de seances', icon: NotebookPen, href: '/premium/messages' },
        { label: 'services', icon: Slack, href: '/premium/services' },
        { label: 'evenements', icon: Ticket, href: '/premium/reservations' },
        { label: 'consultations', icon: Calendar, href: '/premium/reservations' },
        { label: 'profil', icon: User, href: '/premium/profile' }
    ]
}

export const guest: IconsByRole = {
    small: [
        { label: 'accueil', icon: Home, href: '/guest/' },
        { label: 'messages', icon: MessageSquare, href: '/guest/messages' },
        { label: 'events', icon: Calendar, href: '/guest/events' },
        { label: 'profil', icon: User, href: '/guest/profile' }
    ],
    large: []
}

export const partners: IconsByRole = {
    small: [
        { label: 'accueil', icon: Home, href: '/partners/' },
        { label: 'messages', icon: MessageSquare, href: '/partners/messages' },
        { label: 'services', icon: Slack, href: '/partners/services' },
        { label: 'reservations', icon: Calendar, href: '/partners/reservations' },
        { label: 'profil', icon: User, href: '/partners/profile' }
    ],
    large: [
        { label: 'accueil', icon: Home, href: '/partners/' },
        { label: 'conversations', icon: MessageSquare, href: '/partners/messages' },
        { label: 'notes de seances', icon: NotebookPen, href: '/partners/messages' },
        { label: 'services', icon: Slack, href: '/partners/services' },
        { label: 'evenements', icon: Ticket, href: '/partners/reservations' },
        { label: 'consultations', icon: Calendar, href: '/partners/reservations' },
        { label: 'profil', icon: User, href: '/partners/profile' }
    ]
}

export const admin: IconsByRole = {
    small: [
        { label: 'accueil', icon: Home, href: '/admin/' },
        { label: 'messages', icon: MessageSquare, href: '/admin/messages' },
        { label: 'services', icon: Slack, href: '/admin/services' },
        { label: 'reservations', icon: Calendar, href: '/admin/reservations' },
        { label: 'users', icon: Users, href: '/admin/users' }
    ],
    large: [
        { label: 'accueil', icon: Home, href: '/admin/' },
        { label: 'conversations', icon: MessageSquare, href: '/admin/messages' },
        { label: 'notes de seances', icon: NotebookPen, href: '/admin/messages' },
        { label: 'services', icon: Slack, href: '/admin/services' },
        { label: 'evenements', icon: Ticket, href: '/admin/reservations' },
        { label: 'consultations', icon: Calendar, href: '/admin/reservations' },
        { label: 'users', icon: Users, href: '/admin/users' }
    ]
}

export const navigationBarIcons: NavigationBarIcons = {
    visitor,
    standard: basic,
    premium,
    guest,
    partner: partners,
    administrator: admin,
    // TODO: add as many role as needed
    // bodyguard: {},
    // photograph: {},
};
