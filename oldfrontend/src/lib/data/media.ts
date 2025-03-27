import type { EventPhotos, EventVideos } from '$lib/types';

export const eventsPhotos: EventPhotos[] = [
	{
		date: '12/12/2024',
		photos: [
			'https://thecinemaholic.com/wp-content/uploads/2018/02/project-x.jpg',
			'https://static.wixstatic.com/media/e5960e_6c4f92a7be8843a3b2ba0c4cbc16ef1b~mv2.jpg/v1/fill/w_892,h_678,al_c/e5960e_6c4f92a7be8843a3b2ba0c4cbc16ef1b~mv2.jpg',
			'https://imgsrc.cineserie.com/2021/11/wordpress-cs-format-image-39.png?ver=1',
			'https://i.kinja-img.com/gawker-media/image/upload/s--UNs3EDm6--/c_scale,f_auto,fl_progressive,q_80,w_800/1354909709948073574.jpg',
			'https://synapsetrading.com/wp-content/uploads/2016/12/rich-people.jpg'
		]
	},
	{
		date: '16/12/2024',
		photos: [
			'https://thelistlove.com/wp-content/uploads/2014/12/the-wolf-of-wall-street-dancing.jpg',
			'https://assets.bwbx.io/images/i5QPIpGBKntc/v1/-1x-1.jpg',
			'https://fandomwire.com/wp-content/uploads/2021/06/4-24.jpg'
		]
	}
];
export const eventsVideos: EventVideos[] = [
	{
		date: '12/12/2024',
		videos: [
			{
				thumbnail: 'https://thecinemaholic.com/wp-content/uploads/2018/02/project-x.jpg',
				duration: '00:30'
			},
			{
				thumbnail:
					'https://static.wixstatic.com/media/e5960e_6c4f92a7be8843a3b2ba0c4cbc16ef1b~mv2.jpg/v1/fill/w_892,h_678,al_c/e5960e_6c4f92a7be8843a3b2ba0c4cbc16ef1b~mv2.jpg',
				duration: '00:34'
			},
			{
				thumbnail: 'https://imgsrc.cineserie.com/2021/11/wordpress-cs-format-image-39.png?ver=1',
				duration: '01:14'
			},
			{
				thumbnail:
					'https://i.kinja-img.com/gawker-media/image/upload/s--UNs3EDm6--/c_scale,f_auto,fl_progressive,q_80,w_800/1354909709948073574.jpg',
				duration: '00:17'
			},
			{
				thumbnail: 'https://synapsetrading.com/wp-content/uploads/2016/12/rich-people.jpg',
				duration: '02:53'
			}
		]
	},
	{
		date: '16/12/2024',
		videos: [
			{
				thumbnail:
					'https://thelistlove.com/wp-content/uploads/2014/12/the-wolf-of-wall-street-dancing.jpg',
				duration: '01:29'
			},
			{
				thumbnail: 'https://assets.bwbx.io/images/i5QPIpGBKntc/v1/-1x-1.jpg',
				duration: '01:33'
			},
			{
				thumbnail: 'https://fandomwire.com/wp-content/uploads/2021/06/4-24.jpg',
				duration: '00:08'
			}
		]
	}
];
