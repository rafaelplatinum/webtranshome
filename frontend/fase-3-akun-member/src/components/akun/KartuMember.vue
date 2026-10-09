<script setup>
import { computed, ref } from 'vue'
import CopyButton from '@/components/ui/CopyButton.vue'
import { useUiStore } from '@/stores/ui'
import { formatTanggal } from '@/utils/format'

const props = defineProps({
  // user dari GET /me: { full_name, member_code, tier, member_since }
  user: { type: Object, required: true },
})

const LABEL_TIER = { TAHAP_1: 'Tahap 1' }

const ui = useUiStore()
const teksKode = ref(null)
const tier = computed(() => LABEL_TIER[props.user.tier] ?? props.user.tier ?? '')

// Clipboard ditolak browser: kode dipilih otomatis supaya tinggal disalin manual.
function salinGagal() {
  if (teksKode.value) window.getSelection()?.selectAllChildren(teksKode.value)
  ui.tampilkanToast({ pesan: 'Kode member sudah dipilih. Salin secara manual.', jenis: 'info' })
}
</script>

<template>
  <!-- Kartu member digital: kode member disebut ke kasir sebagai pengganti kartu fisik (KEPUTUSAN #2). -->
  <section aria-label="Kartu member Trans Family" class="flex flex-col gap-5 rounded-2xl bg-secondary p-6 text-secondary-foreground sm:p-7">
    <div class="flex items-center justify-between gap-3">
      <strong class="text-base">Trans Family</strong>
      <span v-if="tier" class="rounded bg-accent px-2.5 py-0.5 text-xs font-bold text-accent-foreground">{{ tier }}</span>
    </div>
    <div class="flex flex-col gap-1">
      <span class="text-xs text-on-dark-muted">Kode member · sebutkan ke kasir saat belanja</span>
      <div class="flex flex-wrap items-center gap-3">
        <span ref="teksKode" class="font-mono text-3xl font-semibold tracking-wider" data-kode-member>{{ user.member_code }}</span>
        <CopyButton gelap :text="user.member_code" label="Salin" copied-label="Disalin" @failed="salinGagal" />
      </div>
    </div>
    <span class="text-sm text-on-dark">
      {{ user.full_name }}<template v-if="user.member_since"> · member sejak {{ formatTanggal(user.member_since) }}</template>
    </span>
  </section>
</template>
