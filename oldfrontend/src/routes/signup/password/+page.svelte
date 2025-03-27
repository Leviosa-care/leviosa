<script>
	import PasswordInput from '$lib/components/forms/atoms/PasswordInput.svelte';
	import Button from '$lib/components/Button.svelte';
	import { ArrowRight } from 'lucide-svelte';
	import Header from '../Header.svelte';
	import { goto } from '$app/navigation';
	function handleSubmit() {
		// TODO: store that data in some store, so that I can use it whenever I want, whatever route I am...
		goto('verify-password');
	}
</script>

<div class="container grid" style="--gap: 2rem;">
	<Header selectedElement={2} pathname="address" />
	<h2>Ton mot de passe</h2>
	<form action="" class="grid" style="--gap: 3rem;">
		<div class="grid" style="--gap:3rem;">
			<h3 class="title">Mot de passe</h3>
			<div class="grid" style="--gap: 2rem;">
				<PasswordInput id="password1" />
				<PasswordInput id="password2" label="Confirmation de mot de passe" />
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
	h3 {
		font-weight: 600;
	}
</style>
