<script setup>
import { computed, onMounted, ref } from 'vue'
import { Save } from 'lucide-vue-next'
import Button from '@/components/ui/Button.vue'
import Checkbox from '@/components/ui/Checkbox.vue'
import ErrorState from '@/components/ui/ErrorState.vue'
import Input from '@/components/ui/Input.vue'
import Skeleton from '@/components/ui/Skeleton.vue'
import { listBrands } from '@/services/catalogService'
import { getAdminSettings, updateSettings } from '@/services/adminKontenService'
import { pesanError } from '@/services/errors'
import { useSettingsStore } from '@/stores/settings'
import { useUiStore } from '@/stores/ui'
import { formatAngka } from '@/utils/format'

const MAKS_BRAND = 18
const status = ref('loading')
const pesan = ref('')
const brands = ref([])
const pilihan = ref([])
const awal = ref('[]')
const cari = ref('')
const menyimpan = ref(false)
const galatSimpan = ref('')
const settings = useSettingsStore()
const ui = useUiStore()

const urutRekomendasi = (daftar) => [...daftar].sort(
  (a, b) => Number(Boolean(b.logo_url)) - Number(Boolean(a.logo_url)) ||
    (b.product_count ?? 0) - (a.product_count ?? 0) ||
    a.name.localeCompare(b.name, 'id'),
)
const daftarTampil = computed(() => brands.value
  .filter((brand) => brand.name.toLocaleLowerCase('id').includes(cari.value.trim().toLocaleLowerCase('id')))
  .sort((a, b) => a.name.localeCompare(b.name, 'id')))
const jumlahDipilih = computed(() => pilihan.value.length)
const berubah = computed(() => JSON.stringify([...pilihan.value].sort((a, b) => a - b)) !== awal.value)

function mulai(daftar, setting) {
  // listBrands() already returns active brands and omits is_active from its public shape.
  brands.value = daftar.filter((brand) => brand.is_active !== false)
  const tersimpan = setting.home_featured_brand_ids
  const awalDisimpan = Array.isArray(tersimpan)
    ? tersimpan.map(Number).filter((id) => brands.value.some((brand) => brand.id === id)).slice(0, MAKS_BRAND)
    : urutRekomendasi(brands.value).slice(0, MAKS_BRAND).map((brand) => brand.id)
  pilihan.value = awalDisimpan
  awal.value = JSON.stringify([...awalDisimpan].sort((a, b) => a - b))
}

async function muat() {
  status.value = 'loading'
  pesan.value = ''
  try {
    const [daftar, setting] = await Promise.all([listBrands(), getAdminSettings()])
    mulai(daftar, setting)
    status.value = 'ready'
  } catch (error) {
    pesan.value = pesanError(error)
    status.value = 'error'
  }
}

function ubahPilihan(id, aktif) {
  galatSimpan.value = ''
  if (aktif && !pilihan.value.includes(id) && pilihan.value.length < MAKS_BRAND) pilihan.value = [...pilihan.value, id]
  if (!aktif) pilihan.value = pilihan.value.filter((pilihanId) => pilihanId !== id)
}

function batalkan() {
  pilihan.value = JSON.parse(awal.value)
  galatSimpan.value = ''
}

async function simpan() {
  if (!berubah.value || jumlahDipilih.value < 1 || jumlahDipilih.value > MAKS_BRAND) return
  menyimpan.value = true
  galatSimpan.value = ''
  try {
    const home_featured_brand_ids = brands.value.filter((brand) => pilihan.value.includes(brand.id)).map((brand) => brand.id)
    const hasil = await updateSettings({ home_featured_brand_ids })
    mulai(brands.value, hasil)
    await settings.load({ force: true })
    ui.tampilkanToast({ pesan: 'Pilihan brand Beranda disimpan', jenis: 'success' })
  } catch (error) {
    galatSimpan.value = pesanError(error)
  } finally {
    menyimpan.value = false
  }
}

onMounted(muat)
</script>

<template>
  <section aria-labelledby="judul-brand-beranda" class="flex max-w-3xl flex-col gap-4 rounded-xl border border-border bg-surface p-4 sm:p-5" data-pilihan-brand-admin>
    <div class="flex flex-wrap items-start justify-between gap-3">
      <div class="flex flex-col gap-1">
        <h2 id="judul-brand-beranda" class="text-base font-bold">Brand pilihan Beranda</h2>
        <p class="max-w-2xl text-sm text-muted">Pilih brand yang muncul di lingkaran Beranda. Maksimal 18 brand; di layar ponsel tampil hingga 10. Brand lain tetap ada di halaman Brand.</p>
      </div>
      <span v-if="status === 'ready'" class="shrink-0 rounded-full bg-primary-soft px-3 py-1 text-sm font-semibold text-primary-hover" aria-live="polite">
        {{ jumlahDipilih }} / {{ MAKS_BRAND }} dipilih
      </span>
    </div>

    <div v-if="status === 'loading'" class="flex flex-col gap-3" aria-busy="true">
      <Skeleton class="h-11 max-w-md" />
      <Skeleton v-for="n in 5" :key="n" class="h-11" />
    </div>
    <ErrorState v-else-if="status === 'error'" :message="pesan" @retry="muat" />
    <template v-else>
      <Input v-model="cari" label="Cari brand" type="search" placeholder="Nama brand" autocomplete="off" />
      <p v-if="galatSimpan" role="alert" class="rounded-md bg-danger-soft px-3 py-2 text-sm font-semibold text-danger-ink">{{ galatSimpan }}</p>
      <p v-if="!jumlahDipilih" role="alert" class="rounded-md bg-danger-soft px-3 py-2 text-sm text-danger-ink">Pilih minimal satu brand agar bagian Brand pilihan tetap tampil di Beranda.</p>

      <div v-if="daftarTampil.length" class="grid max-h-80 grid-cols-1 gap-x-4 overflow-y-auto rounded-lg border border-border px-3 sm:grid-cols-2" aria-label="Daftar brand aktif">
        <Checkbox
          v-for="brand in daftarTampil"
          :key="brand.id"
          :model-value="pilihan.includes(brand.id)"
          :label="brand.name"
          :description="`${formatAngka(brand.product_count ?? 0)} produk`"
          :disabled="!pilihan.includes(brand.id) && jumlahDipilih >= MAKS_BRAND"
          @update:model-value="(aktif) => ubahPilihan(brand.id, aktif)"
        />
      </div>
      <p v-else class="rounded-lg border border-dashed border-border-strong px-4 py-6 text-center text-sm text-muted">Tidak ada brand yang cocok dengan pencarian.</p>

      <div class="flex flex-wrap justify-end gap-2">
        <Button v-if="berubah" variant="outline" :disabled="menyimpan" @click="batalkan">Batalkan perubahan</Button>
        <Button :disabled="!berubah || !jumlahDipilih || menyimpan" :loading="menyimpan" @click="simpan">
          <template #icon><Save class="size-4.5" aria-hidden="true" /></template>
          Simpan pilihan
        </Button>
      </div>
    </template>
  </section>
</template>
