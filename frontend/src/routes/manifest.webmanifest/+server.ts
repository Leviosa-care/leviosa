import { json, type RequestHandler } from '@sveltejs/kit';
import { isAdminDomain, isStaffDomain } from '$lib/server/hostname';

export const GET: RequestHandler = async ({ url }) => {
	const hostname = url.hostname;
	const isAdmin = isAdminDomain(hostname);
	const isStaff = isStaffDomain(hostname);

	if (isAdmin) {
		return json({
			name: 'Leviosa Admin',
			short_name: 'Admin',
			description: 'Leviosa admin dashboard for care service management',
			start_url: '/',
			display: 'standalone',
			background_color: '#1a1a1a',
			theme_color: '#1a1a1a',
			lang: 'en',
			scope: '/',
			icons: [
				{
					src: '/pwa-admin-192x192.png',
					sizes: '192x192',
					type: 'image/png'
				},
				{
					src: '/pwa-admin-512x512.png',
					sizes: '512x512',
					type: 'image/png'
				},
				{
					src: '/pwa-admin-512x512.png',
					sizes: '512x512',
					type: 'image/png',
					purpose: 'any maskable'
				}
			]
		});
	}

	if (isStaff) {
		return json({
			name: 'Leviosa Staff',
			short_name: 'Staff',
			description: 'Leviosa staff portal for care service management',
			start_url: '/',
			display: 'standalone',
			background_color: '#475569',
			theme_color: '#475569',
			lang: 'en',
			scope: '/',
			icons: [
				{
					src: '/pwa-staff-192x192.png',
					sizes: '192x192',
					type: 'image/png'
				},
				{
					src: '/pwa-staff-512x512.png',
					sizes: '512x512',
					type: 'image/png'
				},
				{
					src: '/pwa-staff-512x512.png',
					sizes: '512x512',
					type: 'image/png',
					purpose: 'any maskable'
				}
			]
		});
	}

	return json({
		name: 'Leviosa',
		short_name: 'Leviosa',
		description: 'Care service platform',
		start_url: '/',
		display: 'standalone',
		background_color: '#ffffff',
		theme_color: '#ffffff',
		lang: 'en',
		scope: '/',
		icons: [
			{
				src: '/pwa-192x192.png',
				sizes: '192x192',
				type: 'image/png'
			},
			{
				src: '/pwa-512x512.png',
				sizes: '512x512',
				type: 'image/png'
			},
			{
				src: '/pwa-512x512.png',
				sizes: '512x512',
				type: 'image/png',
				purpose: 'any maskable'
			}
		]
	});
};
