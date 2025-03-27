import type { FieldValue, FieldConstructor, Field } from '$lib/types';
import { fieldsConstructors } from '$lib/constructor';

type PageRes = { fields: Field[] };

export function load({ data }): PageRes {
	// TODO: need to join these things
	const fields: Field[] = fieldsConstructors.map((constructor) => {
		const value = data.values.find((value: FieldValue) => value.name === constructor.name);
		return {
			...constructor,
			...value
		};
	});
	return { fields };
}
