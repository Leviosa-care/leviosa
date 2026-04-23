import type { PageServerLoad } from './$types';

interface RecentBooking {
    id: string;
    clientName: string;
    productName: string;
    therapistName: string;
    startTime: string;
    status: 'confirmed' | 'pending' | 'cancelled';
}

interface UpcomingBooking {
    id: string;
    clientName: string;
    productName: string;
    roomName: string;
    startTime: string;
    duration: number;
}

interface DashboardStats {
    revenueThisWeek: number;
    bookingsThisWeek: number;
    upcomingBookingsCount: number;
    pendingBookingsCount: number;
    activeProductsCount: number;
}

// TODO: Replace with real API calls when backend is ready
async function getMockDashboardData(): Promise<{
    stats: DashboardStats;
    recentBookings: RecentBooking[];
    upcomingBookings: UpcomingBooking[];
}> {
    // Mock stats
    const stats: DashboardStats = {
        revenueThisWeek: 12500, // cents (€125.00)
        bookingsThisWeek: 12,
        upcomingBookingsCount: 8,
        pendingBookingsCount: 3,
        activeProductsCount: 15
    };

    // Mock recent bookings
    const recentBookings: RecentBooking[] = [
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

    // Mock upcoming bookings
    const upcomingBookings: UpcomingBooking[] = [
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

export const load: PageServerLoad = async () => {
    // TODO: When backend is ready, switch to real API calls:
    // const stats = await fetch(`${API_URL}/admin/stats`).then(r => r.json());
    // const recentBookings = await fetch(`${API_URL}/admin/bookings?recent=5`).then(r => r.json());
    // const upcomingBookings = await fetch(`${API_URL}/admin/bookings?upcoming=3`).then(r => r.json());

    return await getMockDashboardData();
};
