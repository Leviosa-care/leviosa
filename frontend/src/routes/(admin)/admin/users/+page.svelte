<script lang="ts">
	import type { PageData, ActionData } from "./$types";
	import { enhance } from "$app/forms";
	import Tabs from "$lib/ui/bits-components/Tabs.svelte";
	import TabsList from "$lib/ui/bits-components/TabsList.svelte";
	import TabsTrigger from "$lib/ui/bits-components/TabsTrigger.svelte";
	import TabsContent from "$lib/ui/bits-components/TabsContent.svelte";
	import Modal from "$lib/ui/Modal.svelte";
	import Drawer from "$lib/ui/Drawer.svelte";
	import { cn } from "$lib/utils/design-system";
	import { ROLES } from "$lib/types/role";
	import {
		Shield,
		ShieldCheck,
		Clock,
		MoreVertical,
		UserCheck,
		Trash2,
		KeyRound,
		X,
		CheckCircle,
		AlertCircle,
		BadgeCheck,
		Building2,
		FileText,
		CreditCard
	} from "@lucide/svelte";

	let { data, form }: { data: PageData; form: ActionData } = $props();

	// Mobile detection
	let isMobile = $state(false);

	$effect(() => {
		isMobile = window.innerWidth < 768;
		const onResize = () => { isMobile = window.innerWidth < 768; };
		window.addEventListener("resize", onResize);
		return () => window.removeEventListener("resize", onResize);
	});

	// Dialog states
	let deleteDialogOpen = $state(false);
	let editDialogOpen = $state(false);
	let approveDialogOpen = $state(false);
	let verifyPartnerDialogOpen = $state(false);
	let deletePartnerDialogOpen = $state(false);
	let selectedUser = $state<typeof data.users[0] | null>(null);
	let selectedPartner = $state<typeof data.partners[0] | null>(null);
	let selectedRole = $state<string>("");

	interface Trigger {
		value: string;
		name: string;
		count: number;
	}

	let triggers = $derived<Trigger[]>([
		{ value: "all", name: "Tous les utilisateurs", count: data.users.length },
		{ value: "pending", name: "En attente", count: data.pendingUsers.length },
		{ value: "partners", name: "Partenaires", count: data.partners.length }
	]);

	function getRoleBadgeClass(role: string) {
		return cn(
			"inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-xs font-medium",
			role === "administrator"
				? "bg-purple-100 text-purple-700 dark:bg-purple-900/30 dark:text-purple-300"
				: role === "partner"
					? "bg-blue-100 text-blue-700 dark:bg-blue-900/30 dark:text-blue-300"
					: role === "premium"
						? "bg-amber-100 text-amber-700 dark:bg-amber-900/30 dark:text-amber-300"
						: "bg-gray-100 text-gray-700 dark:bg-gray-800 dark:text-gray-300"
		);
	}

	function getStatusBadgeClass(status: string) {
		return cn(
			"inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-xs font-medium",
			status === "approved"
				? "bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-300"
				: "bg-amber-100 text-amber-700 dark:bg-amber-900/30 dark:text-amber-300"
		);
	}

	function getVerificationBadgeClass(isVerified: boolean) {
		return cn(
			"inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-xs font-medium",
			isVerified
				? "bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-300"
				: "bg-amber-100 text-amber-700 dark:bg-amber-900/30 dark:text-amber-300"
		);
	}

	function getStripeBadgeClass(onboardingComplete: boolean, status: string) {
		return cn(
			"inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-xs font-medium",
			onboardingComplete && status === "active"
				? "bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-300"
				: status === "pending"
					? "bg-gray-100 text-gray-700 dark:bg-gray-800 dark:text-gray-300"
					: "bg-amber-100 text-amber-700 dark:bg-amber-900/30 dark:text-amber-300"
		);
	}

	function formatDate(dateString: string) {
		const date = new Date(dateString);
		return date.toLocaleDateString("fr-FR", {
			day: "numeric",
			month: "short",
			year: "numeric"
		});
	}

	function truncate(text: string, maxLength: number = 80) {
		if (text.length <= maxLength) return text;
		return text.substring(0, maxLength) + "…";
	}

	/** Look up user info for a partner by cross-referencing userId with the users list. */
	function getUserForPartner(partner: typeof data.partners[0]) {
		return data.users.find((u) => u.id === partner.userId) ?? null;
	}

	function openDeleteDialog(user: typeof data.users[0]) {
		selectedUser = user;
		deleteDialogOpen = true;
	}

	function openEditDialog(user: typeof data.users[0]) {
		selectedUser = user;
		selectedRole = user.role;
		editDialogOpen = true;
	}

	function openApproveDialog(user: typeof data.users[0]) {
		selectedUser = user;
		approveDialogOpen = true;
	}

	function openVerifyPartnerDialog(partner: typeof data.partners[0]) {
		selectedPartner = partner;
		verifyPartnerDialogOpen = true;
	}

	function openDeletePartnerDialog(partner: typeof data.partners[0]) {
		selectedPartner = partner;
		deletePartnerDialogOpen = true;
	}

	const availableRoles = Object.entries(ROLES).map(([key, value]) => ({
		label: key.charAt(0).toUpperCase() + key.slice(1),
		value: value
	}));

	// Close dialogs on success
	$effect(() => {
		if (form?.success) {
			deleteDialogOpen = false;
			editDialogOpen = false;
			approveDialogOpen = false;
			verifyPartnerDialogOpen = false;
			deletePartnerDialogOpen = false;
		}
	});
</script>

<svelte:head>
	<title>Utilisateurs | Admin</title>
</svelte:head>

<div class="container mx-auto px-4 py-8 lg:py-12">
	<div class="mb-8">
		<h1 class="text-3xl lg:text-4xl font-bold mb-2 text-foreground">Utilisateurs</h1>
		<p class="text-muted-foreground">Gérer les utilisateurs et leurs permissions</p>
	</div>

	{#if form?.error}
		<div class="mb-6 flex items-center gap-3 px-4 py-3 bg-red-50 border border-red-200 text-red-800 rounded-lg">
			<AlertCircle size={20} />
			<p class="text-sm font-medium">{form.error}</p>
		</div>
	{/if}

	{#if form?.success}
		<div class="mb-6 flex items-center gap-3 px-4 py-3 bg-green-50 border border-green-200 text-green-800 rounded-lg">
			<CheckCircle size={20} />
			<p class="text-sm font-medium">{form.success}</p>
		</div>
	{/if}

	<Tabs value="all" class="space-y-4">
		<TabsList class="inline-flex items-center w-fit bg-transparent gap-2 text-sm font-semibold border-b border-border-card p-1">
			{#each triggers as trigger}
				<TabsTrigger
					value={trigger.value}
					class="px-4 py-2 rounded-none bg-transparent border-b-2 data-[state=active]:shadow-none mb-[-2px] data-[state=active]:border-b-foreground data-[state=active]:text-foreground data-[state=inactive]:border-transparent data-[state=inactive]:text-muted-foreground data-[state=inactive]:hover:bg-transparent data-[state=inactive]:hover:text-foreground-alt transition-colors cursor-pointer flex items-center gap-2"
				>
					{trigger.name}
					<span class="text-xs bg-muted px-2 py-0.5 rounded-full">{trigger.count}</span>
				</TabsTrigger>
			{/each}
		</TabsList>

		<!-- All Users Tab -->
		<TabsContent value="all" class="p-6">
			<div class="bg-background border border-border-card rounded-lg overflow-hidden">
				<div class="overflow-x-auto">
					<table class="w-full">
						<thead>
							<tr class="border-b border-border-card bg-muted/50">
								<th class="text-left px-4 py-3 text-sm font-semibold text-foreground">Utilisateur</th>
								<th class="text-left px-4 py-3 text-sm font-semibold text-foreground">Rôle</th>
								<th class="text-left px-4 py-3 text-sm font-semibold text-foreground">Statut</th>
								<th class="text-left px-4 py-3 text-sm font-semibold text-foreground">Créé le</th>
								<th class="text-right px-4 py-3 text-sm font-semibold text-foreground">Actions</th>
							</tr>
						</thead>
						<tbody>
							{#if data.users.length === 0}
								<tr>
									<td colspan="5" class="px-4 py-8 text-center text-muted-foreground">
										Aucun utilisateur trouvé
									</td>
								</tr>
							{:else}
								{#each data.users as user (user.id)}
									<tr class="border-b border-border-card hover:bg-muted/30 transition-colors">
										<td class="px-4 py-3">
											<div class="flex items-center gap-3">
												<div class="w-9 h-9 rounded-full bg-muted flex items-center justify-center">
													<span class="text-xs font-semibold text-foreground-alt uppercase">
														{user.email[0].toUpperCase()}
													</span>
												</div>
												<div>
													<p class="text-sm font-medium text-foreground">{user.email}</p>
													{#if user.firstname || user.lastname}
														<p class="text-xs text-muted-foreground">
															{user.firstname} {user.lastname}
														</p>
													{/if}
												</div>
											</div>
										</td>
										<td class="px-4 py-3">
											<span class={getRoleBadgeClass(user.role)}>
												{#if user.role === "administrator"}
													<Shield size={12} />
												{:else if user.role === "partner"}
													<ShieldCheck size={12} />
												{:else}
													<KeyRound size={12} />
												{/if}
												{user.role}
											</span>
										</td>
										<td class="px-4 py-3">
											<span class={getStatusBadgeClass(user.status)}>
												{#if user.status === "approved"}
													<UserCheck size={12} />
												{:else}
													<Clock size={12} />
												{/if}
												{user.status === "approved" ? "Approuvé" : "En attente"}
											</span>
										</td>
										<td class="px-4 py-3 text-sm text-muted-foreground">
											{formatDate(user.createdAt)}
										</td>
										<td class="px-4 py-3">
											<div class="flex items-center justify-end gap-1">
												<button
													onclick={() => openEditDialog(user)}
													class="p-2 rounded-md text-foreground-alt hover:text-foreground hover:bg-muted transition-colors"
													title="Modifier le rôle"
												>
													<Shield size={16} />
												</button>
												<button
													onclick={() => openDeleteDialog(user)}
													class="p-2 rounded-md text-red-600 hover:text-red-700 hover:bg-red-50 dark:hover:bg-red-900/20 transition-colors"
													title="Supprimer"
												>
													<Trash2 size={16} />
												</button>
											</div>
										</td>
									</tr>
								{/each}
							{/if}
						</tbody>
					</table>
				</div>
			</div>
		</TabsContent>

		<!-- Pending Users Tab -->
		<TabsContent value="pending" class="p-6">
			<div class="bg-background border border-border-card rounded-lg overflow-hidden">
				<div class="overflow-x-auto">
					<table class="w-full">
						<thead>
							<tr class="border-b border-border-card bg-muted/50">
								<th class="text-left px-4 py-3 text-sm font-semibold text-foreground">Utilisateur</th>
								<th class="text-left px-4 py-3 text-sm font-semibold text-foreground">Téléphone</th>
								<th class="text-left px-4 py-3 text-sm font-semibold text-foreground">Demandé le</th>
								<th class="text-right px-4 py-3 text-sm font-semibold text-foreground">Actions</th>
							</tr>
						</thead>
						<tbody>
							{#if data.pendingUsers.length === 0}
								<tr>
									<td colspan="4" class="px-4 py-8 text-center text-muted-foreground">
										Aucun utilisateur en attente
									</td>
								</tr>
							{:else}
								{#each data.pendingUsers as user (user.id)}
									<tr class="border-b border-border-card hover:bg-muted/30 transition-colors">
										<td class="px-4 py-3">
											<div class="flex items-center gap-3">
												<div class="w-9 h-9 rounded-full bg-muted flex items-center justify-center">
													<span class="text-xs font-semibold text-foreground-alt uppercase">
														{user.email[0].toUpperCase()}
													</span>
												</div>
												<div>
													<p class="text-sm font-medium text-foreground">{user.email}</p>
													{#if user.firstname || user.lastname}
														<p class="text-xs text-muted-foreground">
															{user.firstname} {user.lastname}
														</p>
													{/if}
												</div>
											</div>
										</td>
										<td class="px-4 py-3 text-sm text-muted-foreground">
											{user.telephone || "-"}
										</td>
										<td class="px-4 py-3 text-sm text-muted-foreground">
											{formatDate(user.createdAt)}
										</td>
										<td class="px-4 py-3">
											<div class="flex items-center justify-end gap-1">
												<button
													onclick={() => openApproveDialog(user)}
													class="inline-flex items-center gap-1.5 px-3 py-1.5 text-sm font-medium text-green-700 bg-green-100 hover:bg-green-200 dark:bg-green-900/30 dark:text-green-300 dark:hover:bg-green-900/50 rounded-md transition-colors"
													title="Approuver"
												>
													<UserCheck size={14} />
													Approuver
												</button>
												<button
													onclick={() => openDeleteDialog(user)}
													class="p-2 rounded-md text-red-600 hover:text-red-700 hover:bg-red-50 dark:hover:bg-red-900/20 transition-colors"
													title="Supprimer"
												>
													<Trash2 size={16} />
												</button>
											</div>
										</td>
									</tr>
								{/each}
							{/if}
						</tbody>
					</table>
				</div>
			</div>
		</TabsContent>

		<!-- Partners Tab -->
		<TabsContent value="partners" class="p-6">
			{#if data.partnersError}
				<div class="mb-4 flex items-center gap-3 px-4 py-3 bg-red-50 border border-red-200 text-red-800 rounded-lg">
					<AlertCircle size={20} />
					<p class="text-sm font-medium">Impossible de charger les partenaires. Veuillez rafraîchir la page.</p>
				</div>
			{/if}
			<div class="bg-background border border-border-card rounded-lg overflow-hidden">
				<div class="overflow-x-auto">
					<table class="w-full">
						<thead>
							<tr class="border-b border-border-card bg-muted/50">
								<th class="text-left px-4 py-3 text-sm font-semibold text-foreground">Partenaire</th>
								<th class="text-left px-4 py-3 text-sm font-semibold text-foreground">Bio</th>
								<th class="text-left px-4 py-3 text-sm font-semibold text-foreground">Vérifié</th>
								<th class="text-left px-4 py-3 text-sm font-semibold text-foreground">Stripe</th>
								<th class="text-left px-4 py-3 text-sm font-semibold text-foreground">Catégories / Produits</th>
								<th class="text-left px-4 py-3 text-sm font-semibold text-foreground">Créé le</th>
								<th class="text-right px-4 py-3 text-sm font-semibold text-foreground">Actions</th>
							</tr>
						</thead>
						<tbody>
							{#if data.partners.length === 0}
								<tr>
									<td colspan="7" class="px-4 py-8 text-center text-muted-foreground">
										{data.partnersError ? "Erreur de chargement" : "Aucun partenaire trouvé"}
									</td>
								</tr>
							{:else}
								{#each data.partners as partner (partner.id)}
									{@const user = getUserForPartner(partner)}
									<tr class="border-b border-border-card hover:bg-muted/30 transition-colors">
										<td class="px-4 py-3">
											<div class="flex items-center gap-3">
												<div class="w-9 h-9 rounded-full bg-blue-100 dark:bg-blue-900/30 flex items-center justify-center">
													<Building2 size={16} class="text-blue-600 dark:text-blue-400" />
												</div>
												<div>
													<p class="text-sm font-medium text-foreground">
														{user ? user.email : partner.id}
													</p>
													{#if user && (user.firstname || user.lastname)}
														<p class="text-xs text-muted-foreground">
															{user.firstname} {user.lastname}
														</p>
													{/if}
												</div>
											</div>
										</td>
										<td class="px-4 py-3">
											<p class="text-sm text-muted-foreground max-w-[240px]">
												{truncate(partner.bio, 80)}
											</p>
										</td>
										<td class="px-4 py-3">
											<span class={getVerificationBadgeClass(partner.stripeOnboardingComplete && partner.stripeAccountStatus === "active")}>
												{#if partner.stripeOnboardingComplete && partner.stripeAccountStatus === "active"}
													<BadgeCheck size={12} />
													Vérifié
												{:else}
													<Clock size={12} />
													Non vérifié
												{/if}
											</span>
										</td>
										<td class="px-4 py-3">
											<span class={getStripeBadgeClass(partner.stripeOnboardingComplete, partner.stripeAccountStatus)}>
												<CreditCard size={12} />
												{#if partner.stripeOnboardingComplete && partner.stripeAccountStatus === "active"}
													Complet
												{:else if partner.stripeAccountStatus === "pending"}
													En attente
												{:else}
													Incomplet
												{/if}
											</span>
										</td>
										<td class="px-4 py-3">
											<div class="flex items-center gap-2">
												<span class="inline-flex items-center gap-1 px-2 py-0.5 rounded bg-muted text-xs text-foreground-alt">
													<FileText size={10} />
													{partner.categoryCount} cat.
												</span>
												<span class="inline-flex items-center gap-1 px-2 py-0.5 rounded bg-muted text-xs text-foreground-alt">
													{partner.productCount} prod.
												</span>
											</div>
										</td>
										<td class="px-4 py-3 text-sm text-muted-foreground">
											{formatDate(partner.createdAt)}
										</td>
										<td class="px-4 py-3">
											<div class="flex items-center justify-end gap-1">
												{#if !(partner.stripeOnboardingComplete && partner.stripeAccountStatus === "active")}
													<button
														onclick={() => openVerifyPartnerDialog(partner)}
														class="inline-flex items-center gap-1.5 px-3 py-1.5 text-sm font-medium text-blue-700 bg-blue-100 hover:bg-blue-200 dark:bg-blue-900/30 dark:text-blue-300 dark:hover:bg-blue-900/50 rounded-md transition-colors"
														title="Vérifier"
													>
														<BadgeCheck size={14} />
														Vérifier
													</button>
												{/if}
												<button
													onclick={() => openDeletePartnerDialog(partner)}
													class="p-2 rounded-md text-red-600 hover:text-red-700 hover:bg-red-50 dark:hover:bg-red-900/20 transition-colors"
													title="Supprimer"
												>
													<Trash2 size={16} />
												</button>
											</div>
										</td>
									</tr>
								{/each}
							{/if}
						</tbody>
					</table>
				</div>
			</div>
		</TabsContent>
	</Tabs>
</div>

<!-- Delete User Dialog/Drawer -->
{#if selectedUser}
	{#if isMobile}
		<Drawer bind:isOpen={deleteDialogOpen}>
			<div class="sticky top-0 bg-background pb-4 border-b border-border-card -mx-4 px-4 -mt-4 pt-4 z-10">
				<div class="flex items-center justify-between mb-2">
					<h2 class="text-xl font-semibold tracking-tight">
						Supprimer l'Utilisateur
					</h2>
					<button
						type="button"
						onclick={() => (deleteDialogOpen = false)}
						class="p-2 hover:bg-muted rounded-md transition-colors"
					>
						<X class="text-foreground size-5" />
					</button>
				</div>
				<p class="text-muted-foreground text-sm">
					Êtes-vous sûr de vouloir supprimer "{selectedUser.email}" ? Cette
					action ne peut pas être annulée.
				</p>
			</div>

			<div class="pt-4">
				<form method="POST" action="?/deleteUser" use:enhance={() => {
					return async ({ update }) => {
						await update({ reset: false });
					};
				}} class="flex w-full justify-end gap-3">
					<input type="hidden" name="id" value={selectedUser.id} />
					<button
						type="button"
						onclick={() => (deleteDialogOpen = false)}
						class="px-4 py-2 border border-border-input text-foreground-alt rounded-lg hover:bg-muted transition-colors font-medium text-sm"
					>
						Annuler
					</button>
					<button
						type="submit"
						class="px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700 transition-colors font-medium text-sm"
					>
						Supprimer
					</button>
				</form>
			</div>
		</Drawer>
	{:else}
		<Modal
			bind:isOpen={deleteDialogOpen}
			title="Supprimer l'Utilisateur"
			description="Êtes-vous sûr de vouloir supprimer '{selectedUser.email}' ? Cette action ne peut pas être annulée."
		>
			<form method="POST" action="?/deleteUser" use:enhance={() => {
				return async ({ update }) => {
					await update({ reset: false });
				};
			}} class="mt-6">
				<input type="hidden" name="id" value={selectedUser.id} />
				<div class="flex w-full justify-end gap-3">
					<button
						type="button"
						onclick={() => (deleteDialogOpen = false)}
						class="px-6 py-2.5 border border-border-input text-foreground-alt rounded-lg hover:bg-muted transition-colors font-medium"
					>
						Annuler
					</button>
					<button
						type="submit"
						class="px-6 py-2.5 bg-red-600 text-white rounded-lg hover:bg-red-700 transition-colors font-medium"
					>
						Supprimer
					</button>
				</div>
			</form>
		</Modal>
	{/if}

	<!-- Edit User Role Dialog/Drawer -->
	{#if isMobile}
		<Drawer bind:isOpen={editDialogOpen}>
			<div class="sticky top-0 bg-background pb-4 border-b border-border-card -mx-4 px-4 -mt-4 pt-4 z-10">
				<div class="flex items-center justify-between mb-2">
					<h2 class="text-xl font-semibold tracking-tight">
						Modifier le Rôle
					</h2>
					<button
						type="button"
						onclick={() => (editDialogOpen = false)}
						class="p-2 hover:bg-muted rounded-md transition-colors"
					>
						<X class="text-foreground size-5" />
					</button>
				</div>
				<p class="text-muted-foreground text-sm">
					Modifier le rôle de "{selectedUser.email}"
				</p>
			</div>

			<div class="pt-4">
				<form method="POST" action="?/updateRole" use:enhance={() => {
					return async ({ update }) => {
						await update({ reset: false });
					};
				}} class="space-y-4">
					<input type="hidden" name="id" value={selectedUser.id} />
					<div>
						<label for="role" class="block text-sm font-medium text-foreground mb-1.5">
							Rôle
						</label>
						<select
							id="role"
							name="role"
							bind:value={selectedRole}
							class="w-full px-3 py-2 border border-border-input rounded-lg bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-foreground/20"
						>
							{#each availableRoles as role}
								<option value={role.value}>{role.label}</option>
							{/each}
						</select>
					</div>
					<div class="flex w-full justify-end gap-3 pt-2">
						<button
							type="button"
							onclick={() => (editDialogOpen = false)}
							class="px-4 py-2 border border-border-input text-foreground-alt rounded-lg hover:bg-muted transition-colors font-medium text-sm"
						>
							Annuler
						</button>
						<button
							type="submit"
							class="px-4 py-2 bg-foreground text-background rounded-lg hover:bg-foreground/90 transition-colors font-medium text-sm"
						>
							Enregistrer
						</button>
					</div>
				</form>
			</div>
		</Drawer>
	{:else}
		<Modal
			bind:isOpen={editDialogOpen}
			title="Modifier le Rôle"
			description="Modifier le rôle de '{selectedUser.email}'"
		>
			<form method="POST" action="?/updateRole" use:enhance={() => {
				return async ({ update }) => {
					await update({ reset: false });
				};
			}} class="mt-6 space-y-4">
				<input type="hidden" name="id" value={selectedUser.id} />
				<div>
					<label for="role" class="block text-sm font-medium text-foreground mb-1.5">
						Rôle
					</label>
					<select
						id="role"
						name="role"
						bind:value={selectedRole}
						class="w-full px-3 py-2 border border-border-input rounded-lg bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-foreground/20"
					>
						{#each availableRoles as role}
							<option value={role.value}>{role.label}</option>
						{/each}
					</select>
				</div>
				<div class="flex w-full justify-end gap-3 pt-2">
					<button
						type="button"
						onclick={() => (editDialogOpen = false)}
						class="px-6 py-2.5 border border-border-input text-foreground-alt rounded-lg hover:bg-muted transition-colors font-medium"
					>
						Annuler
					</button>
					<button
						type="submit"
						class="px-6 py-2.5 bg-foreground text-background rounded-lg hover:bg-foreground/90 transition-colors font-medium"
					>
						Enregistrer
					</button>
				</div>
			</form>
		</Modal>
	{/if}

	<!-- Approve User Dialog/Drawer -->
	{#if isMobile}
		<Drawer bind:isOpen={approveDialogOpen}>
			<div class="sticky top-0 bg-background pb-4 border-b border-border-card -mx-4 px-4 -mt-4 pt-4 z-10">
				<div class="flex items-center justify-between mb-2">
					<h2 class="text-xl font-semibold tracking-tight">
						Approuver l'Utilisateur
					</h2>
					<button
						type="button"
						onclick={() => (approveDialogOpen = false)}
						class="p-2 hover:bg-muted rounded-md transition-colors"
					>
						<X class="text-foreground size-5" />
					</button>
				</div>
				<p class="text-muted-foreground text-sm">
					Êtes-vous sûr de vouloir approuver "{selectedUser?.email}" ?
				</p>
			</div>

			<div class="pt-4">
				<form method="POST" action="?/approveUser" use:enhance={() => {
					return async ({ update }) => {
						await update({ reset: false });
					};
				}} class="flex w-full justify-end gap-3">
					<input type="hidden" name="id" value={selectedUser?.id} />
					<button
						type="button"
						onclick={() => (approveDialogOpen = false)}
						class="px-4 py-2 border border-border-input text-foreground-alt rounded-lg hover:bg-muted transition-colors font-medium text-sm"
					>
						Annuler
					</button>
					<button
						type="submit"
						class="px-4 py-2 bg-green-600 text-white rounded-lg hover:bg-green-700 transition-colors font-medium text-sm"
					>
						Approuver
					</button>
				</form>
			</div>
		</Drawer>
	{:else}
		<Modal
			bind:isOpen={approveDialogOpen}
			title="Approuver l'Utilisateur"
			description="Êtes-vous sûr de vouloir approuver '{selectedUser?.email}' ?"
		>
			<form method="POST" action="?/approveUser" use:enhance={() => {
				return async ({ update }) => {
					await update({ reset: false });
				};
			}} class="mt-6">
				<input type="hidden" name="id" value={selectedUser?.id} />
				<div class="flex w-full justify-end gap-3">
					<button
						type="button"
						onclick={() => (approveDialogOpen = false)}
						class="px-6 py-2.5 border border-border-input text-foreground-alt rounded-lg hover:bg-muted transition-colors font-medium"
					>
						Annuler
					</button>
					<button
						type="submit"
						class="px-6 py-2.5 bg-green-600 text-white rounded-lg hover:bg-green-700 transition-colors font-medium"
					>
						Approuver
					</button>
				</div>
			</form>
		</Modal>
	{/if}
{/if}

<!-- Verify Partner Dialog/Drawer -->
{#if selectedPartner}
	{@const partnerUser = getUserForPartner(selectedPartner)}
	{#if isMobile}
		<Drawer bind:isOpen={verifyPartnerDialogOpen}>
			<div class="sticky top-0 bg-background pb-4 border-b border-border-card -mx-4 px-4 -mt-4 pt-4 z-10">
				<div class="flex items-center justify-between mb-2">
					<h2 class="text-xl font-semibold tracking-tight">
						Vérifier le Partenaire
					</h2>
					<button
						type="button"
						onclick={() => (verifyPartnerDialogOpen = false)}
						class="p-2 hover:bg-muted rounded-md transition-colors"
					>
						<X class="text-foreground size-5" />
					</button>
				</div>
				<p class="text-muted-foreground text-sm">
					Êtes-vous sûr de vouloir vérifier le partenaire "{partnerUser?.email ?? selectedPartner.id}" ?
					Cela activera son compte et lui attribuera le rôle partenaire.
				</p>
			</div>

			<div class="pt-4">
				<form method="POST" action="?/verifyPartner" use:enhance={() => {
					return async ({ update }) => {
						await update({ reset: false });
					};
				}} class="flex w-full justify-end gap-3">
					<input type="hidden" name="id" value={selectedPartner.id} />
					<button
						type="button"
						onclick={() => (verifyPartnerDialogOpen = false)}
						class="px-4 py-2 border border-border-input text-foreground-alt rounded-lg hover:bg-muted transition-colors font-medium text-sm"
					>
						Annuler
					</button>
					<button
						type="submit"
						class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors font-medium text-sm"
					>
						Vérifier
					</button>
				</form>
			</div>
		</Drawer>
	{:else}
		<Modal
			bind:isOpen={verifyPartnerDialogOpen}
			title="Vérifier le Partenaire"
			description="Êtes-vous sûr de vouloir vérifier le partenaire '{partnerUser?.email ?? selectedPartner.id}' ? Cela activera son compte et lui attribuera le rôle partenaire."
		>
			<form method="POST" action="?/verifyPartner" use:enhance={() => {
				return async ({ update }) => {
					await update({ reset: false });
				};
			}} class="mt-6">
				<input type="hidden" name="id" value={selectedPartner.id} />
				<div class="flex w-full justify-end gap-3">
					<button
						type="button"
						onclick={() => (verifyPartnerDialogOpen = false)}
						class="px-6 py-2.5 border border-border-input text-foreground-alt rounded-lg hover:bg-muted transition-colors font-medium"
					>
						Annuler
					</button>
					<button
						type="submit"
						class="px-6 py-2.5 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors font-medium"
					>
						Vérifier
					</button>
				</div>
			</form>
		</Modal>
	{/if}

	<!-- Delete Partner Dialog/Drawer -->
	{#if isMobile}
		<Drawer bind:isOpen={deletePartnerDialogOpen}>
			<div class="sticky top-0 bg-background pb-4 border-b border-border-card -mx-4 px-4 -mt-4 pt-4 z-10">
				<div class="flex items-center justify-between mb-2">
					<h2 class="text-xl font-semibold tracking-tight">
						Supprimer le Partenaire
					</h2>
					<button
						type="button"
						onclick={() => (deletePartnerDialogOpen = false)}
						class="p-2 hover:bg-muted rounded-md transition-colors"
					>
						<X class="text-foreground size-5" />
					</button>
				</div>
				<p class="text-muted-foreground text-sm">
					Êtes-vous sûr de vouloir supprimer le partenaire "{partnerUser?.email ?? selectedPartner.id}" ?
					Cette action ne peut pas être annulée.
				</p>
			</div>

			<div class="pt-4">
				<form method="POST" action="?/deletePartner" use:enhance={() => {
					return async ({ update }) => {
						await update({ reset: false });
					};
				}} class="flex w-full justify-end gap-3">
					<input type="hidden" name="id" value={selectedPartner.id} />
					<button
						type="button"
						onclick={() => (deletePartnerDialogOpen = false)}
						class="px-4 py-2 border border-border-input text-foreground-alt rounded-lg hover:bg-muted transition-colors font-medium text-sm"
					>
						Annuler
					</button>
					<button
						type="submit"
						class="px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700 transition-colors font-medium text-sm"
					>
						Supprimer
					</button>
				</form>
			</div>
		</Drawer>
	{:else}
		<Modal
			bind:isOpen={deletePartnerDialogOpen}
			title="Supprimer le Partenaire"
			description="Êtes-vous sûr de vouloir supprimer le partenaire '{partnerUser?.email ?? selectedPartner.id}' ? Cette action ne peut pas être annulée."
		>
			<form method="POST" action="?/deletePartner" use:enhance={() => {
				return async ({ update }) => {
					await update({ reset: false });
				};
			}} class="mt-6">
				<input type="hidden" name="id" value={selectedPartner.id} />
				<div class="flex w-full justify-end gap-3">
					<button
						type="button"
						onclick={() => (deletePartnerDialogOpen = false)}
						class="px-6 py-2.5 border border-border-input text-foreground-alt rounded-lg hover:bg-muted transition-colors font-medium"
					>
						Annuler
					</button>
					<button
						type="submit"
						class="px-6 py-2.5 bg-red-600 text-white rounded-lg hover:bg-red-700 transition-colors font-medium"
					>
						Supprimer
					</button>
				</div>
			</form>
		</Modal>
	{/if}
{/if}
