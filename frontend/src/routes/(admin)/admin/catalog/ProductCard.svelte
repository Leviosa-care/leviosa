<script lang="ts">
	import { Clock, Calendar, MapPin, SquarePen, Trash2, Archive } from "@lucide/svelte";
	import { type SuperValidated } from "sveltekit-superforms";
	import { superForm } from "sveltekit-superforms";
	import ProductModal from "./ProductModal.svelte";
	import AlertDialog from "./AlertDialog.svelte";
	import { type CardType, type Category } from "./products";
	import type { DeleteProduct, type product } from "./schemas";

	type Props = {
		card: CardType;
		statuses: Set<string>;
		categories: Category[];
		availabilities: Set<string>;
		deleteProductForm: SuperValidated<DeleteProduct>;
		updateProductForm: SuperValidated<product>;
	};

	let {
		card,
		statuses,
		categories,
		availabilities,
		deleteProductForm,
		updateProductForm,
	}: Props = $props();

	let { id, name, price, category, description, duration, image, published } =
		card;

	const {
		form: deleteForm,
		enhance: deleteEnhance,
	} = superForm(deleteProductForm, {
		resetForm: true,
	});

	function formatMinutes(totalMinutes: number) {
		const hours = Math.floor(totalMinutes / 60);
		const minutes = totalMinutes % 60;

		let result = "";
		if (hours > 0) {
			result += `${hours}h`;
		}
		if (minutes > 0 || hours === 0) {
			result += `${minutes}min`;
		}

		return result;
	}

	const getStatusBadge = (status: string) => {
		switch (status) {
			case "published":
				return "bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-400";
			case "draft":
				return "bg-yellow-100 text-yellow-700 dark:bg-yellow-900/30 dark:text-yellow-400";
			case "archived":
				return "bg-gray-100 text-gray-700 dark:bg-gray-800 dark:text-gray-400";
			default:
				return "bg-gray-100 text-gray-700 dark:bg-gray-800 dark:text-gray-400";
		}
	};

	const getAvailabilityBadge = (availability: string) => {
		switch (availability) {
			case "online":
				return "bg-blue-100 text-blue-700 dark:bg-blue-900/30 dark:text-blue-400";
			case "in-person":
				return "bg-purple-100 text-purple-700 dark:bg-purple-900/30 dark:text-purple-400";
			case "hybrid":
				return "bg-orange-100 text-orange-700 dark:bg-orange-900/30 dark:text-orange-400";
			default:
				return "bg-gray-100 text-gray-700 dark:bg-gray-800 dark:text-gray-400";
		}
	};

	$effect(() => {
		$deleteForm.id = card.id;
	});
</script>

<div
	{id}
	class="bg-card border border-border-card rounded-lg overflow-hidden hover:shadow-md transition-shadow"
>
	<div class="aspect-video bg-muted relative">
		<img src={image} alt={name} class="w-full h-full object-cover" />
		<div class="absolute top-3 right-3 flex gap-2">
			<span
				class="inline-flex items-center px-2.5 py-0.5 text-xs font-medium rounded-full {getStatusBadge(
					published
				)}"
			>
				{published}
			</span>
		</div>
	</div>

	<div class="p-5">
		<div class="flex items-start justify-between mb-3">
			<h3 class="font-semibold text-foreground line-clamp-1">{name}</h3>
			<span class="text-lg font-bold text-primary">{price}€</span>
		</div>

		<p class="text-sm text-foreground-alt line-clamp-2 mb-3">{description}</p>

		<div class="flex flex-wrap gap-2 mb-4">
			<span
				class="inline-flex items-center gap-1 px-2 py-1 text-xs font-medium rounded-md bg-muted text-foreground-alt"
			>
				<Clock size={12} />
				{formatMinutes(duration)}
			</span>
			<span
				class="inline-flex items-center gap-1 px-2 py-1 text-xs font-medium rounded-md bg-muted text-foreground-alt"
			>
				<MapPin size={12} />
				{card.availability}
			</span>
		</div>

		<div class="flex items-center justify-between pt-3 border-t border-border-card">
			<ProductModal
				{statuses}
				{categories}
				{availabilities}
				{card}
				modalForm={updateProductForm}
			>
				<button
					type="button"
					class="inline-flex items-center gap-1.5 px-3 py-1.5 text-sm font-medium text-foreground-alt hover:text-foreground hover:bg-muted rounded-lg transition-colors"
				>
					<SquarePen size={14} />
					<span>Modifier</span>
				</button>
			</ProductModal>

			<div class="flex gap-1">
				<button
					type="button"
					class="p-1.5 text-muted-foreground hover:text-foreground hover:bg-muted rounded-lg transition-colors"
					aria-label="Archiver"
				>
					<Archive size={14} />
				</button>
				<AlertDialog>
					<form
						method="POST"
						action="?/deleteProduct"
						use:deleteEnhance
						class="contents"
					>
						<input type="hidden" name="id" value={$deleteForm.id} />
						<button
							type="submit"
							class="p-1.5 text-muted-foreground hover:text-red-600 hover:bg-red-50 dark:hover:bg-red-950 rounded-lg transition-colors"
							aria-label="Supprimer"
						>
							<Trash2 size={14} />
						</button>
					</form>
				</AlertDialog>
			</div>
		</div>
	</div>
</div>
