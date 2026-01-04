<script lang="ts">
</script>

<div class="grid grid-cols-4 gap-8">
    {@render card("consultations (jour)", 42)}
    {@render card("ca du mois", "28k $")}
    {@render card("practiciens actifs", 154)}
    {@render card("a traiter", 8)}
</div>

{#snippet card(title: string, figure: string | number)}
    <div
        class="bg-white border border-dark-50 rounded-lg p-4 w-full grid gap-2 shadow-mini"
    >
        <p class="text-dark-700 text-md font-sm uppercase">{title}</p>
        <p class="text-dark-900 text-3xl font-bold">{figure}</p>
        <p class="text-sm">+12 % vs semaine derniere</p>
    </div>
{/snippet}
