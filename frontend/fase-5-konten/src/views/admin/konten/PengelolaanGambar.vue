<script setup>
import { computed, onMounted, reactive, ref } from 'vue'
import { ArrowDown, ArrowUp, ImagePlus, RefreshCw, Trash2 } from 'lucide-vue-next'
import Button from '@/components/ui/Button.vue'
import EmptyState from '@/components/ui/EmptyState.vue'
import ErrorState from '@/components/ui/ErrorState.vue'
import Skeleton from '@/components/ui/Skeleton.vue'
import { errorField, pesanError } from '@/services/errors'
import { createImageAsset, deleteImageAsset, listImageAssets, updateImageAsset } from '@/services/imageAssetService'

const KELOMPOK = [
  { judul: 'Beranda', slot: [
    { key: 'home.hero', label: 'Hero beranda', rasio: '5:4', saran: '1920 × 800 px' },
    { key: 'home.promo_banner', label: 'Banner promo beranda', rasio: '5:2', saran: '1600 × 640 px' },
    { key: 'home.side_banner', label: 'Banner samping promo', rasio: '20:13', saran: '800 × 520 px' },
  ] },
  { judul: 'Katalog dan layanan', slot: [
    { key: 'category.tile', label: 'Tile kategori', rasio: '5:4', saran: '800 × 640 px' },
    { key: 'product.main', label: 'Foto utama produk', rasio: '1:1', saran: '1200 × 1200 px' },
    { key: 'product.gallery', label: 'Galeri produk', rasio: '1:1', saran: '1200 × 1200 px' },
    { key: 'brand.logo', label: 'Logo brand', rasio: '1:1', saran: '512 × 512 px, transparan' },
    { key: 'room.cover', label: 'Cover ruangan', rasio: '4:3', saran: '1200 × 900 px' },
    { key: 'service.cover', label: 'Foto jasa', rasio: '16:9', saran: '1200 × 675 px' },
  ] },
  { judul: 'Konten dan toko', slot: [
    { key: 'article.cover', label: 'Cover artikel', rasio: '16:9', saran: '1200 × 675 px' },
    { key: 'promo.cover', label: 'Cover promo', rasio: '16:9', saran: '1200 × 675 px' },
    { key: 'event.cover', label: 'Cover event', rasio: '16:9', saran: '1200 × 675 px' },
    { key: 'store.front', label: 'Foto tampak depan toko', rasio: '12:5', saran: '1920 × 800 px' },
    { key: 'store.map', label: 'Foto atau preview peta', rasio: '4:3', saran: '1200 × 900 px' },
  ] },
  { judul: 'Identitas situs', slot: [
    { key: 'site.logo', label: 'Logo situs', rasio: '1:1', saran: '512 × 512 px, transparan' },
    { key: 'site.favicon', label: 'Favicon', rasio: '1:1', saran: '32 × 32 px' },
    { key: 'site.og_image', label: 'Gambar pratinjau tautan', rasio: '1.91:1', saran: '1200 × 630 px' },
  ] },
]

const aset = ref([])
const alt = reactive({})
const loading = ref(true)
const error = ref('')
const busy = ref('')
const fileInputs = new Map()
const totalSlot = computed(() => KELOMPOK.reduce((sum, group) => sum + group.slot.length, 0))
const dalamSlot = (key) => aset.value.filter((item) => item.slot_key === key).sort((a, b) => a.sort_order - b.sort_order)

async function muat() {
  loading.value = true
  error.value = ''
  try {
    aset.value = await listImageAssets()
    aset.value.forEach((item) => { alt[item.id] = item.alt_text })
  } catch (err) {
    error.value = pesanError(err)
  } finally {
    loading.value = false
  }
}

function bukaPilih(key) { fileInputs.get(key)?.click() }

async function unggah(key, event) {
  const file = event.target.files?.[0]
  event.target.value = ''
  if (!file) return
  if (!['image/jpeg', 'image/png', 'image/webp'].includes(file.type) || file.size > 2 * 1024 * 1024) {
    error.value = 'Pilih gambar JPG, PNG, atau WebP dengan ukuran maksimal 2 MB.'
    return
  }
  const slot = KELOMPOK.flatMap((group) => group.slot).find((item) => item.key === key)
  busy.value = key
  error.value = ''
  try {
    const record = await createImageAsset({ file, slotKey: key, altText: `${slot.label} Transhome` })
    aset.value.push(record)
    alt[record.id] = record.alt_text
  } catch (err) {
    error.value = errorField(err).file ?? pesanError(err)
  } finally {
    busy.value = ''
  }
}

async function simpanAlt(item) {
  const nilai = String(alt[item.id] ?? '').trim()
  if (!nilai) { error.value = 'Teks alternatif wajib diisi.'; alt[item.id] = item.alt_text; return }
  try {
    const terbaru = await updateImageAsset(item.id, { alt_text: nilai })
    Object.assign(item, terbaru)
    error.value = ''
  } catch (err) { error.value = errorField(err).alt_text ?? pesanError(err) }
}

async function ubahStatus(item) {
  try { Object.assign(item, await updateImageAsset(item.id, { is_active: !item.is_active })) }
  catch (err) { error.value = pesanError(err) }
}

async function geser(item, arah) {
  const satuSlot = dalamSlot(item.slot_key)
  const index = satuSlot.findIndex((row) => row.id === item.id)
  const target = satuSlot[index + arah]
  if (!target) return
  const urutanItem = item.sort_order
  try {
    await Promise.all([
      updateImageAsset(item.id, { sort_order: target.sort_order }),
      updateImageAsset(target.id, { sort_order: urutanItem }),
    ])
    await muat()
  } catch (err) { error.value = pesanError(err) }
}

async function hapus(item) {
  busy.value = `delete-${item.id}`
  try {
    await deleteImageAsset(item.id)
    aset.value = aset.value.filter((row) => row.id !== item.id)
    delete alt[item.id]
  } catch (err) { error.value = pesanError(err) }
  finally { busy.value = '' }
}

onMounted(muat)
</script>

<template>
  <div class="flex flex-col gap-6">
    <div class="flex flex-wrap items-start justify-between gap-3">
      <div>
        <p class="text-sm text-muted">Kelola gambar per lokasi, teks alternatif, status tayang, dan urutannya.</p>
        <p class="mt-1 text-sm text-muted">{{ totalSlot }} slot gambar · JPG, PNG, atau WebP · maksimal 2 MB per file</p>
      </div>
      <Button variant="outline" :loading="loading" @click="muat">
        <template #icon><RefreshCw class="size-4" aria-hidden="true" /></template>
        Muat ulang
      </Button>
    </div>

    <p v-if="error" role="alert" class="rounded-lg border border-danger/30 bg-danger-soft px-4 py-3 text-sm text-danger-ink">{{ error }}</p>
    <div v-if="loading" class="grid gap-4 md:grid-cols-2 xl:grid-cols-3" aria-busy="true"><Skeleton v-for="n in 6" :key="n" class="h-52" /></div>
    <ErrorState v-else-if="!aset" message="Gambar belum dapat dimuat." @retry="muat" />

    <section v-for="group in KELOMPOK" :key="group.judul" class="flex flex-col gap-3">
      <h2 class="text-lg font-bold">{{ group.judul }}</h2>
      <div class="grid gap-4 md:grid-cols-2 xl:grid-cols-3">
        <article v-for="slot in group.slot" :key="slot.key" class="rounded-xl border border-border bg-surface p-4">
          <div class="mb-3 flex items-start justify-between gap-2">
            <div><h3 class="font-semibold">{{ slot.label }}</h3><p class="mt-1 text-xs text-muted">{{ slot.key }} · {{ slot.rasio }} · {{ slot.saran }}</p></div>
            <span class="shrink-0 rounded-full bg-subtle px-2.5 py-1 text-xs font-semibold text-muted">{{ dalamSlot(slot.key).length }}</span>
          </div>
          <div class="mb-3 flex flex-col gap-3">
            <div v-for="item in dalamSlot(slot.key)" :key="item.id" class="rounded-lg border border-border p-3">
              <div class="mb-3 flex aspect-video items-center justify-center overflow-hidden rounded-md bg-subtle">
                <img :src="item.image_url" :alt="alt[item.id] || item.alt_text" class="size-full object-cover" />
              </div>
              <label class="mb-1 block text-sm font-medium" :for="`alt-${item.id}`">Teks alternatif</label>
              <input :id="`alt-${item.id}`" v-model="alt[item.id]" maxlength="255" class="min-h-10 w-full rounded-md border border-border bg-surface px-3 text-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring" @blur="simpanAlt(item)" />
              <div class="mt-3 flex flex-wrap items-center gap-2">
                <Button variant="outline" size="sm" :disabled="dalamSlot(slot.key).indexOf(item) === 0" aria-label="Naikkan urutan" @click="geser(item, -1)"><template #icon><ArrowUp class="size-4" aria-hidden="true" /></template></Button>
                <Button variant="outline" size="sm" :disabled="dalamSlot(slot.key).indexOf(item) === dalamSlot(slot.key).length - 1" aria-label="Turunkan urutan" @click="geser(item, 1)"><template #icon><ArrowDown class="size-4" aria-hidden="true" /></template></Button>
                <Button variant="outline" size="sm" @click="ubahStatus(item)">{{ item.is_active ? 'Nonaktifkan' : 'Aktifkan' }}</Button>
                <Button variant="ghost" size="sm" :loading="busy === `delete-${item.id}`" @click="hapus(item)"><template #icon><Trash2 class="size-4" aria-hidden="true" /></template>Hapus</Button>
              </div>
            </div>
          </div>
          <input :ref="(el) => el && fileInputs.set(slot.key, el)" type="file" accept="image/jpeg,image/png,image/webp" class="sr-only" tabindex="-1" aria-hidden="true" @change="unggah(slot.key, $event)" />
          <Button variant="outline" size="sm" :loading="busy === slot.key" @click="bukaPilih(slot.key)">
            <template #icon><ImagePlus class="size-4" aria-hidden="true" /></template>
            Tambah gambar
          </Button>
        </article>
      </div>
    </section>
    <EmptyState v-if="!loading && !aset.length" class="sr-only" title="Belum ada gambar yang diunggah" description="Pilih slot lalu tambahkan gambar." />
  </div>
</template>
