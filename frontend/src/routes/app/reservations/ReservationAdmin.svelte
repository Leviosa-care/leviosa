<script lang="ts">
	import Tabs from '$lib/components/Tabs.svelte';
	// import { Tabs } from 'bits-ui';
	// import { Plane } from 'lucide-svelte';
	let { cards } = $props();
	console.log('the card information is:', cards);
	// NOTE: il me faut le state de reservations
	// I need a state in here
	type EventState = 'create' | 'all' | 'past';
	let eventState: EventState = 'all';
</script>

<div>
	<div class="bg-red-400">Espace reservations</div>
	<div class="flex bg-red-500">
		<button class="bg-red-500 text-xl">Creer</button>
		<button>Voir les evenements</button>
		<button>Past events</button>
	</div>
	{#if eventState === 'all'}
		<p>display all events in here</p>
	{:else if eventState === 'past'}
		<p>display all past events in herer</p>
	{/if}
	<div style="max-width: 400px;">
		<Tabs />
	</div>
</div>

<!-- <div class="pt-6" style="margin-top: 40px;"> -->
<!-- 	<Tabs.Root -->
<!-- 		value="outbound" -->
<!-- 		class="rounded-card bg-background-alt w-[390px] border border-muted p-3 shadow-card" -->
<!-- 	> -->
<!-- 		<Tabs.List -->
<!-- 			class="rounded-9px bg-dark-10 shadow-mini-inset grid w-full grid-cols-2 gap-1 p-1 text-sm font-semibold leading-[0.01em] dark:border dark:border-neutral-600/30 dark:bg-background" -->
<!-- 		> -->
<!-- 			<Tabs.Trigger -->
<!-- 				value="outbound" -->
<!-- 				class="data-[state=active]:shadow-mini h-8 rounded-[7px] bg-transparent py-2 data-[state=active]:bg-white dark:data-[state=active]:bg-muted" -->
<!-- 				>Outbound</Tabs.Trigger -->
<!-- 			> -->
<!-- 			<Tabs.Trigger -->
<!-- 				value="inbound" -->
<!-- 				class="data-[state=active]:shadow-mini h-8 rounded-[7px] bg-transparent py-2 data-[state=active]:bg-white dark:data-[state=active]:bg-muted" -->
<!-- 				>Inbound</Tabs.Trigger -->
<!-- 			> -->
<!-- 		</Tabs.List> -->
<!-- 		<Tabs.Content value="outbound" class="select-none pt-3"> -->
<!-- 			<div class="grid grid-cols-3 grid-rows-2 gap-0 p-4 pb-1"> -->
<!-- 				<div class="text-left"> -->
<!-- 					<h4 class="mb-2 text-[20px] font-semibold leading-none tracking-[-0.01em]">Prague</h4> -->
<!-- 					<p class="text-sm font-medium text-muted-foreground">06:05</p> -->
<!-- 				</div> -->
<!-- 				<div class="self-end text-center"> -->
<!-- 					<p class="text-sm font-medium text-muted-foreground">3h 30m</p> -->
<!-- 				</div> -->
<!-- 				<div class="text-right"> -->
<!-- 					<h4 class="mb-2 text-[20px] font-semibold leading-none tracking-[-0.01em]">Malaga</h4> -->
<!-- 					<p class="text-sm font-medium text-muted-foreground">06:05</p> -->
<!-- 				</div> -->
<!-- 				<div class="relative col-span-3"> -->
<!-- 					<hr class="border-1 border-border-input relative top-4 h-px border-dashed" /> -->
<!---->
<!-- 					<div class="bg-background-alt absolute left-1/2 -translate-x-1/2 p-1"> -->
<!-- 						<Plane class="size-6 rotate-90 text-muted-foreground" /> -->
<!-- 					</div> -->
<!-- 				</div> -->
<!-- 			</div> -->
<!-- 		</Tabs.Content> -->
<!-- 		<Tabs.Content value="inbound" class="select-none pt-3"> -->
<!-- 			<div class="grid grid-cols-3 grid-rows-2 gap-0 p-4 pb-1"> -->
<!-- 				<div class="text-left"> -->
<!-- 					<h4 class="mb-2 text-[20px] font-semibold leading-none tracking-[-0.01em]">Malaga</h4> -->
<!-- 					<p class="text-sm font-medium text-muted-foreground">07:25</p> -->
<!-- 				</div> -->
<!-- 				<div class="self-end text-center"> -->
<!-- 					<p class="text-sm font-medium text-muted-foreground">3h 20m</p> -->
<!-- 				</div> -->
<!-- 				<div class="text-right"> -->
<!-- 					<h4 class="mb-2 text-[20px] font-semibold leading-none tracking-[-0.01em]">Prague</h4> -->
<!-- 					<p class="text-sm font-medium text-muted-foreground">10:45</p> -->
<!-- 				</div> -->
<!-- 				<div class="relative col-span-3"> -->
<!-- 					<hr class="border-1 border-border-input relative top-4 h-px border-dashed" /> -->
<!---->
<!-- 					<div class="bg-background-alt absolute left-1/2 -translate-x-1/2 p-1"> -->
<!-- 						<Plane class="size-6 rotate-90 text-muted-foreground" /> -->
<!-- 					</div> -->
<!-- 				</div> -->
<!-- 			</div> -->
<!-- 		</Tabs.Content> -->
<!-- 	</Tabs.Root> -->
<!-- </div> -->

<style>
</style>
