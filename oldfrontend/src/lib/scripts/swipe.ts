// types
type swipeHorizontalCallBack = (direction: 'left' | 'right') => void;
type swipeVerticalCallBack = (direction: 'top' | 'bottom') => void;

// functions
export function createHorizontalSwipeHandler(
	onSwipe: swipeHorizontalCallBack,
	swipeTriggerValue: number = 50
) {
	let startX = 0;
	let endX = 0;

	function handleTouchStart(event: TouchEvent): void {
		startX = event.touches[0].clientX;
	}
	function handleTouchMove(event: TouchEvent): void {
		endX = event.touches[0].clientX;
	}
	function handleTouchEnd(): void {
		if (startX - endX > swipeTriggerValue && endX > 0) onSwipe('left');
		else if (endX - startX > swipeTriggerValue) onSwipe('right');

		endX = 0;
	}

	return {
		action: (node: HTMLElement) => {
			node.style.touchAction = 'pan-y';
			node.addEventListener('touchstart', handleTouchStart);
			node.addEventListener('touchmove', handleTouchMove);
			node.addEventListener('touchend', handleTouchEnd);
			return {
				destroy: () => {
					node.removeEventListener('touchstart', handleTouchStart);
					node.removeEventListener('touchmove', handleTouchMove);
					node.removeEventListener('touchend', handleTouchEnd);
				}
			};
		}
	};
}

export function createVerticalSwipeHandler(
	onSwipe: swipeVerticalCallBack,
	swipeTriggerValue: number = 50
) {
	let startY = 0;
	let endY = 0;

	function handleTouchStart(event: TouchEvent): void {
		startY = event.touches[0].clientY;
	}
	function handleTouchMove(event: TouchEvent): void {
		endY = event.touches[0].clientY;
	}
	function handleTouchEnd(): void {
		if (startY - endY > swipeTriggerValue && endY > 0) onSwipe('top');
		else if (endY - startY > swipeTriggerValue) onSwipe('bottom');
		// reinitialise the value of endX so that I always catch new swipe that do not depend on previous values
		endY = 0;
	}

	return {
		action: (node: HTMLElement) => {
			node.style.touchAction = 'pan-x';
			node.addEventListener('touchstart', handleTouchStart);
			node.addEventListener('touchmove', handleTouchMove);
			node.addEventListener('touchend', handleTouchEnd);
			return {
				destroy: () => {
					node.removeEventListener('touchstart', handleTouchStart);
					node.removeEventListener('touchmove', handleTouchMove);
					node.removeEventListener('touchend', handleTouchEnd);
				}
			};
		}
	};
}
