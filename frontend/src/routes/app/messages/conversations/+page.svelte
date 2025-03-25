<script lang="ts">
	import { onMount } from 'svelte';
	import { redirectTo } from '$lib/scripts/redirect';
	import { NAV_STATES } from '$lib/types';
	import { navigationState } from '$lib/stores/persisted_stores.svelte';

	navigationState.set(NAV_STATES.Conversations); // just to forget the value stored in localstore when reconecting and I had the page to another link.

	onMount(() => {
		if (window.matchMedia('(max-width: 500px)').matches) {
			console.log('need to redirect brother');
			redirectTo('/app/messages');
		}
	});
</script>
