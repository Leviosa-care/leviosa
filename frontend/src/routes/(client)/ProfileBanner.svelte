<script lang="ts">
	import { X, UserCog } from '@lucide/svelte';

	interface Props {
		show: boolean;
	}

	let { show }: Props = $props();
	let dismissed = $state(false);
	let visible = $derived(show && !dismissed);
</script>

{#if visible}
	<div class="bg-amber-50 border-b border-amber-200">
		<div class="max-w-5xl mx-auto px-4 py-3 flex items-center justify-between gap-4">
			<div class="flex items-center gap-3 min-w-0">
				<UserCog size={18} class="text-amber-600 flex-shrink-0" />
				<p class="text-sm text-amber-800">
					Complétez votre profil pour une meilleure expérience.
					<a href="/client/profile" class="font-semibold underline hover:text-amber-900">Compléter maintenant</a>
				</p>
			</div>
			<button
				onclick={() => (dismissed = true)}
				class="flex-shrink-0 p-1 rounded-md text-amber-600 hover:bg-amber-100 transition-colors"
				aria-label="Masquer"
			>
				<X size={16} />
			</button>
		</div>
	</div>
{/if}
