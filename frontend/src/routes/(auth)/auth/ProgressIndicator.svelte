<script lang="ts">
    import { Check } from "@lucide/svelte";

    import type { Step } from "$lib/types/step";

    type Props = {
        steps: Step[];
    };

    let { steps }: Props = $props();
</script>

<div class="px-2 pt-2">
    <div class="flex items-center justify-between mb-2">
        {#each steps as step (step.id)}
            <div class="flex flex-col items-center">
                <div
                    class="flex h-8 w-8 sm:h-10 sm:w-10 items-center justify-center rounded-full text-sm font-medium {step.status ===
                    'complete'
                        ? 'bg-green-600 text-white'
                        : ''} {step.status === 'current'
                        ? 'border-2 border-green-600 text-green-600'
                        : ''} {step.status === 'upcoming'
                        ? 'border-2 border-gray-300 text-gray-500'
                        : ''}"
                >
                    {#if step.status === "complete"}
                        <Check class="w-4 h-4 sm:w-6 sm:h-6" />
                    {:else}
                        {step.id}
                    {/if}
                </div>
                <div class="mt-1 text-xs sm:text-base font-medium text-center">
                    <span
                        class="{step.status === 'complete'
                            ? 'text-green-600'
                            : ''} {step.status === 'current'
                            ? 'text-green-600'
                            : ''} {step.status === 'upcoming'
                            ? 'text-gray-500'
                            : ''}"
                    >
                        {step.name}
                    </span>
                </div>
            </div>
        {/each}
    </div>
    <div class="hidden sm:flex items-center" aria-hidden="true">
        {#each steps as step, i}
            <div id={`line-${i}`} class="flex-1 flex">
                {#if i != 0}
                    <div
                        class="h-0.5 w-full {steps[i - 1].status === 'complete'
                            ? 'bg-green-600'
                            : 'bg-gray-300'}"
                    ></div>
                {/if}
                {#if i != steps.length - 1}
                    <div
                        class="h-0.5 w-full {step.status === 'complete'
                            ? 'bg-green-600'
                            : 'bg-gray-300'}"
                    ></div>
                {/if}
            </div>
        {/each}
    </div>
</div>
