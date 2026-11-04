import { formatAngka, formatHargaPerSatuan, formatRupiah } from '@/utils/format'
import { hargaDisembunyikan } from '@/utils/whatsapp'

// Template pesan WhatsApp yang diatur admin di Konten › Template WhatsApp (KEPUTUSAN #23):
// - produk: tombol "Tanya via WhatsApp" satu produk (kartu, daftar, dan detail produk);
// - keranjang + baris: tombol "Pesan via WhatsApp" di /keranjang. {daftar_produk} berisi semua produk keranjang,
//   masing-masing ditulis memakai template baris.
// Kode {nama} diganti nilainya. Baris yang semua kodenya kosong (mis. {pelanggan} saat belum masuk) dibuang.
// Template kosong dari server (backend belum siap) → memakai bawaan di bawah.

export const TEMPLATE_BAWAAN = {
  produk: ['Halo Transhome, {pertanyaan}', '{nama_produk}', 'SKU: {sku}', '{link}'].join('\n'),
  keranjang: [
    'Halo Transhome, saya ingin memesan produk berikut:',
    '',
    '{daftar_produk}',
    '',
    'Jumlah: {jumlah_produk} produk',
    'Perkiraan total: {perkiraan_total}',
    'Atas nama: {pelanggan}',
    '',
    'Mohon dibantu cek stok dan total harganya. Terima kasih.',
  ].join('\n'),
  baris: ['{no}. {nama_produk}', '   SKU {sku} · {jumlah} {satuan} · {harga}', '   {catatan}'].join('\n'),
}

export const BATAS_TEMPLATE = { produk: 1000, keranjang: 1500, baris: 300 }

const PELANGGAN = { kode: 'pelanggan', arti: 'Nama dan kode member bila pelanggan sudah masuk; kosong bila belum.' }
const HARGA = { kode: 'harga', arti: 'Harga per satuan, atau "harga via WhatsApp" bila harganya disembunyikan.' }

// Kode yang boleh dipakai per template. `wajib`: template ditolak bila kode ini tidak ada.
export const KODE_TEMPLATE = {
  produk: [
    { kode: 'pertanyaan', arti: 'Kalimat tanya otomatis: tanya produk, tanya harga (harga disembunyikan), atau tanya ketersediaan (stok habis).' },
    { kode: 'nama_produk', arti: 'Nama produk.', wajib: true },
    { kode: 'sku', arti: 'Kode SKU produk.', wajib: true },
    { kode: 'link', arti: 'Tautan halaman produk di website.', wajib: true },
    { kode: 'brand', arti: 'Nama brand; kosong bila produk tanpa brand.' },
    HARGA,
    PELANGGAN,
  ],
  keranjang: [
    { kode: 'daftar_produk', arti: 'Semua produk di keranjang, masing-masing memakai format baris produk.', wajib: true },
    { kode: 'jumlah_produk', arti: 'Banyaknya produk berbeda di keranjang.' },
    { kode: 'perkiraan_total', arti: 'Total harga produk yang harganya tampil di website; kosong bila semua harga via WhatsApp.' },
    PELANGGAN,
  ],
  baris: [
    { kode: 'no', arti: 'Nomor urut: 1, 2, 3, …' },
    { kode: 'nama_produk', arti: 'Nama produk.', wajib: true },
    { kode: 'sku', arti: 'Kode SKU produk.', wajib: true },
    { kode: 'jumlah', arti: 'Jumlah yang dipesan.', wajib: true },
    { kode: 'satuan', arti: 'Satuan jual, mis. sak atau dus.' },
    HARGA,
    { kode: 'subtotal', arti: 'Harga × jumlah; kosong bila harganya via WhatsApp.' },
    { kode: 'brand', arti: 'Nama brand; kosong bila produk tanpa brand.' },
    { kode: 'link', arti: 'Tautan halaman produk (pesan jadi lebih panjang).' },
    { kode: 'catatan', arti: 'Catatan stok: stok habis, pre-order, atau sisa stok terbatas; kosong bila tersedia.' },
  ],
}

const CATATAN_STOK = { HABIS: 'Stok habis, mohon info ketersediaan', PRE_ORDER: 'Pre-order', SISA_STOK: 'Sisa stok terbatas' }
const POLA_KODE = /\{([^{}\n]*)\}/g

/** Ganti {kode} dengan nilainya; buang baris yang semua kodenya kosong. Kode tak dikenal dibiarkan apa adanya. */
export function isiTemplate(template, nilai) {
  const dikenal = (kode) => Object.hasOwn(nilai, kode)
  const kosong = (kode) => String(nilai[kode] ?? '').trim() === ''
  return String(template ?? '')
    .replace(/\r\n?/g, '\n')
    .split('\n')
    .filter((baris) => {
      const kode = [...baris.matchAll(POLA_KODE)].map((m) => m[1]).filter(dikenal)
      return !kode.length || !kode.every(kosong)
    })
    .map((baris) => baris.replace(POLA_KODE, (asli, kode) => (dikenal(kode) ? String(nilai[kode] ?? '') : asli)))
    .join('\n')
    .replace(/\n{3,}/g, '\n\n')
    .trim()
}

/** Pesan galat untuk template admin, atau '' bila sah. Dipakai halaman admin dan (mock) server. */
export function cekTemplate(jenis, teks) {
  const isi = String(teks ?? '')
  if (!isi.trim()) return 'Template wajib diisi'
  if (isi.length > BATAS_TEMPLATE[jenis]) return `Maksimal ${BATAS_TEMPLATE[jenis]} karakter`
  const daftar = KODE_TEMPLATE[jenis]
  const dipakai = [...isi.matchAll(POLA_KODE)].map((m) => m[1])
  const asing = [...new Set(dipakai.filter((kode) => !daftar.some((k) => k.kode === kode)))]
  if (asing.length) return `Kode tidak dikenal: ${asing.map((k) => `{${k}}`).join(', ')}. Pakai kode dari daftar.`
  const kurang = daftar.filter((k) => k.wajib && !dipakai.includes(k.kode))
  if (kurang.length) return `Wajib memuat ${kurang.map((k) => `{${k.kode}}`).join(', ')}`
  return ''
}

const pilih = (template, jenis) => (String(template ?? '').trim() ? template : TEMPLATE_BAWAAN[jenis])
const teksHarga = (produk) => (hargaDisembunyikan(produk) ? 'harga via WhatsApp' : formatHargaPerSatuan(produk.price_general, produk.unit_sale))

/** Isi {pertanyaan}: stok HABIS paling penting, lalu harga yang disembunyikan (sama dengan label tombolnya). */
export function pertanyaanProduk(produk) {
  if (produk?.stock_status === 'HABIS') return 'apakah produk ini akan tersedia lagi?'
  if (hargaDisembunyikan(produk)) return 'saya ingin menanyakan harga produk ini:'
  return 'saya ingin bertanya tentang produk ini:'
}

/** "Budi Santoso (member TF-000126)", "Budi Santoso", atau '' bila belum masuk. */
export function teksPelanggan(user) {
  const nama = user?.full_name?.trim()
  if (!nama) return ''
  return user.member_code ? `${nama} (member ${user.member_code})` : nama
}

/** Pesan "Tanya via WhatsApp" satu produk. */
export function pesanProduk(produk, { template, link, pelanggan = '' } = {}) {
  return isiTemplate(pilih(template, 'produk'), {
    pertanyaan: pertanyaanProduk(produk),
    nama_produk: produk.name,
    sku: produk.sku,
    link,
    brand: produk.brand?.name ?? '',
    harga: teksHarga(produk),
    pelanggan,
  })
}

/** Total dari produk yang harganya tampil; produk harga via WhatsApp dihitung terpisah. items: [{ produk, jumlah }]. */
export function hitungTotal(items) {
  let rupiah = 0
  let tanpaHarga = 0
  for (const { produk, jumlah } of items) {
    if (hargaDisembunyikan(produk)) tanpaHarga += 1
    else rupiah += Number(produk.price_general) * jumlah
  }
  return { rupiah, tanpaHarga, adaHarga: items.length > tanpaHarga }
}

function teksPerkiraanTotal({ rupiah, tanpaHarga, adaHarga }) {
  if (!adaHarga) return ''
  return tanpaHarga ? `${formatRupiah(rupiah)} (belum termasuk ${tanpaHarga} produk dengan harga via WhatsApp)` : formatRupiah(rupiah)
}

/** Pesan "Pesan via WhatsApp" dari keranjang. items: [{ produk, jumlah }] yang masih dijual; linkProduk(produk) → URL. */
export function pesanKeranjang(items, { template, templateBaris, linkProduk, pelanggan = '' } = {}) {
  const baris = pilih(templateBaris, 'baris')
  const daftar = items.map(({ produk, jumlah }, i) =>
    isiTemplate(baris, {
      no: String(i + 1),
      nama_produk: produk.name,
      sku: produk.sku,
      jumlah: formatAngka(jumlah),
      satuan: produk.unit_sale ?? '',
      harga: teksHarga(produk),
      subtotal: hargaDisembunyikan(produk) ? '' : formatRupiah(Number(produk.price_general) * jumlah),
      brand: produk.brand?.name ?? '',
      link: linkProduk(produk),
      catatan: CATATAN_STOK[produk.stock_status] ?? '',
    }),
  )
  return isiTemplate(pilih(template, 'keranjang'), {
    daftar_produk: daftar.join('\n'),
    jumlah_produk: String(items.length),
    perkiraan_total: teksPerkiraanTotal(hitungTotal(items)),
    pelanggan,
  })
}
