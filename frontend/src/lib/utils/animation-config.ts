export const ANIMATION_PRESETS = {
	'fade-in': {
		initial: { opacity: 0 },
		animate: { opacity: 1 },
		duration: 500,
		easing: 'cubic-bezier(0.4, 0, 0.2, 1)',
	},
	'fade-up': {
		initial: { opacity: 0, transform: 'translateY(12px)' },
		animate: { opacity: 1, transform: 'translateY(0)' },
		duration: 600,
		easing: 'cubic-bezier(0.4, 0, 0.2, 1)',
	},
	'fade-in-scale': {
		initial: { opacity: 0, transform: 'scale(0.98)' },
		animate: { opacity: 1, transform: 'scale(1)' },
		duration: 500,
		easing: 'cubic-bezier(0.4, 0, 0.2, 1)',
	},
	'fade-down': {
		initial: { opacity: 0, transform: 'translateY(-12px)' },
		animate: { opacity: 1, transform: 'translateY(0)' },
		duration: 600,
		easing: 'cubic-bezier(0.4, 0, 0.2, 1)',
	},
} as const;

export const STAGGER_DELAY = 60; // 60ms between items

export const OBSERVER_CONFIG = {
	threshold: 0.15,
	rootMargin: '0px 0px -50px 0px', // Early trigger
};

export type AnimationPreset = keyof typeof ANIMATION_PRESETS;
