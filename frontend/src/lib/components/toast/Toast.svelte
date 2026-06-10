<script lang="ts">
	import { fade, fly } from 'svelte/transition';
	import { CheckCircle, AlertCircle, Info, AlertTriangle, X } from '@lucide/svelte';
	import type { Toast, ToastType } from './types';

	type Props = {
		toast: Toast;
		onClose: (id: number) => void;
	};

	let { toast, onClose }: Props = $props();

	const icons: Record<ToastType, typeof CheckCircle> = {
		success: CheckCircle,
		error: AlertCircle,
		info: Info,
		warning: AlertTriangle
	};

	const colors: Record<ToastType, string> = {
		success: 'bg-green-50 border-green-200 dark:bg-green-900/30 dark:border-green-800 text-green-800 dark:text-green-200',
		error: 'bg-red-50 border-red-200 dark:bg-red-900/30 dark:border-red-800 text-red-800 dark:text-red-200',
		info: 'bg-blue-50 border-blue-200 dark:bg-blue-900/30 dark:border-blue-800 text-blue-800 dark:text-blue-200',
		warning: 'bg-yellow-50 border-yellow-200 dark:bg-yellow-900/30 dark:border-yellow-800 text-yellow-800 dark:text-yellow-200'
	};

	const IconComponent = icons[toast.type];
	const colorClass = colors[toast.type];
</script>

<div
	class="pointer-events-auto flex items-center gap-3 p-4 rounded-lg border shadow-popover {colorClass}"
	transition:fly|fade={{ y: -20, duration: 300 }}
>
	<IconComponent size={20} class="flex-shrink-0" />
	<div class="flex-1">
		<p class="font-semibold text-sm">{toast.title}</p>
		{#if toast.message}
			<p class="text-sm opacity-90 mt-0.5">{toast.message}</p>
		{/if}
	</div>
	<button
		type="button"
		onclick={() => onClose(toast.id)}
		class="flex-shrink-0 opacity-70 hover:opacity-100 transition-opacity"
	>
		<X size={16} />
	</button>
</div>
