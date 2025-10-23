import { Home, Mail, Package, ChartSpline, CalendarDays, NotebookPen, Ticket, Users, Server, HandCoins } from "@lucide/svelte";

export type Icon = {
    icon: typeof import("@lucide/svelte").Icon;
    name: string;
    link: string;
};

export type IconType = "small" | "large"

export const icons: Record<IconType, Icon[]> = {
    small: [],
    large:
        [
            {
                icon: Home,
                name: "Accueil",
                link: "/admin",
            },
            {
                icon: Users,
                name: "Utilisateurs",
                link: "/admin/users",
            },
            {
                icon: CalendarDays,
                name: "Planning",
                link: "/admin/planning",
            },
            {
                icon: NotebookPen,
                name: "Notes de seance",
                link: "/admin/bookings/consultations",
            },
            {
                icon: Mail,
                name: "Messages",
                link: "/admin/messages",
            },
            {
                icon: Ticket,
                name: "Evenements",
                link: "/admin/bookings/events",
            },
            {
                icon: Package,
                name: "Catalogue",
                link: "/admin/products",
            },
            {
                icon: ChartSpline,
                name: "Analytics",
                link: "/admin/analytics",
            },
            {
                icon: Server,
                name: "Infrastructure",
                link: "/admin/infra",
            },
            {
                icon: HandCoins,
                name: "Comptabilité",
                link: "/admin/compta",
            },
        ],
} as const;
