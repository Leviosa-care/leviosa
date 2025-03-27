<script lang="ts">
	import TextInput from './atoms/TextInput.svelte';
	import PasswordInput from './atoms/PasswordInput.svelte';
	interface Props {
		placeholder?: string | undefined;
		label?: string;
		id?: string;
		name?: string;
		type?: string;
		// TODO: remove that when the app goes to production
		value?: string;
	}

	let {
		placeholder = 'Entre ton mot de passe',
		label = 'some label',
		id = 'some id',
		name = 'some name',
		type = 'password',
		value = '123soleil'
	}: Props = $props();
</script>

<div class="wrapper">
	{#if type === 'password'}
		<PasswordInput {id} />
	{:else}
		<TextInput {id} {label} {name} {placeholder} {value} />
	{/if}
</div>

<style>
	.wrapper {
		position: relative;
		--padding-left: calc(3 * var(--fs-0) / 4);
	}
</style>
