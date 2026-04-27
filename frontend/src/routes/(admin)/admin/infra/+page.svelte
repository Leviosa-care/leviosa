<script lang="ts">
	import type { PageProps } from './$types';
	import { Server, Activity, Clock, Zap, CheckCircle2, XCircle, AlertTriangle, Rocket } from '@lucide/svelte';

	let { data }: PageProps = $props();

	function formatTime(isoString: string): string {
		const date = new Date(isoString);
		const now = new Date();
		const diffMs = now.getTime() - date.getTime();
		const diffSecs = Math.floor(diffMs / 1000);
		const diffMins = Math.floor(diffMs / (1000 * 60));
		const diffHours = Math.floor(diffMs / (1000 * 60 * 60));

		if (diffSecs < 60) {
			return `Il y a ${diffSecs}s`;
		} else if (diffMins < 60) {
			return `Il y a ${diffMins}min`;
		} else {
			return `Il y a ${diffHours}h`;
		}
	}

	function getDeploymentTime(isoString: string): string {
		const date = new Date(isoString);
		return date.toLocaleDateString('fr-FR', {
			day: '2-digit',
			month: '2-digit',
			year: '2-digit',
			hour: '2-digit',
			minute: '2-digit'
		});
	}

	function getStatusIcon(status: string) {
		switch (status) {
			case 'healthy':
				return CheckCircle2;
			case 'degraded':
				return AlertTriangle;
			case 'down':
				return XCircle;
			default:
				return Activity;
		}
	}

	function getStatusColor(status: string) {
		switch (status) {
			case 'healthy':
				return 'text-green-500';
			case 'degraded':
				return 'text-yellow-500';
			case 'down':
				return 'text-red-500';
			default:
				return 'text-gray-500';
		}
	}

	function getDeploymentStatusBadge(status: string) {
		switch (status) {
			case 'live':
				return 'bg-green-100 text-green-700 dark:bg-green-900 dark:text-green-300';
			case 'pending':
				return 'bg-yellow-100 text-yellow-700 dark:bg-yellow-900 dark:text-yellow-300';
			case 'failed':
				return 'bg-red-100 text-red-700 dark:bg-red-900 dark:text-red-300';
			default:
				return 'bg-gray-100 text-gray-700 dark:bg-gray-800 dark:text-gray-300';
		}
	}

	function getDeploymentStatusLabel(status: string): string {
		switch (status) {
			case 'live': return 'En ligne';
			case 'pending': return 'En cours';
			case 'failed': return 'Échoué';
			default: return status;
		}
	}
</script>

<svelte:head>
	<title>Infrastructure | Admin</title>
</svelte:head>

<div class="container mx-auto px-4 py-8 lg:py-12">
	<div class="mb-8">
		<h1 class="text-3xl lg:text-4xl font-bold mb-1 text-foreground">
			Infrastructure
		</h1>
		<p class="text-muted-foreground">
			État de santé des services et déploiements récents
		</p>
	</div>

	<!-- Services Grid -->
	<div class="mb-12">
		<h2 class="text-xl font-semibold text-foreground mb-6 flex items-center gap-2">
			<Server size={20} />
			État des services
		</h2>
		<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
			{#each data.services as service}
				{@const StatusIcon = getStatusIcon(service.status)}
				<div class="bg-card rounded-lg border border-border p-6">
					<div class="flex items-start justify-between mb-4">
						<div>
							<h3 class="font-semibold text-foreground mb-1">{service.name}</h3>
							<p class="text-xs text-muted-foreground">
								Vérifié {formatTime(service.lastChecked)}
							</p>
						</div>
						<StatusIcon size={24} class={getStatusColor(service.status)} />
					</div>
					<div class="space-y-3">
						<div class="flex items-center justify-between">
							<div class="flex items-center gap-2 text-sm text-muted-foreground">
								<Zap size={14} />
								<span>Latence</span>
							</div>
							<span class="text-sm font-medium text-foreground">{service.latencyMs}ms</span>
						</div>
						<div class="flex items-center justify-between">
							<div class="flex items-center gap-2 text-sm text-muted-foreground">
								<Activity size={14} />
								<span>Disponibilité</span>
							</div>
							<span class="text-sm font-medium text-foreground">{service.uptimePct}%</span>
						</div>
					</div>
					<div class="mt-4 pt-4 border-t border-border">
						<span class="inline-flex items-center gap-1 px-3 py-1 rounded-full text-xs font-medium {service.status === 'healthy'
							? 'bg-green-100 text-green-700 dark:bg-green-900 dark:text-green-300'
							: service.status === 'degraded'
								? 'bg-yellow-100 text-yellow-700 dark:bg-yellow-900 dark:text-yellow-300'
								: 'bg-red-100 text-red-700 dark:bg-red-900 dark:text-red-300'}">
							{service.status === 'healthy' ? 'Sain' : service.status === 'degraded' ? 'Dégradé' : 'Hors service'}
						</span>
					</div>
				</div>
			{/each}
		</div>
	</div>

	<!-- Deployments -->
	<div>
		<h2 class="text-xl font-semibold text-foreground mb-6 flex items-center gap-2">
			<Rocket size={20} />
			Déploiements récents
		</h2>
		<div class="bg-card rounded-lg border border-border overflow-hidden">
			<div class="overflow-x-auto">
				<table class="w-full">
					<thead>
						<tr class="bg-muted/50 border-b border-border">
							<th class="text-left py-4 px-6 text-sm font-medium text-muted-foreground">Environnement</th>
							<th class="text-left py-4 px-6 text-sm font-medium text-muted-foreground">Version</th>
							<th class="text-left py-4 px-6 text-sm font-medium text-muted-foreground">Déployé à</th>
							<th class="text-left py-4 px-6 text-sm font-medium text-muted-foreground">Statut</th>
						</tr>
					</thead>
					<tbody>
						{#each data.deployments as deployment}
							<tr class="border-b border-border hover:bg-muted/30 transition-colors">
								<td class="py-4 px-6">
									<div class="flex items-center gap-2">
										<div class="w-2 h-2 rounded-full {deployment.environment === 'Production'
											? 'bg-red-500'
											: 'bg-blue-500'}"></div>
										<span class="font-medium text-foreground">{deployment.environment}</span>
									</div>
								</td>
								<td class="py-4 px-6">
									<span class="text-sm font-mono text-foreground">{deployment.version}</span>
								</td>
								<td class="py-4 px-6">
									<div class="flex items-center gap-2 text-sm text-muted-foreground">
										<Clock size={14} />
										<span>{getDeploymentTime(deployment.deployedAt)}</span>
									</div>
								</td>
								<td class="py-4 px-6">
									<span class="px-3 py-1 rounded-full text-xs font-medium {getDeploymentStatusBadge(deployment.status)}">
										{getDeploymentStatusLabel(deployment.status)}
									</span>
								</td>
							</tr>
						{/each}
					</tbody>
				</table>
			</div>
		</div>
	</div>
</div>
