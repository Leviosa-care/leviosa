<script lang="ts">
	import type { EventPickerMonth } from '$lib/types';

	import TimePickerButton from '$lib/components/forms/TimePickerButton.svelte';
	import Drawer from '$lib/components/Drawer.svelte';
	import { XIcon, ChevronLeft, ChevronRight } from 'lucide-svelte';

	import { createHorizontalSwipeHandler } from '$lib/scripts/swipe';
	import { convertMonthToInt } from '$lib/scripts/date';

	const leftArrowClass = 'leftArrow';
	const rightArrowClass = 'rightArrow';
	interface Props {
		monthsData: EventPickerMonth[];
	}

	let { monthsData }: Props = $props();

	// TODO: get that thing in the local storage somewhat ?
	let slidePosition: number = $state(0);

	function formatTime(time: number): string {
		if (time > 9) {
			return time.toString();
		}
		return '0' + time.toString();
	}

	// a function to make the month appear in the carousel
	function handleCarousel(e: MouseEvent) {
		const arrow = e.currentTarget as HTMLButtonElement;
		const swipeLeftCondition = arrow.classList.contains(leftArrowClass) && slidePosition > 0;
		const swipeRightCondition =
			arrow.classList.contains(rightArrowClass) && slidePosition < monthsData.length - 1;
		if (swipeLeftCondition) slidePosition -= 1;
		else if (swipeRightCondition) slidePosition += 1;
	}

	function formatDay(day: string): string {
		return day.slice(0, 2);
	}

	let isOpen = $state(false);

	// TODO: export that thing
	let value: string | undefined = $state('Selectionne un creneau');
	function handleSpotReservation(e: MouseEvent) {
		const target = e.currentTarget as HTMLButtonElement;
		const hour = target.querySelector('.hour')?.textContent?.trim();
		const month = target.querySelector('.month')?.textContent?.trim();
		const day = target.querySelector('.day')?.textContent?.trim();
		const date = target.querySelector('.date')?.textContent?.trim();

		value = `${day} ${date} ${month} - ${hour}`;

		isOpen = false;
	}
	// TODO: make a function that can converts the months into its index. It is for the label that gives the date under the day

	function swipeCalendar(direction: 'left' | 'right') {
		if (direction == 'left') {
			const rightButton = document.querySelector(`.${rightArrowClass}`) as HTMLButtonElement;
			rightButton.click();
		} else if (direction == 'right') {
			const leftButton = document.querySelector(`.${leftArrowClass}`) as HTMLButtonElement;
			leftButton.click();
		}
	}
	const swipeCalendarAction = createHorizontalSwipeHandler(swipeCalendar);
</script>

<TimePickerButton {value} onClick={() => (isOpen = !isOpen)} />
<Drawer bind:isOpen>
	<div class="holder stroke">
		{#each monthsData as month, monthIndex}
			<div
				class="slide flow"
				style="transform: translateX(calc(-100% * {slidePosition})); --flow-space: 1rem;"
			>
				<div class="grid">
					<div class="flex">
						<p class="drawer-title dark-primary-content">Choisis un creneau</p>
						<button type="button" class="close flex" onclick={() => (isOpen = false)}>
							<XIcon />
						</button>
					</div>
					<div class="separator"></div>
				</div>
				<div class="header" use:swipeCalendarAction.action>
					<button
						type="button"
						onclick={handleCarousel}
						class="arrow stroke {leftArrowClass} dark-ternary-content"
						style:visibility={data[monthIndex - 1] != undefined ? 'visibile' : 'hidden'}
					>
						<ChevronLeft />
					</button>
					<p class="fw-600 fs-h3">{month.name}</p>
					<button
						type="button"
						onclick={handleCarousel}
						class="arrow stroke {rightArrowClass} dark-ternary-content"
						style:visibility={data[monthIndex + 1] != undefined ? 'visibile' : 'hidden'}
					>
						<ChevronRight />
					</button>
				</div>
				<div class="calendar flex" use:swipeCalendarAction.action>
					{#each month.info as content}
						<div class="column flex">
							<div class="calendar-header flex" style="--gap: 0.1rem">
								<p class="dark-primary-content">
									{content.day}
								</p>
								<p class="dark-ternary-content">
									{formatTime(content.date)}/{convertMonthToInt(month.name)}
								</p>
							</div>
							<div class="horaires flex">
								{#each content.hours as horraire}
									<button type="button" class="creneau" onclick={handleSpotReservation}>
										<p class="hour">
											{formatTime(horraire)}:00
										</p>
										<p class="date" hidden>
											{content.date}
										</p>
										<p class="month" hidden>{month.name}</p>
										<p class="day" hidden>{content.day}</p>
									</button>
								{/each}
							</div>
						</div>
					{/each}
				</div>
			</div>
		{/each}
	</div>
</Drawer>

<style>
	.holder {
		width: calc(100vw - 1.5rem);
		position: absolute;
		bottom: 0;
		left: 0;
		right: 0;
		width: 100vw;
		display: flex;
		overflow-x: hidden;
		padding-bottom: 2rem;
		padding-block: 0rem;
		background-color: hsl(var(--clr-light-primary));
		border: none;
	}
	.close {
		width: 100%;
		background: transparent;
		align-items: center;
		justify-content: right;
		opacity: 0.6;
	}
	.drawer-title {
		text-wrap: nowrap;
		font-weight: 600;
	}
	.separator {
		width: 100%;
		height: 1px;
		background-color: hsl(var(--clr-stroke));
	}
	.holder,
	.slide {
		margin-inline: auto;
		border-radius: 0.5rem;
	}

	.slide {
		min-width: 100%;
		margin-inline: auto;
		padding: 1rem;
		/* NOTE: the flex part that allow to hide the next months */
		flex-grow: 1;
		transition: transform 0.3s ease;
	}

	.header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		gap: 1rem;
	}
	.arrow {
		aspect-ratio: 1;
		background: transparent;
		padding: 0.25rem;
		border-radius: 0.5rem;
	}

	.column,
	.calendar-header,
	.horaires {
		flex-direction: column;
		align-items: center;
	}
	.horaires button {
		padding: 0.5rem;
		border-radius: 0.5rem;
		background-color: hsl(var(--clr-dark-primary));
		color: hsl(var(--clr-light-primary));
	}
</style>
