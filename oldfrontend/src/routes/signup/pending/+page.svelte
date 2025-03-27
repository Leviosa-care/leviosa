<script lang="ts">
	import Button from '$lib/components/Button.svelte';
	import { Slack, ArrowRight } from 'lucide-svelte';

	import { goto } from '$app/navigation';
	function handleSubmit() {
		goto('/app');
	}
</script>

<div class="content">
	<div class="container grid" style="--gap: 3rem;">
		<div class="grid" style="--gap: 2rem;">
			<div class="logo_placeholder">
				<Slack size={24} />
			</div>
			<div class="grid" style="--gap: 1rem;">
				<h3 class="title center">Merci de ton inscription</h3>
				<div>
					<p>
						Ton inscription est maintenant en attente de validation par un des membres de l'equipe
						Leviosa Tu recevras un mail et une notification des que ton mail sera disponible
					</p>
				</div>
			</div>
		</div>
		<Button rightIcon={ArrowRight} onClick={handleSubmit} style="primary">Visiter le site</Button>
	</div>
</div>

<style>
	.content {
		display: grid;
		place-content: center;
		min-height: 100vh;
	}
	.logo_placeholder {
		width: 40px;
		margin-inline: auto;
		aspect-ratio: 1;
		background-color: hsl(var(--clr-light-secondary));
		display: grid;
		place-content: center;
		border-radius: 0.5rem;
	}
	.title {
		color: black;
		font-size: 1.1rem;
		font-weight: 500;
	}
	h3 {
		font-weight: 600;
	}
</style>
