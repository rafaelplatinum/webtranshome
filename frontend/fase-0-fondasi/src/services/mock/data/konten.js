// Data konten dari snapshot database backup Transhome.
// URL gambar dari backup diganti dengan SVG dummy lokal sampai aset aslinya tersedia.
import backup from './backup.json'

const xml = (teks) => String(teks ?? '').replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;')

function gambarContoh(judul, keterangan, lebar, tinggi, latar) {
  const svg = `<svg xmlns="http://www.w3.org/2000/svg" width="${lebar}" height="${tinggi}" viewBox="0 0 ${lebar} ${tinggi}"><rect width="${lebar}" height="${tinggi}" fill="${latar}"/><text x="50%" y="47%" text-anchor="middle" font-family="sans-serif" font-size="${Math.round(tinggi / 10)}" font-weight="700" fill="#1C2B3A">${xml(judul)}</text><text x="50%" y="59%" text-anchor="middle" font-family="sans-serif" font-size="${Math.round(tinggi / 20)}" fill="#57534E">${xml(keterangan)}</text><text x="96%" y="94%" text-anchor="end" font-family="sans-serif" font-size="${Math.round(tinggi / 30)}" fill="#57534E">GAMBAR CONTOH</text></svg>`
  return `data:image/svg+xml;charset=utf-8,${encodeURIComponent(svg)}`
}

const warna = ['#E7D7C1', '#DCE3D2', '#D4E2EA', '#FDF0C2', '#D6D3CD']

function bannerHeroKosong(indeks) {
  const palet = [
    ['#1C2B3A', '#C2410C', '#FACC15'],
    ['#263B4E', '#D97706', '#FAF8F3'],
    ['#243447', '#B9380B', '#F4C430'],
  ][indeks % 3]
  const svg = `<svg xmlns="http://www.w3.org/2000/svg" width="1920" height="800" viewBox="0 0 1920 800"><defs><linearGradient id="latar" x1="0" y1="0" x2="1" y2="1"><stop stop-color="${palet[0]}"/><stop offset="1" stop-color="${palet[1]}"/></linearGradient><linearGradient id="cahaya" x1="0" y1="0" x2="1" y2="1"><stop stop-color="${palet[2]}" stop-opacity=".7"/><stop offset="1" stop-color="${palet[2]}" stop-opacity="0"/></linearGradient></defs><rect width="1920" height="800" fill="url(#latar)"/><circle cx="1490" cy="170" r="390" fill="url(#cahaya)"/><path d="M0 650 520 330l475 470H0Z" fill="#FAF8F3" fill-opacity=".08"/><path d="m960 800 430-480 530 390v90Z" fill="${palet[2]}" fill-opacity=".18"/><path d="M1050 800 1500 0h180L1230 800Z" fill="#FAF8F3" fill-opacity=".07"/><g fill="none" stroke="#FAF8F3" stroke-opacity=".2" stroke-width="3"><path d="M1270 800a300 300 0 0 1 600 0"/><path d="M1350 800a220 220 0 0 1 440 0"/><path d="M1430 800a140 140 0 0 1 280 0"/></g></svg>`
  return `data:image/svg+xml;charset=utf-8,${encodeURIComponent(svg)}`
}

export const banners = backup.banners.map((item, i) => ({
  ...item,
  image_url: item.position === 'HOME_HERO' ? bannerHeroKosong(i) : gambarContoh(item.title, 'Banner Transhome', 1600, 640, warna[i % warna.length]),
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
