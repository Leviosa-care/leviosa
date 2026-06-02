<script lang="ts">
	import type { PageProps } from './$types';
	import type { ThreadSummary, ThreadMessage } from './+page.server';
	import { Send, ArrowLeft, MessageCircle } from '@lucide/svelte';
	import { browser } from '$app/environment';
	import { tick } from 'svelte';

	let { data }: PageProps = $props();

	let activeThreadId = $state<string | null>(data.activeThreadId);
	let newMessage = $state('');
	let sending = $state(false);
	let threads = $state<ThreadSummary[]>([...data.threads]);
	let pollingMessages = $state<ThreadMessage[]>([...data.activeMessages]);
	let messagesEl = $state<HTMLDivElement | null>(null);

	async function refreshThreads() {
		try {
			const threadsRes = await fetch('/api/threads');
			if (threadsRes.ok) threads = await threadsRes.json();
		} catch {
			// non-critical polling failure
		}
	}

	// SSE: open EventSource when a thread is selected, close on deselect or change.
	$effect(() => {
		if (!browser || !activeThreadId) return;

		const es = new EventSource(`/api/threads/${activeThreadId}/events`);
		es.onmessage = (e) => {
			try {
				const msg: ThreadMessage = JSON.parse(e.data);
				if (!pollingMessages.some((m) => m.id === msg.id)) {
					pollingMessages = [...pollingMessages, msg];
				}
			} catch {
				// Ignore malformed events.
			}
		};
		es.onerror = () => {
			// Browser's built-in retry handles reconnection.
		};

		return () => es.close();
	});

	// Periodic refresh for thread list (unread counts, last messages).
	$effect(() => {
		if (!browser) return;
		const interval = setInterval(refreshThreads, 10_000);
		const onVisible = () => {
			if (!document.hidden) refreshThreads();
		};
		document.addEventListener('visibilitychange', onVisible);
		return () => {
			clearInterval(interval);
			document.removeEventListener('visibilitychange', onVisible);
		};
	});

	// Scroll to the latest message whenever the list grows or the thread changes.
	$effect(() => {
		void pollingMessages.length;
		void activeThreadId;
		if (!browser || !messagesEl) return;
		tick().then(() => {
			if (messagesEl) messagesEl.scrollTop = messagesEl.scrollHeight;
		});
	});

	const activeThread = $derived(
		threads.find((t) => t.thread_id === activeThreadId) ?? null
	);

	function formatRelative(iso: string): string {
		const diff = (Date.now() - new Date(iso).getTime()) / 1000;
		if (diff < 60) return "À l'instant";
		if (diff < 3600) return `Il y a ${Math.floor(diff / 60)} min`;
		if (diff < 86400) return `Il y a ${Math.floor(diff / 3600)} h`;
		return new Date(iso).toLocaleDateString('fr-FR', {
			day: 'numeric',
			month: 'short'
		});
	}

	function formatTime(iso: string): string {
		return new Date(iso).toLocaleTimeString('fr-FR', {
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
				pollingMessages = pollingMessages.filter(
					(m) => m.id !== optimistic.id
				);
				newMessage = body;
			}
		} catch {
			pollingMessages = pollingMessages.filter(
				(m) => m.id !== optimistic.id
			);
			newMessage = body;
		} finally {
			sending = false;
		}
	}

	function selectThread(threadId: string) {
		activeThreadId = threadId;
		pollingMessages = [];

		// Clear the unread badge immediately so the UI responds without waiting
		// for the next poll cycle.
		threads = threads.map((t) =>
			t.thread_id === threadId ? { ...t, unread_count: 0 } : t
		);

		fetch(`/api/threads/${threadId}/messages?limit=100`)
			.then((res) => (res.ok ? res.json() : null))
			.then((d) => {
				if (d) {
					const fetched: ThreadMessage[] = d.messages ?? [];
					const sseOnly = pollingMessages.filter(
						(m) => !fetched.some((f) => f.id === m.id)
					);
					pollingMessages = [...fetched, ...sseOnly];
				}
			})
			.catch(() => {});

		// Update URL without full page navigation
		const url = new URL(window.location.href);
		url.searchParams.set('thread', threadId);
		window.history.replaceState({}, '', url.toString());

		// Mark thread as read
		fetch(`/api/threads/${threadId}/read`, { method: 'POST' }).catch(() => {});
	}

	function goBack() {
		activeThreadId = null;
		const url = new URL(window.location.href);
		url.searchParams.delete('thread');
		window.history.replaceState({}, '', url.toString());
	}

	function handleKeydown(e: KeyboardEvent) {
		if (e.key === 'Enter' && !e.shiftKey) {
			e.preventDefault();
			sendMessage();
		}
	}
</script>

<svelte:head>
	<title>Messages | Espace client</title>
</svelte:head>

{#if threads.length === 0 && !activeThreadId}
	<!-- Empty state: no threads yet -->
	<div class="flex flex-col items-center justify-center py-24 text-center">
		<div
			class="w-16 h-16 rounded-full bg-muted flex items-center justify-center mb-4"
		>
			<MessageCircle size={28} class="text-muted-foreground" />
		</div>
		<h2 class="text-xl font-semibold text-foreground mb-2">Aucun message</h2>
		<p class="text-sm text-muted-foreground max-w-sm">
			La messagerie sera disponible après votre première réservation. Votre
			partenaire pourra vous contacter ici.
		</p>
	</div>
{:else}
	<div class="h-[calc(100vh-4rem)] flex flex-col -mx-4 -my-6 lg:-my-10">
		<!-- Header -->
		<div class="border-b border-border px-6 py-4">
			<h1 class="text-2xl font-bold text-foreground">Messages</h1>
		</div>

		<div class="flex flex-1 overflow-hidden">
			<!-- Conversation List -->
			<div
				class="w-full lg:w-80 xl:w-96 border-r border-border flex-shrink-0 overflow-y-auto {activeThreadId
					? 'hidden lg:block'
					: 'block'}"
			>
				{#each threads as convo (convo.thread_id)}
					<button
						class="w-full text-left px-4 py-4 border-b border-border hover:bg-muted/40 transition-colors {activeThreadId ===
						convo.thread_id
							? 'bg-muted/60'
							: ''}"
						onclick={() => selectThread(convo.thread_id)}
					>
						<div class="flex items-start gap-3">
							<div
								class="w-10 h-10 rounded-full bg-muted flex items-center justify-center flex-shrink-0"
							>
								<span class="text-sm font-semibold text-foreground"
									>{getInitials(convo.participant_name)}</span
								>
							</div>
							<div class="flex-1 min-w-0">
								<div class="flex items-center justify-between mb-0.5">
									<p
										class="font-semibold text-sm text-foreground truncate"
									>
										{convo.participant_name}
									</p>
									<span
										class="text-xs text-muted-foreground flex-shrink-0 ml-2"
									>
										{formatRelative(convo.last_message_at)}
									</span>
								</div>
								<div class="flex items-center justify-between">
									<p
										class="text-xs text-muted-foreground truncate"
									>
										{convo.last_message}
									</p>
									{#if convo.unread_count > 0}
										<span
											class="ml-2 flex-shrink-0 w-5 h-5 rounded-full bg-foreground text-background text-xs flex items-center justify-center font-medium"
										>
											{convo.unread_count}
										</span>
									{/if}
								</div>
							</div>
						</div>
					</button>
				{/each}
			</div>

			<!-- Message Thread -->
			{#if activeThread}
				<div class="flex-1 flex flex-col overflow-hidden {activeThreadId ? 'block' : 'hidden lg:block'}">
					<!-- Thread Header -->
					<div
						class="px-6 py-4 border-b border-border flex items-center gap-3"
					>
						<button
							class="lg:hidden p-1 rounded-md hover:bg-muted transition-colors text-muted-foreground"
							onclick={goBack}
						>
							<ArrowLeft size={20} />
						</button>
						<div
							class="w-9 h-9 rounded-full bg-muted flex items-center justify-center"
						>
							<span class="text-sm font-semibold text-foreground"
								>{getInitials(activeThread.participant_name)}</span
							>
						</div>
						<div>
							<p class="font-semibold text-foreground">
								{activeThread.participant_name}
							</p>
							<p class="text-xs text-muted-foreground">Partenaire</p>
						</div>
					</div>

					<!-- Messages -->
					<div bind:this={messagesEl} class="flex-1 overflow-y-auto px-6 py-4 space-y-4">
						{#each pollingMessages as msg (msg.id)}
							<div
								class="flex {msg.sender_id === data.userId
									? 'justify-end'
									: 'justify-start'}"
							>
								<div
									class="max-w-xs lg:max-w-md px-4 py-2.5 rounded-2xl text-sm {msg.sender_id ===
									data.userId
										? 'bg-foreground text-background rounded-br-sm'
										: 'bg-muted text-foreground rounded-bl-sm'}"
								>
									<p>{msg.body}</p>
									<p
										class="text-xs mt-1 {msg.sender_id === data.userId
											? 'text-background/60'
											: 'text-muted-foreground'}"
									>
										{formatTime(msg.created_at)}
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
								onkeydown={handleKeydown}
								placeholder="Écrire un message…"
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
				<div
					class="flex-1 hidden lg:flex items-center justify-center text-muted-foreground"
				>
					Sélectionnez une conversation
				</div>
			{/if}
		</div>
	</div>
{/if}
