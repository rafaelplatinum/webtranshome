<script setup>
import { ref } from 'vue'
import Badge from '@/components/ui/Badge.vue'
import CopyButton from '@/components/ui/CopyButton.vue'
import { useUiStore } from '@/stores/ui'
import { inisialNama } from '@/utils/format'

const props = defineProps({
  // Detail GET /admin/members/:id
  member: { type: Object, required: true },
})

const LABEL_TIER = { TAHAP_1: 'Tahap 1' }

const ui = useUiStore()
const teksKode = ref(null)

// Clipboard ditolak browser: kode dipilih otomatis supaya tinggal disalin manual.
function salinGagal() {
  if (teksKode.value) window.getSelection()?.selectAllChildren(teksKode.value)
  ui.tampilkanToast({ pesan: `Kode ${props.member.member_code} sudah dipilih. Salin secara manual.`, jenis: 'info' })
}
</script>

<template>
  <section aria-label="Profil member" class="flex flex-wrap items-center gap-4 rounded-xl border border-border bg-surface p-5">
    <span class="flex size-14 shrink-0 items-center justify-center rounded-full bg-primary-soft text-lg font-extrabold text-primary-hover" aria-hidden="true">
      {{ inisialNama(member.full_name) }}
    </span>
    <div class="flex min-w-0 flex-[1_1_14rem] flex-col gap-1.5">
      <h2 class="text-xl font-extrabold [overflow-wrap:anywhere]">{{ member.full_name }}</h2>
      <div class="flex flex-wrap items-center gap-2">
        <span ref="teksKode" class="font-mono text-base font-semibold" data-kode-member>{{ member.member_code }}</span>
        <CopyButton :text="member.member_code" label="Salin" copied-label="Disalin" @failed="salinGagal" />
        <Badge variant="accent">{{ LABEL_TIER[member.tier] ?? member.tier }}</Badge>
        <Badge :variant="member.is_active ? 'success' : 'danger'" data-status-akun>{{ member.is_active ? 'Aktif' : 'Nonaktif' }}</Badge>
      </div>
    </div>
    <div v-if="$slots.aksi" class="flex w-full flex-wrap gap-2 sm:w-auto">
      <slot name="aksi" />
    </div>
  </section>
</template>
