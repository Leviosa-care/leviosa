<script lang="ts">
	import type { PageProps } from './$types';
	import type { Conversation } from './+page.server';
	import { Send } from '@lucide/svelte';

	let { data }: PageProps = $props();

	let selectedId = $state<string | null>(data.conversations[0]?.id ?? null);
	let newMessage = $state('');

	const selected = $derived(data.conversations.find((c) => c.id === selectedId) ?? null);

	function formatRelative(iso: string): string {
		const diff = (Date.now() - new Date(iso).getTime()) / 1000;
		if (diff < 60) return 'À l\'instant';
		if (diff < 3600) return `Il y a ${Math.floor(diff / 60)} min`;
		if (diff < 86400) return `Il y a ${Math.floor(diff / 3600)} h`;
		return new Date(iso).toLocaleDateString('fr-FR', { day: 'numeric', month: 'short' });
	}

	function formatTime(iso: string): string {
		return new Date(iso).toLocaleTimeString('fr-FR', { hour: '2-digit', minute: '2-digit' });
	}

	const totalUnread = $derived(data.conversations.reduce((s, c) => s + c.unreadCount, 0));
</script>

<svelte:head>
	<title>Messages | Staff</title>
</svelte:head>

<div class="h-[calc(100vh-4rem)] lg:h-screen flex flex-col">
	<div class="border-b border-border px-6 py-4">
		<h1 class="text-2xl font-bold text-foreground">
			Messages
			{#if totalUnread > 0}
				<span
					class="ml-2 px-2 py-0.5 text-xs font-medium bg-foreground text-background rounded-full"
				>
					{totalUnread}
				</span>
			{/if}
		</h1>
	</div>

	<div class="flex flex-1 overflow-hidden">
		<!-- Conversation List -->
		<div
			class="w-full lg:w-80 xl:w-96 border-r border-border flex-shrink-0 overflow-y-auto {selectedId
				? 'hidden lg:block'
				: 'block'}"
		>
			{#each data.conversations as convo (convo.id)}
				<button
					class="w-full text-left px-4 py-4 border-b border-border hover:bg-muted/40 transition-colors {selectedId ===
					convo.id
						? 'bg-muted/60'
						: ''}"
					onclick={() => (selectedId = convo.id)}
				>
					<div class="flex items-start gap-3">
						<div
							class="w-10 h-10 rounded-full bg-muted flex items-center justify-center flex-shrink-0"
						>
							<span class="text-sm font-semibold text-foreground">{convo.clientInitials}</span>
						</div>
						<div class="flex-1 min-w-0">
							<div class="flex items-center justify-between mb-0.5">
								<p class="font-semibold text-sm text-foreground truncate">{convo.clientName}</p>
								<span class="text-xs text-muted-foreground flex-shrink-0 ml-2">
									{formatRelative(convo.lastMessageAt)}
								</span>
							</div>
							<div class="flex items-center justify-between">
								<p class="text-xs text-muted-foreground truncate">{convo.lastMessage}</p>
								{#if convo.unreadCount > 0}
									<span
										class="ml-2 flex-shrink-0 w-5 h-5 rounded-full bg-foreground text-background text-xs flex items-center justify-center font-medium"
									>
										{convo.unreadCount}
									</span>
								{/if}
							</div>
						</div>
					</div>
				</button>
			{/each}
		</div>

		<!-- Message Thread -->
		{#if selected}
			<div class="flex-1 flex flex-col overflow-hidden {selectedId ? 'block' : 'hidden lg:block'}">
				<!-- Thread Header -->
				<div class="px-6 py-4 border-b border-border flex items-center gap-3">
					<button
						class="lg:hidden p-1 rounded-md hover:bg-muted transition-colors text-muted-foreground"
						onclick={() => (selectedId = null)}
					>
						←
					</button>
					<div class="w-9 h-9 rounded-full bg-muted flex items-center justify-center">
						<span class="text-sm font-semibold text-foreground">{selected.clientInitials}</span>
					</div>
					<div>
						<p class="font-semibold text-foreground">{selected.clientName}</p>
						<p class="text-xs text-muted-foreground">Client</p>
					</div>
				</div>

				<!-- Messages -->
				<div class="flex-1 overflow-y-auto px-6 py-4 space-y-4">
					{#each selected.messages as msg (msg.id)}
						<div class="flex {msg.fromPartner ? 'justify-end' : 'justify-start'}">
							<div
								class="max-w-xs lg:max-w-md px-4 py-2.5 rounded-2xl text-sm {msg.fromPartner
									? 'bg-foreground text-background rounded-br-sm'
									: 'bg-muted text-foreground rounded-bl-sm'}"
							>
								<p>{msg.content}</p>
								<p
									class="text-xs mt-1 {msg.fromPartner
										? 'text-background/60'
										: 'text-muted-foreground'}"
								>
									{formatTime(msg.sentAt)}
								</p>
							</div>
						</div>
					{/each}
				</div>

				<!-- Compose -->
				<div class="px-6 py-4 border-t border-border">
					<div class="flex items-center gap-3">
						<input
							type="text"
							bind:value={newMessage}
							placeholder="Écrire un message…"
							class="flex-1 px-4 py-2.5 rounded-lg border border-border bg-background text-sm focus:outline-none focus:ring-2 focus:ring-foreground/20"
						/>
						<button
							class="p-2.5 rounded-lg bg-foreground text-background hover:opacity-90 transition-opacity disabled:opacity-40"
							disabled={!newMessage.trim()}
						>
							<Send size={16} />
						</button>
					</div>
				</div>
			</div>
		{:else}
			<div class="flex-1 hidden lg:flex items-center justify-center text-muted-foreground">
				Sélectionnez une conversation
			</div>
		{/if}
	</div>
</div>
