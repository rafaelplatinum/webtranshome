<script setup>
import { computed, nextTick, onMounted, ref } from 'vue'
import { LoaderCircle, ShoppingCart, TriangleAlert } from 'lucide-vue-next'
import Breadcrumb from '@/components/ui/Breadcrumb.vue'
import Button from '@/components/ui/Button.vue'
import ConfirmModal from '@/components/ui/ConfirmModal.vue'
import EmptyState from '@/components/ui/EmptyState.vue'
import WhatsAppButton from '@/components/layout/WhatsAppButton.vue'
import BarisKeranjang from '@/components/keranjang/BarisKeranjang.vue'
import RingkasanKeranjang from '@/components/keranjang/RingkasanKeranjang.vue'
import { useSeo } from '@/composables/useSeo'
import { useWhatsApp } from '@/composables/useWhatsApp'
import { listProducts } from '@/services/productService'
import { MAKS_PRODUK_KERANJANG, useKeranjangStore } from '@/stores/keranjang'
import { useUiStore } from '@/stores/ui'
import { formatRupiah } from '@/utils/format'
import { hitungTotal } from '@/utils/templateWhatsApp'

// /keranjang (KEPUTUSAN #23): daftar produk + jumlah, lalu dikirim sekaligus lewat WhatsApp memakai template
// "Pesan keranjang" dari admin. Isi keranjang dari browser; harga & stok diperbarui dari API setiap halaman dibuka.
const keranjang = useKeranjangStore()
const ui = useUiStore()
const { keranjangLink, pesanDariKeranjang } = useWhatsApp()

useSeo({ title: 'Keranjang belanja', robots: 'noindex, follow' })

const status = ref('siap') // memuat | siap | gagal (pembaruan harga & stok)
const konfirmasiKosong = ref(false)
const judul = ref(null)

async function perbarui() {
  const ids = keranjang.items.map((i) => i.produk.id)
  if (!ids.length) return
  status.value = 'memuat'
  try {
    const { data } = await listProducts({ ids: ids.join(','), per_page: MAKS_PRODUK_KERANJANG })
    keranjang.perbarui(data)
    status.value = 'siap'
  } catch {
    status.value = 'gagal'
  }
}
onMounted(perbarui)

const dikirim = computed(() => keranjang.dijual.map(({ produk, jumlah }) => ({ produk, jumlah })))
const jumlahTidakDijual = computed(() => keranjang.items.length - dikirim.value.length)
const pesan = computed(() => (dikirim.value.length ? pesanDariKeranjang(dikirim.value) : ''))
const link = computed(() => keranjangLink(dikirim.value))
const total = computed(() => hitungTotal(dikirim.value))

// Hapus satu produk: langsung, dengan "Urungkan" di toast. Fokus pindah ke tombol hapus berikutnya.
async function hapus(item) {
  const posisi = keranjang.items.indexOf(item)
  const terhapus = keranjang.hapus(item.produk.id)
  if (!terhapus) return
  ui.tampilkanToast({
    pesan: `${item.produk.name} dihapus dari keranjang.`,
    jenis: 'info',
    durasi: 6000,
    aksi: { label: 'Urungkan', jalankan: () => keranjang.kembalikan(terhapus) },
  })
  await nextTick()
  const tombol = [...document.querySelectorAll('[data-hapus-keranjang]')]
  ;(tombol[posisi] ?? tombol[posisi - 1] ?? judul.value)?.focus()
}

function kosongkan() {
  keranjang.kosongkan()
  konfirmasiKosong.value = false
  ui.tampilkanToast({ pesan: 'Keranjang dikosongkan.', jenis: 'info' })
  nextTick(() => judul.value?.focus())
}
</script>

<template>
  <div class="mx-auto max-w-7xl px-4 pt-4 pb-12 sm:px-6 lg:pt-6 lg:pb-16">
    <Breadcrumb :items="[{ label: 'Beranda', to: '/' }, { label: 'Keranjang belanja' }]" />
    <div class="mt-2 flex flex-wrap items-baseline justify-between gap-x-4 gap-y-1">
      <h1 ref="judul" tabindex="-1" class="text-2xl font-extrabold focus:outline-none sm:text-3xl">Keranjang belanja</h1>
      <p v-if="keranjang.jumlahProduk" class="text-muted" data-jumlah-keranjang>{{ keranjang.jumlahProduk }} dari maks. {{ MAKS_PRODUK_KERANJANG }} produk</p>
    </div>
    <p class="mt-2 max-w-2xl text-muted">
      Kumpulkan produk yang ingin dibeli, lalu kirim daftarnya sekaligus ke WhatsApp Transhome. Tim kami membalas dengan stok dan total harganya.
    </p>

    <EmptyState
      v-if="!keranjang.jumlahProduk"
      class="mt-6"
      title="Keranjang masih kosong"
      description="Tambahkan produk dari katalog dengan tombol keranjang, lalu kirim daftarnya lewat WhatsApp."
    >
      <template #icon><ShoppingCart class="size-10" aria-hidden="true" /></template>
      <template #actions><Button to="/katalog">Lihat katalog</Button></template>
    </EmptyState>

    <div v-else class="mt-6 grid items-start gap-6 lg:grid-cols-[minmax(0,1fr)_22rem]">
      <section aria-labelledby="judul-daftar-keranjang" class="flex min-w-0 flex-col gap-3">
        <h2 id="judul-daftar-keranjang" class="sr-only">Produk di keranjang</h2>
        <p v-if="status === 'memuat'" class="flex items-center gap-2 text-sm text-muted" role="status">
          <LoaderCircle class="size-4 animate-spin motion-reduce:animate-none" aria-hidden="true" />
          Memperbarui harga dan stok…
        </p>
        <div v-else-if="status === 'gagal'" class="flex flex-wrap items-center gap-3 rounded-lg bg-warning-soft px-4 py-3 text-sm text-warning-ink" role="alert">
          <TriangleAlert class="size-4 shrink-0" aria-hidden="true" />
          <span class="min-w-0 flex-1">Harga dan stok belum bisa diperbarui. Yang tampil adalah data saat produk ditambahkan.</span>
          <Button variant="outline" size="sm" @click="perbarui">Coba lagi</Button>
        </div>

        <ul class="flex flex-col divide-y divide-border rounded-xl border border-border bg-surface" data-daftar-keranjang>
          <BarisKeranjang v-for="item in keranjang.items" :key="item.produk.id" :item="item" @hapus="hapus" />
        </ul>

        <div class="flex flex-wrap items-center justify-between gap-2">
          <Button variant="ghost" to="/katalog">Tambah produk lain</Button>
          <Button variant="ghost" data-kosongkan-keranjang @click="konfirmasiKosong = true">Kosongkan keranjang</Button>
        </div>
      </section>

      <RingkasanKeranjang
        class="lg:sticky lg:top-[calc(var(--tinggi-header,0px)+1.5rem)]"
        :items="dikirim"
        :pesan="pesan"
        :link="link"
        :jumlah-tidak-dijual="jumlahTidakDijual"
        @kosongkan="konfirmasiKosong = true"
      />

      <!-- HP: total + tombol kirim selalu terjangkau (meta.barAksiHp menyembunyikan tombol WA melayang). -->
      <div class="fixed inset-x-0 bottom-0 z-40 border-t border-border bg-surface px-4 pt-3 pb-[calc(0.75rem+env(safe-area-inset-bottom))] lg:hidden" data-bar-keranjang>
        <div class="mx-auto flex max-w-7xl items-center gap-3">
          <p class="flex min-w-0 flex-1 flex-col">
            <span class="text-xs text-muted">{{ dikirim.length }} produk · perkiraan</span>
            <strong class="truncate text-lg tabular-nums">{{ total.adaHarga ? formatRupiah(total.rupiah) : 'Via WhatsApp' }}</strong>
          </p>
          <WhatsAppButton class="shrink-0" :href="link" label="Pesan via WhatsApp" />
        </div>
      </div>
    </div>

    <ConfirmModal
      :open="konfirmasiKosong"
      title="Kosongkan keranjang?"
      :message="`Semua ${keranjang.jumlahProduk} produk akan dihapus dari keranjang.`"
      confirm-label="Kosongkan"
      danger
      @update:open="(buka) => !buka && (konfirmasiKosong = false)"
      @confirm="kosongkan"
    />
  </div>
</template>
