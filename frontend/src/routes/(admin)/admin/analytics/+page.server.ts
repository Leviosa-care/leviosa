import { env } from '$env/dynamic/private';
import type { PageServerLoad } from './$types';

interface DashboardStats {
	totalRevenueInCents: number;
	bookingsThisMonth: number;
	newUsersThisMonth: number;
	avgBookingValue: number;
	repeatRate: number;
	conversionRate: number;
}

interface MonthlyRevenue {
	month: string;
	amountInCents: number;
}

interface TopProduct {
	name: string;
	bookings: number;
	revenueInCents: number;
}

interface AnalyticsData {
	stats: DashboardStats;
	monthlyRevenue: MonthlyRevenue[];
	topProducts: TopProduct[];
}

async function getMockAnalyticsData(): Promise<AnalyticsData> {
	const stats: DashboardStats = {
		totalRevenueInCents: 48750,
		bookingsThisMonth: 156,
		newUsersThisMonth: 23,
		avgBookingValue: 312,
		repeatRate: 42,
		conversionRate: 18
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

export const load: PageServerLoad = async () => {
	if (env.USE_MOCK_DATA === 'true') {
		return await getMockAnalyticsData();
	}
	return {
		stats: {
			totalRevenueInCents: 0,
			bookingsThisMonth: 0,
			newUsersThisMonth: 0,
			avgBookingValue: 0,
			repeatRate: 0,
			conversionRate: 0,
		},
		monthlyRevenue: [],
		topProducts: [],
	};
};
