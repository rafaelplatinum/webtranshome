<script setup>
import { computed } from 'vue'
import { CircleAlert, CircleCheck, CircleX, Clock } from 'lucide-vue-next'
import { LABEL_STOK } from '@/utils/produk'

const props = defineProps({
  // Nilai products.stock_status
  status: {
    type: String,
    required: true,
    validator: (v) => ['TERSEDIA', 'SISA_STOK', 'PRE_ORDER', 'HABIS'].includes(v),
  },
  // products.stock_qty_label, contoh "3 dus"
  keterangan: { type: String, default: '' },
})

// Status dibedakan dengan label + ikon, bukan warna saja (CLAUDE.md §7).
const STATUS = {
  TERSEDIA: { ikon: CircleCheck, kelas: 'bg-success-soft text-success-ink' },
  SISA_STOK: { ikon: CircleAlert, kelas: 'bg-warning-soft text-warning-ink' },
  PRE_ORDER: { ikon: Clock, kelas: 'bg-info-soft text-info-ink' },
  HABIS: { ikon: CircleX, kelas: 'bg-danger-soft text-danger-ink' },
}

const info = computed(() => ({ ...STATUS[props.status], label: LABEL_STOK[props.status] }))
</script>

<template>
  <span class="inline-flex items-center gap-1.5 self-start rounded px-2 py-1 text-xs font-semibold" :class="info.kelas">
    <component :is="info.ikon" class="size-3.5 shrink-0" aria-hidden="true" />
    <span>{{ info.label }}<template v-if="keterangan"> · {{ keterangan }}</template></span>
  </span>
</template>
