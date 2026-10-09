<script setup>
import { ref, watch } from 'vue'
import { ChevronDown, CircleCheck, CircleX } from 'lucide-vue-next'
import Button from '@/components/ui/Button.vue'
import ErrorState from '@/components/ui/ErrorState.vue'
import Modal from '@/components/ui/Modal.vue'
import Skeleton from '@/components/ui/Skeleton.vue'
import { pesanError } from '@/services/errors'
import { getContactLogs } from '@/services/qontakService'
import { formatTanggalJam } from '@/utils/format'

const open = defineModel('open', { type: Boolean, default: false })

const props = defineProps({
  // Baris kontak yang dibuka: { id, name, member_code }
  kontak: { type: Object, default: null },
})

const status = ref('loading') // loading | ready | error
const logs = ref([])
const pesan = ref('')
let urutan = 0

async function muat() {
  if (!props.kontak) return
  const nomor = ++urutan
  status.value = 'loading'
  try {
    const hasil = await getContactLogs(props.kontak.id)
    if (nomor !== urutan) return
    logs.value = hasil
    status.value = 'ready'
  } catch (error) {
    if (nomor !== urutan) return
    pesan.value = pesanError(error)
    status.value = 'error'
  }
}

watch(open, (buka) => buka && muat())

const json = (isi) => (isi == null ? '—' : JSON.stringify(isi, null, 2))
const berhasil = (log) => log.status_code === 200
</script>

<template>
  <Modal
    v-model:open="open"
    side
    size="lg"
    :title="kontak ? `Log sync ${kontak.name}` : 'Log sync'"
    description="Setiap percobaan kirim kontak ke Qontak, terbaru di atas."
  >
    <div v-if="status === 'loading'" class="flex flex-col gap-3" aria-busy="true">
      <span class="sr-only">Memuat log sync…</span>
      <Skeleton v-for="n in 3" :key="n" class="h-28 w-full" />
    </div>
    <ErrorState v-else-if="status === 'error'" :message="pesan" @retry="muat" />
    <p v-else-if="!logs.length" class="text-sm text-muted">Belum ada percobaan kirim untuk kontak ini.</p>
    <ol v-else class="flex flex-col gap-3">
      <li v-for="log in logs" :key="log.id" class="flex flex-col gap-2 rounded-lg border border-border p-3" data-log-sync>
        <div class="flex flex-wrap items-center justify-between gap-2">
          <span
            class="inline-flex items-center gap-1 rounded px-2 py-0.5 text-xs font-semibold"
            :class="berhasil(log) ? 'bg-success-soft text-success-ink' : 'bg-danger-soft text-danger-ink'"
          >
            <component :is="berhasil(log) ? CircleCheck : CircleX" class="size-3.5" aria-hidden="true" />
            {{ berhasil(log) ? 'Berhasil' : 'Gagal' }} · status {{ log.status_code ?? '—' }}
          </span>
          <time :datetime="log.created_at" class="text-sm text-muted tabular-nums">{{ formatTanggalJam(log.created_at) }}</time>
        </div>
        <p class="text-sm">
          <span class="font-mono text-xs">{{ log.event_type }}</span> ·
          {{ log.retry_count === 0 ? 'Kiriman pertama' : `Percobaan ulang ke-${log.retry_count}` }}
        </p>
        <details class="group text-sm">
          <summary class="flex min-h-11 cursor-pointer list-none items-center gap-1.5 rounded-sm font-semibold focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring md:min-h-9 [&::-webkit-details-marker]:hidden">
            <ChevronDown class="size-4 -rotate-90 transition-transform duration-150 group-open:rotate-0 motion-reduce:transition-none" aria-hidden="true" />
            Isi permintaan dan respons
          </summary>
          <div class="mt-2 flex flex-col gap-2">
            <p class="text-xs font-semibold text-muted uppercase">Permintaan</p>
            <pre class="overflow-x-auto rounded-md bg-background p-3 font-mono text-xs">{{ json(log.request_payload) }}</pre>
            <p class="text-xs font-semibold text-muted uppercase">Respons</p>
            <pre class="overflow-x-auto rounded-md bg-background p-3 font-mono text-xs">{{ json(log.response_payload) }}</pre>
          </div>
        </details>
      </li>
    </ol>
    <template #footer>
      <Button variant="outline" @click="open = false">Tutup</Button>
    </template>
  </Modal>
</template>
