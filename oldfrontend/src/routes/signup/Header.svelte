<script lang="ts">
	import BackButton from '$lib/components/navigation/BackButton.svelte';
	import PaginationIndicator from '$lib/components/ui/PaginationIndicator.svelte';
	interface Props {
		selectedElement: number;
		pathname?: string;
	}

	let { selectedElement, pathname = '' }: Props = $props();
</script>

<div class="header">
	<BackButton pathname={pathname !== '' ? pathname : '/'} />
	<PaginationIndicator count={4} {selectedElement} />
	<div class="fake-arrow"></div>
</div>

<style>
	.header {
		padding-block: 2rem;
		display: flex;
		justify-content: space-between;
		align-content: center;
	}
	.fake-arrow {
		width: 32px;
		aspect-ratio: 1;
		visibility: hidden;
	}
</style>
