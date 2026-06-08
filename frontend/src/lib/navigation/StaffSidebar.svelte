<script lang="ts">
	import { page } from "$app/state";
	import {
		Home,
		CalendarClock,
		CalendarCheck,
		Package,
		ChartColumn,
		DollarSign,
		Settings,
		ChevronLeft,
		ChevronRight,
		Mail,
		UserRound,
		LogOut,
		ExternalLink,
	} from "@lucide/svelte";
	import type { Component } from "svelte";
	import type { Permissions } from "$lib/security/permissions";
	import ThemeToggle from "$lib/components/admin/ThemeToggle.svelte";

	interface SidebarProps {
		user: App.User;
		permissions: Permissions;
		unreadCount?: number;
	}

	let { user, permissions, unreadCount = 0 }: SidebarProps = $props();

	// Sidebar collapse state
	let isCollapsed = $state(false);
	let userMenuOpen = $state(false);

	function toggleSidebar() {
		isCollapsed = !isCollapsed;
	}

	function toggleUserMenu() {
		userMenuOpen = !userMenuOpen;
	}

	function closeUserMenu() {
		userMenuOpen = false;
	}

	/**
	 * Navigation item definition
	 * - href: Route path
	 * - label: Display text
	 * - icon: Lucide icon component
	 * - adminOnly: Whether only admin users can see this
	 */
	interface NavItem {
		href: string;
		label: string;
		icon: Component;
		adminOnly?: boolean;
	}

	/**
	 * Desktop Navigation Structure
	 * More detailed for larger screens
	 */
	const desktopNavigation: NavItem[] = [
		{
			href: "/staff",
			label: "Accueil",
			icon: Home,
		},
		{
			href: "/staff/agenda/disponibilites",
			label: "Disponibilités",
			icon: CalendarClock,
		},
		{
			href: "/staff/agenda/reservations",
			label: "Réservations",
			icon: CalendarCheck,
		},
		{
			href: "/staff/messages",
			label: "Messages",
			icon: Mail,
		},
		{
			href: "/staff/profile",
			label: "Mon profil",
			icon: UserRound,
		},
		{
			href: "/staff/statistics/analytics",
			label: "Statistiques",
			icon: ChartColumn,
		},
		{
			href: "/staff/statistics/finances",
			label: "Finances",
			icon: DollarSign,
		},
		// Admin only
		{
			href: "/staff/catalog",
			label: "Catalogue",
			icon: Package,
			adminOnly: true,
		},
	];

	/**
	 * Mobile Navigation Structure
	 * Consolidated to 5 items max for bottom bar
	 */
	const mobileNavigation: NavItem[] = [
		{
			href: "/staff",
			label: "Accueil",
			icon: Home,
		},
		{
			href: "/staff/agenda/reservations",
			label: "Réservations",
			icon: CalendarCheck,
		},
		{
			href: "/staff/agenda/disponibilites",
			label: "Agenda",
			icon: CalendarClock,
		},
		{
			href: "/staff/messages",
			label: "Messages",
			icon: Mail,
		},
		{
			href: "/staff/profile",
			label: "Profil",
			icon: UserRound,
		},
	];

	/**
	 * Filter navigation items based on user role
	 */
	const desktopItems = $derived(
		desktopNavigation.filter((item) => !item.adminOnly || user.role === "administrator"),
	);

	const mobileItems = $derived(
		mobileNavigation.filter((item) => !item.adminOnly || user.role === "administrator"),
	);

	/**
	 * Check if a route is active
	 */
	function isActive(href: string): boolean {
		const currentPath = page.url.pathname;
		if (href === "/staff") {
			return currentPath === "/staff";
		}
		return currentPath.startsWith(href);
	}
</script>

{#if permissions.canAccessOps}
	<!-- Desktop Sidebar Navigation -->
	<aside
		class="hidden lg:flex lg:flex-col lg:border-r bg-background border-border-card sticky top-0 h-screen transition-all duration-300 {isCollapsed
			? 'lg:w-20'
			: 'lg:w-64'}"
		aria-label="Sidebar navigation"
	>
		<!-- Collapse Toggle Button -->
		<button
			onclick={toggleSidebar}
			class="absolute -right-4 bottom-24 z-10 w-8 h-8 rounded-full bg-background flex items-center justify-center text-foreground-alt hover:text-foreground hover:bg-muted transition-all duration-200 border border-border-card shadow-sm"
			aria-label={isCollapsed ? "Agrandir la barre latérale" : "Réduire la barre latérale"}
			title={isCollapsed ? "Agrandir la barre latérale" : "Réduire la barre latérale"}
		>
			{#if isCollapsed}
				<ChevronRight size={16} strokeWidth={2} />
			{:else}
				<ChevronLeft size={16} strokeWidth={2} />
			{/if}
		</button>

		<!-- Sidebar Header -->
		<div
			class="flex items-center justify-between border-b border-border-card py-6 {isCollapsed
				? 'px-3'
				: 'px-6'}"
		>
			{#if !isCollapsed}
				<div>
					<h1 class="font-display text-lg font-bold text-foreground tracking-tight">L</h1>
					<p class="text-[9px] text-muted-foreground uppercase tracking-[0.15em] mt-0.5">Leviosa</p>
				</div>
				<div class="flex items-center gap-1">
					<ThemeToggle />
					<a
						href="/staff/settings"
						class="p-2 rounded-md transition-all duration-200 text-foreground-alt hover:bg-muted hover:text-foreground"
						aria-label="Paramètres"
						title="Paramètres"
					>
						<Settings size={18} strokeWidth={1.5} />
					</a>
				</div>
			{:else}
				<div class="mx-auto flex flex-col items-center gap-2">
					<ThemeToggle />
					<a
						href="/staff/settings"
						class="p-2 rounded-md transition-all duration-200 text-foreground-alt hover:bg-muted"
						aria-label="Paramètres"
						title="Paramètres"
					>
						<Settings size={18} strokeWidth={1.5} />
					</a>
				</div>
			{/if}
		</div>

		<!-- Navigation Items -->
		<nav class="flex-1 py-6 overflow-y-auto {isCollapsed ? 'px-3' : 'px-4'}">
			<ul class="space-y-1">
				{#each desktopItems as item (item.href)}
					{@const active = isActive(item.href)}
					{@const badge = item.href === '/staff/messages' ? unreadCount : 0}
					<li>
						<a
							href={item.href}
							class="flex items-center text-sm font-medium transition-all duration-200 rounded-lg
			                       {isCollapsed
								? 'justify-center px-3 py-3'
								: 'gap-3 px-3 py-2.5'}
			                       {active
								? 'text-foreground bg-muted'
								: 'text-foreground-alt hover:text-foreground hover:bg-muted'}"
							aria-current={active ? "page" : undefined}
							title={isCollapsed ? item.label : undefined}
						>
							<item.icon strokeWidth={active ? 2 : 1.5} size={18} />
							{#if !isCollapsed}
								<span class="flex-1">{item.label}</span>
								{#if badge > 0}
									<span class="px-1.5 py-0.5 text-xs font-bold bg-destructive text-white rounded-full">
										{badge > 99 ? '99+' : badge}
									</span>
								{/if}
							{/if}
						</a>
					</li>
				{/each}
			</ul>
		</nav>

		<!-- Sidebar Footer -->
		<div class="py-5 border-t border-border-card {isCollapsed ? 'px-3' : 'px-6'}">
			{#if !isCollapsed}
				<div class="relative">
					{#if userMenuOpen}
						<div class="fixed inset-0 z-10" onclick={closeUserMenu}></div>
						<div
							class="absolute bottom-full left-0 right-0 mb-2 z-20 bg-background border border-border-card rounded-lg shadow-lg py-1 overflow-hidden"
						>
							<a
								href="https://leviosa.care"
								target="_blank"
								rel="noopener noreferrer"
								onclick={closeUserMenu}
								class="flex items-center gap-2.5 px-3 py-2 text-sm text-foreground-alt hover:text-foreground hover:bg-muted transition-colors"
							>
								<ExternalLink size={15} />
								<span>Voir le site</span>
							</a>
							<div class="my-1 border-t border-border-card"></div>
							<form method="POST" action="/logout" class="contents">
								<button
									type="submit"
									class="w-full flex items-center gap-2.5 px-3 py-2 text-sm text-red-600 hover:text-red-700 hover:bg-red-50 dark:hover:bg-red-950 transition-colors cursor-pointer"
								>
									<LogOut size={15} />
									<span>Déconnexion</span>
								</button>
							</form>
						</div>
					{/if}
					<div class="flex items-center gap-3">
						<div
							class="w-9 h-9 flex-shrink-0 rounded-full flex items-center justify-center bg-muted border border-border-card"
						>
							<span class="text-xs font-semibold text-foreground-alt uppercase">
								{user.firstname?.[0] ?? user.email?.[0]?.toUpperCase() ?? "A"}
							</span>
						</div>
						<div class="flex-1 min-w-0">
							<p class="text-sm font-medium text-foreground truncate tracking-tight">
								{user.firstname
									? `${user.firstname} ${user.lastname ?? ""}`.trim()
									: user.email?.split("@")[0] ?? "Staff"}
							</p>
							<p class="text-xs text-muted-foreground capitalize">
								{user.role === "administrator" ? "Administrateur" : "Partenaire"}
							</p>
						</div>
						<button
							onclick={toggleUserMenu}
							class="flex-shrink-0 p-1.5 rounded-md text-foreground-alt hover:text-foreground hover:bg-muted transition-colors"
							aria-label="Options utilisateur"
						>
							<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="currentColor" stroke="none">
								<circle cx="12" cy="5" r="1.5" />
								<circle cx="12" cy="12" r="1.5" />
								<circle cx="12" cy="19" r="1.5" />
							</svg>
						</button>
					</div>
				</div>
			{:else}
				<div class="flex flex-col items-center gap-2">
					<div
						class="w-9 h-9 rounded-full flex items-center justify-center bg-muted border border-border-card"
						title={user.firstname ?? user.email}
					>
						<span class="text-xs font-semibold text-foreground-alt uppercase">
							{user.firstname?.[0] ?? user.email?.[0]?.toUpperCase() ?? "A"}
						</span>
					</div>
					<a
						href="https://leviosa.care"
						target="_blank"
						rel="noopener noreferrer"
						class="flex items-center justify-center w-9 h-9 text-foreground-alt hover:text-foreground hover:bg-muted rounded-lg transition-colors"
						title="Voir le site public"
					>
						<ExternalLink size={14} />
					</a>
				</div>
			{/if}
		</div>
	</aside>

	<!-- Mobile Bottom Navigation -->
	<nav
		class="lg:hidden fixed bottom-0 left-0 right-0 z-50 bg-background/80 backdrop-blur-sm border-t border-border-card"
		aria-label="Bottom navigation"
	>
		<ul class="flex items-center justify-around px-1 py-1">
			{#each mobileItems as item (item.href)}
				{@const active = isActive(item.href)}
				{@const badge = item.href === '/staff/messages' ? unreadCount : 0}
				<li class="flex-1">
					<a
						href={item.href}
						class="flex flex-col items-center gap-1.5 py-3 px-2 rounded-lg transition-all duration-200
		                   {active ? 'text-foreground' : 'text-muted-foreground'}"
						aria-current={active ? "page" : undefined}
					>
						<span class="relative">
							<item.icon size={20} strokeWidth={active ? 2 : 1.5} />
							{#if badge > 0}
								<span class="absolute -top-1 -right-1 w-4 h-4 bg-destructive text-white text-xs rounded-full flex items-center justify-center font-bold leading-none">
									{badge > 9 ? '9+' : badge}
								</span>
							{/if}
						</span>
						<span class="text-xs font-medium">
							{item.label}
						</span>
					</a>
				</li>
			{/each}
		</ul>
	</nav>
{/if}
