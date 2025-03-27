import type { EventState, MessageState, ConsultationState } from './Store';

export const ROLES = {
    Unknown: 'unknown',
    Anonymous: 'anonymous',
    Basic: 'basic',
    Premium: 'premium',
    Guest: 'guest',
    Admin: 'admin',
    Freelance: 'freelance',
} as const;
export type Role = typeof ROLES[keyof typeof ROLES];
export const ROLES_ARRAY = Object.values(ROLES);


export const NAVIGATION_BAR_SIZES = {
    Small: 'small',
    Large: 'large',
} as const;
export type NavigationBarSize = typeof NAVIGATION_BAR_SIZES[keyof typeof NAVIGATION_BAR_SIZES];
export const NAVIGATION_BAR_SIZES_ARRAY = Object.values(NAVIGATION_BAR_SIZES);

export type NavigationBarElement = {
    label: string;
    href: string;
    icon: typeof import('lucide-svelte').Icon;
};

export type NavigationBarIcons = Record<Role, Record<NavigationBarSize, NavigationBarElement[]>>;

type EventTabType = {
    name: EventState;
    href: string;
};

export type EventTabs = Record<Role, EventTabType[]>;

type MessageTabType = {
    name: MessageState;
    href: string;
};

export type MessageTabs = Record<Role, MessageTabType[]>;

type ConsultationTabType = {
    name: ConsultationState;
    href: string;
};

export type ConsultationTabs = Record<Role, ConsultationTabType[]>;
