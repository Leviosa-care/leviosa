<script lang="ts">
	import { QrCode } from 'lucide-svelte';

	interface Props {
		name: string | undefined;
		hasEvent?: boolean;
		toggleDrawer: () => any;
	}

	// TODO: pass the hasEvent to false by default when done
	let { name, hasEvent = true, toggleDrawer }: Props = $props();
</script>

<div class="content relative container flex">
	<div class="header flex" style="--gap: 0.75rem;">
		<img class="header-avatar" src="https://pbs.twimg.com/media/FA9tIzxUUAUy80c.jpg" alt="denzel" />
		<div class="">
			<p class="salut">Bonjour {name},</p>
			<p class="title">Content de vous revoir !</p>
		</div>
	</div>
	{#if hasEvent}
		<button class="qrcode relative" onclick={toggleDrawer()}>
			<QrCode color="hsl(var(--clr-light-primary))" size={20} />
		</button>
	{/if}
</div>

<style>
	.content {
		padding-top: 2rem;
		justify-content: space-between;
		align-items: center;
	}
	.qrcode {
		background-color: hsl(var(--clr-dark-primary));
		padding: 0.5rem;
		border-radius: 100%;
		z-index: 10;
	}
	@keyframes pulse {
		0% {
			transform: scale(1);
			opacity: 0.75;
		}
		100% {
			transform: scale(1.5);
			opacity: 0.2;
		}
	}
	/* pulsing animation */
	.qrcode::before {
		content: '';
		position: absolute;
		top: 0;
		bottom: 0;
		left: 0;
		right: 0;
		border-radius: 50%;
		background-color: hsl(var(--clr-dark-primary));

		/* TODO: change the linear for something more sophisticated */
		animation: pulse 1250ms linear 3;
		z-index: -1;
	}
	.header {
		background-color: hsl(var(--clr-light-primary));
		align-items: center;
	}
	.header-avatar {
		width: 40px;
		aspect-ratio: 1;
		border-radius: 100%;
	}
	.salut,
	.title {
		line-height: 1.2rem;
	}
	.salut {
		color: hsl(var(--clr-grey-600));
		font-size: var(--fs-1);
		font-weight: 600;
	}
	.title {
		text-wrap: nowrap;
	}
</style>
