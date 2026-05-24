import { env } from '$env/dynamic/private';
import { error, isRedirect, redirect } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';

// API response types (snake_case from backend)
interface APICurrentMonth {
	revenue_cents: number;
	bookings_count: number;
	new_clients_count: number;
	avg_booking_value_cents: number;
}

interface APIMonthlyRevenue {
	month: string;
	revenue_cents: number;
	bookings_count: number;
}

interface APITopProduct {
	product_id: string;
	name: string;
	bookings_count: number;
	revenue_cents: number;
}

interface APIAnalyticsSummary {
	current_month: APICurrentMonth;
	monthly_revenue: APIMonthlyRevenue[];
	top_products: APITopProduct[];
}

// UI types (camelCase)
export interface DashboardStats {
	totalRevenueInCents: number;
	bookingsThisMonth: number;
	newUsersThisMonth: number;
	avgBookingValue: number;
}

export interface MonthlyRevenue {
	month: string;
	amountInCents: number;
}

export interface TopProduct {
	name: string;
	bookings: number;
	revenueInCents: number;
}

export interface AnalyticsData {
	stats: DashboardStats;
	monthlyRevenue: MonthlyRevenue[];
	topProducts: TopProduct[];
}

async function getMockAnalyticsData(): Promise<AnalyticsData> {
	const stats: DashboardStats = {
		totalRevenueInCents: 48750,
		bookingsThisMonth: 156,
		newUsersThisMonth: 23,
		avgBookingValue: 312
	};

	const monthlyRevenue: MonthlyRevenue[] = [
		{ month: 'Nov', amountInCents: 35000 },
		{ month: 'Déc', amountInCents: 42000 },
		{ month: 'Jan', amountInCents: 38500 },
		{ month: 'Fév', amountInCents: 45100 },
		{ month: 'Mar', amountInCents: 51200 },
		{ month: 'Avr', amountInCents: 48750 }
	];

	const topProducts: TopProduct[] = [
		{ name: 'Massage Relaxant 60min', bookings: 45, revenueInCents: 22500 },
		{ name: 'Soin du Dos 60min', bookings: 38, revenueInCents: 24700 },
		{ name: 'Drainage Lymphatique', bookings: 28, revenueInCents: 19600 },
		{ name: 'Consultation Kiné 45min', bookings: 25, revenueInCents: 8750 },
		{ name: 'Massage Relaxant 90min', bookings: 20, revenueInCents: 14000 }
	];

	return { stats, monthlyRevenue, topProducts };
}

export const load: PageServerLoad = async ({ fetch }) => {
	if (env.USE_MOCK_DATA === 'true') {
		return await getMockAnalyticsData();
	}

	try {
		const res = await fetch(`${env.API_URL}/admin/analytics/summary?months=6`);
		if (res.status === 401 || res.status === 403) {
			throw redirect(302, '/auth');
		}
		if (!res.ok) {
			throw new Error(`Failed to fetch analytics: ${res.status} ${res.statusText}`);
		}

		const api: APIAnalyticsSummary = await res.json();

		const stats: DashboardStats = {
			totalRevenueInCents: api.current_month.revenue_cents,
			bookingsThisMonth: api.current_month.bookings_count,
			newUsersThisMonth: api.current_month.new_clients_count,
			avgBookingValue: api.current_month.avg_booking_value_cents
		};

		const monthlyRevenue: MonthlyRevenue[] = (api.monthly_revenue ?? []).map((m) => ({
			month: formatMonthLabel(m.month),
			amountInCents: m.revenue_cents
		}));

		const topProducts: TopProduct[] = (api.top_products ?? []).map((p) => ({
			name: p.name,
			bookings: p.bookings_count,
			revenueInCents: p.revenue_cents
		}));

		return { stats, monthlyRevenue, topProducts };
	} catch (err) {
		if (isRedirect(err)) throw err;
		console.error('Error loading analytics data:', err);
		throw error(503, 'Impossible de charger les données analytiques. Veuillez réessayer.');
	}
};

/** Convert "2026-04" → "Avr" */
function formatMonthLabel(iso: string): string {
	const [year, month] = iso.split('-');
	const date = new Date(parseInt(year), parseInt(month) - 1, 1);
	return new Intl.DateTimeFormat('fr-FR', { month: 'short' }).format(date);
}
