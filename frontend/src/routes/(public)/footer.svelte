<script lang="ts">
    import { Instagram } from "@lucide/svelte";
</script>

<footer class="bg-dark-900 text-white/80 pt-16 pb-8">
    <div class="container mx-auto px-4">
        <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-[2fr_1fr_1fr] gap-12 pb-12 border-b border-white/10">
            <div class="flex flex-col gap-4">
                <a href="/" class="text-xl font-bold text-white w-fit">
                    Leviosa
                </a>
                <p class="text-white/40 text-sm max-w-xs leading-relaxed">
                    Prenez soin de vous avec Leviosa Care.
                </p>
                <div class="mt-2 flex flex-col gap-1.5">
                    <p class="uppercase text-white/30 text-[10px] tracking-widest">Contact</p>
                    <a
                        href="mailto:contact@leviosa.care"
                        class="text-white/60 hover:text-white text-sm transition-colors w-fit"
                    >
                        contact@leviosa.care
                    </a>
                </div>
                <a
                    class="mt-2 flex items-center gap-2 text-white/60 hover:text-white text-sm transition-colors w-fit"
                    target="_blank"
                    href="https://www.instagram.com/leviosa_care"
                    rel="noopener noreferrer"
                >
                    <Instagram class="w-3.5 h-3.5" />
                    <span>@leviosa_care</span>
                </a>
            </div>

            <div class="flex flex-col gap-4">
                <p class="uppercase text-white/30 text-[10px] tracking-widest">Explorer</p>
                <nav class="flex flex-col gap-3">
                    <a href="/services" class="text-white/60 hover:text-white text-sm transition-colors w-fit">
                        Services
                    </a>
                    <a href="/team" class="text-white/60 hover:text-white text-sm transition-colors w-fit">
                        L'équipe
                    </a>
                    <a href="/about" class="text-white/60 hover:text-white text-sm transition-colors w-fit">
                        À propos
                    </a>
                    <a href="/book" class="text-white/60 hover:text-white text-sm transition-colors w-fit">
                        Réserver
                    </a>
                </nav>
            </div>

            <div class="flex flex-col gap-4">
                <p class="uppercase text-white/30 text-[10px] tracking-widest">Légal</p>
                <nav class="flex flex-col gap-3">
                    <a href="/legal/privacy" class="text-white/60 hover:text-white text-sm transition-colors w-fit">
                        Politique de confidentialité
                    </a>
                    <a href="/legal/terms" class="text-white/60 hover:text-white text-sm transition-colors w-fit">
                        Conditions générales
                    </a>
                </nav>
            </div>
        </div>

        <div class="pt-6 text-white/30 text-xs">
            <p>&copy; {new Date().getFullYear()} Leviosa Care. Tous droits réservés.</p>
        </div>
    </div>
</footer>
