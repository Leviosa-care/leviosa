<script lang="ts">
	import BackButton from '$lib/components/navigation/BackButton.svelte';
	import Field from './Field.svelte';

	import type { PageData } from './$types';
	interface Props {
		data: PageData;
	}
	let { data }: Props = $props();
	const { fields } = data;

	let active: boolean = $state(false);
	let activeFieldIndex: number = $state(-1);
	function setActiveState(index: number) {
		active = index !== -1;
		activeFieldIndex = index;
	}

	let hasShadow = $state(false);

	function handleScroll() {
		hasShadow = window.scrollY > 0;
	}
	import { createHorizontalSwipeHandler } from '$lib/scripts/swipe';

	function swipeBack(direction: 'left' | 'right') {
		const backButton = document.querySelector('.back') as HTMLButtonElement;
		if (direction === 'right') {
			backButton.click();
		}
	}
	const swipeBackAction = createHorizontalSwipeHandler(swipeBack);
</script>

<svelte:window onscroll={handleScroll} />

<div class="content grid" use:swipeBackAction.action>
	<div class:shadow={hasShadow} class="back-header container">
		<div class="back-container">
			<BackButton />
		</div>
	</div>
	<div class="header container flex">
		<h3 class="page-title">Informations personnelles</h3>
		<div class="ghost"></div>
	</div>
	<div class="fields container">
		{#each fields as field, index}
			<Field
				{active}
				isActive={active && activeFieldIndex === index}
				activeFn={() => setActiveState(index)}
				inactiveFn={() => setActiveState(-1)}
				name={field.fieldname}
				value={field.value}
				missingLabel={field.missingLabel}
				modifyLabel={field.modifyLabel}
				modifiedSlot={field.modifiedSlot}
				addLabel={field.addLabel}
				properties={field.properties}
			/>
		{/each}
	</div>
</div>

<style>
	.content {
		padding-bottom: 1rem;
		position: relative;
	}
	.back-header {
		padding-block: 1rem;
		background-color: hsl(var(--clr-light-primary));
		position: sticky;
		top: 0;
		left: 0;
		right: 0;
	}
	.shadow {
		box-shadow: rgba(28, 24, 24, 0.16) 0px -2px 8px;
	}
	.back-container {
		width: fit-content;
		height: fit-content;
		display: grid;
		padding: 0.75rem;
		place-content: center;
		border-radius: 100%;
		background-color: #f7f7f9;
	}
	.header {
		align-items: center;
		padding-block: 1rem;
		justify-content: space-between;
	}
	.ghost {
		height: 40px;
		width: 40px;
	}
</style>
