<script lang="ts">
	// TODO: change the button type to submit when I do the work with form handling
	import Google from '@assets/google.svelte';
	import Apple from '@assets/apple.svelte';
	import { Mail, LogIn, Slack } from 'lucide-svelte';

	import FormInput from '$lib/components/forms/FormInput.svelte';
	import Button from '$lib/components/Button.svelte';

	import Tabs from '$lib/components/Tabs.svelte';
	type pageType = 'signin' | 'signup';
	let pageState: pageType = $state('signin');

	const handleTab = (): void => {
		if (pageState === 'signin') pageState = 'signup';
		else pageState = 'signin';
	};

	const offers = [
		{ name: 'Se connecter', action: () => handleTab() },
		{ name: "S'enregister", action: () => handleTab() }
	];
	// TODO: add the little floating element to indicate that this a pwa with cta to tutorial

	import { goto } from '$app/navigation';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();
	$inspect(data);
	function getSignPath(): string {
		console.log('the state is:', pageState);
		if (pageState === 'signin') {
			console.log('go to app');
			return '/app';
		} else if (pageState === 'signup') {
			console.log('go to sign up');
			return '/signup/general';
		} else {
			console.log('stay on the same path');
			return '/';
		}
	}
	let { SignInFormControls, SignUpFormControl } = $derived(data);
	let formControls = $derived(pageState === 'signin' ? SignInFormControls : SignUpFormControl);
</script>

<div class="page-container container">
	<!-- <div class="page-container"> -->
	<div class="left-content flow container" style="--flow-space: 2.5rem;">
		<div class="header grid" style="--gap: 2rem;">
			<div class="flex" style="--gap: 0.5rem; align-items: center;">
				<div class="logo_placeholder">
					<Slack size={24} />
				</div>
				<p style="font-size: 1.2rem; font-weight: 500; color; black">Leviosa</p>
			</div>
			<div class="grid" style="--gap: -0.1rem;">
				<p style="color: black; font-size: 1.5rem; font-weight: 600;">Content de te revoir</p>
				<p>
					{pageState === 'signin' ? 'Connecte toi a ton compte.' : 'Cree toi un compte Leviosa'}
				</p>
			</div>
		</div>
		<Tabs {offers} isSecondary={true} />
		<form>
			<div class="flow grid" style="--flow-space: 2rem;">
				<div class="grid" style="--gap: 2rem;">
					{#each formControls as formcontrol}
						<FormInput
							id={formcontrol.name}
							name={formcontrol.name}
							type={formcontrol.type}
							placeholder={formcontrol.placeholder}
							value={formcontrol.value}
							label={formcontrol.label}
						/>
					{/each}
					{#if pageState === 'signin'}
						<Button buttonType="button" leftIcon={LogIn} onClick={() => goto(getSignPath())}>
							Se connecter
						</Button>
					{:else if pageState === 'signup'}
						<Button buttonType="button" leftIcon={Mail} onClick={() => goto(getSignPath())}>
							Continuer avec ton email
						</Button>
					{/if}
				</div>
			</div>
		</form>
		<div class="separator-block">
			<div class="separator"></div>
			<p class="or fs-paragraph stroke">ou</p>
			<div class="separator"></div>
		</div>
		<form class="oauths grid" method="GET">
			<Button
				style="secondary"
				buttonType="submit"
				onClick={() => console.log('do the google oauth')}
				leftIcon={Google}
			>
				{pageState === 'signin' ? 'Se connecter' : "S'enregistrer"} avec Google
			</Button>
			<Button
				style="secondary"
				buttonType="submit"
				onClick={() => console.log('do the apple oauth')}
				leftIcon={Apple}
			>
				{pageState === 'signin' ? 'Se connecter' : "S'enregistrer"} avec Apple
			</Button>
		</form>
	</div>
	<div class="right-content stroke">
		<img
			class="right-img"
			src="https://st.depositphotos.com/2627021/4738/i/950/depositphotos_47383969-stock-photo-spa-in-garden-vertical-composition.jpg"
			alt="some illustration for the page"
		/>
	</div>
</div>

<style>
	.page-container {
		display: grid;
		grid-template-columns: 1fr;
		align-items: center;
		justify-content: center;

		min-height: 100vh;
		height: 100vh;
		/* NOTE: not the best way to deal with that... */
		overflow: hidden;
	}
	.left-content {
		/* outline: 2px solid red; */
		width: 100%;
		background-color: hsl(var(--clr-light-primary));
		padding-bottom: 1rem;
		padding-top: 3rem;
		height: 100%;
	}

	.right-content {
		/* outline: 2px solid blue; */
		width: 100%;
		height: 100%;
		display: none;
		grid-column: 2;
		border-radius: 0.5rem;
		background-color: #bdbdbd;
		border-top-left-radius: 0;
		border-bottom-left-radius: 0;
	}

	@media (min-width: 768px) {
		.page-container {
			grid-template-columns: repeat(2, 1fr);
		}
		.left-content {
			max-width: none;
			width: 100%;
		}
		.right-content {
			display: block;
		}
		.right-img {
			object-fit: cover;
			height: 100%;
			width: 100%;
		}
	}

	.logo_placeholder {
		width: 40px;
		aspect-ratio: 1;
		background-color: hsl(var(--clr-light-secondary));
		display: grid;
		place-content: center;
		border-radius: 0.5rem;
	}

	.separator-block {
		--or-width: 3.6ch;
		display: flex;
		align-items: center;
	}

	.separator {
		width: 100%;
		height: 1px;
		background-color: hsl(var(--clr-stroke));
	}

	.or {
		color: hsl(var(--clr-stroke));
		width: var(--or-width);
		aspect-ratio: 1;
		border-radius: 100%;
		display: flex;
		align-items: center;
		justify-content: center;
		flex-grow: 0;
		flex-shrink: 0;
	}

	::placeholder {
		font-weight: 300;
	}
</style>
