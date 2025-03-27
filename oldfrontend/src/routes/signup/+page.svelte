<script lang="ts">
	import Button from '$lib/components/Button.svelte';
	import FormInput from '$lib/components/forms/FormInput.svelte';
	import { ArrowRight } from 'lucide-svelte';

	const formControls = [
		{
			name: 'email',
			label: 'Email',
			type: 'email',
			placeholder: 'Entrer votre email',
			// TODO: change that value when testing the backend
			// value: 'admin@gmail.fr'
			value: ''
		}
		// {
		//     name: 'password',
		//     label: 'Mot de passe',
		//     type: 'password',
		//     placeholder: 'Entrer votre mot de passe',
		// }
	];
</script>

<div class="flow container" style="--gap: 3rem;">
	<h2 class="fs-h2" style="color: black;">Faisons connaissance</h2>
	<form action="" class="grid" method="POST">
		<div class="flow" style="--flow-space: 2rem;">
			<p class="title">Informations generales</p>
			<div class="flow">
				<div class="two-column">
					<FormInput
						id="name"
						name="name"
						type="text"
						placeholder="Entre ton nom"
						value=""
						label="Nom"
					/>
					<FormInput
						id="firstname"
						name="fistname"
						type="text"
						placeholder="Entre ton prenom"
						value=""
						label="Prenom"
					/>
				</div>
				<div class="two-column">
					<FormInput
						id="gender"
						name="gender"
						type="text"
						placeholder="Entre ton genre"
						value=""
						label="Genre"
					/>
					<FormInput
						id="bday"
						name="bday"
						type="text"
						placeholder="Entre ta date de naissance"
						value=""
						label="Date de naissance"
					/>
				</div>
			</div>
		</div>

		<div class="flow" style="--flow-space: 2rem;">
			<p class="title">Ton addresse</p>
			<div class="flow">
				<FormInput
					id="address1"
					name="address1"
					type="text"
					placeholder="Entre ton addresse"
					value=""
					label="Addresse 1"
				/>
				<FormInput
					id="address2"
					name="address2"
					type="text"
					placeholder="Entre ton addresse"
					value=""
					label="Addresse 2"
				/>
				<div class="two-column">
					<FormInput
						id="postalcode"
						name="postalcode"
						type="text"
						placeholder="Entre ton code postal"
						value=""
						label="Code Postal"
					/>
					<FormInput
						id="ville"
						name="ville"
						type="text"
						placeholder="Entre ta ville"
						value=""
						label="Ville"
					/>
				</div>
			</div>
		</div>
		<div class="flow" style="--flow-space: 2rem">
			<p class="title">Choisis un mot de passe</p>
			<div class="flow">
				<FormInput
					id="mdp"
					name="mdp"
					type="text"
					placeholder="Entre ton mot de passe"
					value=""
					label="Mot de passe"
				/>
				<FormInput
					id="mdp2"
					name="mdp2"
					type="text"
					placeholder="Confirme ton mot de passe"
					value=""
					label="Confirmation de mot de passe"
				/>
			</div>
		</div>
		<Button
			style="primary"
			rightIcon={ArrowRight}
			buttonType="submit"
			onClick={() => console.log('do something with that cta brother')}
		>
			Poursuivre l'inscription
		</Button>
	</form>
</div>

<style>
	.title {
		font-size: 1.2rem;
		font-weight: 600;
		color: black;
	}
	.container {
		padding-top: 2rem;
		padding-inline: 0.75rem;
		min-height: 100vh;
	}

	form {
		margin-top: 2rem;
	}
	.two-column {
		display: flex;
		align-content: center;
		gap: 1rem;
	}
	h2 {
		font-weight: 800;
	}
</style>
