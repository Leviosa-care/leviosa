import type { EventState, MessageState, ConsultationState } from './Store';

export type Role = 'unknown' | 'anonymous' | ' basic' | 'premium' | 'guest' | 'admin' | 'freelance';

export type NavigationBarSize = 'small' | 'large';

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
