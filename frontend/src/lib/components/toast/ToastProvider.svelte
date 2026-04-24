<script lang="ts">
	import { setContext } from 'svelte';
	import ToastContainer from './ToastContainer.svelte';
	import type { Toast, ToastType } from './types';
	import type { ToastContext } from './types';

	const TOAST_KEY = Symbol('toast');

	let toasts = $state<Array<Toast>>([]);
	let toastId = 0;

	function addToast(type: ToastType, title: string, message?: string, duration = 4000) {
		const id = toastId++;
		const toast: Toast = { id, type, title, message, duration };
		toasts.push(toast);

		if (duration > 0) {
			setTimeout(() => {
				removeToast(id);
			}, duration);
		}

		return id;
	}

	function removeToast(id: number) {
		toasts = toasts.filter((t) => t.id !== id);
	}

	const context: ToastContext = {
		get toasts() {
			return toasts;
		},
		success: (title, message, duration) => addToast('success', title, message, duration),
		error: (title, message, duration) => addToast('error', title, message, duration),
		info: (title, message, duration) => addToast('info', title, message, duration),
		warning: (title, message, duration) => addToast('warning', title, message, duration)
	};

	setContext(TOAST_KEY, context);
</script>

<ToastContainer />
