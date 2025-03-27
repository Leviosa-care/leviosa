<script lang="ts">
	import BackButton from '$lib/components/navigation/BackButton.svelte';
	let scrollY: number = $state(0);
	// let hasShadow = $state(false);
	let hasShadow = $derived(scrollY > 0);
	// function handleScroll() {
	//     hasShadow = window.scrollY > 0;
	// }
	import Field from './Field.svelte';
	const fields = [
		{
			title: 'Connection et securite',
			subtitle: 'Mot de passe et historique de conversation',
			pathname: 'parameters/login-and-security'
		},
		{
			title: 'Notifications',
			subtitle: 'Comment nous vous contactons',
			pathname: 'parameters/notifications'
		}
	];
	import { createHorizontalSwipeHandler } from '$lib/scripts/swipe';
	function swipeBack(direction: 'left' | 'right') {
		const backButton = document.querySelector('.back') as HTMLButtonElement;
		if (direction === 'right') {
			backButton.click();
		}
	}
	const swipeBackAction = createHorizontalSwipeHandler(swipeBack);
</script>

<!-- <svelte:window bind:scrollY onscroll={handleScroll} /> -->
<svelte:window bind:scrollY />

<div class="content grid" use:swipeBackAction.action>
	<!-- <div class:shadow={hasShadow} class="container back-container"> -->
	<div class:hasShadow class="back-container container">
		<div class="back">
			<BackButton />
		</div>
	</div>
	<div class="header container flex">
		<h3 class="page-title">Parametres du compte</h3>
		<div class="ghost"></div>
	</div>
	<div class="fields container">
		{#each fields as field, index}
			<Field
				title={field.title}
				subtitle={field.subtitle}
				pathname="parameters/login-and-security"
			/>
			{#if index != fields.length - 1}
				<div class="separator"></div>
			{/if}
		{/each}
	</div>
</div>

<style>
	.content {
		padding-bottom: 1rem;
		position: relative;
	}
	.back-container {
		padding-block: 1rem;
		background-color: hsl(var(--clr-light-primary));
		position: sticky;
		top: 0;
		left: 0;
		right: 0;
	}
	/* .shadow { */
	/* 	box-shadow: rgba(28, 24, 24, 0.16) 0px -2px 8px; */
	/* } */
	.back {
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
