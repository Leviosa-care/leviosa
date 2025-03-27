import type { FieldConstructor } from '$lib/types';

import ModifiedFullName from '../../routes/app/profile/personal-information/ModifiedFullname.svelte';
import ModifiedCity from '../../routes/app/profile/personal-information/ModifiedCity.svelte';
import ModifiedGender from '../../routes/app/profile/personal-information/ModifiedGender.svelte';
import ModifiedPassword from '../../routes/app/profile/personal-information//ModifiedPassword.svelte';
import ModifiedAddress1 from '../../routes/app/profile/personal-information//ModifiedAddress1.svelte';
import ModifiedAddress2 from '../../routes/app/profile/personal-information//ModifiedAddress2.svelte';
import ModifiedMail from '../../routes/app/profile/personal-information//ModifiedMail.svelte';
import ModifiedPhoneNumber from '../../routes/app/profile/personal-information//ModifiedPhoneNumber.svelte';
import ModifiedPostalCode from '../../routes/app/profile/personal-information//ModifiedPostalCode.svelte';

export const fieldsConstructors: FieldConstructor[] = [
	{
		name: 'fullname',
		fieldname: 'Nom officiel',
		missingLabel: 'Information non fournie',
		modifyLabel:
			'Une modification de votre nom officiel entraine une suspension temporaire de votre compte en attendant que notre equipe revalide celui-ci',
		modifiedSlot: ModifiedFullName
	},
	{
		name: 'phonenumber',
		fieldname: 'Numero de telephone',
		missingLabel:
			'Ajoutez un numero de telephone afin que nous puissions vous contacter le plus rapidement possible',
		addLabel: 'Renseignez un numero de telephone sur lequel on peut vous joindre',
		modifyLabel: 'Utilisez un numero de telephone sur lequel on peut vous joindre rapidement',
		modifiedSlot: ModifiedPhoneNumber
	},
	{
		name: 'gender',
		fieldname: 'Genre',
		missingLabel: 'Ajoutez un genre...',
		addLabel: 'add label for gender',
		modifyLabel: 'modify label for gender',
		modifiedSlot: ModifiedGender
	},
	{
		name: 'mail',
		fieldname: 'Adresse e-mail',
		modifyLabel:
			'Une modification de votre adresse e-mail entraine une suspension temporaire de votre compte en attendant que notre equipe revalide celui-ci',
		modifiedSlot: ModifiedMail
	},
	{
		name: 'password',
		fieldname: 'Mot de passe',
		modifyLabel:
			'Une modification de votre mot de passe entraine une suspension temporaire de votre compte en attendant que notre equipe revalide celui-ci',
		modifiedSlot: ModifiedPassword
	},
	{
		name: 'city',
		fieldname: 'Ville',
		missingLabel: 'Information non fournie',
		addLabel:
			'Ajoutez votre ville afin que nous puissions vous proposer de meilleures recommendations',
		modifyLabel:
			"Utilisez le code postal d'une adresse permanente ou vous pouvez recevoir du courrier",
		modifiedSlot: ModifiedCity
	},
	{
		name: 'postalCode',
		fieldname: 'Code postal',
		missingLabel: 'Information non fournie',
		addLabel:
			'Ajoutez votre code postal afin que nous puissions vous proposer de meilleures recommendations',
		modifyLabel:
			"Utilisez le code postal d'une adresse permanente ou vous pouvez recevoir du courrier",
		modifiedSlot: ModifiedPostalCode
	},
	{
		name: 'address1',
		fieldname: 'Adresse 1',
		missingLabel: 'Information non fournie',
		addLabel:
			'Ajoutez votre adresse afin que nous puissions vous proposer de meilleures recommendations',
		modifyLabel: 'Utilisez une adresse permanente ou vous pouvez recevoir du courrier',
		modifiedSlot: ModifiedAddress1
	},
	{
		name: 'address2',
		fieldname: 'Adresse 2',
		missingLabel: 'Information non fournie',
		addLabel:
			'Ajoutez votre adresse afin que nous puissions vous proposer de meilleures recommendations',
		modifyLabel: 'Utilisez une adresse permanente ou vous pouvez recevoir du courrier',
		modifiedSlot: ModifiedAddress2
	}
];
