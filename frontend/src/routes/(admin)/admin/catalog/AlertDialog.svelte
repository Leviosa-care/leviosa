<script lang="ts">
    import { Dialog, Button } from "bits-ui";
    import { X, Info } from "@lucide/svelte";

    type Props = {
        children: import("svelte").Snippet;
    };

    let { children }: Props = $props();
</script>

<Dialog.Root>
    <Dialog.Trigger>
        {@render children()}
    </Dialog.Trigger>
    <Dialog.Portal>
        <Dialog.Overlay
            class="data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 fixed inset-0 z-50 bg-black/80"
        />
        <Dialog.Content
            class="rounded-card-lg bg-background shadow-popover data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 data-[state=closed]:zoom-out-95 data-[state=open]:zoom-in-95 outline-hidden fixed left-[50%] top-[50%] z-50 w-full max-w-[calc(100%-2rem)] translate-x-[-50%] translate-y-[-50%] border p-8 sm:max-w-[540px] md:w-full"
        >
            <div class="flex items-start gap-8">
                <Info size={48} />
                <div>
                    <Dialog.Title
                        class="w-full text-xl font-semibold tracking-tight"
                    >
                        Êtes-vous sûr ?
                    </Dialog.Title>
                    <div class="pb-4">
                        Êtes-vous sûr de vouloir supprimer la carte pour '' ?
                        Toutes les données de celle ci seront supprimées.
                    </div>
                </div>
            </div>
            <div class="flex w-full justify-end mt-4">
                <Dialog.Close
                    class="h-input rounded-input bg-dark text-background shadow-mini hover:bg-dark/95 focus-visible:ring-dark focus-visible:ring-offset-background focus-visible:outline-hidden inline-flex items-center justify-center px-8 text-[15px] font-semibold focus-visible:ring-2 focus-visible:ring-offset-2 active:scale-[0.98] cursor-pointer"
                >
                    <Button.Root type="submit" class="cursor-pointer"
                        >Confirmer</Button.Root
                    >
                </Dialog.Close>
            </div>
            <Dialog.Close
                class="focus-visible:ring-foreground focus-visible:ring-offset-background focus-visible:outline-hidden absolute right-5 top-5 rounded-md focus-visible:ring-2 focus-visible:ring-offset-2 active:scale-[0.98] cursor-pointer"
            >
                <Button.Root type="button" class="cursor-pointer">
                    <X class="text-foreground size-5" />
                    <span class="sr-only">Close</span>
                </Button.Root>
            </Dialog.Close>
        </Dialog.Content>
    </Dialog.Portal>
</Dialog.Root>
