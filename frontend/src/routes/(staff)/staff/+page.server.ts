import type { PageServerLoad } from './$types';
import { redirect, isRedirect } from '@sveltejs/kit';
import { env } from '$env/dynamic/private';

// ── Types ────────────────────────────────────────────────────────

export interface TodaySlot {
	id: string;
	startTime: string;
	endTime: string;
	productName: string;
	clientName: string;
	status: string;
}

export interface KpiData {
	revenueCents: number;
	revenueGrowthPct: string; // e.g. "+12,5%" or "—"
	bookingsCount: number;
	bookingsGrowthPct: string;
	occupationPct: string; // e.g. "94%" or "—"
}

export interface ActivityEvent {
	title: string;
	subtitle: string;
	relativeTime: string;
	colorKey: 'blue' | 'green' | 'purple' | 'amber' | 'red';
}

export interface VolumeDay {
	label: string; // single-letter day abbreviation
	pct: number; // 0–100 normalised to the week's maximum
}

export interface ActionItem {
	title: string;
	description: string;
	age: string;
	href: string;
}

// ── Backend shapes ───────────────────────────────────────────────

interface BackendBooking {
	id: string;
	client_name: string;
	product_name: string;
	slot_start_time: string;
	slot_end_time: string;
	status: string;
	payment_status: string;
	total_price_cents: number;
	currency: string;
}

interface BackendDailyMetrics {
	date: string;
	total_minutes_booked: number;
	utilization_percent: number;
}

interface BackendRoomMetrics {
	daily_metrics: BackendDailyMetrics[];
}

interface BackendPartnerMetrics {
	room_metrics: BackendRoomMetrics[];
	summary: {
		average_utilization: number;
	};
}

// ── Helpers ──────────────────────────────────────────────────────

const DAY_LABELS = ['L', 'M', 'M', 'J', 'V', 'S', 'D'] as const;

function isToday(isoDate: string): boolean {
	const d = new Date(isoDate);
	const now = new Date();
	return (
		d.getFullYear() === now.getFullYear() &&
		d.getMonth() === now.getMonth() &&
		d.getDate() === now.getDate()
	);
}

function formatTime(iso: string): string {
	return new Date(iso).toLocaleTimeString('fr-FR', { hour: '2-digit', minute: '2-digit' });
}

/** ISO date string (YYYY-MM-DD) for a Date offset by `days` from today. */
function dateOffset(days: number): string {
	const d = new Date();
	d.setDate(d.getDate() + days);
	return d.toISOString().split('T')[0];
}

/** Check if an ISO datetime falls within [startDay, endDay] inclusive (date-only comparison). */
function isInPeriod(isoDateTime: string, startDay: string, endDay: string): boolean {
	const d = new Date(isoDateTime).toISOString().split('T')[0];
	return d >= startDay && d <= endDay;
}

/** Format a growth percentage as a French-localised string, or "—" when not applicable. */
function formatGrowth(current: number, previous: number): string {
	if (previous === 0) return '—';
	const pct = ((current - previous) / previous) * 100;
	const sign = pct >= 0 ? '+' : '';
	return `${sign}${pct.toLocaleString('fr-FR', { maximumFractionDigits: 1 })}%`;
}

/** Relative time string in French (e.g. "il y a 5 min", "il y a 2h", "il y a 3j"). */
function relativeTime(isoDate: string): string {
	const now = Date.now();
	const then = new Date(isoDate).getTime();
	const diffMs = now - then;
	if (diffMs < 0) return 'à venir';

	const minutes = Math.floor(diffMs / 60_000);
	if (minutes < 1) return "il y a moins d'1 min";
	if (minutes < 60) return `il y a ${minutes} min`;

	const hours = Math.floor(minutes / 60);
	if (hours < 24) return `il y a ${hours}h`;

	const days = Math.floor(hours / 24);
	return `il y a ${days}j`;
}

/** Build activity events from recent bookings, most recent first. */
function buildActivityEvents(bookings: BackendBooking[], limit: number): ActivityEvent[] {
	// Sort by slot_start_time descending (most recent first)
	const sorted = [...bookings].sort(
		(a, b) => new Date(b.slot_start_time).getTime() - new Date(a.slot_start_time).getTime()
	);

	const events: ActivityEvent[] = [];

	for (const b of sorted) {
		if (events.length >= limit) break;

		const name = b.client_name || 'Client inconnu';
		const relTime = relativeTime(b.slot_start_time);

		if (b.payment_status === 'paid') {
			const amount = (b.total_price_cents / 100).toLocaleString('fr-FR', {
				minimumFractionDigits: 2,
				maximumFractionDigits: 2,
			});
			events.push({
				title: 'Paiement reçu',
				subtitle: `de ${amount} €`,
				relativeTime: relTime,
				colorKey: 'green',
			});
		}

		if (events.length >= limit) break;

		// Map booking status to an activity event
		if (b.status === 'cancelled') {
			events.push({
				title: 'Annulation',
				subtitle: `par ${name}`,
				relativeTime: relTime,
				colorKey: 'red',
			});
		} else if (b.status === 'no_show') {
			events.push({
				title: 'Absence',
				subtitle: name,
				relativeTime: relTime,
				colorKey: 'amber',
			});
		} else {
			events.push({
				title: 'Nouvelle réservation',
				subtitle: `par ${name}`,
				relativeTime: relTime,
				colorKey: 'blue',
			});
		}
	}

	return events;
}

/** Build 7-day volume chart data from metrics daily data. */
function buildVolumeDays(metrics: BackendPartnerMetrics | null): VolumeDay[] {
	// Initialise 7 days with zero values and proper labels
	const today = new Date();
	const days: { date: string; label: string; booked: number }[] = [];

	for (let i = 6; i >= 0; i--) {
		const d = new Date(today);
		d.setDate(d.getDate() - i);
		const iso = d.toISOString().split('T')[0];
		const dayIdx = d.getDay(); // 0=Sun … 6=Sat
		// Map to Mon=0 … Sun=6 to match DAY_LABELS
		const labelIdx = dayIdx === 0 ? 6 : dayIdx - 1;
		days.push({ date: iso, label: DAY_LABELS[labelIdx], booked: 0 });
	}

	if (!metrics) return days.map((d) => ({ label: d.label, pct: 0 }));

	// Collect daily total_minutes_booked from all rooms
	const dailyMap = new Map<string, number>();
	for (const rm of metrics.room_metrics ?? []) {
		for (const dm of rm.daily_metrics ?? []) {
			const key = new Date(dm.date).toISOString().split('T')[0];
			dailyMap.set(key, (dailyMap.get(key) ?? 0) + dm.total_minutes_booked);
		}
	}

	for (const day of days) {
		day.booked = dailyMap.get(day.date) ?? 0;
	}

	const maxBooked = Math.max(...days.map((d) => d.booked), 0);

	return days.map((d) => ({
		label: d.label,
		pct: maxBooked > 0 ? Math.round((d.booked / maxBooked) * 100) : 0,
	}));
}

// ── Load ─────────────────────────────────────────────────────────

export const load: PageServerLoad = async ({ locals, fetch }) => {
	if (!locals.user?.id) {
		throw redirect(302, '/auth');
	}

	const partnerId = locals.user.id;
	const todaySlots: TodaySlot[] = [];
	let kpi: KpiData = {
		revenueCents: 0,
		revenueGrowthPct: '—',
		bookingsCount: 0,
		bookingsGrowthPct: '—',
		occupationPct: '—',
	};
	let activityEvents: ActivityEvent[] = [];
	let volumeDays: VolumeDay[] = DAY_LABELS.map((label) => ({ label, pct: 0 }));

	if (env.USE_MOCK_DATA === 'true') {
		return { todaySlots, kpi, activityEvents, volumeDays, actions: [] as ActionItem[] };
	}

	// ── Fetch partner profile (for Stripe onboarding status) ─────
	let stripeOnboardingComplete = true;

	try {
		const profileRes = await fetch(`${env.API_URL}/partners/me`);

		if (profileRes.status === 401) {
			throw redirect(302, '/auth');
		}

		if (profileRes.ok) {
			const profileData = await profileRes.json();
			stripeOnboardingComplete = profileData.stripe_onboarding_complete ?? true;
		} else {
			console.error(`Failed to fetch partner profile: ${profileRes.status} ${profileRes.statusText}`);
		}
	} catch (err) {
		if (isRedirect(err)) throw err;
		console.error('Error loading partner profile:', err);
	}

	// ── Fetch bookings ────────────────────────────────────────────
	let allBookings: BackendBooking[] = [];

	try {
		const res = await fetch(`${env.API_URL}/partners/${partnerId}/bookings`);

		if (res.status === 401) {
			throw redirect(302, '/auth');
		}

		if (res.ok) {
			const data = await res.json();
			if (Array.isArray(data)) {
				allBookings = data as BackendBooking[];
			}
		} else {
			console.error(`Failed to fetch bookings: ${res.status} ${res.statusText}`);
		}
	} catch (err) {
		if (isRedirect(err)) throw err;
		console.error('Error loading bookings:', err);
	}

	// ── Today's agenda slots ──────────────────────────────────────
	const todayBookings = allBookings
		.filter((b) => isToday(b.slot_start_time))
		.sort(
			(a, b) =>
				new Date(a.slot_start_time).getTime() - new Date(b.slot_start_time).getTime()
		);

	for (const b of todayBookings) {
		todaySlots.push({
			id: b.id,
			startTime: formatTime(b.slot_start_time),
			endTime: formatTime(b.slot_end_time),
			productName: b.product_name || 'Produit inconnu',
			clientName: b.client_name || 'Client inconnu',
			status: b.status || 'confirmed',
		});
	}

	// ── KPI: revenue & bookings for current & previous 7-day periods ─
	const currentStart = dateOffset(-6);
	const currentEnd = dateOffset(0);
	const prevStart = dateOffset(-13);
	const prevEnd = dateOffset(-7);

	const currentBookings = allBookings.filter((b) => isInPeriod(b.slot_start_time, currentStart, currentEnd));
	const prevBookings = allBookings.filter((b) => isInPeriod(b.slot_start_time, prevStart, prevEnd));

	const currentRevenue = currentBookings
		.filter((b) => b.payment_status === 'paid')
		.reduce((sum, b) => sum + (b.total_price_cents || 0), 0);
	const prevRevenue = prevBookings
		.filter((b) => b.payment_status === 'paid')
		.reduce((sum, b) => sum + (b.total_price_cents || 0), 0);

	kpi = {
		revenueCents: currentRevenue,
		revenueGrowthPct: formatGrowth(currentRevenue, prevRevenue),
		bookingsCount: currentBookings.length,
		bookingsGrowthPct: formatGrowth(currentBookings.length, prevBookings.length),
		occupationPct: '—', // will be filled from metrics below
	};

	// ── Activity feed (last 5 bookings) ──────────────────────────
	activityEvents = buildActivityEvents(allBookings, 5);

	// ── Fetch 7-day metrics ──────────────────────────────────────
	let metrics: BackendPartnerMetrics | null = null;

	try {
		const mRes = await fetch(
			`${env.API_URL}/partners/${partnerId}/metrics?start_date=${currentStart}&end_date=${currentEnd}`
		);

		if (mRes.status === 401) {
			throw redirect(302, '/auth');
		}

		if (mRes.ok) {
			metrics = await mRes.json();
		} else {
			console.error(`Failed to fetch metrics: ${mRes.status} ${mRes.statusText}`);
		}
	} catch (err) {
		if (isRedirect(err)) throw err;
		console.error('Error loading metrics:', err);
	}

	// ── Occupation rate from metrics ─────────────────────────────
	if (metrics?.summary?.average_utilization != null) {
		kpi.occupationPct = `${Math.round(metrics.summary.average_utilization)}%`;
	}

	// ── Volume chart data ────────────────────────────────────────
	volumeDays = buildVolumeDays(metrics);

	// ── Build action items ───────────────────────────────────────
	const actions: ActionItem[] = [];

	if (!stripeOnboardingComplete) {
		actions.push({
			title: 'Configuration Stripe requise',
			description: 'Complétez votre inscription Stripe pour recevoir vos paiements.',
			age: 'Maintenant',
			href: '/staff/profile',
		});
	}

	const pendingBookingsCount = allBookings.filter(
		(b) => b.status === 'pending'
	).length;

	if (pendingBookingsCount > 0) {
		const plural = pendingBookingsCount > 1 ? 's' : '';
		actions.push({
			title: 'Réservations en attente',
			description: `${pendingBookingsCount} réservation${plural} nécessite${plural} votre confirmation.`,
			age: 'À confirmer',
			href: '/staff/agenda/reservations',
		});
	}

	return { todaySlots, kpi, activityEvents, volumeDays, actions };
};
