<script lang="ts">
	interface Props {
		children?: import('svelte').Snippet;
	}

	let { children }: Props = $props();
	// TODO: page transition
</script>

{@render children?.()}

<style></style>
