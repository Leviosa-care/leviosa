<script lang="ts">
	import { Send } from 'lucide-svelte';
	import BackButton from '$lib/components/navigation/BackButton.svelte';

	import type { Message } from '$lib/types';
	interface Props {
		messages: Message[];
	}

	let { messages }: Props = $props();
	let messageContent: string = $state('');

	let freelancer = $derived(messages.find((message) => message.author != 'Toi')?.author);
	// TODO: make the page scroll all the way down upon opening, so that I get to the last message directly.
	// TODO: I can do more complex work by grouping the message by date to display the date and the time but I do not care for now
	let wrapperWidth: number = $state(360);
	const inputPadding = '0.75rem';
</script>

<div class="wrapper" bind:clientWidth={wrapperWidth}>
	<div class="underlay"></div>
	<div class="header container flex">
		<BackButton />
		<div class="header-content flex">
			<div class="profile flex">
				<img
					class="img"
					src="https://pbs.twimg.com/media/FA9tIzxUUAUy80c.jpg"
					alt="placeholder avatar"
				/>
				<div class="profile-text flex" style="--gap: 0rem;">
					<p class="author dark-primary-content">{freelancer}</p>
					<p class="team">Equipe massage</p>
				</div>
			</div>
			<button
				onclick={() => console.log('open the modal or the drawer brother')}
				class="cta-profile">Voir profil</button
			>
		</div>
	</div>
	<div class="messages grid" style="--gap: 1.5rem;">
		{#each messages as message}
			{@const { author, content } = message}
			<div class="message-content" style:margin-left={author === freelancer ? '0' : 'auto'}>
				<div class="flex" style="justify-content: space-between;">
					<p class="author">{author}</p>
					<p class="date">Mardi 12 Juillet</p>
				</div>
				<p class:you={author != freelancer} class="message">
					{content}
				</p>
			</div>
		{/each}
	</div>
	<form
		class="footer container flex"
		style:width={`calc(${wrapperWidth}px - 2 * ${inputPadding})`}
		style:transform={`translateX(${inputPadding})`}
	>
		<span class="footer-textarea" role="textbox" contenteditable bind:innerText={messageContent}
		></span>
		<input type="hidden" name="message" bind:value={messageContent} />
		<button class="send-button" type="submit">
			<Send color="hsl(var(--clr-light-primary))" size={20} />
		</button>
	</form>
</div>

<style>
	.wrapper {
		min-height: 100vh;
		position: relative;
		padding-bottom: 5rem;
	}
	.message-content {
		max-width: 400px;
		width: min(75%, 400px);
	}
	.profile {
		align-items: center;
	}
	.team {
		font-size: var(--fs--1);
	}
	.author {
		font-weight: 500;
	}
	.date {
		font-size: var(--fs--1);
	}
	.message {
		padding: 1rem;
		border-radius: 0.5rem;
		width: fit-content;
		width: 100%;
		background-color: hsl(var(--clr-light-primary));
		text-wrap: pretty;
	}
	.you {
		color: hsl(var(--clr-light-primary));
		background-color: hsl(var(--clr-dark-primary));
	}
	.underlay {
		position: absolute;
		top: 0;
		bottom: 0;
		left: 0;
		right: 0;
		background-color: hsl(var(--clr-light-secondary));
		z-index: -1;
	}
	.header {
		background-color: hsl(var(--clr-light-primary));
		padding-block: 1rem;
		align-items: center;
		box-shadow: rgba(17, 17, 26, 0.1) 0px 1px 0px;
		position: sticky;
		top: 0;
		min-width: 100%;
	}
	.header-content {
		justify-content: space-between;
		align-items: center;
		flex-grow: 1;
	}
	.messages {
		margin-top: 2rem;
		background-color: hsl(var(--clr-light-secondary));
		padding-inline: clamp(0.75rem, 2vw, 1.5rem);
	}
	.img {
		--img-dimension: 2.5rem;
		width: var(--img-dimension);
		height: var(--img-dimension);
		border-radius: 100%;
	}
	.profile-text {
		line-height: 1.2;
		flex-direction: column;
		/* align-items: center; */
	}
	.cta-profile {
		/* background-color: hsl(var(--clr-dark-primary)); */
		/* background-color: hsl(var(--clr-accent)); */
		/* color: hsl(var(--clr-light-primary)); */

		color: hsl(var(--clr-dark-primary));
		background: transparent;
		border: 1px solid hsl(var(--clr-dark-primary));

		padding: 0.4rem 0.8rem;
		border-radius: 0.5rem;
		/* TODO: make that thing a little bigger with the vw */
		font-size: clamp(0.8rem, 1vw, 1rem);
		font-weight: 500;
	}
	.cta-profile:is(:hover, :focus) {
		opacity: 0.9;
	}
	/* the position fixed should take into account the wrapper to be placed right ? */
	.footer {
		background-color: hsl(var(--clr-light-primary));
		padding-block: clamp(0.25rem, 1vw, 0.75rem);
		padding-right: 0.5rem;
		width: calc(100% - 0.75rem);
		border-radius: 2rem;
		box-shadow: rgba(0, 0, 0, 0.05) 0px 0px 0px 1px;

		justify-content: space-between;
		align-items: center;

		position: fixed;
		bottom: 0.5rem;

		/* --minmax: 300px; */

		.footer-textarea {
			padding: 0.5rem 1rem;
			display: block;
			width: 100%;
			/* max-width: var(--minmax); */
			overflow: hidden;
			resize: both;
			min-height: 40px;
			line-height: 20px;
		}
		.footer-textarea:is(:hover, :focus) {
			border: none;
			outline: none;
		}
		.footer-textarea[contenteditable]:empty::before {
			content: 'Tom message';
			color: hsl(var(--clr-grey-400));
		}
		.send-button {
			background-color: hsl(var(--clr-dark-primary));
			padding: clamp(0.5rem, 0.5vw, 0.75rem);
			border-radius: 100%;
			aspect-ratio: 1;
			display: grid;
			place-content: center;
		}
	}
</style>
