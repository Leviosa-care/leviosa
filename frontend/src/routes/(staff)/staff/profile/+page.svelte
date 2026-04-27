<script lang="ts">
	import type { PageProps } from './$types';
	import { BadgeCheck, Pencil, Mail, MapPin, Phone, Tag, Package } from '@lucide/svelte';

	let { data }: PageProps = $props();

	let editingBio = $state(false);
	let editingExperience = $state(false);
	let bioValue = $state(data.profile.bio);
	let experienceValue = $state(data.profile.experience);

	function formatDate(iso: string): string {
		return new Date(iso).toLocaleDateString('fr-FR', {
			day: 'numeric',
			month: 'long',
			year: 'numeric',
		});
	}

	const initials = $derived(
		data.user.firstname && data.user.lastname
			? `${data.user.firstname[0]}${data.user.lastname[0]}`.toUpperCase()
			: data.user.email[0].toUpperCase(),
	);
</script>

<svelte:head>
	<title>Mon profil | Staff</title>
</svelte:head>

<div class="container mx-auto px-4 py-8 lg:py-12 max-w-4xl">
	<div class="mb-8">
		<h1 class="text-3xl lg:text-4xl font-bold mb-1 text-foreground">Mon profil</h1>
		<p class="text-muted-foreground">Gérez vos informations de praticien</p>
	</div>

	<!-- Profile Header -->
	<div class="bg-card rounded-lg border border-border p-6 mb-6">
		<div class="flex flex-col sm:flex-row items-start sm:items-center gap-5">
			<div
				class="w-20 h-20 rounded-full bg-muted flex items-center justify-center flex-shrink-0 text-2xl font-bold text-foreground"
			>
				{initials}
			</div>
			<div class="flex-1 min-w-0">
				<div class="flex flex-wrap items-center gap-2 mb-1">
					<h2 class="text-xl font-bold text-foreground">
						{data.user.firstname}
						{data.user.lastname}
					</h2>
					{#if data.profile.isVerified}
						<span
							class="inline-flex items-center gap-1 px-2.5 py-1 text-xs font-medium bg-green-100 text-green-700 rounded-full"
						>
							<BadgeCheck size={13} />
							Vérifié
						</span>
					{/if}
				</div>
				<p class="text-sm text-muted-foreground mb-3">
					Partenaire depuis {formatDate(data.profile.joinedAt)}
				</p>
				<div class="flex flex-wrap gap-4 text-sm text-muted-foreground">
					{#if data.user.email}
						<span class="flex items-center gap-1.5">
							<Mail size={14} />
							{data.user.email}
						</span>
					{/if}
					{#if data.user.telephone}
						<span class="flex items-center gap-1.5">
							<Phone size={14} />
							{data.user.telephone}
						</span>
					{/if}
					{#if data.user.city}
						<span class="flex items-center gap-1.5">
							<MapPin size={14} />
							{data.user.city}
						</span>
					{/if}
				</div>
			</div>
		</div>
	</div>

	<div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
		<!-- Left Column: Bio + Experience -->
		<div class="lg:col-span-2 space-y-6">
			<!-- Bio -->
			<div class="bg-card rounded-lg border border-border p-6">
				<div class="flex items-center justify-between mb-4">
					<h3 class="font-semibold text-foreground">Biographie</h3>
					<button
						class="p-1.5 rounded-md text-muted-foreground hover:text-foreground hover:bg-muted transition-colors"
						onclick={() => (editingBio = !editingBio)}
						title="Modifier"
					>
						<Pencil size={15} />
					</button>
				</div>
				{#if editingBio}
					<div class="space-y-3">
						<textarea
							bind:value={bioValue}
							rows="5"
							maxlength="1000"
							class="w-full px-3 py-2 rounded-lg border border-border bg-background text-sm resize-none focus:outline-none focus:ring-2 focus:ring-foreground/20"
						></textarea>
						<div class="flex items-center justify-between">
							<span class="text-xs text-muted-foreground">{bioValue.length}/1000</span>
							<div class="flex gap-2">
								<button
									class="px-3 py-1.5 text-sm rounded-md border border-border hover:bg-muted transition-colors"
									onclick={() => { bioValue = data.profile.bio; editingBio = false; }}
								>
									Annuler
								</button>
								<button
									class="px-3 py-1.5 text-sm rounded-md bg-foreground text-background hover:opacity-90 transition-opacity"
									onclick={() => (editingBio = false)}
								>
									Enregistrer
								</button>
							</div>
						</div>
					</div>
				{:else}
					<p class="text-sm text-muted-foreground leading-relaxed">{bioValue}</p>
				{/if}
			</div>

			<!-- Experience -->
			<div class="bg-card rounded-lg border border-border p-6">
				<div class="flex items-center justify-between mb-4">
					<h3 class="font-semibold text-foreground">Expérience</h3>
					<button
						class="p-1.5 rounded-md text-muted-foreground hover:text-foreground hover:bg-muted transition-colors"
						onclick={() => (editingExperience = !editingExperience)}
						title="Modifier"
					>
						<Pencil size={15} />
					</button>
				</div>
				{#if editingExperience}
					<div class="space-y-3">
						<textarea
							bind:value={experienceValue}
							rows="7"
							maxlength="2000"
							class="w-full px-3 py-2 rounded-lg border border-border bg-background text-sm resize-none focus:outline-none focus:ring-2 focus:ring-foreground/20"
						></textarea>
						<div class="flex items-center justify-between">
							<span class="text-xs text-muted-foreground">{experienceValue.length}/2000</span>
							<div class="flex gap-2">
								<button
									class="px-3 py-1.5 text-sm rounded-md border border-border hover:bg-muted transition-colors"
									onclick={() => { experienceValue = data.profile.experience; editingExperience = false; }}
								>
									Annuler
								</button>
								<button
									class="px-3 py-1.5 text-sm rounded-md bg-foreground text-background hover:opacity-90 transition-opacity"
									onclick={() => (editingExperience = false)}
								>
									Enregistrer
								</button>
							</div>
						</div>
					</div>
				{:else}
					<p class="text-sm text-muted-foreground leading-relaxed">{experienceValue}</p>
				{/if}
			</div>
		</div>

		<!-- Right Column: Categories + Products + Status -->
		<div class="space-y-6">
			<!-- Categories -->
			<div class="bg-card rounded-lg border border-border p-6">
				<div class="flex items-center gap-2 mb-4">
					<Tag size={16} class="text-muted-foreground" />
					<h3 class="font-semibold text-foreground">Spécialités</h3>
				</div>
				<div class="flex flex-wrap gap-2">
					{#each data.profile.categories as cat (cat.id)}
						<span class="px-3 py-1.5 text-xs font-medium bg-muted text-foreground rounded-full">
							{cat.name}
						</span>
					{/each}
				</div>
			</div>

			<!-- Products -->
			<div class="bg-card rounded-lg border border-border p-6">
				<div class="flex items-center gap-2 mb-4">
					<Package size={16} class="text-muted-foreground" />
					<h3 class="font-semibold text-foreground">Prestations</h3>
				</div>
				<ul class="space-y-2">
					{#each data.profile.products as prod (prod.id)}
						<li class="text-sm text-muted-foreground flex items-center gap-2">
							<span class="w-1.5 h-1.5 rounded-full bg-muted-foreground/50 flex-shrink-0"></span>
							{prod.name}
						</li>
					{/each}
				</ul>
			</div>

			<!-- Status -->
			<div class="bg-card rounded-lg border border-border p-6">
				<h3 class="font-semibold text-foreground mb-4">Statut du compte</h3>
				<div class="space-y-3">
					<div class="flex items-center justify-between">
						<span class="text-sm text-muted-foreground">Vérification</span>
						<span
							class="text-xs font-medium px-2 py-1 rounded-full {data.profile.isVerified
								? 'bg-green-100 text-green-700'
								: 'bg-yellow-100 text-yellow-700'}"
						>
							{data.profile.isVerified ? 'Vérifié' : 'En attente'}
						</span>
					</div>
					<div class="flex items-center justify-between">
						<span class="text-sm text-muted-foreground">Paiements Stripe</span>
						<span
							class="text-xs font-medium px-2 py-1 rounded-full {data.profile.stripeOnboardingComplete
								? 'bg-green-100 text-green-700'
								: 'bg-yellow-100 text-yellow-700'}"
						>
							{data.profile.stripeOnboardingComplete ? 'Actif' : 'À configurer'}
						</span>
					</div>
				</div>
			</div>
		</div>
	</div>
</div>
