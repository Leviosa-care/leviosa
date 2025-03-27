<script lang="ts">
	interface Props {
		figure: number | string;
		label: string;
	}

	let { figure, label }: Props = $props();
</script>

<div class="content card">
	<p class="figure">{figure}</p>
	<p class="label">{label}</p>
	<div class="card-overlay"></div>
</div>

<style>
	.content {
		color: hsl(var(--clr-grey-100));
	}
	.figure {
		font-weight: 700;
		font-size: var(--fs-3);
	}
	.card {
		--card-border-radius: 1rem;
		--card-padding: 1rem;
		padding: var(--card-padding);
		border-radius: var(--card-border-radius);
		backdrop-filter: blur(15px) saturate(150%);
		color: hsl(var(--clr-light-ternary));
		box-shadow:
			rgba(60, 64, 67, 0.3) 0px 1px 2px 0px,
			rgba(60, 64, 67, 0.15) 0px 1px 3px 1px;
		position: relative;
		overflow: hidden;
	}
	.card-overlay {
		position: absolute;
		top: 0;
		left: 0;
		width: 100%;
		height: 100%;
		background-color: rgba(0, 0, 0, 0.25);
		z-index: -1;
	}
</style>
