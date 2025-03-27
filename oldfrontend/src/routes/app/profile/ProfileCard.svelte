<script>
	import { ChevronRight } from 'lucide-svelte';
	// TODO: make the avatar some color with the letter of the first name of that person if they do not have a photo
	import { goto } from '$app/navigation';
</script>

<button
	onclick={() => goto('profile/user-profile')}
	class="field-left flex"
	style="--gap: 1rem;"
	role="link"
	aria-label="Aller vers les informations de son profil"
>
	<img class="avatar" src="https://pbs.twimg.com/media/FA9tIzxUUAUy80c.jpg" alt="avatar profile" />
	<div>
		<p class="name" aria-label="Nom complet">John DOE</p>
		<p class="cta">Afficher le profil</p>
	</div>
</button>
<ChevronRight color="black" aria-label="Chevron droit" />

<style>
	/* TODO: make that responsive brother */
	.avatar {
		width: 3rem;
		aspect-ratio: 1;
		background-color: hsl(var(--clr-dark-primary));
		border-radius: 100%;
	}
	/* NOTE: the css sizing using the utopia thing */
	.field-left {
		font-size: var(--fs-0);
		background: transparent;
		text-align: left;
	}
	.name {
		color: hsl(var(--clr-grey-700));
		font-weight: 500;
	}
	.cta {
		font-size: var(--fs--1);
	}
</style>
