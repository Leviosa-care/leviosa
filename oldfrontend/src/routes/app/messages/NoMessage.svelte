<script>
	import { MessageSquareText } from 'lucide-svelte';
</script>

<div class="empty-content center container grid" style="--gap: 2rem;" role="region">
	<div class="icon">
		<MessageSquareText size={36} color="hsl(var(--clr-grey-700))" aria-label="Icon message" />
	</div>
	<div class="grid" style="--gap: 0.5rem;">
		<p class="no-message" aria-label="Titre aucun message">Vous n'avez aucun message</p>
		<p class="no-message-body" aria-label="Description aucun message">
			Lorsque vous recevrez un nouveau message, il apparaitra par ici.
		</p>
	</div>
</div>

<style>
	.empty-content {
		padding-top: 4rem;
		text-wrap: balance;
	}
	.no-message {
		color: hsl(var(--clr-grey-600));
		font-weight: 600;
	}
	.no-message-body {
		color: hsl(var(--clr-grey-400));
		font-size: var(--fs--1);
	}
	.icon {
		width: fit-content;
		margin-inline: auto;
	}
</style>
