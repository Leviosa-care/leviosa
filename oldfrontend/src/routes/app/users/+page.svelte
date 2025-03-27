<script lang="ts">
	interface Props {
		data: import('./$types').PageData;
	}
	let { data }: Props = $props();
</script>

<div>
	<p>this page contains all the elemens about users</p>
	<p>reserved for admin users only</p>
</div>
