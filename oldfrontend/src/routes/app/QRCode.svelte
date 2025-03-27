<script lang="ts">
	import Button from '$lib/components/Button.svelte';
	function addToAppleWallet() {
		console.log('yes les gars');
	}
	import { Wallet } from 'lucide-svelte';
	interface Props {
		qrcode: string;
	}

	let { qrcode }: Props = $props();
	// TODO: add the apple wallet icon like in Luma for the same page
</script>

<div class="content container flex">
	<!-- {#await data.qrcode} -->
	{#await qrcode}
		<p>loading the qrcode</p>
	{:then qrcode}
		{#if qrcode !== ''}
			<div class="center">
				<p class="title">Ton ticket</p>
				<p class="subtitle">Ticket a presenter a l'accueil</p>
			</div>
			<img class="qrcode" src={qrcode} alt="qrcode" />
			<Button leftIcon={Wallet} onClick={addToAppleWallet}>Ajoute a Apple Wallet</Button>
		{:else}
			<div>
				<p>There is not qrcode brother</p>
			</div>
		{/if}
	{/await}
</div>

<style>
	.content {
		flex-direction: column;
		justify-content: space-between;
		padding-top: 1rem;
	}
	.title {
		color: hsl(var(--clr-grey-700));
		font-size: var(--fs-1);
		font-weight: 600;
	}
	.qrcode {
		border: 1px dashed hsl(var(--clr-stroke));
		display: block;
	}
</style>
