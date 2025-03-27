<script lang="ts">
	import { goto } from '$app/navigation';
	interface Props {
		id: string;
		author: string;
		date: string;
		time: string;
		content: string;
	}

	let { id, author, date, time, content }: Props = $props();
	function getConversationPath(id: string) {
		return `/app/messages/${id}`;
	}
</script>

<button
	onclick={() => goto(getConversationPath(id))}
	class="message-content flex"
	style="--gap: var(--fs--1);"
	role="link"
	aria-label={`Aller à la conversation avec ${author}`}
>
	<img class="img" src="https://pbs.twimg.com/media/FA9tIzxUUAUy80c.jpg" alt="avatar du profil" />
	<div class="message-text" role="region" aria-label="Contenu du message">
		<div class="header flex">
			<p class="author dark-primary-content" aria-label="Author">
				{author}
			</p>
			<p class="time" aria-label="Timestamp">
				{date} - {time}
			</p>
		</div>
		<p class="message-body" aria-label="Texte du message">
			{content}
		</p>
	</div>
</button>

<style>
	.message-content,
	.header {
		align-items: center;
	}
	.message-content {
		--p-0: clamp(0.5rem, 1.5vw + 0.1rem, 1rem);
		--p-1: clamp(1rem, 2vw + 0.2rem, 2rem);
		/* padding: 1rem 0.5rem; */
		padding: var(--p-1) var(--p-0);
		background-color: hsl(var(--clr-light-primary));
		border-bottom: 1px solid hsl(var(--clr-grey-200));
	}
	.message-text {
		flex-grow: 1;
	}
	.header {
		justify-content: space-between;
	}
	.img {
		background-color: hsl(var(--clr-light-secondary));
		/* TODO: make that some css variable ? */
		width: clamp(2.5rem, 7vw + 1rem, 3rem);
		aspect-ratio: 1;
		border-radius: 100%;
		flex-shrink: 0;
	}
	.author {
		font-weight: 500;
	}
	.time {
		font-size: var(--fs--1);
		color: hsl(var(--clr-grey-500));
	}
	.message-body {
		display: -webkit-box;
		-webkit-box-orient: vertical;
		-webkit-line-clamp: 1;
		line-clamp: 1;
		overflow: hidden;
		text-align: left;
		font-size: var(--fs--1);
		color: hsl(var(--clr-grey-500));
	}
</style>
