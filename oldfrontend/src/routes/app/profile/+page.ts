import { Settings, Bell, Slack, User, BookOpenText } from 'lucide-svelte';

type Field = {
	name: string;
	icon: typeof import('lucide-svelte').Icon;
	pathname: string;
};

function getPath(end: string) {
	return `/app/profile/${end}`;
}

const juridique: Field[] = [
	{ name: 'Conditions de service', icon: BookOpenText, pathname: getPath('service-condition') },
	{
		name: 'Politique de confidentialite',
		icon: BookOpenText,
		pathname: getPath('confidential-policy')
	},
	{
		name: "Consentement droit a l'image",
		icon: BookOpenText,
		pathname: getPath('image-right-consent')
	}
];

const perso: Field[] = [
	{ name: 'Informations personnelles', icon: User, pathname: getPath('personal-information') },
	{ name: 'Parametres du compte', icon: Settings, pathname: getPath('parameters') }
];

const parameters: Field[] = [
	// TODO: change the next one for the Leviosa logo
	{ name: 'A propos', icon: Slack, pathname: getPath('about') },
	{ name: 'Notifications', icon: Bell, pathname: getPath('notifications') }
];
type PageRes = { perso: Field[]; juridique: Field[]; parameters: Field[] };
export function load(): PageRes {
	return {
		perso,
		juridique,
		parameters
	};
}
