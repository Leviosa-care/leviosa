import type { EventPickerMonth } from '$lib/types';

export const monthsData: EventPickerMonth[] = [
	{
		month: 'Octobre',
		days: [
			{ day: 'Lundi', date: 23, hours: [8, 12, 16, 17] },
			{ day: 'Jeudi', date: 13, hours: [4, 7, 10, 15] },
			{ day: 'Jeudi', date: 13, hours: [4, 7, 10, 15, 18, 19, 21] }
		]
	},
	{
		month: 'Novembre',
		days: [
			{ day: 'Samedi', date: 5, hours: [9, 10, 11] },
			{ day: 'Mercredi', date: 13, hours: [15, 16, 17, 18] }
		]
	},
	{
		month: 'Decembre',
		days: [
			{ day: 'Mardi', date: 6, hours: [11, 15, 16] },
			{ day: 'Jeudi', date: 8, hours: [8, 10, 12, 18] }
		]
	}
];
