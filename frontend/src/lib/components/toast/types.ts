import { getContext } from 'svelte';

const TOAST_KEY = Symbol('toast');

export type ToastType = 'success' | 'error' | 'info' | 'warning';

export type Toast = {
	id: number;
	type: ToastType;
	title: string;
	message?: string;
	duration?: number;
};

export type ToastContext = {
	toasts: Array<Toast>;
	success: (title: string, message?: string, duration?: number) => void;
	error: (title: string, message?: string, duration?: number) => void;
	info: (title: string, message?: string, duration?: number) => void;
	warning: (title: string, message?: string, duration?: number) => void;
};

export function getToastContext(): ToastContext {
	return getContext(TOAST_KEY);
}
