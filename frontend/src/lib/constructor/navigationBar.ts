import type { NavigationBarIcons } from '$lib/types';

import {
    Home,
    Slack,
    MessageSquare,
    Calendar,
    User,
    NotebookPen,
    Ticket,
    Users
} from 'lucide-svelte';

const user = {
    small: [
        { label: 'accueil', icon: Home, href: '/app/' },
        {
            label: 'messages',
            icon: MessageSquare,
            href: '/app/messages'
        },
        { label: 'services', icon: Slack, href: '/app/services' },
        {
            label: 'reservations',
            icon: Calendar,
            href: '/app/reservations'
        },
        { label: 'profil', icon: User, href: '/app/profile' }
    ],
    large: [
        { label: 'accueil', icon: Home, href: '/app/' },
        {
            label: 'conversations',
            icon: MessageSquare,
            href: '/app/messages'
        },
        {
            label: 'Note de seance',
            icon: NotebookPen,
            href: '/app/messages'
        },
        { label: 'services', icon: Slack, href: '/app/services' },
        {
            label: 'reservations',
            icon: Calendar,
            href: '/app/reservations'
        },
        { label: 'profil', icon: User, href: '/app/profile' }
    ]
}
const userPremium = {
    small: [
        { label: 'accueil', icon: Home, href: '/app/' },
        {
            label: 'messages',
            icon: MessageSquare,
            href: '/app/messages'
        },
        { label: 'services', icon: Slack, href: '/app/services' },
        {
            label: 'reservations',
            icon: Calendar,
            href: '/app/reservations'
        },
        { label: 'profil', icon: User, href: '/app/profile' }
    ],
    large: [
        { label: 'accueil', icon: Home, href: '/app/' },
        {
            label: 'conversations',
            icon: MessageSquare,
            href: '/app/messages'
        },
        {
            label: 'notes de seances',
            icon: NotebookPen,
            href: '/app/messages'
        },
        { label: 'services', icon: Slack, href: '/app/services' },
        {
            label: 'evenements',
            icon: Ticket,
            href: '/app/reservations'
        },
        {
            label: 'consultations',
            icon: Calendar,
            href: '/app/reservations'
        },
        { label: 'profil', icon: User, href: '/app/profile' }
    ]
}

const helper = {
    small: [
        { label: 'accueil', icon: Home, href: '/app/' },
        {
            label: 'messages',
            icon: MessageSquare,
            href: '/app/messages'
        },
        {
            label: 'events',
            icon: Calendar,
            href: '/app/events'
        },
        { label: 'profil', icon: User, href: '/app/profile' }
    ],
    large: []
}

const freelance = {
    small: [
        { label: 'accueil', icon: Home, href: '/app/' },
        {
            label: 'messages',
            icon: MessageSquare,
            href: '/app/messages'
        },
        { label: 'services', icon: Slack, href: '/app/services' },
        {
            label: 'reservations',
            icon: Calendar,
            href: '/app/reservations'
        },
        { label: 'profil', icon: User, href: '/app/profile' }
    ],
    large: [
        { label: 'accueil', icon: Home, href: '/app/' },
        {
            label: 'conversations',
            icon: MessageSquare,
            href: '/app/messages'
        },
        {
            label: 'notes de seances',
            icon: NotebookPen,
            href: '/app/messages'
        },
        { label: 'services', icon: Slack, href: '/app/services' },
        {
            label: 'evenements',
            icon: Ticket,
            href: '/app/reservations'
        },
        {
            label: 'consultations',
            icon: Calendar,
            href: '/app/reservations'
        },
        { label: 'profil', icon: User, href: '/app/profile' }
    ]
}

const admin = {
    small: [
        { label: 'accueil', icon: Home, href: '/app/' },
        {
            label: 'messages',
            icon: MessageSquare,
            href: '/app/messages'
        },
        { label: 'services', icon: Slack, href: '/app/services' },
        {
            label: 'reservations',
            icon: Calendar,
            href: '/app/reservations'
        },
        { label: 'users', icon: Users, href: '/app/users' }
    ],
    large: [
        { label: 'accueil', icon: Home, href: '/app/' },
        {
            label: 'conversations',
            icon: MessageSquare,
            href: '/app/messages'
        },
        {
            label: 'notes de seances',
            icon: NotebookPen,
            href: '/app/messages'
        },
        { label: 'services', icon: Slack, href: '/app/services' },
        {
            label: 'evenements',
            icon: Ticket,
            href: '/app/reservations'
        },
        {
            label: 'consultations',
            icon: Calendar,
            href: '/app/reservations'
        },
        { label: 'users', icon: Users, href: '/app/users' }
    ]
}
export const navigationBarIcons: NavigationBarIcons = {
    user,
    userPremium,
    helper,
    freelance,
    admin
    // TODO: add as many role as needed
    // bodyguard: {},
    // photograph: {},
};
