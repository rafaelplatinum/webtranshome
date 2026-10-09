import { normalisasiNomorHp } from './format'

/** Link wa.me dengan pesan yang sudah di-encode. Nomor diambil dari pengaturan toko, jangan hardcode. */
export function buatLinkWhatsApp(nomorWa, pesan) {
  return `https://wa.me/${normalisasiNomorHp(nomorWa)}?text=${encodeURIComponent(pesan)}`
}

/** "Tanya via WhatsApp" di produk; produk HABIS memakai pesan "Tanya ketersediaan". */
export function pesanTanyaProduk({ sku, nama, url, habis = false }) {
  const pembuka = habis
    ? 'Halo Transhome, apakah produk ini akan tersedia lagi?'
    : 'Halo Transhome, saya ingin bertanya tentang produk ini:'
  return [pembuka, nama, `SKU: ${sku}`, url].join('\n')
}

/** "Chat tim proyek" untuk kebutuhan jumlah besar. */
export function pesanJumlahBesar({ sku, nama, url }) {
  return [
    'Halo tim proyek Transhome, saya butuh produk ini dalam jumlah besar:',
    nama,
    `SKU: ${sku}`,
    url,
    'Perkiraan kebutuhan: ',
  ].join('\n')
}

/** Dari katalog/pencarian yang kosong: barang yang dicari belum ketemu di website. */
export function pesanCariProduk(topik) {
  return `Halo Transhome, saya mencari ${topik} tapi belum ketemu di website. Apakah ada?`
}

/** Dari halaman promo, event, atau artikel: sebut judul dan tautannya. */
export function pesanTanyaKonten({ judul, url }) {
  return ['Halo Transhome, saya ingin bertanya tentang ini:', judul, url].join('\n')
}

/** Pesan umum ke CS; bila sudah masuk, diawali nama dan member_code (FR-S). */
export function pesanUmum({ nama, memberCode } = {}) {
  if (nama && memberCode) return `Halo Transhome, saya ${nama} (member ${memberCode}).`
  if (nama) return `Halo Transhome, saya ${nama}.`
  return 'Halo Transhome, saya ingin bertanya.'
}
