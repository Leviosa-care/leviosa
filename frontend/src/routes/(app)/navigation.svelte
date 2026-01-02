<script lang="ts">
    import { page } from "$app/state";
    import {
        Home,
        Briefcase,
        CalendarCheck,
        Activity,
        User,
    } from "@lucide/svelte";
    import type { Component } from "svelte";

    interface Props {
        user: App.User;
    }

    let { user }: Props = $props();

    interface NavItem {
        href: string;
        label: string;
    }

    interface MobileNavItem extends NavItem {
        icon: Component;
    }

    /**
     * Desktop Navigation - all items
     */
    const desktopNavItems: NavItem[] = [
        { href: "/app", label: "Accueil" },
        { href: "/services", label: "Services" },
        { href: "/team", label: "Équipe" },
        { href: "/app/bookings", label: "Mes rendez-vous" },
        { href: "/app/account", label: "Compte" },
    ];

    /**
     * Mobile Navigation - consolidated to 5 items max
     */
    const mobileNavItems: MobileNavItem[] = [
        { href: "/app", label: "Accueil", icon: Home },
        { href: "/services", label: "Services", icon: Briefcase },
        { href: "/app/bookings", label: "Réservations", icon: CalendarCheck },
        { href: "/app/activity", label: "Activité", icon: Activity },
        { href: "/app/account", label: "Compte", icon: User },
    ];

    function isActive(href: string): boolean {
        const currentPath = page.url.pathname;

        // Exact match for /app home
        if (href === "/app") {
            return currentPath === "/app";
        }

        // For other routes, check if current path starts with href
        // and ensure we're matching full path segments (not partial)
        return currentPath === href || currentPath.startsWith(href + "/");
    }
</script>

<!-- Desktop Top Navigation (hidden below 1280px) -->
<div
    class="hidden xl:flex fixed top-0 left-0 right-0 z-50 bg-white/80 backdrop-blur-md border-b border-dark-100/50 items-center justify-between px-4 py-4 xl:grid xl:grid-cols-[1fr_auto_1fr] xl:py-4 xl:px-12"
>
    <!-- Logo/Brand -->
    <div class="xl:justify-start">
        <a
            href="/app"
            class="text-lg font-semibold tracking-tight text-dark-900 hover:text-dark-700 transition-colors"
        >
            Leviosa
        </a>
    </div>

    <!-- Desktop Navigation -->
    <nav>
        <ul class="flex gap-12 items-center list-none p-0 m-0">
            {#each desktopNavItems as item}
                <li>
                    <a
                        class="transition-colors {isActive(item.href)
                            ? 'text-dark-900 font-semibold'
                            : 'text-dark-500 hover:text-dark-900'}"
                        href={item.href}
                    >
                        {item.label}
                    </a>
                </li>
            {/each}
        </ul>
    </nav>

    <!-- User Section -->
    <div class="flex items-center justify-end gap-3">
        <span class="text-sm text-dark-500">
            {user.firstname}
            {user.lastname}
        </span>
        {#if user.picture}
            <img
                src={user.picture}
                alt={`${user.firstname} ${user.lastname}`}
                class="w-8 h-8 rounded-full object-cover ring-1 ring-dark-200"
            />
        {:else}
            <div
                class="w-8 h-8 rounded-full bg-dark-100 flex items-center justify-center ring-1 ring-dark-200"
            >
                <span class="text-xs font-medium text-dark-700">
                    {user.firstname.charAt(0)}{user.lastname.charAt(0)}
                </span>
            </div>
        {/if}
    </div>
</div>

<!-- Mobile Bottom Navigation (hidden on desktop) -->
<nav
    class="xl:hidden fixed bottom-0 left-0 right-0 z-50 bg-white/95 backdrop-blur-md border-t border-dark-100/50"
    aria-label="Bottom navigation"
>
    <ul
        class="flex items-center justify-around px-1 py-1 safe-area-inset-bottom"
    >
        {#each mobileNavItems as item (item.href)}
            {@const active = isActive(item.href)}
            <li class="flex-1">
                <a
                    href={item.href}
                    class="flex flex-col items-center gap-1.5 py-3 px-2 rounded-lg transition-all duration-200
                           {active ? 'text-dark-900' : 'text-dark-400'}"
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

<!-- Spacer to prevent content from being hidden under fixed nav -->
<div class="h-0 xl:h-[64px]" aria-hidden="true"></div>
