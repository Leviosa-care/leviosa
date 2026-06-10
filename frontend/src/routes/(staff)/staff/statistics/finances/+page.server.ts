import type { PageServerLoad } from './$types';
import { redirect } from '@sveltejs/kit';
import { env } from '$env/dynamic/private';

export interface Transaction {
	id: string;
	slotStartTime: string;
	productId: string;
	productName: string;
	amountCents: number;
	paymentStatus: 'paid' | 'pending' | 'refunded';
	bookingStatus: 'confirmed' | 'cancelled' | 'completed' | 'no_show';
}

export interface EarningsSummary {
	currentMonthCents: number;
	lastMonthCents: number;
	pendingCents: number;
	nextPayoutDate: string;
	nextPayoutCents: number;
	transactions: Transaction[];
}

interface BackendTransaction {
	id: string;
	slot_start_time: string;
	product_id: string;
	product_name: string;
	amount_cents: number;
	payment_status: 'paid' | 'pending' | 'refunded';
	booking_status: 'confirmed' | 'cancelled' | 'completed' | 'no_show';
}

interface BackendEarningsSummary {
	current_month_cents: number;
	last_month_cents: number;
	pending_cents: number;
	next_payout_date: string;
	next_payout_cents: number;
	transactions: BackendTransaction[];
}

function defaultMonth(): string {
	const now = new Date();
	return `${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, '0')}`;
}

export const load: PageServerLoad = async ({ locals, fetch, url }) => {
	const partnerId = locals.user?.id;
	if (!partnerId) {
		throw redirect(302, '/auth');
	}

	const monthParam = url.searchParams.get('month');
	const selectedMonth = monthParam && /^\d{4}-\d{2}$/.test(monthParam) ? monthParam : defaultMonth();

	const earningsRes = await fetch(`${env.API_URL}/partners/${partnerId}/earnings`, {
		headers: {
			'Content-Type': 'application/json',
		},
	});

	if (earningsRes.status === 401) {
		throw redirect(302, '/auth');
	}

	if (!earningsRes.ok) {
		return {
			summary: null,
			selectedMonth,
			error: earningsRes.status >= 500
				? 'Erreur serveur. Veuillez réessayer dans quelques instants.'
				: 'Impossible de charger vos données financières.',
		};
	}

	const backend: BackendEarningsSummary = await earningsRes.json();

	const summary: EarningsSummary = {
		currentMonthCents: backend.current_month_cents,
		lastMonthCents: backend.last_month_cents,
		pendingCents: backend.pending_cents,
		nextPayoutDate: backend.next_payout_date,
		nextPayoutCents: backend.next_payout_cents,
		transactions: backend.transactions.map((t) => ({
			id: t.id,
			slotStartTime: t.slot_start_time,
			productId: t.product_id,
			productName: t.product_name,
			amountCents: t.amount_cents,
			paymentStatus: t.payment_status,
			bookingStatus: t.booking_status,
		})),
	};

	return { summary, selectedMonth };
};
