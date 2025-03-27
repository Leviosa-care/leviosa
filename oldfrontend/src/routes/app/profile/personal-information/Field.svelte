<script lang="ts">
	// todo: each field should be a form
	import Button from '$lib/components/Button.svelte';
	import type { Component } from 'svelte';

	type Props = {
		name: string;
		value: string | number;
		missingLabel: string | undefined;
		activeFn: () => void;
		inactiveFn: () => void;
		isActive: boolean;
		active: boolean;
		modifyLabel: string;
		modifiedSlot: Component | undefined;
		addLabel: string | undefined;
		properties: any;
	};
	let {
		name,
		value,
		missingLabel,
		activeFn,
		inactiveFn,
		isActive,
		active,
		modifyLabel,
		modifiedSlot: ModifiedSlot,
		addLabel,
		properties
	}: Props = $props();

	let action: 'ajouter' | 'modifier' | 'annuler' = $state(value !== '' ? 'modifier' : 'ajouter');

	function handleClick(e: MouseEvent) {
		const button = e.currentTarget as HTMLButtonElement;
		switch (button.innerText) {
			case 'ajouter':
			case 'modifier':
				activeFn();
				action = 'annuler';
				break;
			case 'annuler':
				inactiveFn();
				if (value === '' || value === 0) {
					action = 'ajouter';
				} else {
					action = 'modifier';
				}
				break;
		}
	}
	function sendData() {
		console.log('here send data to server');
	}

	let disabled = $derived(active && (action === 'ajouter' || action === 'modifier'));
</script>

<form class="content grid" style="--gap:1.5rem;" style:opacity={disabled ? 0.2 : 1}>
	<div class="grid" style="--gap:0.3725rem">
		<div class="header flex">
			<p class="name">{name}</p>
			<button {disabled} type="button" onclick={handleClick} class="action">{action}</button>
		</div>
		{#if isActive}
			<p class="value">
				{value === '' || value === 0 ? addLabel : modifyLabel}
			</p>
		{:else}
			<p class="value">
				{value === '' || value === 0 ? missingLabel : value}
			</p>
		{/if}
	</div>
	{#if isActive}
		{#if value === ''}
			<p>I need to add the thing</p>
		{:else}
			<ModifiedSlot {properties} />
		{/if}
		<Button onClick={sendData} buttonType="submit">Enregistrer</Button>
	{/if}
</form>

<style>
	.content {
		padding-block: 1rem;
		border-bottom: 1px solid hsl(var(--clr-grey-100));
	}
	.header {
		justify-content: space-between;
		align-items: center;
	}
	.action {
		text-decoration: underline;
		color: hsl(var(--clr-grey-700));
		font-weight: 500;
		background: transparent;
	}
	.name {
		color: hsl(var(--clr-grey-700));
		font-weight: 400;
	}
	.value {
		max-width: 75%;
		line-height: 1.1;
		font-size: var(--fs--1);
	}
</style>
