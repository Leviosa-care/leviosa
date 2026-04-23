<script lang="ts">
	import { page } from "$app/state";
	import {
		Home,
		Users,
		CalendarDays,
		NotebookPen,
		Mail,
		Ticket,
		Package,
		ChartSpline,
		Server,
		HandCoins,
		ExternalLink,
	} from "@lucide/svelte";
	import type { Component } from "svelte";
	import type { Permissions } from "$lib/security/permissions";
	import ThemeToggle from "$lib/components/admin/ThemeToggle.svelte";

	interface SidebarProps {
		user: App.User;
		permissions: Permissions;
	}

	let { user, permissions }: SidebarProps = $props();

	// Sidebar collapse state
	let isCollapsed = $state(false);

	// User menu state
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
	 */
	interface NavItem {
		href: string;
		label: string;
		icon: Component;
	}

	/**
	 * Desktop Navigation Structure
	 * Using the original links from the old NavigationBar
	 */
	const desktopNavigation: NavItem[] = [
		{
			href: "/admin",
			label: "Accueil",
			icon: Home,
		},
		{
			href: "/admin/users",
			label: "Utilisateurs",
			icon: Users,
		},
		{
			href: "/admin/planning",
			label: "Planning",
			icon: CalendarDays,
		},
		{
			href: "/admin/bookings/consultations",
			label: "Notes de seance",
			icon: NotebookPen,
		},
		{
			href: "/admin/messages",
			label: "Messages",
			icon: Mail,
		},
		{
			href: "/admin/bookings/events",
			label: "Evenements",
			icon: Ticket,
		},
		{
			href: "/admin/products",
			label: "Catalogue",
			icon: Package,
		},
		{
			href: "/admin/analytics",
			label: "Analytics",
			icon: ChartSpline,
		},
		{
			href: "/admin/infra",
			label: "Infrastructure",
			icon: Server,
		},
		{
			href: "/admin/compta",
			label: "Comptabilite",
			icon: HandCoins,
		},
	];

	/**
	 * Mobile Navigation Structure
	 * Limited to 5 items for bottom bar
	 */
	const mobileNavigation: NavItem[] = [
		{
			href: "/admin",
			label: "Accueil",
			icon: Home,
		},
		{
			href: "/admin/users",
			label: "Utilisateurs",
			icon: Users,
		},
		{
			href: "/admin/planning",
			label: "Planning",
			icon: CalendarDays,
		},
		{
			href: "/admin/bookings/consultations",
			label: "Notes",
			icon: NotebookPen,
		},
		{
			href: "/admin/products",
			label: "Catalogue",
			icon: Package,
		},
	];

	/**
	 * Check if a route is active
	 */
	function isActive(href: string): boolean {
		const currentPath = page.url.pathname;
		if (href === "/admin") {
			return currentPath === "/admin" || currentPath === "/admin/";
		}
		return currentPath.startsWith(href);
	}
</script>

{#if permissions.canAccessOps}
	<!-- Desktop Sidebar Navigation -->
	<aside
		class="hidden lg:flex lg:flex-col lg:border-r bg-white dark:bg-dark-50 border-dark-100 dark:border-dark-300 sticky top-0 h-screen transition-all duration-300 {isCollapsed
			? 'lg:w-20'
			: 'lg:w-64'}"
		aria-label="Sidebar navigation"
	>
		<!-- Collapse Toggle Button -->
		<button
			onclick={toggleSidebar}
			class="absolute -right-4 bottom-24 z-10 w-8 h-8 rounded-full bg-white dark:bg-dark-100 flex items-center justify-center text-dark-700 hover:text-dark-900 dark:text-dark-800 dark:hover:text-dark-900 hover:bg-dark-100 dark:hover:bg-dark-200 transition-all duration-200 border border-dark-200 dark:border-dark-300 shadow-sm"
			aria-label={isCollapsed ? "Agrandir la barre laterale" : "Reduire la barre laterale"}
			title={isCollapsed ? "Agrandir la barre laterale" : "Reduire la barre laterale"}
		>
			{#if isCollapsed}
				<svg
					xmlns="http://www.w3.org/2000/svg"
					width="16"
					height="16"
					viewBox="0 0 24 24"
					fill="none"
					stroke="currentColor"
					stroke-width="2"
					stroke-linecap="round"
					stroke-linejoin="round"
				>
					<polyline points="9 18 15 12 9 6"></polyline>
				</svg>
			{:else}
				<svg
					xmlns="http://www.w3.org/2000/svg"
					width="16"
					height="16"
					viewBox="0 0 24 24"
					fill="none"
					stroke="currentColor"
					stroke-width="2"
					stroke-linecap="round"
					stroke-linejoin="round"
				>
					<polyline points="15 18 9 12 15 6"></polyline>
				</svg>
			{/if}
		</button>

		<!-- Sidebar Header -->
		<div
			class="flex items-center justify-between border-b border-dark-100 dark:border-dark-300 py-6 {isCollapsed
				? 'px-3'
				: 'px-6'}"
		>
			{#if !isCollapsed}
				<div>
					<h1 class="text-sm font-semibold tracking-tight text-dark-900 dark:text-dark-100 uppercase">Leviosa</h1>
					<p class="text-xs text-dark-500 dark:text-dark-400">Panneau d'Administration</p>
				</div>
				<ThemeToggle />
			{:else}
				<div class="mx-auto flex flex-col items-center gap-2">
					<span class="text-lg font-bold text-dark-900 dark:text-dark-100">L</span>
					<ThemeToggle />
				</div>
			{/if}
		</div>

		<!-- Navigation Items -->
		<nav class="flex-1 py-6 overflow-y-auto {isCollapsed ? 'px-3' : 'px-4'}">
			<ul class="space-y-1">
				{#each desktopNavigation as item (item.href)}
					{@const active = isActive(item.href)}
					<li>
						<a
							href={item.href}
							class="flex items-center text-sm font-medium transition-all duration-200 rounded-lg {isCollapsed
								? 'justify-center px-3 py-3'
								: 'gap-3 px-3 py-2.5'} {active
								? 'text-dark-900 dark:text-dark-900 bg-dark-100 dark:bg-dark-100'
								: 'text-dark-600 dark:text-dark-500 hover:text-dark-900 dark:hover:text-dark-800 hover:bg-dark-50 dark:hover:bg-dark-100'}"
							aria-current={active ? "page" : undefined}
							title={isCollapsed ? item.label : undefined}
						>
							<item.icon strokeWidth={active ? 2 : 1.5} size={20} />
							{#if !isCollapsed}
								<span class="tracking-tight">{item.label}</span>
							{/if}
						</a>
					</li>
				{/each}
			</ul>
		</nav>

		<!-- Sidebar Footer -->
		<div class="py-5 border-t border-dark-100 dark:border-dark-300 {isCollapsed ? 'px-3' : 'px-6'}">
			{#if !isCollapsed}
				<!-- User row with kebab menu -->
				<div class="relative">
					{#if userMenuOpen}
						<!-- Backdrop -->
						<div class="fixed inset-0 z-10" onclick={closeUserMenu}></div>
						<!-- Dropdown (opens upward) -->
						<div
							class="absolute bottom-full left-0 right-0 mb-2 z-20 bg-white dark:bg-dark-100 border border-dark-200 dark:border-dark-300 rounded-lg shadow-lg py-1 overflow-hidden"
						>
							<a
								href="https://leviosa.com"
								target="_blank"
								rel="noopener noreferrer"
								onclick={closeUserMenu}
								class="flex items-center gap-2.5 px-3 py-2 text-sm text-dark-600 dark:text-dark-500 hover:text-dark-900 dark:hover:text-dark-800 hover:bg-dark-50 dark:hover:bg-dark-200 transition-colors"
							>
								<ExternalLink size={15} />
								<span>Voir le site</span>
							</a>
							<div class="my-1 border-t border-dark-100 dark:border-dark-300"></div>
							<form method="POST" action="/logout">
								<button
									type="submit"
									class="w-full flex items-center gap-2.5 px-3 py-2 text-sm text-red-600 hover:text-red-700 hover:bg-red-50 transition-colors"
								>
									<svg
										xmlns="http://www.w3.org/2000/svg"
										width="15"
										height="15"
										viewBox="0 0 24 24"
										fill="none"
										stroke="currentColor"
										stroke-width="2"
										stroke-linecap="round"
										stroke-linejoin="round"
									>
										<path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"></path>
										<polyline points="16 17 21 12 16 7"></polyline>
										<line x1="21" x2="9" y1="12" y2="12"></line>
									</svg>
									<span>Deconnexion</span>
								</button>
							</form>
						</div>
					{/if}
					<div class="flex items-center gap-3">
						<div
							class="w-9 h-9 flex-shrink-0 rounded-full flex items-center justify-center bg-dark-100 dark:bg-dark-200 border border-dark-200 dark:border-dark-300"
						>
							<span class="text-xs font-semibold text-dark-700 uppercase">
								{user.email?.[0]?.toUpperCase() ?? "A"}
							</span>
						</div>
						<div class="flex-1 min-w-0">
							<p class="text-sm font-medium text-dark-900 dark:text-dark-100 truncate tracking-tight">
								{user.email}
							</p>
							<p class="text-xs text-dark-500 dark:text-dark-400 capitalize">{user.role}</p>
						</div>
						<button
							onclick={toggleUserMenu}
							class="flex-shrink-0 p-1.5 rounded-md text-dark-600 dark:text-dark-500 hover:text-dark-900 dark:hover:text-dark-800 hover:bg-dark-100 dark:hover:bg-dark-200 transition-colors"
							aria-label="Options utilisateur"
							title="Options"
						>
							<svg
								xmlns="http://www.w3.org/2000/svg"
								width="16"
								height="16"
								viewBox="0 0 24 24"
								fill="currentColor"
								stroke="none"
							>
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
						class="w-9 h-9 rounded-full flex items-center justify-center bg-dark-100 dark:bg-dark-200 border border-dark-200 dark:border-dark-300"
						title={user.email}
					>
						<span class="text-xs font-semibold text-dark-700 dark:text-dark-900 uppercase">
							{user.email?.[0]?.toUpperCase() ?? "A"}
						</span>
					</div>
					<a
						href="https://leviosa.com"
						target="_blank"
						rel="noopener noreferrer"
						class="flex items-center justify-center w-9 h-9 text-dark-600 dark:text-dark-500 hover:text-dark-900 dark:hover:text-dark-800 hover:bg-dark-100 dark:hover:bg-dark-200 rounded-lg transition-colors"
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
		class="lg:hidden fixed bottom-0 left-0 right-0 z-50 bg-white dark:bg-dark-50 backdrop-blur-sm border-t border-dark-100 dark:border-dark-300"
		aria-label="Bottom navigation"
	>
		<ul class="flex items-center justify-around px-1 py-1">
			{#each mobileNavigation as item (item.href)}
				{@const active = isActive(item.href)}
				<li class="flex-1">
					<a
						href={item.href}
						class="flex flex-col items-center gap-1.5 py-3 px-2 rounded-lg transition-all duration-200 {active
							? 'text-dark-900 dark:text-dark-800'
							: 'text-dark-500 dark:text-dark-500'}"
						aria-current={active ? "page" : undefined}
					>
						<item.icon size={20} strokeWidth={active ? 2 : 1.5} />
						<span class="text-xs font-medium">{item.label}</span>
					</a>
				</li>
			{/each}
		</ul>
	</nav>
{/if}
