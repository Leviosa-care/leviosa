<script>
    import { ArrowLeft } from "@lucide/svelte";
    import { Button } from "bits-ui";
</script>

<Button.Root
    class="mb-4 justify-center gap-2 items-center h-input rounded-input hover:bg-muted focus-visible:ring-dark focus-visible:ring-offset-background focus-visible:outline-hidden inline-flex px-4 text-[15px] font-bold focus-visible:ring-2 focus-visible:ring-offset-2 active:scale-[0.98] cursor-pointer"
    type="button"
    onclick={() => history.back()}
>
    <ArrowLeft size={16} />
    <span class="text-xs sm:text-base">Retour à l'étape précédente</span>
</Button.Root>
