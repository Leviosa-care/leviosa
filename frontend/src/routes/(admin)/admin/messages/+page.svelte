<script lang="ts">
	import type { PageProps } from './$types';
	import { MessageCircle, Search, Circle, Check } from '@lucide/svelte';

	let { data }: PageProps = $props();

	let searchQuery = $state('');
	let selectedConversationId = $state(data.conversations[0]?.id ?? '');

	function formatTime(isoString: string): string {
		const date = new Date(isoString);
		const now = new Date();
		const diffMs = now.getTime() - date.getTime();
		const diffHours = diffMs / (1000 * 60 * 60);
		const diffDays = diffMs / (1000 * 60 * 60 * 24);

		if (diffHours < 1) {
			return "À l'instant";
		} else if (diffHours < 24) {
			return `Il y a ${Math.floor(diffHours)}h`;
		} else if (diffDays < 7) {
			return `Il y a ${Math.floor(diffDays)}j`;
		} else {
			return date.toLocaleDateString('fr-FR', { day: '2-digit', month: '2-digit', year: '2-digit' });
		}
	}

	function formatMessageTime(isoString: string): string {
		const date = new Date(isoString);
		return date.toLocaleTimeString('fr-FR', { hour: '2-digit', minute: '2-digit' });
	}

	const filteredConversations = $derived(
		data.conversations.filter(c =>
			c.clientName.toLowerCase().includes(searchQuery.toLowerCase()) ||
			c.lastMessage.toLowerCase().includes(searchQuery.toLowerCase())
		)
	);

	const selectedConversation = $derived(
		data.conversations.find(c => c.id === selectedConversationId)
	);
</script>

<svelte:head>
	<title>Messages | Admin</title>
</svelte:head>

<div class="container mx-auto px-4 py-8 lg:py-12">
	<div class="mb-8">
		<h1 class="text-3xl lg:text-4xl font-bold mb-1 text-foreground">
			Messages
		</h1>
		<p class="text-muted-foreground">
			Gérez vos conversations avec les clients
		</p>
	</div>

	<div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
		<!-- Conversation List -->
		<div class="lg:col-span-1 bg-card rounded-lg border border-border overflow-hidden">
			<!-- Search -->
			<div class="p-4 border-b border-border">
				<div class="relative">
					<Search size={18} class="absolute left-3 top-1/2 -translate-y-1/2 text-muted-foreground" />
					<input
						type="text"
						placeholder="Rechercher..."
						bind:value={searchQuery}
						class="w-full pl-10 pr-4 py-2 bg-background border border-border rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-primary"
					/>
				</div>
			</div>

			<!-- Conversations -->
			<div class="divide-y divide-border max-h-[600px] overflow-y-auto">
				{#each filteredConversations as conversation}
					<button
						class="w-full p-4 text-left hover:bg-muted/50 transition-colors {selectedConversationId === conversation.id
							? 'bg-muted border-l-4 border-l-primary'
							: ''}"
						onclick={() => selectedConversationId = conversation.id}
					>
						<div class="flex items-start justify-between gap-2 mb-1">
							<div class="flex items-center gap-2">
								<div class="relative">
									<div class="w-10 h-10 bg-primary/10 rounded-full flex items-center justify-center">
										<span class="text-sm font-medium text-primary">
											{conversation.clientName.split(' ').map(n => n[0]).join('')}
										</span>
									</div>
									{#if conversation.isOnline}
										<Circle size={10} fill="currentColor" class="absolute -bottom-0.5 -right-0.5 text-green-500" />
									{/if}
								</div>
								<span class="font-medium text-foreground">{conversation.clientName}</span>
							</div>
							<span class="text-xs text-muted-foreground whitespace-nowrap">
								{formatTime(conversation.lastMessageAt)}
							</span>
						</div>
						<p class="text-sm text-muted-foreground truncate mb-2">
							{conversation.lastMessage}
						</p>
						{#if conversation.unreadCount > 0}
							<span class="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium bg-primary text-primary-foreground">
								{conversation.unreadCount} non lu{conversation.unreadCount > 1 ? 's' : ''}
							</span>
						{/if}
					</button>
				{:else}
					<div class="p-8 text-center text-muted-foreground">
						<MessageCircle size={32} class="mx-auto mb-2 opacity-50" />
						Aucune conversation trouvée
					</div>
				{/each}
			</div>
		</div>

		<!-- Message Thread (Desktop) -->
		{#if selectedConversation}
			<div class="hidden lg:block lg:col-span-2 bg-card rounded-lg border border-border overflow-hidden">
				<!-- Header -->
				<div class="p-4 border-b border-border bg-muted/30">
					<div class="flex items-center gap-3">
						<div class="relative">
							<div class="w-10 h-10 bg-primary/10 rounded-full flex items-center justify-center">
								<span class="text-sm font-medium text-primary">
									{selectedConversation.clientName.split(' ').map(n => n[0]).join('')}
								</span>
							</div>
							{#if selectedConversation.isOnline}
								<Circle size={10} fill="currentColor" class="absolute -bottom-0.5 -right-0.5 text-green-500" />
							{/if}
						</div>
						<div>
							<h3 class="font-semibold text-foreground">{selectedConversation.clientName}</h3>
							<span class="text-sm text-muted-foreground">
								{selectedConversation.isOnline ? 'En ligne' : 'Hors ligne'}
							</span>
						</div>
					</div>
				</div>

				<!-- Messages -->
				<div class="p-4 space-y-4 max-h-[500px] overflow-y-auto">
					{#each selectedConversation.messages as message}
						<div class="flex {message.sender === 'admin' ? 'justify-end' : 'justify-start'}">
							<div class="max-w-[70%]">
								<div class="rounded-2xl px-4 py-2 {message.sender === 'admin'
									? 'bg-primary text-primary-foreground rounded-br-sm'
									: 'bg-muted text-foreground rounded-bl-sm'}">
									<p class="text-sm">{message.content}</p>
								</div>
								<div class="flex items-center gap-1 mt-1 px-2">
									<span class="text-xs text-muted-foreground">
										{formatMessageTime(message.sentAt)}
									</span>
									{#if message.sender === 'admin'}
										<Check size={12} class="text-primary" />
									{/if}
								</div>
							</div>
						</div>
					{/each}
				</div>

				<!-- Input (disabled, read-only) -->
				<div class="p-4 border-t border-border bg-muted/20">
					<div class="flex items-center gap-3">
						<input
							type="text"
							placeholder="Répondre à {selectedConversation.clientName}..."
							disabled
							class="flex-1 px-4 py-2 bg-background border border-border rounded-md text-sm opacity-50 cursor-not-allowed"
						/>
						<button
							class="px-4 py-2 bg-primary text-primary-foreground rounded-md opacity-50 cursor-not-allowed"
							disabled
						>
							Envoyer
						</button>
					</div>
				</div>
			</div>
		{/if}
	</div>
</div>
