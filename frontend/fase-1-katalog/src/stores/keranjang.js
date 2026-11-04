import { computed, ref, watch } from 'vue'
import { defineStore } from 'pinia'

// Keranjang belanja (KEPUTUSAN #23): daftar produk + jumlah yang dikirim sekaligus lewat WhatsApp. Tanpa checkout,
// pembayaran, atau ongkir. Disimpan di browser (localStorage) supaya bisa dipakai tanpa login dan tetap ada setelah
// halaman ditutup; tab lain ikut berubah. Harga & stok yang disimpan hanya salinan, diperbarui saat /keranjang dibuka.
const KUNCI = 'transhome.keranjang'
const VERSI = 1
export const MAKS_PRODUK_KERANJANG = 30
export const MAKS_JUMLAH = 9999

/** Jumlah minimal per produk: minimal pembelian dari admin (min_order), paling sedikit 1. */
export const jumlahMinimal = (produk) => Math.max(1, Math.floor(Number(produk?.min_order) || 1))
const rapikanJumlah = (jumlah, produk) => Math.min(MAKS_JUMLAH, Math.max(jumlahMinimal(produk), Math.floor(Number(jumlah) || 0)))

// Data produk yang cukup untuk menampilkan baris keranjang dan menyusun pesan WhatsApp.
function ringkasProduk(p) {
  return {
    id: p.id, slug: p.slug, sku: p.sku, name: p.name, unit_sale: p.unit_sale ?? '', min_order: p.min_order ?? 1,
    price_general: p.price_general ?? null, show_price: p.show_price !== false, stock_status: p.stock_status ?? 'TERSEDIA',
    primary_image_url: p.primary_image_url ?? null, brand: p.brand ? { name: p.brand.name, slug: p.brand.slug } : null,
  }
}

function baca() {
  try {
    const simpanan = JSON.parse(localStorage.getItem(KUNCI) ?? 'null')
    if (simpanan?.versi !== VERSI || !Array.isArray(simpanan.items)) return []
    return simpanan.items
      .filter((i) => i?.produk?.id != null && i.produk.slug && i.produk.sku)
      .slice(0, MAKS_PRODUK_KERANJANG)
      .map((i) => ({ produk: ringkasProduk(i.produk), jumlah: rapikanJumlah(i.jumlah, i.produk), tidakDijual: Boolean(i.tidakDijual) }))
  } catch {
    return []
  }
}

export const useKeranjangStore = defineStore('keranjang', () => {
  // [{ produk, jumlah, tidakDijual }]; tidakDijual = produk sudah nonaktif/tidak ditemukan saat diperbarui.
  const items = ref(baca())

  const jumlahProduk = computed(() => items.value.length)
  const penuh = computed(() => items.value.length >= MAKS_PRODUK_KERANJANG)
  const dijual = computed(() => items.value.filter((i) => !i.tidakDijual))

  const cari = (id) => items.value.find((i) => i.produk.id === id) ?? null

  /** Tambah produk (atau tambah jumlahnya bila sudah ada). → { hasil: 'baru' | 'tambah' | 'penuh' | 'habis', jumlah } */
  function tambah(produk, jumlah = jumlahMinimal(produk)) {
    if (produk.stock_status === 'HABIS') return { hasil: 'habis', jumlah: 0 }
    const ada = cari(produk.id)
    if (ada) {
      ada.produk = ringkasProduk(produk)
      ada.tidakDijual = false
      ada.jumlah = rapikanJumlah(ada.jumlah + Number(jumlah), produk)
      return { hasil: 'tambah', jumlah: ada.jumlah }
    }
    if (penuh.value) return { hasil: 'penuh', jumlah: 0 }
    const baru = { produk: ringkasProduk(produk), jumlah: rapikanJumlah(jumlah, produk), tidakDijual: false }
    items.value.push(baru)
    return { hasil: 'baru', jumlah: baru.jumlah }
  }

  function ubahJumlah(id, jumlah) {
    const ada = cari(id)
    if (ada) ada.jumlah = rapikanJumlah(jumlah, ada.produk)
  }

  /** Hapus satu produk; hasilnya bisa dikembalikan lewat `kembalikan` (tombol Urungkan). */
  function hapus(id) {
    const posisi = items.value.findIndex((i) => i.produk.id === id)
    if (posisi < 0) return null
    const [item] = items.value.splice(posisi, 1)
    return { item, posisi }
  }

  function kembalikan(terhapus) {
    if (!terhapus || cari(terhapus.item.produk.id) || penuh.value) return
    items.value.splice(Math.min(terhapus.posisi, items.value.length), 0, terhapus.item)
  }

  function kosongkan() {
    items.value = []
  }

  /** Data terbaru dari GET /products?ids=…; produk yang tidak kembali dianggap sudah tidak dijual. */
  function perbarui(daftar) {
    const terbaru = new Map(daftar.map((p) => [p.id, p]))
    for (const item of items.value) {
      const p = terbaru.get(item.produk.id)
      item.tidakDijual = !p
      if (p) {
        item.produk = ringkasProduk(p)
        item.jumlah = rapikanJumlah(item.jumlah, item.produk)
      }
    }
  }

  // Simpan setiap perubahan. localStorage diblokir/penuh: keranjang hanya bertahan selama tab terbuka.
  let terakhir = null
  watch(
    items,
    () => {
      const teks = JSON.stringify({ versi: VERSI, items: items.value })
      if (teks === terakhir) return
      terakhir = teks
      try {
        localStorage.setItem(KUNCI, teks)
      } catch {
        // diabaikan
      }
    },
    { deep: true },
  )
  window.addEventListener('storage', (event) => {
    if (event.key !== KUNCI && event.key !== null) return
    terakhir = event.newValue
    items.value = baca()
  })

  return { items, jumlahProduk, penuh, dijual, cari, tambah, ubahJumlah, hapus, kembalikan, kosongkan, perbarui }
})
