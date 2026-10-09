import { computed } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useSettingsStore } from '@/stores/settings'
import { buatLinkWhatsApp, pesanCariProduk, pesanJumlahBesar, pesanTanyaKonten, pesanTanyaProduk, pesanUmum } from '@/utils/whatsapp'

/**
 * Link WhatsApp ke CS Transhome. Nomor dari pengaturan toko (API).
 * Setiap fungsi mengembalikan null selama nomor belum tersedia; tombol menampilkan state nonaktif.
 */
export function useWhatsApp() {
  const settings = useSettingsStore()
  const auth = useAuthStore()
  if (settings.status === 'idle') settings.load()

  const nomor = computed(() => settings.waNumber)
  const urlProduk = (produk) => new URL(`/produk/${produk.slug}`, window.location.origin).href

  function productLink(produk) {
    if (!nomor.value || !produk) return null
    const pesan = pesanTanyaProduk({
      sku: produk.sku,
      nama: produk.name,
      url: urlProduk(produk),
      habis: produk.stock_status === 'HABIS',
    })
    return buatLinkWhatsApp(nomor.value, pesan)
  }

  function bulkLink(produk) {
    if (!nomor.value || !produk) return null
    return buatLinkWhatsApp(nomor.value, pesanJumlahBesar({ sku: produk.sku, nama: produk.name, url: urlProduk(produk) }))
  }

  function generalLink() {
    if (!nomor.value) return null
    return buatLinkWhatsApp(nomor.value, pesanUmum({ nama: auth.user?.full_name, memberCode: auth.user?.member_code }))
  }

  /** EmptyState katalog/pencarian: tanyakan barang yang belum ketemu. Tanpa topik → pesan umum. */
  function searchLink(topik) {
    if (!topik) return generalLink()
    if (!nomor.value) return null
    return buatLinkWhatsApp(nomor.value, pesanCariProduk(topik))
  }

  /** Promo, event, atau artikel (Fase 5): pesan berisi judul dan tautan halamannya. */
  function kontenLink(artikel) {
    if (!nomor.value || !artikel) return null
    const url = new URL(`/artikel/${artikel.slug}`, window.location.origin).href
    return buatLinkWhatsApp(nomor.value, pesanTanyaKonten({ judul: artikel.title, url }))
  }

  return { tersedia: computed(() => Boolean(nomor.value)), productLink, bulkLink, generalLink, searchLink, kontenLink }
}
