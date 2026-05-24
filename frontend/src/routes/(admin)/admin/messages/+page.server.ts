import { env } from '$env/dynamic/private';
import { redirect, isRedirect } from '@sveltejs/kit';
import type { PageServerLoad, Actions } from './$types';

export interface ThreadMessage {
	id: string;
	thread_id: string;
	sender_id: string;
	body: string;
	created_at: string;
	read_at: string | null;
}

export interface ThreadSummary {
	thread_id: string;
	participant_id: string;
	participant_name: string;
	last_message: string;
	last_message_at: string;
	unread_count: number;
	is_online: boolean;
}

export interface UserSearchResult {
	id: string;
	firstname: string;
	lastname: string;
	email: string;
}

export const load: PageServerLoad = async ({ locals, fetch, url }) => {
	if (!locals.user?.id) {
		throw redirect(302, '/auth');
	}

	if (env.USE_MOCK_DATA === 'true') {
		return getMockData();
	}

	try {
		// Fetch threads
		const threadsRes = await fetch(`${env.API_URL}/threads`, {
			headers: { 'Content-Type': 'application/json' }
		});

		if (threadsRes.status === 401) {
			throw redirect(302, '/auth');
		}

		let threads: ThreadSummary[] = [];
		if (threadsRes.ok) {
			threads = await threadsRes.json();
		}

		// If a thread is selected via query param, fetch its messages
		const activeThreadId = url.searchParams.get('thread');
		let activeMessages: ThreadMessage[] = [];

		if (activeThreadId) {
			const messagesRes = await fetch(
				`${env.API_URL}/threads/${activeThreadId}/messages?limit=100`,
				{ headers: { 'Content-Type': 'application/json' } }
			);
			if (messagesRes.ok) {
				const data = await messagesRes.json();
				activeMessages = data.messages ?? [];
			}

			// Mark thread as read (fire and forget)
			fetch(`${env.API_URL}/threads/${activeThreadId}/read`, {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' }
			}).catch(() => {});
		}

		return {
			threads,
			activeThreadId: activeThreadId ?? null,
			activeMessages,
			userId: locals.user.id
		};
	} catch (err) {
		if (isRedirect(err)) throw err;
		console.error('Error loading admin messages:', err);
		return {
			threads: [],
			activeThreadId: null,
			activeMessages: [],
			userId: locals.user.id
		};
	}
};

export const actions: Actions = {
	createThread: async ({ request, locals, fetch }) => {
		if (!locals.user?.id) return { success: false };

		const formData = await request.formData();
		const participantId = formData.get('participant_id') as string;
		if (!participantId) return { success: false };

		try {
			const res = await fetch(`${env.API_URL}/threads`, {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ participant_id: participantId })
			});

			if (res.ok) {
				const data = await res.json();
				return { success: true, threadId: data.id };
			}
			return { success: false };
		} catch {
			return { success: false };
		}
	}
};

function getMockData() {
	const now = new Date();
	const t = (offsetMin: number) => new Date(now.getTime() - offsetMin * 60 * 1000).toISOString();

	return {
		threads: [
			{
				thread_id: '1',
				participant_id: 'u1',
				participant_name: 'Marie Dupont',
				last_message: 'Merci pour le rendez-vous, à bientôt!',
				last_message_at: t(5),
				unread_count: 0,
				is_online: true
			},
			{
				thread_id: '2',
				participant_id: 'u2',
				participant_name: 'Jean Durand',
				last_message: 'Quelle est la durée de la séance?',
				last_message_at: t(30),
				unread_count: 2,
				is_online: false
			},
			{
				thread_id: '3',
				participant_id: 'u3',
				participant_name: 'Claire Bernard',
				last_message: "D'accord, merci!",
				last_message_at: t(120),
				unread_count: 0,
				is_online: true
			},
			{
				thread_id: '4',
				participant_id: 'u4',
				participant_name: 'Lucas Petit',
				last_message: 'Les prix des massages?',
				last_message_at: t(1440),
				unread_count: 1,
				is_online: false
			},
			{
				thread_id: '5',
				participant_id: 'u5',
				participant_name: 'Emma Moreau',
				last_message: 'Parfait, je confirme',
				last_message_at: t(4320),
				unread_count: 0,
				is_online: false
			}
		],
		activeThreadId: null as string | null,
		activeMessages: [] as ThreadMessage[],
		userId: 'admin-mock-id'
	};
}
