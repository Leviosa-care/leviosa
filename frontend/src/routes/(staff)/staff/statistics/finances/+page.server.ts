import type { PageServerLoad } from './$types';
import { redirect } from '@sveltejs/kit';
import { env } from '$env/dynamic/private';

export interface Transaction {
	id: string;
	slotStartTime: string;
	productId: string;
	amountCents: number;
	paymentStatus: 'paid' | 'pending' | 'refunded';
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
	amount_cents: number;
	payment_status: 'paid' | 'pending' | 'refunded';
}

interface BackendEarningsSummary {
	current_month_cents: number;
	last_month_cents: number;
	pending_cents: number;
	next_payout_date: string;
	next_payout_cents: number;
	transactions: BackendTransaction[];
}

export const load: PageServerLoad = async ({ locals, fetch }) => {
	const partnerId = locals.user?.id;
	if (!partnerId) {
		throw redirect(302, '/auth');
	}

	const earningsRes = await fetch(`${env.API_URL}/partners/${partnerId}/earnings`, {
		headers: {
			'Content-Type': 'application/json',
		},
	});

	if (earningsRes.status === 401) {
		throw redirect(302, '/auth');
	}

	if (!earningsRes.ok) {
		if (earningsRes.status === 500) {
			return {
				summary: null,
				error: 'Erreur serveur. Veuillez réessayer dans quelques instants.',
			};
		}
		throw redirect(302, '/auth');
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
			amountCents: t.amount_cents,
			paymentStatus: t.payment_status,
		})),
	};

	return { summary };
};
