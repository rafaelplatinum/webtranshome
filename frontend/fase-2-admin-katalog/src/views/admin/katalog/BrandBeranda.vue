<script setup>
import { computed, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { RouterLink } from 'vue-router'
import { Monitor, Plus, Search, Smartphone, X } from 'lucide-vue-next'
import Button from '@/components/ui/Button.vue'
import ConfirmModal from '@/components/ui/ConfirmModal.vue'
import ErrorState from '@/components/ui/ErrorState.vue'
import Input from '@/components/ui/Input.vue'
import Skeleton from '@/components/ui/Skeleton.vue'
import TombolUrut from '@/components/admin/TombolUrut.vue'
import SusunanLingkaranBrand from '@/components/beranda/SusunanLingkaranBrand.vue'
import { useBrandBeranda } from '@/composables/useBrandBeranda'
import { useKonfirmasiKeluar } from '@/composables/useKonfirmasiKeluar'
import { usePermission } from '@/composables/usePermission'
import { brandLingkaran, LINGKARAN_HP, MAKS_LINGKARAN_BRAND } from '@/utils/brandBeranda'

// /admin/brand-beranda (KEPUTUSAN #21): fokus mengatur brand di lingkaran "Brand pilihan" Beranda.
// Kiri: urutan tampil (1–8 juga tampil di HP). Kanan: tambah brand aktif. Atas: pratinjau Beranda (desktop / HP).
const d = useBrandBeranda()
const { can } = usePermission()
const bolehUbah = computed(() => can('brand.update'))
const keluar = useKonfirmasiKeluar(d.kotor)

// Pratinjau memakai komponen yang sama dengan Beranda, dirender selebar frame asli lalu diperkecil agar muat.
const LEBAR = { desktop: 1232, hp: 343 } // frame Beranda: max-w-7xl − padding; HP 375 px − 2 × 16 px
const mode = ref('desktop')
const bingkai = ref(null)
const isiPratinjau = ref(null)
const lebarBingkai = ref(0)
const tinggiIsi = ref(0)
const skala = computed(() => (lebarBingkai.value ? Math.min(1, lebarBingkai.value / LEBAR[mode.value]) : 1))
let pengamat = null
onMounted(() => {
  pengamat = new ResizeObserver(() => {
    lebarBingkai.value = bingkai.value?.clientWidth ?? 0
    tinggiIsi.value = isiPratinjau.value?.offsetHeight ?? 0
  })
})
watch([bingkai, isiPratinjau], ([b, i]) => {
  pengamat?.disconnect()
  if (b) pengamat?.observe(b)
  if (i) pengamat?.observe(i)
})
onBeforeUnmount(() => pengamat?.disconnect())

const brandPratinjau = computed(() => (d.dipilih.value.length ? d.dipilih.value : brandLingkaran(d.aktif.value).brands))

async function simpan() {
  await d.simpan()
}
</script>

<template>
  <div class="flex flex-col gap-6">
    <p class="max-w-3xl text-sm text-muted">
      Pilih dan urutkan brand untuk seksi "Siap untuk menjelajah lebih banyak Brand?" di Beranda. Desktop menampilkan
      <strong class="text-foreground">{{ MAKS_LINGKARAN_BRAND }} lingkaran</strong> (8 kiri, 8 kanan, utuh di dalam frame); HP dan tablet menampilkan
      <strong class="text-foreground">{{ LINGKARAN_HP }} teratas</strong>. Logo diunggah di
      <RouterLink to="/admin/brand" class="rounded-sm font-semibold text-primary hover:underline focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring">Katalog › Brand</RouterLink>.
    </p>

    <div v-if="d.status.value === 'loading'" class="flex flex-col gap-4" aria-busy="true">
      <span class="sr-only">Memuat brand…</span>
      <Skeleton class="h-72 rounded-xl" />
      <Skeleton class="h-64 rounded-xl" />
    </div>
    <ErrorState v-else-if="d.status.value === 'error'" :message="d.pesan.value" @retry="d.muat" />

    <template v-else>
      <section aria-labelledby="judul-pratinjau-lingkaran" class="flex flex-col gap-3 rounded-xl border border-border bg-surface p-4" data-pratinjau-lingkaran>
        <div class="flex flex-wrap items-center justify-between gap-3">
          <h2 id="judul-pratinjau-lingkaran" class="text-base font-bold">Pratinjau Beranda</h2>
          <div role="group" aria-label="Ukuran pratinjau" class="inline-flex rounded-md border border-border-strong p-0.5">
            <button
              v-for="(ikon, kunciMode) in { desktop: Monitor, hp: Smartphone }"
              :key="kunciMode"
              type="button"
              :aria-pressed="mode === kunciMode ? 'true' : 'false'"
              class="inline-flex min-h-10 items-center gap-1.5 rounded px-3 text-sm font-semibold transition-colors duration-150 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
              :class="mode === kunciMode ? 'bg-secondary text-secondary-foreground' : 'hover:bg-subtle'"
              @click="mode = kunciMode"
            >
              <component :is="ikon" class="size-4" aria-hidden="true" />
              {{ kunciMode === 'desktop' ? 'Desktop' : 'HP' }}
            </button>
          </div>
        </div>
        <p v-if="!d.dipilih.value.length" class="rounded-md bg-info-soft px-3 py-2 text-sm text-info-ink">
          Belum ada brand yang dipilih, jadi Beranda memilih otomatis seperti pratinjau ini.
        </p>
        <!-- Pratinjau hanya gambaran: tidak bisa diklik/difokus dan disembunyikan dari pembaca layar. -->
        <div ref="bingkai" class="relative overflow-hidden rounded-lg bg-background" :style="{ height: `${Math.round(tinggiIsi * skala) + 24}px` }">
          <div
            ref="isiPratinjau"
            class="absolute top-3 left-1/2"
            :style="{ width: `${LEBAR[mode]}px`, transform: `translateX(-50%) scale(${skala})`, transformOrigin: 'top center' }"
            inert
            aria-hidden="true"
          >
            <SusunanLingkaranBrand :brands="brandPratinjau" :total="d.aktif.value.length" id-judul="judul-pratinjau-brand" />
          </div>
        </div>
      </section>

      <div class="grid items-start gap-6 xl:grid-cols-2">
        <section aria-labelledby="judul-urutan-lingkaran" class="flex flex-col gap-3 rounded-xl border border-border bg-surface p-4">
          <div class="flex flex-wrap items-baseline justify-between gap-2">
            <h2 id="judul-urutan-lingkaran" class="text-base font-bold">Urutan tampil</h2>
            <span class="text-sm font-semibold" :class="d.penuh.value ? 'text-warning-ink' : 'text-muted'" aria-live="polite" data-jumlah-lingkaran>
              {{ d.dipilih.value.length }} dari {{ MAKS_LINGKARAN_BRAND }} lingkaran{{ d.penuh.value ? ' (penuh)' : '' }}
            </span>
          </div>
          <p class="text-sm text-muted">Nomor 1 paling dekat judul. Nomor 1–{{ LINGKARAN_HP }} juga tampil di HP.</p>
          <p v-if="d.nonaktifDicentang.value.length" class="rounded-md bg-warning-soft px-3 py-2 text-sm text-warning-ink">
            {{ d.nonaktifDicentang.value.map((b) => b.name).join(', ') }} dicentang tetapi nonaktif, jadi tidak tampil dan akan dilepas saat disimpan.
          </p>
          <p v-if="!d.dipilih.value.length" class="rounded-md border border-dashed border-border-strong px-4 py-6 text-center text-sm text-muted">
            Belum ada brand. Tambahkan dari daftar "Tambah brand".
          </p>
          <ol v-else class="flex flex-col divide-y divide-border" data-urutan-lingkaran>
            <li v-for="(b, i) in d.dipilih.value" :key="b.id" class="flex items-center gap-3 py-2" :data-brand-terpilih="b.id">
              <span class="w-6 shrink-0 text-right font-mono text-sm font-semibold text-muted">{{ i + 1 }}</span>
              <span class="flex size-10 shrink-0 items-center justify-center overflow-hidden rounded-full border border-border bg-subtle text-sm font-extrabold text-muted">
                <img v-if="b.logo_url" :src="b.logo_url" alt="" width="40" height="40" class="size-full object-contain p-1" />
                <span v-else aria-hidden="true">{{ b.name.charAt(0).toUpperCase() }}</span>
              </span>
              <span class="flex min-w-0 flex-1 flex-col">
                <span class="truncate font-semibold">{{ b.name }}</span>
                <span class="flex flex-wrap gap-1.5 text-xs">
                  <span v-if="i < LINGKARAN_HP" class="rounded bg-info-soft px-1.5 py-0.5 font-semibold text-info-ink">Tampil di HP</span>
                  <span v-if="!b.logo_url" class="rounded bg-subtle px-1.5 py-0.5 font-semibold text-muted">Tanpa logo (tampil nama)</span>
                </span>
              </span>
              <template v-if="bolehUbah">
                <TombolUrut :id="b.id" :nama="b.name" :pertama="i === 0" :terakhir="i === d.dipilih.value.length - 1" @geser="(arah) => d.geser(b, arah)" />
                <Button variant="ghost" size="sm" :data-lepas-brand="b.id" @click="d.lepas(b)">
                  <template #icon><X class="size-4" aria-hidden="true" /></template>
                  Lepas<span class="sr-only"> {{ b.name }}</span>
                </Button>
              </template>
            </li>
          </ol>
        </section>

        <section aria-labelledby="judul-tambah-lingkaran" class="flex flex-col gap-3 rounded-xl border border-border bg-surface p-4">
          <h2 id="judul-tambah-lingkaran" class="text-base font-bold">Tambah brand</h2>
          <div data-cari-brand-beranda>
            <Input v-model="d.cari.value" label="Cari brand aktif" type="search" placeholder="Nama brand" autocomplete="off">
              <template #akhir><Search class="mr-2.5 size-4.5 text-muted" aria-hidden="true" /></template>
            </Input>
          </div>
          <p v-if="d.penuh.value" class="rounded-md bg-warning-soft px-3 py-2 text-sm text-warning-ink">
            Lingkaran sudah penuh ({{ MAKS_LINGKARAN_BRAND }}). Lepas satu brand dulu untuk menambah yang lain.
          </p>
          <p v-if="!d.tersedia.value.length" class="px-1 py-4 text-sm text-muted">
            {{ d.cari.value.trim() ? `Tidak ada brand aktif bernama “${d.cari.value.trim()}” yang belum dipilih.` : 'Semua brand aktif sudah dipilih.' }}
          </p>
          <ul v-else class="flex max-h-[28rem] flex-col divide-y divide-border overflow-y-auto pr-1" data-brand-tersedia>
            <li v-for="b in d.tersedia.value" :key="b.id" class="flex items-center gap-3 py-2">
              <span class="flex size-10 shrink-0 items-center justify-center overflow-hidden rounded-full border border-border bg-subtle text-sm font-extrabold text-muted">
                <img v-if="b.logo_url" :src="b.logo_url" alt="" width="40" height="40" class="size-full object-contain p-1" />
                <span v-else aria-hidden="true">{{ b.name.charAt(0).toUpperCase() }}</span>
              </span>
              <span class="flex min-w-0 flex-1 flex-col">
                <span class="truncate font-semibold">{{ b.name }}</span>
                <span class="text-xs text-muted">{{ b.product_count }} produk{{ b.logo_url ? '' : ' · tanpa logo' }}</span>
              </span>
              <Button v-if="bolehUbah" variant="outline" size="sm" :disabled="d.penuh.value" :data-tambah-brand="b.id" @click="d.tambah(b)">
                <template #icon><Plus class="size-4" aria-hidden="true" /></template>
                Tambah<span class="sr-only"> {{ b.name }}</span>
              </Button>
            </li>
          </ul>
        </section>
      </div>

      <p v-if="!bolehUbah" class="text-sm text-muted">Akun ini hanya bisa melihat. Perubahan butuh izin mengubah brand.</p>

      <div v-if="bolehUbah" class="sticky bottom-0 z-20 -mx-4 flex flex-wrap items-center justify-end gap-2 border-t border-border bg-surface px-4 py-3 sm:-mx-7 sm:px-7">
        <p v-if="d.galat.value" role="alert" class="mr-auto text-sm font-semibold text-danger-ink" data-galat>{{ d.galat.value }}</p>
        <span v-else-if="d.kotor.value" class="mr-auto hidden text-sm text-muted sm:inline">Ada perubahan yang belum disimpan</span>
        <Button variant="outline" :disabled="!d.kotor.value || d.menyimpan.value" @click="d.batal">Batalkan perubahan</Button>
        <Button :loading="d.menyimpan.value" :disabled="!d.kotor.value" data-simpan-lingkaran @click="simpan">Simpan</Button>
      </div>

      <ConfirmModal
        :open="keluar.terbuka.value"
        title="Buang perubahan?"
        message="Urutan brand di Beranda yang belum disimpan akan hilang."
        confirm-label="Buang perubahan"
        cancel-label="Tetap di sini"
        danger
        @confirm="keluar.buang"
        @cancel="keluar.tetap"
      />
    </template>
  </div>
</template>
