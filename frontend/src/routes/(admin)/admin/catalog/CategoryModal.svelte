<script lang="ts">
	import type { Snippet } from "svelte";
	import { Dialog, Label, Separator, Button } from "bits-ui";
	import { X } from "@lucide/svelte";

	import type { SuperValidated } from "sveltekit-superforms";
	import { superForm } from "sveltekit-superforms";
	import type { category } from "./schemas";

	type Props = {
		children: import("svelte").Snippet;
		modalForm: SuperValidated<category>;
	};

	let { children, modalForm = $bindable() }: Props = $props();

	const {
		form,
		errors,
		enhance,
	} = superForm(modalForm, {
		resetForm: true,
	});
</script>

<Dialog.Root>
	<Dialog.Trigger>
		{@render children()}
	</Dialog.Trigger>
	<Dialog.Portal>
		<Dialog.Overlay
			class="data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 fixed inset-0 z-50 bg-black/80"
		/>
		<Dialog.Content
			class="rounded-card-lg bg-background shadow-popover data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 data-[state=closed]:zoom-out-95 data-[state=open]:zoom-in-95 outline-hidden fixed left-[50%] top-[50%] z-50 w-full max-w-[calc(100%-2rem)] translate-x-[-50%] translate-y-[-50%] border p-8 sm:max-w-[540px] md:w-full"
		>
			<Dialog.Title class="w-full text-xl font-semibold tracking-tight">
				Nouvelle Catégorie
			</Dialog.Title>
			<Dialog.Description class="text-foreground-alt !mt-1 text-sm">
				Créez une nouvelle catégorie pour organiser vos produits
			</Dialog.Description>

			<Separator.Root class="bg-muted mx-5 !mb-2 !mt-5 block h-px" />

			<form
				method="POST"
				action="?/createCategory"
				class="grid gap-4"
				use:enhance
				enctype="multipart/form-data"
			>
				<div>
					<Label.Root for="name" class="text-sm font-medium text-foreground-alt">
						Nom
					</Label.Root>
					<div class="relative w-full">
						<input
							id="name"
							name="name"
							type="text"
							bind:value={$form.name}
							required
							placeholder="Massage"
							class="w-full px-4 py-2.5 border border-border-input rounded-lg focus:outline-none focus:ring-2 focus:ring-foreground focus:border-transparent text-sm"
						/>
					</div>
				</div>

				<div>
					<Label.Root for="description" class="text-sm font-medium text-foreground-alt">
						Description
					</Label.Root>
					<div class="relative w-full">
						<textarea
							id="description"
							name="description"
							bind:value={$form.description}
							rows="4"
							required
							placeholder="Décrivez cette catégorie..."
							class="w-full px-4 py-2.5 border border-border-input rounded-lg focus:outline-none focus:ring-2 focus:ring-foreground focus:border-transparent text-sm resize-none"
						></textarea>
					</div>
				</div>

				<div>
					<Label.Root for="status" class="text-sm font-medium text-foreground-alt">
						Statut
					</Label.Root>
					<div class="relative w-full">
						<select
							id="status"
							name="status"
							bind:value={$form.status}
							class="w-full px-4 py-2.5 border border-border-input rounded-lg focus:outline-none focus:ring-2 focus:ring-foreground focus:border-transparent text-sm appearance-none bg-background"
						>
							<option value="draft">Brouillon</option>
							<option value="published">Publié</option>
							<option value="archived">Archivé</option>
						</select>
					</div>
				</div>

				<div class="flex w-full justify-end pt-2">
					<Button.Root
						type="submit"
						class="cursor-pointer bg-foreground text-background hover:opacity-90 px-6 py-2.5 rounded-lg font-medium"
					>
						Créer la Catégorie
					</Button.Root>
				</div>
			</form>

			<Button.Root
				type="button"
				class="cursor-pointer absolute right-5 top-5"
			>
				<Dialog.Close
					class="focus-visible:ring-foreground focus-visible:ring-offset-background focus-visible:outline-hidden rounded-md focus-visible:ring-2 focus-visible:ring-offset-2 active:scale-[0.98] cursor-pointer"
				>
					<X class="text-foreground size-5" />
					<span class="sr-only">Close</span>
				</Dialog.Close>
			</Button.Root>
		</Dialog.Content>
	</Dialog.Portal>
</Dialog.Root>
