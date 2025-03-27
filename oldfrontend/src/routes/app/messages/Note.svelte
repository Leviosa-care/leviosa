<script lang="ts">
	import type { SessionNote } from '$lib/types';

	import { redirectTo } from '$lib/scripts/redirect';
	import { Check } from 'lucide-svelte';
	interface Props {
		note: SessionNote;
	}

	let { note }: Props = $props();
</script>

<button class="note flex" onclick={() => redirectTo(`notes/${note.id}`)}>
	<div class="grid" style="--gap: 0.5rem;">
		<p class="note-title">Session {note.offer}</p>
		<div class="status flex" style="--gap: 0.5rem;">
			<div class="check">
				<Check size={12} />
			</div>
			<p>Consultation realisee</p>
		</div>
		<div class="text flex" style="--gap: 0.5rem;">
			<img class="avatar" src="https://pbs.twimg.com/media/FA9tIzxUUAUy80c.jpg" alt="avatar" />
			<div>
				<p>avec {note.author}</p>
			</div>
		</div>
	</div>
	<p>{note.date}</p>
</button>

<style>
	.note {
		/* width: 100%; */
		margin-inline: auto;
		background: transparent;
		border-radius: 1rem;
		padding: 0.5rem 1rem;
		padding: 1rem;
		text-align: left;
		color: hsl(var(--clr-grey-400));
		justify-content: space-between;
		box-shadow:
			rgba(0, 0, 0, 0.05) 0px 6px 24px 0px,
			rgba(0, 0, 0, 0.08) 0px 0px 0px 1px;
	}
	.note-title {
		color: hsl(var(--clr-grey-600));
		font-size: var(--fs-1);
	}
	.text {
		align-items: center;
	}
	.avatar {
		width: 2rem;
		height: fit-content;
		background-color: hsl(var(--clr-dark-primary));
		border-radius: 100%;
		flex: none;
		/* flex-grow: 0; */
	}
	.status {
		align-items: center;
	}
	.check {
		border-radius: 100%;
		border: 1px solid hsl(var(--clr-success));
		padding: 0.2rem;
		color: hsl(var(--clr-success));
	}
</style>
