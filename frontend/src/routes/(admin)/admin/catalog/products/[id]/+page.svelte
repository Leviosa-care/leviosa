<script lang="ts">
	import { enhance } from "$app/forms";
	import { goto } from "$app/navigation";
	import { ArrowLeft, Package, Clock, MapPin, Calendar, FileText, Tag, DollarSign, X, ImageIcon } from "@lucide/svelte";
	import type { PageData, ActionData } from "./$types";
	import { getToastContext } from "$lib/components/toast";

	let { data, form }: { data: PageData, form: ActionData } = $props();

	const toast = getToastContext();

	// Form state - initialize with product data
	let name = $state(data.product.name);
	let description = $state(data.product.description);
	let categoryId = $state(data.updateProductForm.data.category);
	let duration = $state(data.product.duration);
	let price = $state(data.product.price);
	let status = $state<"published" | "draft" | "archived">(data.product.published);
	let availability = $state<"online" | "in-person" | "hybrid">(data.product.availability);
	let bufferTime = $state(data.product.bufferTime);
	let cancellationHours = $state(data.product.cancellationHours);
	let stripeProductId = $state("");
	let imageUrl = $state(data.product.image);

	// Image upload state
	let imagePreview = $state<string | null>(data.product.image || null);
	let isUploading = $state(false);
	let imageFileInput = $state<HTMLInputElement>();

	function updateProductEnhance() {
		return async ({ result }: { result: import('@sveltejs/kit').ActionResult }) => {
			if (result.type === 'redirect') {
				toast.success('Succès', 'Produit mis à jour avec succès');
				goto(result.location);
			} else if (result.type === 'failure') {
				toast.error('Erreur', (result.data as { error?: string })?.error ?? 'Une erreur est survenue');
			}
		};
	}

	function uploadImageEnhance() {
		return async ({ result }: { result: import('@sveltejs/kit').ActionResult }) => {
			isUploading = false;
			if (result.type === 'redirect') {
				goto(result.location);
			} else if (result.type === 'success' && result.data && typeof result.data === 'object' && 'url' in result.data) {
				imageUrl = (result.data as { url: string }).url;
				imagePreview = imageUrl;
				toast.success('Succès', 'Image téléchargée avec succès');
			} else if (result.type === 'failure') {
				toast.error('Erreur', 'Le téléchargement de l\'image a échoué');
				imagePreview = null;
			}
		};
	}

	const categoriesWithoutDefault = data.categories.filter(c => c.id !== "default");

	// Handle image file selection
	function handleImageSelect(event: Event) {
		const input = event.target as HTMLInputElement;
		const file = input.files?.[0];
		if (file) {
			isUploading = true;
			// Preview immediately
			imagePreview = URL.createObjectURL(file);
		}
	}

	// Handle paste of image URL
	function handleImagePaste(event: ClipboardEvent) {
		const items = event.clipboardData?.items;
		if (items) {
			for (const item of items) {
				if (item.type.indexOf('image') !== -1) {
					const file = item.getAsFile();
					if (file && imageFileInput) {
						const dataTransfer = new DataTransfer();
						dataTransfer.items.add(file);
						imageFileInput.files = dataTransfer.files;
						isUploading = true;
						imagePreview = URL.createObjectURL(file);
						// Trigger form submission
						imageFileInput.form?.requestSubmit();
					}
				}
			}
		}
	}

	function removeImage() {
		imagePreview = null;
		imageUrl = "";
	}
</script>

<svelte:head>
	<title>Modifier {data.product.name} | Admin</title>
</svelte:head>

<div class="px-4 py-8 lg:py-12">
	<div class="mb-8">
		<a
			href="/admin/catalog"
			class="inline-flex items-center gap-2 text-foreground-alt hover:text-foreground transition-colors mb-4"
		>
			<ArrowLeft size={20} />
			<span>Retour au Catalogue</span>
		</a>
		<h1 class="text-3xl lg:text-4xl font-bold mb-2">Modifier {data.product.name}</h1>
		<p class="text-muted-foreground">
			Mettez à jour les informations du produit
		</p>
	</div>

	<!-- Image Upload Form (separate from main form) -->
	{#if !imagePreview}
		<div class="mb-6">
			<form method="POST" action="?/uploadImage" use:enhance={uploadImageEnhance}>
				<div class="bg-background rounded-lg border border-border-card p-6">
					<label class="block text-sm font-medium text-foreground-alt mb-2">
						Image du Produit
					</label>
					<p class="text-xs text-muted-foreground mb-3">
						Téléchargez une image ou collez une URL. Max 10 Mo.
					</p>
					<div
						class="relative w-full aspect-video max-h-96 bg-muted rounded-lg border-2 border-dashed border-border-card flex flex-col items-center justify-center cursor-pointer hover:border-foreground/30 transition-colors"
						onpaste={handleImagePaste}
					>
						<input
							bind:this={imageFileInput}
							type="file"
							name="image"
							accept="image/*"
							onchange={handleImageSelect}
							class="absolute inset-0 w-full h-full opacity-0 cursor-pointer"
							aria-label="Télécharger une image"
						/>
						<div class="flex flex-col items-center gap-3 pointer-events-none">
							<div class="w-16 h-16 rounded-full bg-foreground/10 flex items-center justify-center">
								<ImageIcon size={32} class="text-foreground-alt" />
							</div>
							<div class="text-center">
								<p class="text-sm font-medium text-foreground">
									Cliquez pour télécharger ou collez une image
								</p>
								<p class="text-xs text-muted-foreground mt-1">
									ou glissez-déposez un fichier ici
								</p>
							</div>
							<p class="text-xs text-muted-foreground">
								JPEG, PNG, WebP · max 10 Mo
							</p>
							{#if isUploading}
								<p class="text-xs text-foreground">Téléchargement en cours...</p>
							{/if}
						</div>
					</div>
				</div>
			</form>
		</div>
	{/if}

	<div class="bg-background rounded-lg border border-border-card p-6 lg:p-8">
		<form method="POST" use:enhance={updateProductEnhance} class="space-y-6">
			<!-- Image Display (when uploaded) -->
			{#if imagePreview}
				<div>
					<label class="block text-sm font-medium text-foreground-alt mb-2">
						Image du Produit
					</label>
					<div class="relative w-full aspect-video max-h-96 bg-muted rounded-lg overflow-hidden border border-border-card">
						<img
							src={imagePreview}
							alt="Aperçu de l'image"
							class="w-full h-full object-contain"
						/>
						<button
							type="button"
							onclick={removeImage}
							class="absolute top-3 right-3 p-2 bg-black/60 hover:bg-red-600 text-white rounded-lg transition-colors"
							aria-label="Supprimer l'image"
						>
							<X size={16} />
						</button>
					</div>
					<input type="hidden" name="imageUrl" value={imageUrl} />
				</div>
			{/if}

			<!-- Name -->
			<div>
				<label for="name" class="block text-sm font-medium text-foreground-alt mb-2">
					Nom du Produit <span class="text-red-500">*</span>
				</label>
				<div class="relative">
					<Package
						size={18}
						class="absolute left-3 top-1/2 -translate-y-1/2 text-muted-foreground"
					/>
					<input
						id="name"
						name="name"
						type="text"
						bind:value={name}
						required
						placeholder="Massage Relaxant"
						class="w-full pl-10 pr-4 py-3 border border-border-input rounded-lg focus:outline-none focus:ring-2 focus:ring-foreground focus:border-transparent"
					/>
				</div>
			</div>

			<!-- Description -->
			<div>
				<label for="description" class="block text-sm font-medium text-foreground-alt mb-2">
					Description <span class="text-red-500">*</span>
				</label>
				<div class="relative">
					<FileText
						size={18}
						class="absolute left-3 top-3 text-muted-foreground"
					/>
					<textarea
						id="description"
						name="description"
						bind:value={description}
						required
						rows="4"
						placeholder="Décrivez ce produit..."
						class="w-full pl-10 pr-4 py-3 border border-border-input rounded-lg focus:outline-none focus:ring-2 focus:ring-foreground focus:border-transparent resize-none"
					></textarea>
				</div>
			</div>

			<div class="grid grid-cols-1 md:grid-cols-2 gap-6">
				<!-- Category -->
				<div>
					<label for="categoryId" class="block text-sm font-medium text-foreground-alt mb-2">
						Catégorie <span class="text-red-500">*</span>
					</label>
					<div class="relative">
						<Tag
							size={18}
							class="absolute left-3 top-1/2 -translate-y-1/2 text-muted-foreground"
						/>
						<select
							id="categoryId"
							name="categoryId"
							bind:value={categoryId}
							required
							class="w-full pl-10 pr-4 py-3 border border-border-input rounded-lg focus:outline-none focus:ring-2 focus:ring-foreground focus:border-transparent appearance-none bg-background"
						>
							<option value="">Sélectionner une catégorie</option>
							{#each categoriesWithoutDefault as cat}
								<option value={cat.id}>{cat.name}</option>
							{/each}
						</select>
					</div>
				</div>

				<!-- Duration -->
				<div>
					<label for="duration" class="block text-sm font-medium text-foreground-alt mb-2">
						Durée (minutes) <span class="text-red-500">*</span>
					</label>
					<div class="relative">
						<Clock
							size={18}
							class="absolute left-3 top-1/2 -translate-y-1/2 text-muted-foreground"
						/>
						<input
							id="duration"
							name="duration"
							type="number"
							bind:value={duration}
							required
							min="20"
							step="10"
							placeholder="60"
							class="w-full pl-10 pr-4 py-3 border border-border-input rounded-lg focus:outline-none focus:ring-2 focus:ring-foreground focus:border-transparent"
						/>
					</div>
					<p class="mt-1 text-xs text-muted-foreground">
						Minimum 20 minutes, par incréments de 10
					</p>
				</div>
			</div>

			<div class="grid grid-cols-1 md:grid-cols-2 gap-6">
				<!-- Price -->
				<div>
					<label for="price" class="block text-sm font-medium text-foreground-alt mb-2">
						Prix (€)
					</label>
					<div class="relative">
						<DollarSign
							size={18}
							class="absolute left-3 top-1/2 -translate-y-1/2 text-muted-foreground"
						/>
						<input
							id="price"
							name="price"
							type="text"
							value={price}
							readonly
							tabindex="-1"
							class="w-full pl-10 pr-4 py-3 border border-border-input rounded-lg bg-muted text-muted-foreground cursor-not-allowed"
						/>
						<p class="mt-1 text-xs text-muted-foreground">
							Prix en lecture seule — géré via Stripe
						</p>
					</div>
				</div>
				<!-- Stripe Product ID -->
				<div>
					<label for="stripeProductId" class="block text-sm font-medium text-foreground-alt mb-2">
						Stripe Product ID (Optionnel)
					</label>
					<div class="relative">
						<DollarSign
							size={18}
							class="absolute left-3 top-1/2 -translate-y-1/2 text-muted-foreground"
						/>
						<input
							id="stripeProductId"
							name="stripeProductId"
							type="text"
							bind:value={stripeProductId}
							placeholder="prod_xxxxxxxxx"
							class="w-full pl-10 pr-4 py-3 border border-border-input rounded-lg focus:outline-none focus:ring-2 focus:ring-foreground focus:border-transparent"
						/>
					</div>
					<p class="mt-1 text-xs text-muted-foreground">
						ID du produit Stripe pour les paiements
					</p>
				</div>
			</div>

			<!-- Availability -->
			<div>
				<label for="availability" class="block text-sm font-medium text-foreground-alt mb-2">
					Disponibilité <span class="text-red-500">*</span>
				</label>
				<div class="relative">
					<MapPin
						size={18}
						class="absolute left-3 top-1/2 -translate-y-1/2 text-muted-foreground"
					/>
					<select
						id="availability"
						name="availability"
						bind:value={availability}
						required
						class="w-full pl-10 pr-4 py-3 border border-border-input rounded-lg focus:outline-none focus:ring-2 focus:ring-foreground focus:border-transparent appearance-none bg-background"
					>
						<option value="online">En ligne</option>
						<option value="in-person">En présentiel</option>
						<option value="hybrid">Hybride</option>
					</select>
				</div>
			</div>

			<div class="grid grid-cols-1 md:grid-cols-3 gap-6">
				<!-- Buffer Time -->
				<div>
					<label for="bufferTime" class="block text-sm font-medium text-foreground-alt mb-2">
						Tampon (minutes)
					</label>
					<div class="relative">
						<Clock
							size={18}
							class="absolute left-3 top-1/2 -translate-y-1/2 text-muted-foreground"
						/>
						<input
							id="bufferTime"
							name="bufferTime"
							type="number"
							bind:value={bufferTime}
							min="0"
							placeholder="15"
							class="w-full pl-10 pr-4 py-3 border border-border-input rounded-lg focus:outline-none focus:ring-2 focus:ring-foreground focus:border-transparent"
						/>
					</div>
					<p class="mt-1 text-xs text-muted-foreground">
						Temps de pause entre les réservations
					</p>
				</div>

				<!-- Cancellation Hours -->
				<div>
					<label for="cancellationHours" class="block text-sm font-medium text-foreground-alt mb-2">
						Annulation (heures)
					</label>
					<div class="relative">
						<Calendar
							size={18}
							class="absolute left-3 top-1/2 -translate-y-1/2 text-muted-foreground"
						/>
						<input
							id="cancellationHours"
							name="cancellationHours"
							type="number"
							bind:value={cancellationHours}
							min="0"
							placeholder="24"
							class="w-full pl-10 pr-4 py-3 border border-border-input rounded-lg focus:outline-none focus:ring-2 focus:ring-foreground focus:border-transparent"
						/>
					</div>
					<p class="mt-1 text-xs text-muted-foreground">
						Délai d'annulation avant le service
					</p>
				</div>

				<!-- Status -->
				<div>
					<label for="status" class="block text-sm font-medium text-foreground-alt mb-2">
						Statut <span class="text-red-500">*</span>
					</label>
					<div class="relative">
						<Tag
							size={18}
							class="absolute left-3 top-1/2 -translate-y-1/2 text-muted-foreground"
						/>
						<select
							id="status"
							name="status"
							bind:value={status}
							required
							class="w-full pl-10 pr-4 py-3 border border-border-input rounded-lg focus:outline-none focus:ring-2 focus:ring-foreground focus:border-transparent appearance-none bg-background"
						>
							<option value="draft">Brouillon</option>
							<option value="published">Publié</option>
							<option value="archived">Archivé</option>
						</select>
					</div>
				</div>
			</div>

			<!-- Submit Buttons -->
			<div class="flex flex-col sm:flex-row gap-3 pt-4 border-t border-border-card">
				<button
					type="submit"
					class="flex-1 px-6 py-3 bg-foreground text-background rounded-lg hover:opacity-90 focus:outline-none focus:ring-2 focus:ring-foreground focus:ring-offset-2 transition-colors font-medium"
				>
					Mettre à jour le Produit
				</button>
				<a
					href="/admin/catalog"
					class="flex-1 px-6 py-3 bg-background border border-border-input text-foreground-alt rounded-lg hover:bg-muted focus:outline-none focus:ring-2 focus:ring-foreground focus:ring-offset-2 transition-colors font-medium text-center"
				>
					Annuler
				</a>
			</div>
		</form>
	</div>
</div>
