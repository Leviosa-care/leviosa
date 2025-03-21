import { browser } from '$app/environment';
import { writable, get } from 'svelte/store';

import type { NavState } from '$lib/types';
// export type navState = 'ghost' | 'home' | 'messages' | 'ghost' | 'events' | 'profil';

// function createNavbarStore(key: string, initValue: navState) {
function createNavbarStore(key: string, initValue: NavState) {
    // const store = writable<navState>(initValue);
    const store = writable<NavState>(initValue);
    if (!browser) return store;

    const storedValueStr = localStorage.getItem(key);
    if (storedValueStr != null) store.set(JSON.parse(storedValueStr));

    store.subscribe((val) => {
        if ([null, undefined].includes(val)) localStorage.removeItem(key);
        else localStorage.setItem(key, JSON.stringify(val));
    });

    window.addEventListener('storage', () => {
        const storedValueStr = localStorage.getItem(key);
        if (storedValueStr == null) return;

        // const localValue: navState = JSON.parse(storedValueStr);
        const localValue: NavState = JSON.parse(storedValueStr);
        if (localValue !== get(store)) store.set(localValue);
    });
    return store;
}

export const navstate = createNavbarStore('navState', 'accueil');
