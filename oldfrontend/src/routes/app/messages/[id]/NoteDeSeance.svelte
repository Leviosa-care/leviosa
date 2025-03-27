<script lang="ts">
	// TODO: put all the note de seance content in here

	import type { SessionNote } from '$lib/types';
	interface Props {
		notes: SessionNote[];
	}
	let { notes }: Props = $props();
</script>

<div class="content">
	{#each notes as note}
		<div>je teste du texte les gars</div>
		<div>{note}</div>
	{/each}
</div>

<style></style>
