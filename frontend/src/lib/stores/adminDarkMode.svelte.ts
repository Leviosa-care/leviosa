/**
 * Admin-only dark mode store
 * Persists preference to localStorage
 */

import { browser } from '$app/environment';

const STORAGE_KEY = 'admin-dark-mode';

// Initialize from localStorage or default to false (light mode)
const initialValue = browser
	? localStorage.getItem(STORAGE_KEY) === 'true'
	: false;

// Create a reactive store
let value = $state(initialValue);

export const adminDarkMode = {
	get(): boolean {
		return value;
	},
	set(newValue: boolean) {
		value = newValue;
		if (browser) {
			localStorage.setItem(STORAGE_KEY, String(newValue));
		}
	},
	toggle() {
		value = !value;
		if (browser) {
			localStorage.setItem(STORAGE_KEY, String(value));
		}
	}
};
