<script setup>
import { computed, useId } from 'vue'
import Button from '@/components/ui/Button.vue'
import StatusSync from '@/components/admin/StatusSync.vue'
import { usePermission } from '@/composables/usePermission'
import { formatTanggalJam } from '@/utils/format'

const props = defineProps({
  // `qontak` dari detail member: { id, sync_status, last_synced_at, retry_count, status_code, sedang? }
  kontak: { type: Object, default: null },
  punyaHp: { type: Boolean, default: false },
})

// `coba-lagi`: induk menjalankan retry (useSyncQontak) dan memperbarui `kontak`.
const emit = defineEmits(['coba-lagi'])

const { can } = usePermission()
const idJudul = useId()

const keterangan = computed(() => {
  const k = props.kontak
  if (!k) return 'Belum ada data sync untuk member ini.'
  if (k.sync_status === 'PENDING') return 'Sedang dikirim ke Qontak…'
  if (k.sync_status === 'SYNCED') return k.last_synced_at ? `Terakhir tersinkron ${formatTanggalJam(k.last_synced_at)}` : 'Tersinkron'
  return `Percobaan ke-${k.retry_count} gagal · status ${k.status_code ?? '-'}`
})
</script>

<template>
  <section :aria-labelledby="idJudul" class="flex flex-col gap-3 rounded-xl border border-border bg-surface p-5">
    <div class="flex items-center justify-between gap-3">
      <h2 :id="idJudul" tabindex="-1" data-kartu-sync class="rounded-sm text-base font-bold">Sync Qontak</h2>
      <StatusSync :status="kontak?.sync_status ?? null" :memproses="Boolean(kontak?.sedang)" />
    </div>
    <p class="text-sm text-muted" aria-live="polite">{{ keterangan }}</p>
    <p v-if="!punyaHp" class="text-sm text-muted">Tanpa nomor HP, kontak tersinkron dengan email saja.</p>
    <!-- Hanya untuk yang gagal (tetap tampil selama proses supaya fokus tidak hilang). -->
    <Button
      v-if="kontak && can('qontak.retry') && (kontak.sync_status === 'FAILED' || kontak.sedang)"
      variant="outline"
      size="sm"
      class="self-start"
      :loading="Boolean(kontak.sedang)"
      @click="emit('coba-lagi')"
    >
      Coba sync ulang
    </Button>
  </section>
</template>
