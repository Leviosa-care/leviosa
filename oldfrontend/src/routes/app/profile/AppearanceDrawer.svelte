<script>
	import RadioButton from './RadioButton.svelte';
</script>

<div class="accroche"></div>
<div class="grid" style="--gap: 2rem;">
	<p class="title">Appearance</p>
	<div class="btns" style="--gap: 0.5rem;">
		<RadioButton
			options={[
				{ id: '1', text: 'Mode clair' },
				{ id: '2', text: 'Mode sombre' }
			]}
			name="appearance"
		/>
	</div>
</div>

<style>
	.accroche {
		margin-inline: auto;
		background-color: hsl(var(--clr-dark-ternary));
		height: 5px;
		border-radius: 0.5rem;
		width: 8rem;
		margin-bottom: 2.5rem;
	}
	.title {
		font-size: var(--fs-1);
		font-weight: 600;
	}
</style>
