<script setup>
import { ref, useId } from 'vue'
import { ImagePlus, X } from 'lucide-vue-next'
import Button from '@/components/ui/Button.vue'
import { errorField, pesanError } from '@/services/errors'
import { ATURAN_UNGGAH, cekFile, uploadFile } from '@/services/uploadService'

// Satu gambar: gambar kategori, logo brand, cover ruangan.
const url = defineModel({ type: String, default: '' })

defineProps({
  label: { type: String, required: true },
  required: { type: Boolean, default: false },
  error: { type: String, default: '' },
  hint: { type: String, default: 'JPG, PNG, atau WebP, maksimal 2 MB.' },
  // Rasio pratinjau, mis. 'aspect-4/3' untuk cover ruangan, 'aspect-square' untuk logo.
  rasio: { type: String, default: 'aspect-square' },
  // Lebar pratinjau; gambar lebar (banner) memakai pratinjau yang lebih besar.
  lebar: { type: String, default: 'w-28' },
})

const id = useId()
const input = ref(null)
const sedang = ref(false)
const galat = ref('')

async function onPilih(event) {
  const file = event.target.files?.[0]
  event.target.value = ''
  if (!file) return
  galat.value = cekFile(file, 'gambar')
  if (galat.value) return
  sedang.value = true
  try {
    url.value = (await uploadFile(file, 'gambar')).url
  } catch (error) {
    galat.value = errorField(error).file ?? pesanError(error)
  } finally {
    sedang.value = false
  }
}
</script>

<template>
  <div class="flex flex-col gap-1.5" role="group" :aria-labelledby="`${id}-label`">
    <span :id="`${id}-label`" class="text-sm font-semibold">{{ label }}<span v-if="required" class="text-danger" aria-hidden="true"> *</span></span>
    <div class="flex flex-wrap items-start gap-3">
      <div class="max-w-full shrink-0 overflow-hidden rounded-lg border bg-subtle" :class="[rasio, lebar, error ? 'border-danger' : 'border-border']">
        <img v-if="url" :src="url" alt="" class="size-full object-cover" />
        <span v-else class="flex size-full items-center justify-center text-faint"><ImagePlus class="size-6" aria-hidden="true" /></span>
      </div>
      <div class="flex min-w-0 flex-1 basis-48 flex-col items-start gap-1.5">
        <Button variant="outline" size="sm" :loading="sedang" @click="input?.click()">{{ url ? 'Ganti gambar' : 'Unggah gambar' }}</Button>
        <button
          v-if="url && !required"
          type="button"
          class="inline-flex min-h-10 items-center gap-1 rounded-md px-1 text-sm font-semibold text-muted transition-colors duration-150 hover:text-danger focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
          @click="url = ''"
        >
          <X class="size-4" aria-hidden="true" />
          Hapus gambar
        </button>
        <p class="text-xs text-muted">{{ hint }}</p>
      </div>
    </div>
    <input ref="input" type="file" class="sr-only" :accept="ATURAN_UNGGAH.gambar.accept" tabindex="-1" aria-hidden="true" @change="onPilih" />
    <p v-if="galat || error" role="alert" class="text-sm font-medium text-danger-ink" data-galat>{{ galat || error }}</p>
  </div>
</template>
