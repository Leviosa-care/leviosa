<script lang="ts">
	import type { LayoutProps } from "./$types";
	import { setUserContext } from "$lib/context/user.svelte";
	import AdminSidebar from "$lib/navigation/AdminSidebar.svelte";

	let { data, children }: LayoutProps = $props();

	setUserContext(data.user);
</script>

<div class="flex flex-col lg:flex-row">
	<AdminSidebar user={data.user} permissions={data.permissions} />
	{@render children()}
</div>
