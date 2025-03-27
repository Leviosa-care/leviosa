<script lang="ts">
	// TODO: do the export for the partner thing
	import { Slack, Globe, Instagram, Facebook } from 'lucide-svelte';
	let { partner } = $props();
</script>

<div class="partner flow">
	<div class="partner-header container flex">
		<div class="company-icon">
			<Slack />
		</div>
		<div>
			<p style="font-size: 1.1rem;">{partner.name}</p>
			<p class="dark-ternary-content" style="font-size: 0.9rem">
				{partner.description}
			</p>
		</div>
	</div>
	<div class="partner-icons container flex">
		<a class="partner-icons-container flex" href={partner.website}>
			<Globe />
			<p>site internet</p>
		</a>
		<a class="partner-icons-container flex" href={partner.instagram}>
			<Instagram />
			<p>instagram</p>
		</a>
		<a class="partner-icons-container flex" href={partner.faceboook}>
			<Facebook />
			<p>facebook</p>
		</a>
	</div>
</div>

<style>
	.partner-header {
		gap: 1rem;
		align-items: flex-start;
		align-items: center;
	}
	.company-icon {
		padding: 0.5rem;
		border-radius: 0.5rem;
		background-color: hsl(var(--clr-light-secondary));
		padding-block: 0.5rem;
	}
	.partner {
		padding-top: 0.5rem;
		padding-bottom: 2rem;
	}
	.partner:not(:last-child) {
		border-bottom: 1px solid hsl(var(--clr-stroke));
	}
	.partner-icons {
		justify-content: space-between;
	}
	.partner-icons-container {
		flex-direction: column;
		align-items: center;
		gap: 0.5rem;
		justify-content: center;
	}
	a {
		text-decoration: none;
		color: inherit;
	}
</style>
