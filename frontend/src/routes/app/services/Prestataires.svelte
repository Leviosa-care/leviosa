<script lang="ts">
	import type { Freelancer } from '$lib/types';

	import Drawer from '$lib/components/Drawer.svelte';
	import Profile from './Profile.svelte';

	interface Props {
		freelancers: Freelancer[];
	}

	let { freelancers }: Props = $props();

	let selectedFreelancer: (typeof freelancers)[0] | undefined = $state();
	let isOpen = $state(false);
	// TODO: changes the badges for this thing ?
</script>

<div class="container grid" style="--gap: 1rem;">
	<div class="card grid" style="--gap: 0.5rem;">
		<h3 class="title">Les prestataires</h3>
		<p>
			De nombreux prestataires sont disponibles, decouvre leur profil en cliquant sur un des badges
			ci-dessous
		</p>
		<div class="card-overlay"></div>
	</div>
	<div class="badges flex" style="--gap: 2rem;">
		{#each freelancers as freelancer}
			<button
				onclick={() => {
					isOpen = true;
					selectedFreelancer = freelancer;
				}}
				class="badge flex"
				style="--gap:0.2rem;"
			>
				<img class="photo" src={freelancer.avatar} alt="" />
				<p class="name">{freelancer.firstname}</p>
			</button>
		{/each}
	</div>
</div>

<Drawer bind:isOpen>
	<Profile freelancer={selectedFreelancer} />
</Drawer>

<style>
	.photo {
		border-radius: 100%;
		width: 3rem;
		aspect-ratio: 1;
	}
	/* TODO: for the badges a grid container might be better since it make everything in columns and add rows if necessary */
	.badge {
		background: transparent;
		flex-direction: column;
		align-items: center;
	}
	.name {
		color: hsl(var(--clr-light-secondary));
	}
	.title {
		color: hsl(var(--clr-grey-100));
		font-size: var(--fs-1);
	}
	.card {
		--card-border-radius: 1rem;
		--card-padding: 1rem;
		padding: var(--card-padding);
		border-radius: var(--card-border-radius);
		backdrop-filter: blur(15px) saturate(150%);
		color: hsl(var(--clr-light-ternary));
		box-shadow:
			rgba(60, 64, 67, 0.3) 0px 1px 2px 0px,
			rgba(60, 64, 67, 0.15) 0px 1px 3px 1px;
		position: relative;
		overflow: hidden;
	}
	.card-overlay {
		position: absolute;
		top: 0;
		left: 0;
		width: 100%;
		height: 100%;
		background-color: rgba(0, 0, 0, 0.25);
		z-index: -1;
	}
</style>
