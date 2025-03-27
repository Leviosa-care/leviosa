<script lang="ts">
	import { Slack, ArrowRight } from 'lucide-svelte';
	import Button from '$lib/components/Button.svelte';
	import Header from '../Header.svelte';
	const otpCount = 6;
	// NOTE: useful links for these things :
	// - https://codingtorque.com/otp-input-field-using-html-css/
	// - https://www.geeksforgeeks.org/create-otp-input-field-using-html-css-and-javascript/

	function focusNextTarget(e) {
		const target = e.target;
		const val = target.value;
		if (isNaN(val)) {
			target.value = '';
			return;
		}
		if (val != '') {
			const next = target?.nextElementSibling;
			if (next) {
				next.focus();
			}
		}
	}

	function focusPrevTarget(e) {
		const target = e.target;
		const key = e.key.toLowerCase();
		if (key === 'backspace' || key === 'delete') {
			// target?.value = ""
			if (target) target.value = '';
			const prev = target?.previousElementSibling;
			if (prev) {
				prev.focus();
			}
			return;
		}
	}

	function handleOnPasteOTP(e) {
		const data = e.clipboardData.getData('text');
		const value = data.split('');
		if (value.length === otpCount) {
			// TODO: get all the data to send using some function in the backend
			// inputs.forEach((input, index) => (input.value = value[index]));
			// TODO: use some sort action instead of this
			// submit();
		}
	}
	import { goto } from '$app/navigation';
	function handleSubmit() {
		console.log('click');
		goto('pending');
	}
</script>

<div class="content container flex">
	<div class="header">
		<Header selectedElement={3} pathname="password" />
	</div>
	<div class="body">
		<div class="grid" style="--gap: 2rem;">
			<div class="logo_placeholder">
				<Slack size={24} />
			</div>
			<div class="flow center" style="--flow-space:0.5rem;">
				<h3 class="title">Renseigne le code recu par mail</h3>
				<p>
					Un code vient d'etre envoye au mail <strong>john.doe@gmail.com</strong>
				</p>
			</div>
		</div>
		<form action="" class="grid" style="--gap: 3rem;">
			<div class="numbers flex" style="--gap: 1rem">
				{#each Array(6) as _, index}
					<input
						oninput={focusNextTarget}
						onkeyup={focusPrevTarget}
						onpaste={handleOnPasteOTP}
						type="text"
						inputmode="numeric"
						maxlength="1"
						id={`${index}`}
						class="number-input stroke fs-h2"
					/>
				{/each}
			</div>
			<Button rightIcon={ArrowRight} onClick={handleSubmit} style="primary">
				Finalise ton inscription
			</Button>
		</form>
	</div>
	<div></div>
</div>

<style>
	.content {
		flex-direction: column;
		justify-content: space-between;
		align-items: center;
		min-height: 100vh;
	}
	.header {
		flex: none;
		width: 100%;
	}
	.body {
		flex: 1;
		display: grid;
		place-content: center;
	}
	.title {
		color: black;
		font-size: 1.1rem;
		font-weight: 500;
	}
	.logo_placeholder {
		width: 40px;
		margin-inline: auto;
		aspect-ratio: 1;
		background-color: hsl(var(--clr-light-secondary));
		display: grid;
		place-content: center;
		border-radius: 0.5rem;
	}
	form {
		margin-top: 5rem;
	}
	.numbers {
		justify-content: center;
		align-items: center;
	}
	.number-input {
		/* TODO: make the cells responsive in size using some sort of clamp brother */
		--dimension: 2.5rem;
		width: var(--dimension);
		aspect-ratio: 1;
		/* margin-top: 2rem; */
		border-radius: 0.5rem;
		text-align: center;
	}
	.number-input:is(:global(:hover, :focus)) {
		border-color: hsl(var(--clr-dark-ternary));
	}
	h3 {
		font-weight: 600;
	}
</style>
