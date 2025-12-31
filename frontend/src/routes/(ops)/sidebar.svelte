<script lang="ts">
    import { page } from "$app/state";
    import {
        Home,
        Calendar,
        CalendarClock,
        CalendarCheck,
        Package,
        BarChart3,
        DollarSign,
        Users,
        Settings,
        ChevronLeft,
        ChevronRight,
    } from "@lucide/svelte";
    import type { Component } from "svelte";
    import type { Permissions } from "$lib/security/permissions";

    interface SidebarProps {
        data: {
            user: App.User;
            permissions: Permissions;
        };
    }

    let { data }: SidebarProps = $props();

    let { user, permissions } = data;

    // Sidebar collapse state
    let isCollapsed = $state(false);

    function toggleSidebar() {
        isCollapsed = !isCollapsed;
    }

    /**
     * Navigation item definition
     * - href: Route path
     * - label: Display text
     * - icon: Lucide icon component
     * - roles: Which roles can see this item
     */
    interface NavItem {
        href: string;
        label: string;
        icon: Component;
        roles: Array<"admin" | "partner">;
    }

    /**
     * Desktop Navigation Structure
     * More detailed and granular for larger screens
     */
    const desktopNavigation: NavItem[] = [
        {
            href: "/ops",
            label: "Accueil",
            icon: Home,
            roles: ["admin", "partner"],
        },
        // Agenda section - expanded on desktop
        {
            href: "/ops/agenda/disponibilites",
            label: "Disponibilités",
            icon: CalendarClock,
            roles: ["admin", "partner"],
        },
        {
            href: "/ops/agenda/reservations",
            label: "Réservations",
            icon: CalendarCheck,
            roles: ["admin", "partner"],
        },
        // Catalogue
        {
            href: "/ops/catalog",
            label: "Catalogue",
            icon: Package,
            roles: ["admin", "partner"],
        },
        // Statistics section - expanded on desktop
        {
            href: "/ops/statistics/analytics",
            label: "Analytics",
            icon: BarChart3,
            roles: ["admin", "partner"],
        },
        {
            href: "/ops/statistics/finances",
            label: "Finances",
            icon: DollarSign,
            roles: ["admin", "partner"],
        },
        // Users - admin only
        {
            href: "/ops/users",
            label: "Utilisateurs",
            icon: Users,
            roles: ["admin"],
        },
    ];

    /**
     * Mobile Navigation Structure
     * Consolidated and limited to 5 items max
     */
    const mobileNavigation: NavItem[] = [
        {
            href: "/ops",
            label: "Accueil",
            icon: Home,
            roles: ["admin", "partner"],
        },
        {
            href: "/ops/agenda",
            label: "Agenda",
            icon: Calendar,
            roles: ["admin", "partner"],
        },
        {
            href: "/ops/catalog",
            label: "Catalogue",
            icon: Package,
            roles: ["admin", "partner"],
        },
        {
            href: "/ops/statistics",
            label: "Statistiques",
            icon: BarChart3,
            roles: ["admin", "partner"],
        },
        {
            href: "/ops/users",
            label: "Utilisateurs",
            icon: Users,
            roles: ["admin"],
        },
    ];

    /**
     * Filter navigation items based on user role
     */
    const desktopItems = $derived(
        desktopNavigation.filter((item) =>
            item.roles.includes(user.role as "admin" | "partner"),
        ),
    );

    const mobileItems = $derived(
        mobileNavigation.filter((item) =>
            item.roles.includes(user.role as "admin" | "partner"),
        ),
    );

    /**
     * Check if a route is active
     * Exact match for home, starts-with for other routes
     */
    function isActive(href: string, currentPath: string): boolean {
        if (href === "/ops") {
            return currentPath === "/ops";
        }
        return currentPath.startsWith(href);
    }
</script>

{#if permissions.canAccessOps}
    <!-- Desktop Sidebar Navigation -->
    <aside
        class="hidden lg:flex lg:flex-col lg:border-r bg-black relative transition-all duration-300
               {isCollapsed ? 'lg:w-20' : 'lg:w-64'}"
        style="border-color: rgba(255, 255, 255, 0.1);"
        aria-label="Sidebar navigation"
    >
        <!-- Collapse Toggle Button -->
        <button
            onclick={toggleSidebar}
            class="absolute -right-4 bottom-24 z-10 w-8 h-8 rounded-full bg-black flex items-center justify-center
                   text-white/50 hover:text-white hover:bg-dark-900 transition-all duration-200"
            style="border: 1px solid rgba(255, 255, 255, 0.2);"
            aria-label={isCollapsed ? "Expand sidebar" : "Collapse sidebar"}
            title={isCollapsed ? "Expand sidebar" : "Collapse sidebar"}
        >
            {#if isCollapsed}
                <ChevronRight size={16} strokeWidth={2} />
            {:else}
                <ChevronLeft size={16} strokeWidth={2} />
            {/if}
        </button>

        <!-- Sidebar Header -->
        <div
            class="flex items-center justify-between py-6 {isCollapsed
                ? 'px-3'
                : 'px-6'}"
            style="border-bottom: 1px solid rgba(255, 255, 255, 0.1);"
        >
            {#if !isCollapsed}
                <h1
                    class="text-sm font-semibold tracking-tight text-white uppercase"
                >
                    Administration
                </h1>
                <a
                    href="/ops/settings"
                    class="p-2 rounded-md transition-all duration-200 text-white/50 hover:bg-white/5 hover:text-white"
                    aria-label="Paramètres"
                    title="Paramètres"
                >
                    <Settings size={18} strokeWidth={1.5} />
                </a>
            {:else}
                <a
                    href="/ops/settings"
                    class="p-2 rounded-md transition-all duration-200 text-white/50 hover:bg-white/5 hover:text-white mx-auto"
                    aria-label="Paramètres"
                    title="Paramètres"
                >
                    <Settings size={18} strokeWidth={1.5} />
                </a>
            {/if}
        </div>

        <!-- Navigation Items -->
        <nav class="flex-1 py-6 {isCollapsed ? 'px-2' : 'px-4'}">
            <ul class="space-y-0.5">
                {#each desktopItems as item (item.href)}
                    {@const active = isActive(item.href, page.url.pathname)}
                    <li>
                        <a
                            href={item.href}
                            class="flex items-center text-sm font-medium transition-all duration-200 rounded-md
                                   {isCollapsed
                                ? 'justify-center px-3 py-3'
                                : 'gap-3 px-3 py-2.5'}
                                   {active
                                ? 'text-white bg-white/10'
                                : 'text-white/50 hover:text-white/90 hover:bg-white/5'}"
                            aria-current={active ? "page" : undefined}
                            title={isCollapsed ? item.label : undefined}
                        >
                            <item.icon
                                size={18}
                                strokeWidth={active ? 2 : 1.5}
                            />
                            {#if !isCollapsed}
                                <span class="tracking-tight">{item.label}</span>
                            {/if}
                        </a>
                    </li>
                {/each}
            </ul>
        </nav>

        <!-- Sidebar Footer -->
        <div
            class="py-5 {isCollapsed ? 'px-3' : 'px-6'}"
            style="border-top: 1px solid rgba(255, 255, 255, 0.1);"
        >
            {#if !isCollapsed}
                <div class="flex items-center gap-3">
                    <div
                        class="w-9 h-9 rounded-full flex items-center justify-center bg-white/10"
                        style="border: 1px solid rgba(255, 255, 255, 0.15);"
                    >
                        <span
                            class="text-xs font-semibold text-white/70 uppercase"
                        >
                            {user.firstname?.[0] ?? "A"}
                        </span>
                    </div>
                    <div class="flex-1 min-w-0">
                        <p
                            class="text-sm font-medium text-white/95 truncate tracking-tight"
                        >
                            {user.firstname ?? "Admin"}
                        </p>
                        <p class="text-xs text-white/40 truncate">
                            {user.role === "admin"
                                ? "Administrateur"
                                : "Partenaire"}
                        </p>
                    </div>
                </div>
            {:else}
                <div
                    class="w-9 h-9 mx-auto rounded-full flex items-center justify-center bg-white/10"
                    style="border: 1px solid rgba(255, 255, 255, 0.15);"
                    title={user.firstname ?? "Admin"}
                >
                    <span class="text-xs font-semibold text-white/70 uppercase">
                        {user.firstname?.[0] ?? "A"}
                    </span>
                </div>
            {/if}
        </div>
    </aside>

    <!-- Mobile Bottom Navigation -->
    <nav
        class="lg:hidden fixed bottom-0 left-0 right-0 z-50 bg-black backdrop-blur-sm"
        style="border-top: 1px solid rgba(255, 255, 255, 0.1);"
        aria-label="Bottom navigation"
    >
        <ul class="flex items-center justify-around px-1 py-1">
            {#each mobileItems as item (item.href)}
                {@const active = isActive(item.href, page.url.pathname)}
                <li class="flex-1">
                    <a
                        href={item.href}
                        class="flex flex-col items-center gap-1.5 py-3 px-2 rounded-lg transition-all duration-200
                               {active ? 'text-white' : 'text-white/50'}"
                        aria-current={active ? "page" : undefined}
                    >
                        <item.icon size={20} strokeWidth={active ? 2 : 1.5} />
                        <span class="text-xs font-medium tracking-tight">
                            {item.label}
                        </span>
                    </a>
                </li>
            {/each}
        </ul>
    </nav>
{/if}
