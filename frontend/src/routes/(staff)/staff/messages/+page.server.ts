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

export interface BookingContext {
	id: string;
	product_name: string;
	slot_start_time: string;
	slot_end_time: string;
	status: string;
	total_price_cents: number;
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
		let bookingContext: BookingContext[] = [];

		if (activeThreadId) {
			const messagesRes = await fetch(`${env.API_URL}/threads/${activeThreadId}/messages?limit=100`, {
				headers: { 'Content-Type': 'application/json' }
			});
			if (messagesRes.ok) {
				const data = await messagesRes.json();
				activeMessages = data.messages ?? [];
			}

			// Mark thread as read
			fetch(`${env.API_URL}/threads/${activeThreadId}/read`, {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' }
			}).catch(() => {});

			// Fetch booking context for the active thread's participant
			const activeThread = threads.find(t => t.thread_id === activeThreadId);
			if (activeThread) {
				try {
					const bookingsRes = await fetch(
						`${env.API_URL}/partners/bookings/${locals.user!.id}`,
						{ headers: { 'Content-Type': 'application/json' } }
					);
					if (bookingsRes.ok) {
						const allBookings = await bookingsRes.json();
						// Filter bookings for the participant
						bookingContext = (Array.isArray(allBookings) ? allBookings : [])
							.filter((b: any) => b.client_id === activeThread.participant_id)
							.map((b: any) => ({
								id: b.id,
								product_name: b.product_name ?? 'Réservation',
								slot_start_time: b.slot_start_time,
								slot_end_time: b.slot_end_time,
								status: b.status,
								total_price_cents: b.total_price_cents
							}));
					}
				} catch {
					// Booking context is optional, don't fail
				}
			}
		}

		return {
			threads,
			activeThreadId: activeThreadId ?? null,
			activeMessages,
			bookingContext,
			userId: locals.user.id
		};
	} catch (err) {
		if (isRedirect(err)) throw err;
		console.error('Error loading messages:', err);
		return {
			threads: [],
			activeThreadId: null,
			activeMessages: [],
			bookingContext: [],
			userId: locals.user.id
		};
	}
};

function getMockData() {
	const now = new Date();
	const t = (offsetMin: number) => new Date(now.getTime() - offsetMin * 60 * 1000).toISOString();

	return {
		threads: [
			{
				thread_id: 'c1',
				participant_id: 'u1',
				participant_name: 'Marie Dupont',
				last_message: 'Parfait, merci pour la confirmation !',
				last_message_at: t(10),
				unread_count: 2
			},
			{
				thread_id: 'c2',
				participant_id: 'u2',
				participant_name: 'Jean Durand',
				last_message: 'Je serai peut-être un peu en retard demain.',
				last_message_at: t(90),
				unread_count: 1
			},
			{
				thread_id: 'c3',
				participant_id: 'u3',
				participant_name: 'Claire Bernard',
				last_message: 'À samedi alors !',
				last_message_at: t(240),
				unread_count: 0
			},
			{
				thread_id: 'c4',
				participant_id: 'u4',
				participant_name: 'Lucas Petit',
				last_message: "D'accord, je prendrai rendez-vous la semaine prochaine.",
				last_message_at: t(1440),
				unread_count: 0
			}
		],
		activeThreadId: null,
		activeMessages: [],
		bookingContext: [],
		userId: 'partner-mock-id'
	};
}
