import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vitest/config';
import type { Plugin } from 'vite';
import tailwindcss from '@tailwindcss/vite'
import { kitRoutes } from 'vite-plugin-kit-routes';

const add_browser_onmount: Plugin = {
    name: 'vite-plugin-onmount',
    config(config) {
        if (process.env.VITEST) {
            if (!config.resolve?.conditions) {
                if (!config.resolve) {
                    config.resolve = {};
                }
                config.resolve.conditions = [];
            }
            config.resolve.conditions.unshift('browser');
        }
    }
};

export default defineConfig({
    plugins: [
        sveltekit(),
        add_browser_onmount,
        kitRoutes(),
        tailwindcss(),
    ],
    test: {
        include: ['src/**/*.{test,spec}.{js,ts}', 'tests/**/*.{test,spec}.{js,ts}'],
        environment: 'jsdom',
        setupFiles: ['src/vitest.setup.ts'],
        globals: true
    },
});
