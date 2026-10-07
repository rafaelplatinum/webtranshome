<script setup>
import { computed, ref } from 'vue'
import { FileText, FileUp, X } from 'lucide-vue-next'
import Button from '@/components/ui/Button.vue'
import { errorField, pesanError } from '@/services/errors'
import { ATURAN_UNGGAH, cekFile, uploadFile } from '@/services/uploadService'

// products.datasheet_pdf_url
const url = defineModel({ type: String, default: '' })

const input = ref(null)
const sedang = ref(false)
const galat = ref('')
const namaFile = ref('')

const nama = computed(() => {
  if (namaFile.value) return namaFile.value
  if (!url.value || url.value.startsWith('data:')) return 'Datasheet.pdf'
  return decodeURIComponent(url.value.split('/').pop() || 'Datasheet.pdf')
})

async function onPilih(event) {
  const file = event.target.files?.[0]
  event.target.value = ''
  if (!file) return
  galat.value = cekFile(file, 'dokumen')
  if (galat.value) return
  sedang.value = true
  try {
    const hasil = await uploadFile(file, 'dokumen')
    url.value = hasil.url
    namaFile.value = hasil.name
  } catch (error) {
    galat.value = errorField(error).file ?? pesanError(error)
  } finally {
    sedang.value = false
  }
}

function hapus() {
  url.value = ''
  namaFile.value = ''
}
</script>

<template>
  <div class="flex flex-col gap-2">
    <div v-if="url" class="flex flex-wrap items-center gap-3 rounded-lg border border-border bg-background p-3">
      <FileText class="size-6 shrink-0 text-primary" aria-hidden="true" />
      <span class="min-w-0 flex-1 truncate text-sm font-semibold">{{ nama }}</span>
      <a
        :href="url"
        target="_blank"
        rel="noopener"
        class="inline-flex h-10 items-center rounded-md px-2.5 text-sm font-semibold text-primary transition-colors duration-150 hover:bg-primary-soft focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
      >
        Lihat<span class="sr-only"> datasheet (tab baru)</span>
      </a>
      <button
        type="button"
        class="inline-flex size-10 items-center justify-center rounded-md text-muted transition-colors duration-150 hover:bg-subtle hover:text-danger focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
        aria-label="Hapus datasheet"
        @click="hapus"
      >
        <X class="size-4.5" aria-hidden="true" />
      </button>
    </div>
    <div v-else class="flex flex-wrap items-center gap-3">
      <Button variant="outline" :loading="sedang" @click="input?.click()">
        <template #icon><FileUp class="size-4.5" aria-hidden="true" /></template>
        Unggah datasheet
      </Button>
      <span class="text-sm text-muted">PDF, maksimal 10 MB. Tampil sebagai tombol "Unduh datasheet" di website.</span>
    </div>
    <input ref="input" type="file" class="sr-only" :accept="ATURAN_UNGGAH.dokumen.accept" tabindex="-1" aria-hidden="true" @change="onPilih" />
    <p v-if="galat" role="alert" class="text-sm font-medium text-danger-ink">{{ galat }}</p>
  </div>
</template>
