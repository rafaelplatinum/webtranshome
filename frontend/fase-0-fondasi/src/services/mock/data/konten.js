// Data contoh banner dan artikel. Kolom mengikuti docs/transhome_postgres.sql. TODO: ganti ke API.
// Gambar berupa SVG contoh berlabel (bukan foto asli). Tanggal dibuat relatif terhadap hari ini (WIB),
// supaya status banner (Tayang, Terjadwal, Berakhir) dan "banner hampir berakhir" di dashboard selalu bisa dicoba.

const tanggalRelatif = (hari) => new Date(Date.now() + 7 * 3600 * 1000 + hari * 86400 * 1000).toISOString().slice(0, 10)

// Teks di dalam SVG harus di-escape: "&" mentah membuat gambar rusak.
const xml = (teks) => String(teks).replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;')

// `judul` kosong: banner (judulnya sudah ditampilkan di atas gambar oleh slider).
function gambarContoh(judul, keterangan, latar, { lebar = 1600, tinggi = 640 } = {}) {
  const svg =
    `<svg xmlns="http://www.w3.org/2000/svg" width="${lebar}" height="${tinggi}" viewBox="0 0 ${lebar} ${tinggi}">` +
    `<rect width="${lebar}" height="${tinggi}" fill="${latar}"/>` +
    (judul ? `<text x="${lebar / 2}" y="${tinggi / 2 - 10}" text-anchor="middle" font-family="sans-serif" font-size="${Math.round(tinggi / 9)}" font-weight="700" fill="#1C2B3A">${xml(judul)}</text>` : '') +
    `<text x="${lebar / 2}" y="${tinggi / 2 + tinggi / 9}" text-anchor="middle" font-family="sans-serif" font-size="${Math.round(tinggi / 16)}" fill="#57534E">${xml(keterangan)}</text>` +
    `<text x="${lebar - 40}" y="${tinggi - 32}" text-anchor="end" font-family="sans-serif" font-size="${Math.round(tinggi / 26)}" fill="#57534E">CONTOH GAMBAR</text>` +
    `</svg>`
  return `data:image/svg+xml;charset=utf-8,${encodeURIComponent(svg)}`
}

const banner = (id, title, latar, link_url, position, mulai, selesai, sort_order, is_active = true) => ({
  id, title, image_url: gambarContoh('', 'Gambar banner promo', latar, position === 'HOME_SIDE' ? { lebar: 800, tinggi: 520 } : {}),
  link_url, position, start_at: `${mulai}T00:00:00`, end_at: selesai ? `${selesai}T23:59:59` : null, sort_order, is_active,
  created_at: '2026-09-01T09:00:00', updated_at: '2026-09-01T09:00:00',
})

// position: HOME_SLIDER (slider utama Beranda) atau HOME_SIDE (kartu di samping slider).
export const banners = [
  banner(1, 'Promo keramik lantai', '#E7D7C1', '/katalog/keramik-lantai', 'HOME_SLIDER', tanggalRelatif(-6), tanggalRelatif(3), 1),
  banner(2, 'Paket atap baja ringan', '#F3D9C4', '/katalog/atap', 'HOME_SLIDER', tanggalRelatif(-10), tanggalRelatif(30), 2),
  banner(3, 'Lampu LED hemat energi', '#FDF0C2', '/katalog/lampu', 'HOME_SIDE', tanggalRelatif(-3), tanggalRelatif(20), 1),
  banner(4, 'Pekan cat tembok', '#DCE3D2', '/katalog/cat-tembok', 'HOME_SLIDER', tanggalRelatif(20), tanggalRelatif(34), 3),
  banner(5, 'Promo genteng metal', '#E9DFCE', '/katalog/genteng', 'HOME_SLIDER', tanggalRelatif(-14), tanggalRelatif(-1), 4),
  banner(6, 'Promo sanitasi September', '#D4E2EA', '/katalog/sanitasi', 'HOME_SLIDER', '2026-09-01', '2026-09-30', 5),
  banner(7, 'Promo kran dapur', '#D6D3CD', '/katalog/kran', 'HOME_SLIDER', tanggalRelatif(-2), tanggalRelatif(12), 6, false),
]

// Isi artikel: paragraf dipisah baris kosong, "## " untuk subjudul, "- " untuk daftar (lihat utils/artikel.js).
const artikel = (id, type, title, slug, status, hari, latar, content, preview_token) => ({
  id, type, title, slug, status, content, preview_token,
  thumbnail_url: latar ? gambarContoh(title, type === 'ARTIKEL' ? 'Tips bangunan' : 'Promo & event', latar, { lebar: 800, tinggi: 450 }) : null,
  published_at: status === 'PUBLISHED' ? `${tanggalRelatif(hari)}T08:00:00` : null, author_id: 4,
  created_at: `${tanggalRelatif(hari - 1)}T15:00:00`, updated_at: `${tanggalRelatif(hari)}T08:00:00`,
})

export const articles = [
  artikel(1, 'PROMO', 'Promo keramik lantai bulan ini', 'promo-keramik-lantai', 'PUBLISHED', -6, '#E7D7C1',
    'Contoh isi promo. Harga dan syarat promo diisi Admin Konten.\n\n## Berlaku untuk\n- Keramik lantai pilihan\n- Pembelian langsung di toko\n\nTanyakan stok dan harga promo ke tim kami lewat WhatsApp sebelum datang ke toko.',
    'a41f0c9e7b2d4c18'),
  artikel(2, 'EVENT', 'Demo pemasangan bata ringan di toko', 'demo-pemasangan-bata-ringan', 'PUBLISHED', -12, '#DCE3D2',
    'Contoh isi event. Tanggal, jam, dan tempat diisi Admin Konten.\n\n## Yang akan dibahas\n- Menyiapkan perekat bata ringan\n- Memasang dan meratakan susunan bata\n- Menghitung kebutuhan untuk satu dinding\n\nAcara terbuka untuk umum. Tanyakan jadwalnya lewat WhatsApp.',
    'b53e1dae8c3f5d29'),
  artikel(3, 'PROMO', 'Paket atap baja ringan', 'paket-atap-baja-ringan', 'PUBLISHED', -2, '#F3D9C4',
    'Contoh isi promo paket. Isi paket dan harga diisi Admin Konten.\n\nKirim ukuran atap rumahmu lewat WhatsApp, tim kami bantu hitung kebutuhan rangka dan genteng.',
    'c62f2ebf9d406e3a'),
  artikel(4, 'ARTIKEL', 'Tips memilih cat tembok untuk ruang luar', 'tips-cat-tembok-ruang-luar', 'PUBLISHED', -20, '#FDF0C2',
    'Dinding luar terkena panas dan hujan setiap hari, jadi pilih cat yang memang dibuat untuk eksterior.\n\n## Yang perlu dicek\n- Label eksterior atau tahan cuaca pada kemasan\n- Warna: warna terang membantu ruangan tidak cepat panas\n- Jumlah: hitung luas dinding sebelum membeli\n\nBelum yakin? Kirim foto dindingmu ke tim kami lewat WhatsApp.',
    'd7304fc0ae517f4b'),
  artikel(5, 'ARTIKEL', 'Menghitung kebutuhan semen untuk pondasi', 'menghitung-kebutuhan-semen-pondasi', 'PUBLISHED', -30, '#E9DFCE',
    'Kebutuhan semen bergantung pada ukuran pondasi dan campuran yang dipakai tukang.\n\n## Siapkan dulu\n- Panjang, lebar, dan kedalaman pondasi\n- Jenis pondasi yang dipakai\n\nKirim ukurannya ke tim kami lewat WhatsApp untuk dibantu menghitung jumlah sak semen, pasir, dan batu.',
    'e8415ad1bf62805c'),
  artikel(6, 'ARTIKEL', 'Cara memilih keramik lantai untuk teras', 'cara-memilih-keramik-teras', 'DRAFT', -1, '#D4E2EA',
    'Draf artikel. Teras sering basah, jadi pilih keramik yang permukaannya tidak licin.\n\n## Yang perlu dicek\n- Permukaan doff atau bertekstur\n- Ukuran yang sesuai luas teras',
    'f9526be2c073916d'),
  artikel(7, 'ARTIKEL', 'Merawat kran dan wastafel agar awet', 'merawat-kran-wastafel', 'PUBLISHED', -40, null,
    'Kran dan wastafel yang rutin dibersihkan lebih jarang bermasalah.\n\n## Perawatan sederhana\n- Lap bagian kran setelah dipakai\n- Bersihkan saringan kran bila air mulai kecil\n- Hindari pembersih yang terlalu keras\n\nButuh suku cadang? Tanyakan ketersediaannya lewat WhatsApp.',
    '0a637cf3d184a27e'),
]
