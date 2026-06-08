<script lang="ts">
	import type { PageProps } from './$types';
	import { replaceState } from '$app/navigation';
	import { BadgeCheck, Pencil, Mail, MapPin, Phone, Tag, Package, AlertCircle, Link, Unlink, CreditCard, ExternalLink, Briefcase, MessageSquareQuote, Tags } from '@lucide/svelte';

	let { data }: PageProps = $props();

	let editingBio = $state(false);
	let editingExperience = $state(false);
	let editingOccupation = $state(false);
	let editingQuote = $state(false);
	let editingTags = $state(false);
	let bioValue = $state(data.profile?.bio ?? '');
	let experienceValue = $state(data.profile?.experience ?? '');
	let occupationValue = $state(data.profile?.occupation ?? '');
	let quoteValue = $state(data.profile?.quote ?? '');
	let tagsInput = $state((data.profile?.tags ?? []).join(', '));
	let tagsValue = $state<string[]>(data.profile?.tags ?? []);
	// Track last successfully saved values so cancel reverts to the last save, not page load.
	let savedBio = $state(data.profile?.bio ?? '');
	let savedExperience = $state(data.profile?.experience ?? '');
	let savedOccupation = $state(data.profile?.occupation ?? '');
	let savedQuote = $state(data.profile?.quote ?? '');
	let savedTagsInput = $state((data.profile?.tags ?? []).join(', '));
	let bioSaving = $state(false);
	let experienceSaving = $state(false);
	let occupationSaving = $state(false);
	let quoteSaving = $state(false);
	let tagsSaving = $state(false);
	let bioError = $state<string | null>(null);
	let experienceError = $state<string | null>(null);
	let occupationError = $state<string | null>(null);
	let quoteError = $state<string | null>(null);
	let tagsError = $state<string | null>(null);
	let oauthLoading = $state<string | null>(null);
	let oauthError = $state<string | null>(null);
	let linkedProviders = $state(data.linkedProviders ?? { google: false, apple: false });
	let stripeLoading = $state(false);
	let stripeError = $state<string | null>(null);

	// Strip the ?stripe=return query param from the URL without re-running the load.
	// The server already fetched fresh status on this request.
	$effect(() => {
		if (data.stripeCallback === 'return') {
			replaceState('/staff/profile', {});
		}
	});

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

	async function saveBio() {
		bioSaving = true;
		bioError = null;
		try {
			const res = await fetch(`/api/partners/me`, {
				method: 'PUT',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ bio: bioValue }),
			});
			if (res.ok) {
				savedBio = bioValue;
				editingBio = false;
			} else {
				bioError = 'Impossible d\'enregistrer. Veuillez réessayer.';
			}
		} catch {
			bioError = 'Une erreur est survenue. Veuillez réessayer.';
		} finally {
			bioSaving = false;
		}
	}

	async function saveExperience() {
		experienceSaving = true;
		experienceError = null;
		try {
			const res = await fetch(`/api/partners/me`, {
				method: 'PUT',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ experience: experienceValue }),
			});
			if (res.ok) {
				savedExperience = experienceValue;
				editingExperience = false;
			} else {
				experienceError = 'Impossible d\'enregistrer. Veuillez réessayer.';
			}
		} catch {
			experienceError = 'Une erreur est survenue. Veuillez réessayer.';
		} finally {
			experienceSaving = false;
		}
	}

	async function saveOccupation() {
		occupationSaving = true;
		occupationError = null;
		try {
			const res = await fetch(`/api/partners/me`, {
				method: 'PUT',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ occupation: occupationValue }),
			});
			if (res.ok) {
				savedOccupation = occupationValue;
				editingOccupation = false;
			} else {
				occupationError = 'Impossible d\'enregistrer. Veuillez réessayer.';
			}
		} catch {
			occupationError = 'Une erreur est survenue. Veuillez réessayer.';
		} finally {
			occupationSaving = false;
		}
	}

	async function saveQuote() {
		quoteSaving = true;
		quoteError = null;
		try {
			const res = await fetch(`/api/partners/me`, {
				method: 'PUT',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ quote: quoteValue }),
			});
			if (res.ok) {
				savedQuote = quoteValue;
				editingQuote = false;
			} else {
				quoteError = 'Impossible d\'enregistrer. Veuillez réessayer.';
			}
		} catch {
			quoteError = 'Une erreur est survenue. Veuillez réessayer.';
		} finally {
			quoteSaving = false;
		}
	}

	async function saveTags() {
		tagsSaving = true;
		tagsError = null;
		const parsedTags = tagsInput
			.split(',')
			.map((t: string) => t.trim())
			.filter((t: string) => t.length > 0);
		try {
			const res = await fetch(`/api/partners/me`, {
				method: 'PUT',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ tags: parsedTags }),
			});
			if (res.ok) {
				tagsValue = parsedTags;
				savedTagsInput = tagsInput;
				editingTags = false;
			} else {
				tagsError = 'Impossible d\'enregistrer. Veuillez réessayer.';
			}
		} catch {
			tagsError = 'Une erreur est survenue. Veuillez réessayer.';
		} finally {
			tagsSaving = false;
		}
	}

	const needsStripeOnboarding = $derived(
		data.profile && (data.profile.stripeAccountStatus === 'pending' || data.profile.stripeAccountStatus === 'restricted')
	);

	async function startStripeOnboarding() {
		stripeLoading = true;
		stripeError = null;
		try {
			const res = await fetch('/api/partners/me/stripe/onboarding-link', { method: 'POST' });
			if (res.ok) {
				const { url } = await res.json();
				window.location.href = url;
			} else {
				stripeError = 'Impossible de générer le lien Stripe. Veuillez réessayer.';
			}
		} catch {
			stripeError = 'Une erreur est survenue. Veuillez réessayer.';
		} finally {
			stripeLoading = false;
		}
	}

	async function linkProvider(provider: string) {
		oauthLoading = provider;
		oauthError = null;
		try {
			const res = await fetch(`/api/users/me/oauth/${provider}/link`, {
				method: 'POST',
			});
			if (res.ok) {
				const linkData = await res.json();
				if (linkData.authorization_url) {
					window.location.href = linkData.authorization_url;
				}
			} else {
				oauthError = `Impossible de lier le compte ${provider}. Veuillez réessayer.`;
			}
		} catch {
			oauthError = 'Une erreur est survenue. Veuillez réessayer.';
		} finally {
			oauthLoading = null;
		}
	}

	async function unlinkProvider(provider: string) {
		oauthLoading = provider;
		oauthError = null;
		try {
			const res = await fetch(`/api/users/me/oauth/${provider}/unlink`, {
				method: 'DELETE',
			});
			if (res.ok) {
				linkedProviders = { ...linkedProviders, [provider]: false };
			} else {
				oauthError = `Impossible de délier le compte ${provider}. Veuillez réessayer.`;
			}
		} catch {
			oauthError = 'Une erreur est survenue. Veuillez réessayer.';
		} finally {
			oauthLoading = null;
		}
	}
</script>

<svelte:head>
	<title>Mon profil | Staff</title>
</svelte:head>

<div class="p-6 lg:p-10 max-w-4xl">
	<div class="mb-10">
		<h1 class="font-display text-3xl lg:text-4xl font-semibold tracking-tight text-foreground">Mon profil</h1>
		<p class="text-muted-foreground mt-1.5 text-base">Gérez vos informations de praticien</p>
	</div>

	{#if data.error}
		<div class="bg-red-50 border border-red-200 rounded-lg p-6 mb-6">
			<div class="flex items-center gap-3">
				<AlertCircle class="text-red-600" size={20} />
				<div>
					<h2 class="font-semibold text-red-900">Erreur de chargement</h2>
					<p class="text-sm text-red-700">{data.error}</p>
				</div>
			</div>
		</div>
	{:else if !data.profile}
		<div class="bg-yellow-50 border border-yellow-200 rounded-lg p-6 mb-6">
			<div class="flex items-center gap-3">
				<AlertCircle class="text-yellow-600" size={20} />
				<div>
					<h2 class="font-semibold text-yellow-900">Profil non disponible</h2>
					<p class="text-sm text-yellow-700">Votre profil n'a pas pu être chargé. Veuillez réessayer.</p>
				</div>
			</div>
		</div>
	{:else}

	{#if needsStripeOnboarding}
		<div class="bg-amber-50 border border-amber-200 rounded-lg p-6 mb-6">
			<div class="flex flex-col sm:flex-row items-start gap-4">
				<div class="flex items-center gap-3 flex-1">
					<CreditCard class="text-amber-600 flex-shrink-0" size={22} />
					<div>
						<h2 class="font-semibold text-amber-900">Configuration Stripe requise</h2>
						<p class="text-sm text-amber-700 mt-1">
							{data.stripeCallback === 'refresh'
								? 'Le lien précédent a expiré. Veuillez générer un nouveau lien pour poursuivre la configuration de votre compte Stripe.'
								: 'Pour recevoir vos paiements, vous devez configurer votre compte Stripe en complétant les informations bancaires et d\'identité.'}
						</p>
					</div>
				</div>
				<button
					class="flex items-center gap-2 px-4 py-2.5 text-sm font-medium rounded-lg bg-amber-600 text-white hover:bg-amber-700 transition-colors disabled:opacity-50 flex-shrink-0"
					onclick={startStripeOnboarding}
					disabled={stripeLoading}
				>
					{stripeLoading ? 'Chargement...' : 'Configurer Stripe'}
					{#if !stripeLoading}
						<ExternalLink size={14} />
					{/if}
				</button>
			</div>
			{#if stripeError}
				<p class="text-xs text-red-600 mt-3">{stripeError}</p>
			{/if}
		</div>
	{/if}

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
						{#if bioError}
							<p class="text-xs text-red-600">{bioError}</p>
						{/if}
						<div class="flex items-center justify-between">
							<span class="text-xs text-muted-foreground">{bioValue.length}/1000</span>
							<div class="flex gap-2">
								<button
									class="px-3 py-1.5 text-sm rounded-md border border-border hover:bg-muted transition-colors"
									onclick={() => { bioValue = savedBio; editingBio = false; bioError = null; }}
									disabled={bioSaving}
								>
									Annuler
								</button>
								<button
									class="px-3 py-1.5 text-sm rounded-md bg-foreground text-background hover:opacity-90 transition-opacity disabled:opacity-50"
									onclick={saveBio}
									disabled={bioSaving}
								>
									{bioSaving ? 'Enregistrement...' : 'Enregistrer'}
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
						{#if experienceError}
							<p class="text-xs text-red-600">{experienceError}</p>
						{/if}
						<div class="flex items-center justify-between">
							<span class="text-xs text-muted-foreground">{experienceValue.length}/2000</span>
							<div class="flex gap-2">
								<button
									class="px-3 py-1.5 text-sm rounded-md border border-border hover:bg-muted transition-colors"
									onclick={() => { experienceValue = savedExperience; editingExperience = false; experienceError = null; }}
									disabled={experienceSaving}
								>
									Annuler
								</button>
								<button
									class="px-3 py-1.5 text-sm rounded-md bg-foreground text-background hover:opacity-90 transition-opacity disabled:opacity-50"
									onclick={saveExperience}
									disabled={experienceSaving}
								>
									{experienceSaving ? 'Enregistrement...' : 'Enregistrer'}
								</button>
							</div>
						</div>
					</div>
				{:else}
					<p class="text-sm text-muted-foreground leading-relaxed">{experienceValue}</p>
				{/if}
			</div>

			<!-- Occupation -->
			<div class="bg-card rounded-lg border border-border p-6">
				<div class="flex items-center justify-between mb-4">
					<div class="flex items-center gap-2">
						<Briefcase size={16} class="text-muted-foreground" />
						<h3 class="font-semibold text-foreground">Profession</h3>
					</div>
					<button
						class="p-1.5 rounded-md text-muted-foreground hover:text-foreground hover:bg-muted transition-colors"
						onclick={() => (editingOccupation = !editingOccupation)}
						title="Modifier"
					>
						<Pencil size={15} />
					</button>
				</div>
				{#if editingOccupation}
					<div class="space-y-3">
						<input
							type="text"
							bind:value={occupationValue}
							maxlength="200"
							placeholder="Ex: Kinésithérapeute du sport"
							class="w-full px-3 py-2 rounded-lg border border-border bg-background text-sm focus:outline-none focus:ring-2 focus:ring-foreground/20"
						/>
						{#if occupationError}
							<p class="text-xs text-red-600">{occupationError}</p>
						{/if}
						<div class="flex items-center justify-between">
							<span class="text-xs text-muted-foreground">{occupationValue.length}/200</span>
							<div class="flex gap-2">
								<button
									class="px-3 py-1.5 text-sm rounded-md border border-border hover:bg-muted transition-colors"
									onclick={() => { occupationValue = savedOccupation; editingOccupation = false; occupationError = null; }}
									disabled={occupationSaving}
								>
									Annuler
								</button>
								<button
									class="px-3 py-1.5 text-sm rounded-md bg-foreground text-background hover:opacity-90 transition-opacity disabled:opacity-50"
									onclick={saveOccupation}
									disabled={occupationSaving}
								>
									{occupationSaving ? 'Enregistrement...' : 'Enregistrer'}
								</button>
							</div>
						</div>
					</div>
				{:else}
					<p class="text-sm text-muted-foreground leading-relaxed">{occupationValue || 'Non renseigné'}</p>
				{/if}
			</div>

			<!-- Quote -->
			<div class="bg-card rounded-lg border border-border p-6">
				<div class="flex items-center justify-between mb-4">
					<div class="flex items-center gap-2">
						<MessageSquareQuote size={16} class="text-muted-foreground" />
						<h3 class="font-semibold text-foreground">Citation</h3>
					</div>
					<button
						class="p-1.5 rounded-md text-muted-foreground hover:text-foreground hover:bg-muted transition-colors"
						onclick={() => (editingQuote = !editingQuote)}
						title="Modifier"
					>
						<Pencil size={15} />
					</button>
				</div>
				{#if editingQuote}
					<div class="space-y-3">
						<input
							type="text"
							bind:value={quoteValue}
							maxlength="300"
							placeholder="Ex: Le mouvement est la vie"
							class="w-full px-3 py-2 rounded-lg border border-border bg-background text-sm focus:outline-none focus:ring-2 focus:ring-foreground/20"
						/>
						{#if quoteError}
							<p class="text-xs text-red-600">{quoteError}</p>
						{/if}
						<div class="flex items-center justify-between">
							<span class="text-xs text-muted-foreground">{quoteValue.length}/300</span>
							<div class="flex gap-2">
								<button
									class="px-3 py-1.5 text-sm rounded-md border border-border hover:bg-muted transition-colors"
									onclick={() => { quoteValue = savedQuote; editingQuote = false; quoteError = null; }}
									disabled={quoteSaving}
								>
									Annuler
								</button>
								<button
									class="px-3 py-1.5 text-sm rounded-md bg-foreground text-background hover:opacity-90 transition-opacity disabled:opacity-50"
									onclick={saveQuote}
									disabled={quoteSaving}
								>
									{quoteSaving ? 'Enregistrement...' : 'Enregistrer'}
								</button>
							</div>
						</div>
					</div>
				{:else}
					<p class="text-sm text-muted-foreground leading-relaxed italic">{quoteValue ? `« ${quoteValue} »` : 'Non renseigné'}</p>
				{/if}
			</div>

			<!-- Tags -->
			<div class="bg-card rounded-lg border border-border p-6">
				<div class="flex items-center justify-between mb-4">
					<div class="flex items-center gap-2">
						<Tags size={16} class="text-muted-foreground" />
						<h3 class="font-semibold text-foreground">Tags de spécialité</h3>
					</div>
					<button
						class="p-1.5 rounded-md text-muted-foreground hover:text-foreground hover:bg-muted transition-colors"
						onclick={() => (editingTags = !editingTags)}
						title="Modifier"
					>
						<Pencil size={15} />
					</button>
				</div>
				{#if editingTags}
					<div class="space-y-3">
						<input
							type="text"
							bind:value={tagsInput}
							placeholder="Ex: sport, rééducation, blessures"
							class="w-full px-3 py-2 rounded-lg border border-border bg-background text-sm focus:outline-none focus:ring-2 focus:ring-foreground/20"
						/>
						<p class="text-xs text-muted-foreground">Séparez les tags par des virgules</p>
						{#if tagsError}
							<p class="text-xs text-red-600">{tagsError}</p>
						{/if}
						<div class="flex items-center justify-between">
							<span class="text-xs text-muted-foreground">{tagsInput.split(',').filter((t: string) => t.trim()).length} tag(s)</span>
							<div class="flex gap-2">
								<button
									class="px-3 py-1.5 text-sm rounded-md border border-border hover:bg-muted transition-colors"
									onclick={() => { tagsInput = savedTagsInput; editingTags = false; tagsError = null; }}
									disabled={tagsSaving}
								>
									Annuler
								</button>
								<button
									class="px-3 py-1.5 text-sm rounded-md bg-foreground text-background hover:opacity-90 transition-opacity disabled:opacity-50"
									onclick={saveTags}
									disabled={tagsSaving}
								>
									{tagsSaving ? 'Enregistrement...' : 'Enregistrer'}
								</button>
							</div>
						</div>
					</div>
				{:else}
					{#if tagsValue.length > 0}
						<div class="flex flex-wrap gap-2">
							{#each tagsValue as tag (tag)}
								<span class="px-3 py-1.5 text-xs font-medium bg-muted text-foreground rounded-full">
									{tag}
								</span>
							{/each}
						</div>
					{:else}
						<p class="text-sm text-muted-foreground">Non renseigné</p>
					{/if}
				{/if}
			</div>
		</div>

		<!-- Right Column: Categories + Products + Status + OAuth -->
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

			<!-- OAuth Account Linking -->
			<div class="bg-card rounded-lg border border-border p-6">
				<h3 class="font-semibold text-foreground mb-4">Identifiants de connexion</h3>
				{#if oauthError}
					<p class="text-xs text-red-600 mb-3">{oauthError}</p>
				{/if}
				<div class="space-y-3">
					<!-- Google -->
					<div class="flex items-center justify-between">
						<div class="flex items-center gap-2">
							<svg class="w-4 h-4" viewBox="0 0 24 24" fill="none">
								<path d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92a5.06 5.06 0 0 1-2.2 3.32v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.1z" fill="#4285F4"/>
								<path d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z" fill="#34A853"/>
								<path d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z" fill="#FBBC05"/>
								<path d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z" fill="#EA4335"/>
							</svg>
							<span class="text-sm text-muted-foreground">Google</span>
							{#if linkedProviders.google}
								<span class="text-xs font-medium px-2 py-0.5 rounded-full bg-green-100 text-green-700">Lié</span>
							{/if}
						</div>
						{#if linkedProviders.google}
							<button
								class="flex items-center gap-1 px-2.5 py-1 text-xs rounded-md border border-border text-muted-foreground hover:text-red-600 hover:border-red-300 transition-colors disabled:opacity-50"
								onclick={() => unlinkProvider('google')}
								disabled={oauthLoading === 'google'}
							>
								<Unlink size={12} />
								Délier
							</button>
						{:else}
							<button
								class="flex items-center gap-1 px-2.5 py-1 text-xs rounded-md bg-foreground text-background hover:opacity-90 transition-opacity disabled:opacity-50"
								onclick={() => linkProvider('google')}
								disabled={oauthLoading === 'google'}
							>
								<Link size={12} />
								Lier
							</button>
						{/if}
					</div>

					<!-- Apple -->
					<div class="flex items-center justify-between">
						<div class="flex items-center gap-2">
							<svg class="w-4 h-4" viewBox="0 0 24 24" fill="currentColor">
								<path d="M17.05 20.28c-.98.95-2.05.88-3.08.4-1.09-.5-2.08-.48-3.24 0-1.44.62-2.2.44-3.06-.4C2.79 15.25 3.51 7.59 9.05 7.31c1.35.07 2.29.74 3.08.8 1.18-.24 2.31-.93 3.57-.84 1.51.12 2.65.72 3.4 1.8-3.12 1.87-2.38 5.98.48 7.13-.57 1.5-1.31 2.99-2.54 4.09zM12.03 7.25c-.15-2.23 1.66-4.07 3.74-4.25.29 2.58-2.34 4.5-3.74 4.25z"/>
							</svg>
							<span class="text-sm text-muted-foreground">Apple</span>
							{#if linkedProviders.apple}
								<span class="text-xs font-medium px-2 py-0.5 rounded-full bg-green-100 text-green-700">Lié</span>
							{/if}
						</div>
						{#if linkedProviders.apple}
							<button
								class="flex items-center gap-1 px-2.5 py-1 text-xs rounded-md border border-border text-muted-foreground hover:text-red-600 hover:border-red-300 transition-colors disabled:opacity-50"
								onclick={() => unlinkProvider('apple')}
								disabled={oauthLoading === 'apple'}
							>
								<Unlink size={12} />
								Délier
							</button>
						{:else}
							<button
								class="flex items-center gap-1 px-2.5 py-1 text-xs rounded-md bg-foreground text-background hover:opacity-90 transition-opacity disabled:opacity-50"
								onclick={() => linkProvider('apple')}
								disabled={oauthLoading === 'apple'}
							>
								<Link size={12} />
								Lier
							</button>
						{/if}
					</div>
				</div>
			</div>
		</div>
	</div>
	{/if}
</div>
