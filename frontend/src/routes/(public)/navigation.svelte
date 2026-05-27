<script lang="ts">
    import Button from "$lib/ui/Button.svelte";
    import SideDrawer from "$lib/ui/SideDrawer.svelte";
    import Logo from "./__logo.svelte";
    import { page } from "$app/state";
    import { type Permissions } from "$lib/security/permissions";
    import { Menu, X } from "@lucide/svelte";

    interface Props {
        permissions: Permissions;
    }
    let { permissions }: Props = $props();

    let isMobileMenuOpen = $state(false);

    interface NavItem {
        link: string;
        title: string;
    }

    const items: NavItem[] = [
        { link: "/services", title: "Services" },
        { link: "/team", title: "L’équipe" },
        { link: "/about", title: "À propos" },
        { link: "/bookings", title: "Mes réservations" },
    ];
    // Services · L’équipe · À propos · Réserver

    function isActive(path: string): boolean {
        return page.url.pathname.startsWith(path);
    }
    const secondaryButtonText =
        permissions.canAccessOps || permissions.canAccessApp
            ? "Voir ton profil"
            : "S'authentifier";

    const secondaryButtonLink = permissions.canAccessOps
        ? "/staff"
        : permissions.canAccessApp
          ? "/app"
          : "/auth";

    // Close drawer when resizing to desktop (xl: 1280px)
    $effect(() => {
        if (typeof window === "undefined") return;

        function handleResize() {
            if (window.innerWidth >= 1280 && isMobileMenuOpen) {
                isMobileMenuOpen = false;
            }
        }

        window.addEventListener("resize", handleResize);
        return () => window.removeEventListener("resize", handleResize);
    });
</script>

<svelte:window
    onkeydown={(e) => {
        if (e.key === "Escape") isMobileMenuOpen = false;
    }}
/>

<div
    class="fixed top-0 left-0 right-0 z-50 bg-white/80 backdrop-blur-md border-b border-dark-100/50 flex items-center justify-between px-4 py-4 xl:grid xl:grid-cols-[1fr_auto_1fr] xl:py-4 xl:px-12"
>
    <!-- Mobile/Tablet Hamburger (hidden on desktop) -->
    <button
        onclick={() => (isMobileMenuOpen = true)}
        aria-label="Open navigation menu"
        aria-expanded={isMobileMenuOpen}
        class="flex items-center justify-center xl:hidden w-11 h-11 cursor-pointer"
    >
        <Menu size={24} />
    </button>

    <!-- Logo -->
    <div class="xl:justify-start">
        <Logo />
    </div>

    <!-- Desktop Navigation (hidden below 1280px) -->
    <nav class="hidden xl:block">
        <ul class="flex gap-12 items-center list-none p-0 m-0">
            {#each items as item}
                <li>
                    <a
                        class="transition-colors {isActive(item.link)
                            ? 'text-dark-900 font-semibold'
                            : 'text-dark-500 hover:text-dark-900'}"
                        href={item.link}
                    >
                        {item.title}
                    </a>
                </li>
            {/each}
        </ul>
    </nav>

    <!-- Desktop Actions (hidden below 1280px) -->
    <div class="hidden xl:flex xl:gap-4 xl:items-center xl:justify-end">
        <a href="/book"><Button class="rounded-xl py-3 text-white">Réserver maintenant</Button></a>
        <Button class="rounded-xl py-3 bg-white border border-dark-900 hover:bg-dark-50">
            <a href={secondaryButtonLink}>
                {secondaryButtonText}
            </a>
        </Button>
    </div>
</div>

<!-- Mobile Drawer -->
<SideDrawer bind:isOpen={isMobileMenuOpen}>
    <div class="flex flex-col h-full">
        <!-- Close button (top-right) -->
        <div class="flex justify-end mb-6">
            <button
                onclick={() => (isMobileMenuOpen = false)}
                aria-label="Close navigation menu"
                class="w-11 h-11 flex items-center justify-center cursor-pointer"
            >
                <X size={24} />
            </button>
        </div>

        <!-- Navigation items (stacked vertically) -->
        <nav class="flex-1">
            <ul class="flex flex-col gap-1 list-none p-0 m-0">
                {#each items as item}
                    <li>
                        <a
                            href={item.link}
                            onclick={() => (isMobileMenuOpen = false)}
                            class="block py-3 px-6 rounded-lg min-h-[44px] transition-colors {isActive(
                                item.link,
                            )
                                ? 'bg-dark-100 text-dark-900 font-semibold'
                                : 'text-dark-500 hover:bg-dark-50 hover:text-dark-900'}"
                        >
                            {item.title}
                        </a>
                    </li>
                {/each}
            </ul>
        </nav>

        <!-- Divider -->
        <div class="h-px bg-dark-200 my-6"></div>

        <!-- Action buttons (full width, stacked) -->
        <div class="flex flex-col gap-3">
            <a href="/book"><Button
                class="w-full inline-flex justify-center items-center px-8 py-3.5 bg-dark-900 hover:bg-dark-800 text-background text-sm sm:text-base font-medium transition-all duration-200 shadow-sm hover:shadow-md ring-offset-2 focus:ring-2 focus:ring-dark-900"
                >Réserver maintenant</Button></a>
            <Button
                class="w-full bg-background border border-dark-900 hover:bg-dark-50 px-8 py-3.5 rounded-xl"
            >
                <a href={secondaryButtonLink} class="block w-full">
                    {secondaryButtonText}
                </a>
            </Button>
        </div>
    </div>
</SideDrawer>
