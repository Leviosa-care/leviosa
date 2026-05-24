import { env } from '$env/dynamic/private';
import { error, isRedirect, redirect } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';

// --- API response types (snake_case from backend) ---

interface APISummary {
	gross_revenue_cents: number;
	refunds_cents: number;
	net_revenue_cents: number;
}

interface APITransaction {
	id: string;
	slot_start_time: string;
	client_name: string;
	partner_name: string;
	product_name: string;
	amount_cents: number;
	payment_status: 'paid' | 'refunded';
	booking_status: 'completed' | 'cancelled' | 'confirmed' | 'no_show';
}

interface APIFinancialSummary {
	summary: APISummary;
	transactions: APITransaction[];
}

// --- UI types (camelCase) ---

export interface SummaryKPIs {
	grossRevenueInCents: number;
	refundsInCents: number;
	netRevenueInCents: number;
}

export interface Transaction {
	id: string;
	date: string;
	description: string;
	clientName: string;
	partnerName: string;
	amountInCents: number;
	type: 'payment' | 'refund';
	paymentStatus: 'paid' | 'refunded';
	bookingStatus: string;
}

export interface ComptaData {
	summary: SummaryKPIs;
	transactions: Transaction[];
	from: string;
	to: string;
}

async function getMockComptaData(from: string, to: string): Promise<ComptaData> {
	const now = new Date();

	const summary: SummaryKPIs = {
		grossRevenueInCents: 52300,
		refundsInCents: 2150,
		netRevenueInCents: 50150
	};

	const transactions: Transaction[] = [
		{ id: 'TXN-001', date: new Date(now.getTime() - 2 * 60 * 60 * 1000).toISOString(), description: 'Massage Relaxant 60min', clientName: 'Marie Dupont', partnerName: 'Dr. Martin', amountInCents: 5000, type: 'payment', paymentStatus: 'paid', bookingStatus: 'completed' },
		{ id: 'TXN-002', date: new Date(now.getTime() - 4 * 60 * 60 * 1000).toISOString(), description: 'Soin du Dos 60min', clientName: 'Jean Durand', partnerName: 'Dr. Martin', amountInCents: 6500, type: 'payment', paymentStatus: 'paid', bookingStatus: 'completed' },
		{ id: 'TXN-003', date: new Date(now.getTime() - 6 * 60 * 60 * 1000).toISOString(), description: 'Remboursement', clientName: 'Claire Bernard', partnerName: 'Dr. Dupont', amountInCents: 5000, type: 'refund', paymentStatus: 'refunded', bookingStatus: 'cancelled' },
		{ id: 'TXN-004', date: new Date(now.getTime() - 24 * 60 * 60 * 1000).toISOString(), description: 'Drainage Lymphatique', clientName: 'Lucas Petit', partnerName: 'Dr. Dupont', amountInCents: 7000, type: 'payment', paymentStatus: 'paid', bookingStatus: 'completed' },
		{ id: 'TXN-005', date: new Date(now.getTime() - 48 * 60 * 60 * 1000).toISOString(), description: 'Consultation Kiné 45min', clientName: 'Thomas Richard', partnerName: 'Dr. Martin', amountInCents: 3500, type: 'payment', paymentStatus: 'paid', bookingStatus: 'completed' },
		{ id: 'TXN-006', date: new Date(now.getTime() - 72 * 60 * 60 * 1000).toISOString(), description: 'Massage Relaxant 60min', clientName: 'Camille Simon', partnerName: 'Dr. Dupont', amountInCents: 5000, type: 'payment', paymentStatus: 'paid', bookingStatus: 'completed' },
		{ id: 'TXN-007', date: new Date(now.getTime() - 96 * 60 * 60 * 1000).toISOString(), description: 'Remboursement', clientName: 'Louis Laurent', partnerName: 'Dr. Martin', amountInCents: 6500, type: 'refund', paymentStatus: 'refunded', bookingStatus: 'cancelled' },
		{ id: 'TXN-008', date: new Date(now.getTime() - 120 * 60 * 60 * 1000).toISOString(), description: 'Massage Relaxant 90min', clientName: 'Lola Bernard', partnerName: 'Dr. Dupont', amountInCents: 7000, type: 'payment', paymentStatus: 'paid', bookingStatus: 'completed' },
	];

	return { summary, transactions, from, to };
}

export const load: PageServerLoad = async ({ fetch, url }) => {
	const from = url.searchParams.get('from') || defaultFrom();
	const to = url.searchParams.get('to') || defaultTo();

	if (env.USE_MOCK_DATA === 'true') {
		return await getMockComptaData(from, to);
	}

	try {
		const res = await fetch(`${env.API_URL}/admin/bookings/financial-summary?from=${from}&to=${to}`);
		if (res.status === 401 || res.status === 403) {
			throw redirect(302, '/auth');
		}
		if (!res.ok) {
			throw new Error(`Failed to fetch financial summary: ${res.status} ${res.statusText}`);
		}

		const api: APIFinancialSummary = await res.json();

		const summary: SummaryKPIs = {
			grossRevenueInCents: api.summary.gross_revenue_cents,
			refundsInCents: api.summary.refunds_cents,
			netRevenueInCents: api.summary.net_revenue_cents
		};

		const transactions: Transaction[] = (api.transactions ?? []).map((t) => ({
			id: t.id,
			date: t.slot_start_time,
			description: t.product_name,
			clientName: t.client_name,
			partnerName: t.partner_name,
			amountInCents: t.amount_cents,
			type: t.payment_status === 'refunded' ? 'refund' as const : 'payment' as const,
			paymentStatus: t.payment_status,
			bookingStatus: t.booking_status
		}));

		return { summary, transactions, from, to };
	} catch (err) {
		if (isRedirect(err)) throw err;
		console.error('Error loading financial summary:', err);
		throw error(503, 'Impossible de charger les données comptables. Veuillez réessayer.');
	}
};

function defaultFrom(): string {
	const now = new Date();
	return `${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, '0')}-01`;
}

function defaultTo(): string {
	const now = new Date();
	return `${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, '0')}-${String(now.getDate()).padStart(2, '0')}`;
}
