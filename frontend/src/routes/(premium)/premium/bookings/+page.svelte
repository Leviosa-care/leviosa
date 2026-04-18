<script lang="ts">
    import type { PageProps } from "./$types";

    import { getUserContext } from "$lib/context/user.svelte";
    import type { Event, Consultation } from "$lib/types/bookings";

    let { data }: PageProps = $props();

    const {
        user: { role },
    } = getUserContext();

    interface Props {
        consultations: Promise<Consultation[]>;
        events?: Promise<Event[]>;
    }
    const { consultations, events }: Props = data;
    // TODO: handle the tab thing display based on the size of the content
</script>

<div>the user role is: {role}</div>

{@render handleConsultation()}
{@render handleEvent()}

{#snippet handleConsultation()}
    {#await consultations}
        <p>waiting for the consultations, put some skeleton UI there</p>
    {:then consultations}
        {#each consultations as consultation}
            <div>{consultation.name}</div>
        {/each}
    {/await}
{/snippet}

{#snippet handleEvent()}
    {#await events}
        <p>waiting for the events, put some skeleton UI there</p>
    {:then events}
        {#if events}
            {#each events as event}
                <div>{event.title}</div>
            {/each}
        {/if}
    {/await}
{/snippet}
