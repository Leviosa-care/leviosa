import { env } from '$env/dynamic/private';
import type { PageServerLoad } from './$types';

interface Message {
	sender: 'admin' | 'client';
	content: string;
	sentAt: string;
}

interface Conversation {
	id: string;
	clientName: string;
	lastMessage: string;
	lastMessageAt: string;
	unreadCount: number;
	isOnline: boolean;
	messages: Message[];
}

async function getMockConversations(): Promise<Conversation[]> {
	const now = new Date();

	return [
		{
			id: '1',
			clientName: 'Marie Dupont',
			lastMessage: 'Merci pour le rendez-vous, à bientôt!',
			lastMessageAt: new Date(now.getTime() - 5 * 60 * 1000).toISOString(),
			unreadCount: 0,
			isOnline: true,
			messages: [
				{ sender: 'client', content: 'Bonjour, je voudrais réserver un massage pour vendredi prochain', sentAt: new Date(now.getTime() - 2 * 60 * 60 * 1000).toISOString() },
				{ sender: 'admin', content: 'Bonjour Marie! Bien sûr, quel créneau préférez-vous?', sentAt: new Date(now.getTime() - 105 * 60 * 1000).toISOString() },
				{ sender: 'client', content: 'Plutôt en matinée, vers 10h si possible', sentAt: new Date(now.getTime() - 100 * 60 * 1000).toISOString() },
				{ sender: 'admin', content: 'C\'est noté! Je vous ai réservé pour vendredi à 10h. À bientôt!', sentAt: new Date(now.getTime() - 95 * 60 * 1000).toISOString() },
				{ sender: 'client', content: 'Merci pour le rendez-vous, à bientôt!', sentAt: new Date(now.getTime() - 5 * 60 * 1000).toISOString() }
			]
		},
		{
			id: '2',
			clientName: 'Jean Durand',
			lastMessage: 'Quelle est la durée de la séance?',
			lastMessageAt: new Date(now.getTime() - 30 * 60 * 1000).toISOString(),
			unreadCount: 2,
			isOnline: false,
			messages: [
				{ sender: 'client', content: 'Bonjour, je suis intéressé par le drainage lymphatique', sentAt: new Date(now.getTime() - 45 * 60 * 1000).toISOString() },
				{ sender: 'admin', content: 'Bonjour Jean! Le drainage lymphatique dure 60 minutes. Cela vous convient-il?', sentAt: new Date(now.getTime() - 40 * 60 * 1000).toISOString() },
				{ sender: 'client', content: 'Quelle est la durée de la séance?', sentAt: new Date(now.getTime() - 30 * 60 * 1000).toISOString() }
			]
		},
		{
			id: '3',
			clientName: 'Claire Bernard',
			lastMessage: 'D\'accord, merci!',
			lastMessageAt: new Date(now.getTime() - 2 * 60 * 60 * 1000).toISOString(),
			unreadCount: 0,
			isOnline: true,
			messages: [
				{ sender: 'client', content: 'Je dois annuler mon rendez-vous de demain', sentAt: new Date(now.getTime() - 3 * 60 * 60 * 1000).toISOString() },
				{ sender: 'admin', content: 'Pas de problème Claire. Voulez-vous le reporter à une autre date?', sentAt: new Date(now.getTime() - 2 * 60 * 60 * 1000 + 30 * 60 * 1000).toISOString() },
				{ sender: 'client', content: 'Oui, la semaine prochaine si possible', sentAt: new Date(now.getTime() - 2 * 60 * 60 * 1000 + 20 * 60 * 1000).toISOString() },
				{ sender: 'admin', content: 'Je vous ai proposé deux créneaux la semaine prochaine. N\'hésitez pas à choisir celui qui vous convient.', sentAt: new Date(now.getTime() - 2 * 60 * 60 * 1000 + 10 * 60 * 1000).toISOString() },
				{ sender: 'client', content: 'D\'accord, merci!', sentAt: new Date(now.getTime() - 2 * 60 * 60 * 1000).toISOString() }
			]
		},
		{
			id: '4',
			clientName: 'Lucas Petit',
			lastMessage: 'Les prix des massages?',
			lastMessageAt: new Date(now.getTime() - 24 * 60 * 60 * 1000).toISOString(),
			unreadCount: 1,
			isOnline: false,
			messages: [
				{ sender: 'client', content: 'Bonjour', sentAt: new Date(now.getTime() - 26 * 60 * 60 * 1000).toISOString() },
				{ sender: 'admin', content: 'Bonjour Lucas! Comment puis-je vous aider?', sentAt: new Date(now.getTime() - 25 * 60 * 60 * 1000).toISOString() },
				{ sender: 'client', content: 'Les prix des massages?', sentAt: new Date(now.getTime() - 24 * 60 * 60 * 1000).toISOString() }
			]
		},
		{
			id: '5',
			clientName: 'Emma Moreau',
			lastMessage: 'Parfait, je confirme',
			lastMessageAt: new Date(now.getTime() - 3 * 24 * 60 * 60 * 1000).toISOString(),
			unreadCount: 0,
			isOnline: false,
			messages: [
				{ sender: 'client', content: 'Je voudrais prendre rendez-vous pour un soin du dos', sentAt: new Date(now.getTime() - 4 * 24 * 60 * 60 * 1000).toISOString() },
				{ sender: 'admin', content: 'Bonjour Emma! Le soin du dos dure 60 minutes et coûte 65€. Je vous propose mercredi prochain à 14h.', sentAt: new Date(now.getTime() - 3 * 24 * 60 * 60 * 1000 + 30 * 60 * 1000).toISOString() },
				{ sender: 'client', content: 'Parfait, je confirme', sentAt: new Date(now.getTime() - 3 * 24 * 60 * 60 * 1000).toISOString() }
			]
		}
	];
}

export const load: PageServerLoad = async () => {
	if (env.USE_MOCK_DATA === 'true') {
		return { conversations: await getMockConversations() };
	}
	return { conversations: [] };
};
