<script lang="ts">
    import type { PageProps } from "./$types";

    let { data }: PageProps = $props();
</script>

<div class="coming-soon">
    <h1>Événements</h1>
    <p>Cette fonctionnalité arrive bientôt.</p>
</div>

<style>
    .coming-soon {
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        padding: 4rem 2rem;
        text-align: center;
    }

    h1 {
        font-size: 2rem;
        font-weight: 600;
        margin-bottom: 1rem;
        color: #1f2937;
    }

    p {
        font-size: 1.125rem;
        color: #6b7280;
    }
</style>
