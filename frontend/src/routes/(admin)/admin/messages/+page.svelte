<script lang="ts">
	import type { PageProps } from './$types';
	import type { ThreadSummary, ThreadMessage } from './+page.server';
	import {
		MessageCircle,
		Search,
		Send,
		ArrowLeft,
		Check,
		Plus
	} from '@lucide/svelte';
	import { browser } from '$app/environment';

	let { data }: PageProps = $props();

	let searchQuery = $state('');
	let activeThreadId = $state<string | null>(data.activeThreadId);
	let newMessage = $state('');
	let sending = $state(false);
	let showNewThreadDialog = $state(false);
	let userSearchQuery = $state('');
	let userSearchResults: { id: string; name: string; email: string }[] = $state([]);
	let searchingUsers = $state(false);
	let threads = $state<ThreadSummary[]>([...data.threads]);
	let pollingMessages = $state<ThreadMessage[]>([...data.activeMessages]);

	const totalUnread = $derived(threads.reduce((s, c) => s + c.unread_count, 0));

	const filteredConversations = $derived(
		threads.filter(
			(c) =>
				c.participant_name.toLowerCase().includes(searchQuery.toLowerCase()) ||
				c.last_message.toLowerCase().includes(searchQuery.toLowerCase())
		)
	);

	const activeConversation = $derived(threads.find((c) => c.thread_id === activeThreadId) ?? null);

	function formatTime(isoString: string): string {
		const date = new Date(isoString);
		const now = new Date();
		const diffMs = now.getTime() - date.getTime();
		const diffHours = diffMs / (1000 * 60 * 60);
		const diffDays = diffMs / (1000 * 60 * 60 * 24);

		if (diffHours < 1) return "À l'instant";
		if (diffHours < 24) return `Il y a ${Math.floor(diffHours)}h`;
		if (diffDays < 7) return `Il y a ${Math.floor(diffDays)}j`;
		return date.toLocaleDateString('fr-FR', {
			day: '2-digit',
			month: '2-digit',
			year: '2-digit'
		});
	}

	function formatMessageTime(isoString: string): string {
		return new Date(isoString).toLocaleTimeString('fr-FR', {
			hour: '2-digit',
			minute: '2-digit'
		});
	}

	function getInitials(name: string): string {
		return name
			.split(' ')
			.map((n) => n[0])
			.join('')
			.toUpperCase()
			.slice(0, 2);
	}

	async function refreshData() {
		try {
			const [threadsRes, msgsRes] = await Promise.all([
				fetch('/api/threads'),
				activeThreadId
					? fetch(`/api/threads/${activeThreadId}/messages?limit=100`)
					: Promise.resolve(null)
			]);
			if (threadsRes.ok) threads = await threadsRes.json();
			if (msgsRes?.ok) {
				const d = await msgsRes.json();
				pollingMessages = d.messages ?? [];
			}
		} catch {
			// non-critical polling failure
		}
	}

	$effect(() => {
		if (!browser) return;
		const interval = setInterval(refreshData, 10_000);
		const onVisible = () => { if (!document.hidden) refreshData(); };
		document.addEventListener('visibilitychange', onVisible);
		return () => {
			clearInterval(interval);
			document.removeEventListener('visibilitychange', onVisible);
		};
	});

	async function sendMessage() {
		if (!newMessage.trim() || !activeThreadId || sending) return;
		sending = true;
		const body = newMessage.trim();
		newMessage = '';

		const optimistic: ThreadMessage = {
			id: `temp-${Date.now()}`,
			thread_id: activeThreadId,
			sender_id: data.userId,
			body,
			created_at: new Date().toISOString(),
			read_at: null
		};
		pollingMessages = [...pollingMessages, optimistic];

		try {
			const res = await fetch(`/api/threads/${activeThreadId}/messages`, {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ body })
			});
			if (res.ok) {
				const saved: ThreadMessage = await res.json();
				pollingMessages = pollingMessages.map((m) =>
					m.id === optimistic.id ? saved : m
				);
			} else {
				pollingMessages = pollingMessages.filter((m) => m.id !== optimistic.id);
				newMessage = body;
			}
		} catch {
			pollingMessages = pollingMessages.filter((m) => m.id !== optimistic.id);
			newMessage = body;
		} finally {
			sending = false;
		}
	}

	function handleKeydown(e: KeyboardEvent) {
		if (e.key === 'Enter' && !e.shiftKey) {
			e.preventDefault();
			sendMessage();
		}
	}

	function selectThread(threadId: string) {
		activeThreadId = threadId;
		pollingMessages = [];
		const url = new URL(window.location.href);
		url.searchParams.set('thread', threadId);
		window.location.href = url.toString();
	}

	async function searchUsers() {
		if (!userSearchQuery.trim()) return;
		searchingUsers = true;
		try {
			const res = await fetch(
				`/api/users?search=${encodeURIComponent(userSearchQuery)}`
			);
			if (res.ok) {
				const users = await res.json();
				userSearchResults = users.map((u: { id: string; first_name?: string; last_name?: string; email?: string }) => {
					const fullName = `${u.first_name ?? ''} ${u.last_name ?? ''}`.trim();
					return {
						id: u.id,
						name: fullName || (u.email ?? ''),
						email: u.email ?? ''
					};
				});
			}
		} catch {
			userSearchResults = [];
		} finally {
			searchingUsers = false;
		}
	}

	async function createThread(participantId: string) {
		try {
			const res = await fetch(`/api/threads`, {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ participant_id: participantId })
			});
			if (res.ok) {
				const d = await res.json();
				showNewThreadDialog = false;
				userSearchQuery = '';
				userSearchResults = [];
				selectThread(d.id);
			}
		} catch {
			// Silently fail
		}
	}
</script>

<svelte:head>
	<title>Messages | Admin</title>
</svelte:head>

<div class="h-[calc(100vh-4rem)] lg:h-screen flex flex-col">
	<!-- Header -->
	<div class="border-b border-border px-6 py-4 flex items-center justify-between">
		<div>
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
			<p class="text-sm text-muted-foreground">Gérez vos conversations</p>
		</div>
		<button
			onclick={() => (showNewThreadDialog = true)}
			class="flex items-center gap-2 px-4 py-2 rounded-lg bg-foreground text-background text-sm font-medium hover:opacity-90 transition-opacity"
		>
			<Plus size={16} />
			Nouveau message
		</button>
	</div>

	<div class="flex flex-1 overflow-hidden">
		<!-- Conversation List -->
		<div
			class="w-full lg:w-80 xl:w-96 border-r border-border flex-shrink-0 flex flex-col overflow-hidden {activeThreadId
				? 'hidden lg:flex'
				: 'flex'}"
		>
			<!-- Search -->
			<div class="p-4 border-b border-border">
				<div class="relative">
					<Search
						size={18}
						class="absolute left-3 top-1/2 -translate-y-1/2 text-muted-foreground"
					/>
					<input
						type="text"
						placeholder="Rechercher..."
						bind:value={searchQuery}
						class="w-full pl-10 pr-4 py-2 bg-background border border-border rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-foreground/20"
					/>
				</div>
			</div>

			<!-- Conversations -->
			<div class="flex-1 overflow-y-auto">
				{#each filteredConversations as conversation (conversation.thread_id)}
					<button
						class="w-full text-left px-4 py-4 border-b border-border hover:bg-muted/40 transition-colors {activeThreadId ===
						conversation.thread_id
							? 'bg-muted/60'
							: ''}"
						onclick={() => selectThread(conversation.thread_id)}
					>
						<div class="flex items-start gap-3">
							<div
								class="w-10 h-10 rounded-full bg-muted flex items-center justify-center flex-shrink-0"
							>
								<span class="text-sm font-semibold text-foreground">
									{getInitials(conversation.participant_name)}
								</span>
							</div>
							<div class="flex-1 min-w-0">
								<div class="flex items-center justify-between mb-0.5">
									<p class="font-semibold text-sm text-foreground truncate">
										{conversation.participant_name}
									</p>
									<span class="text-xs text-muted-foreground flex-shrink-0 ml-2">
										{formatTime(conversation.last_message_at)}
									</span>
								</div>
								<div class="flex items-center justify-between">
									<p class="text-xs text-muted-foreground truncate">
										{conversation.last_message}
									</p>
									{#if conversation.unread_count > 0}
										<span
											class="ml-2 flex-shrink-0 w-5 h-5 rounded-full bg-foreground text-background text-xs flex items-center justify-center font-medium"
										>
											{conversation.unread_count}
										</span>
									{/if}
								</div>
							</div>
						</div>
					</button>
				{:else}
					<div class="p-8 text-center text-muted-foreground">
						<MessageCircle size={32} class="mx-auto mb-2 opacity-50" />
						Aucune conversation trouvée
					</div>
				{/each}
			</div>
		</div>

		<!-- Message Thread -->
		{#if activeConversation}
			<div
				class="flex-1 flex flex-col overflow-hidden {activeThreadId
					? 'block'
					: 'hidden lg:block'}"
			>
				<!-- Thread Header -->
				<div class="px-6 py-4 border-b border-border flex items-center gap-3">
					<button
						class="lg:hidden p-1 rounded-md hover:bg-muted transition-colors text-muted-foreground"
						onclick={() => (activeThreadId = null)}
					>
						<ArrowLeft size={20} />
					</button>
					<div class="w-9 h-9 rounded-full bg-muted flex items-center justify-center">
						<span class="text-sm font-semibold text-foreground">
							{getInitials(activeConversation.participant_name)}
						</span>
					</div>
					<div>
						<p class="font-semibold text-foreground">
							{activeConversation.participant_name}
						</p>
						<p class="text-xs text-muted-foreground">Utilisateur</p>
					</div>
				</div>

				<!-- Messages -->
				<div class="flex-1 overflow-y-auto px-6 py-4 space-y-4">
					{#each pollingMessages as message (message.id)}
						<div
							class="flex {message.sender_id === data.userId
								? 'justify-end'
								: 'justify-start'}"
						>
							<div class="max-w-[70%]">
								<div
									class="rounded-2xl px-4 py-2 text-sm {message.sender_id ===
									data.userId
										? 'bg-foreground text-background rounded-br-sm'
										: 'bg-muted text-foreground rounded-bl-sm'}"
								>
									<p>{message.body}</p>
								</div>
								<div class="flex items-center gap-1 mt-1 px-2">
									<span class="text-xs text-muted-foreground">
										{formatMessageTime(message.created_at)}
									</span>
									{#if message.sender_id === data.userId}
										<Check size={12} class="text-muted-foreground" />
									{/if}
								</div>
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
							onkeydown={handleKeydown}
							placeholder="Répondre à {activeConversation.participant_name}..."
							class="flex-1 px-4 py-2.5 rounded-lg border border-border bg-background text-sm focus:outline-none focus:ring-2 focus:ring-foreground/20"
						/>
						<button
							onclick={sendMessage}
							class="p-2.5 rounded-lg bg-foreground text-background hover:opacity-90 transition-opacity disabled:opacity-40"
							disabled={!newMessage.trim() || sending}
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

<!-- New Thread Dialog -->
{#if showNewThreadDialog}
	<div
		class="fixed inset-0 z-50 flex items-center justify-center bg-black/50"
		onclick={() => (showNewThreadDialog = false)}
	>
		<div
			class="bg-background rounded-xl border border-border shadow-lg w-full max-w-md mx-4"
			onclick={(e) => e.stopPropagation()}
		>
			<div class="px-6 py-4 border-b border-border">
				<h2 class="text-lg font-semibold text-foreground">Nouveau message</h2>
				<p class="text-sm text-muted-foreground">Rechercher un utilisateur</p>
			</div>
			<div class="p-6">
				<div class="flex gap-2 mb-4">
					<input
						type="text"
						placeholder="Nom ou email..."
						bind:value={userSearchQuery}
						onkeydown={(e) => e.key === 'Enter' && searchUsers()}
						class="flex-1 px-4 py-2 rounded-lg border border-border bg-background text-sm focus:outline-none focus:ring-2 focus:ring-foreground/20"
					/>
					<button
						onclick={searchUsers}
						disabled={searchingUsers}
						class="px-4 py-2 rounded-lg bg-foreground text-background text-sm font-medium hover:opacity-90 disabled:opacity-40"
					>
						{searchingUsers ? '...' : 'Chercher'}
					</button>
				</div>
				{#if userSearchResults.length > 0}
					<div class="space-y-1 max-h-64 overflow-y-auto">
						{#each userSearchResults as user (user.id)}
							<button
								class="w-full text-left px-4 py-3 rounded-lg hover:bg-muted/40 transition-colors border border-border"
								onclick={() => createThread(user.id)}
							>
								<p class="text-sm font-medium text-foreground">{user.name}</p>
								<p class="text-xs text-muted-foreground">{user.email}</p>
							</button>
						{/each}
					</div>
				{/if}
			</div>
		</div>
	</div>
{/if}
