import type { EventState, MessageState, ConsultationState } from './Store';

export type Role = 'user' | 'userPremium' | 'helper' | 'admin' | 'freelance';
export type NavigationBarSize = 'small' | 'large';

export type NavigationBarElement = {
    label: string;
    href: string;
    icon: typeof import('lucide-svelte').Icon;
};

// NOTE: the old one that worked fine
// export type NavigationBarIcons = {
//     [key in Role]: {
//         [key in NavigationBarSize]: NavigationBarElement[]
//     }
// }

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
