<script lang="ts">
	import type { PageProps } from './$types';
	import type { ThreadSummary, ThreadMessage, BookingContext } from './+page.server';
	import { Send, ArrowLeft, Calendar, Clock } from '@lucide/svelte';
	import { browser } from '$app/environment';

	let { data }: PageProps = $props();

	let activeThreadId = $state<string | null>(data.activeThreadId);
	let newMessage = $state('');
	let sending = $state(false);
	let threads = $state<ThreadSummary[]>([...data.threads]);
	let pollingMessages = $state<ThreadMessage[]>([...data.activeMessages]);

	const activeThread = $derived(threads.find((t) => t.thread_id === activeThreadId) ?? null);

	const totalUnread = $derived(threads.reduce((s, c) => s + c.unread_count, 0));

	// Bookings are already filtered server-side for the active participant.
	const upcomingBookings = $derived(
		data.bookingContext.filter(
			(b) =>
				(b.status === 'confirmed' || b.status === 'upcoming') &&
				new Date(b.slot_start_time) > new Date()
		)
	);

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

	function formatRelative(iso: string): string {
		const diff = (Date.now() - new Date(iso).getTime()) / 1000;
		if (diff < 60) return "À l'instant";
		if (diff < 3600) return `Il y a ${Math.floor(diff / 60)} min`;
		if (diff < 86400) return `Il y a ${Math.floor(diff / 3600)} h`;
		return new Date(iso).toLocaleDateString('fr-FR', { day: 'numeric', month: 'short' });
	}

	function formatTime(iso: string): string {
		return new Date(iso).toLocaleTimeString('fr-FR', { hour: '2-digit', minute: '2-digit' });
	}

	function formatDate(iso: string): string {
		return new Date(iso).toLocaleDateString('fr-FR', {
			day: 'numeric',
			month: 'long',
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

	function selectThread(threadId: string) {
		activeThreadId = threadId;
		pollingMessages = [];
		const url = new URL(window.location.href);
		url.searchParams.set('thread', threadId);
		window.location.href = url.toString();
	}

	function handleKeydown(e: KeyboardEvent) {
		if (e.key === 'Enter' && !e.shiftKey) {
			e.preventDefault();
			sendMessage();
		}
	}
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
								<p class="font-semibold text-sm text-foreground truncate">
									{convo.participant_name}
								</p>
								<span class="text-xs text-muted-foreground flex-shrink-0 ml-2">
									{formatRelative(convo.last_message_at)}
								</span>
							</div>
							<div class="flex items-center justify-between">
								<p class="text-xs text-muted-foreground truncate">{convo.last_message}</p>
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
			<div class="flex-1 flex overflow-hidden">
				<!-- Main message area -->
				<div
					class="flex-1 flex flex-col overflow-hidden {activeThreadId ? 'block' : 'hidden lg:block'}"
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
							<span class="text-sm font-semibold text-foreground"
								>{getInitials(activeThread.participant_name)}</span
							>
						</div>
						<div>
							<p class="font-semibold text-foreground">{activeThread.participant_name}</p>
							<p class="text-xs text-muted-foreground">Client</p>
						</div>
					</div>

					<!-- Messages -->
					<div class="flex-1 overflow-y-auto px-6 py-4 space-y-4">
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

				<!-- Booking Context Sidebar (desktop only) -->
				{#if upcomingBookings.length > 0}
					<div class="hidden xl:block w-72 border-l border-border overflow-y-auto bg-muted/20">
						<div class="p-4">
							<h3 class="text-sm font-semibold text-foreground mb-3 flex items-center gap-2">
								<Calendar size={14} />
								Prochaines séances
							</h3>
							{#each upcomingBookings as booking (booking.id)}
								<div class="p-3 rounded-lg border border-border bg-background mb-2">
									<p class="text-sm font-medium text-foreground">
										{booking.product_name}
									</p>
									<p class="text-xs text-muted-foreground flex items-center gap-1 mt-1">
										<Clock size={12} />
										{formatDate(booking.slot_start_time)}
									</p>
									<p class="text-xs text-muted-foreground mt-0.5">
										Statut: {booking.status === 'confirmed' ? 'Confirmé' : booking.status}
									</p>
								</div>
							{/each}
						</div>
					</div>
				{/if}
			</div>
		{:else}
			<div class="flex-1 hidden lg:flex items-center justify-center text-muted-foreground">
				Sélectionnez une conversation
			</div>
		{/if}
	</div>
</div>
