<script>
	import { MessageSquareText } from 'lucide-svelte';
</script>

<div class="empty-content center container grid" style="--gap: 2rem;">
	<div class="icon">
		<MessageSquareText size={36} color="hsl(var(--clr-grey-700))" />
	</div>
	<div class="grid" style="--gap: 0.5rem;">
		<p class="no-note">Vous n'avez aucune note</p>
		<p class="no-note-body">Lorsque vous recevrez une nouvelle note, elle apparaitra par ici.</p>
	</div>
</div>

<style>
	.empty-content {
		padding-top: 4rem;
		text-wrap: balance;
	}
	.no-note {
		color: hsl(var(--clr-grey-600));
		font-weight: 600;
	}
	.no-note-body {
		color: hsl(var(--clr-grey-400));
		font-size: var(--fs--1);
	}
	.icon {
		width: fit-content;
		margin-inline: auto;
	}
</style>
