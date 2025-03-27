<script lang="ts">
	import { X } from 'lucide-svelte';
	// TODO: make sure that space does not close the drawer
	// TODO: store the messages in a store or in localstorage so that the user do not loose everything
	// TODO: how to send the message to the contact ?
	type Recipient = {
		imgUrl: string;
		firstname: string;
		lastname: string;
	};

	interface Props {
		// TODO: that part is supposed to be found through filtering the list of potential contact possible for a user
		recipients?: Recipient[];
	}

	let {
		recipients = [
			{
				imgUrl: 'https://pbs.twimg.com/media/FA9tIzxUUAUy80c.jpg',
				firstname: 'Livio',
				lastname: 'HENRY'
			},
			{
				imgUrl: 'https://pbs.twimg.com/media/FA9tIzxUUAUy80c.jpg',
				firstname: 'Gary',
				lastname: 'HENRY'
			}
		]
	}: Props = $props();
	let modalIsOpen = $state(false);
	function toggleModalState() {
		modalIsOpen = !modalIsOpen;
		console.log('modalIsOpen is now:', modalIsOpen);
	}
	function updateRecipient() {
		console.log('need to update the recipient list');
	}
	let value: string = $state('');
	function cleanInput() {
		value = '';
	}
</script>

<div class="modal">
	<div class="bar"></div>
	<div class="header flex">
		<p class="title">Nouvelle conversation</p>
	</div>
	<div class="to flex">
		<p>A</p>
		<input
			onfocus={toggleModalState}
			onblur={toggleModalState}
			onchange={updateRecipient}
			bind:value
			type="text"
		/>
		<button aria-label="Ferme le modal" type="button" onclick={cleanInput} class="close">
			<X size={12} />
		</button>
		<div class="recipients grid" class:reveal={modalIsOpen} style="--gap: 0rem;">
			{#each recipients as recipient}
				<button class="recipient flex" style="--gap: 0.5rem;">
					<div class="avatar">
						<img src={recipient.imgUrl} alt="profile denzel" />
					</div>
					<p>{recipient.firstname} {recipient.lastname}</p>
				</button>
			{/each}
		</div>
	</div>
	<div class="message">
		<textarea name="" id=""></textarea>
	</div>
</div>

<style>
	.modal {
		min-height: 50px;
	}
	.header {
		font-size: var(--fs--1);
		align-items: center;
		padding-bottom: 0.5rem;
		margin-inline: auto;
		justify-content: center;
	}
	.avatar {
		width: 32px;
		aspect-ratio: 1;
		border-radius: 100%;
		overflow: hidden;
	}
	.title {
		font-weight: 600;
	}
	.header,
	.to,
	.to .recipient {
		border-bottom: 1px solid hsl(var(--clr-stroke));
	}
	.to {
		padding: 0.5rem;
		padding-right: 0;
		position: relative;
		align-items: center;
	}
	.to .recipients {
		position: absolute;
		bottom: 0;
		transform: translateY(100%);
		left: 0;
		right: 0;
		width: 100%;
	}
	.recipients.reveal > .recipient {
		display: flex;
		visibility: visible;
	}
	.to .recipient {
		padding: 0.75rem 1rem;
		background: transparent;
		align-items: center;
		display: none;
		visibility: hidden;
	}
	.to input {
		border: 1x solid red;
	}
	.to input:is(:global(:focus, :hover)) {
		outline: none;
	}
	.close {
		border-radius: 100%;
		padding: 0.5rem;
	}
	.message {
		min-height: 30vh;
	}
	/* TODO: use a text area or some sort of input ? */
	.message textarea {
		padding: 1rem;
		border: none;
		width: 100%;
		height: 300px;
	}
	.message textarea:is(:global(:hover, :focus)) {
		outline: none;
	}
</style>
