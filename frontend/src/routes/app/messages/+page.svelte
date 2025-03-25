<script lang="ts">
	// TODO: add a conversation suggestion when the user has no conversation like on instagram
	import { MESSAGE_STATES } from '$lib/types';
	import { messageState } from '$lib/stores/persisted_stores.svelte';

	import Message from './Message.svelte';
	import Note from './Note.svelte';
	import NoMessage from './NoMessage.svelte';
	import NoNote from './NoNote.svelte';

	function newMessage() {
		console.log('send a new message brother');
	}
	import type { PageData } from './$types';
	interface Props {
		data: PageData;
	}
	let { data }: Props = $props();
	const { conversations, notes } = data;
</script>

<div class="content" style="--flow-space: 3rem;">
	<div class="underlay"></div>
	{#if $messageState === MESSAGE_STATES.NotesDeSeances}
		{#if notes.length === 0}
			<NoNote />
		{:else}
			<div class="notes">
				{#each notes as note}
					<Note {note} />
				{/each}
			</div>
		{/if}
	{:else if notes.length === 0}
		<NoMessage />
	{:else}
		<div class="messages grid" style="--gap: 1px;">
			{#each conversations as conversation}
				<Message
					id={conversation.id}
					author={conversation.author}
					date={conversation.date}
					time={conversation.time}
					content={conversation.content}
				/>
			{/each}
		</div>
	{/if}
</div>

<style>
	.underlay {
		position: absolute;
		top: 0;
		bottom: 0;
		left: 0;
		right: 0;
		z-index: -1;
	}

	.messages,
	.notes {
		padding-bottom: 7rem;
	}
	.notes {
		padding-top: 1rem;
	}
</style>
