import { computed } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useSettingsStore } from '@/stores/settings'
import { tautanArtikel } from '@/utils/artikel'
import { buatLinkWhatsApp, pesanCariProduk, pesanJumlahBesar, pesanKonsultasiJasa, pesanTanyaKonten, pesanUmum } from '@/utils/whatsapp'
import { pesanKeranjang, pesanProduk, teksPelanggan } from '@/utils/templateWhatsApp'

/**
 * Link WhatsApp ke CS Transhome. Nomor dan template pesan dari pengaturan toko (API).
 * Setiap fungsi mengembalikan null selama nomor belum tersedia; tombol menampilkan state nonaktif.
 */
export function useWhatsApp() {
  const settings = useSettingsStore()
  const auth = useAuthStore()
  if (settings.status === 'idle') settings.load()

  const nomor = computed(() => settings.waNumber)
  const urlProduk = (produk) => new URL(`/produk/${produk.slug}`, window.location.origin).href

  /** "Tanya via WhatsApp" satu produk: template admin "Pesan per produk". */
  function productLink(produk) {
    if (!nomor.value || !produk) return null
    const pesan = pesanProduk(produk, { template: settings.templateWa.produk, link: urlProduk(produk), pelanggan: teksPelanggan(auth.user) })
    return buatLinkWhatsApp(nomor.value, pesan)
  }

  /** Isi pesan keranjang (juga untuk pratinjau di /keranjang). items: [{ produk, jumlah }] yang masih dijual. */
  function pesanDariKeranjang(items) {
    return pesanKeranjang(items, {
      template: settings.templateWa.keranjang,
      templateBaris: settings.templateWa.baris,
      linkProduk: urlProduk,
      pelanggan: teksPelanggan(auth.user),
    })
  }

  /** "Pesan via WhatsApp" di /keranjang: semua produk keranjang dalam satu pesan. */
  function keranjangLink(items) {
    if (!nomor.value || !items?.length) return null
    return buatLinkWhatsApp(nomor.value, pesanDariKeranjang(items))
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
    const url = new URL(tautanArtikel(artikel), window.location.origin).href
    return buatLinkWhatsApp(nomor.value, pesanTanyaKonten({ judul: artikel.title, url }))
  }

  /** Konsultasi jasa; `produk` diisi bila tombolnya di halaman produk (sebut SKU & nama produk). */
  function jasaLink(jasa, produk = null) {
    if (!nomor.value || !jasa) return null
    const url = new URL(`/jasa/${jasa.slug}`, window.location.origin).href
    const pesan = pesanKonsultasiJasa({
      jasa: jasa.name,
      url,
      produk: produk ? { sku: produk.sku, nama: produk.name } : null,
      nama: auth.user?.full_name,
      memberCode: auth.user?.member_code,
    })
    return buatLinkWhatsApp(nomor.value, pesan)
  }

  return {
    tersedia: computed(() => Boolean(nomor.value)),
    productLink, keranjangLink, pesanDariKeranjang, bulkLink, generalLink, searchLink, kontenLink, jasaLink,
  }
}
