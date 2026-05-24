import type { PageServerLoad } from './$types';
import { redirect, isRedirect } from '@sveltejs/kit';
import { env } from '$env/dynamic/private';

export interface Availability {
	id: string;
	date: string;
	startTime: string;
	endTime: string;
	status: 'available' | 'booked' | 'cancelled' | 'blocked';
	roomId: string;
	roomName: string;
	maxCapacity: number;
	serviceType?: string;
	priceCents?: number;
	notes?: string;
	isRecurring: boolean;
	recurrencePattern?: {
		type: string;
		interval: number;
		until?: string;
		days_of_week?: number[];
	};
	createdAt: string;
	updatedAt: string;
}

export interface RoomAllocation {
	id: string;
	roomId: string;
	allocationType: 'dedicated' | 'shared';
	startDate?: string;
	endDate?: string;
	isActive: boolean;
	roomName?: string;
}

export interface AvailabilityDay {
	date: string;
	dayName: string;
	slots: Availability[];
}

export const load: PageServerLoad = async ({ locals, url, fetch }) => {
	if (!locals.user?.id) {
		throw redirect(302, '/auth');
	}

	const partnerId = locals.user.id;

	try {
		// Fetch availabilities and allocations in parallel
		const [availRes, allocRes] = await Promise.all([
			fetch(`${env.API_URL}/partners/${partnerId}/availabilities`, {
				headers: { 'Content-Type': 'application/json' }
			}),
			fetch(`${env.API_URL}/partners/${partnerId}/allocations`, {
				headers: { 'Content-Type': 'application/json' }
			})
		]);

		if (availRes.status === 401) {
			throw redirect(302, '/auth');
		}

		// Process availabilities
		let availabilities: AvailabilityDay[] = [];
		if (availRes.ok) {
			const data = await availRes.json();
			if (Array.isArray(data)) {
				availabilities = groupByDay(data);
			}
		} else {
			console.error(`Failed to fetch availabilities: ${availRes.status} ${availRes.statusText}`);
		}

		// Process allocations
		let allocations: RoomAllocation[] = [];
		if (allocRes.status === 401) {
			throw redirect(302, '/auth');
		}
		if (allocRes.ok) {
			const allocData = await allocRes.json();
			if (Array.isArray(allocData)) {
				allocations = allocData
					.filter((a: any) => a.is_active)
					.map((a: any): RoomAllocation => ({
						id: a.id,
						roomId: a.room_id,
						allocationType: a.allocation_type,
						startDate: a.start_date ?? undefined,
						endDate: a.end_date ?? undefined,
						isActive: a.is_active,
						roomName: a.room_name ?? undefined
					}));
			}
		} else {
			console.error(`Failed to fetch allocations: ${allocRes.status} ${allocRes.statusText}`);
		}

		return { availabilities, allocations, partnerId };
	} catch (err) {
		if (isRedirect(err)) throw err;
		console.error('Error loading availabilities:', err);
		return { availabilities: [], allocations: [], partnerId };
	}
};

function groupByDay(slots: any[]): AvailabilityDay[] {
	const dayNames = ['Dimanche', 'Lundi', 'Mardi', 'Mercredi', 'Jeudi', 'Vendredi', 'Samedi'];
	const today = new Date();
	today.setHours(0, 0, 0, 0);

	const slotsByDate = new Map<string, Availability[]>();

	for (const slot of slots) {
		const startDate = new Date(slot.start_time);
		const dateKey = startDate.toISOString().split('T')[0];

		const availability: Availability = {
			id: slot.id,
			date: dateKey,
			startTime: slot.start_time,
			endTime: slot.end_time,
			status: slot.status || 'available',
			roomId: slot.room_id,
			roomName: slot.room_id || 'Salle inconnue',
			maxCapacity: slot.max_capacity || 1,
			serviceType: slot.service_type,
			priceCents: slot.price_cents ?? undefined,
			notes: slot.notes,
			isRecurring: !!slot.recurrence_pattern,
			recurrencePattern: slot.recurrence_pattern ?? undefined,
			createdAt: slot.created_at,
			updatedAt: slot.updated_at
		};

		if (!slotsByDate.has(dateKey)) {
			slotsByDate.set(dateKey, []);
		}
		slotsByDate.get(dateKey)!.push(availability);
	}

	const availabilities: AvailabilityDay[] = [];
	for (let i = 0; i < 7; i++) {
		const d = new Date(today.getTime());
		d.setDate(d.getDate() + i);
		const dateKey = d.toISOString().split('T')[0];
		const dayName = dayNames[d.getDay()];

		const slots = slotsByDate.get(dateKey) || [];
		if (slots.length > 0 || d.getDay() !== 0) {
			availabilities.push({ date: dateKey, dayName, slots });
		}
	}

	return availabilities;
}
