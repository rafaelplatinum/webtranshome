<script setup>
import { computed, nextTick, ref } from 'vue'
import { ArrowLeft, ArrowRight, ImagePlus, LoaderCircle, Trash2 } from 'lucide-vue-next'
import ConfirmModal from '@/components/ui/ConfirmModal.vue'
import { errorField, pesanError } from '@/services/errors'
import { ATURAN_UNGGAH, cekFile, uploadFile } from '@/services/uploadService'

// [{ id? (foto lama), kunci? (foto baru), image_url, is_primary }]; urutan array = sort_order.
const foto = defineModel({ type: Array, default: () => [] })

defineProps({
  error: { type: String, default: '' },
})

const input = ref(null)
const wadah = ref(null)
const galatFile = ref([])
const sedangUnggah = ref(0)
const indeksHapus = ref(null)
const diseret = ref(false)
let nomor = 0

const kunci = (f) => f.id ?? f.kunci

async function terima(daftar) {
  galatFile.value = []
  const valid = []
  for (const file of daftar) {
    const galat = cekFile(file, 'gambar')
    if (galat) galatFile.value.push(galat)
    else valid.push(file)
  }
  sedangUnggah.value += valid.length
  await Promise.all(
    valid.map(async (file) => {
      try {
        const hasil = await uploadFile(file, 'gambar')
        foto.value = [...foto.value, { kunci: `baru-${++nomor}`, image_url: hasil.url, is_primary: !foto.value.some((f) => f.is_primary) }]
      } catch (error) {
        galatFile.value.push(`${file.name}: ${errorField(error).file ?? pesanError(error)}`)
      } finally {
        sedangUnggah.value--
      }
    }),
  )
}

function onPilih(event) {
  terima([...event.target.files])
  event.target.value = ''
}

function onLepas(event) {
  diseret.value = false
  terima([...event.dataTransfer.files])
}

function jadikanUtama(i) {
  foto.value = foto.value.map((f, j) => ({ ...f, is_primary: j === i }))
}

// Geser dengan tombol (bisa keyboard); fokus ikut foto yang dipindah.
async function geser(i, arah) {
  const j = i + arah
  const baru = [...foto.value]
  ;[baru[i], baru[j]] = [baru[j], baru[i]]
  foto.value = baru
  await nextTick()
  const tombol = wadah.value?.querySelector(`[data-geser="${j}:${arah}"]:not(:disabled)`) ?? wadah.value?.querySelector(`[data-geser^="${j}:"]:not(:disabled)`)
  tombol?.focus()
}

// Foto utama dihapus → foto pertama yang tersisa jadi utama.
function hapus() {
  const i = indeksHapus.value
  const baru = foto.value.filter((_, j) => j !== i)
  if (foto.value[i]?.is_primary && baru.length) baru[0] = { ...baru[0], is_primary: true }
  foto.value = baru
  indeksHapus.value = null
}

const pesanHapus = computed(() =>
  foto.value[indeksHapus.value]?.is_primary && foto.value.length > 1
    ? 'Ini foto utama. Setelah dihapus, foto pertama yang tersisa menjadi foto utama.'
    : 'Foto akan dilepas dari produk ini saat kamu menyimpan.',
)

const kelasIkon =
  'inline-flex size-11 items-center justify-center rounded-md text-muted transition-colors duration-150 hover:bg-subtle hover:text-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-40 md:size-9'
</script>

<template>
  <div ref="wadah" class="flex flex-col gap-3">
    <ul class="grid grid-cols-2 gap-3 sm:grid-cols-3 2xl:grid-cols-4">
      <li v-for="(f, i) in foto" :key="kunci(f)" class="flex flex-col overflow-hidden rounded-lg border-2" :class="f.is_primary ? 'border-primary' : 'border-border'">
        <div class="relative aspect-square bg-subtle">
          <img :src="f.image_url" :alt="`Foto ${i + 1}${f.is_primary ? ' (utama)' : ''}`" width="240" height="240" class="size-full object-cover" />
          <span v-if="f.is_primary" class="absolute top-2 left-2 rounded bg-primary px-2 py-0.5 text-xs font-bold text-primary-foreground">Utama</span>
        </div>
        <div class="flex flex-wrap items-center justify-between gap-1 p-1">
          <button
            v-if="!f.is_primary"
            type="button"
            class="inline-flex min-h-11 items-center rounded-md px-2 text-xs font-semibold text-primary transition-colors duration-150 hover:bg-primary-soft focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring md:min-h-9"
            @click="jadikanUtama(i)"
          >
            Jadikan utama
          </button>
          <span v-else class="px-2 text-xs font-semibold text-primary-hover">Foto utama</span>
          <div class="flex">
            <button type="button" :class="kelasIkon" :data-geser="`${i}:-1`" :disabled="i === 0" :aria-label="`Geser foto ${i + 1} ke kiri`" @click="geser(i, -1)">
              <ArrowLeft class="size-4" aria-hidden="true" />
            </button>
            <button type="button" :class="kelasIkon" :data-geser="`${i}:1`" :disabled="i === foto.length - 1" :aria-label="`Geser foto ${i + 1} ke kanan`" @click="geser(i, 1)">
              <ArrowRight class="size-4" aria-hidden="true" />
            </button>
            <button type="button" :class="[kelasIkon, 'hover:text-danger']" :aria-label="`Hapus foto ${i + 1}`" @click="indeksHapus = i">
              <Trash2 class="size-4" aria-hidden="true" />
            </button>
          </div>
        </div>
      </li>
      <li>
        <button
          type="button"
          class="flex aspect-square w-full flex-col items-center justify-center gap-1.5 rounded-lg border-2 border-dashed p-3 text-center transition-colors duration-150 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
          :class="diseret ? 'border-primary bg-primary-soft' : 'border-border-strong bg-background hover:border-foreground'"
          :aria-busy="sedangUnggah ? 'true' : undefined"
          @click="input?.click()"
          @dragover.prevent="diseret = true"
          @dragleave="diseret = false"
          @drop.prevent="onLepas"
        >
          <LoaderCircle v-if="sedangUnggah" class="size-6 animate-spin text-muted" aria-hidden="true" />
          <ImagePlus v-else class="size-6 text-muted" aria-hidden="true" />
          <span class="text-sm font-semibold">{{ sedangUnggah ? 'Mengunggah…' : 'Unggah foto' }}</span>
          <span class="text-xs text-muted">JPG, PNG, WebP · maks 2 MB · bisa beberapa sekaligus</span>
        </button>
      </li>
    </ul>
    <input ref="input" type="file" class="sr-only" :accept="ATURAN_UNGGAH.gambar.accept" multiple tabindex="-1" aria-hidden="true" @change="onPilih" />
    <ul v-if="galatFile.length" role="alert" class="flex flex-col gap-1 text-sm font-medium text-danger-ink">
      <li v-for="g in galatFile" :key="g">{{ g }}</li>
    </ul>
    <p v-if="error" class="text-sm font-medium text-danger-ink" data-galat>{{ error }}</p>

    <ConfirmModal
      :open="indeksHapus !== null"
      title="Hapus foto ini?"
      :message="pesanHapus"
      confirm-label="Hapus foto"
      danger
      @update:open="(buka) => !buka && (indeksHapus = null)"
      @confirm="hapus"
    />
  </div>
</template>
