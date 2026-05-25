<script lang="ts">
	import { page } from '$app/state';

	import { Home, CalendarDays, PlusCircle, Menu, X } from '@lucide/svelte';
	import type { Component } from 'svelte';

	interface Props {
		user: App.User;
	}

	let { user }: Props = $props();
	let isMobileMenuOpen = $state(false);

	interface NavItem {
		href: string;
		label: string;
		icon: Component;
	}

	const navItems: NavItem[] = [
		{ href: '/client', label: 'Accueil', icon: Home },
		{ href: '/client/bookings', label: 'Mes réservations', icon: CalendarDays },
		{ href: '/book', label: 'Réserver', icon: PlusCircle },
	];

	function isActive(href: string): boolean {
		if (href === '/client') return page.url.pathname === '/client';
		return page.url.pathname.startsWith(href);
	}

	const initials = ((user.firstname?.[0] ?? '') + (user.lastname?.[0] ?? '')).toUpperCase();
	const displayName = `${user.firstname ?? ''} ${user.lastname ?? ''}`.trim() || 'Client';
</script>

<svelte:window onkeydown={(e) => { if (e.key === 'Escape') isMobileMenuOpen = false; }} />

<!-- Top navigation bar -->
<header class="sticky top-0 z-40 bg-white/80 backdrop-blur-md border-b border-dark-10">
	<div class="max-w-5xl mx-auto flex items-center justify-between px-4 h-16">
		<!-- Mobile hamburger -->
		<button
			onclick={() => (isMobileMenuOpen = true)}
			class="lg:hidden flex items-center justify-center w-10 h-10"
			aria-label="Ouvrir le menu"
		>
			<Menu size={22} />
		</button>

		<!-- Logo / brand -->
		<a href="/client" class="font-semibold text-lg tracking-tight text-foreground">
			Espace client
		</a>

		<!-- Desktop nav -->
		<nav class="hidden lg:flex items-center gap-6">
			{#each navItems as item (item.href)}
				<a
					href={item.href}
					class="flex items-center gap-1.5 text-sm font-medium transition-colors {isActive(item.href)
						? 'text-foreground'
						: 'text-muted-foreground hover:text-foreground'}"
				>
					<item.icon size={16} strokeWidth={isActive(item.href) ? 2 : 1.5} />
					{item.label}
				</a>
			{/each}
		</nav>

		<!-- User avatar -->
		<div class="flex items-center gap-3">
			<div class="hidden sm:block text-right">
				<p class="text-sm font-medium text-foreground leading-tight">{displayName}</p>
				<p class="text-xs text-muted-foreground">Client</p>
			</div>
			<div class="w-9 h-9 rounded-full bg-muted flex items-center justify-center border border-border-card">
				{#if user.picture}
					<img src={user.picture} alt={displayName} class="w-9 h-9 rounded-full object-cover" />
				{:else}
					<span class="text-xs font-semibold text-foreground">{initials || '?'}</span>
				{/if}
			</div>
		</div>
	</div>
</header>

<!-- Mobile drawer -->
{#if isMobileMenuOpen}
	<div class="fixed inset-0 z-50 lg:hidden">
		<!-- Backdrop -->
		<div class="absolute inset-0 bg-black/40" onclick={() => (isMobileMenuOpen = false)}></div>

		<!-- Panel -->
		<div class="absolute left-0 top-0 bottom-0 w-72 bg-background shadow-xl flex flex-col">
			<div class="flex items-center justify-between p-4 border-b border-border-card">
				<p class="font-semibold text-foreground">Menu</p>
				<button onclick={() => (isMobileMenuOpen = false)} class="w-10 h-10 flex items-center justify-center">
					<X size={20} />
				</button>
			</div>

			<nav class="flex-1 p-4">
				<ul class="space-y-1">
					{#each navItems as item (item.href)}
						<li>
							<a
								href={item.href}
								onclick={() => (isMobileMenuOpen = false)}
								class="flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm font-medium transition-colors {isActive(item.href)
									? 'bg-muted text-foreground'
									: 'text-muted-foreground hover:bg-muted/50 hover:text-foreground'}"
							>
								<item.icon size={18} strokeWidth={isActive(item.href) ? 2 : 1.5} />
								{item.label}
							</a>
						</li>
					{/each}
				</ul>
			</nav>

			<div class="p-4 border-t border-border-card">
				<div class="flex items-center gap-3">
					<div class="w-9 h-9 rounded-full bg-muted flex items-center justify-center">
						{#if user.picture}
							<img src={user.picture} alt={displayName} class="w-9 h-9 rounded-full object-cover" />
						{:else}
							<span class="text-xs font-semibold text-foreground">{initials || '?'}</span>
						{/if}
					</div>
					<div>
						<p class="text-sm font-medium text-foreground">{displayName}</p>
						<p class="text-xs text-muted-foreground">{user.email}</p>
					</div>
				</div>
			</div>
		</div>
	</div>
{/if}
