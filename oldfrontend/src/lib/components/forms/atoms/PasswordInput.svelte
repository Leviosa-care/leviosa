<script lang="ts">
	// TODO: remove that when the app goes to production
	let value: string = '123soleil';
	import { Eye, EyeOff } from 'lucide-svelte';
	interface Props {
		// TODO: change the type so that it can be a enum
		id: string;
		label?: string;
	}

	let { id, label = 'Mot de passe' }: Props = $props();
	let visible = $state(false);
</script>

<div class="wrapper">
	<label for={id} class="label">{label}</label>
	<input
		{id}
		class="stroke"
		name="password"
		type={visible ? 'text' : 'password'}
		placeholder="Entre ton mot de passe"
		{value}
		required
	/>
	<button type="button" class="icon" onclick={() => (visible = !visible)}>
		{#if visible}
			<EyeOff strokeWidth={1} absoluteStrokeWidth={true} />
		{:else}
			<Eye strokeWidth={1} absoluteStrokeWidth={true} />
		{/if}
	</button>
</div>

<style>
	/* TODO: the thing with that icon that I need to work on brother */
	.icon {
		position: absolute;
		border-radius: 0.5rem;
		top: 0;
		right: 0;
		display: flex;
		align-items: center;
		justify-content: center;
		height: 100%;
		background: transparent;
		padding-inline: var(--padding-left);
	}
	.wrapper {
		position: relative;
		/* max-width: 100vw; */
		--padding-left: calc(3 * var(--fs-0) / 4);
	}

	.label {
		--paddinginline: 0.2rem;
		position: absolute;
		top: calc((-1 * var(--fs-0) / 2) - 0.3rem);
		left: calc(1rem - var(--paddinginline));
		left: calc(var(--padding-left) - var(--paddinginline));
		background-color: white;
		padding-inline: var(--paddinginline);
	}

	input {
		width: 100%;
		border-radius: 0.5rem;
		padding: var(--padding-left);
	}

	/* input:not(:placeholder-shown):valid { outline: 1px solid hsl(var(--clr-success)); } */
	input:not(:placeholder-shown):valid {
		outline: 1px solid hsl(var(--clr-dark-ternary));
	}
	input:not(:placeholder-shown):invalid {
		outline: 1px solid hsl(var(--clr-error));
	}
	input:focus:invalid {
		outline: 1px solid yellow;
	}

	/* ::placeholder, .label { color: hsl(var(--clr-stroke)); } */
	::placeholder,
	.label {
		color: hsl(var(--clr-dark-ternary) / 0.8);
	}
</style>
