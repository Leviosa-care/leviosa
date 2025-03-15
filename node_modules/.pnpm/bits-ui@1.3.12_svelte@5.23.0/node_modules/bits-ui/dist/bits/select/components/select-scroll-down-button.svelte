<script lang="ts">
	import { box, mergeProps } from "svelte-toolbelt";
	import type { SelectScrollDownButtonProps } from "../types.js";
	import { useSelectScrollDownButton } from "../select.svelte.js";
	import { useId } from "../../../internal/use-id.js";
	import { Mounted } from "../../utilities/index.js";

	let {
		id = useId(),
		ref = $bindable(null),
		child,
		children,
		...restProps
	}: SelectScrollDownButtonProps = $props();

	const scrollButtonState = useSelectScrollDownButton({
		id: box.with(() => id),
		ref: box.with(
			() => ref,
			(v) => (ref = v)
		),
	});

	const mergedProps = $derived(mergeProps(restProps, scrollButtonState.props));
</script>

{#if scrollButtonState.canScrollDown}
	<Mounted bind:mounted={scrollButtonState.state.mounted} />
	{#if child}
		{@render child({ props: restProps })}
	{:else}
		<div {...mergedProps}>
			{@render children?.()}
		</div>
	{/if}
{/if}
