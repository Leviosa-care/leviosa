<script lang="ts">
    import BackButton from './BackButton.svelte';
    import Card from './Card.svelte';
    import FormInput from './FormInput.svelte';
    import { BUTTON_VARIANTS, CARD_VARIANTS, INPUT_VARIANTS } from './constants';

    function handleBack() {
        console.log('Constants test - Back clicked');
    }

    function handleCardClick() {
        console.log('Constants test - Card clicked');
    }

    let testInput = '';
</script>

<!-- Test BackButton with different constants -->
<div class="space-y-6 p-6">
    <h2 class="text-lg font-semibold mb-4">Constants Test Components</h2>

    <div class="space-y-4">
        <h3 class="text-md font-medium">BackButton (using constants)</h3>
        <BackButton onclick={handleBack} label="Back with Constants" />

        <h3 class="text-md font-medium">Card Components (using constants)</h3>
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            <!-- Default Card -->
            <Card variant={CARD_VARIANTS.DEFAULT} onclick={handleCardClick}>
                <div class="p-4">
                    <h4 class="font-medium">Default Card</h4>
                    <p class="text-sm text-muted-foreground">Using CARD_VARIANTS.DEFAULT</p>
                </div>
            </Card>

            <!-- Elevated Card -->
            <Card variant={CARD_VARIANTS.ELEVATED} onclick={handleCardClick}>
                <div class="p-4">
                    <h4 class="font-medium">Elevated Card</h4>
                    <p class="text-sm text-muted-foreground">Using CARD_VARIANTS.ELEVATED</p>
                </div>
            </Card>

            <!-- Outline Card -->
            <Card variant={CARD_VARIANTS.OUTLINE} onclick={handleCardClick}>
                <div class="p-4">
                    <h4 class="font-medium">Outline Card</h4>
                    <p class="text-sm text-muted-foreground">Using CARD_VARIANTS.OUTLINE</p>
                </div>
            </Card>
        </div>

        <h3 class="text-md font-medium">Form Input (using constants)</h3>
        <FormInput
            id="test-input"
            name="testInput"
            label="Test Input (using constants)"
            placeholder="Type something..."
            bind:value={testInput}
            variant={INPUT_VARIANTS.DEFAULT}
        />
    </div>
</div>