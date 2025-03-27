<script>
	import { Slack, X } from 'lucide-svelte';
	let installHeaderVisible = $state(true);
</script>

<div>
	<div class="install-container flex" style:display={installHeaderVisible ? 'visible' : 'none'}>
		<div class="flex" style="flex-shrink: 0;--gap: 0.75rem; align-items: center;">
			<div class="flex" style="--gap: 0.2rem; align-items: center;">
				<button class="close" onclick={() => (installHeaderVisible = false)}>
					<X size={16} />
				</button>
				<div
					class="logo"
					style="height: fit-content; padding: 0.5rem; border-radius: 0.5rem; background-color: hsl(var(--clr-accent));"
				>
					<Slack color="white" />
				</div>
			</div>
			<div>
				<p class="install">Installez l'application</p>
				<p class="access">Acces facile et rapide a Leviosa</p>
			</div>
		</div>
		<button class="use">Utiliser</button>
	</div>
	<div class="separator"></div>
</div>

<style>
	.install-container {
		align-items: center;
		padding: 1rem 0.5rem;
		justify-content: space-between;

		@media only screen and (min-width: 500px) {
			display: none;
			visibility: hidden;
		}
	}
	.close {
		background: transparent;
	}
	.access,
	.install {
		line-height: 1.2;
	}
	.install {
		font-size: var(--fs--1);
		color: hsl(var(--clr-grey-700));
		font-weight: 500;
	}
	.access {
		font-size: var(--fs--2);
		color: hsl(var(--clr-grey-400));
	}
	.use {
		background-color: green;
		color: hsl(var(--clr-light-primary));
		border-radius: 2rem;
		padding: 0.25rem 0.5rem;
		font-size: var(--fs--1);
		font-weight: 500;
	}
	.separator {
		height: 1px;
		width: 100%;
		background-color: hsl(var(--clr-grey-200));
	}
</style>
