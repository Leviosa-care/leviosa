import type { Component } from "svelte";

import Google from "@assets/google.svelte";
import Apple from "@assets/apple.svelte";

export type OAuthButton = { icon: Component; name: string; };

export const buttons: OAuthButton[] = [
    { icon: Google, name: "google" },
    { icon: Apple, name: "apple" },
];
