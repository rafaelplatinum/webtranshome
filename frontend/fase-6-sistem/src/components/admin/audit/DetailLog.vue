<script setup>
import { computed } from 'vue'
import Button from '@/components/ui/Button.vue'
import Modal from '@/components/ui/Modal.vue'
import PerbandinganData from './PerbandinganData.vue'
import { formatTanggalJam } from '@/utils/format'
import { labelAksi } from '@/utils/log'

const open = defineModel('open', { type: Boolean, default: false })

const props = defineProps({
  // Satu item GET /admin/activity-logs (null saat drawer tertutup)
  log: { type: Object, default: null },
})

const info = computed(() => {
  const l = props.log
  if (!l) return []
  return [
    ['Admin', l.admin_name],
    ['Waktu', formatTanggalJam(l.created_at)],
    ['Modul', l.module],
    ['Kode aksi', l.action],
    ['Alamat IP', l.ip_address ?? '—'],
  ]
})
</script>

<template>
  <Modal v-model:open="open" side size="lg" :title="log ? labelAksi(log.action) : 'Detail log'" description="Perbandingan data sebelum dan sesudah aksi ini.">
    <template v-if="log">
      <dl class="grid grid-cols-[auto_minmax(0,1fr)] gap-x-4 gap-y-1.5 text-sm">
        <template v-for="[label, nilai] in info" :key="label">
          <dt class="text-muted">{{ label }}</dt>
          <dd class="font-semibold break-words" :class="label === 'Kode aksi' && 'font-mono text-xs font-normal'">{{ nilai }}</dd>
        </template>
      </dl>
      <PerbandinganData :sebelum="log.before_data" :sesudah="log.after_data" />
    </template>
    <template #footer>
      <Button variant="outline" @click="open = false">Tutup</Button>
    </template>
  </Modal>
</template>
