<script lang="ts">
	interface Props {
		placeholder?: string | undefined;
		label?: string;
		// TODO: change the type so that it can be a enum
		id?: string;
		name?: string;
		// TODO: remove that when the app goes to production
		value: string | number; // open eye if the password is visible and the opposite otherwise
	}

	let {
		placeholder = 'Entre ton mot de passe',
		label = 'Mot de passe',
		id = 'some id',
		name = 'some name',
		value
	}: Props = $props();
</script>

<div class="wrapper">
	<label for={id} class="label capitalize">{label}</label>
	<input {id} class="stroke" {name} type="text" {placeholder} {value} required />
</div>

<style>
	.wrapper {
		position: relative;
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
