<script lang="ts">
	import { goto } from '$app/navigation';
	import { ArrowRight } from 'lucide-svelte';

	import FormInput from '$lib/components/forms/FormInput.svelte';
	import Button from '$lib/components/Button.svelte';
	import TimePickerButton from '$lib/components/forms/TimePickerButton.svelte';
	import GenderRadio from './GenderRadio.svelte';
	import Drawer from '$lib/components/Drawer.svelte';
	import Header from '../Header.svelte';
	let isOpen = $state(false);

	function handleSubmit() {
		// TODO: store that data in some store, so that I can use it whenever I want, whatever route I am...
		goto('address');
	}
</script>

<div class="container grid" style="--gap: 2rem;">
	<Header selectedElement={0} />
	<h2>Informations personnelles</h2>
	<form action="" class="grid" style="--gap: 3rem;">
		<div class="grid" style="--gap: 2rem;">
			<div class="grid" style="--gap: 1.5rem;">
				<p class="title">Nom et prenom</p>
				<div class="grid" style="--gap: 0.5rem; grid-auto-flow: column;">
					<FormInput id="name" name="name" type="text" placeholder="DUPONT" value="" label="Nom" />
					<FormInput
						id="firstname"
						name="fistname"
						type="text"
						placeholder="Jean"
						value=""
						label="Prenom"
					/>
				</div>
			</div>
			<div class="grid" style="--gap: 1.5rem;">
				<p class="title">Genre</p>
				<GenderRadio />
			</div>
			<div class="grid" style="--gap: 1.5rem;">
				<p class="title">Date de naissance</p>
				<TimePickerButton onClick={() => (isOpen = !isOpen)} />
			</div>
		</div>
		<Button rightIcon={ArrowRight} onClick={handleSubmit}>Continue ton inscription</Button>
	</form>
</div>

<Drawer bind:isOpen>
	<div>soem content that should be in the drawer to help me select the right date guys</div>
	<div>soem content that should be in the drawer to help me select the right date guys</div>
	<div>soem content that should be in the drawer to help me select the right date guys</div>
	<div>soem content that should be in the drawer to help me select the right date guys</div>
</Drawer>

<style>
	h2 {
		color: black;
		font-weight: 500;
		font-size: 1.2rem;
	}
	.title {
		color: black;
		font-size: 1.1rem;
		font-weight: 500;
	}
	form {
		margin-top: 2rem;
	}
	h2 {
		font-weight: 800;
	}
</style>
