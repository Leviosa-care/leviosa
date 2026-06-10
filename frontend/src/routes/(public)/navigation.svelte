<script lang="ts">
    import Button from "$lib/ui/Button.svelte";
    import SideDrawer from "$lib/ui/SideDrawer.svelte";
    import Logo from "./__logo.svelte";
    import { page } from "$app/state";
    import { onMount, onDestroy } from "svelte";
    import { browser } from "$app/environment";
    import { type Permissions } from "$lib/security/permissions";
    import { Menu, X } from "@lucide/svelte";

    interface Props {
        permissions: Permissions;
    }
    let { permissions }: Props = $props();

    let isMobileMenuOpen = $state(false);
    let scrolled = $state(false);

    const isHeroPage = $derived(page.url.pathname === "/");
    const showScrolledState = $derived(!isHeroPage || scrolled);

    function handleScroll() {
        if (!isHeroPage) { scrolled = false; return; }
        scrolled = window.scrollY > window.innerHeight;
    }

    onMount(() => {
        if (!browser) return;
        window.addEventListener("scroll", handleScroll, { passive: true });
        handleScroll();
    });

    onDestroy(() => {
        if (!browser) return;
        window.removeEventListener("scroll", handleScroll);
    });

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

    function isActive(path: string): boolean {
        return page.url.pathname.startsWith(path);
    }

    const secondaryButtonText =
        permissions.canAccessOps || permissions.canAccessApp
            ? "Voir ton profil"
            : "S’authentifier";

    const secondaryButtonLink = permissions.canAccessOps
        ? "/staff"
        : permissions.canAccessApp
          ? "/app"
          : "/auth";

    $effect(() => {
        if (!browser) return;
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
    class="fixed top-0 left-0 right-0 z-50 flex items-center justify-between px-4 py-4 xl:grid xl:grid-cols-[1fr_auto_1fr] xl:py-4 xl:px-12 border-b transition-all duration-300 {showScrolledState
        ? 'bg-white/80 backdrop-blur-md border-border-input/50'
        : 'bg-transparent border-transparent'}"
>
    <!-- Mobile/Tablet Hamburger (hidden on desktop) -->
    <button
        onclick={() => (isMobileMenuOpen = true)}
        aria-label="Open navigation menu"
        aria-expanded={isMobileMenuOpen}
        class="flex items-center justify-center xl:hidden w-11 h-11 cursor-pointer"
    >
        <Menu size={24} class="transition-colors duration-300 {showScrolledState ? 'text-foreground' : 'text-white'}" />
    </button>

    <!-- Logo -->
    <div class="xl:justify-start">
        <Logo scrolled={showScrolledState} />
    </div>

    <!-- Desktop Navigation (hidden below 1280px) -->
    <nav class="hidden xl:block">
        <ul class="flex gap-12 items-center list-none p-0 m-0">
            {#each items as item}
                <li>
                    <a
                        class="transition-colors duration-300 {isActive(item.link)
                            ? showScrolledState ? 'text-foreground font-semibold' : 'text-white font-semibold'
                            : showScrolledState ? 'text-muted-foreground hover:text-foreground' : 'text-white/80 hover:text-white'}"
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
        <a href="/book">
            <Button class="rounded-xl py-3 transition-colors duration-300 {showScrolledState
                ? 'bg-foreground text-white hover:bg-foreground-alt'
                : 'bg-white text-foreground hover:bg-white/90'}">
                Réserver maintenant
            </Button>
        </a>
        <Button class="rounded-xl py-3 transition-colors duration-300 {showScrolledState
            ? 'bg-white border border-foreground hover:bg-surface text-foreground'
            : 'bg-transparent border border-white/60 hover:bg-white/10 text-white'}">
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
                                ? 'bg-surface-hover text-foreground font-semibold'
                                : 'text-muted-foreground hover:bg-surface hover:text-foreground'}"
                        >
                            {item.title}
                        </a>
                    </li>
                {/each}
            </ul>
        </nav>

        <!-- Divider -->
        <div class="h-px bg-border-input-hover my-6"></div>

        <!-- Action buttons (full width, stacked) -->
        <div class="flex flex-col gap-3">
            <a href="/book"><Button
                class="w-full inline-flex justify-center items-center px-8 py-3.5 bg-foreground hover:bg-foreground-alt text-background text-sm sm:text-base font-medium transition-all duration-200 shadow-mini hover:shadow-card ring-offset-2 focus:ring-2 focus:ring-foreground"
                >Réserver maintenant</Button></a>
            <Button
                class="w-full bg-background border border-foreground hover:bg-surface px-8 py-3.5 rounded-xl"
            >
                <a href={secondaryButtonLink} class="block w-full">
                    {secondaryButtonText}
                </a>
            </Button>
        </div>
    </div>
</SideDrawer>
