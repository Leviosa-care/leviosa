import type { FieldValue } from '$lib/types';

export const values: FieldValue[] = [
	{
		name: 'fullname',
		value: 'John DOE',
		props: { firstname: 'John', lastname: 'DOE' }
	},
	{
		name: 'phonenumber',
		value: '',
		props: { phoneNumber: '0612345678' }
	},
	{
		name: 'gender',
		value: 'Homme'
	},
	{
		name: 'mail',
		value: 'john.doe@gmail.com',
		props: { email: 'john.doe@gmail.com' }
	},
	{
		name: 'password',
		value: '************'
	},
	{
		name: 'city',
		value: 'Ivry-Sur-Seine, Paris',
		props: { city: 'Ivry-Sur-Seine, Paris' }
	},
	{
		name: 'postalCode',
		value: 94200,
		props: { postalCode: 94200 }
	},
	{
		name: 'address1',
		value: '',
		props: { address: '' }
	},
	{
		name: 'address2',
		value: '',
		props: { address: '' }
	}
];
