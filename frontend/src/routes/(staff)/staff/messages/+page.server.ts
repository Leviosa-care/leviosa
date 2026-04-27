import type { PageServerLoad } from './$types';

export interface Message {
	id: string;
	content: string;
	sentAt: string;
	fromPartner: boolean;
}

export interface Conversation {
	id: string;
	clientName: string;
	clientInitials: string;
	lastMessage: string;
	lastMessageAt: string;
	unreadCount: number;
	messages: Message[];
}

function getMockConversations(): Conversation[] {
	const now = new Date();
	const t = (offsetMin: number) => new Date(now.getTime() - offsetMin * 60 * 1000).toISOString();

	return [
		{
			id: 'c1',
			clientName: 'Marie Dupont',
			clientInitials: 'MD',
			lastMessage: 'Parfait, merci pour la confirmation !',
			lastMessageAt: t(10),
			unreadCount: 2,
			messages: [
				{ id: 'm1', content: 'Bonjour, je souhaiterais savoir si vous êtes disponible jeudi à 10h ?', sentAt: t(60), fromPartner: false },
				{ id: 'm2', content: 'Bonjour Marie ! Oui, je suis disponible jeudi à 10h. Je vous confirme le créneau.', sentAt: t(45), fromPartner: true },
				{ id: 'm3', content: 'Parfait, merci pour la confirmation !', sentAt: t(10), fromPartner: false },
			],
		},
		{
			id: 'c2',
			clientName: 'Jean Durand',
			clientInitials: 'JD',
			lastMessage: 'Je serai peut-être un peu en retard demain.',
			lastMessageAt: t(90),
			unreadCount: 1,
			messages: [
				{ id: 'm4', content: 'Je serai peut-être un peu en retard demain.', sentAt: t(90), fromPartner: false },
			],
		},
		{
			id: 'c3',
			clientName: 'Claire Bernard',
			clientInitials: 'CB',
			lastMessage: 'À samedi alors !',
			lastMessageAt: t(240),
			unreadCount: 0,
			messages: [
				{ id: 'm5', content: 'Pouvez-vous me recommander des exercices pour les tensions cervicales ?', sentAt: t(300), fromPartner: false },
				{ id: 'm6', content: 'Bien sûr ! Je vous enverrai une fiche lors de votre prochaine séance samedi.', sentAt: t(270), fromPartner: true },
				{ id: 'm7', content: 'À samedi alors !', sentAt: t(240), fromPartner: false },
			],
		},
		{
			id: 'c4',
			clientName: 'Lucas Petit',
			clientInitials: 'LP',
			lastMessage: 'D\'accord, je prendrai rendez-vous la semaine prochaine.',
			lastMessageAt: t(1440),
			unreadCount: 0,
			messages: [
				{ id: 'm8', content: 'Bonjour, je dois malheureusement annuler ma séance de demain.', sentAt: t(1500), fromPartner: false },
				{ id: 'm9', content: 'Pas de problème Lucas, j\'ai annulé votre réservation. N\'hésitez pas à reprendre un créneau quand vous souhaitez.', sentAt: t(1470), fromPartner: true },
				{ id: 'm10', content: 'D\'accord, je prendrai rendez-vous la semaine prochaine.', sentAt: t(1440), fromPartner: false },
			],
		},
		{
			id: 'c5',
			clientName: 'Emma Moreau',
			clientInitials: 'EM',
			lastMessage: 'Super, merci beaucoup !',
			lastMessageAt: t(2880),
			unreadCount: 0,
			messages: [
				{ id: 'm11', content: 'Super, merci beaucoup !', sentAt: t(2880), fromPartner: false },
			],
		},
	];
}

export const load: PageServerLoad = async () => {
	// TODO: Replace with messaging API endpoints
	return {
		conversations: getMockConversations(),
	};
};
