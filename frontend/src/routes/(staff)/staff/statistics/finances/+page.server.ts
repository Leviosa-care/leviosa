import type { PageServerLoad } from './$types';

export interface EarningsSummary {
	currentMonthInCents: number;
	lastMonthInCents: number;
	pendingInCents: number;
	nextPayoutDate: string;
	nextPayoutInCents: number;
}

export interface Transaction {
	id: string;
	date: string;
	clientName: string;
	productName: string;
	amountInCents: number;
	status: 'paid' | 'pending' | 'refunded';
}

export interface MonthlyEarning {
	month: string;
	amountInCents: number;
}

function getMockFinances() {
	const summary: EarningsSummary = {
		currentMonthInCents: 285000,
		lastMonthInCents: 256000,
		pendingInCents: 37500,
		nextPayoutDate: new Date(Date.now() + 5 * 24 * 60 * 60 * 1000).toISOString(),
		nextPayoutInCents: 247500,
	};

	const now = new Date();
	const transactions: Transaction[] = [
		{ id: 't1', date: new Date(now.getTime() - 2 * 3600000).toISOString(), clientName: 'Marie Dupont', productName: 'Massage Relaxant 60min', amountInCents: 7500, status: 'pending' },
		{ id: 't2', date: new Date(now.getTime() - 26 * 3600000).toISOString(), clientName: 'Thomas Richard', productName: 'Drainage Lymphatique 90min', amountInCents: 11000, status: 'paid' },
		{ id: 't3', date: new Date(now.getTime() - 50 * 3600000).toISOString(), clientName: 'Camille Simon', productName: 'Massage Relaxant 60min', amountInCents: 7500, status: 'paid' },
		{ id: 't4', date: new Date(now.getTime() - 74 * 3600000).toISOString(), clientName: 'Hugo Michel', productName: 'Soin du Dos 60min', amountInCents: 9000, status: 'paid' },
		{ id: 't5', date: new Date(now.getTime() - 98 * 3600000).toISOString(), clientName: 'Léa Fontaine', productName: 'Massage Relaxant 60min', amountInCents: 7500, status: 'refunded' },
		{ id: 't6', date: new Date(now.getTime() - 122 * 3600000).toISOString(), clientName: 'Antoine Garnier', productName: 'Massage Sportif 60min', amountInCents: 8500, status: 'paid' },
		{ id: 't7', date: new Date(now.getTime() - 146 * 3600000).toISOString(), clientName: 'Sophie Blanc', productName: 'Soin du Dos 60min', amountInCents: 9000, status: 'paid' },
		{ id: 't8', date: new Date(now.getTime() - 170 * 3600000).toISOString(), clientName: 'Paul Mercier', productName: 'Drainage Lymphatique 90min', amountInCents: 11000, status: 'paid' },
	];

	const monthlyEarnings: MonthlyEarning[] = [
		{ month: 'Nov', amountInCents: 198000 },
		{ month: 'Déc', amountInCents: 224000 },
		{ month: 'Jan', amountInCents: 215000 },
		{ month: 'Fév', amountInCents: 241000 },
		{ month: 'Mar', amountInCents: 256000 },
		{ month: 'Avr', amountInCents: 285000 },
	];

	return { summary, transactions, monthlyEarnings };
}

export const load: PageServerLoad = async () => {
	// TODO: Replace with earnings endpoints when available
	return getMockFinances();
};
