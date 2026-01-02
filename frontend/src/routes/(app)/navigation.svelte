<script lang="ts">
    import { page } from "$app/state";

    interface Props {
        user: App.User;
    }

    let { user }: Props = $props();

    interface NavItem {
        href: string;
        label: string;
    }

    const navItems: NavItem[] = [
        { href: "/app", label: "Accueil" },
        { href: "/services", label: "Services" },
        { href: "/team", label: "Équipe" },
        { href: "/app/bookings", label: "Mes rendez-vous" },
        { href: "/app/account", label: "Compte" },
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

<div
    class="fixed top-0 left-0 right-0 z-50 bg-white/80 backdrop-blur-md border-b border-dark-100/50 flex items-center justify-between px-4 py-4 xl:grid xl:grid-cols-[1fr_auto_1fr] xl:py-4 xl:px-12"
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

    <!-- Desktop Navigation (hidden below 1280px) -->
    <nav class="hidden xl:block">
        <ul class="flex gap-12 items-center list-none p-0 m-0">
            {#each navItems as item}
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
    <div class="hidden xl:flex xl:items-center xl:justify-end xl:gap-3">
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

<!-- Spacer to prevent content from being hidden under fixed nav -->
<div class="h-[72px] xl:h-[64px]" aria-hidden="true"></div>
