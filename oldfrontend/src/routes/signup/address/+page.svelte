<script>
	import FormInput from '$lib/components/forms/FormInput.svelte';
	import Button from '$lib/components/Button.svelte';
	import { ArrowLeftIcon, ArrowRight } from 'lucide-svelte';
	import PaginationIndicator from '$lib/components/ui/PaginationIndicator.svelte';
	import Header from '../Header.svelte';
	// TODO: handle spacing like in the general page
	import { goto } from '$app/navigation';
	function handleSubmit() {
		// TODO: store that data in some store, so that I can use it whenever I want, whatever route I am...
		goto('password');
	}
</script>

<div class="container grid" style="--gap: 2rem;">
	<Header selectedElement={1} pathname="general" />
	<h2>Informations de domiciliation</h2>
	<form action="" class="grid" style="--gap: 4rem;">
		<div class="flow" style="--gap: 4rem;">
			<div class="grid" style="--gap: 3rem;">
				<h3 class="title">Adresse</h3>
				<div class="grid" style="--gap: 2rem;">
					<FormInput
						id="address1"
						name="address1"
						type="text"
						placeholder="Entre ton addresse"
						value=""
						label="Adresse 1"
					/>
					<FormInput
						id="address2"
						name="address2"
						type="text"
						placeholder="Entre ton adresse"
						value=""
						label="Adresse 2"
					/>
				</div>
				<div class="grid" style="--gap: 2rem;">
					<h3 class="title">Code postal et ville</h3>
					<div class="two-column">
						<FormInput
							id="postalcode"
							name="postalcode"
							type="text"
							placeholder="Entre ton code postal"
							value=""
							label="Code postal"
						/>
						<FormInput
							id="city"
							name="city"
							type="text"
							placeholder="Entre ta ville"
							value=""
							label="Ville"
						/>
					</div>
				</div>
			</div>
		</div>
		<Button rightIcon={ArrowRight} onClick={handleSubmit}>Continue ton inscription</Button>
	</form>
</div>

<style>
	h2 {
		color: black;
		font-weight: 500;
		font-size: 1.2rem;
		font-weight: 800;
	}
	h3 {
		font-weight: 600;
	}
	.title {
		color: black;
		font-size: 1.1rem;
		font-weight: 500;
	}
	form {
		margin-top: 2rem;
	}
	.two-column {
		display: flex;
		align-content: center;
		gap: 1rem;
	}
</style>
