import { env } from '$env/dynamic/private';
import type { PageServerLoad } from './$types';

interface Service {
	name: string;
	status: 'healthy' | 'degraded' | 'down';
	latencyMs: number;
	uptimePct: number;
	lastChecked: string;
}

interface Deployment {
	environment: string;
	version: string;
	deployedAt: string;
	status: 'live' | 'pending' | 'failed';
}

interface InfraData {
	services: Service[];
	deployments: Deployment[];
}

async function getMockInfraData(): Promise<InfraData> {
	const now = new Date();

	const services: Service[] = [
		{ name: 'API', status: 'healthy', latencyMs: 12, uptimePct: 99.95, lastChecked: new Date(now.getTime() - 30 * 1000).toISOString() },
		{ name: 'PostgreSQL', status: 'healthy', latencyMs: 2, uptimePct: 99.99, lastChecked: new Date(now.getTime() - 45 * 1000).toISOString() },
		{ name: 'Redis', status: 'healthy', latencyMs: 1, uptimePct: 99.98, lastChecked: new Date(now.getTime() - 20 * 1000).toISOString() },
		{ name: 'RabbitMQ', status: 'degraded', latencyMs: 85, uptimePct: 98.5, lastChecked: new Date(now.getTime() - 60 * 1000).toISOString() },
		{ name: 'S3', status: 'healthy', latencyMs: 45, uptimePct: 99.9, lastChecked: new Date(now.getTime() - 15 * 1000).toISOString() },
		{ name: 'Vault', status: 'healthy', latencyMs: 8, uptimePct: 99.97, lastChecked: new Date(now.getTime() - 35 * 1000).toISOString() }
	];

	const deployments: Deployment[] = [
		{ environment: 'Production', version: 'v1.2.3', deployedAt: new Date(now.getTime() - 2 * 60 * 60 * 1000).toISOString(), status: 'live' },
		{ environment: 'Staging', version: 'v1.2.4-rc1', deployedAt: new Date(now.getTime() - 24 * 60 * 60 * 1000).toISOString(), status: 'live' },
		{ environment: 'Production', version: 'v1.2.2', deployedAt: new Date(now.getTime() - 5 * 24 * 60 * 60 * 1000).toISOString(), status: 'live' },
		{ environment: 'Staging', version: 'v1.2.5-rc1', deployedAt: new Date(now.getTime() - 6 * 60 * 60 * 1000).toISOString(), status: 'pending' }
	];

	return { services, deployments };
}

export const load: PageServerLoad = async () => {
	if (env.USE_MOCK_DATA === 'true') {
		return await getMockInfraData();
	}
	return { services: [], deployments: [] };
};
