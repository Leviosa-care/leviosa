import { describe, it, expect } from 'vitest';
import { convertMonthToInt } from './date';

describe('convertMonthToInt', () => {
	it('should convert "Janvier" to 1', () => {
		expect(convertMonthToInt('Janvier')).toBe(1);
	});

	it('should convert "Fevrier" to 2', () => {
		expect(convertMonthToInt('Fevrier')).toBe(2);
	});

	it('should convert "Mars" to 3', () => {
		expect(convertMonthToInt('Mars')).toBe(3);
	});

	it('should convert "Avril" to 4', () => {
		expect(convertMonthToInt('Avril')).toBe(4);
	});

	it('should convert "Mai" to 5', () => {
		expect(convertMonthToInt('Mai')).toBe(5);
	});

	it('should convert "Juin" to 6', () => {
		expect(convertMonthToInt('Juin')).toBe(6);
	});

	it('should convert "Juillet" to 7', () => {
		expect(convertMonthToInt('Juillet')).toBe(7);
	});

	it('should convert "Aout" to 8', () => {
		expect(convertMonthToInt('Aout')).toBe(8);
	});

	it('should convert "Septembre" to 9', () => {
		expect(convertMonthToInt('Septembre')).toBe(9);
	});

	it('should convert "Octobre" to 10', () => {
		expect(convertMonthToInt('Octobre')).toBe(10);
	});

	it('should convert "Novembre" to 11', () => {
		expect(convertMonthToInt('Novembre')).toBe(11);
	});

	it('should convert "Decembre" to 12', () => {
		expect(convertMonthToInt('Decembre')).toBe(12);
	});
});
