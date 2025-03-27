<script lang="ts">
	let name: string = 'prestataires';

	import type { Prestation } from '$lib/types';
	interface Props {
		prestations: Prestation[];
	}

	let { prestations }: Props = $props();

	let focusedButton: number = $state(0);
	function toggleChecked(e: MouseEvent) {
		const target = e.currentTarget as HTMLButtonElement;
		focusedButton = parseInt(target.id);
		const input = target.querySelector("input[type='radio']") as HTMLInputElement;
		if (input) input.click();
	}
	// TODO: change the icon when I have the right one or change the whole component altogether
</script>

<div class="components flex" style="--gap: 0.5rem;">
	{#each prestations as prestation, index}
		<button
			id={String(index)}
			type="button"
			onclick={toggleChecked}
			class="option stroke flex"
			style="background-image: url({prestation.imageUrl});"
		>
			<input class="radioinput" checked={index === 0} type="radio" {name} id={prestation.id} />
			<div class="content flex" style="--gap:-0.1rem;">
				{#if focusedButton !== index}
					<div class="overlay"></div>
				{/if}
				<div class="blured-bg" style:height={focusedButton === index ? '35%' : '100%'}>
					<label for={prestation.id} class="name">
						{prestation.text}
					</label>
				</div>
			</div>
		</button>
	{/each}
</div>

<style>
	.option:has(:global(input[type='radio']:checked)) {
		border: 2px solid hsl(var(--clr-dark-primary));
	}
	.radioinput[type='radio'] {
		visibility: hidden;
		display: none;
	}
	.components {
		border-radius: 0.5rem;
		align-items: center;
		width: 100%;
		margin-inline: auto;
	}
	.overlay {
		position: absolute;
		background-color: rgba(0, 0, 0, 0.75);
		top: 0;
		left: 0;
		height: 100%;
		width: 100%;
	}
	.option {
		padding: 0.5rem;
		position: relative;
		width: 100%;
		background: transparent;
		align-items: center;
		justify-content: center;
		border-radius: 0.5rem;
		border: 2px solid hsl(var(--clr-light-primary));
		max-width: 120px;
		color: hsl(var(--clr-dark-secondary));
		/* NOTE:  for the background  */
		background-size: cover;
		background-position: center;
		background-repeat: no-repeat;
		position: relative;
		overflow: hidden;
		height: 120px;
		color: hsl(var(--clr-light-primary));
		line-height: 0.95;
		font-weight: 500;
	}
	.blured-bg {
		position: absolute;
		bottom: 0;
		height: 35%;
		width: 100%;
		backdrop-filter: blur(15px);
		display: grid;
		place-content: center;
		padding: 0.5rem;
		transition: height 0.3s ease;
	}
	.content {
		flex-direction: column;
		align-items: center;
		gap: 0.75rem;
	}
	.name {
		text-align: center;
		font-size: var(--fs--1);
	}
</style>
