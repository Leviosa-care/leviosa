<script lang="ts">
	import type { PageProps } from './$types';
	import { enhance } from '$app/forms';
	import { CheckCircle, ArrowLeft, Save } from '@lucide/svelte';

	let { data, form }: PageProps = $props();

	// Use updated user from form response if available, otherwise use initial data
	let user = $derived(
		form?.success && form.user ? form.user : data.user
	);

	let gender = $state(user?.gender ?? '');
	let birthdate = $state(user?.birthdate ? user.birthdate.split('T')[0] : '');
	let address1 = $state(user?.address1 ?? '');
	let address2 = $state(user?.address2 ?? '');
	let postalCode = $state(user?.postal_code ?? '');
	let city = $state(user?.city ?? '');
</script>

<svelte:head>
	<title>Mon profil | Leviosa</title>
</svelte:head>

<div class="space-y-8">
	<!-- Header -->
	<div>
		<div class="flex items-center gap-3 mb-1">
			<a href="/client" class="text-muted-foreground hover:text-foreground transition-colors">
				<ArrowLeft size={20} />
			</a>
			<h1 class="text-2xl lg:text-3xl font-bold text-foreground">Mon profil</h1>
		</div>
		<p class="text-muted-foreground mt-1">
			{#if user?.profile_incomplete}
				<span class="text-amber-600 font-medium">Profil incomplet</span> — Complétez vos informations personnelles
			{:else}
				Gérez vos informations personnelles
			{/if}
		</p>
	</div>

	{#if form?.success}
		<div class="bg-green-50 border border-green-200 rounded-lg p-4 flex items-center gap-3">
			<CheckCircle size={18} class="text-green-600 flex-shrink-0" />
			<p class="text-sm text-green-800">Profil mis à jour avec succès.</p>
		</div>
	{/if}

	{#if form?.error}
		<div class="bg-red-50 border border-red-200 rounded-lg p-4 flex items-center gap-3">
			<p class="text-sm text-red-700">{form.error}</p>
		</div>
	{/if}

	{#if user}
		<form method="POST" use:enhance class="space-y-8">
			<!-- Personal info -->
			<div class="bg-card rounded-lg border border-border-card p-6 shadow-card space-y-6">
				<h2 class="text-lg font-semibold text-foreground">Informations personnelles</h2>

				<!-- Read-only fields -->
				<div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
					<div>
						<label class="block text-sm font-medium text-muted-foreground mb-1">Prénom</label>
						<p class="text-foreground">{user.first_name || '—'}</p>
					</div>
					<div>
						<label class="block text-sm font-medium text-muted-foreground mb-1">Nom</label>
						<p class="text-foreground">{user.last_name || '—'}</p>
					</div>
					<div>
						<label class="block text-sm font-medium text-muted-foreground mb-1">Email</label>
						<p class="text-foreground">{user.email || '—'}</p>
					</div>
					<div>
						<label class="block text-sm font-medium text-muted-foreground mb-1">Téléphone</label>
						<p class="text-foreground">{user.telephone || '—'}</p>
					</div>
				</div>

				<div class="border-t border-border-card pt-6 space-y-4">
					<!-- Gender -->
					<div>
						<label for="gender" class="block text-sm font-medium text-foreground mb-1.5">
							Genre
							{#if !gender}
								<span class="text-amber-600 text-xs">(requis)</span>
							{/if}
						</label>
						<select
							name="gender"
							id="gender"
							bind:value={gender}
							class="w-full px-3 py-2.5 rounded-lg border border-border-card bg-background text-foreground text-sm focus:outline-none focus:ring-2 focus:ring-foreground/20"
						>
							<option value="">Sélectionnez votre genre</option>
							<option value="man">Homme</option>
							<option value="woman">Femme</option>
							<option value="non_binary">Non binaire</option>
							<option value="prefer_not_to_say">Je préfère ne pas le dire</option>
							<option value="custom">Je préfère décrire mon genre</option>
						</select>
					</div>

					<!-- Birthdate -->
					<div>
						<label for="birthdate" class="block text-sm font-medium text-foreground mb-1.5">
							Date de naissance
							{#if !birthdate}
								<span class="text-amber-600 text-xs">(requis)</span>
							{/if}
						</label>
						<input
							type="date"
							name="birthdate"
							id="birthdate"
							bind:value={birthdate}
							class="w-full px-3 py-2.5 rounded-lg border border-border-card bg-background text-foreground text-sm focus:outline-none focus:ring-2 focus:ring-foreground/20"
						/>
					</div>
				</div>
			</div>

			<!-- Address -->
			<div class="bg-card rounded-lg border border-border-card p-6 shadow-card space-y-6">
				<h2 class="text-lg font-semibold text-foreground">
					Adresse
					{#if !address1}
						<span class="text-amber-600 text-sm font-normal">(requis)</span>
					{/if}
				</h2>

				<div class="space-y-4">
					<div>
						<label for="address1" class="block text-sm font-medium text-foreground mb-1.5">Adresse</label>
						<input
							type="text"
							name="address1"
							id="address1"
							bind:value={address1}
							placeholder="Numéro et nom de rue"
							class="w-full px-3 py-2.5 rounded-lg border border-border-card bg-background text-foreground text-sm focus:outline-none focus:ring-2 focus:ring-foreground/20"
						/>
					</div>
					<div>
						<label for="address2" class="block text-sm font-medium text-foreground mb-1.5">Complément d'adresse</label>
						<input
							type="text"
							name="address2"
							id="address2"
							bind:value={address2}
							placeholder="Bâtiment, appartement, etc."
							class="w-full px-3 py-2.5 rounded-lg border border-border-card bg-background text-foreground text-sm focus:outline-none focus:ring-2 focus:ring-foreground/20"
						/>
					</div>
					<div class="grid grid-cols-2 gap-4">
						<div>
							<label for="postalCode" class="block text-sm font-medium text-foreground mb-1.5">Code postal</label>
							<input
								type="text"
								name="postalCode"
								id="postalCode"
								bind:value={postalCode}
								placeholder="75000"
								class="w-full px-3 py-2.5 rounded-lg border border-border-card bg-background text-foreground text-sm focus:outline-none focus:ring-2 focus:ring-foreground/20"
							/>
						</div>
						<div>
							<label for="city" class="block text-sm font-medium text-foreground mb-1.5">Ville</label>
							<input
								type="text"
								name="city"
								id="city"
								bind:value={city}
								placeholder="Paris"
								class="w-full px-3 py-2.5 rounded-lg border border-border-card bg-background text-foreground text-sm focus:outline-none focus:ring-2 focus:ring-foreground/20"
							/>
						</div>
					</div>
				</div>
			</div>

			<!-- Submit -->
			<div class="flex justify-end">
				<button
					type="submit"
					class="inline-flex items-center gap-2 px-6 py-2.5 rounded-lg bg-foreground text-background text-sm font-medium hover:bg-foreground/90 transition-colors disabled:opacity-50"
				>
					<Save size={16} />
					Enregistrer
				</button>
			</div>
		</form>
	{:else}
		<div class="bg-card rounded-lg border border-border-card p-8 text-center text-muted-foreground shadow-card">
			Impossible de charger vos informations. Veuillez réessayer.
		</div>
	{/if}
</div>
