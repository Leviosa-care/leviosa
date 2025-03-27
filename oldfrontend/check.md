> frontend@0.0.1 check /home/henga/Documents/projects/livio/leviosa/.worktrees/frontend/frontend
> svelte-kit sync && svelte-check --tsconfig ./tsconfig.json

====================================
Loading svelte-check in workspace: /home/henga/Documents/projects/livio/leviosa/.worktrees/frontend/frontend
Getting Svelte diagnostics...

/home/henga/Documents/projects/livio/leviosa/.worktrees/frontend/frontend/src/lib/components/ui/form/form-description.svelte:2:33
Error: Cannot find module 'formsnap' or its corresponding type declarations. (ts)

<script lang="ts">
	import * as FormPrimitive from 'formsnap';
	import type { HTMLAttributes } from 'svelte/elements';


/home/henga/Documents/projects/livio/leviosa/.worktrees/frontend/frontend/src/lib/components/ui/form/form-label.svelte:3:33
Error: Cannot find module 'formsnap' or its corresponding type declarations. (ts)
	import type { Label as LabelPrimitive } from 'bits-ui';
	import { getFormControl } from 'formsnap';
	import { cn } from '$lib/utils.js';


/home/henga/Documents/projects/livio/leviosa/.worktrees/frontend/frontend/src/lib/components/ui/form/form-field-errors.svelte:2:33
Error: Cannot find module 'formsnap' or its corresponding type declarations. (ts)
<script lang="ts">
	import * as FormPrimitive from 'formsnap';
	import { cn } from '$lib/utils.js';


/home/henga/Documents/projects/livio/leviosa/.worktrees/frontend/frontend/src/lib/components/ui/form/form-field.svelte:2:43
Error: Cannot find module 'sveltekit-superforms' or its corresponding type declarations. (ts)
<script lang="ts" context="module">
	import type { FormPath, SuperForm } from 'sveltekit-superforms';
	type T = Record<string, unknown>;


/home/henga/Documents/projects/livio/leviosa/.worktrees/frontend/frontend/src/lib/components/ui/form/form-field.svelte:9:33
Error: Cannot find module 'formsnap' or its corresponding type declarations. (ts)
	import type { HTMLAttributes } from 'svelte/elements';
	import * as FormPrimitive from 'formsnap';
	import { cn } from '$lib/utils.js';


/home/henga/Documents/projects/livio/leviosa/.worktrees/frontend/frontend/src/lib/components/ui/form/form-fieldset.svelte:2:43
Error: Cannot find module 'sveltekit-superforms' or its corresponding type declarations. (ts)
<script lang="ts" context="module">
	import type { FormPath, SuperForm } from 'sveltekit-superforms';
	type T = Record<string, unknown>;


/home/henga/Documents/projects/livio/leviosa/.worktrees/frontend/frontend/src/lib/components/ui/form/form-fieldset.svelte:8:33
Error: Cannot find module 'formsnap' or its corresponding type declarations. (ts)
<script lang="ts" generics="T extends Record<string, unknown>, U extends FormPath<T>">
	import * as FormPrimitive from 'formsnap';
	import { cn } from '$lib/utils.js';


/home/henga/Documents/projects/livio/leviosa/.worktrees/frontend/frontend/src/lib/components/ui/form/form-legend.svelte:2:33
Error: Cannot find module 'formsnap' or its corresponding type declarations. (ts)
<script lang="ts">
	import * as FormPrimitive from 'formsnap';
	import { cn } from '$lib/utils.js';


/home/henga/Documents/projects/livio/leviosa/.worktrees/frontend/frontend/src/lib/components/ui/form/form-element-field.svelte:2:49
Error: Cannot find module 'sveltekit-superforms' or its corresponding type declarations. (ts)
<script lang="ts" context="module">
	import type { FormPathLeaves, SuperForm } from 'sveltekit-superforms';
	type T = Record<string, unknown>;


/home/henga/Documents/projects/livio/leviosa/.worktrees/frontend/frontend/src/lib/components/ui/form/form-element-field.svelte:9:33
Error: Cannot find module 'formsnap' or its corresponding type declarations. (ts)
	import type { HTMLAttributes } from 'svelte/elements';
	import * as FormPrimitive from 'formsnap';
	import { cn } from '$lib/utils.js';


/home/henga/Documents/projects/livio/leviosa/.worktrees/frontend/frontend/src/lib/components/ui/form/index.ts:1:32
Error: Cannot find module 'formsnap' or its corresponding type declarations. 
import * as FormPrimitive from 'formsnap';
import Description from './form-description.svelte';


/home/henga/Documents/projects/livio/leviosa/.worktrees/frontend/frontend/src/lib/scripts/credentials.test.ts:6:23
Error: Cannot find name 'validate'. 
		(email: string, password: string) => {
			expect(async () => validate(email, password)).rejects.toThrowError();
		}


/home/henga/Documents/projects/livio/leviosa/.worktrees/frontend/frontend/src/lib/scripts/parseCookie.test.ts:5:24
Error: Cannot find name 'parseCookie'. 
	it.todo('should do nothing for now', () => {
		const cookieParsed = parseCookie(cookie);
		// TODO: find the right format for the cookoie that I am going to exploit


/home/henga/Documents/projects/livio/leviosa/.worktrees/frontend/frontend/src/lib/stores/navbar.ts:16:34
Error: Argument of type 'navState' is not assignable to parameter of type 'null | undefined'.
  Type '"events"' is not assignable to type 'null | undefined'. 
	store.subscribe((val) => {
		if ([null, undefined].includes(val)) {
			localStorage.removeItem(key);


/home/henga/Documents/projects/livio/leviosa/.worktrees/frontend/frontend/src/routes/app/admin/+page.server.ts:8:18
Error: Property 'role' does not exist on type 'never'. 
	}
	if (locals.user.role !== 'admin') {
		throw redirect(301, '/');


/home/henga/Documents/projects/livio/leviosa/.worktrees/frontend/frontend/src/routes/app/admin/users/+page.server.ts:42:26
Error: Property 'get' does not exist on type 'Promise<FormData>'. 
	const formData = request.formData();
	const userId = formData.get('userid');
	const res = await fetch(`${API_URL}/api/v1/admin/users`, {


/home/henga/Documents/projects/livio/leviosa/.worktrees/frontend/frontend/src/routes/app/settings/profile/+page.ts:72:17
Error: Element implicitly has an 'any' type because expression of type 'string' can't be used to index type '{ email: string; role: string; lastname: string; firstname: string; gender: string; birthdate: string; telephone: string; address: string; city: string; postalcard: string; }'.
  No index signature with a parameter of type 'string' was found on type '{ email: string; role: string; lastname: string; firstname: string; gender: string; birthdate: string; telephone: string; address: string; city: string; postalcard: string; }'. 
	for (const field of userfields) {
		const value = user[field] !== '' ? user[field] : missingValues[field];
		fields[field] = { ...fieldsStatic[field], value };


/home/henga/Documents/projects/livio/leviosa/.worktrees/frontend/frontend/src/routes/app/settings/profile/+page.ts:72:38
Error: Element implicitly has an 'any' type because expression of type 'string' can't be used to index type '{ email: string; role: string; lastname: string; firstname: string; gender: string; birthdate: string; telephone: string; address: string; city: string; postalcard: string; }'.
  No index signature with a parameter of type 'string' was found on type '{ email: string; role: string; lastname: string; firstname: string; gender: string; birthdate: string; telephone: string; address: string; city: string; postalcard: string; }'. 
	for (const field of userfields) {
		const value = user[field] !== '' ? user[field] : missingValues[field];
		fields[field] = { ...fieldsStatic[field], value };


/home/henga/Documents/projects/livio/leviosa/.worktrees/frontend/frontend/src/routes/app/settings/profile/+page.ts:72:52
Error: Element implicitly has an 'any' type because expression of type 'string' can't be used to index type '{ readonly email: "Aucune adresse email précisé"; readonly telephone: "Aucun numero de telephone précisé"; readonly lastname: "Aucun nom précisé"; readonly firstname: "Aucun prenom précisé"; readonly birthdate: "Aucune date de naissance précisé"; readonly address: "Aucune adresse précisé"; readonly city: "Aucune vil...'.
  No index signature with a parameter of type 'string' was found on type '{ readonly email: "Aucune adresse email précisé"; readonly telephone: "Aucun numero de telephone précisé"; readonly lastname: "Aucun nom précisé"; readonly firstname: "Aucun prenom précisé"; readonly birthdate: "Aucune date de naissance précisé"; readonly address: "Aucune adresse précisé"; readonly city: "Aucune vil...'. 
	for (const field of userfields) {
		const value = user[field] !== '' ? user[field] : missingValues[field];
		fields[field] = { ...fieldsStatic[field], value };


/home/henga/Documents/projects/livio/leviosa/.worktrees/frontend/frontend/src/routes/app/votes/[year]/[month]/+page.server.ts:46:2
Error: Object literal may only specify known properties, and 'default' does not exist in type 'Action'. 
export const actions: Action = {
	default: async ({ request, cookies, params }) => {
		console.log('Sending the data to the backend !');


/home/henga/Documents/projects/livio/leviosa/.worktrees/frontend/frontend/src/routes/app/votes/[year]/[month]/+page.server.ts:46:20
Error: Binding element 'request' implicitly has an 'any' type. 
export const actions: Action = {
	default: async ({ request, cookies, params }) => {
		console.log('Sending the data to the backend !');


/home/henga/Documents/projects/livio/leviosa/.worktrees/frontend/frontend/src/routes/app/votes/[year]/[month]/+page.server.ts:46:29
Error: Binding element 'cookies' implicitly has an 'any' type. 
export const actions: Action = {
	default: async ({ request, cookies, params }) => {
		console.log('Sending the data to the backend !');


/home/henga/Documents/projects/livio/leviosa/.worktrees/frontend/frontend/src/routes/app/votes/[year]/[month]/+page.server.ts:46:38
Error: Binding element 'params' implicitly has an 'any' type. 
export const actions: Action = {
	default: async ({ request, cookies, params }) => {
		console.log('Sending the data to the backend !');


/home/henga/Documents/projects/livio/leviosa/.worktrees/frontend/frontend/src/routes/signup/+page.server.ts:9:2
Error: Object literal may only specify known properties, and 'default' does not exist in type 'Action'. 
export const actions: Action = {
	default: async ({ request, cookies }) => {
		const formData = await request.formData();


/home/henga/Documents/projects/livio/leviosa/.worktrees/frontend/frontend/src/routes/signup/+page.server.ts:9:20
Error: Binding element 'request' implicitly has an 'any' type. 
export const actions: Action = {
	default: async ({ request, cookies }) => {
		const formData = await request.formData();


/home/henga/Documents/projects/livio/leviosa/.worktrees/frontend/frontend/src/routes/signup/+page.server.ts:9:29
Error: Binding element 'cookies' implicitly has an 'any' type. 
export const actions: Action = {
	default: async ({ request, cookies }) => {
		const formData = await request.formData();


/home/henga/Documents/projects/livio/leviosa/.worktrees/frontend/frontend/src/routes/signup/+page.server.ts:29:42
Error: Property 'sessionId' does not exist on type 'CookieParsed'. 
			const cookieParsed = parseCookie(res.headers.getSetCookie()[0]);
			cookies.set('sessionId', cookieParsed.sessionId, {
				path: '/'


/home/henga/Documents/projects/livio/leviosa/.worktrees/frontend/frontend/src/lib/components/events/EventComponent.svelte:30:19
Error: Argument of type 'number' is not assignable to parameter of type 'Date'. (ts)
					<p class="text-sm">
						{formatDate(Date.parse(beginat))} -
						<span class="placecount">


/home/henga/Documents/projects/livio/leviosa/.worktrees/frontend/frontend/src/lib/components/events/EventComponent.svelte:30:30
Error: Argument of type 'Date' is not assignable to parameter of type 'string'. (ts)
					<p class="text-sm">
						{formatDate(Date.parse(beginat))} -
						<span class="placecount">


/home/henga/Documents/projects/livio/leviosa/.worktrees/frontend/frontend/src/lib/components/events/OldTableEvent.svelte:3:32
Error: Cannot find module './$types' or its corresponding type declarations. (ts)
	import * as Table from '$lib/components/ui/table';
	import type { PageData } from './$types';
	export let data: PageData;


/home/henga/Documents/projects/livio/leviosa/.worktrees/frontend/frontend/src/lib/components/home/NextVote.svelte:33:2
Warn: Do not use empty rulesets (css)
	}
	.places {
	}


====================================
svelte-check found 30 errors and 1 warning in 19 files
 ELIFECYCLE  Command failed with exit code 1.
