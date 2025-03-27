<script lang="ts">
	let options = [
		{ id: 'homme', text: 'Monsieur', icon: '♂' },
		{ id: 'femme', text: 'Madame', icon: '♀' },
		{ id: 'autre', text: 'Autre', icon: '○' }
	];
	let name: string = 'gender';

	function toggleChecked(e: MouseEvent) {
		const target = e.currentTarget as HTMLButtonElement;
		const input = target.querySelector("input[type='radio']") as HTMLInputElement;
		if (input) input.click();
	}
</script>

<div class="components flex">
	{#each options as option, index}
		<button type="button" onclick={toggleChecked} class="option stroke flex">
			<input class="radioinput" checked={index === 0} type="radio" {name} id={option.id} />
			<div class="content flex" style="--gap:-0.1rem;">
				<div class="icon">
					{option.icon}
				</div>
				<label for={option.id} style="font-size:0.9rem;" class="text">
					{option.text}
				</label>
			</div>
		</button>
	{/each}
</div>

<style>
	.option:has(:global(input[type='radio']:checked)) {
		border: 1px solid hsl(var(--clr-dark-primary));
		color: hsl(var(--clr-dark-primary));
	}
	.radioinput[type='radio'] {
		visibility: hidden;
		display: none;
	}
	.components {
		border-radius: 0.5rem;
		align-items: center;
	}
	.option {
		padding: 0.5rem;
		position: relative;
		width: 100%;
		background: transparent;
		align-items: center;
		justify-content: center;
		border-radius: 0.5rem;
		max-width: 120px;
		color: hsl(var(--clr-dark-secondary));
	}
	.content {
		flex-direction: column;
		align-items: center;
	}
	.icon {
		--size: 3rem;
		line-height: var(--size);
		font-size: var(--size);
		width: fit-content;
	}
</style>
