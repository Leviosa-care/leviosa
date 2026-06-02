import { env } from '$env/dynamic/private';
import { redirect, isRedirect } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';

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
}

export const load: PageServerLoad = async ({ locals, fetch, url }) => {
	if (!locals.user?.id) {
		throw redirect(302, '/auth');
	}

	if (env.USE_MOCK_DATA === 'true') {
		return getMockData(locals.user.id);
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
		console.error('Error loading messages:', err);
		return {
			threads: [],
			activeThreadId: null,
			activeMessages: [],
			userId: locals.user.id
		};
	}
};

function getMockData(userId: string) {
	const now = new Date();
	const t = (offsetMin: number) =>
		new Date(now.getTime() - offsetMin * 60 * 1000).toISOString();

	return {
		threads: [
			{
				thread_id: 'c1',
				participant_id: 'p1',
				participant_name: 'Dr. Sophie Martin',
				last_message:
					'Votre rendez-vous de jeudi est confirmé. N\'oubliez pas d\'apporter vos documents.',
				last_message_at: t(10),
				unread_count: 2
			},
			{
				thread_id: 'c2',
				participant_id: 'p2',
				participant_name: 'Marc Lefevre',
				last_message: 'Merci pour votre message, à bientôt !',
				last_message_at: t(90),
				unread_count: 0
			}
		],
		activeThreadId: null as string | null,
		activeMessages: [] as ThreadMessage[],
		userId
	};
}
