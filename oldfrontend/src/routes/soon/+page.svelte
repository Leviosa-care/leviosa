<script lang="ts">
	import { Instagram } from 'lucide-svelte';

	function redirectInstagram() {
		window.location.href = 'https://www.instagram.com/leviosa_care/';
	}
</script>

<div class="container">
	<div class="content">
		<h1 class="title">Bientôt, un nouveau souffle pour votre bien-être.</h1>
		<p class="subtitle">
			Découvrez une nouvelle façon de prendre soin de vous et de partager des moments uniques.
			<br />Rejoignez nous sur Instagram dès maintenant.
		</p>
	</div>
	<button class="btn" onclick={() => redirectInstagram()}>
		<Instagram size="40" />
	</button>
</div>

<style>
	.btn {
		display: grid;
		place-content: center;
		margin-top: 2rem;
		padding: 1rem;
		border-radius: 0.5rem;
		background: none;
	}
	.container {
		display: grid;
		place-content: center;
		height: 100vh;
		/* width: 65%; */
		margin-inline: auto;
	}

	.content {
		text-align: center;
		padding: 2rem;
	}

	/* .content { */
	/* 	background-color: #202127; */
	/* 	padding: 4rem; */
	/* 	padding-block: 3rem; */
	/* 	border-radius: 0.5rem; */
	/* 	outline: 1px solid #bdbdbd; */
	/* 	box-shadow: */
	/* 		rgba(0, 0, 0, 0.02) 0px 1px 3px 0px, */
	/* 		rgba(27, 31, 35, 0.15) 0px 0px 0px 1px; */
	/* } */

	.title,
	.subtitle {
		text-align: center;
	}

	.title {
		/* color: #f67373; */
		color: black;
		font-size: clamp(1.15rem, 4vw, 2.5rem);
		font-weight: 800;
	}
	.subtitle {
		color: #3c3c43;
		/* color: rgba(255, 255, 245, 0.86); */
		margin-top: 2rem;
		font-size: clamp(1rem, 1.5vw, 1.5rem);
	}

	::placeholder {
		font-weight: 300;
	}
</style>
