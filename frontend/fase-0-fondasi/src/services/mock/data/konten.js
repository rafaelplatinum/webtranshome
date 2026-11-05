// Data konten dari snapshot database backup Transhome.
// URL gambar dari backup diganti dengan SVG dummy lokal sampai aset aslinya tersedia.
import backup from './backup.json'

const xml = (teks) => String(teks ?? '').replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;')

function gambarContoh(judul, keterangan, lebar, tinggi, latar) {
  const svg = `<svg xmlns="http://www.w3.org/2000/svg" width="${lebar}" height="${tinggi}" viewBox="0 0 ${lebar} ${tinggi}"><rect width="${lebar}" height="${tinggi}" fill="${latar}"/><text x="50%" y="47%" text-anchor="middle" font-family="sans-serif" font-size="${Math.round(tinggi / 10)}" font-weight="700" fill="#1C2B3A">${xml(judul)}</text><text x="50%" y="59%" text-anchor="middle" font-family="sans-serif" font-size="${Math.round(tinggi / 20)}" fill="#57534E">${xml(keterangan)}</text><text x="96%" y="94%" text-anchor="end" font-family="sans-serif" font-size="${Math.round(tinggi / 30)}" fill="#57534E">GAMBAR CONTOH</text></svg>`
  return `data:image/svg+xml;charset=utf-8,${encodeURIComponent(svg)}`
}

const warna = ['#E7D7C1', '#DCE3D2', '#D4E2EA', '#FDF0C2', '#D6D3CD']

export const banners = backup.banners.map((item, i) => ({
  ...item,
  image_url: gambarContoh(item.title, 'Banner Transhome', 1600, 640, warna[i % warna.length]),
}))

export const articles = backup.articles.map((item, i) => ({
  ...item,
  thumbnail_url: item.thumbnail_url
    ? gambarContoh(item.title, item.type === 'ARTIKEL' ? 'Panduan isi rumah' : 'Promo & event', 800, 450, warna[i % warna.length])
    : null,
}))

export const services = backup.services.map((item, i) => ({
  ...item,
  image_url: gambarContoh(item.name, 'Jasa Transhome', 800, 450, warna[i % warna.length]),
  category_ids: backup.service_categories.filter((link) => link.service_id === item.id).map((link) => link.category_id),
}))
