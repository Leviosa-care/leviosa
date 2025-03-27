export type Conversation = {
	id: string;
	author: string;
	image: string;
	content: string;
	date: string;
	time: string;
};

// TODO: preciser ce type plus tard avec tout ce qu'il faut, c'est a Livio de preciser
export type SessionNote = {
	id: string;
	offer: string;
	author: string;
	date: string;
};

export type Message = {
	author: string;
	team?: string;
	content: string;
};
