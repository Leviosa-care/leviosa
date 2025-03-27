import type { Month } from '$lib/types';
// TODO: make a  type for number that goes from 0 to 12
export function convertMonthToInt(month: Month): number {
	switch (month) {
		case 'Janvier':
			return 1;
		case 'Fevrier':
			return 2;
		case 'Mars':
			return 3;
		case 'Avril':
			return 4;
		case 'Mai':
			return 5;
		case 'Juin':
			return 6;
		case 'Juillet':
			return 7;
		case 'Aout':
			return 8;
		case 'Septembre':
			return 9;
		case 'Octobre':
			return 10;
		case 'Novembre':
			return 11;
		case 'Decembre':
			return 12;
		default:
			return 0;
	}
}
