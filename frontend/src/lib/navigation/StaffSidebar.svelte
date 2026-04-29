<script lang="ts">
	import { page } from "$app/state";
	import {
		Home,
		CalendarClock,
		CalendarCheck,
		Package,
		ChartColumn,
		DollarSign,
		Users,
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

	interface SidebarProps {
		user: App.User;
		permissions: Permissions;
	}

	let { user, permissions }: SidebarProps = $props();

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
		{
			href: "/staff/users",
			label: "Utilisateurs",
			icon: Users,
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
		class="hidden lg:flex lg:flex-col lg:border-r bg-white border-dark-100 relative transition-all duration-300 min-h-screen
		           {isCollapsed ? 'lg:w-20' : 'lg:w-64'}"
		aria-label="Sidebar navigation"
	>
		<!-- Collapse Toggle Button -->
		<button
			onclick={toggleSidebar}
			class="absolute -right-4 bottom-24 z-10 w-8 h-8 rounded-full bg-white flex items-center justify-center
		           text-dark-700 hover:text-dark-900 hover:bg-dark-100 transition-all duration-200 border border-dark-200"
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
			class="flex items-center justify-between border-b border-dark-100 py-6 {isCollapsed
				? 'px-3'
				: 'px-6'}"
		>
			{#if !isCollapsed}
				<div>
					<h1 class="text-sm font-semibold tracking-tight text-dark-900 uppercase">
						{user.role === "administrator" ? "Admin" : "Staff"}
					</h1>
					<p class="text-xs text-dark-500">Leviosa</p>
				</div>
				<a
					href="/staff/settings"
					class="p-2 rounded-md transition-all duration-200 text-dark-500 hover:bg-dark-100 hover:text-dark-900"
					aria-label="Paramètres"
					title="Paramètres"
				>
					<Settings size={18} strokeWidth={1.5} />
				</a>
			{:else}
				<a
					href="/staff/settings"
					class="p-2 rounded-md transition-all duration-200 text-dark-500 hover:bg-dark-100 mx-auto"
					aria-label="Paramètres"
					title="Paramètres"
				>
					<Settings size={18} strokeWidth={1.5} />
				</a>
			{/if}
		</div>

		<!-- Navigation Items -->
		<nav class="flex-1 py-6 {isCollapsed ? 'px-3' : 'px-4'}">
			<ul class="space-y-1">
				{#each desktopItems as item (item.href)}
					{@const active = isActive(item.href)}
					<li>
						<a
							href={item.href}
							class="flex items-center text-sm font-medium transition-all duration-200 rounded-md
			                       {isCollapsed
								? 'justify-center px-3 py-3'
								: 'gap-3 px-3 py-2.5'}
			                       {active
								? 'text-dark-900 bg-dark-100'
								: 'text-dark-600 hover:text-dark-900 hover:bg-dark-50'}"
							aria-current={active ? "page" : undefined}
							title={isCollapsed ? item.label : undefined}
						>
							<item.icon strokeWidth={active ? 2 : 1.5} size={18} />
							{#if !isCollapsed}
								<span>{item.label}</span>
							{/if}
						</a>
					</li>
				{/each}
			</ul>
		</nav>

		<!-- Sidebar Footer -->
		<div class="py-5 border-t border-dark-100 {isCollapsed ? 'px-3' : 'px-6'}">
			{#if !isCollapsed}
				<div class="relative">
					{#if userMenuOpen}
						<div class="fixed inset-0 z-10" onclick={closeUserMenu}></div>
						<div
							class="absolute bottom-full left-0 right-0 mb-2 z-20 bg-white border border-dark-100 rounded-lg shadow-lg py-1 overflow-hidden"
						>
							<a
								href="/"
								onclick={closeUserMenu}
								class="flex items-center gap-2.5 px-3 py-2 text-sm text-dark-600 hover:text-dark-900 hover:bg-dark-50 transition-colors"
							>
								<ExternalLink size={15} />
								<span>Voir le site</span>
							</a>
							<div class="my-1 border-t border-dark-100"></div>
							<form method="POST" action="/staff/logout?/logout" class="contents">
								<button
									type="submit"
									onclick={closeUserMenu}
									class="w-full flex items-center gap-2.5 px-3 py-2 text-sm text-red-600 hover:text-red-700 hover:bg-red-50 transition-colors cursor-pointer"
								>
									<LogOut size={15} />
									<span>Déconnexion</span>
								</button>
							</form>
						</div>
					{/if}
					<div class="flex items-center gap-3">
						<div class="w-9 h-9 rounded-full flex items-center justify-center bg-dark-100">
							<span class="text-xs font-semibold text-dark-700 uppercase">
								{user.firstname?.[0] ?? user.email?.[0]?.toUpperCase() ?? "A"}
							</span>
						</div>
						<div class="flex-1 min-w-0">
							<p class="text-sm font-medium text-dark-900 truncate">
								{user.firstname
									? `${user.firstname} ${user.lastname ?? ""}`.trim()
									: user.email?.split("@")[0] ?? "Staff"}
							</p>
							<p class="text-xs text-dark-500">
								{user.role === "administrator" ? "Administrateur" : "Partenaire"}
							</p>
						</div>
						<button
							onclick={toggleUserMenu}
							class="flex-shrink-0 p-1.5 rounded-md text-dark-400 hover:text-dark-900 hover:bg-dark-100 transition-colors"
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
				<div
					class="w-9 h-9 mx-auto rounded-full flex items-center justify-center bg-dark-100"
					title={user.firstname ?? user.email}
				>
					<span class="text-xs font-semibold text-dark-700 uppercase">
						{user.firstname?.[0] ?? user.email?.[0]?.toUpperCase() ?? "A"}
					</span>
				</div>
			{/if}
		</div>
	</aside>

	<!-- Mobile Bottom Navigation -->
	<nav
		class="lg:hidden fixed bottom-0 left-0 right-0 z-50 bg-white backdrop-blur-sm border-t border-dark-100"
		aria-label="Bottom navigation"
	>
		<ul class="flex items-center justify-around px-1 py-1">
			{#each mobileItems as item (item.href)}
				{@const active = isActive(item.href)}
				<li class="flex-1">
					<a
						href={item.href}
						class="flex flex-col items-center gap-1.5 py-3 px-2 rounded-lg transition-all duration-200
		                   {active ? 'text-dark-900' : 'text-dark-500'}"
						aria-current={active ? "page" : undefined}
					>
						<item.icon size={20} strokeWidth={active ? 2 : 1.5} />
						<span class="text-xs font-medium">
							{item.label}
						</span>
					</a>
				</li>
			{/each}
		</ul>
	</nav>
{/if}
