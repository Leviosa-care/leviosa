import { env } from '$env/dynamic/private';
import type { PageServerLoad } from './$types';

interface SummaryKPIs {
	grossRevenueInCents: number;
	refundsInCents: number;
	netRevenueInCents: number;
}

interface Transaction {
	id: string;
	date: string;
	description: string;
	clientName: string;
	amountInCents: number;
	type: 'payment' | 'refund';
	status: 'completed' | 'pending' | 'failed';
	paymentMethod: 'card' | 'cash' | 'transfer';
}

interface MonthlyBreakdown {
	month: string;
	payments: number;
	refunds: number;
	net: number;
}

interface ComptaData {
	summary: SummaryKPIs;
	transactions: Transaction[];
	monthlyBreakdown: MonthlyBreakdown[];
}

async function getMockComptaData(): Promise<ComptaData> {
	const now = new Date();

	const summary: SummaryKPIs = {
		grossRevenueInCents: 52300,
		refundsInCents: 2150,
		netRevenueInCents: 50150
	};

	const transactions: Transaction[] = [
		{ id: 'TXN-001', date: new Date(now.getTime() - 2 * 60 * 60 * 1000).toISOString(), description: 'Massage Relaxant 60min', clientName: 'Marie Dupont', amountInCents: 5000, type: 'payment', status: 'completed', paymentMethod: 'card' },
		{ id: 'TXN-002', date: new Date(now.getTime() - 4 * 60 * 60 * 1000).toISOString(), description: 'Soin du Dos 60min', clientName: 'Jean Durand', amountInCents: 6500, type: 'payment', status: 'completed', paymentMethod: 'card' },
		{ id: 'TXN-003', date: new Date(now.getTime() - 6 * 60 * 60 * 1000).toISOString(), description: 'Remboursement', clientName: 'Claire Bernard', amountInCents: -5000, type: 'refund', status: 'completed', paymentMethod: 'card' },
		{ id: 'TXN-004', date: new Date(now.getTime() - 24 * 60 * 60 * 1000).toISOString(), description: 'Drainage Lymphatique', clientName: 'Lucas Petit', amountInCents: 7000, type: 'payment', status: 'completed', paymentMethod: 'cash' },
		{ id: 'TXN-005', date: new Date(now.getTime() - 26 * 60 * 60 * 1000).toISOString(), description: 'Massage Relaxant 90min', clientName: 'Emma Moreau', amountInCents: 7000, type: 'payment', status: 'pending', paymentMethod: 'card' },
		{ id: 'TXN-006', date: new Date(now.getTime() - 48 * 60 * 60 * 1000).toISOString(), description: 'Consultation Kiné 45min', clientName: 'Thomas Richard', amountInCents: 3500, type: 'payment', status: 'completed', paymentMethod: 'transfer' },
		{ id: 'TXN-007', date: new Date(now.getTime() - 50 * 60 * 60 * 1000).toISOString(), description: 'Massage Relaxant 60min', clientName: 'Camille Simon', amountInCents: 5000, type: 'payment', status: 'completed', paymentMethod: 'card' },
		{ id: 'TXN-008', date: new Date(now.getTime() - 72 * 60 * 60 * 1000).toISOString(), description: 'Soin du Dos 60min', clientName: 'Hugo Michel', amountInCents: 6500, type: 'payment', status: 'failed', paymentMethod: 'card' },
		{ id: 'TXN-009', date: new Date(now.getTime() - 74 * 60 * 60 * 1000).toISOString(), description: 'Drainage Lymphatique', clientName: 'Chloe Garcia', amountInCents: 7000, type: 'payment', status: 'completed', paymentMethod: 'card' },
		{ id: 'TXN-010', date: new Date(now.getTime() - 96 * 60 * 60 * 1000).toISOString(), description: 'Remboursement', clientName: 'Louis Laurent', amountInCents: -6500, type: 'refund', status: 'completed', paymentMethod: 'card' },
		{ id: 'TXN-011', date: new Date(now.getTime() - 98 * 60 * 60 * 1000).toISOString(), description: 'Massage Relaxant 60min', clientName: 'Jade Roux', amountInCents: 5000, type: 'payment', status: 'completed', paymentMethod: 'cash' },
		{ id: 'TXN-012', date: new Date(now.getTime() - 120 * 60 * 60 * 1000).toISOString(), description: 'Soin du Dos 60min', clientName: 'Nathan Girard', amountInCents: 6500, type: 'payment', status: 'completed', paymentMethod: 'card' },
		{ id: 'TXN-013', date: new Date(now.getTime() - 122 * 60 * 60 * 1000).toISOString(), description: 'Massage Relaxant 90min', clientName: 'Lola Bernard', amountInCents: 7000, type: 'payment', status: 'completed', paymentMethod: 'card' },
		{ id: 'TXN-014', date: new Date(now.getTime() - 144 * 60 * 60 * 1000).toISOString(), description: 'Consultation Kiné 45min', clientName: 'Enzo Dubois', amountInCents: 3500, type: 'payment', status: 'completed', paymentMethod: 'transfer' },
		{ id: 'TXN-015', date: new Date(now.getTime() - 168 * 60 * 60 * 1000).toISOString(), description: 'Drainage Lymphatique', clientName: 'Sarah Richard', amountInCents: 7000, type: 'payment', status: 'completed', paymentMethod: 'card' }
	];

	const monthlyBreakdown: MonthlyBreakdown[] = [
		{ month: 'Nov', payments: 38500, refunds: 1500, net: 37000 },
		{ month: 'Déc', payments: 45200, refunds: 0, net: 45200 },
		{ month: 'Jan', payments: 41800, refunds: 3500, net: 38300 },
		{ month: 'Fév', payments: 49600, refunds: 1800, net: 47800 },
		{ month: 'Mar', payments: 52100, refunds: 2150, net: 49950 },
		{ month: 'Avr', payments: 52300, refunds: 2150, net: 50150 }
	];

	return { summary, transactions, monthlyBreakdown };
}

export const load: PageServerLoad = async () => {
	if (env.USE_MOCK_DATA === 'true') {
		return await getMockComptaData();
	}
	return {
		summary: { grossRevenueInCents: 0, refundsInCents: 0, netRevenueInCents: 0 },
		transactions: [],
		monthlyBreakdown: [],
	};
};
