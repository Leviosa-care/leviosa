<script>
	// TODO: find everything as variable that I need for this component but I need design first
</script>

<div class="container grid" style="--gap: 2rem;">
	<div class="card grid" style="--gap:0.5rem;">
		<h3 class="title">Deroule de la seance</h3>
		<p>
			Lorem ipsum dolor sit amet consectetur adipisicing elit. Vitae aliquam debitis aperiam
			exercitationem? Distinctio ipsa excepturi quibusdam
		</p>
		<div class="card-overlay"></div>
	</div>
	<div class="video">
		<div class="play"></div>
	</div>
	<div class="card grid" style="--gap:0.5rem;">
		<h3 class="title">De quel equipement as tu besoin</h3>
		<p>
			Lorem ipsum dolor sit amet consectetur adipisicing elit. Vitae aliquam debitis aperiam
			exercitationem? Distinctio ipsa excepturi quibusdam
		</p>
		<div class="card-overlay"></div>
	</div>
</div>

<style>
	.title {
		color: hsl(var(--clr-grey-100));
		font-size: var(--fs-1);
	}
	.video {
		background-color: #d9d9d9;
		width: 100%;
		aspect-ratio: 16/9;
		border-radius: 0.5rem;
		/* NOTE: just for the wireframe */
		display: grid;
		place-content: center;
	}
	.play {
		background-color: hsl(var(--clr-accent));
		width: 3.5rem;
		aspect-ratio: 1;
		border-radius: 100%;
		position: relative;
	}
	.play::before {
		content: '';
		position: absolute;
		top: 50%;
		left: 50%;
		transform: translate(-50%, -50%);
		background-color: white;
		border-top: 12px solid hsl(var(--clr-accent));
		border-left: 24px solid white;
		border-bottom: 12px solid hsl(var(--clr-accent));
	}
	.card {
		--card-border-radius: 1rem;
		--card-padding: 1rem;
		padding: var(--card-padding);
		border-radius: var(--card-border-radius);
		backdrop-filter: blur(15px) saturate(150%);
		box-shadow:
			rgba(60, 64, 67, 0.3) 0px 1px 2px 0px,
			rgba(60, 64, 67, 0.15) 0px 1px 3px 1px;
		position: relative;
		overflow: hidden;
	}
	.card-overlay {
		position: absolute;
		top: 0;
		left: 0;
		width: 100%;
		height: 100%;
		background-color: rgba(0, 0, 0, 0.25);
		z-index: -1;
	}
</style>
