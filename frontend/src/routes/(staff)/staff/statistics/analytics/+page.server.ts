import type { PageServerLoad } from './$types';

export interface PartnerStats {
	sessionsThisMonth: number;
	sessionsLastMonth: number;
	attendanceRate: number;
	utilizationRate: number;
	avgSessionDurationMin: number;
	uniqueClientsThisMonth: number;
}

export interface WeeklyVolume {
	week: string;
	sessions: number;
}

export interface TopService {
	name: string;
	sessions: number;
	percentage: number;
}

function getMockAnalytics() {
	const stats: PartnerStats = {
		sessionsThisMonth: 38,
		sessionsLastMonth: 32,
		attendanceRate: 92,
		utilizationRate: 74,
		avgSessionDurationMin: 68,
		uniqueClientsThisMonth: 24,
	};

	const weeklyVolume: WeeklyVolume[] = [
		{ week: 'S-5', sessions: 7 },
		{ week: 'S-4', sessions: 9 },
		{ week: 'S-3', sessions: 8 },
		{ week: 'S-2', sessions: 11 },
		{ week: 'S-1', sessions: 10 },
		{ week: 'Cette sem.', sessions: 6 },
	];

	const topServices: TopService[] = [
		{ name: 'Massage Relaxant 60min', sessions: 16, percentage: 42 },
		{ name: 'Drainage Lymphatique 90min', sessions: 10, percentage: 26 },
		{ name: 'Soin du Dos 60min', sessions: 7, percentage: 18 },
		{ name: 'Massage Sportif 60min', sessions: 5, percentage: 13 },
	];

	return { stats, weeklyVolume, topServices };
}

export const load: PageServerLoad = async () => {
	// TODO: Replace with GET /partners/{partnerId}/metrics
	return getMockAnalytics();
};
