import type { PageServerLoad } from './$types';
import { env } from '$env/dynamic/private';
import { error, redirect, isRedirect } from '@sveltejs/kit';

interface RecentBooking {
    id: string;
    client_name: string;
    product_name: string;
    partner_name: string;
    start_time: string;
    status: string;
}

interface UpcomingBooking {
    id: string;
    client_name: string;
    product_name: string;
    room_name: string;
    start_time: string;
    duration_min: number;
}

interface DashboardStatsAPI {
    bookings_this_week: number;
    revenue_this_week: number;
    upcoming_bookings_count: number;
    pending_bookings_count: number;
    active_products_count: number;
    recent_bookings: RecentBooking[];
    upcoming_bookings: UpcomingBooking[];
}

interface DashboardStats {
    revenueThisWeek: number;
    bookingsThisWeek: number;
    upcomingBookingsCount: number;
    pendingBookingsCount: number;
    activeProductsCount: number;
}

export interface RecentBookingUI {
    id: string;
    clientName: string;
    productName: string;
    therapistName: string;
    startTime: string;
    status: 'confirmed' | 'pending' | 'cancelled' | 'completed' | 'no_show';
}

export interface UpcomingBookingUI {
    id: string;
    clientName: string;
    productName: string;
    roomName: string;
    startTime: string;
    duration: number;
}

export const load: PageServerLoad = async ({ fetch }) => {
    if (env.USE_MOCK_DATA === 'true') {
        return getMockDashboardData();
    }

    try {
        const statsRes = await fetch(`${env.API_URL}/admin/dashboard/stats`);
        if (statsRes.status === 401) {
            throw redirect(302, '/auth');
        }
        if (!statsRes.ok) {
            throw new Error(`Failed to fetch dashboard stats: ${statsRes.status} ${statsRes.statusText}`);
        }
        const apiData: DashboardStatsAPI = await statsRes.json();
        const stats: DashboardStats = {
            revenueThisWeek: apiData.revenue_this_week,
            bookingsThisWeek: apiData.bookings_this_week,
            upcomingBookingsCount: apiData.upcoming_bookings_count,
            pendingBookingsCount: apiData.pending_bookings_count,
            activeProductsCount: apiData.active_products_count
        };

        const recentBookings: RecentBookingUI[] = (apiData.recent_bookings ?? []).map(
            (b: RecentBooking): RecentBookingUI => ({
                id: b.id,
                clientName: b.client_name,
                productName: b.product_name,
                therapistName: b.partner_name,
                startTime: b.start_time,
                status: b.status as RecentBookingUI['status']
            })
        );

        const upcomingBookings: UpcomingBookingUI[] = (apiData.upcoming_bookings ?? []).map(
            (b: UpcomingBooking): UpcomingBookingUI => ({
                id: b.id,
                clientName: b.client_name,
                productName: b.product_name,
                roomName: b.room_name,
                startTime: b.start_time,
                duration: b.duration_min
            })
        );

        return {
            stats,
            recentBookings,
            upcomingBookings
        };
    } catch (err) {
        if (isRedirect(err)) throw err;
        console.error('Error loading dashboard data:', err);
        throw error(503, 'Impossible de charger les données du tableau de bord. Veuillez réessayer.');
    }
};

async function getMockDashboardData(): Promise<{
    stats: DashboardStats;
    recentBookings: RecentBookingUI[];
    upcomingBookings: UpcomingBookingUI[];
}> {
    const stats: DashboardStats = {
        revenueThisWeek: 12500,
        bookingsThisWeek: 12,
        upcomingBookingsCount: 8,
        pendingBookingsCount: 3,
        activeProductsCount: 15
    };

    const recentBookings: RecentBookingUI[] = [
        {
            id: '1',
            clientName: 'Marie Dupont',
            productName: 'Massage Relaxant 60min',
            therapistName: 'Sophie Martin',
            startTime: new Date(Date.now() - 2 * 60 * 60 * 1000).toISOString(),
            status: 'confirmed'
        },
        {
            id: '2',
            clientName: 'Jean Durand',
            productName: 'Consultation Kiné 45min',
            therapistName: 'Pierre Leroy',
            startTime: new Date(Date.now() - 5 * 60 * 60 * 1000).toISOString(),
            status: 'confirmed'
        },
        {
            id: '3',
            clientName: 'Claire Bernard',
            productName: 'Massage Relaxant 90min',
            therapistName: 'Sophie Martin',
            startTime: new Date(Date.now() - 24 * 60 * 60 * 1000).toISOString(),
            status: 'confirmed'
        },
        {
            id: '4',
            clientName: 'Lucas Petit',
            productName: 'Soin du Dos 60min',
            therapistName: 'Marie Dubois',
            startTime: new Date(Date.now() - 36 * 60 * 60 * 1000).toISOString(),
            status: 'pending'
        },
        {
            id: '5',
            clientName: 'Emma Moreau',
            productName: 'Massage Relaxant 60min',
            therapistName: 'Sophie Martin',
            startTime: new Date(Date.now() - 48 * 60 * 60 * 1000).toISOString(),
            status: 'confirmed'
        }
    ];

    const upcomingBookings: UpcomingBookingUI[] = [
        {
            id: '6',
            clientName: 'Thomas Richard',
            productName: 'Consultation Kiné 45min',
            roomName: 'Cabinet 1',
            startTime: new Date(Date.now() + 1 * 60 * 60 * 1000).toISOString(),
            duration: 45
        },
        {
            id: '7',
            clientName: 'Camille Simon',
            productName: 'Massage Relaxant 60min',
            roomName: 'Cabinet 2',
            startTime: new Date(Date.now() + 3 * 60 * 60 * 1000).toISOString(),
            duration: 60
        },
        {
            id: '8',
            clientName: 'Hugo Michel',
            productName: 'Soin du Dos 60min',
            roomName: 'Cabinet 1',
            startTime: new Date(Date.now() + 5 * 60 * 60 * 1000).toISOString(),
            duration: 60
        }
    ];

    return { stats, recentBookings, upcomingBookings };
}

