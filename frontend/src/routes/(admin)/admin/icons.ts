import { Home, Mail, Package, ChartSpline, CalendarDays, NotebookPen, Users, HandCoins, Building2 } from "@lucide/svelte";

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
                icon: Package,
                name: "Catalogue",
                link: "/admin/catalog",
            },
            {
                icon: Building2,
                name: "Bâtiments",
                link: "/admin/buildings",
            },
            {
                icon: ChartSpline,
                name: "Analytics",
                link: "/admin/analytics",
            },
            {
                icon: HandCoins,
                name: "Comptabilité",
                link: "/admin/compta",
            },
        ],
} as const;
