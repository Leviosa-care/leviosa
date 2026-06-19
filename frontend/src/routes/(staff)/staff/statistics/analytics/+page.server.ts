import type { PageServerLoad } from './$types';
import { redirect } from '@sveltejs/kit';
import { env } from '$env/dynamic/private';

export interface PartnerMetrics {
	partnerId: string;
	startDate: string;
	endDate: string;
	summary: {
		averageUtilization: number;
		totalFragmentation: number;
		totalIdleMinutes: number;
		averageEfficiency: number;
		daysAnalyzed: number;
	};
	roomMetrics: Array<{
		roomId: string;
		startDate: string;
		endDate: string;
		summary: {
			averageUtilization: number;
			totalFragmentation: number;
			totalIdleMinutes: number;
			averageEfficiency: number;
			daysAnalyzed: number;
		};
		dailyMetrics: Array<{
			date: string;
			totalMinutesOpen: number;
			totalMinutesBooked: number;
			utilizationPercent: number;
			fragmentationCount: number;
			idleMinutes: number;
			averageGapMinutes: number;
			efficiencyScore: number;
		}>;
	}>;
}

interface BackendDailyMetrics {
	date: string;
	total_minutes_open: number;
	total_minutes_booked: number;
	utilization_percent: number;
	fragmentation_count: number;
	idle_minutes: number;
	average_gap_minutes: number;
	efficiency_score: number;
}

interface BackendMetricsSummary {
	average_utilization: number;
	total_fragmentation: number;
	total_idle_minutes: number;
	average_efficiency: number;
	days_analyzed: number;
}

interface BackendRoomMetrics {
	room_id: string;
	start_date: string;
	end_date: string;
	daily_metrics: BackendDailyMetrics[];
	summary: BackendMetricsSummary;
}

interface BackendPartnerMetrics {
	partner_id: string;
	start_date: string;
	end_date: string;
	room_metrics: BackendRoomMetrics[];
	summary: BackendMetricsSummary;
}

function mapMetrics(b: BackendPartnerMetrics): PartnerMetrics {
	return {
		partnerId: b.partner_id,
		startDate: b.start_date,
		endDate: b.end_date,
		summary: {
			averageUtilization: b.summary.average_utilization,
			totalFragmentation: b.summary.total_fragmentation,
			totalIdleMinutes: b.summary.total_idle_minutes,
			averageEfficiency: b.summary.average_efficiency,
			daysAnalyzed: b.summary.days_analyzed,
		},
		roomMetrics: b.room_metrics.map((rm) => ({
			roomId: rm.room_id,
			startDate: rm.start_date,
			endDate: rm.end_date,
			summary: {
				averageUtilization: rm.summary.average_utilization,
				totalFragmentation: rm.summary.total_fragmentation,
				totalIdleMinutes: rm.summary.total_idle_minutes,
				averageEfficiency: rm.summary.average_efficiency,
				daysAnalyzed: rm.summary.days_analyzed,
			},
			dailyMetrics: rm.daily_metrics.map((d) => ({
				date: d.date,
				totalMinutesOpen: d.total_minutes_open,
				totalMinutesBooked: d.total_minutes_booked,
				utilizationPercent: d.utilization_percent,
				fragmentationCount: d.fragmentation_count,
				idleMinutes: d.idle_minutes,
				averageGapMinutes: d.average_gap_minutes,
				efficiencyScore: d.efficiency_score,
			})),
		})),
	};
}

export const load: PageServerLoad = async ({ locals, fetch, url }) => {
	const partnerId = locals.user?.id;
	if (!partnerId) {
		throw redirect(302, '/auth');
	}

	// Parse date range from URL query params, default to last 30 days
	const now = new Date();
	const defaultStart = new Date(now.getFullYear(), now.getMonth(), now.getDate() - 29);
	const startDateParam = url.searchParams.get('start_date') ?? defaultStart.toISOString().split('T')[0];
	const endDateParam = url.searchParams.get('end_date') ?? now.toISOString().split('T')[0];

	// Validate the date format (YYYY-MM-DD)
	const dateRegex = /^\d{4}-\d{2}-\d{2}$/;
	if (!dateRegex.test(startDateParam) || !dateRegex.test(endDateParam)) {
		return {
			metrics: null,
			startDate: startDateParam,
			endDate: endDateParam,
			error: 'Format de date invalide. Utilisez AAAA-MM-JJ.',
		};
	}

	// Validate start <= end
	if (startDateParam > endDateParam) {
		return {
			metrics: null,
			startDate: startDateParam,
			endDate: endDateParam,
			error: 'La date de début doit être antérieure ou égale à la date de fin.',
		};
	}

	if (env.USE_MOCK_DATA === 'true') {
		return { metrics: null, startDate: startDateParam, endDate: endDateParam };
	}

	const metricsRes = await fetch(`${env.API_URL}/partners/metrics/${partnerId}?start_date=${startDateParam}&end_date=${endDateParam}`, {
		headers: {
			'Content-Type': 'application/json',
		},
	});

	if (metricsRes.status === 401) {
		throw redirect(302, '/auth');
	}

	if (!metricsRes.ok) {
		return {
			metrics: null,
			startDate: startDateParam,
			endDate: endDateParam,
			error:
				metricsRes.status === 500
					? 'Erreur serveur. Veuillez réessayer dans quelques instants.'
					: 'Impossible de charger les statistiques.',
		};
	}

	const backend: BackendPartnerMetrics = await metricsRes.json();

	return {
		metrics: mapMetrics(backend),
		startDate: startDateParam,
		endDate: endDateParam,
	};
};
