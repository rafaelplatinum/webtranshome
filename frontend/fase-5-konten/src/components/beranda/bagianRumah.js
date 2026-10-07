// Enam bagian rumah di hero Beranda (teks dari sketsa "Landing3D"). Kategori tujuan "Lihat produk …"
// datang dari API (GET /home/house-parts) lewat `key`. `titik`/`tinta`: token warna nomor di daftar bagian.
export const BAGIAN_RUMAH = [
  {
    key: 'pondasi', nomor: 1, nama: 'Pondasi dan struktur', singkat: 'struktur', isi: 'Semen, besi beton, batu, pasir',
    deskripsi: 'Pondasi telapak dan plat lantai: semen, besi beton, batu, pasir, dan agregat.', titik: 'rumah-beton', tinta: 'foreground',
  },
  {
    key: 'lantai', nomor: 2, nama: 'Lantai dan teras', singkat: 'keramik', isi: 'Keramik, granit, perekat, nat',
    deskripsi: 'Penutup lantai dan teras: keramik, granit, perekat semen, dan nat.', titik: 'rumah-lantai', tinta: 'foreground',
  },
  {
    key: 'dinding', nomor: 3, nama: 'Dinding dan kolom', singkat: 'dinding dan cat', isi: 'Bata ringan, plester, cat tembok',
    deskripsi: 'Dinding, kolom, dan sekat: bata ringan, semen instan, plester, dan cat tembok.', titik: 'rumah-dinding', tinta: 'foreground',
  },
  {
    key: 'bukaan', nomor: 4, nama: 'Pintu dan jendela', singkat: 'pintu dan jendela', isi: 'Kusen, kaca, engsel, kunci',
    deskripsi: 'Bukaan rumah: pintu, kusen, kaca jendela, engsel, handle, dan kunci.', titik: 'rumah-kayu', tinta: 'surface',
  },
  {
    key: 'sanitasi', nomor: 5, nama: 'Sanitasi dan listrik', singkat: 'sanitasi', isi: 'Kloset, wastafel, pipa, lampu',
    deskripsi: 'Kamar mandi dan instalasi: bathtub, kloset, wastafel, pemanas air, pipa, panel listrik, dan lampu.', titik: 'rumah-air', tinta: 'foreground',
  },
  {
    key: 'atap', nomor: 6, nama: 'Atap', singkat: 'atap', isi: 'Genteng, bubungan, talang',
    deskripsi: 'Penutup atap: genteng, bubungan, talang, rangka baja ringan, dan cerobong.', titik: 'rumah-genteng', tinta: 'surface',
  },
]
