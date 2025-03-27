import type { Offer } from '$lib/types';

export const offers: Offer[] = [
	{
		type: 'Massage',
		services: [
			{
				name: 'Standard',
				label: "L'offre la plus appreciee de nos clients",
				description:
					'Venez decouvrir les massages standard de Levioisa afin de vous relaxer et blabla pour pouvoir remplir une zone de texte assez consequente et ainsi avoir une idee de comment tout cela va se mettre en place sur la page.',
				image: 'image standard massage',
				positive_responses: 92,
				clients_count: 172,
				duration: 30,
				freelancers: [
					{
						id: '1',
						firstname: 'Livio',
						lastname: 'HENRY',
						avatar: 'https://pbs.twimg.com/media/FA9tIzxUUAUy80c.jpg',
						reviews: [
							{
								note: 1,
								count: 1
							},
							{
								note: 3,
								count: 23
							},
							{
								note: 4,
								count: 32
							},
							{
								note: 5,
								count: 12
							}
						],
						ratings: 8.8,
						bio: 'Lorem ipsum dolor sit amet consectetur adipisicing elit. Ipsam quaerat blanditiis illum fugit. Reprehenderit ab quo quam ipsam cum dolore expedita tempore ratione. Repellat esse a, labore minima est quod sit repellendus?'
					},
					{
						id: '2',
						firstname: 'Kelsy',
						lastname: 'CHERDIEU',
						avatar: 'https://pbs.twimg.com/media/FA9tIzxUUAUy80c.jpg',
						reviews: [
							{
								note: 1,
								count: 1
							},
							{
								note: 3,
								count: 23
							},
							{
								note: 4,
								count: 32
							},
							{
								note: 5,
								count: 12
							}
						],
						ratings: 8.8,
						bio: 'Lorem ipsum dolor sit amet consectetur adipisicing elit. Ipsam quaerat blanditiis illum fugit. Reprehenderit ab quo quam ipsam cum dolore expedita tempore ratione. Repellat esse a, labore minima est quod sit repellendus?'
					}
				],
				bgurl:
					'https://media1.popsugar-assets.com/files/thumbor/spbNmLSC9dFfgLIN41jVpX-1w5E/fit-in/2048xorig/filters:format_auto-!!-:strip_icc-!!-/2019/12/20/718/n/1922729/tmp_Eg85Vj_65809b0747406623_GettyImages-925256722.jpg',
				bgblur: 200,
				isLight: true
			},
			{
				name: 'Premium',
				label: 'Une offre premium pour les gens styles',
				description: 'Some long ass description for the service',
				image: 'image premium massage',
				positive_responses: 92,
				clients_count: 89,
				duration: 60,
				freelancers: [
					{
						id: '3',
						firstname: 'Livio',
						lastname: 'HENRY',
						avatar: 'https://pbs.twimg.com/media/FA9tIzxUUAUy80c.jpg',
						reviews: [
							{
								note: 1,
								count: 1
							},
							{
								note: 3,
								count: 23
							},
							{
								note: 4,
								count: 32
							},
							{
								note: 5,
								count: 12
							}
						],
						ratings: 8.8,
						bio: 'Lorem ipsum dolor sit amet consectetur adipisicing elit. Ipsam quaerat blanditiis illum fugit. Reprehenderit ab quo quam ipsam cum dolore expedita tempore ratione. Repellat esse a, labore minima est quod sit repellendus?'
					},
					{
						id: '4',
						firstname: 'Kelsy',
						lastname: 'CHERDIEU',
						avatar: 'https://pbs.twimg.com/media/FA9tIzxUUAUy80c.jpg',
						reviews: [
							{
								note: 1,
								count: 1
							},
							{
								note: 3,
								count: 23
							},
							{
								note: 4,
								count: 32
							},
							{
								note: 5,
								count: 12
							}
						],
						ratings: 8.8,
						bio: 'Lorem ipsum dolor sit amet consectetur adipisicing elit. Ipsam quaerat blanditiis illum fugit. Reprehenderit ab quo quam ipsam cum dolore expedita tempore ratione. Repellat esse a, labore minima est quod sit repellendus?'
					}
				],
				bgurl:
					'https://images.healthshots.com/healthshots/en/uploads/2022/12/19171903/foot-massage.jpg',
				bgblur: 200,
				isLight: true
			},
			{
				name: 'Elite',
				label: 'Une offre premium pour les gens styles',
				description: 'Some long ass description for the service',
				image: 'image premium massage',
				positive_responses: 92,
				clients_count: 89,
				duration: 60,
				freelancers: [
					{
						id: '4',
						firstname: 'Livio',
						lastname: 'HENRY',
						avatar: 'https://pbs.twimg.com/media/FA9tIzxUUAUy80c.jpg',
						reviews: [
							{
								note: 1,
								count: 1
							},
							{
								note: 3,
								count: 23
							},
							{
								note: 4,
								count: 32
							},
							{
								note: 5,
								count: 12
							}
						],
						ratings: 8.8,
						bio: 'Lorem ipsum dolor sit amet consectetur adipisicing elit. Ipsam quaerat blanditiis illum fugit. Reprehenderit ab quo quam ipsam cum dolore expedita tempore ratione. Repellat esse a, labore minima est quod sit repellendus?'
					},
					{
						id: '5',
						firstname: 'Kelsy',
						lastname: 'CHERDIEU',
						avatar: 'https://pbs.twimg.com/media/FA9tIzxUUAUy80c.jpg',
						reviews: [
							{
								note: 1,
								count: 1
							},
							{
								note: 3,
								count: 23
							},
							{
								note: 4,
								count: 32
							},
							{
								note: 5,
								count: 12
							}
						],
						ratings: 8.8,
						bio: 'Lorem ipsum dolor sit amet consectetur adipisicing elit. Ipsam quaerat blanditiis illum fugit. Reprehenderit ab quo quam ipsam cum dolore expedita tempore ratione. Repellat esse a, labore minima est quod sit repellendus?'
					}
				],
				bgurl: 'https://media.timeout.com/images/103170265/image.jpg',
				bgblur: 200,
				isLight: true
			}
		]
	},
	{
		type: 'Coaching mental',
		services: [
			{
				name: 'Standard',
				label: 'Pour une premiere experience de coaching mental',
				description: 'Some long ass description for the service',
				image: 'image standard coaching mental',
				positive_responses: 75,
				clients_count: 34,
				duration: 30,
				freelancers: [
					{
						id: '6',
						firstname: 'Livio',
						lastname: 'HENRY',
						avatar: 'https://pbs.twimg.com/media/FA9tIzxUUAUy80c.jpg',
						reviews: [
							{
								note: 1,
								count: 1
							},
							{
								note: 3,
								count: 23
							},
							{
								note: 4,
								count: 32
							},
							{
								note: 5,
								count: 12
							}
						],
						ratings: 8.8,
						bio: 'Lorem ipsum dolor sit amet consectetur adipisicing elit. Ipsam quaerat blanditiis illum fugit. Reprehenderit ab quo quam ipsam cum dolore expedita tempore ratione. Repellat esse a, labore minima est quod sit repellendus?'
					},
					{
						id: '1',
						firstname: 'Kelsy',
						lastname: 'CHERDIEU',
						avatar: 'https://pbs.twimg.com/media/FA9tIzxUUAUy80c.jpg',
						reviews: [
							{
								note: 1,
								count: 1
							},
							{
								note: 3,
								count: 23
							},
							{
								note: 4,
								count: 32
							},
							{
								note: 5,
								count: 12
							}
						],
						ratings: 8.8,
						bio: 'Lorem ipsum dolor sit amet consectetur adipisicing elit. Ipsam quaerat blanditiis illum fugit. Reprehenderit ab quo quam ipsam cum dolore expedita tempore ratione. Repellat esse a, labore minima est quod sit repellendus?'
					}
				],
				bgurl:
					'https://media1.popsugar-assets.com/files/thumbor/spbNmLSC9dFfgLIN41jVpX-1w5E/fit-in/2048xorig/filters:format_auto-!!-:strip_icc-!!-/2019/12/20/718/n/1922729/tmp_Eg85Vj_65809b0747406623_GettyImages-925256722.jpg',
				bgblur: 200,
				isLight: false
			},
			{
				name: 'Premium',
				label: "Pour les plus experimentes de l'experience de coaching mental",
				description: 'Some long ass description for the service',
				image: 'image premium coaching mental',
				positive_responses: 87,
				clients_count: 23,
				duration: 45,
				freelancers: [
					{
						id: '1',
						firstname: 'Livio',
						lastname: 'HENRY',
						avatar: 'https://pbs.twimg.com/media/FA9tIzxUUAUy80c.jpg',
						reviews: [
							{
								note: 1,
								count: 1
							},
							{
								note: 3,
								count: 23
							},
							{
								note: 4,
								count: 32
							},
							{
								note: 5,
								count: 12
							}
						],
						ratings: 8.8,
						bio: 'Lorem ipsum dolor sit amet consectetur adipisicing elit. Ipsam quaerat blanditiis illum fugit. Reprehenderit ab quo quam ipsam cum dolore expedita tempore ratione. Repellat esse a, labore minima est quod sit repellendus?'
					},
					{
						id: '1',
						firstname: 'Kelsy',
						lastname: 'CHERDIEU',
						avatar: 'https://pbs.twimg.com/media/FA9tIzxUUAUy80c.jpg',
						reviews: [
							{
								note: 1,
								count: 1
							},
							{
								note: 3,
								count: 23
							},
							{
								note: 4,
								count: 32
							},
							{
								note: 5,
								count: 12
							}
						],
						ratings: 8.8,
						bio: 'Lorem ipsum dolor sit amet consectetur adipisicing elit. Ipsam quaerat blanditiis illum fugit. Reprehenderit ab quo quam ipsam cum dolore expedita tempore ratione. Repellat esse a, labore minima est quod sit repellendus?'
					}
				],
				bgurl:
					'https://media1.popsugar-assets.com/files/thumbor/spbNmLSC9dFfgLIN41jVpX-1w5E/fit-in/2048xorig/filters:format_auto-!!-:strip_icc-!!-/2019/12/20/718/n/1922729/tmp_Eg85Vj_65809b0747406623_GettyImages-925256722.jpg',
				bgblur: 200,
				isLight: false
			}
		]
	}
];
