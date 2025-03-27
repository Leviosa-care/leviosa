<script lang="ts">
	import { ChevronRight } from 'lucide-svelte';
	import { goto } from '$app/navigation';
	interface Props {
		title: string;
		subtitle: string;
		pathname: string;
	}

	let { title, subtitle, pathname }: Props = $props();
</script>

<button class="field flex" onclick={() => goto(pathname)}>
	<div>
		<p class="title">{title}</p>
		<p class="subtitle">{subtitle}</p>
	</div>
	<ChevronRight />
</button>

<style>
	.field {
		/* width: 100%; */
		align-items: center;
		justify-content: space-between;
		padding-block: 1.5rem;
		background: transparent;
		text-align: left;
		color: inherit;
	}
	.subtitle {
		line-height: 1.1;
		font-size: var(--fs--1);
	}
	.title {
		color: hsl(var(--clr-grey-700));
		font-weight: 400;
	}
</style>
